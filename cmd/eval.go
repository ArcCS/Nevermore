package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(evaluate{},
		"Usage:  evaluate target\n\n  Examine a monster or item to find it's properties. ",
		permissions.Anyone,
		"evaluate", "eval")
}

type evaluate cmd

func (evaluate) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What do you want to evaluate?")
		return
	}

	if s.actor.Evals <= 0 {
		s.msg.Actor.SendBad("You cannot perform anymore evaluations today.")
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

	// Check mobs
	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	// It was a mob!
	if whatMob != nil {
		s.actor.Evals -= 1
		s.msg.Actor.SendInfo(whatMob.Eval())
		return
	}

	// Check items
	whatItem := s.where.Items.Search(name, nameNum)

	// Item in the room?
	if whatItem != nil {
		s.actor.Evals -= 1
		s.msg.Actor.SendInfo(whatItem.Eval())
		return
	}

	whatItem = s.actor.Inventory.Search(name, nameNum)

	// It was on you the whole time
	if whatItem != nil {
		s.actor.Evals -= 1
		s.msg.Actor.SendInfo(whatItem.Eval())
		return
	}

	whatItem = s.actor.Equipment.Search(name)

	// Check your equipment
	if whatItem != nil {
		s.actor.Evals -= 1
		s.msg.Actor.SendInfo(whatItem.Eval())
		return
	}

	s.ok = true
	s.msg.Actor.SendInfo("Could not find anything to evaluate based on your input.")
	return

}
