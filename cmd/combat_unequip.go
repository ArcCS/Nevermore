package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(unequip{},
		"Usage:  unequip item # \n\n Try to unequip something you're wearing",
		permissions.Player,
		"unequip", "remove", "rem")
}

type unequip cmd

func (unequip) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendBad("What did you want to equip?")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	name := s.input[0]
	/*nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	*/

	s.actor.RunHook("combat")
	_, what := s.actor.Equipment.Unequip(name)
	if what != nil {
		s.actor.Inventory.Lock()
		s.actor.Inventory.Add(what)
		s.actor.Inventory.Unlock()
		if what.ItemType == 16 && s.actor.CheckFlag("singing") {
			s.actor.RemoveEffect("sing")
			s.msg.Observers.SendInfo(s.actor.Name + " stops singing.")
		}
		s.msg.Actor.SendGood("You unequip " + what.Name)
		s.msg.Observer.SendInfo(s.actor.Name + " unequips " + what.Name)
		s.ok = true
		return
	}
	s.msg.Actor.SendInfo("What did you want to unequip?")
	s.ok = true
}
