package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

// Syntax: DROP item
func init() {
	addHandler(roll{},
		"Usage:  roll sides num_dice \n \n Roll a number of specified sided dice",
		permissions.Player,
		"roll")
}

type roll cmd

func (roll) process(s *state) {
	rollSides := 20
	rollDice := 1
	if len(s.words) > 0 {
		rollSides, _ = strconv.Atoi(s.words[0])
	}
	if len(s.words) > 1 {
		rollDice, _ = strconv.Atoi(s.words[1])
	}

	dVal := utils.Roll(rollSides, rollDice, 0)

	s.msg.Actor.SendGood("You rolled: " + strconv.Itoa(dVal))
	s.msg.Observers.SendGood(s.actor.Name + " rolled: " + strconv.Itoa(dVal))
	s.ok = true
}
