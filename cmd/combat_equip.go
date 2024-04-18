package cmd

import (
	"github.com/ArcCS/Nevermore/config"
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
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	what := s.actor.Inventory.Search(name, nameNum)
	if what != nil {
		s.actor.RunHook("combat")
		if ok, msg := s.actor.CanEquip(what); !ok {
			s.msg.Actor.SendBad(msg)
			s.ok = true
			return
		}

		if what.MaxUses <= 0 {
			s.msg.Actor.SendInfo("The " + what.DisplayName() + " is broken")
			return
		}

		if s.actor.Equipment.Equip(what, s.actor.Class) {
			s.msg.Actor.SendGood("You equip " + what.DisplayName())
			s.msg.Observers.SendInfo(s.actor.Name + " equips " + what.DisplayName())
			if err := s.actor.Inventory.Remove(what); err != nil {
				s.msg.Actor.SendBad("Failure to equip item.")
				return
			}
			s.actor.SetTimer("combat", config.CombatCooldown)
		} else {
			s.msg.Actor.SendBad("You cannot equip that.")
		}

		s.ok = true
		return
	}
	s.msg.Actor.SendInfo("What did you want to equip?")
	s.ok = true
}
