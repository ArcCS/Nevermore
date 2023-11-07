package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
)

func init() {
	addHandler(remspawn{},
		"Usage:  remspawn id \n \n Remove a spawn from the encounter table \n",
		permissions.Builder,
		"remspawn")
}

type remspawn cmd

func (remspawn) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Remove what?")
		return
	}

	var mobId int
	val, err := strconv.Atoi(s.words[0])
	if err != nil {
		log.Println(err)
	}
	mobId = val

	if _, ok := s.where.EncounterTable[mobId]; ok {
		delete(s.where.EncounterTable, mobId)
		data.DeleteEncounter(mobId, s.actor.ParentId)
		s.msg.Actor.SendGood("Mob removed from this room's encounter table.")
	} else {
		s.msg.Actor.SendBad("That mob ID doesn't encounter here.")
		return
	}

	s.ok = true
	return
}
