package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(search{},
		"Usage:  search \n\n Search the room for hidden things",
		permissions.Player,
		"search")
}

type search cmd

func (search) process(s *state) {
	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything!")
		return
	}

	ready, msg := s.actor.TimerReady("search")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	s.msg.Actor.SendInfo("You search about in the area.")

	if !s.actor.CheckFlag("hidden") && !s.actor.CheckFlag("invisible") {
		s.msg.Observers.SendInfo(s.actor.Name + " searches the area.")
	}
	s.actor.SetTimer("search", 16)

	// Use Int to determine the success roll for searching
	searchChance := s.actor.Int.Current * config.SearchPerInt

	// Look for hidden exits
	for _, exit := range s.where.Exits {
		if exit.Flags["hidden"] {
			if utils.Roll(100, 1, 0) <= searchChance {
				s.msg.Actor.SendGood("You find hidden exit: ", exit.Name, "!")
			}
		}
	}

	// Look for hidden people
	for _, char := range s.where.Chars.ListHiddenChars(s.actor) {
		if utils.Roll(100, 1, 0) <= searchChance {
			s.msg.Actor.SendGood("You find ", char.Name, " deftly hidden away from view!")
		}
	}

	// TODO: Should we make hidden items again?

	// Look for hidden mobs
	for _, mob := range s.where.Chars.ListHiddenChars(s.actor) {
		if utils.Roll(100, 1, 0) <= searchChance {
			s.msg.Actor.SendGood("You find ", mob.Name, " deftly hidden away from view!")
		}
	}

	s.ok = true
	return
}
