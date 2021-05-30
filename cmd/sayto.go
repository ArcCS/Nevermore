package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

// Syntax: SAY <message> | " <message>
func init() {
	addHandler(sayto{},
		"Usage:  sayto [character] \n \n Say something out loud!",
		permissions.Player,
		"SAYTO")
}

type sayto cmd

func (sayto) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	for _, loc := range s.where.Exits {
		room := objects.Rooms[loc.ToId]
		room.Chars.Lock()
		room.MessageAll("You hear someone speaking nearby.")
		room.Chars.Unlock()
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
		s.msg.Actor.SendGood("You ask "+who.Name+": \"", msg, "\"")
		s.msg.Participant.SendInfo(whoSays, " asks you: \"", msg, "\"")
		s.msg.Observers.SendInfo(whoSays, " asks "+who.Name+": \"", msg, "\"")
	} else if msg[len(msg)-1:] == "!" {
		s.msg.Actor.SendGood("You exclaim to "+who.Name+": \"", msg, "\"")
		s.msg.Participant.SendInfo(whoSays, " exclaims to you: \"", msg, "\"")
		s.msg.Observers.SendInfo(whoSays, " exclaims to "+who.Name+": \"", msg, "\"")
	} else {
		s.msg.Actor.SendGood("You say to "+who.Name+": \"", msg, "\"")
		s.msg.Participant.SendInfo(whoSays, " says to you: \"", msg, "\"")
		s.msg.Observers.SendInfo(whoSays, " says to "+who.Name+": \"", msg, "\"")
	}


	s.ok = true
	return
}
