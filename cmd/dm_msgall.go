package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/stats"
	"strings"
)

func init() {
	addHandler(msgall{},
		"Usage:  msgall A thunderstorm rolls in from the east",
		permissions.Dungeonmaster,
		"msgall")
}

type msgall cmd

func (msgall) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to tell the realms?")
		return
	}

	stats.ActiveCharacters.MessageAll("###: " + strings.Join(s.input, " "))

	s.ok = true
	return
}
