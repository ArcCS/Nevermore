package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/spells"
	"strconv"
	"strings"
)

func init() {
	addHandler(cast{},
           "Usage:  cast spell_name target # \n\n Attempts to cast a known spell from your spellbook",
           permissions.Player,
           "cast")
}

type cast cmd

func (cast) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("What do you want to cast and what on?")
		return
	}

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 2 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[2]); err == nil {
			nameNum = val
		}
	}

	spellInstance, ok := spells.Spells[strings.ToLower(s.input[0])]; if !ok {
		s.msg.Actor.SendBad("What spell do you want to cast?")
		return
	}

	// Try Mobs First
	var whatMob *objects.Mob
	if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		whatMob = s.where.Mobs.Search(name, nameNum,true)
	}else{
		whatMob = s.where.Mobs.Search(name, nameNum,false)
	}
	// It was a mob!
	if whatMob != nil {
		spells.Effects[spellInstance.Effect](s.actor, whatMob, spellInstance.Magnitude)
		return
	}

	// Are we casting on a character
	var whatChar *objects.Character
	if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		whatChar = s.where.Chars.Search(name, true)
	}else{
		whatChar = s.where.Chars.Search(name, false)
	}
	// It was a person!
	if whatChar != nil {
		if strings.Contains(spellInstance.Effect, "damage") {
			//TODO PVP flags etc.
			s.msg.Actor.SendBad("No PVP implemented yet. ")
		}else{
			whatChar = s.actor
		}
	}


	s.msg.Actor.SendInfo("")
	s.ok = true
}
