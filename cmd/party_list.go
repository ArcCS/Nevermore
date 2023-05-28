package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(party{},
		"Usage:  party # \n\n List your current party information.",
		permissions.Player,
		"party")
}

type party cmd

func (party) process(s *state) {

	if s.actor.PartyFollow != "" {
		s.msg.Actor.SendInfo("You are currently following " + s.actor.PartyFollow + ".")
	} else {
		s.msg.Actor.SendInfo("You aren't following anyone.")
	}

	var followerList []string

	for _, findPlayer := range s.actor.PartyFollowers {
		player := objects.ActiveCharacters.Find(findPlayer)
		if player != nil {
			if !player.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				followerList = append(followerList, player.Name)
			}
		}
	}

	if len(followerList) > 0 {
		s.msg.Actor.SendInfo("Current party:")
		for _, player := range followerList {
			s.msg.Actor.SendInfo("\t " + player)
		}

	} else {
		s.msg.Actor.SendInfo("There is no one following you.")
	}

	s.ok = true
}
