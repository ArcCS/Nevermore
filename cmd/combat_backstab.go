package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"math"
	"strconv"
	"strings"
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

	if s.actor.Tier < 10 {
		s.msg.Actor.SendBad("You must be at least tier 10 to use this skill.")
		return
	}

	if s.actor.Flags["hidden"] != true {
		s.msg.Actor.SendBad("You must be hidden to backstab.")
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

		_, ok := whatMob.ThreatTable[s.actor.Name]
		if ok {
			s.msg.Actor.SendBad("You have already engaged ", whatMob.Name, " in combat!")
			return
		}

		curChance := config.BackStabChance + (config.BackStabChancePerLevel * (s.actor.Tier - whatMob.Level))

		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			curChance = 100
		}

		curChance += s.actor.Dex.Current * config.BackStabChancePerPoint

		s.actor.Victim = whatMob
		s.actor.RunHook("combat")
		if curChance >= 100 || utils.Roll(100, 1, 0) <= curChance {
			actualDamage, _ := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * float64(config.CombatModifiers["backstab"]))))
			s.actor.AdvanceSkillExp(int((float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience)) * config.Classes[config.AvailableClasses[s.actor.Class]].WeaponAdvancement))
			whatMob.AddThreatDamage(actualDamage, s.actor)
			s.msg.Actor.SendInfo("You backstabbed the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
			s.msg.Observers.SendInfo(s.actor.Name + " backstabs " + whatMob.Name)
			DeathCheck(s, whatMob)
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
			if utils.Roll(100, 1, 0) <= config.MobBSRevengeVitalChance {
				whatMob.CurrentTarget = s.actor.Name
				s.msg.Actor.SendInfo(whatMob.Name + " turns it's attention to you.")
				s.msg.Observers.SendInfo(whatMob.Name + " turns to " + s.actor.Name + ".")
				vitDamage := s.actor.ReceiveVitalDamage(int(math.Ceil(float64(whatMob.InflictDamage() * config.VitalStrikeScale))))
				if vitDamage == 0 {
					s.msg.Actor.SendGood(whatMob.Name, " vital strike bounces off of you for no damage!")
				} else {
					s.msg.Actor.SendInfo(whatMob.Name, " attacks you for "+strconv.Itoa(vitDamage)+" points of vitality damage!")
				}
				s.actor.DeathCheck("was slain while trying to backstab a " + strings.Title(whatMob.Name))
			}
			s.ok = true
			return
		}
	}

	s.msg.Actor.SendInfo("Backstab what?")
	s.ok = true
}
