package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
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

	for _, char := range objects.Rooms[s.actor.ParentId].Chars.ListChars(s.actor) {
		s.msg.Actor.SendInfo(char.Name + char.ReturnState() + "," + utils.WhereAt(char.Placement, s.actor.Placement))
	}
	s.ok = true

	return
}
