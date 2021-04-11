package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(circle{},
		"Usage:  circle target # \n\n Try to circle a mob to apply a short duration stun and generate threat",
		permissions.Fighter&permissions.Barbarian,
		"circle", "cir")
}

type circle cmd

func (circle) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Circle what exactly?")
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("combat_circle")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
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
	whatMob = s.where.Mobs.Search(name, nameNum, true)
	if whatMob != nil {

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == nil {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		// Shortcut Missile weapon:
		if s.actor.Equipment.Main.ItemType == 4 {
			s.msg.Actor.SendBad("You cannot circle with a ranged weapon.")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to circle them.")
			return
		}

		//skillLevel := config.WeaponLevel(s.actor.Skills[s.actor.Equipment.Main.ItemType].Value)

		//TODO: Parry/Miss/Resist being circled?
		whatMob.MobStunned = config.CircleStuns
		whatMob.AddThreatDamage(whatMob.Stam.Max/10, s.actor.Name)
		s.actor.SetTimer("combat_circle", config.CircleTimer)
		s.actor.SetTimer("combat", config.CombatCooldown)
		s.msg.Actor.SendInfo("You circled " + whatMob.Name)
		s.msg.Observers.SendInfo(s.actor.Name + " circles " + whatMob.Name)
		return

	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}
