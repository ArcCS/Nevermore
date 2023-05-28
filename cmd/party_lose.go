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
	if len(s.input) < 1 || strings.ToLower(s.words[0]) == "all" {
		if len(s.actor.PartyFollowers) > 0 {
			s.actor.LoseParty()
			s.msg.Actor.SendInfo("You lose everyone following you.")
		} else {
			s.msg.Actor.SendInfo("There is no one following you.")
		}
		return
	}

	name := s.input[0]
	var whatChar *objects.Character
	whatChar = objects.ActiveCharacters.Find(name)
	if whatChar != nil {
		for c, player := range s.actor.PartyFollowers {
			if player == whatChar.Name {
				if !whatChar.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					s.actor.PartyFollowers = append(s.actor.PartyFollowers[:c], s.actor.PartyFollowers[c+1:]...)
					whatChar.PartyFollow = ""
					s.msg.Actor.SendInfo("You lose " + player + ".")
					whatChar.Write([]byte(text.Info + s.actor.Name + " loses you."))
					return
				} else {
					s.msg.Actor.SendInfo("That person isn't in your party.")
				}
			}
		}

	} else {
		s.msg.Actor.SendBad("That person isn't in your party.")
	}

	s.ok = true
}
