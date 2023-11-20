package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strings"
)

func init() {
	addHandler(appeal{},
		"Usage:  appeal HELP ME OH GODS OR CREATORS\n \n Appeal a message to higher powers.  Note: Preppend OOC to clarify it's a non RP issue", permissions.Player,
		"APPEAL")

}

type appeal cmd

func (appeal) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	message := "### " + s.actor.Name + " appeals: " + strings.Join(s.input, " ")
	objects.ActiveCharacters.MessageGM(message)
	if objects.DiscordSession != nil {
		if _, err := objects.DiscordSession.ChannelMessageSend(config.AppealChannel, message); err != nil {
			log.Println("Error sending message to discord: ", err)
		}
	}

	s.ok = true
	return
}
