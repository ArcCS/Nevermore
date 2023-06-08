package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(unfollow{},
		"Usage:  unfollow \n\n Lose a person or everyone from your party.",
		permissions.Player,
		"unfollow")
}

type unfollow cmd

func (unfollow) process(s *state) {
	if s.actor.PartyFollow == "" {
		s.msg.Actor.SendBad("You aren't following anyone.")
	} else {
		s.actor.Unfollow()
	}
	s.ok = true
}
