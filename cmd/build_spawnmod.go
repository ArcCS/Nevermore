package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
)

func init() {
	addHandler(modspawn{},
	"Usage:  modspawn 452 39 \n Modify a current spawn with a new value \n",
	permissions.Builder,
	"modspawn")
}

type modspawn cmd

func (modspawn) process(s *state) {
	if len(s.words) < 2{
		s.msg.Actor.SendInfo("Add what, where?")
		return
	}

	var mob_id, mob_rate int64
	val, err := strconv.Atoi(s.words[0])
	if err != nil {
		log.Println(err)
	}
	mob_id = int64(val)

	val2, err2 := strconv.Atoi(s.words[1])
	if err != nil {
		log.Println(err2)
	}
	mob_rate = int64(val2)

	if _, ok := s.where.EncounterTable[mob_id]; ok {
		previousRate := s.where.EncounterTable[mob_id]
		s.where.EncounterTable[mob_id] = mob_rate
		var sumVals int64
		for _, v := range s.where.EncounterTable {
			sumVals += v
		}
		if sumVals > 100 {
			s.where.EncounterTable[mob_id] = previousRate
			s.msg.Actor.SendBad("The sum of the encounter rates is more than 100% with the new value")
		}else{
			data.UpdateEncounter(map[string]interface{}{
				"mobId": mob_id,
				"roomId":  s.where.RoomId,
				"chance": mob_rate,})
			s.msg.Actor.SendGood("Mob spawn rate updated")
		}

	}else{
		s.msg.Actor.SendBad("That mob ID doesn't exist")
		return
	}

	s.ok = true
	return
}
