package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() { addHandler(hamstring{},
           "Usage:  hamstring target # \n\n Try to hamstring a mob to generate a large amount of threat",
           permissions.Fighter, //TODO: Add Barbarian here
           "hamstring", "ham")
}

type hamstring cmd

func (hamstring) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Hamstring what exactly?")
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

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum,true)
	if whatMob != nil {

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == nil && s.actor.Class != 0 {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Equipment.Main.ItemType != 4 && (s.actor.Placement != whatMob.Placement) {
			s.msg.Actor.SendBad("You are too far away to hamstring them.")
			return
		}

		//skillLevel := config.WeaponLevel(s.actor.Skills[s.actor.Equipment.Main.ItemType].Value)

		//TODO: Parry/Miss/Resist being circled?
		whatMob.AddThreatDamage(whatMob.Stam.Max/2, s.actor.Name)
		s.actor.SetTimer("combat", config.CombatCooldown)
		s.msg.Actor.SendInfo("You hamstring " + whatMob.Name)
		s.msg.Observers.SendInfo(s.actor.Name + " hamstrings " + whatMob.Name )
		return

	}

	s.msg.Actor.SendInfo("Hamstring what?")
	s.ok = true
}
