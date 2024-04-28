package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

// Syntax: SAY <message> | " <message>
func init() {
	addHandler(osay{},
		"Usage:  osay \n \n Say something out of character",
		permissions.Player,
		"OSAY")
}

type osay cmd

func (osay) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	who := s.actor.Name

	if s.actor.Flags["invisible"] {
		who = "Someone"
	}

	msg := strings.Join(s.input, " ")
	s.actor.RunHook("say")
	if msg[len(msg)-1:] == "?" {
		s.msg.Actor.SendGood("You ask, out of character: \"", msg, "\"")
		s.msg.Observers.SendInfo(who, " asks, out of character: \"", msg, "\"")
	} else if msg[len(msg)-1:] == "!" {
		s.msg.Actor.SendGood("You exclaim, out of character: \"", msg, "\"")
		s.msg.Observers.SendInfo(who, " exclaims, out of character: \"", msg, "\"")
	} else {
		s.msg.Actor.SendGood("You say, out of character: \"", msg, "\"")
		s.msg.Observers.SendInfo(who, " says, out of character:  \"", msg, "\"")
	}

	s.ok = true
	return
}
