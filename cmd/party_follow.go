package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(follow{},
		"Usage:  follow player # \n\n Become a part of another players party",
		permissions.Player,
		"follow")
}

type follow cmd

func (follow) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendInfo("Who ya followin'??")
		return
	}

	name := s.input[0]
	var whatChar *objects.Character
	whatChar = s.where.Chars.Search(name, s.actor)
	if whatChar != nil {
		if s.actor.PartyFollow != nil {
			s.msg.Actor.SendInfo(s.actor.Name + " stops folowing you.")
		}
		s.participant = whatChar
		s.actor.PartyFollow = whatChar
		s.msg.Actor.SendGood("You follow " + whatChar.Name)
		whatChar.PartyFollowers = append(whatChar.PartyFollowers, s.actor)
		if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			s.msg.Participant.SendInfo(s.actor.Name, " follows you.")
		}
	}else{
		s.msg.Actor.SendBad("Who ya followin'??")
	}

	s.ok = true
}
