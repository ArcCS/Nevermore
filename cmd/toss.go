package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
)

func init() {
	addHandler(toss{},
		"Usage:  toss itemName # \n \n Toss an item away, this is a permanent deletion.",
		permissions.Player,
		"toss")
	addHandler(confirmToss{},
		"",
		permissions.Player,
		"$confirm_toss")

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

	s.msg.Actor.SendGood("Are you sure you want to toss away ", what.Name, "?  This cannot be undone. (y/n)")
	s.actor.AddCommands("yes", "$confirm_toss "+name+" "+strconv.Itoa(nameNum))
	s.actor.AddCommands("y", "$confirm_toss "+name+" "+strconv.Itoa(nameNum))
	what = nil
	s.ok = true
}

type confirmToss cmd

func (j confirmToss) process(s *state) {
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

	if err := s.actor.Inventory.Remove(what); err != nil {
		s.msg.Actor.SendBad("You can't remove that item from your inventory.")
		log.Println("Error removing item from inventory: ", err)
		return
	}

	s.msg.Actor.SendGood("You toss away ", what.Name, ".")
	s.msg.Observers.SendInfo("You see ", s.actor.Name, " toss away ", what.Name, ".")
	what = nil
	s.ok = true
}
