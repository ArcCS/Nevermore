package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

// Syntax: SAY <message> | " <message>
func init() {
	addHandler(pose{},
		"Usage:  pose \n \n Place your character into a passive RP pose!  (Skip the 'IS', auto appended)",
		permissions.Player,
		"POSE")
}

type pose cmd

func (pose) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("How do you want to pose?")
		return
	}

	msg := strings.Join(s.input, " ")
	s.actor.RunHook("say")
	data.StoreChatLog(0, s.actor.CharId, 0, msg)
	s.actor.Pose = msg
	s.msg.Actor.SendGood("You pose: \"", msg, "\"")

	s.ok = true
	return
}
