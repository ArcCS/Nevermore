package cmd

import "strconv"

func init() {
	addHandler(move{},
		 "SPRINT", "MOVE", "MV", "SPR", "F", "FORWARD", "BACK", "B",
	)
	addHelp("Usage:  move|sprint forward|backwards # \n \n Change combat stance in the room's 5 space grid.  Sprint will move 2, move will move 1", 0, "move", "sprint")
}


type move cmd


// Move allows characters to change their position along a 5 point grid system in a room.
func (move) process(s *state) {

	spaces := 1
	if s.cmd == "SPRINT" || s.cmd == "SPR"{
		spaces = 2
		s.msg.Actor.SendInfo("You wanted to move: ", strconv.Itoa(spaces))
	}
	return
}
