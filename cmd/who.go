package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
)

// Syntax: WHO
func init() {
	addHandler(who{},
		"Usage:  who \n \n Display other currently logged in characters.",
		permissions.Player,
		"WHO")
}

type who cmd

func (who) process(s *state) {
	var players []string
	if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		players = objects.ActiveCharacters.GMList()
	} else {
		players = objects.ActiveCharacters.List()
	}

	s.msg.Actor.SendInfo("You sense the presence of " + config.TextNumbers[len(players)-1] + " other beings (Tiers 1-25):\n")

	for _, player := range players {
		s.msg.Actor.Send(player)
	}

	s.ok = true
}
