package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(open{},
		"Usage:  open exit_name \n \n Open the specified exit.",
		permissions.Player,
		"OPEN")
}

type open cmd

func (open) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to open?")
		return
	}

	name := s.words[0]

	// Search for item to close in the inventory where we are
	what := s.where.FindExit(name)

	// Was item to get found?
	if what == nil {
		s.msg.Actor.SendBad("You see no '", name, "' to open.")
		return
	}

	// Get item's proper name
	name = what.Name

	// Is item a door that can be close
	if !what.Flags["closeable"] {
		s.msg.Actor.SendBad(name, " is already open.")
		return
	}

	if !what.Flags["closed"] {
		s.msg.Actor.SendInfo(strings.ToTitle(name), " is already open.")
		return
	}

	s.actor.RunHook("act")

	what.Open()

	if s.actor.Flags["invisible"] == false {
		who := s.actor.Name
		s.msg.Actor.SendGood("You open ", name, ".")
		s.msg.Observers.SendInfo(who, " opens ", name, ".")
	} else {
		s.msg.Actor.SendGood("You open ", name, ".")
	}

	return
}
