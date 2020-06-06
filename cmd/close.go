package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

// Syntax: CLOSE <door>
func init() {
	addHandler(closeExit{},
           "Usage:  close exitName \n\n Close the specified exit so no one can pass through it.",
           permissions.Player,
           "CLOSE")
}

type closeExit cmd

func (closeExit) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to close?")
		return
	}

	name := s.words[0]

	// Search for item to close in the inventory where we are
	what := s.where.FindExit(name)

	// Was item to get found?
	if what == nil {
		s.msg.Actor.SendBad("You see no '", name, "' to close.")
		return
	}

	// Get item's proper name
	name = what.Name

	// Is item a door that can be closed
	if !what.Flags["closeable"] {
		s.msg.Actor.SendBad("You cannot close ", name, ".")
		return
	}

	if what.Flags["closed"] {
		s.msg.Actor.SendInfo(strings.ToTitle(name), " is already closed.")
		return
	}

	what.Close()

	if s.actor.Flags["invisible"] == false {
		who := s.actor.Name
		s.msg.Actor.SendGood("You close ", name, ".")
		s.msg.Observer.SendInfo(who, " closes ", name, ".")
	}else{
		s.msg.Actor.SendGood("You close ", name, ".")
	}


	s.ok = true
}
