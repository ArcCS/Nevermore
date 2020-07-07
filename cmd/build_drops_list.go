package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(listdrops{},
		"Usage:  listdrops mob_id \n \n List drops for the given mob. \n",
		permissions.Builder,
		"listdrops")
}

type listdrops cmd

func (listdrops) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Which mob did you want to list drops for?")
		return
	}
	mobId, _ := strconv.Atoi(s.words[0])
	mob, ok := objects.Mobs[mobId]
	if ok {
		s.msg.Actor.SendInfo("List of items dropped by " + mob.Name + "\n =================")
		for k, v := range mob.ItemList {
			if _, ok := objects.Items[k]; ok {
				s.msg.Actor.SendInfo("ItemId: ", strconv.Itoa(k), " ", objects.Mobs[k].Name, "   Rate: ", strconv.Itoa(v))
			} else {
				delete(s.where.EncounterTable, k)
			}
		}

		s.ok = true
		return
	}
}
