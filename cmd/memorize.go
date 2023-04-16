package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(memorize{},
		"Usage:  memorize item_name # \n \n Memorize the sheet music so you can sing it later",
		permissions.Bard,
		"MEMORIZE", "MEM")
}

type memorize cmd

func (memorize) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("You need to specify sheet music to read from")
		return
	}
	s.ok = true

	name := s.words[0]
	nameNum := 1

	// Try searching inventory where we are
	what := s.actor.Inventory.Search(name, nameNum)

	// Was item to read found?
	if what == nil {
		s.msg.Actor.SendBad("You couldn't find anything like that to memorize.")
		return
	} else {
		if what.ItemType == 18 {
			if what.Spell == "" {
				s.msg.Actor.SendBad("You read the music but find that it contains no discernible music of value.")
				return
			}
			_, ok := objects.Songs[what.Spell]
			if !ok {
				s.msg.Actor.SendBad("The spell contained does not exist in this world.")
				return
			}
			if utils.StringIn(what.Spell, s.actor.Spells) {
				s.msg.Actor.SendBad("You already know this spell.")
				return
			}
			s.msg.Actor.SendGood("You study ", what.Name, " and learn the song "+what.Spell)
			s.actor.Spells = append(s.actor.Spells, what.Spell)
			s.msg.Observers.SendInfo("You see ", s.actor.Name, " memorize a ", name, ".")
			s.actor.Inventory.Remove(what)
			s.msg.Actor.SendInfo("The " + what.Name + " disintegrates.")
			return
		} else {
			s.msg.Actor.SendBad("That's not sheet music.")
		}
	}
}
