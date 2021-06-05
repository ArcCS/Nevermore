package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(buy{},
	"",
	permissions.Player,
	"$BUY")
}

type buy cmd

func (buy) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Buy what?")
		return
	}


	targetStr := s.words[0]
	targetNum := 1

	if len(s.words) > 1 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		}
	}

	if len(s.where.StoreInventory.Contents) > 0 {
		purchaseItem := s.where.StoreInventory.Search(targetStr, targetNum)
		if purchaseItem != nil {
			if s.actor.Gold.Value > purchaseItem.StorePrice {
				if (s.actor.Inventory.TotalWeight + purchaseItem.GetWeight()) <= s.actor.MaxWeight() {
					s.actor.RunHook("act")
					s.actor.Inventory.Lock()
					s.where.StoreInventory.Lock()
					s.actor.Gold.Subtract(purchaseItem.StorePrice)
					purchased := s.where.BuyStoreItem(purchaseItem)
					s.actor.Inventory.Add(purchased)
					s.where.StoreInventory.Unlock()
					s.actor.Inventory.Unlock()
					s.msg.Actor.SendGood("You purchase ", purchaseItem.Name, ".")
					s.msg.Observers.SendInfo(s.actor.Name, " purchases ", purchaseItem.Name, ".")
					return
				} else {
					s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
					return
				}
			}else {
				s.msg.Actor.SendInfo ("You don't have enough gold to purchase that item.")
				return
			}
		}else{
			s.msg.Actor.SendInfo ("You don't see that item to purchase.")
			return
		}
	}else{
		s.msg.Actor.SendBad("There's nothing to purchase at this store.")
		return
	}
}
