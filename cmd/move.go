package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(move{},
		"Usage:  move|sprint forward|backwards # \n \n Change combat stance in the room's 5 space grid.  Sprint will move 2, move will move 1",
		permissions.Player,
		"SPRINT", "MOVE", "MV", "SPR", "F", "FORWARD", "BACK", "B")
}

type move cmd

// Move allows characters to change their position along a 5 point grid system in a room.
func (move) process(s *state) {
	if s.cmd != "F" && s.cmd != "B" && len(s.words) < 1 {
		s.msg.Actor.SendInfo("You wanted to move where?")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	if s.actor.Equipment.Weight > s.actor.MaxWeight() {
		s.msg.Actor.SendBad("You are carrying too much to move.")
		return
	}

	s.actor.RunHook("gridmove")

	previous := s.actor.Placement
	if s.cmd == "SPRINT" || s.cmd == "SPR" {
		direction := string(s.words[0][0])
		if direction == "F" {
			if 5-s.actor.Placement >= 2 {
				s.actor.Placement += 2
				s.actor.SetTimer("global", 8)
				if s.actor.Placement == 5 {
					s.msg.Actor.SendGood("You sprint forward, to the front of the room.")
				} else {
					s.msg.Actor.SendGood("You sprint forward.")
				}
				if s.actor.Flags["hidden"] != true && s.actor.Flags["invisible"] != true {
					for _, char := range s.where.Chars.Contents {
						if char != s.actor {
							char.WriteMovement(previous, s.actor.Placement, s.actor.Name)
						}
					}
				}
			} else {
				s.msg.Actor.SendBad("There's not enough room to sprint forward")
			}
		} else if direction[0:] == "B" {
			if s.actor.Placement-1 >= 2 {
				s.actor.Placement -= 2
				s.actor.SetTimer("global", 8)
				if s.actor.Placement == 1 {
					s.msg.Actor.SendGood("You sprint back, to the back of the room.")
				} else {
					s.msg.Actor.SendGood("You sprint backward.")
				}
				if s.actor.Flags["hidden"] != true && s.actor.Flags["invisible"] != true {
					for _, char := range s.where.Chars.Contents {
						if char != s.actor {
							char.WriteMovement(previous, s.actor.Placement, s.actor.Name)
						}
					}
				}
			} else {
				s.msg.Actor.SendBad("There's not enough room to sprint backward")
			}
		}
	} else if s.cmd == "F" || s.cmd == "FORWARD" || (len(s.words) > 0 && string(s.words[0][0]) == "F") {
		if 5-s.actor.Placement >= 1 {
			s.actor.Placement += 1
			s.actor.SetTimer("global", 4)
			if s.actor.Placement == 5 {
				s.msg.Actor.SendGood("You move forward, to the front of the room.")
			} else {
				s.msg.Actor.SendGood("You move forward.")
			}
			if s.actor.Flags["hidden"] != true && s.actor.Flags["invisible"] != true {
				for _, char := range s.where.Chars.Contents {
					if char != s.actor {
						char.WriteMovement(previous, s.actor.Placement, s.actor.Name)
					}
				}
			}
		} else {
			s.msg.Actor.SendBad("There's not enough room to move forward")
		}
	} else if s.cmd == "B" || s.cmd == "BACK" || (len(s.words) > 0 && string(s.words[0][0]) == "B") {
		if s.actor.Placement-1 >= 1 {
			s.actor.Placement -= 1
			s.actor.SetTimer("global", 4)
			if s.actor.Placement == 1 {
				s.msg.Actor.SendGood("You move back, to the back of the room.")
			} else {
				s.msg.Actor.SendGood("You move backward.")
			}
			if s.actor.Flags["hidden"] != true && s.actor.Flags["invisible"] != true {
				for _, char := range s.where.Chars.Contents {
					if char != s.actor {
						char.WriteMovement(previous, s.actor.Placement, s.actor.Name)
					}
				}
			}
		} else {
			s.msg.Actor.SendBad("There's not enough room to move backwards.")
		}
	}
	return
}
