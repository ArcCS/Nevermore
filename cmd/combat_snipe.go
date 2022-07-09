package cmd

import (
	"github.com/ArcCS/Nevermore/config"
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

	if s.actor.Tier < 10 {
		s.msg.Actor.SendBad("You must be at least tier 10 to use this skill.")
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
		curChance := config.SnipeChance + (config.SnipeChancePerLevel * (s.actor.Tier - whatMob.Level))

		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			curChance = 100
		}

		curChance += s.actor.Dex.Current * config.SnipeChancePerPoint
		whatMob.AddThreatDamage(whatMob.Stam.Max/10, s.actor)
		if curChance >= 100 || utils.Roll(100, 1, 0) <= curChance {
			actualDamage, _ := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * float64(config.CombatModifiers["snipe"]))))
			s.msg.Actor.SendInfo("You sniped the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
			s.actor.AdvanceSkillExp(int((float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience)) * config.Classes[config.AvailableClasses[s.actor.Class]].WeaponAdvancement))
			s.msg.Observers.SendInfo(s.actor.Name + " snipes " + whatMob.Name)
			DeathCheck(s, whatMob)
			s.actor.SetTimer("combat", config.CombatCooldown)
			return
		} else {
			s.msg.Actor.SendBad("You failed to snipe ", whatMob.Name, "!")
			s.msg.Observers.SendBad(s.actor.Name+" failed to snipe ", whatMob.Name, "!")
			if utils.Roll(100, 1, 0) < config.SnipeFumbleChance {
				s.msg.Actor.SendBad("You fumbled your weapon!")
				s.msg.Observer.SendInfo(s.actor.Name + " fails to snipe and fumbles their weapon. ")
				s.actor.SetTimer("global", 25)
				_, what := s.actor.Equipment.Unequip(s.actor.Equipment.Main.Name)
				if what != nil {
					s.actor.Inventory.Lock()
					s.actor.Inventory.Add(what)
					s.actor.Inventory.Unlock()
				}
			} else {
				s.msg.Observer.SendInfo(s.actor.Name + " fails to snipe " + whatMob.Name)
			}
			s.ok = true
			return
		}
	}

	s.msg.Actor.SendInfo("Bash what?")
	s.ok = true
}
