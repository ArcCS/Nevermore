package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(removeitem{},
	"Usage: removeitem name #  Remove an item from the store inventory.",
	permissions.Player,
	"removeitem")
}

type removeitem cmd

func (removeitem) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Remove what item?")
		return
	}

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	whatItem := s.where.StoreInventory.Search(name, nameNum)
	if whatItem != nil {
		s.where.StoreInventory.Lock()
		s.actor.Inventory.Lock()
		s.where.StoreInventory.Remove(whatItem)
		s.actor.Inventory.Add(whatItem)
		s.where.StoreInventory.Unlock()
		s.actor.Inventory.Unlock()
		s.where.Save()
		s.msg.Actor.SendGood("You remove " + whatItem.Name + " from the store front.")
	}else{
		s.msg.Actor.SendBad("There's no matching item.")
	}
}
