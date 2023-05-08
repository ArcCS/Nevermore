package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/jinzhu/copier"
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
			if s.actor.Gold.Value >= purchaseItem.StorePrice {
				if (s.actor.GetCurrentWeight() + purchaseItem.GetWeight()) <= s.actor.MaxWeight() {
					s.actor.RunHook("act")
					s.actor.Gold.Subtract(purchaseItem.StorePrice)
					if purchaseItem.Flags["infinite"] {
						newItem := objects.Item{}
						copier.CopyWithOption(&newItem, objects.Items[purchaseItem.ItemId], copier.Option{DeepCopy: true})
						s.actor.Inventory.Add(&newItem)
					} else {
						s.where.StoreInventory.Remove(purchaseItem)
						s.actor.Inventory.Add(purchaseItem)
					}
					s.msg.Actor.SendGood("You purchase ", purchaseItem.Name, ".")
					s.msg.Observers.SendInfo(s.actor.Name, " purchases ", purchaseItem.Name, ".")
					return
				} else {
					s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
					return
				}
			} else {
				s.msg.Actor.SendInfo("You don't have enough gold to purchase that item.")
				return
			}
		} else {
			s.msg.Actor.SendInfo("You don't see that item to purchase.")
			return
		}
	} else {
		s.msg.Actor.SendBad("There's nothing to purchase at this store.")
		return
	}
}
