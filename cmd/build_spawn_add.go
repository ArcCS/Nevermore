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
		"Usage:  addspawn 452 39 \n Add a spawn to a room with a whole number chance of encounter when an encounter is triggered \n",
		permissions.Builder,
		"addspawn")
}

type addspawn cmd

func (addspawn) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Add what, where?")
		return
	}

	var mobId, mobRate int
	val, err := strconv.Atoi(s.words[0])
	if err != nil {
		log.Println(err)
	}
	mobId = val

	val2, err2 := strconv.Atoi(s.words[1])
	if err != nil {
		log.Println(err2)
	}
	mobRate = val2

	if _, ok := objects.Mobs[mobId]; ok {
		curSpawn := data.SumEncounters(s.where.RoomId)
		if curSpawn+mobRate <= 100 {
			data.CreateEncounter(map[string]interface{}{
				"mobId":  mobId,
				"roomId": s.actor.ParentId,
				"chance": mobRate})
			s.where.EncounterTable[mobId] = mobRate
			s.msg.Actor.SendGood("Mob added to this room's encounter table.")
		} else {
			s.msg.Actor.SendBad("The addition of this spawn rate would exceed 100%, mob not added to the encounter table")
		}
	} else {
		s.msg.Actor.SendBad("That mob ID doesn't exist")
		return
	}

	s.ok = true
	return
}
