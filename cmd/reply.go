package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(reply{},
		"Usage:  reply \n \n Send a telepathic message to the last person that messaged you",
		permissions.Player,
		"r", "reply", "rep")
}

type reply cmd

func (reply) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Send what?")
		return
	}

	if s.actor.LastMessenger == "" {
		s.msg.Actor.SendBad("No one has sent you a message recently.")
		return
	}

	s.scriptActor("tell " + s.actor.LastMessenger + " " + s.original)

	s.ok = true
	return
}
