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
	addHandler(snipe{},
		"Usage:  snipe target # \n\n Snipe the target, can only be done while hidden",
		permissions.Ranger,
		"snipe")
}

type snipe cmd

func (snipe) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Snipe what exactly?")
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
		s.msg.Actor.SendBad("You must be hidden to snipe.")
		return
	}

	ready, msg := s.actor.TimerReady("combat")
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
		s.actor.Victim = whatMob
		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == nil {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		// Shortcut weapon not being blunt
		if s.actor.Equipment.Main.ItemType != 4 {
			s.msg.Actor.SendBad("You can only snipe with a ranged weapon.")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Placement == whatMob.Placement {
			s.msg.Actor.SendBad("You are too close to snipe them.")
			return
		}

		_, ok := whatMob.ThreatTable[s.actor.Name]
		if ok {
			s.msg.Actor.SendBad("You have already engaged ", whatMob.Name, " in combat!")
			return
		}

		s.actor.RunHook("combat")

		curChance := config.SnipeChance + (s.actor.Dex.Current * config.SnipeChancePerPoint) + (config.SnipeChancePerLevel * (s.actor.Tier - whatMob.Level))

		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			curChance = 100
		}

		whatMob.AddThreatDamage(whatMob.Stam.Max/10, s.actor)
		if curChance >= 100 || utils.Roll(100, 1, 0) <= curChance {
			actualDamage, _, resisted := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * float64(config.CombatModifiers["snipe"]))))
			data.StoreCombatMetric("snipe", 0, 0, actualDamage+resisted, resisted, actualDamage, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			s.msg.Actor.SendInfo("You sniped the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
			s.actor.AdvanceSkillExp(int((float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience)) * config.Classes[config.AvailableClasses[s.actor.Class]].WeaponAdvancement))
			s.msg.Observers.SendInfo(s.actor.Name + " snipes " + whatMob.Name)
			if whatMob.CheckFlag("reflection") {
				reflectDamage := int(float64(actualDamage) * config.ReflectDamageFromMob)
				stamDamage, vitDamage, resisted := s.actor.ReceiveDamage(reflectDamage)
				data.StoreCombatMetric("snipe_mob_reflect", 0, 0, stamDamage+vitDamage+resisted, resisted, stamDamage+vitDamage, 1, whatMob.MobId, whatMob.Level, 0, s.actor.CharId)
				s.msg.Actor.Send("The " + whatMob.Name + " reflects " + strconv.Itoa(reflectDamage) + " damage back at you!")
				s.actor.DeathCheck(" was killed by reflection!")
			}
			DeathCheck(s, whatMob)
			s.actor.SetTimer("combat", config.CombatCooldown)
			return
		} else {
			s.msg.Actor.SendBad("You failed to snipe ", whatMob.Name, "!")
			s.msg.Observers.SendBad(s.actor.Name+" failed to snipe ", whatMob.Name, "!")
			data.StoreCombatMetric("snipe-miss", 0, 0, 0, 0, 0, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			if utils.Roll(100, 1, 0) < config.SnipeFumbleChance {
				s.msg.Actor.SendBad("You fumbled your weapon!")
				s.msg.Observer.SendInfo(s.actor.Name + " fails to snipe and fumbles their weapon. ")
				_, what := s.actor.Equipment.Unequip(s.actor.Equipment.Main.Name)
				if what != nil {
					s.actor.Inventory.Add(what)
				}
			} else {
				s.msg.Observer.SendInfo(s.actor.Name + " fails to snipe " + whatMob.Name)
			}
			s.ok = true
			return
		}
	}

	s.msg.Actor.SendInfo("Snipe what?")
	s.ok = true
}
