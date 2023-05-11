package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

// Syntax: ( EXPERIENCE | EXP )
func init() {
	addHandler(experience{},
		"Displays your character's experience points required for the next tier.",
		permissions.Player,
		"EXP", "EXPERIENCE", "XP")
}

type experience cmd

func (experience) process(s *state) {
	s.msg.Actor.SendGood("You require " + strconv.Itoa(config.TierExpLevels[s.actor.Tier+1]-s.actor.Experience.Value) + " additional experience pts for your next tier.")
	s.msg.Actor.SendGood("You are carying " + strconv.Itoa(s.actor.Gold.Value) + " gold marks")
	s.ok = true
}
