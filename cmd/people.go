package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

// Overloaded Look object for all of your looking pleasure
// Syntax: ( LOOK | L ) has.Thing
func init() {
	addHandler(people{},
		"Usage:  people \n \n List everyone in the room and their status",
		permissions.Player,
		"people", "peo", "p")
}

type people cmd

func (people) process(s *state) {
	//TODO This should list some health and stuff
	var others []string
	others = objects.Rooms[s.actor.ParentId].Chars.List(s.actor)
	if len(others) == 1 {
		s.msg.Actor.SendInfo(strings.Join(others, ", "), " is also here.")
	} else if len(others) > 1 {
		s.msg.Actor.SendInfo(strings.Join(others, ", "), " are also here.")
	}

	s.ok = true

	return
}
