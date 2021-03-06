package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(follow{},
		"Usage:  ptell # \n\n Send a message to your whole party.",
		permissions.Player,
		"ptell", "partytell")
}

type ptell cmd

func (ptell) process(s *state) {
	message := strings.Join(s.input[1:], " ")
	if s.actor.PartyFollow == nil && len(s.actor.PartyFollowers) == 0 {
		s.msg.Actor.SendBad("You have no party to telepathically communicate with.")
	}
	if s.actor.PartyFollow != nil {
		s.actor.PartyFollow.MessageParty(message)
	}
	if len(s.actor.PartyFollowers) > 0 {
			s.actor.MessageParty(message)
	}

	s.ok = true
}