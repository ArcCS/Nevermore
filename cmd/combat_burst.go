package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(burst{},
		"Usage:  burst target # \n\n Call upon the power of your faith to bathe the target in holy light.",
		permissions.Paladin,
		"burst", "flash", "bur")
}

type burst cmd

func (burst) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Burst what exactly?")
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
	ready, msg := s.actor.TimerReady("combat_burst")
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

		missChance := DetermineMissChance(s, whatMob.Level-s.actor.Tier)
		if objects.DayTime {
			missChance -= 20
		} else {
			missChance += 20
		}
		// Check for a miss
		if utils.Roll(100, 1, 0) <= missChance {
			s.msg.Actor.SendBad("Your burst of light fails!!")
			s.actor.SetTimer("combat_burst", config.CircleTimer)
			s.actor.SetTimer("combat", config.CombatCooldown)
			whatMob.AddThreatDamage(1, s.actor)
			whatMob.CurrentTarget = s.actor.Name
			s.msg.Observers.SendBad(s.actor.Name + " burst of light fails to dazzle " + whatMob.Name)
			return
		}

		whatMob.Stun(config.CircleStuns)
		whatMob.AddThreatDamage(whatMob.Stam.Max/2, s.actor)
		whatMob.CurrentTarget = s.actor.Name
		s.actor.SetTimer("combat_burst", config.CircleTimer)
		s.actor.SetTimer("combat", config.CombatCooldown)
		s.msg.Actor.SendInfo("You bathed " + whatMob.Name + " in holy light.")
		s.msg.Observers.SendInfo(s.actor.Name + " bathes " + whatMob.Name + " in holy light.")
		return

	}

	s.msg.Actor.SendInfo("Burst what?")
	s.ok = true
}
