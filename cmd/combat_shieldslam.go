package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
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

	if s.actor.Tier < config.SpecialAbilityTier {
		s.msg.Actor.SendBad("You must be at least tier " + strconv.Itoa(config.SpecialAbilityTier) + " to use this skill.")
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
		s.actor.RunHook("combat")
		s.actor.Victim = whatMob
		// Shortcut a missing weapon:
		if s.actor.Equipment.Off == nil {
			s.msg.Actor.SendBad("You have nothing equipped in your offhand!")
			return
		}

		// Shortcut weapon not being blunt
		if s.actor.Equipment.Off.ItemType != 23 {
			s.msg.Actor.SendBad("You can only slam with a shield!")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to slam them.")
			return
		}

		actualDamage, _, resisted := whatMob.ReceiveDamage(s.actor.GetStat("str") * config.ShieldDamage)
		data.StoreCombatMetric("shieldslam", 0, 0, actualDamage+resisted, resisted, actualDamage, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
		whatMob.AddThreatDamage(whatMob.Stam.Max/10, s.actor)
		whatMob.Stun(int(config.ShieldStun * float64(s.actor.GetStat("pie"))))
		whatMob.CurrentTarget = s.actor.Name
		s.msg.Actor.SendInfo("You slammed the " + whatMob.Name + " with your shield for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
		s.msg.Observers.SendInfo(s.actor.Name + " slams " + config.TextPosPronoun[s.actor.Gender] + " shield into " + whatMob.Name)
		if whatMob.CheckFlag("reflection") {
			reflectDamage := int(float64(actualDamage) * config.ReflectDamageFromMob)
			stamDamage, vitDamage, resisted := s.actor.ReceiveDamage(reflectDamage)
			data.StoreCombatMetric("shieldslam_mob_reflect", 0, 0, stamDamage+vitDamage+resisted, resisted, stamDamage+vitDamage, 1, whatMob.MobId, whatMob.Level, 0, s.actor.CharId)
			s.msg.Actor.Send("The " + whatMob.Name + " reflects " + strconv.Itoa(reflectDamage) + " damage back at you!")
			s.actor.DeathCheck(" was killed by reflection!")
		}
		DeathCheck(s, whatMob)
		s.actor.SetTimer("combat_shieldslam", config.SlamTimer)
		s.actor.SetTimer("combat", config.CombatCooldown)
		return
	}

	s.msg.Actor.SendInfo("Slam what?")
	s.ok = true
}
