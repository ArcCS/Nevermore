package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(broadcast{},
		"Usage:  broadcast I have so many things to talk about! \n \n Broadcast messages to the entire realm at the cost of a broadcast point.",
		permissions.Player,
		"BROADCAST")
}

type broadcast cmd

func (broadcast) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}
	message := "### " + s.actor.Name + ": " + strings.Join(s.input, " ")
	if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.God, permissions.NPC, permissions.Dungeonmaster, permissions.Gamemaster) {
		if s.actor.Broadcasts < 1 {
			s.msg.Actor.SendBad("You're out of broadcasts today.")
		} else {
			s.actor.Broadcasts -= 1
			objects.ActiveCharacters.MessageAll(message, config.BroadcastChannel)
		}
	} else {
		objects.ActiveCharacters.MessageAll(message, config.BroadcastChannel)
	}

	s.ok = true
	return
}
