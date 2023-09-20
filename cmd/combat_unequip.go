package cmd

import (
	"github.com/ArcCS/Nevermore/config"
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

	// Check some timers
	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	name := s.input[0]
	s.actor.RunHook("combat")

	if s.actor.CheckFlag("singing") {
		if s.actor.Equipment.FindLocation(name) != "main" {
			s.msg.Actor.SendBad("You may only remove your main hand weapon while performing.")
			s.ok = true
			return
		}
	}

	_, what := s.actor.Equipment.Unequip(name)
	if what != nil {
		s.actor.Inventory.Add(what)
		s.msg.Actor.SendGood("You unequip " + what.Name)
		s.msg.Observer.SendInfo(s.actor.Name + " unequips " + what.Name)
		s.actor.SetTimer("combat", config.UnequipCooldown)
		s.ok = true
		return
	}
	s.msg.Actor.SendInfo("What did you want to unequip?")
	s.ok = true
}
