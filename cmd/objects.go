package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

// Overloaded Look object for all of your looking pleasure
// Syntax: ( LOOK | L ) has.Thing
func init() {
	addHandler(obj{},
		"Usage:  objects \n \n List objects and locations.",
		permissions.Player,
		"objects", "obj", "items", "ob")
}

type obj cmd

func (obj) process(s *state) {

	for _, obj:= range objects.Rooms[s.actor.ParentId].Items.ListItems() {
		s.msg.Actor.SendInfo(obj.Name + utils.WhereAt(obj.Placement, s.actor.Placement))
	}
	s.ok = true

	return
}
