package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/jinzhu/copier"
	"log"
	"strconv"
	"strings"
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

	targetStr := ""
	targetNum := 1

	for _, word := range s.words {
		if val, err := strconv.Atoi(word); err == nil {
			targetNum = val
		} else {
			targetStr += " " + word
		}
	}

	targetStr = strings.Trim(targetStr, " ")

	if len(s.where.StoreInventory.Contents) > 0 {
		purchaseItem := s.where.StoreInventory.Search(targetStr, targetNum)
		if purchaseItem != nil {
			if s.actor.Gold.Value >= purchaseItem.StorePrice {
				if (s.actor.GetCurrentWeight() + purchaseItem.GetWeight()) <= s.actor.MaxWeight() {
					s.actor.RunHook("act")
					s.actor.Gold.Subtract(purchaseItem.StorePrice)
					if purchaseItem.Flags["infinite"] {
						newItem := objects.Item{}
						if err := copier.CopyWithOption(&newItem, objects.Items[purchaseItem.ItemId], copier.Option{DeepCopy: true}); err != nil {
							log.Println("Error copying item: ", err)
						}
						s.actor.Inventory.Add(&newItem)
					} else {
						if err := s.where.StoreInventory.Remove(purchaseItem); err != nil {
							log.Println("Error removing item from store: ", err)
						}
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
