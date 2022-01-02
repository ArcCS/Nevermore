package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(whisper{},
		"Usage:  whisper \n \n Get the specified item.",
		permissions.Player,
		"WHISPER")
}

type whisper cmd

func (whisper) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	whoSays := s.actor.Name
	whoStr := s.words[0]

	if s.actor.Flags["invisible"] {
		whoSays = "Someone"
	}

	var who *objects.Character
	who = s.where.Chars.Search(whoStr, s.actor)
	if who == nil {
		s.msg.Actor.SendInfo("Give who what???")
		return
	}
	s.participant = who

	msg := strings.Join(s.input[1:], " ")

	if msg[len(msg)-1:] == "?" {
		s.msg.Actor.SendGood("You whisper to "+who.Name+": \"", msg, "\"")
		s.msg.Participant.SendInfo(whoSays, " whispers to  you: \"", msg, "\"")
		s.msg.Observers.SendInfo(whoSays, " whispers to "+who.Name)
	}

	s.ok = true
	return
}
