package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
)

func init() {
	addHandler(quit{},
		"Usage:  quit \n \n GTFO ",
		permissions.Anyone,
		"QUIT")
}

type quit cmd

// The Quit command acts as a hook for processing - such as cleanup - to be
// done when a player quits the game.
func (quit) process(s *state) {
	log.Println("Quit Command: Processing Quit")
	s.actor.SetPromptStyle(objects.StyleNone)
	s.msg.Actor.SendGood("You leave this world behind.")
	s.where.Chars.Remove(s.actor)
	s.actor.PrepareUnload()

	s.ok = true
}
