package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(afk{},
		"Usage:  afk \n \n While OOC, set yourself AFK.  This extends the timer for when you will be logged out to 45 minutes.",
		permissions.Player,
		"afk")
}

type afk cmd

func (afk) process(s *state) {

	if s.actor.Flags["ooc"] && s.actor.Flags["afk"] {
		s.actor.Flags["afk"] = false
		s.msg.Actor.SendInfo("You are no longer AFK.")
	} else if s.actor.Flags["ooc"] && !s.actor.Flags["afk"] {
		s.actor.Flags["afk"] = true
		s.msg.Actor.SendInfo("You are now AFK.")
	}else{
		s.msg.Actor.SendBad("You must be OOC to use this command.")
	}
	s.ok = true
	return
}
