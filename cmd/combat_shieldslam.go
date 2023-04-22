package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"strconv"
)

func init() {
	addHandler(slam{},
		"Usage:  shield-slam target # \n\n Slam your shield into the target",
		permissions.Paladin,
		"shield", "shield-slam", "slam")
}

type slam cmd

func (slam) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Slam what exactly?")
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

	if s.actor.Tier < 7 {
		s.msg.Actor.SendBad("You must be at least tier 7 to use this skill.")
		return
	}
	// Check some timers
	ready, msg := s.actor.TimerReady("combat_shieldslam")
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
		_, ok := whatMob.ThreatTable[s.actor.Name]
		if ok {
			s.msg.Actor.SendBad("You have already engaged ", whatMob.Name, " in combat!")
			return
		}
		s.actor.RunHook("combat")
		s.actor.Victim = whatMob
		// Shortcut a missing weapon:
		if s.actor.Equipment.Off == nil {
			s.msg.Actor.SendBad("You have nothing equipped in your offhand!")
			return
		}

		// Shortcut weapon not being blunt
		if s.actor.Equipment.Off.ItemType != 23 {
			s.msg.Actor.SendBad("You can only bash with a shield!")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to bash them.")
			return
		}

		actualDamage, _ := whatMob.ReceiveDamage(s.actor.GetStat("str") * config.ShieldDamage)
		whatMob.AddThreatDamage(whatMob.Stam.Max/10, s.actor)
		whatMob.Stun(config.ShieldStun * s.actor.GetStat("pie"))
		s.actor.AdvanceSkillExp(int((float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience)) * config.Classes[config.AvailableClasses[s.actor.Class]].WeaponAdvancement))
		s.msg.Actor.SendInfo("You slammed the " + whatMob.Name + " with your shield for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
		s.msg.Observers.SendInfo(s.actor.Name + " slams " + config.TextPosPronoun[s.actor.Gender] + " shield into " + whatMob.Name)
		DeathCheck(s, whatMob)
		s.actor.SetTimer("combat_shieldslam", config.SlamTimer)
		s.actor.SetTimer("combat", config.CombatCooldown)
		return
	}

	s.msg.Actor.SendInfo("Slam what?")
	s.ok = true
}
