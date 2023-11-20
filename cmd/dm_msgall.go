package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
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

	message := "###: " + strings.Join(s.input, " ")
	objects.ActiveCharacters.MessageAll(message, config.BroadcastChannel)

	s.ok = true
	return
}
