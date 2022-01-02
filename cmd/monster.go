package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

// Overloaded Look object for all of your looking pleasure
// Syntax: ( LOOK | L ) has.Thing
func init() {
	addHandler(monster{},
		"Usage:  monster \n \n List all mobs in the room and their status",
		permissions.Player,
		"monster", "monst", "mobs")
}

type monster cmd

func (monster) process(s *state) {
	mobList := objects.Rooms[s.actor.ParentId].Mobs.ListMobs(s.actor)
	if len(mobList) > 0 {
		for _, mob := range objects.Rooms[s.actor.ParentId].Mobs.ListMobs(s.actor) {
			s.msg.Actor.SendInfo(mob.Name + mob.ReturnState() + "," + utils.WhereAt(mob.Placement, s.actor.Placement))
		}
	}else{
		s.msg.Actor.SendInfo("There is no one visible.")
	}
	s.ok = true

	return
}
