package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"math"
	"strconv"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(bash{},
		"Usage:  bash target # \n\n Bash the target",
		permissions.Barbarian,
		"bash")
}

type bash cmd

func (bash) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Bash what exactly?")
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

	if s.actor.Tier < config.MinorAbilityTier {
		s.msg.Actor.SendBad("You must be at least tier " + strconv.Itoa(config.MinorAbilityTier) + " to use this skill.")
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("combat_bash")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
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

		// Shortcut weapon not being blunt
		if s.actor.Equipment.Main.ItemType != 2 {
			s.msg.Actor.SendBad("You can only bash with a blunt weapon.")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to bash them.")
			return
		}

		// Check for a miss
		if utils.Roll(100, 1, 0) <= DetermineMissChance(s, whatMob.Level-s.actor.Tier) {
			s.msg.Actor.SendBad("You missed!!")
			s.msg.Observers.SendBad(s.actor.Name + " fails to bash " + whatMob.Name)
			s.actor.SetTimer("combat", config.CombatCooldown)
			data.StoreCombatMetric("bash-miss", 0, 0, 0, 0, 0, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			return
		}

		s.actor.Victim = whatMob
		// Check the rolls in reverse order from hardest to lowest for bash rolls.
		damageModifier, stunModifier, bashMsg := config.RollBash(config.WeaponLevel(s.actor.Skills[2].Value, s.actor.Class))
		whatMob.Stun(config.BashStuns * stunModifier)
		actualDamage, _, resisted := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * float64(damageModifier))))
		data.StoreCombatMetric("bash", 0, 0, actualDamage+resisted, resisted, actualDamage, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
		whatMob.AddThreatDamage(actualDamage, s.actor)
		s.actor.AdvanceSkillExp(int((float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience)) * config.Classes[config.AvailableClasses[s.actor.Class]].WeaponAdvancement))
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
		DeathCheck(s, whatMob)
		s.actor.SetTimer("combat_bash", config.BashTimer)
		s.actor.SetTimer("combat", config.CombatCooldown)
		return
	}

	s.msg.Actor.SendInfo("Bash what?")
	s.ok = true
}
