package cmd

import (
	"strings"
)

func init() {
	addHandler(roomall{}, "roomall")
	addHelp("Usage:  roomall A mouse darts by in the corner of the square", 60, "roomall")
}

type roomall cmd

func (roomall) process(s *state) {
	if s.actor.Class < 60 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What happened in this room?")
		return
	}

	s.msg.Actor.SendInfo(strings.Join(s.input, " "))
	s.msg.Observer.SendInfo(strings.Join(s.input, " "))
	s.ok = true
	return
}
