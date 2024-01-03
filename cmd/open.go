package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
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

	// Test for partial exit names
	exitTxt := strings.ToLower(strings.Join(s.words, " "))
	if !utils.StringIn(strings.ToUpper(exitTxt), directionals) {
		exitS := s.where.FindExit(exitTxt, s.actor)
		if exitS != nil {
			exitTxt = exitS.Name
		}
	}

	if what := s.where.FindExit(exitTxt, s.actor); what != nil {
		// Is item a door that can be close
		if !what.Flags["closeable"] {
			s.msg.Actor.SendBad(what.Name, " is already open.")
			return
		}

		if !what.Flags["closed"] {
			s.msg.Actor.SendInfo(utils.Title(what.Name), " is already open.")
			return
		}

		if what.Placement != s.actor.Placement {
			s.msg.Actor.SendBad("You must be next to it to open it.")
			return
		}

		s.actor.RunHook("act")

		what.Open()

		if s.actor.Flags["invisible"] == false {
			who := s.actor.Name
			s.msg.Actor.SendGood("You open ", what.Name, ".")
			s.msg.Observers.SendInfo(who, " opens ", what.Name, ".")
		} else {
			s.msg.Actor.SendGood("You open ", what.Name, ".")
		}

		return
	} else {
		s.msg.Actor.SendBad("You see no '", exitTxt, "' to open.")
		return
	}
}
