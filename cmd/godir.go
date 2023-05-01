package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
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
			if !utils.IntIn(toE.ToId, s.cLocks) {
				s.AddCharLock(toE.ToId)
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

					if toE.Flags["invisible"] && !s.actor.CheckFlag("detect_invisible") {
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

					if s.actor.Equipment.Weight > s.actor.MaxWeight() {
						s.msg.Actor.SendBad("You are carrying too much to move.")
						return
					}

					if toE.Flags["levitate"] && !s.actor.CheckFlag("levitate") {
						s.msg.Actor.Send("You fall while trying to go that way!  You take 20 points of damage!")
						s.actor.ReceiveDamage(20)
						return
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

					if len(s.actor.PartyFollowers) > 0 {
						for _, party := range s.actor.PartyFollowers {
							if party.ParentId == s.where.RoomId {
								go Script(party, s.cmd+" "+strings.Join(s.input, " "))
							}
						}
					}

					// Character has been removed, invoke any follows
					for _, mob := range followList {
						go func() { mob.MobCommands <- "follow " + s.actor.Name }()
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

					if len(s.actor.PartyFollowers) > 0 {
						for _, party := range s.actor.PartyFollowers {
							if party.ParentId == s.where.RoomId {
								go func() { party.CharCommands <- "go " + exitTxt }()
							}
						}
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
