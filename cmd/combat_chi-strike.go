package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"math"
	"strconv"
)

func init() {
	addHandler(chiStrike{},
		"Usage:  Use your chi to strike a target",
		permissions.Monk,
		"chi-strike", "cs")
}

type chiStrike cmd

func (chiStrike) process(s *state) {
	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything!")
		return
	}

	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Strike what exactily?")
		return
	}
	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	if s.actor.Tier < 3 {
		s.msg.Actor.SendBad("You must be at least tier 3 to use this skill.")
		return
	}

	// Check some timers

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

		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to perform a chi-strike on them.")
			return
		}

		if s.actor.Mana.Current < config.ChiStrikeCost {
			s.msg.Actor.SendBad("You do not have enough chi to perform a chi-strike.")
			return
		}
		s.actor.Victim = whatMob

		s.actor.RunHook("combat")
		s.actor.Mana.Subtract(config.ChiStrikeCost)
		actualDamage, _, resisted := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.Tier) * float64(s.actor.Dex.Current) / float64((15 - config.WeaponLevel(s.actor.Skills[5].Value, 8) + utils.Roll(5, 1, 0))))))
		data.StoreCombatMetric("Chi-Strike", 0, 0, actualDamage+resisted, resisted, actualDamage, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
		s.actor.AdvanceSkillExp(int((float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience)) * config.Classes[config.AvailableClasses[s.actor.Class]].WeaponAdvancement))
		whatMob.AddThreatDamage(actualDamage, s.actor)
		s.msg.Actor.SendInfo("You strike the " + whatMob.Name + " with your chi for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
		s.msg.Observers.SendInfo(s.actor.Name + " chi-strikes " + whatMob.Name)
		if whatMob.CheckFlag("reflection") {
			reflectDamage := int(float64(actualDamage) * config.ReflectDamageFromMob)
			stamDamage, vitDamage, resisted := s.actor.ReceiveDamage(reflectDamage)
			data.StoreCombatMetric("backstab_mob_reflect", 0, 0, stamDamage+vitDamage+resisted, resisted, stamDamage+vitDamage, 1, whatMob.MobId, whatMob.Level, 0, s.actor.CharId)
			s.msg.Actor.Send("The " + whatMob.Name + " reflects " + strconv.Itoa(reflectDamage) + " damage back at you!")
			s.actor.DeathCheck(" was killed by reflection!")
		}
		DeathCheck(s, whatMob)
		return
	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}
