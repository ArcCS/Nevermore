package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(read{},
		"Usage:  read item_name # \n \n Read the specified scroll into your spellbook",
		permissions.Player,
		"READ", "STUDY")
}

type read cmd

func (read) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("You need to specify a scroll to read from")
		return
	}
	s.ok=true

	name := s.words[0]
	nameNum := 1

	// Try searching inventory where we are
	what := s.actor.Inventory.Search(name, nameNum)

	// Was item to read found?
	if what == nil {
		s.msg.Actor.SendBad("You couldn't find anything like that to study.")
		return
	}else{
		if what.ItemType == 7 {
			if what.Spell == "" {
				s.msg.Actor.SendBad("You study the scroll but find that it contains no spell.")
				return
			}
			spellInstance, ok := objects.Spells[what.Spell]
			if !ok {
				s.msg.Actor.SendBad("The spell contained does not exist in this world.")
				return
			}
			if minLevel, ok :=  spellInstance.Classes[config.AvailableClasses[s.actor.Class]]; !ok {
				s.msg.Actor.SendBad("The comprehension of this spell is beyond you.")
				return
			}else if s.actor.Tier < minLevel {
				s.msg.Actor.SendBad("You are not high enough level to learn this spell.")
				return
			}
			if utils.StringIn(what.Spell, s.actor.Spells){
				s.msg.Actor.SendBad("You already know this spell.")
				return
			}
			s.msg.Actor.SendGood("You study ", what.Name, " and learn the spell " + what.Spell)
			s.actor.Spells = append(s.actor.Spells, what.Spell)
			s.msg.Observers.SendInfo("You see ", s.actor.Name, " study a ", name, ".")
			s.actor.Inventory.Remove(what)
			s.msg.Actor.SendInfo("The " + what.Name + " disintegrates.")
			return
		}else{
			s.msg.Actor.SendBad("That's not a scroll.")
		}
	}
}
