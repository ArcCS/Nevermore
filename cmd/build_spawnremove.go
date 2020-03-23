package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"log"
	"strconv"
)

func init() {
	addHandler(remspawn{}, "remspawn")
	addHelp("Usage:  remspawn id \n \n Remove a spawn from the encounter table \n" ,50, "remspawn")
}

type remspawn cmd

func (remspawn) process(s *state) {
	// Handle Permissions
	if s.actor.Class < 50 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) < 1{
		s.msg.Actor.SendInfo("Remove what?")
		return
	}

	var mob_id int64
	val, err := strconv.Atoi(s.words[0])
	if err != nil {
		log.Println(err)
	}
	mob_id = int64(val)

	if _, ok := s.where.EncounterTable[mob_id]; ok {
		delete(s.where.EncounterTable, mob_id)
		data.DeleteEncounter(mob_id, s.actor.ParentId)
		s.msg.Actor.SendGood("Mob removed from this room's encounter table.")
	}else{
		s.msg.Actor.SendBad("That mob ID doesn't encounter here.")
		return
	}

	s.ok = true
	return
}
