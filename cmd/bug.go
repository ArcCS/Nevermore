package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strings"
)

func init() {
	addHandler(errorReport{},
		"Usage:  error_report Disconnected when attacking all the rats\n \n Send a bug report to the GM's, be sure to include a time if you are reporting a bug much later, it will be a tad bit easier to fix.", permissions.Player,
		"error_report")

}

type errorReport cmd

func (errorReport) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	message := "### " + s.actor.Name + " reporting bug: " + strings.Join(s.input, " ")
	objects.ActiveCharacters.MessageGM(message)
	if objects.DiscordSession != nil {
		if _, err := objects.DiscordSession.ChannelMessageSend("729467777416691712", message); err != nil {
			log.Println("Error sending message to discord: ", err)
		}
	}

	s.ok = true
	return
}
