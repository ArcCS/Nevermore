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

	who := s.actor.Name

	s.actor.SetPromptStyle(objects.StyleNone)
	if s.actor.Flags["invisible"] == false && s.actor.Flags["hidden"] == false {
		s.msg.Observer.SendInfo(who, " vanishes in a puff of smoke.")
	}
	s.where.Chars.Remove(s.actor)
	s.msg.Actor.SendGood("You leave this world behind.")
	s.actor.Save()
	stats.ActiveCharacters.Remove(s.actor)
	s.ok = true
}
