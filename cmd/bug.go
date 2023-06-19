package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(error_report{},
		"Usage:  error_report Disconnected when attacking all the rats\n \n Send a bug report to the GM's, be sure to include a time if you are reporting a bug much later, it will be a tad bit easier to fix.", permissions.Player,
		"error_report")

}

type error_report cmd

func (error_report) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	message := "### " + s.actor.Name + " reporting bug: " + strings.Join(s.input, " ")
	objects.ActiveCharacters.MessageGM(message)
	if objects.DiscordSession != nil {
		objects.DiscordSession.ChannelMessageSend("729467777416691712", message)
	}

	s.ok = true
	return
}
