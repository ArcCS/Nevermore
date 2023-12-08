package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(description{},
		"Usage:  description Your super fancy description \n \n If your description is empty you can set it the first time, otherwise after death you can change it while in the healing hand.",
		permissions.Player,
		"DESCRIPTION", "DESC")
}

type description cmd

func (description) process(s *state) {
	if len(s.words) <= 12 {
		s.msg.Actor.SendInfo("Your description should contain more than 12 words.")
		return
	}

	if s.actor.Description != "" && s.actor.ParentId != config.HealingHand {
		s.msg.Actor.SendInfo("You can only change your features with the expert help of the healing hand.")
		return
	}

	s.actor.Description = strings.Join(s.input, " ")

	s.msg.Actor.SendGood("You have set your description to: \n", s.actor.Description)

	s.ok = true
}
