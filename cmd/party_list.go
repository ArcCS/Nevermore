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

	if s.actor.PartyFollow != nil {
		s.msg.Actor.SendInfo("You are currently following " +  s.actor.PartyFollow.Name)
	}else{
		s.msg.Actor.SendInfo("You aren't following anyone.")
	}

	var followerList []*objects.Character

	for _, player := range s.actor.PartyFollowers {
		if !player.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			followerList = append(followerList, player)
		}
	}

	if len(followerList) > 0 {
		s.msg.Actor.SendInfo("Current party:")
		for _, player := range followerList{
			s.msg.Actor.SendInfo("\t " + player.Name)
		}

	}else{
		s.msg.Actor.SendInfo("There is no one following you.")
	}


	s.ok = true
}
