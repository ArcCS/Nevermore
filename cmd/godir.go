package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
	"strings"
)

func init() {
	addHandler(godir{},
		"Usage:  go direction # \n \n Proceed to the specified exit.   The cardinal directions can also be used without the use of go",
		permissions.Player,
		"N", "NE", "E", "SE", "S", "SW", "W", "NW", "U", "D",
		"NORTH", "NORTHEAST", "EAST", "SOUTHEAST",
		"SOUTH", "SOUTHWEST", "WEST", "NORTHWEST",
		"UP", "DOWN", "GO", "OUT")
}

var (
	directionals = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW", "U", "D", "NORTH", "NORTHEAST",
		"EAST", "SOUTHEAST", "SOUTH", "SOUTHWEST", "WEST", "NORTHWEST", "UP", "DOWN", "OUT"}

	directionIndex = map[string]string{
		"N":         "NORTH",
		"NORTH":     "NORTH",
		"NE":        "NORTHEAST",
		"NORTHEAST": "NORTHEAST",
		"E":         "EAST",
		"EAST":      "EAST",
		"SE":        "SOUTHEAST",
		"SOUTHEAST": "SOUTHEAST",
		"S":         "SOUTH",
		"SOUTH":     "SOUTH",
		"SW":        "SOUTHWEST",
		"SOUTHWEST": "SOUTHWEST",
		"W":         "WEST",
		"WEST":      "WEST",
		"NW":        "NORTHWEST",
		"NORTHWEST": "NORTHWEST",
		"U":         "UP",
		"UP":        "UP",
		"D":         "DOWN",
		"DOWN":      "DOWN",
		"OUT":       "OUT",
	}
)

type godir cmd

