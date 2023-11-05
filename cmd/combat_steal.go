package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"math"
	"strconv"
)

// Syntax: ( INVENTORY | INV )
func init() {
	addHandler(steal{},
		"Usage:  steal target item \n \n Try to steal an item from a targets inventory",
		permissions.Thief,
		"steal")
}

type steal cmd

func (steal) process(s *state) {
	if s.actor.Tier < config.MinorAbilityTier {
		s.msg.Actor.SendBad("You must be at least tier " + strconv.Itoa(config.MinorAbilityTier) + " to use this skill.")
		return
	}

	if len(s.input) < 2 {
		s.msg.Actor.SendBad("Steal what from who")
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

	// Check some timers
	ready, msg := s.actor.TimerReady("steal")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	if s.actor.Flags["hidden"] != true {
		s.msg.Actor.SendBad("You can't steal while standing out in the open.")
		return
	}

	if s.actor.Tier < 5 {
		s.msg.Actor.SendBad("You must be level 5 before you can steal.")
		return
	}

	targetStr := s.words[0]
	targetNum := 1
	nameStr := ""
	nameNum := 1

	if len(s.words) > 1 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		} else {
			nameStr = s.words[1]
		}
	}
	if len(s.words) > 2 {
		if val2, err2 := strconv.Atoi(s.words[2]); err2 == nil {
			nameNum = val2
		} else {
			nameStr = s.words[2]
		}
	}

	if len(s.words) > 3 {
		if val3, err3 := strconv.Atoi(s.words[3]); err3 == nil {
			nameNum = val3
		}
	}

	if nameStr == "" {
		s.msg.Actor.SendBad("Steal what from who?")
		return
	}

	// TODO: Steal from players inventory if PvP flag is set

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(targetStr, targetNum, s.actor)
	if whatMob != nil {
		if whatMob.CheckFlag("no_steal") {
			s.msg.Actor.SendBad("Try as you might you can not find a way to steal from this enemy.")
			return
		}

		if whatMob.Placement != s.actor.Placement {
			s.msg.Actor.SendBad("You are too far away to steal from ", whatMob.Name)
			return
		}

		if _, ok := whatMob.ThreatTable[s.actor.Name]; ok {
			s.msg.Actor.SendBad("This mob is already aware of you, you can't steal from them.")
			return
		}

		what := whatMob.Inventory.Search(nameStr, nameNum)
		if what != nil {
			s.actor.SetTimer("steal", config.StealCD)
			if (s.actor.GetCurrentWeight() + what.GetWeight()) <= s.actor.MaxWeight() {
				// base chance is 15% to hide
				//curChance := config.StealChance + (config.StealChancePerLevel * (s.actor.Tier - whatMob.Level))
				//curChance += s.actor.Dex.Current * config.StealChancePerPoint

				curChance := config.StealChance + (s.actor.Dex.Current * config.StealChancePerPoint) + (config.StealthLevel(s.actor.Skills[11].Value) * config.StealChancePerSkillLevel)
				curChance -= int(what.Weight / 2)
				lvlDiff := float64(whatMob.Level - s.actor.Tier)
				if lvlDiff > 2 {
					lvlDiff = (lvlDiff - 2) * 0.2
					curChance -= int(float64(curChance) * lvlDiff)
				}

				//s.msg.Actor.SendInfo("Steal chance = " + strconv.Itoa(curChance))

				if curChance > 95 {
					curChance = 95
				}

				if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					curChance = 100
				}

				//log.Println(s.actor.Name+"Steal Chance Roll: ", curChance)

				if utils.Roll(100, 1, 0) <= curChance {
					whatMob.Inventory.Remove(what)
					s.actor.Inventory.Add(what)
					s.actor.AdvanceStealthExp(int(float64(what.Value) * 0.5))
					s.msg.Actor.SendGood("You steal a ", what.Name, " from ", whatMob.Name, ".")
					return
				} else {
					s.msg.Actor.SendBad("You failed to steal from ", whatMob.Name, ", and stumble out of the shadows.")
					s.msg.Observers.SendBad(s.actor.Name + " fails to steal from " + whatMob.Name)
					s.actor.RemoveHook("combat", "hide")
					whatMob.AddThreatDamage(whatMob.Stam.Max/4, s.actor)
					if utils.Roll(100, 1, 0) <= config.MobStealRevengeVitalChance {
						whatMob.CurrentTarget = s.actor.Name
						s.msg.Actor.SendInfo(whatMob.Name + " turns to you.")
						s.msg.Observers.SendInfo(whatMob.Name + " turns to " + s.actor.Name + ".")
						vitDamage, resisted := s.actor.ReceiveVitalDamage(int(math.Ceil(float64(whatMob.InflictDamage() * config.VitalStrikeScale))))
						data.StoreCombatMetric("steal_fail_vital", 0, 0, vitDamage+resisted, resisted, vitDamage, 1, whatMob.MobId, whatMob.Level, 0, s.actor.CharId)
						if vitDamage == 0 {
							s.msg.Actor.SendGood(whatMob.Name, " vital strike bounces off of you for no damage!")
						} else {
							s.msg.Actor.SendInfo(whatMob.Name, " attacks you for "+strconv.Itoa(vitDamage)+" points of vitality damage!")
						}
						s.actor.DeathCheck("was slain trying to steal from " + whatMob.Name + ".")
					}
					return
				}
			} else {
				s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
				return
			}
		} else {
			s.msg.Actor.SendInfo("That item isn't on the target.")
			return
		}
	} else {
		s.msg.Actor.SendBad("What are you trying to steal from?")
		s.ok = true
		return
	}

}
