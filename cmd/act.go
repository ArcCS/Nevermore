package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(act{},
		"Usage:  act performs for all to see \n \n Perform actions.",
		permissions.Player,
		"act", "emote")
}

type act cmd

func (act) process(s *state) {

	// Did they send an action?
	if len(s.words) == 0 {
		s.msg.Actor.SendBad("... what were you trying to do???")
		return
	}

	action := strings.Join(s.input, " ")

	s.msg.Actor.SendInfo(s.actor.Name, " ", action)
	s.msg.Observer.SendInfo(s.actor.Name, " ", action)

	s.ok = true
}
