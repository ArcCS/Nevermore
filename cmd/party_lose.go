package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"strings"
)

func init() {
	addHandler(lose{},
		"Usage:  lose name|all \n\n Lose a person or everyone from your party.",
		permissions.Player,
		"lose")
}

type lose cmd

func (lose) process(s *state) {
	if len(s.input) < 1  || strings.ToLower(s.words[0]) == "all"{
		if len(s.actor.PartyFollowers) > 0 {
			s.actor.LoseParty()
			s.msg.Actor.SendInfo("You lose everyone following you.")
		}else{
			s.msg.Actor.SendInfo("There is no one following you.")
		}
		return
	}

	name := s.input[1]
	var whatChar *objects.Character
	whatChar = s.where.Chars.Search(name, s.actor)
	if whatChar != nil {
		for c, player := range s.actor.PartyFollowers {
			if player == whatChar {
				if !player.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					copy(s.actor.PartyFollowers[c:], s.actor.PartyFollowers[c+1:])
					s.actor.PartyFollowers[len(s.actor.PartyFollowers)-1] = nil
					s.actor.PartyFollowers= s.actor.PartyFollowers[:len(s.actor.PartyFollowers)-1]
					s.msg.Actor.SendInfo("You lose " + player.Name + ".")
					player.Write([]byte(text.Info + s.actor.Name + " loses you."))
					return
				}
			}
			s.msg.Actor.SendInfo("That person isn't in your party.")
		}
		if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			s.msg.Participant.SendInfo(s.actor.Name, " follows you.")
		}
	}else{
		s.msg.Actor.SendBad("Who ya followin'??")
	}

	s.ok = true
}