func (godir) process(s *state) {

	var exitName string
	from := s.where
	// Does this place even have exits?
	if len(from.Exits) == 0 {
		s.msg.Actor.SendInfo("You can't see anywhere to go from here.")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	// Decide what exit we are going to
	if utils.StringIn(s.cmd, directionals) {
		exitName = directionIndex[s.cmd]
	} else {
		if len(s.words) > 0 {
			// Join the strings together for exits with spaces
			exitName = strings.Join(s.words, " ")
		} else {
			s.msg.Actor.SendBad("Go where?")
		}
	}

	// Test for partial exit names
	exitTxt := strings.ToLower(exitName)
	if !utils.StringIn(strings.ToUpper(exitTxt), directionals) {
		for txtE := range from.Exits {
			if strings.Contains(txtE, exitTxt) {
				exitTxt = txtE
			}
		}
	}
	if toE, ok := from.Exits[exitTxt]; ok {
		s.actor.RunHook("move")
		// Check that the room ID exists
		if to, ok := objects.Rooms[toE.ToId]; ok {
			// Apply a lock
			if !utils.IntIn(toE.ToId, s.rLocks) {
				s.AddLocks(toE.ToId)
				s.ok = false
				return
			} else {
				if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {

					// Check some timers
					ready, msg := s.actor.TimerReady("evade")
					if !ready {
						s.msg.Actor.SendBad(msg)
						return
					}

					if !objects.Rooms[toE.ToId].Flags["active"] {
						s.msg.Actor.SendBad("Go where?")
						return
					}

					if toE.Flags["invisible"] && !s.actor.CheckFlag("detectj-invisible") {
						s.msg.Actor.SendBad("Go where?")
						return
					}

					if toE.Flags["placement_dependent"] && s.actor.Placement != toE.Placement {
						s.msg.Actor.SendBad("You must be next to the exit to use it.")
						return
					}

					if toE.Flags["closed"] {
						s.msg.Actor.SendBad("The way is closed.")
						return
					}

					if toE.Flags["day_only"] && !objects.DayTime {
						s.msg.Actor.SendBad("You can only go there at night.")
						return
					}

					if toE.Flags["night_only"] && objects.DayTime {
						s.msg.Actor.SendBad("You can only go there during the day.")
						return
					}

					if s.actor.Equipment.GetWeight() > s.actor.MaxWeight() {
						s.msg.Actor.SendBad("You are carrying too much to move.")
						return
					}

					hasRope := false
					if s.actor.Equipment.Off != (*objects.Item)(nil) {
						if s.actor.Equipment.Off.ItemId == 1463 {
							hasRope = true
						}
					}
					if toE.Flags["levitate"] && !s.actor.CheckFlag("levitate") && !hasRope {
						chanceToPass := s.actor.GetStat("dex")/45 + 10
						if utils.Roll(100, 1, 0) >= chanceToPass {
							fallDamageStam := int(config.FallDamage*float64(s.actor.Stam.Max)) -
								(config.ConFallDamageMod * s.actor.GetStat("con")) -
								(config.DexFallDamageMod * s.actor.GetStat("dex"))
							fallDamageVit := int(config.FallDamage*float64(s.actor.Stam.Max)) -
								(config.ConFallDamageMod * s.actor.GetStat("con")) -
								(config.DexFallDamageMod * s.actor.GetStat("dex"))
							totStam, totVit := 0, 0
							if fallDamageStam > 0 {
								totStam, totVit = s.actor.ReceiveDamageNoArmor(fallDamageStam)
							}
							if fallDamageVit > 0 {
								totVit += s.actor.ReceiveVitalDamageNoArmor(fallDamageVit)
							}
							buildStr := ""
							if totStam <= 0 && totVit <= 0 {
								buildStr = "You take no damage in the fall."
							} else {
								if totStam >= 1 {
									buildStr += "You take " + strconv.Itoa(totStam) + " points of stamina"
								}
								if totVit >= 1 {
									if totStam >= 1 {
										buildStr += " and "
									}
									buildStr += strconv.Itoa(totVit) + " points of vitality"
								}
								buildStr += " damage in the fall."
							}
							s.msg.Actor.Send("You fall while trying to go that way! " + buildStr)
							go s.actor.DeathCheck("fell to their death.")
							return
						}
					}

					if objects.Rooms[toE.ToId].Crowded() {
						s.msg.Actor.SendInfo("That area is crowded.")
						s.ok = true
						return
					}

					evasiveMan := 0
					followList := make([]*objects.Mob, 0)
					// Check if anyone blocks.
					for _, mob := range s.where.Mobs.Contents {
						// Check if a mob blocks.
						if _, inList := mob.ThreatTable[s.actor.Name]; inList {
							if mob.CheckFlag("block_exit") && mob.Placement == s.actor.Placement && mob.MobStunned == 0 && !mob.CheckFlag("run_away") {
								curChance := config.MobBlock - ((s.actor.Tier - mob.Level) * config.MobBlockPerLevel)
								if curChance > 85 {
									curChance = 85
								}
								if utils.Roll(100, 1, 0) <= curChance {
									s.msg.Actor.SendBad(mob.Name + " blocks your way.")
									s.actor.SetTimer("global", 8)
									return
								}
							}
							if mob.CurrentTarget == s.actor.Name {
								// Now check if they follow.
								if mob.CheckFlag("follows") && !mob.CheckFlag("curious_canticle") {
									followList = append(followList, mob)
								}
								evasiveMan = 2
								if mob.Placement == s.actor.Placement {
									evasiveMan = 4
								}
							}
						}
					}
					from.Chars.Remove(s.actor)
					// If they were evasive, add a global timer
					s.actor.SetTimer("evade", evasiveMan)
					to.Chars.Add(s.actor)
					s.actor.Victim = nil
					s.actor.Placement = 3
					s.actor.ParentId = toE.ToId

					// Broadcast leaving and arrival notifications
					if s.actor.Flags["invisible"] == false {
						s.msg.Observers[from.RoomId].SendInfo("You see ", s.actor.Name, " go to the ", strings.ToLower(toE.Name), ".")
						s.msg.Observers[to.RoomId].SendInfo(s.actor.Name, " just arrived.")
					}

					// Character has been removed, invoke any follows for them..  this should be fine as the mob should take over locks
					for _, mob := range followList {
						mobProc := mob
						go func() { mobProc.MobCommands <- "follow " + s.actor.Name }()
					}

					// Do not invoke player state, just move them within this state lock
					if len(s.actor.PartyFollowers) > 0 {
						for _, peo := range s.actor.PartyFollowers {
							follChar := s.where.Chars.SearchAll(peo)
							endFollProc := false
							if follChar != nil {
								// Check some timers
								if !follChar.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
									ready, msg := follChar.TimerReady("evade")
									if !ready {
										follChar.Write([]byte(text.Bad + msg))
										break
									}

									if s.actor.Stam.Current <= 0 {
										follChar.Write([]byte(text.Bad + "You are far too tired to follow."))
										break
									}

									follChar.RunHook("move")

									evasiveMan = 0
									followList = make([]*objects.Mob, 0)

									if !objects.Rooms[toE.ToId].Flags["active"] {
										follChar.Write([]byte(text.Bad + "Go where?"))
										break
									}

									if toE.Flags["invisible"] && !follChar.CheckFlag("detect-invisible") {
										follChar.Write([]byte(text.Bad + "Go where?"))
										break
									}

									if toE.Flags["placement_dependent"] && follChar.Placement != toE.Placement {
										follChar.Write([]byte(text.Bad + "You must be next to the exit to use it."))
										break
									}

									if follChar.Equipment.GetWeight() > follChar.MaxWeight() {
										follChar.Write([]byte(text.Bad + "You are carrying too much to move."))
										break
									}

									hasRope := false
									if follChar.Equipment.Off != (*objects.Item)(nil) {
										if follChar.Equipment.Off.ItemId == 1463 {
											hasRope = true
										}
									}

									if toE.Flags["levitate"] && !follChar.CheckFlag("levitate") && !hasRope {
										chanceToPass := follChar.GetStat("dex")/45 + 10
										if utils.Roll(100, 1, 0) >= chanceToPass {
											fallDamageStam := int(config.FallDamage*float64(follChar.Stam.Max)) -
												(config.ConFallDamageMod * follChar.GetStat("con")) -
												(config.DexFallDamageMod * follChar.GetStat("dex"))
											fallDamageVit := int(config.FallDamage*float64(follChar.Stam.Max)) -
												(config.ConFallDamageMod * follChar.GetStat("con")) -
												(config.DexFallDamageMod * follChar.GetStat("dex"))
											totStam, totVit := 0, 0
											if fallDamageStam > 0 {
												totStam, totVit = follChar.ReceiveDamageNoArmor(fallDamageStam)
											}
											if fallDamageVit > 0 {
												totVit += follChar.ReceiveVitalDamageNoArmor(fallDamageVit)
											}
											buildStr := ""
											if totStam <= 0 && totVit <= 0 {
												buildStr = "You take no damage in the fall."
											} else {
												if totStam >= 1 {
													buildStr += "You take " + strconv.Itoa(totStam) + " points of stamina"
												}
												if totVit >= 1 {
													if totStam >= 1 {
														buildStr += " and "
													}
													buildStr += strconv.Itoa(totVit) + " points of vitality"
												}
												buildStr += " damage in the fall."
											}
											follChar.Write([]byte(text.Bad + "You fall while trying to go that way! " + buildStr))
											go follChar.DeathCheck("fell to their death.")
											break
										}
									}

									// Check if anyone blocks.
									for _, mob := range s.where.Mobs.Contents {
										// Check if a mob blocks.
										if _, inList := mob.ThreatTable[follChar.Name]; inList {
											if mob.CheckFlag("block_exit") && mob.Placement == follChar.Placement && mob.MobStunned == 0 && !mob.CheckFlag("run_away") {
												curChance := config.MobBlock - ((follChar.Tier - mob.Level) * config.MobBlockPerLevel)
												if curChance > 85 {
													curChance = 85
												}
												if utils.Roll(100, 1, 0) <= curChance {
													follChar.Write([]byte(mob.Name + " blocks you from following." + "\n"))
													follChar.SetTimer("global", 8)
													endFollProc = true
													break
												}
											}
											if mob.CurrentTarget == follChar.Name {
												// Now check if they follow.
												if mob.CheckFlag("follows") && !mob.CheckFlag("curious_canticle") {
													followList = append(followList, mob)
												}
												evasiveMan = 2
												if mob.Placement == follChar.Placement {
													evasiveMan = 4
												}
											}
										}
									}
									if endFollProc {
										continue
									}
								}
								from.Chars.Remove(follChar)
								// If they were evasive, add a global timer
								follChar.SetTimer("evade", evasiveMan)
								to.Chars.Add(follChar)
								follChar.Victim = nil
								follChar.Placement = 3
								follChar.ParentId = toE.ToId

								if s.actor.CheckFlag("blind") {
									s.msg.Actor.SendBad("You can't see anything!")
									return
								} else {
									follChar.Write([]byte(objects.Rooms[to.RoomId].Look(follChar)))
								}

								// Broadcast leaving and arrival notifications
								if follChar.Flags["invisible"] == false {
									s.msg.Observers[from.RoomId].SendInfo("You see ", follChar.Name, " follow "+s.actor.Name+" to the ", strings.ToLower(toE.Name), ".")
									s.msg.Observers[to.RoomId].SendInfo(follChar.Name, " just arrived.")
								}

								// Character has been removed, invoke any follows
								for _, mob := range followList {
									mobProc := mob
									go func() { mobProc.MobCommands <- "follow " + follChar.Name }()
								}
							}
						}
					}

					s.scriptActor("LOOK")
					s.ok = true
					return
				} else {
					from.Chars.Remove(s.actor)
					to.Chars.Add(s.actor)
					s.actor.Victim = nil
					s.actor.Placement = 3
					s.actor.ParentId = toE.ToId

					// Broadcast leaving and arrival notifications
					if s.actor.Flags["invisible"] == false {
						s.msg.Observers[from.RoomId].SendInfo("You see ", s.actor.Name, " go to the ", strings.ToLower(toE.Name), ".")
						s.msg.Observers[to.RoomId].SendInfo(s.actor.Name, " just arrived.")
					}

					s.scriptActor("LOOK")
					s.ok = true
					return
				}
			}
		} else {
			s.msg.Actor.SendInfo("You can't go that direction.")
			s.ok = true
			return
		}
	} else {
		s.msg.Actor.SendInfo("You can't go that direction.")
		s.ok = true
		return
	}
}
