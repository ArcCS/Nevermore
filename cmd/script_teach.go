package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
)

func init() {
	addHandler(scriptTeach{},
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
		log.Printf("That's not a known spell.")
	}

	s.ok = true
}
