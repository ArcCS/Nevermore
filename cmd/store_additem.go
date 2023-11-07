package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
)

func init() {
	addHandler(additem{},
		"Usage: additem name # price, adds an item to a store if you own it at the given price",
		permissions.Player,
		"additem")
}

type additem cmd

func (additem) process(s *state) {
	if s.where.StoreOwner == s.actor.Name || s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		if len(s.words) < 2 {
			s.msg.Actor.SendBad("There aren't enough arguments to process.")
			return
		}

		targetStr := s.words[0]
		targetNum := 1
		priceStr := 0

		if len(s.words) == 3 {
			if val, err := strconv.Atoi(s.words[1]); err == nil {
				targetNum = val
			}
			if val2, err2 := strconv.Atoi(s.words[2]); err2 == nil {
				priceStr = val2
			}
		} else {
			if val2, err2 := strconv.Atoi(s.words[1]); err2 == nil {
				priceStr = val2
			}
		}

		whatItem := s.actor.Inventory.Search(targetStr, targetNum)
		if whatItem != nil {
			if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				if err := s.actor.Inventory.Remove(whatItem); err != nil {
					s.msg.Actor.SendBad("You can't remove that item from your inventory.")
					log.Println("Error removing item from inventory: ", err)
					return
				}
				s.where.AddStoreItem(whatItem, priceStr, true)
				s.where.Save()
				s.msg.Actor.SendGood("You add " + whatItem.Name + " to the store.")
			} else {
				if err := s.actor.Inventory.Remove(whatItem); err != nil {
					s.msg.Actor.SendBad("You can't remove that item from your inventory.")
					log.Println("Error removing item from inventory: ", err)
					return
				}
				s.where.AddStoreItem(whatItem, priceStr, false)
				s.where.Save()
				s.msg.Actor.SendGood("You add " + whatItem.Name + " to the store.")
			}
		} else {
			s.msg.Actor.SendBad("That's not an item in your inventory.")
		}
	} else {
		s.msg.Actor.SendBad("This isn't your store to modify.")
	}
}
