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

type scriptTeach cmd

func (scriptTeach) process(s *state) {

	spell := s.words[0]

	if _, ok := objects.Spells[spell]; ok {
		s.actor.Spells = append(s.actor.Spells, spell)
		s.msg.Participant.SendGood("You learned ", spell)
		return
	} else {
		s.msg.Actor.SendBad("That's not a known spell.")
	}

	s.ok = true
}
