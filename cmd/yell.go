package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

// Syntax: SAY <message> | " <message>
func init() {
	addHandler(yell{},
		"Usage:  yell \n \n Yell something to everyone around!",
		permissions.Player,
		"yell")
}

type yell cmd

func (yell) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to yell?")
		return
	}

	who := s.actor.Name

	if s.actor.Flags["invisible"] {
		who = "Someone"
	}

	msg := strings.Join(s.input, " ")
	s.actor.RunHook("say")
	s.msg.Actor.SendGood("You yell: \"", msg, "\"")
	s.msg.Observers.SendInfo(who, " yells: \"", msg, "\"")
	for _, loc := range s.where.Exits {
		room := objects.Rooms[loc.ToId]
		room.MessageAll("Someone yells: \"" + msg + "\"")
	}

	s.ok = true
	return
}
