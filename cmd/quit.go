package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/stats"
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
	s.actor.SetPromptStyle(objects.StyleNone)
	s.where.Chars.Remove(s.actor)
	s.msg.Actor.SendGood("You leave this world behind.")
	s.actor.Save()
	s.actor.Unfollow()
	s.actor.LoseParty()
	s.actor.PurgeEffects()
	stats.ActiveCharacters.Remove(s.actor)
	s.ok = true
}
