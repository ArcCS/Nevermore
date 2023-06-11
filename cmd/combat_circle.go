package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(circle{},
		"Usage:  circle target # \n\n Try to circle a mob to apply a short duration stun and generate threat",
		permissions.Fighter|permissions.Barbarian,
		"circle", "cir")
}

type circle cmd

func (circle) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Circle what exactly?")
		return
	}

	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything!")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
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
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	if whatMob != nil {

		s.actor.RunHook("combat")
		s.actor.Victim = whatMob
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

		// Check for a miss
		if utils.Roll(100, 1, 0) <= DetermineMissChance(s, whatMob.Level-s.actor.Tier) {
			s.msg.Actor.SendBad("You missed!!")
			s.actor.SetTimer("combat_circle", config.CircleTimer)
			s.actor.SetTimer("combat", config.CombatCooldown)
			whatMob.AddThreatDamage(1, s.actor)
			whatMob.CurrentTarget = s.actor.Name
			s.msg.Observers.SendBad(s.actor.Name + " fails to circle " + whatMob.Name)
			return
		}

		whatMob.Stun(config.CircleStuns)
		whatMob.AddThreatDamage(whatMob.Stam.Max/2, s.actor)
		whatMob.CurrentTarget = s.actor.Name
		s.actor.SetTimer("combat_circle", config.CircleTimer)
		s.actor.SetTimer("combat", config.CombatCooldown)
		s.msg.Actor.SendInfo("You circled " + whatMob.Name)
		s.msg.Observers.SendInfo(s.actor.Name + " circles " + whatMob.Name)
		return

	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}
