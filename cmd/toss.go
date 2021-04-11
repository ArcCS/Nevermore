package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

// Syntax: JUNK item
func init() {
	addHandler(toss{},
		"Usage:  toss itemName # \n \n Toss an item away, this is a permanent deletion.",
		permissions.Player,
		"toss")
}

type toss cmd

func (j toss) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What do you want to toss?")
		return
	}

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	what := s.actor.Inventory.Search(name, nameNum)

	// Still not found?
	if what == nil {
		s.msg.Actor.SendBad("You don't have a  '", name, "' to throw away.")
		return
	}

	s.actor.Inventory.Lock()
	s.actor.Inventory.Remove(what)
	s.actor.Inventory.Unlock()

	s.msg.Actor.SendGood("You toss away ", what.Name, ".")
	s.msg.Observers.SendInfo("You see ", s.actor.Name, " toss away ", what.Name, ".")
	what = nil
	s.ok = true
}
