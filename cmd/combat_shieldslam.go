package cmd

import (
	"strconv"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(slam{},
		"Usage:  shield-slam target # \n\n Slam your shield into the target to stun and taunt it.  Shield slam is an opener: it only works on a target you have not yet attacked.",
		permissions.Paladin|permissions.Fighter,
		"shield", "shield-slam", "slam")
}

type slam cmd

func (slam) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Slam what exactly?")
		return
	}
	if !s.requireSight() {
		return
	}

	if !s.requireStamina() {
		return
	}

	if !s.requireTier(config.MinorAbilityTier) {
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

		// Shield slam is an opener, like bash: once you have attacked a mob
		// in any way, you cannot slam it.
		if whatMob.AttackedBy(s.actor.Name) {
			s.msg.Actor.SendBad("You've already engaged " + whatMob.Name + "; shield slam only works as an opener.")
			return
		}

		shield := s.actor.Equipment.Off
		actualDamage, _, resisted := whatMob.ReceiveDamage(config.ShieldSlamDamage(config.ShieldSlamStat(s.actor.Class, s.actor.GetStat("str"), s.actor.GetStat("pie")), s.actor.Tier, config.RollShieldArmor(shield.Armor)))
		data.StoreCombatMetric("shieldslam", 0, 0, actualDamage+resisted, resisted, actualDamage, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
		whatMob.AddThreatDamage(actualDamage+config.ThreatPercent(whatMob.Stam.Max, config.LargeThreatPercent), s.actor)
		whatMob.CurrentTarget = s.actor.Name
		s.msg.Actor.SendInfo("You slammed the " + whatMob.Name + " with your shield for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
		// The slam always lands; a mob out of level may shrug off the stun.
		failChance := 0
		if !s.where.InQuestMode() {
			failChance = config.ShieldSlamStunFailChance(whatMob.Level-s.actor.Tier, whatMob.EngagedCount(s.actor))
		}
		if failChance > 0 && utils.Roll(100, 1, 0) <= failChance {
			s.msg.Actor.SendBad(whatMob.Name + " shrugs off the blow, unstunned!")
		} else {
			whatMob.StunBy(config.ShieldSlamStuns, "shieldslam-stun", s.actor)
		}
		s.msg.Observers.SendInfo(s.actor.Name + " slams " + config.PosPronoun(s.actor.Gender) + " shield into " + whatMob.Name)
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
