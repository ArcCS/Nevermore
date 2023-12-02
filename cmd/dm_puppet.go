package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(puppet{},
		"Usage:  puppet <charName> \n \n Moves a character from an account to the puppet list. ",
		permissions.Dungeonmaster|permissions.Gamemaster,
		"puppet")
}

type puppet cmd

func (puppet) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Who are you trying to puppet?")
		return
	}

	if data.PuppetChar(s.words[0]) {
		s.msg.Actor.SendGood("Character has been made a puppet.")
	} else {
		s.msg.Actor.SendBad("Couldn't puppet that character.")
	}

	s.ok = true
	return
}
