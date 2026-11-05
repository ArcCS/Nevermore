package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"strconv"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(bash{},
		"Usage:  bash target # \n\n Bash the target to stun it.  Bash is an opener: it only works on a target you have not yet attacked, and a miss uses it up.",
		permissions.Barbarian,
		"bash")
}

type bash cmd

func (bash) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Bash what exactly?")
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
	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	s.actor.RunHook("combat")

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

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == nil {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		// Bash needs the heft of a blunt weapon or a two-hander.
		if s.actor.Equipment.Main.ItemType != 2 && !s.actor.Equipment.Main.IsTwoHanded() {
			s.msg.Actor.SendBad("You can only bash with a blunt or two-handed weapon.")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to bash them.")
			return
		}

		// Bash is an opener: once you have attacked a mob in any way,
		// including a missed bash, you cannot bash it.
		if whatMob.AttackedBy(s.actor.Name) {
			s.msg.Actor.SendBad("You've already engaged " + whatMob.Name + "; bash only works as an opener.")
			return
		}
		whatMob.MarkAttackedBy(s.actor)

		recklessAttack(s)
		// Check for a miss
		if utils.Roll(100, 1, 0) <= DetermineMissChance(s, whatMob.Level-s.actor.Tier, whatMob.EngagedCount(s.actor)) {
			s.msg.Actor.SendBad("You missed!!")
			s.msg.Observers.SendBad(s.actor.Name + " fails to bash " + whatMob.Name)
			whatMob.AddThreatDamage(1, s.actor)
			whatMob.CurrentTarget = s.actor.Name
			s.actor.SetTimer("combat", config.CombatCooldown)
			data.StoreCombatMetric("bash-miss", 0, 0, 0, 0, 0, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			return
		}

		s.actor.Victim = whatMob
		// The special roll uses the skill of the weapon in hand, blunt or two-handed.
		damageModifier, stunModifier, bashMsg := config.RollBash(attackSkillLevel(s), s.actor.Tier)
		whatMob.StunBy(config.BashStuns*stunModifier, "bash-stun", s.actor)
		actualDamage, _, resisted := whatMob.ReceiveDamage(config.BashDamage(s.actor.InflictDamage(), damageModifier, s.actor.Tier))
		data.StoreCombatMetric("bash", 0, 0, actualDamage+resisted, resisted, actualDamage, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
		whatMob.AddThreatDamage(actualDamage+config.ThreatPercent(whatMob.Stam.Max, config.SmallThreatPercent), s.actor)
		s.actor.AdvanceSkillExp((float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience())))
		s.msg.Actor.SendInfo(bashMsg)
		whatMob.CurrentTarget = s.actor.Name
		s.msg.Actor.SendInfo("You bashed the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
		s.msg.Observers.SendInfo(s.actor.Name + " bashes " + whatMob.Name)
		if whatMob.CheckFlag("reflection") {
			reflectDamage := int(float64(actualDamage) * config.ReflectDamageFromMob)
			stamDamage, vitDamage, resisted := s.actor.ReceiveDamage(reflectDamage)
			data.StoreCombatMetric("bash_mob_reflect", 0, 0, stamDamage+vitDamage+resisted, resisted, stamDamage+vitDamage, 1, whatMob.MobId, whatMob.Level, 0, s.actor.CharId)
			s.msg.Actor.Send("The " + whatMob.Name + " reflects " + strconv.Itoa(reflectDamage) + " damage back at you!")
			s.actor.DeathCheck(" was killed by reflection!")
		}
		if bashMsg == "Thunk!!" {
			shatterWeapon(s)
		}
		DeathCheck(s, whatMob)
		s.actor.SetTimer("combat", config.CombatCooldown)
		return
	}

	s.msg.Actor.SendInfo("Bash what?")
	s.ok = true
}
