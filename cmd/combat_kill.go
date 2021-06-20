package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"math"
	"strconv"
)

func init() {
	addHandler(kill{},
		"Usage:  kill target # \n\n Try to attack something! Can also use attack, or shorthand k",
		permissions.Player,
		"kill", "k")
}

type kill cmd

func (kill) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Attack what exactly?")
		return
	}

	// Check some timers
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

		// This is an override for a GM to delete a mob
		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			s.msg.Actor.SendInfo("You smashed ", whatMob.Name, " out of existence.")
			objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
			whatMob = nil
			return
		}

		s.actor.RunHook("combat")

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == (*objects.Item)(nil) && s.actor.Class != 8 {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		if _, err := whatMob.ThreatTable[s.actor.Name]; !err {
			s.msg.Actor.Send(text.White + "You engaged " + whatMob.Name + " #" + strconv.Itoa(s.where.Mobs.GetNumber(whatMob)) + " in combat.")
			s.msg.Observers.Send(text.White + s.actor.Name + " attacks " + whatMob.Name)
			whatMob.AddThreatDamage(0, s.actor)
		}

		if s.actor.Class != 8 {
			// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
			if s.actor.Equipment.Main.ItemType != 4 && (s.actor.Placement != whatMob.Placement) {
				s.msg.Actor.SendBad("You are too far away to attack.")
				return
			}else if s.actor.Equipment.Main.ItemType == 4 && (s.actor.Placement == whatMob.Placement){
				s.msg.Actor.SendBad("You are too close to attack.")
				return
			}
		}else{
			if s.actor.Placement != whatMob.Placement {
				s.msg.Actor.SendBad("You are too far away to attack.")
				return
			}
		}

		// Lets use a list of attacks,  so we can expand this later if other classes get multi style attacks
		attacks := []float64{
			1.0,
		}

		skillLevel := config.WeaponLevel(s.actor.Skills[5].Value, s.actor.Class)
		if s.actor.Class != 8 {
			skillLevel = config.WeaponLevel(s.actor.Skills[s.actor.Equipment.Main.ItemType].Value, s.actor.Class)
		}
		// Kill is really the fighters realm for specialty..
		if s.actor.Permission.HasAnyFlags(permissions.Fighter) {
			// Did this mofo lethal?
			if config.RollLethal(skillLevel) {
				// Sure did.  Kill this fool and bail.
				s.msg.Actor.SendInfo("You landed a lethal blow on the " + whatMob.Name)
				s.msg.Observers.SendInfo(s.actor.Name + " landed a lethal blow on " + whatMob.Name)
				s.actor.Equipment.DamageWeapon("main", 1)
				whatMob.Stam.Current = 0
				DeathCheck(s, whatMob)
				s.actor.SetTimer("combat", 8)
				return
			}

			if skillLevel >= 4 {
				attacks = append(attacks, .15)
				if skillLevel >= 5 {
					attacks[1] = .3
					if skillLevel >= 6 {
						attacks = append(attacks, .15)
						if skillLevel >= 7 {
							attacks[2] = .30
							if skillLevel >= 8 {
								attacks = append(attacks, .15)
								if skillLevel >= 9 {
									attacks[3] = .3
									if skillLevel >= 10 {
										attacks = append(attacks, .30)
									}
								}
							}
						}
					}
				}
			}
		}

		// Lets start executing the attacks
		weapMsg := ""
		if s.actor.Class != 8 {
			weapMsg = s.actor.Equipment.DamageWeapon("main", 1)
		}
		for _, mult := range attacks {
			// Lets try to crit:
			//TODO: Parry/Miss?
			alwaysCrit := false
			if s.actor.Class != 8 {
				alwaysCrit = s.actor.Equipment.Main.Flags["always_crit"]
			}
			if config.RollCritical(skillLevel) || alwaysCrit {
				mult *= float64(config.CombatModifiers["critical"])
				s.msg.Actor.SendGood("Critical Strike!")
				// TODO: Something something shattered weapons something or other
			} else if config.RollDouble(skillLevel) {
				mult *= float64(config.CombatModifiers["double"])
				s.msg.Actor.SendGood("Double Damage!")
			}
			actualDamage, _ := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * mult)))
			whatMob.AddThreatDamage(actualDamage, s.actor)
			log.Println(strconv.Itoa(whatMob.Stam.Max))
			s.actor.AdvanceSkillExp(int((float64(actualDamage)/float64(whatMob.Stam.Max) * float64(whatMob.Experience))*config.Classes[config.AvailableClasses[s.actor.Class]].WeaponAdvancement))
			s.msg.Actor.SendInfo("You hit the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)

		}
		DeathCheck(s, whatMob)
		if weapMsg != "" {
			s.msg.Actor.SendInfo("weapMsg")
		}
		s.actor.SetTimer("combat", config.CombatCooldown)
		return

	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}

func DeathCheck(s *state, m *objects.Mob) {
	totalExperience := 0
	buildActorString := ""
	if m.Stam.Current <= 0 {
		s.msg.Actor.SendGood("You killed " + m.Name)
		s.msg.Observers.SendGood(s.actor.Name + " killed " + m.Name)
		for k, threat := range m.ThreatTable {
			charClean := s.where.Chars.SearchAll(k)
			if charClean != nil {
				if threat > 0 {
					if m.Level < charClean.Tier {
						totalExperience = int(math.Abs(float64(m.Experience / (6 + (charClean.Tier - m.Level)))))
					} else {
						if threat >= m.Stam.Max/2 {
							totalExperience = m.Experience
						} else if threat >= m.Stam.Max/4 && threat < m.Stam.Max/2 {
							totalExperience = m.Experience/2 + utils.Roll(m.Experience/8, 2, 0)
						} else {
							totalExperience = m.Experience / 8
						}
					}
				} else {
					totalExperience = 0
				}

				if totalExperience == 0 {
					buildActorString += text.Cyan + "You earn no experience for the defeat of the " + m.Name + "\n"
				} else if totalExperience <= m.Experience/8 {
					buildActorString += text.Cyan + "You earn merely " + strconv.Itoa(totalExperience) + " experience for the defeat of the " + m.Name + "\n"
					charClean.Experience.Add(totalExperience)
				} else {
					buildActorString += text.Cyan + "You earn " + strconv.Itoa(totalExperience) + " experience for the defeat of the " + m.Name + "\n"
					charClean.Experience.Add(totalExperience)
				}
				if charClean == s.actor {
					buildActorString += text.Green + m.DropInventory() + "\n"
				}
				if charClean == s.actor {
					s.msg.Actor.Send(buildActorString)
				}else {
					charClean.Write([]byte(buildActorString + "\n" + text.Reset))
				}
				if charClean.Victim == m {
					charClean.Victim = nil
				}
			}
		}

		s.where.Mobs.Remove(m)
	}
}