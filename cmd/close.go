package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strings"
)

// Syntax: CLOSE <door>
func init() {
	addHandler(closeExit{},
		"Usage:  close exitName \n\n Close the specified exit so no one can pass through it.",
		permissions.Player,
		"CLOSE", "CL", "CLO")
}

type closeExit cmd

func (closeExit) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to close?")
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
			s.msg.Actor.SendBad(what.Name, " cannot be closed")
			return
		}

		if what.Flags["closed"] {
			s.msg.Actor.SendInfo(utils.Title(what.Name), " is already closed.")
			return
		}

		if what.Placement != s.actor.Placement {
			s.msg.Actor.SendBad("You must be next to it to close it.")
			return
		}

		s.actor.RunHook("act")

		what.Close()

		if s.actor.Flags["invisible"] == false {
			who := s.actor.Name
			s.msg.Actor.SendGood("You close ", what.Name, ".")
			s.msg.Observers.SendInfo(who, " closes ", what.Name, ".")
		} else {
			s.msg.Actor.SendGood("You close ", what.Name, ".")
		}

		return
	} else {
		s.msg.Actor.SendBad("You see no '", exitTxt, "' to close.")
		return
	}
}
