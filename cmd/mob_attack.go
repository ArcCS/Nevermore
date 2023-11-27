package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(mobAttack{},
		"",
		permissions.None,
		"$MATTACK")
}

type mobAttack cmd

func (mobAttack) process(s *state) {
	/*
		// Am I hostile?  Should I pick a target?
					if m.CurrentTarget == "" && m.Flags["hostile"] && len(Rooms[m.ParentId].Chars.MobList(m)) > 0 {
						for m.CurrentTarget == "" {
							for i := 0; i < 4; i++ {
								if m.CurrentTarget != "" {
									break
								}
								potentials := Rooms[m.ParentId].Chars.MobListAt(m, i)
								if len(potentials) > 0 {
									for _, potential := range potentials {
										if utils.Roll(100, 1, 0) <= config.ProximityChance-(i*config.ProximityStep) {
											m.AddThreatDamage(1, Rooms[m.ParentId].Chars.MobSearch(potential, m))
											Rooms[m.ParentId].MessageAll(m.Name + " attacks " + m.CurrentTarget + text.Reset + "\n")
											break
										}
									}
								}
							}
						}
					}

					if m.CurrentTarget != "" {
						if Rooms[m.ParentId].Chars.SearchAll(m.CurrentTarget) == nil {
							m.CurrentTarget = ""
						}
					}

					// Do I want to change targets? 33% chance if the current target isn't the highest on the threat table
					if len(m.ThreatTable) > 1 {
						rankedThreats := utils.RankMapStringInt(m.ThreatTable)
						if m.CurrentTarget != rankedThreats[0] {
							if utils.Roll(100, 1, 0) <= 5 {
								if utils.StringIn(rankedThreats[0], Rooms[m.ParentId].Chars.MobList(m)) {
									m.CurrentTarget = rankedThreats[0]
									Rooms[m.ParentId].MessageAll(m.Name + " turns to " + m.CurrentTarget + "\n" + text.Reset)
								}
							}
						}
					}

					if m.CurrentTarget == "" && m.Placement != 3 && !m.CheckFlag("immobile") {
						oldPlacement := m.Placement
						if m.Placement > 3 {
							m.Placement--
						} else {
							m.Placement++
						}
						if !m.Flags["hidden"] {
							whichNumber := Rooms[m.ParentId].Mobs.GetNumber(m)
							if len(Rooms[m.ParentId].Mobs.Contents) > 1 && whichNumber > 1 {
								Rooms[m.ParentId].MessageMovement(oldPlacement, m.Placement, m.Name+" #"+strconv.Itoa(whichNumber))
							} else {
								Rooms[m.ParentId].MessageMovement(oldPlacement, m.Placement, m.Name)
							}
						}
						return
					}

					if m.CurrentTarget != "" && m.BreathWeapon != "" &&
						(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m).Placement)) == 1) {

						// Roll to see if we're going to breathe
						if utils.Roll(100, 1, 0) <= 30 {
							target := Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m)
							var targets []*Character
							for _, character := range Rooms[m.ParentId].Chars.Contents {
								if character.Placement == target.Placement && !character.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
									log.Println("Adding target: ", character.Name, " to breath list")
									targets = append(targets, character)
								}
							}

							Rooms[m.ParentId].MessageAll("The " + m.Name + " breathes " + m.BreathWeapon + " at " + target.Name + "\n")
							damageTotal := config.BreatheDamage(m.Level)
							reflectDamage := 0
							for _, t := range targets {
								if utils.StringIn(m.BreathWeapon, []string{"fire", "air", "earth", "water"}) {
									t.RunHook("attacked")
									stamDam, vitDam, resisted := t.ReceiveMagicDamage(damageTotal, m.BreathWeapon)
									if _, err := t.Write([]byte(text.Bad + m.Name + "'s breath  struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality. You resisted " + strconv.Itoa(resisted) + "damage." + text.Reset + "\n")); err != nil {
										log.Println("Error writing to player:", err)
									}
									if target.CheckFlag("reflection") {
										reflectDamage = int(float64(damageTotal) * (float64(target.GetStat("int")) * config.ReflectDamagePerInt))
										m.ReceiveDamage(reflectDamage)
										if _, err := target.Write([]byte(text.Cyan + "You reflect " + strconv.Itoa(reflectDamage) + " damage back to " + m.Name + "!\n" + text.Reset)); err != nil {
											log.Println("Error writing to player:", err)
										}
										m.DeathCheck(target)
									}
									t.DeathCheck("was slain by a " + m.Name + ".")
								} else if m.BreathWeapon == "paralytic" {
									if _, err := t.Write([]byte(text.Gray + m.Name + " breathes paralytic gas on to you.\n")); err != nil {
										log.Println("Error writing to player:", err)
									}
									target.SetTimer("global", 24)
								} else if m.BreathWeapon == "pestilence" {
									if _, err := t.Write([]byte(text.Gray + m.Name + " breathes infectious gas on to you.\n")); err != nil {
										log.Println("Error writing to player:", err)
									}
									Effects["disease"](m, target, m.Level)
								}
							}
							return
						}
					}

					// Calculate Vital/Crit/Double
					multiplier := float64(1)
					vitalStrike := false
					criticalStrike := false
					doubleDamage := false
					penalty := 1

					if m.CurrentTarget != "" && m.Flags["ranged_attack"] &&
						(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m).Placement)) >= 1) {
						target := Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m)
						missChance := 0
						lvlDiff := target.Tier - m.Level
						if lvlDiff >= 1 {
							missChance += lvlDiff * config.MissPerLevel
						}
						missChance += target.GetStat("dex") * config.MissPerDex
						if utils.Roll(100, 1, 0) <= missChance {
							if _, err := target.Write([]byte(text.Green + m.Name + " missed you!!" + "\n" + text.Reset)); err != nil {
								log.Println("Error writing to player:", err)
							}
							data.StoreCombatMetric("range-miss", 0, 1, 0, 0, 0, 1, m.MobId, m.Level, 0, target.CharId)
							return
						}
						// If we made it here, default out and do a range hit.
						stamDamage := 0
						vitDamage := 0
						resisted := 0
						reflectDamage := 0
						actualDamage := m.InflictDamage()
						if utils.Roll(10, 1, 0) <= penalty {
							attackStyleRoll := utils.Roll(10, 1, 0)
							if attackStyleRoll <= config.MobVital {
								multiplier = 2
								vitalStrike = true
							} else if attackStyleRoll <= config.MobCritical {
								multiplier = 4
								criticalStrike = true
							} else if attackStyleRoll <= config.MobDouble {
								multiplier = 2
								doubleDamage = true
							}
						}
						if vitalStrike {
							vitDamage, resisted = target.ReceiveVitalDamage(int(math.Ceil(float64(actualDamage) * multiplier)))
							data.StoreCombatMetric("range_vital", 0, 1, int(math.Ceil(float64(actualDamage)*multiplier)), resisted, vitDamage, 1, m.MobId, m.Level, 0, target.CharId)
							if _, err := target.Write([]byte(text.Red + "Vital Strike!!!\n" + text.Reset)); err != nil {
								log.Println("Error writing to player:", err)
							}
						} else {
							stamDamage, vitDamage, resisted = target.ReceiveDamage(int(math.Ceil(float64(actualDamage) * multiplier)))
							data.StoreCombatMetric("range", 0, 1, int(math.Ceil(float64(actualDamage)*multiplier)), resisted, vitDamage, 1, m.MobId, m.Level, 0, target.CharId)
						}

						buildString := ""
						if stamDamage != 0 {
							buildString += strconv.Itoa(stamDamage) + " stamina"
						}
						if stamDamage != 0 && vitDamage != 0 {
							buildString += " and "
						}
						if vitDamage != 0 {
							buildString += strconv.Itoa(vitDamage) + " vitality"
						}
						if criticalStrike {
							if _, err := target.Write([]byte(text.Red + "Critical Strike!!!\n" + text.Reset)); err != nil {
								log.Println("Error writing to player:", err)
							}
						}
						if doubleDamage {
							if _, err := target.Write([]byte(text.Red + "Double Damage!!!\n" + text.Reset)); err != nil {
								log.Println("Error writing to player:", err)
							}
						}
						if _, err := target.Write([]byte(text.Red + "Thwwip!! " + m.Name + " attacks you for " + buildString + " points of damage!" + "\n" + text.Reset)); err != nil {
							log.Println("Error writing to player:", err)
						}
						if target.CheckFlag("reflection") {
							reflectDamage = int(float64(actualDamage) * (float64(target.GetStat("int")) * config.ReflectDamagePerInt))
							mobFin, _, mobResisted := m.ReceiveDamage(reflectDamage)
							data.StoreCombatMetric("range_player_reflect", 0, 1, reflectDamage, mobResisted, mobFin, 0, target.CharId, target.Tier, 1, m.MobId)
							if _, err := target.Write([]byte(text.Cyan + "You reflect " + strconv.Itoa(reflectDamage) + " damage back to " + m.Name + "!\n" + text.Reset)); err != nil {
								log.Println("Error writing to player:", err)
							}
							m.DeathCheck(target)
						}
						target.RunHook("attacked")
						target.DeathCheck("was slain by a " + m.Name + ".")
						return
					}

					if (m.CurrentTarget != "" && !m.CheckFlag("immobile") &&
						m.Placement != Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m).Placement) ||
						(m.CurrentTarget != "" &&
							(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m).Placement)) > 1)) {
						oldPlacement := m.Placement
						if m.Placement > Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m).Placement {
							m.Placement--
						} else {
							m.Placement++
						}
						if !m.Flags["hidden"] {
							whichNumber := Rooms[m.ParentId].Mobs.GetNumber(m)
							Rooms[m.ParentId].MessageMovement(oldPlacement, m.Placement, m.Name+" #"+strconv.Itoa(whichNumber))
						}
						// Next to attack
					} else if m.CurrentTarget != "" &&
						m.Placement == Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m).Placement {
						// Check to see if the mob misses:
						// Am I against a fighter, and they succeed in a parry roll?
						target := Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m)
						missChance := 0
						lvlDiff := target.Tier - m.Level
						if lvlDiff >= 1 {
							missChance += lvlDiff * config.MissPerLevel
						}
						missChance += target.GetStat("dex") * config.HitPerDex
						if utils.Roll(100, 1, 0) <= missChance {
							if _, err := target.Write([]byte(text.Green + m.Name + " missed you!!" + "\n" + text.Reset)); err != nil {
								log.Println("Error writing to player:", err)
							}
							data.StoreCombatMetric("melee-miss", 0, 1, 0, 0, 0, 1, m.MobId, m.Level, 0, target.CharId)
							return
						}
						target.RunHook("attacked")
						m.CheckForExtraAttack(target)
						if target.Class == 0 && target.Equipment.Main != nil && config.RollParry(config.WeaponLevel(target.Skills[target.Equipment.Main.ItemType].Value, target.Class)) {
							if target.Tier >= config.SpecialAbilityTier {
								// It's a riposte
								actualDamage, _, resisted := m.ReceiveDamage(int(math.Ceil(float64(target.InflictDamage()))))
								data.StoreCombatMetric("melee_player_riposte", 0, 1, actualDamage+resisted, resisted, actualDamage, 0, target.CharId, target.Tier, 1, m.MobId)
								if _, err := target.Write([]byte(text.Green + "You parry and riposte the attack from " + m.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + "\n" + text.Reset)); err != nil {
									log.Println("Error writing to player:", err)
								}
								if m.DeathCheck(target) {
									return
								}
								m.Stun(config.ParryStuns * 8)
							} else {
								if _, err := target.Write([]byte(text.Green + "You parry the attack from " + m.Name + "\n" + text.Reset)); err != nil {
									log.Println("Error writing to player:", err)
								}
								m.Stun(config.ParryStuns * 8)
							}
						} else {
							stamDamage := 0
							vitDamage := 0
							resisted := 0
							actualDamage := m.InflictDamage()
							reflectDamage := 0
							if utils.Roll(10, 1, 0) <= penalty {
								attackStyleRoll := utils.Roll(10, 1, 0)
								if attackStyleRoll <= config.MobVital {
									multiplier = 2
									vitalStrike = true
								} else if attackStyleRoll <= config.MobCritical {
									multiplier = 4
									criticalStrike = true
								} else if attackStyleRoll <= config.MobDouble {
									multiplier = 2
									doubleDamage = true
								}
							}
							if vitalStrike {
								vitDamage, resisted = target.ReceiveVitalDamage(int(math.Ceil(float64(actualDamage) * multiplier)))
								data.StoreCombatMetric("melee_vital", 0, 1, int(math.Ceil(float64(actualDamage)*multiplier)), resisted, vitDamage, 1, m.MobId, m.Level, 0, target.CharId)
								if _, err := target.Write([]byte(text.Red + "Vital Strike!!!\n" + text.Reset)); err != nil {
									log.Println("Error writing to player:", err)
								}
							} else {
								stamDamage, vitDamage, resisted = target.ReceiveDamage(int(math.Ceil(float64(actualDamage) * multiplier)))
								data.StoreCombatMetric("melee", 0, 1, int(math.Ceil(float64(actualDamage)*multiplier)), resisted, stamDamage+vitDamage, 1, m.MobId, m.Level, 0, target.CharId)

							}
							buildString := ""
							if stamDamage != 0 {
								buildString += strconv.Itoa(stamDamage) + " stamina"
							}
							if stamDamage != 0 && vitDamage != 0 {
								buildString += " and "
							}
							if vitDamage != 0 {
								buildString += strconv.Itoa(vitDamage) + " vitality"
							}
							if stamDamage == 0 && vitDamage == 0 {
								if _, err := target.Write([]byte(text.Red + m.Name + " attacks bounces off of you for no damage!" + "\n" + text.Reset)); err != nil {
									log.Println("Error writing to player:", err)
								}
							} else {
								if criticalStrike {
									if _, err := target.Write([]byte(text.Red + "Critical Strike!!!\n" + text.Reset)); err != nil {
										log.Println("Error writing to player:", err)
									}
								}
								if doubleDamage {
									if _, err := target.Write([]byte(text.Red + "Double Damage!!!\n" + text.Reset)); err != nil {
										log.Println("Error writing to player:", err)
									}
								}
								if _, err := target.Write([]byte(text.Red + m.Name + " attacks you for " + buildString + " points of damage!" + "\n" + text.Reset)); err != nil {
									log.Println("Error writing to player:", err)
								}
							}
							if target.CheckFlag("reflection") {
								reflectDamage = int(float64(actualDamage) * (float64(target.GetStat("int")) * config.ReflectDamagePerInt))
								mobFin, _, mobResisted := m.ReceiveDamage(reflectDamage)
								data.StoreCombatMetric("melee_player_reflect", 0, 1, reflectDamage, mobResisted, mobFin, 0, target.CharId, target.Tier, 1, m.MobId)
								if _, err := target.Write([]byte(text.Cyan + "You reflect " + strconv.Itoa(reflectDamage) + " damage back to " + m.Name + "!\n" + text.Reset)); err != nil {
									log.Println("Error writing to player:", err)
								}
								m.DeathCheck(target)
							}
							target.DeathCheck("was slain by a " + m.Name + ".")
						}
					}
	*/
	return
}
