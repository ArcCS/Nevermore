package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(listspawn{},
		"Usage:  listspawn \n \n List the spawns in the current room \n" ,
		permissions.Builder,
		"listspawn")
}

type listspawn cmd

func (listspawn) process(s *state) {
	s.msg.Actor.SendInfo("Spawn Rate: ", strconv.Itoa(int(s.where.EncounterRate)) ,"% per encounter tick\nList of Mobs that spawn here \n =================")
	for k, v := range s.where.EncounterTable {
		if _, ok := objects.Mobs[k]; ok {
			s.msg.Actor.SendInfo("MobId: ", strconv.Itoa(int(k)), " ", objects.Mobs[k].Name, "   Rate: ", strconv.Itoa(int(v)))
		}else{
			delete(s.where.EncounterTable, k)
		}
	}

	s.ok = true
	return
}
