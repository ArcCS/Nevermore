package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
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

	if s.actor.GetStat("int") < config.IntMajorPenalty {
		s.msg.Actor.SendBad("You swing the item around wildly but can't figure out how to use it.")
		return
	}

	if s.actor.GetStat("int") < config.IntMinorPenalty {
		if utils.Roll(100, 1, 0) <= config.FizzleSave {
			s.msg.Actor.SendBad("You tried to invoke the item but it fizzled out.")
			s.actor.SetTimer("use", 8)
			return
		}
	}

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
		if what.Spell != "" && what.MaxUses > 0 {
			spellInstance, ok := objects.Spells[strings.ToLower(what.Spell)]
			if !ok {
				s.msg.Actor.SendBad("Spell doesn't exist in this world. ")
				return
			}
			if utils.StringIn(spellInstance.Name, objects.OffensiveSpells) && s.actor.Victim != nil {
				log.Println("offensive spell, victim is not nil")
				switch s.actor.Victim.(type) {
				case *objects.Character:
					name = s.actor.Victim.(*objects.Character).Name
				case *objects.Mob:
					name = s.actor.Victim.(*objects.Mob).Name
				}
			}
			if name != "" {
				var whatMob *objects.Mob
				if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
				} else {
					whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
				}
				// It was a mob!
				if whatMob != nil {
					s.actor.RunHook("use")
					s.actor.SetTimer("use", 8)
					msg = objects.Cast(s.actor, whatMob, spellInstance.Effect, spellInstance.Magnitude)
					s.msg.Actor.SendGood("You use a  " + what.Name + " on " + whatMob.Name)
					s.msg.Observers.SendGood(s.actor.Name + " used a " + what.Name + " on " + whatMob.Name)
					if strings.Contains(msg, "$CRIPT") {
						go Script(s.actor, strings.Replace(msg, "$CRIPT ", "", 1))
					} else if msg != "" {
						s.msg.Actor.SendGood(msg)
					}
					DeathCheck(s, whatMob)
					what.MaxUses -= 1
					if what.MaxUses <= 0 {
						s.msg.Actor.SendBad("Your " + what.Name + " disintegrates.")
						s.actor.Inventory.Remove(what)
					}
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
					s.actor.RunHook("use")
					s.actor.SetTimer("use", 8)
					msg = objects.Cast(s.actor, whatChar, spellInstance.Effect, spellInstance.Magnitude)
					s.msg.Actor.SendGood("You use a  " + what.Name + " on " + whatChar.Name)
					s.msg.Observers.SendGood(s.actor.Name + " used a " + what.Name + " on " + whatChar.Name)
					s.participant = whatChar
					s.msg.Participant.SendInfo(s.actor.Name + " used a " + what.Name + " on you")
					if strings.Contains(msg, "$CRIPT") {
						go Script(s.actor, strings.Replace(msg, "$CRIPT ", "", 1))
					} else if msg != "" {
						s.msg.Actor.SendGood(msg)
					}
					what.MaxUses -= 1
					if what.MaxUses <= 0 {
						s.msg.Actor.SendBad("Your " + what.Name + " disintegrates.")
						s.actor.Inventory.Remove(what)
					}
					return
				}
			} else {
				s.actor.RunHook("use")
				s.actor.SetTimer("use", 8)
				msg = objects.Cast(s.actor, s.actor, spellInstance.Effect, spellInstance.Magnitude)
				if strings.Contains(msg, "$CRIPT") {
					go Script(s.actor, strings.Replace(msg, "$CRIPT ", "", 1))
				} else {
					s.msg.Actor.SendGood(msg)
				}
				what.MaxUses -= 1
				if what.MaxUses <= 0 {
					s.msg.Actor.SendBad("Your " + what.Name + " disintegrates.")
					s.actor.Inventory.Remove(what)
				}
				return
			}
		}
	}

	s.msg.Actor.SendInfo("Use what?")
	s.ok = true
}
