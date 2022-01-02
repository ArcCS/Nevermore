package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

// Syntax: ( INVENTORY | INV )
func init() {
	addHandler(pick{},
		"Usage:  pick target \n \n Pick a lock.",
		permissions.Thief,
		"pick")
}

type pick cmd

func (pick) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Pick what?")
		return
	}

	targetStr := s.words[0]

	whatExit := s.where.FindExit(strings.ToLower(targetStr))
	if whatExit != nil {
			if whatExit.Placement != s.actor.Placement {
				s.msg.Actor.SendBad("You are too far away to pick  ", whatExit.Name)
				return
			}

			if !whatExit.Flags["locked"] {
				s.msg.Actor.SendBad("The door isn't locked.")
				return
			}

			if whatExit.Flags["unpickable"] {
				s.msg.Actor.SendBad("That lock cannot be picked.")
				return
			}

			whatExit.Flags["locked"] = false
			s.msg.Actor.SendGood("You successfully picked " + whatExit.Name)
			s.msg.Observers.SendInfo(s.actor.Name + " picked " + whatExit.Name)

		}else{
			s.msg.Actor.SendInfo("That item isn't on the target.")
			return
		}
	s.ok = true
}

