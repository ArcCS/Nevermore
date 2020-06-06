package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(use{},
           "Usage:  use item # \n\n Use an item",
           permissions.Player,
           "USE")
}

//TODO: Map out the use of items and the effect they map to under spells

type use cmd

func (use) process(s *state) {
	s.msg.Actor.SendGood("Use what?")
	s.ok = true
}
