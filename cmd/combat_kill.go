package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"strconv"
)

func init() {
	addHandler(kill{}, "kill", "k", "attack")
	addHelp("Usage:  kill target # \n\n Try to attack something! Can also use attack, or shorthand k", 0, "kill")
}

type kill cmd

func (kill) process(s *state) {
	if s.actor.Class == 50 {
		s.msg.Actor.SendInfo("As a builder you can't use these commands.")
		return
	}

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	var whatMob *objects.Mob

	// This is an override for a GM to delete a mob
	if s.actor.Class >= 60 {
		whatMob = s.where.Mobs.Search(name, int64(nameNum),true)
		if whatMob != nil {
			s.msg.Actor.SendInfo("You smashed ", whatMob.Name , " out of existence.")
			objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
			whatMob = nil
			return
		}
	}

	s.msg.Actor.SendInfo("You focus really hard but...  couldn't muster up an attack on anything...")
	s.ok = true
}
