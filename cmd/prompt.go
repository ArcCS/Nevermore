package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(prompt{},
		"Usage:  prompt \n \n Change the prompt style between 'basic' and 'status'",
		permissions.Player,
		"prompt")
}

type prompt cmd

func (prompt) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Change prompt to what")
		return
	}

	if s.words[0] == "BASIC" {
		s.actor.SetPromptStyle(objects.StyleNone)
		s.msg.Actor.SendGood("Style Changed to Basic")
	} else if s.words[0] == "STATUS" {
		s.actor.SetPromptStyle(objects.StyleStat)
		s.msg.Actor.SendGood("Style Changed to Status")
	} else {
		s.msg.Actor.SendBad("Style not changed")
	}

	s.ok = true
	return
}
