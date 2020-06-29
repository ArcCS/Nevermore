package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(equip{},
           "Usage:  equip item # \n\n Try to equip an item from your inventory",
           permissions.Player,
           "equip", "wield", "wear")
}

type equip cmd

func (equip) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendBad("What did you want to equip?")
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

	what := s.actor.Inventory.Search(name, nameNum)
	if what != nil {
		s.actor.Inventory.Lock()
		s.actor.Equipment.Equip(what)
		s.actor.Inventory.Remove(what)
		s.actor.Inventory.Unlock()
		s.msg.Actor.SendGood("You put on " + what.Name)
		s.msg.Observer.SendInfo(s.actor.Name + " equips " + what.Name)
		s.ok = true
		return
	}
	s.msg.Actor.SendInfo("What did you want to equip?")
	s.ok = true
}
