package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(kill{},
           "Usage:  kill target # \n\n Try to attack something! Can also use attack, or shorthand k",
           permissions.Player,
           "kill")
}

type kill cmd

func (kill) process(s *state) {
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
	if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		whatMob = s.where.Mobs.Search(name, nameNum,true)
		if whatMob != nil {
			s.msg.Actor.SendInfo("You smashed ", whatMob.Name , " out of existence.")
			objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
			whatMob = nil
			return
		}
	}

	// Kill is really the fighters realm for specialty.. lets see if this is a fighter if so, lets try to lethal.
	if s.actor.Permission.HasAnyFlags(permissions.Fighter) {

	}
	1 Extra Attack	15%
		Ups 2nd attack to 30%
		2 Extra Attacks	15%
		Ups 3rd attack to 30%
		3 Extra Attacks	15%
		Ups 4th attack to 30%
		4 Extra Attacks	30%

	s.msg.Actor.SendInfo("You focus really hard but...  couldn't muster up an attack on anything...")
	s.ok = true
}
