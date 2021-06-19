package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
	"strings"
)

func init() {
	addHandler(use{},
		"Usage:  use item # \n\n Use an item",
		permissions.Player,
		"USE")
}

type use cmd

func (use) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What do you want to use?")
		return
	}

	ready, msg := s.actor.TimerReady("use")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	itemName := s.words[0]
	itemNum := 1
	name := ""
	nameNum := 1

	if len(s.words) == 4 {
		name = s.words[2]
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			itemNum = val
		}
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[3]); err == nil {
			nameNum = val
		}
	} else if len(s.words) == 3 {
		name = s.words[2]
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			itemNum = val
		} else {
			// Try to snag a number off the list
			name = s.words[1]
			if val, err := strconv.Atoi(s.words[2]); err == nil {
				nameNum = val
			}
		}
	} else if len(s.words) == 2 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			itemNum = val
		} else {
			name = s.words[1]
		}
	}

	what := s.actor.Inventory.Search(itemName, itemNum)

	// It was on you the whole time
	if what != nil {
		s.actor.RunHook("use")
		//log.Println("Arrived here", what.Name, what.Spell)
		s.actor.SetTimer("use", 8)
		if what.Spell != "" && what.MaxUses > 1 {
			spellInstance, ok := objects.Spells[strings.ToLower(what.Spell)]
			if !ok {
				s.msg.Actor.SendBad("Spell doesn't exist in this world. ")
				return
			}
			log.Println("Arrived here", what.Name, what.Spell)
			if name != "" {
				var whatMob *objects.Mob
				if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
				} else {
					whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
				}
				// It was a mob!
				if whatMob != nil {
					msg = objects.Cast(s.actor, whatMob, spellInstance.Effect, spellInstance.Magnitude)
					s.msg.Actor.SendGood("You use a  " + what.Name + " on " + whatMob.Name)
					s.msg.Observers.SendGood(s.actor.Name + " used a " + what.Name + " on " + whatMob.Name)
					if strings.Contains(msg, "$CRIPT"){
						go Script(s.actor, strings.Replace(msg, "$CRIPT ", "",1))
					}else if msg != "" {
						s.msg.Actor.SendGood(msg)
					}
					go whatMob.DeathCheck(s.actor)
					return
				}

				// Are we casting on a character
				var whatChar *objects.Character
				whatChar = s.where.Chars.Search(name, s.actor)
				// It was a person!
				if whatChar != nil {
					if strings.Contains(spellInstance.Effect, "damage") {
						s.msg.Actor.SendBad("No PVP implemented yet. ")
						return
					}
					msg = objects.Cast(s.actor, whatChar, spellInstance.Effect, spellInstance.Magnitude)
					s.msg.Actor.SendGood("You use a  " + what.Name + " on " + whatChar.Name)
					s.msg.Observers.SendGood(s.actor.Name + " used a " + what.Name + " on " + whatChar.Name)
					s.participant = whatChar
					s.msg.Participant.SendInfo(s.actor.Name + " used a " + what.Name + " on you")
					if strings.Contains(msg, "$CRIPT"){
						go Script(s.actor, strings.Replace(msg, "$CRIPT ", "",1))
					}else if msg != "" {
						s.msg.Actor.SendGood(msg)
					}
					return
				}
			} else {
				msg = objects.Cast(s.actor, s.actor, spellInstance.Effect, spellInstance.Magnitude)
				if strings.Contains(msg, "$CRIPT"){
					go Script(s.actor, strings.Replace(msg, "$CRIPT ", "",1))
				}
				return
			}
		}
	}

	s.msg.Actor.SendInfo("Use what?")
	s.ok = true
}
