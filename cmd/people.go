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
	// Pick whether it's a GM or a user looking and go for it.
	if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		others = objects.Rooms[s.actor.ParentId].Chars.List(true, true, s.actor.Name, true)
	} else {
		others = objects.Rooms[s.actor.ParentId].Chars.List(false, false, s.actor.Name, false)
	}
	if len(others) == 1 {
		s.msg.Actor.SendInfo(strings.Join(others, ", "), " is also here.")
	} else if len(others) > 1 {
		s.msg.Actor.SendInfo(strings.Join(others, ", "), " are also here.")
	}

	s.ok = true

	return
}
