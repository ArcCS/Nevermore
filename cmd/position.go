package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

// Overloaded Look object for all of your looking pleasure
// Syntax: ( LOOK | L ) has.Thing
func init() {
	addHandler(position{},
		"Usage:  position \n \n Where you are in the room. ",
		permissions.Player,
		"pos", "position")
}

type position cmd

func (position) process(s *state) {
	locationStr := ""
	switch s.actor.Placement {
	case 5: locationStr = "You are standing at the front of the room."
	case 4: locationStr = "You are standing toward the front of the room."
	case 3: locationStr = "You are standing in the center of the room."
	case 2: locationStr = "You are standing toward the back of the room."
	case 1: locationStr = "You are standing at the back of the room."
	}

	s.msg.Actor.SendInfo(locationStr)
	s.ok = true

	return
}
