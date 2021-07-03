package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(victim{},
		"Usage:  victim \n \n Show your current victim and state",
		permissions.Player,
		"victim", "vic", "v")
}

type victim cmd

func (victim) process(s *state) {

	s.msg.Actor.SendInfo(s.actor.ReturnVictim())

	s.ok = true

	return
}
