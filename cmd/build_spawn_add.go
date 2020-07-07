package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
)

func init() {
	addHandler(addspawn{},
	"Usage:  addspawn 452 39 \n Add a spawn to a room with a whole number chance of encounter when an encounter is triggered \n" ,
	permissions.Builder,
	"addspawn")
}

type addspawn cmd

func (addspawn) process(s *state) {
	if len(s.words) < 2{
		s.msg.Actor.SendInfo("Add what, where?")
		return
	}

	var mob_id, mob_rate int
	val, err := strconv.Atoi(s.words[0])
	if err != nil {
		log.Println(err)
	}
	mob_id = val

	val2, err2 := strconv.Atoi(s.words[1])
	if err != nil {
		log.Println(err2)
	}
	mob_rate = val2

	if _, ok := objects.Mobs[mob_id]; ok {
		curSpawn := data.SumEncounters(s.where.RoomId)
		if curSpawn + mob_rate <= 100 {
			data.CreateEncounter(map[string]interface{}{
				"mobId":  mob_id,
				"roomId": s.actor.ParentId,
				"chance": mob_rate})
			s.where.EncounterTable[mob_id] = mob_rate
			s.msg.Actor.SendGood("Mob added to this room's encounter table.")
		}else{
			s.msg.Actor.SendBad("The addition of this spawn rate would exceed 100%, mob not added to the encounter table")
		}
	}else{
		s.msg.Actor.SendBad("That mob ID doesn't exist")
		return
	}

	s.ok = true
	return
}
