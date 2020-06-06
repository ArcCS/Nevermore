package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(read{},
           "Usage:  read item_name # \n \n Read the specified scroll into your spellbook",
           permissions.Player,
           "READ")
}

type read cmd

func (read) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("You need to specify a scroll to read from")
		return
	}

	name := s.words[0]
	nameNum := 1


	// Try searching inventory where we are
	what := s.actor.Inventory.Search(name, nameNum)

	// Was item to read found?
	if what == nil {
		s.msg.Actor.SendBad("You have no '", name, "' to read.")
		return
	}

	// Get item's proper name
	name = what.Name

	// Make sure that this is a scroll

	s.msg.Actor.Send("You study ", name, " and learn ")

	who := s.actor.Name
	s.msg.Observer.SendInfo("You see ", who, " read ", name, ".")

	s.ok = true
}
