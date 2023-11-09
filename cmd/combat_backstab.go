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
	addHandler(backstab{},
		"Usage:  backstab target # \n\n Backstab the target, can only be done while hidden",
		permissions.Thief,
		"backstab", "bs")
}

type backstab cmd

func (backstab) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Backstab what exactly?")
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

	if s.actor.Flags["hidden"] != true {
		s.msg.Actor.SendBad("You must be hidden to backstab.")
		return
	}

	ready, msg := s.actor.TimerReady("combat_backstab")
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

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == nil {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		// Shortcut weapon not being blunt
		if s.actor.Equipment.Main.ItemType != 0 && s.actor.Equipment.Main.ItemType != 1 {
			s.msg.Actor.SendBad("You can only backstab with sharp and thrust weapons.")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to backstab them.")
			return
		}

		if whatMob.Flags["undead"] == true {
			s.msg.Actor.SendBad("Your target is undead and you cannot find their vitals!")
			return
		}

		_, ok := whatMob.ThreatTable[s.actor.Name]
		if ok {
			s.msg.Actor.SendBad("You have already engaged ", whatMob.Name, " in combat!")
			return
		}

		//curChance := config.BackStabChance + (s.actor.Dex.Current * config.BackStabChancePerPoint) + (config.BackStabChancePerLevel * (s.actor.Tier - whatMob.Level))

		curChance := config.BackStabChance + (s.actor.Dex.Current * config.BackStabChancePerPoint) + (config.StealthLevel(s.actor.Skills[11].Value) * config.BackStabChancePerSkillLevel)
		lvlDiff := float64(whatMob.Level - s.actor.Tier)
		if lvlDiff > 1 {
			lvlDiff = (lvlDiff - 1) * .125
			curChance -= int(float64(curChance) * lvlDiff)
		} else if lvlDiff == 1 {
			curChance -= int(float64(curChance) * 0.05)
		}

		//s.msg.Actor.SendInfo("BS chance = " + strconv.Itoa(curChance))

		if curChance > 95 {
			curChance = 95
		}

		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			curChance = 100
		}

		s.actor.Victim = whatMob
		s.actor.RunHook("combat")
		if curChance >= 100 || utils.Roll(100, 1, 0) <= curChance {

			actualDamage, _, resisted := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * float64(config.CombatModifiers["backstab"]))))
			data.StoreCombatMetric("backstab", 0, 0, actualDamage+resisted, resisted, actualDamage, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			s.actor.AdvanceSkillExp(int((float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience)) * config.Classes[config.AvailableClasses[s.actor.Class]].WeaponAdvancement))
			s.actor.AdvanceStealthExp(int(float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience)))
			whatMob.AddThreatDamage(actualDamage, s.actor)
			s.msg.Actor.SendInfo("You backstabbed the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
			s.msg.Observers.SendInfo(s.actor.Name + " backstabs " + whatMob.Name)
			if whatMob.CheckFlag("reflection") {
				reflectDamage := int(float64(actualDamage) * config.ReflectDamageFromMob)
				stamDamage, vitDamage, resisted := s.actor.ReceiveDamage(reflectDamage)
				data.StoreCombatMetric("backstab_mob_reflect", 0, 0, stamDamage+vitDamage+resisted, resisted, stamDamage+vitDamage, 1, whatMob.MobId, whatMob.Level, 0, s.actor.CharId)
				s.msg.Actor.Send("The " + whatMob.Name + " reflects " + strconv.Itoa(reflectDamage) + " damage back at you!")
				s.actor.DeathCheck(" was killed by reflection!")
			}
			DeathCheck(s, whatMob)
			s.actor.SetTimer("combat_backstab", config.BackstabCooldown)
			s.actor.SetTimer("combat", config.CombatCooldown)
			msg := s.actor.Equipment.DamageWeapon("main", 4)
			if msg != "" {
				s.msg.Actor.SendInfo(msg)
			}
			return
		} else {
			s.msg.Actor.SendBad("You failed to backstab ", whatMob.Name, ", and are vulnerable to attack!")
			s.msg.Observers.SendBad(s.actor.Name+" failed to backstab ", whatMob.Name, ", and is vulnerable to attack!")
			whatMob.AddThreatDamage(whatMob.Stam.Max/2, s.actor)
			s.actor.SetTimer("combat", config.CombatCooldown)
			if utils.Roll(100, 1, 0) <= config.MobBSRevengeVitalChance {
				whatMob.CurrentTarget = s.actor.Name
				s.msg.Actor.SendInfo(whatMob.Name + " turns it's attention to you.")
				s.msg.Observers.SendInfo(whatMob.Name + " turns to " + s.actor.Name + ".")
				vitDamage, resisted := s.actor.ReceiveVitalDamage(int(math.Ceil(float64(whatMob.InflictDamage() * config.VitalStrikeScale))))
				data.StoreCombatMetric("backstab_mob_vital", 0, 0, vitDamage+resisted, resisted, vitDamage, 1, whatMob.MobId, whatMob.Level, 0, s.actor.CharId)
				if vitDamage == 0 {
					s.msg.Actor.SendGood(whatMob.Name, " vital strike bounces off of you for no damage!")
				} else {
					s.msg.Actor.SendInfo(whatMob.Name, " attacks you for "+strconv.Itoa(vitDamage)+" points of vitality damage!")
					if s.actor.CheckFlag("reflection") {
						reflectDamage := int(float64(vitDamage) * (float64(s.actor.GetStat("int")) * config.ReflectDamagePerInt))
						whatMob.ReceiveDamage(reflectDamage)
						s.msg.Actor.Send(text.Cyan + "You reflect " + strconv.Itoa(reflectDamage) + " damage back to " + whatMob.Name + "!\n" + text.Reset)
						whatMob.DeathCheck(s.actor)
					}
				}
				s.actor.DeathCheck("was slain while trying to backstab a " + utils.Title(whatMob.Name))
			} else {
				data.StoreCombatMetric("backstab-miss", 0, 0, 0, 0, 0, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			}
			s.ok = true
			return
		}
	}

	s.msg.Actor.SendInfo("Backstab what?")
	s.ok = true
}
