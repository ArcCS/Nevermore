package cmd

import (
	"github.com/ArcCS/Nevermore/jarvoral"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(bug{},
		"Usage:  bug Disconnected when attacking all the rats\n \n Send a bug report to the GM's, be sure to include a time if you are reporting a bug much later, it will be a tad bit easier to fix.", permissions.Player,
		"BUG")

}

type bug cmd

func (bug) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	message := "### " + s.actor.Name + " reporting bug: " + strings.Join(s.input, " ")
	objects.ActiveCharacters.MessageGM(message)
	if jarvoral.DiscordSession != nil {
		jarvoral.DiscordSession.ChannelMessageSend("729467777416691712", message)
	}

	s.ok = true
	return
}
