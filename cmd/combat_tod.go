package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(tod{},
		"Usage:  touch target # \n\n Attempt the secret art of a touch of death on your living target",
		permissions.Monk,
		"touch", "tod", "touch-of-death")
}

type tod cmd

func (tod) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Turn what exactly?")
		return
	}
	if s.actor.Tier < 10 {
		s.msg.Actor.SendBad("You aren't high enough level to perform that skill.")
		return
	}
	// Check some timers
	ready, msg := s.actor.TimerReady("combat_tod")
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
		if whatMob.Flags["undead"] != false {
			s.msg.Actor.SendBad("Your target is undead and unaffected by your chi!")
			return
		}

		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to perform a touch of death on them.")
			return
		}

		if s.actor.Mana.Current < config.TodCost {
			s.msg.Actor.SendBad("You do not have enough chi to perform that.")
			return
		}
		s.actor.Victim = whatMob

		s.actor.RunHook("combat")
		s.actor.SetTimer("combat_tod", config.TodTimer)
		s.actor.SetTimer("combat", config.CombatCooldown)
		s.actor.Mana.Subtract(config.TodCost)
		// base chance is 15% to hide
		curChance := config.TodMax
		if whatMob.Level > s.actor.Tier {
			curChance -= config.TodScaleDown * (whatMob.Level - s.actor.Tier)
		} else if s.actor.Tier > whatMob.Level {
			curChance += config.TodScaleDown * (s.actor.Tier - whatMob.Level)
		}

		todRoll := utils.Roll(100, 1, 0)
		if todRoll <= config.VitalChance {
			s.msg.Actor.SendInfo("Your chi flows through you and you perform a perfect touch of death on " + whatMob.Name + " and kill them.")
			s.msg.Observers.SendInfo(s.actor.Name + " touches " + whatMob.Name + " and kills them.")
			s.actor.AdvanceSkillExp(int((float64(whatMob.Stam.Max) * float64(whatMob.Experience)) * config.Classes[config.AvailableClasses[s.actor.Class]].WeaponAdvancement))
			whatMob.Stam.Current = 0
			DeathCheck(s, whatMob)
			whatMob = nil
		} else if curChance >= 100 || todRoll <= curChance {
			s.msg.Actor.SendInfo("You focus your chi and perform a touch of death on " + whatMob.Name + "!")
			s.msg.Observers.SendInfo(s.actor.Name + " performed a touch of death on " + whatMob.Name)
			whatMob.AddThreatDamage(whatMob.Stam.Current/2, s.actor)
			s.actor.AdvanceSkillExp(int((float64(whatMob.Stam.Max) / 2 * float64(whatMob.Experience)) * config.Classes[config.AvailableClasses[s.actor.Class]].WeaponAdvancement))
			whatMob.Stam.Subtract(whatMob.Stam.Current / 2)
		} else {
			s.msg.Actor.SendBad("You misperform the touch of death " + whatMob.Name + ".  They charge you!")
			whatMob.CurrentTarget = s.actor.Name
			whatMob.AddThreatDamage(whatMob.Stam.Current, s.actor)
			s.actor.ReceiveDamage(s.actor.Stam.Max / 2)
			s.msg.Observers.SendInfo(s.actor.Name + " turn attempt fails and enrages " + whatMob.Name)
		}
		return
	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}
