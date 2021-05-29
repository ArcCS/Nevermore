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
				return
			} else {
				if !toE.Flags["placement_dependent"] {
					//  Next room needs to be active
					if !objects.Rooms[toE.ToId].Flags["active"] && !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
						s.msg.Actor.SendBad("Go where?")
						return
					}
					//TODO: Night/Day/Level Checks

					//log.Println(rooms.Rooms[toE.ToId].Crowded())
					if !objects.Rooms[toE.ToId].Crowded() {
						for _ , mob := range s.where.Mobs.Contents {
							// Check if a mob blocks.
							if mob.Flags["block_exit"] {
								curChance := config.MobBlock - ((s.actor.Tier - mob.Level)*config.MobBlockPerLevel)
								if curChance > 85 {
									curChance = 85
								}
								if utils.Roll(100, 1, 0) <= curChance{
									s.msg.Actor.SendBad(mob.Name + " blocks your way.")
									return
								}
							}
							if mob.Flags["follows"] {
								curChance := config.MobFollow - ((s.actor.Tier - mob.Level)*config.MobFollowPerLevel)
								if curChance > 85 {
									curChance = 85
								}
								if utils.Roll(100, 1, 0) <= curChance{
									s.msg.Actor.SendBad(mob.Name + " follows you!")
									from.Mobs.Remove(mob)
									//TODO Calculate Vital?
									to.Mobs.Add(mob)
									mob.CurrentTarget = s.actor.Name
								}
							}
						}
						from.Chars.Remove(s.actor)
						to.Chars.Add(s.actor)
						s.actor.Placement = 3
						s.actor.ParentId = toE.ToId
						// Broadcast leaving and arrival notifications
						if s.actor.Flags["invisible"] == false {
							s.msg.Observers[from.RoomId].SendInfo("You see ", s.actor.Name, " go to the ", strings.ToLower(to.Name), ".")
							s.msg.Observers[to.RoomId].SendInfo(s.actor.Name, " just arrived.")
						}
						s.scriptActor("LOOK")
						s.ok = true
						return
					} else {
						s.msg.Actor.SendInfo("That area is crowded.")
						s.ok = true
						return
					}
				}

			}
		}
	} else {
		s.msg.Actor.SendInfo("You can't go that direction.")
		s.ok = true
		return
	}

}
