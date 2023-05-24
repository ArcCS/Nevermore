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
	log.Println("Quit Command: Kill prompt")
	s.actor.SetPromptStyle(objects.StyleNone)
	log.Println("Quit Command: Message user they are leaving")
	s.msg.Actor.SendGood("You leave this world behind.")
	log.Println("Quit Command: Save Character")
	s.actor.Save()
	log.Println("Quit Command: Clear Follows")
	s.actor.Unfollow()
	log.Println("Quit Command: Lose Party")
	s.actor.LoseParty()
	log.Println("Quit Command: Purge Effects")
	s.actor.PurgeEffects()
	log.Println("Quit Command: Clean Room")
	s.where.Chars.Remove(s.actor)
	log.Println("Quit Command: Active Character Removal")
	objects.ActiveCharacters.Remove(s.actor)

	s.ok = true
}
