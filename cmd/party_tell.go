package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"strings"
)

func init() {
	addHandler(ptell{},
		"Usage:  ptell # \n\n Send a message to your whole party.",
		permissions.Player,
		"ptell", "partytell")
}

type ptell cmd

func (ptell) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	msg := strings.Join(s.input, " ")
	msg = text.White + s.actor.Name + " party flashes# \"" + msg + "\""
	if s.actor.PartyFollow == "" && len(s.actor.PartyFollowers) == 0 {
		s.msg.Actor.SendBad("You have no party to telepathically communicate with.")
	}
	if s.actor.PartyFollow != "" {
		leadChar := objects.ActiveCharacters.Find(s.actor.PartyFollow)
		s.participant = leadChar
		s.msg.Participant.Send(msg)
		if leadChar != nil {
			leadChar.MessageParty(msg, s.actor)
		}
	}
	if len(s.actor.PartyFollowers) > 0 {
		s.actor.MessageParty(msg, s.actor)
	}
	s.msg.Actor.Send(text.White + "You sent to party#, \"" + msg + "\"")

	s.ok = true
}
