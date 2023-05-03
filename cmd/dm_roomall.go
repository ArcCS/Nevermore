package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(roomall{},
		"Usage:  roomall A mouse darts by in the corner of the square",
		permissions.Dungeonmaster,
		"roomall")
}

type roomall cmd

func (roomall) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What happened in this room?")
		return
	}

	s.msg.Actor.SendInfo(strings.Join(s.input, " "))
	s.msg.Observers.SendInfo(strings.Join(s.input, " "))
	s.ok = true
	return
}
