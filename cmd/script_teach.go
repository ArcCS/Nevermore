package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(script_teach{},
		"",
		permissions.Player,
		"$TEACH")
}

type script_teach cmd

func (script_teach) process(s *state) {

	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("You need to specify the person and the spell you want to teach.")
		return
	}

	spell := s.words[0]

	if _, ok := objects.Spells[spell]; ok {
		s.actor.Spells = append(s.actor.Spells, spell)
		s.msg.Participant.SendGood("You learned ", spell)
		return
	}else{
		s.msg.Actor.SendBad("That's not a known spell.")
	}

	s.ok = true
}
