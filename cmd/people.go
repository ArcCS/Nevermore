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
	if s.actor.CheckFlag("berserk") {
		s.msg.Actor.SendGood("The red rage has you, there's no time for that.\n")
		return
	}

	for _, char := range objects.Rooms[s.actor.ParentId].Chars.ListPeoChars(s.actor) {
		name := char.Name
		if char.CheckFlag("singing") {
			name += " (singing " + char.SongDescriptor() + ")"
		}
		s.msg.Actor.SendInfo(name + char.ReturnState() + "," + utils.WhereAt(char.Placement, s.actor.Placement))
	}
	s.ok = true

	return
}
