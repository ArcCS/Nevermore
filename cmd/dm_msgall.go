package cmd

import (
	"github.com/ArcCS/Nevermore/stats"
	"strings"
)

func init() {
	addHandler(msgall{}, "msgall")
	addHelp("Usage:  msgall A thunderstorm rolls in from the east", 60, "msgall")
}

type msgall cmd

func (msgall) process(s *state) {
	if s.actor.Class < 60 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to tell the realms?")
		return
	}

	stats.ActiveCharacters.MessageAll("###: " + strings.Join(s.input, " "))

	s.ok = true
	return
}
