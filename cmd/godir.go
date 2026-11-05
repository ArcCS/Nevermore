package cmd

import (
	"log"
	"strconv"
	"strings"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/message"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(godir{},
		"Usage:  go direction # \n \n Proceed to the specified exit.   The cardinal directions can also be used without the use of go",
		permissions.Player,
		"GO", "N", "NE", "E", "SE", "S", "SW", "W", "NW", "U", "D",
		"NORTH", "NORTHEAST", "EAST", "SOUTHEAST",
		"SOUTH", "SOUTHWEST", "WEST", "NORTHWEST",
		"UP", "DOWN", "OUT", "O")
}

var (
	directionals = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW", "U", "D", "NORTH", "NORTHEAST",
		"EAST", "SOUTHEAST", "SOUTH", "SOUTHWEST", "WEST", "NORTHWEST", "UP", "DOWN", "OUT", "O"}

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
		"O":         "OUT",
	}
)

type godir cmd

func (godir) process(s *state) {

	// Moving ends the post-resume death grace, whether or not the move
	// succeeds - a player picking an exit has seen the screen.
	s.actor.ClearResumeGrace()

	var exitName string
	from := s.where
	// Does this place even have exits?
	if len(from.Exits) == 0 {
		s.msg.Actor.SendInfo("You can't see anywhere to go from here.")
		return
	}

	if !s.requireStamina() {
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
		toId := toE.ToFor(s.actor)
		if to, ok := objects.Rooms[toId]; ok {
			// Apply a lock
			if !utils.IntIn(toId, s.rLocks) {
				s.AddLocks(toId)
				s.ok = false
				return
			} else {
				if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {

					allowed, evasive, chasers := canMove(actorMover(s), from, to, toE)
					if !allowed {
						return
					}

					moveChar(s.actor, from, to, evasive)

					// Broadcast leaving and arrival notifications
					if s.actor.Flags["invisible"] == false {
						s.msg.Observers[from.RoomId].SendInfo("You see ", s.actor.Name, " go to the ", strings.ToLower(toE.Name), ".")
						s.msg.Observers[to.RoomId].SendInfo(s.actor.Name, " just arrived.")
					}

					// The actor is out of the room, so anything hunting them can
					// come along. This is done here, inline, because this state
					// already holds both rooms: handing the mob a command to chase
					// them later would have it move itself with nothing locked,
					// after the actor may have moved on or died. FollowChar makes
					// its own roll and refuses rooms mobs have no business in.
					// Only a mob that actually arrives gets a swing, and only the
					// first one to land it does - a chased character pays for the
					// exit once, however big the pack behind them.
					// MOB-CHASE-SWITCH: delete this loop to stop mobs following the actor between rooms.
					vitalSpent := false
					for _, mob := range chasers {
						if !mob.FollowChar(s.actor, from, to) {
							continue
						}
						if vitalSpent {
							continue
						}
						landed, died := followVital(actorMover(s), mob)
						vitalSpent = landed
						if died {
							// The death script has them now, so there is no one
							// left here to lead the party or be shown the room.
							s.ok = true
							return
						}
					}

					// Do not invoke player state, just move them within this state
					// lock - a state of their own would deadlock on the rooms this
					// one is holding. Each follower is judged on their own merits,
					// so one who cannot come along does not strand the rest.
					var arrived []*objects.Character
					for _, peo := range s.actor.PartyFollowers {
						follChar := from.Chars.SearchAll(peo)
						if follChar == nil {
							continue
						}

						follChar.RunHook("move")

						follAllowed, follEvasive, follChasers := true, 0, []*objects.Mob(nil)
						if !follChar.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
							follAllowed, follEvasive, follChasers = canMove(followerMover(follChar), from, to, toE)
						}
						if !follAllowed {
							continue
						}

						moveChar(follChar, from, to, follEvasive)

						// Broadcast leaving and arrival notifications
						if follChar.Flags["invisible"] == false {
							s.msg.Observers[from.RoomId].SendInfo("You see ", follChar.Name, " follow "+s.actor.Name+" to the ", strings.ToLower(toE.Name), ".")
							s.msg.Observers[to.RoomId].SendInfo(follChar.Name, " just arrived.")
						}

						// MOB-CHASE-SWITCH: delete this loop to stop mobs following party followers between rooms.
						follVitalSpent, follDied := false, false
						for _, mob := range follChasers {
							if !mob.FollowChar(follChar, from, to) {
								continue
							}
							if follVitalSpent {
								continue
							}
							landed, died := followVital(followerMover(follChar), mob)
							follVitalSpent = landed
							if died {
								// The death script has them; there is no one to show the room to.
								follDied = true
								break
							}
						}
						if !follDied {
							arrived = append(arrived, follChar)
						}
					}

					// Show the room to everyone who made it, once the whole
					// party is in it, so a follower sees the same thing their
					// leader does: the room, and everyone standing in it.
					for _, follChar := range arrived {
						showRoom(follChar, to)
					}

					s.scriptActor("LOOK")
					s.ok = true
					return
				} else {
					// Builders and staff walk through everything unchecked.
					// SetTimer is a no-op for them, so the evade timer here
					// costs them nothing.
					moveChar(s.actor, from, to, 0)

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

// mover is the character being walked through an exit, together with how to
// talk to them. The actor's output goes through the state's buffers so it
// arrives with the rest of the command; a party follower is not this state's
// actor, so theirs has to be written straight to their connection. Everything
// else about moving the two is identical, which is the point.
type mover struct {
	char *objects.Character
	bad  func(string)
	info func(string)
}

func actorMover(s *state) mover {
	return mover{
		char: s.actor,
		bad:  func(msg string) { s.msg.Actor.SendBad(msg) },
		info: func(msg string) { s.msg.Actor.SendInfo(msg) },
	}
}

func followerMover(char *objects.Character) mover {
	return mover{
		char: char,
		bad:  func(msg string) { writeTo(char, text.Bad+msg) },
		info: func(msg string) { writeTo(char, text.Info+msg) },
	}
}

// showRoom shows a character who is not this state's actor the room they have
// just been led into, exactly as LOOK would. The view is built in a buffer of
// their own and delivered in one write, as the actor's is, rather than a line
// at a time with a prompt after each.
func showRoom(char *objects.Character, room *objects.Room) {
	b := message.AcquireBuffer()
	lookRoom(char, room, b)
	b.Deliver(char)
	message.ReleaseBuffer(b)
}

// writeTo sends text straight to a character, for output that cannot go
// through the state's buffers because the character is not this state's actor.
func writeTo(char *objects.Character, msg string) {
	if _, err := char.Write([]byte(msg + "\n")); err != nil {
		log.Println("Error writing to player: ", err)
	}
}

// canMove runs every per-character check for walking through toE and reports
// whether the character may go, the evade timer their exit earned them, and
// the mobs entitled to give chase once they are gone. One of the checks bites:
// a failed levitate roll damages the character before refusing, and can kill
// them.
//
// Nothing here ends the command or touches the state. The caller decides what
// a refusal means - the actor stops walking, a follower is simply left behind -
// which is what lets the actor and their party share one copy of the rules.
// Both rooms must already be locked by the caller.
func canMove(m mover, from *objects.Room, to *objects.Room, toE *objects.Exit) (allowed bool, evasive int, chasers []*objects.Mob) {
	char := m.char

	if rooted, msg := char.Rooted(); rooted {
		m.bad(msg)
		return false, 0, nil
	}

	// Check some timers
	if ready, msg := char.TimerReady("evade"); !ready {
		m.bad(msg)
		return false, 0, nil
	}

	if char.Stam.Current <= 0 {
		m.bad("You are far too tired to move.")
		return false, 0, nil
	}

	if !to.Flags["active"] {
		m.bad("Go where?")
		return false, 0, nil
	}

	if toE.Flags["invisible"] && !char.CheckFlag("detect-invisible") {
		m.bad("Go where?")
		return false, 0, nil
	}

	if toE.Flags["placement_dependent"] && char.Placement != toE.Placement {
		m.bad("You must be next to the exit to use it.")
		return false, 0, nil
	}

	if toE.Flags["closed"] {
		m.bad("The way is closed.")
		return false, 0, nil
	}

	if objects.QuestExitSealed(from.RoomId, toE.Name) {
		m.bad("The way is sealed.")
		return false, 0, nil
	}

	if toE.Flags["day_only"] && !objects.DayTime {
		m.bad("You can only go there during the day.")
		return false, 0, nil
	}

	if toE.Flags["night_only"] && objects.DayTime {
		m.bad("You can only go there at night.")
		return false, 0, nil
	}

	if char.Equipment.GetWeight() > char.MaxWeight() {
		m.bad("You are carrying too much to move.")
		return false, 0, nil
	}

	// Asked before the fall roll below: no reason to drop someone down a shaft
	// they were never going to be let into.
	if to.Crowded() {
		m.info("That area is crowded.")
		return false, 0, nil
	}

	hasRope := false
	if char.Equipment.Off != (*objects.Item)(nil) {
		if char.Equipment.Off.ItemId == 1463 {
			hasRope = true
		}
	}
	if toE.Flags["levitate"] && !char.CheckFlag("levitate") && !hasRope {
		chanceToPass := char.GetStat("dex")/45 + 10
		if utils.Roll(100, 1, 0) >= chanceToPass {
			m.bad("You fall while trying to go that way! " + takeFallDamage(char))
			go char.DeathCheck("fell to their death.")
			return false, 0, nil
		}
	}

	// Check if anyone blocks.
	for _, mob := range from.Mobs.Contents {
		if _, inList := mob.ThreatTable[char.Name]; !inList {
			continue
		}
		// A stunned mob can't step in the way.
		if mob.CheckFlag("block_exit") && mob.Placement == char.Placement && !mob.Stunned() && !mob.CheckFlag("run_away") {
			evasive = 2
			curChance := config.MobBlock - ((char.Tier - mob.Level) * config.MobBlockPerLevel)
			if curChance > 85 {
				curChance = 85
			}
			if mob.CheckFlag("crippled") {
				curChance = config.CrippledChance(curChance)
			}
			if utils.Roll(100, 1, 0) <= curChance {
				m.bad(mob.Name + " blocks your way.")
				char.SetTimer("global", 4)
				return false, 0, nil
			}
			break
		}
	}

	// No one blocked, so work out who is entitled to give chase. Only the
	// entitlement is decided here - whether each one actually comes along is
	// FollowChar's roll, and the parting strike hangs off that, not off this.
	for _, mob := range from.Mobs.Contents {
		if _, inList := mob.ThreatTable[char.Name]; !inList {
			continue
		}
		if mob.CurrentTarget != char.Name {
			continue
		}
		if !mob.CheckFlag("follows") || mob.CheckFlag("curious_canticle") {
			continue
		}
		// A stunned or freshly spawned mob can't give chase, so it earns the character no evade timer either.
		if mob.Stunned() || mob.FollowLocked(char.Name) {
			continue
		}

		evasive = 4
		chasers = append(chasers, mob)
	}

	return true, evasive, chasers
}

// followVital rolls the parting strike a mob earns for chasing a character
// through an exit, and reports whether it landed and whether it killed them.
//
// This is only ever called for a mob that has already followed. The strike is
// the mob catching up with them on the other side, so one that lost its follow
// roll - or was never in any shape to give chase - has nothing to swing at. It
// lands in the destination room, after the move, for the same reason.
func followVital(m mover, mob *objects.Mob) (landed bool, died bool) {
	char := m.char

	// Whoever took the hit is the one to death check.
	return mob.RevengeVital(char, config.MobFollowVital-(char.GetStat("dex")/2), nil, objects.StrikeOpts{
		Metric:   "follow",
		Mode:     1,
		DeathMsg: "was slain by a " + mob.Name + ".",
	})
}

// takeFallDamage applies the damage for a failed levitate roll and returns the
// line describing it.
func takeFallDamage(char *objects.Character) string {
	fallDamageStam := int(config.FallDamage*float64(char.Stam.Max)) -
		(config.ConFallDamageMod * char.GetStat("con")) -
		(config.DexFallDamageMod * char.GetStat("dex"))
	fallDamageVit := int(config.FallDamage*float64(char.Stam.Max)) -
		(config.ConFallDamageMod * char.GetStat("con")) -
		(config.DexFallDamageMod * char.GetStat("dex"))

	totStam, totVit := 0, 0
	if fallDamageStam > 0 {
		totStam, totVit = char.ReceiveDamageNoArmor(fallDamageStam)
	}
	if fallDamageVit > 0 {
		totVit += char.ReceiveVitalDamageNoArmor(fallDamageVit)
	}

	if totStam <= 0 && totVit <= 0 {
		return "You take no damage in the fall."
	}

	buildStr := ""
	if totStam >= 1 {
		buildStr += "You take " + strconv.Itoa(totStam) + " points of stamina"
	}
	if totVit >= 1 {
		if totStam >= 1 {
			buildStr += " and "
		}
		buildStr += strconv.Itoa(totVit) + " points of vitality"
	}
	return buildStr + " damage in the fall."
}

// moveChar takes the character out of one room and puts them in the other,
// with the evade timer their exit earned. Both rooms must be locked by the
// caller.
func moveChar(char *objects.Character, from *objects.Room, to *objects.Room, evasive int) {
	from.Chars.Remove(char)
	// If they were evasive, add a global timer
	char.SetTimer("evade", evasive)
	to.Chars.Add(char)
	char.Victim = nil
	char.Placement = 3
	char.ParentId = to.RoomId
}
