package cmd

import (
	"github.com/ArcCS/Nevermore/jarvoral"
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

	message := "###: " + strings.Join(s.input, " ")
	stats.ActiveCharacters.MessageAll(message)
	if jarvoral.DiscordSession != nil {
		jarvoral.DiscordSession.ChannelMessageSend("854733320474329088", message)
	}

	s.ok = true
	return
}
