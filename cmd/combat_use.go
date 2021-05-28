package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
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
					msg = objects.PlayerCast(s.actor, whatMob, spellInstance.Effect, map[string]interface{}{"magnitude": spellInstance.Magnitude})
					s.msg.Actor.SendGood("You use a  " + what.Name + " on " + whatMob.Name)
					s.msg.Observers.SendGood(s.actor.Name + " used a " + what.Name + " on " + whatMob.Name)
					s.msg.Actor.SendGood(msg)
					if whatMob.Stam.Current <= 0 {
						s.msg.Actor.SendInfo("You killed " + whatMob.Name + text.Reset)
						s.msg.Observers.SendInfo(s.actor.Name + " killed " + whatMob.Name + text.Reset)
						//TODO Calculate experience
						stringExp := strconv.Itoa(whatMob.Experience)
						for k := range whatMob.ThreatTable {
							s.where.Chars.Search(k, s.actor).Write([]byte(text.Cyan + "You earn " + stringExp + " exp for the defeat of the " + whatMob.Name + "\n" + text.Reset))
							s.where.Chars.Search(k, s.actor).Experience.Add(whatMob.Experience)
						}
						s.msg.Observers.SendInfo(whatMob.Name + " dies.")
						s.msg.Actor.SendInfo(whatMob.DropInventory())
						objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
						whatMob = nil
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
					msg = objects.PlayerCast(s.actor, whatChar, spellInstance.Effect, map[string]interface{}{"magnitude": spellInstance.Magnitude})
					s.msg.Actor.SendGood("You use a  " + what.Name + " on " + whatChar.Name)
					s.msg.Observers.SendGood(s.actor.Name + " used a " + what.Name + " on " + whatChar.Name)
					s.msg.Actor.SendGood(msg)
					s.participant = whatChar
					s.msg.Participant.SendInfo(s.actor.Name + " used a " + what.Name + " on you")
					s.msg.Actor.SendGood(msg)
					return
				}
			} else {
				msg = objects.PlayerCast(s.actor, s.actor, spellInstance.Effect, map[string]interface{}{"magnitude": spellInstance.Magnitude})
				//s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on yourself")
				//s.msg.Actor.SendGood(msg)
				return
			}
		}
	}

	s.msg.Actor.SendInfo("Use what?")
	s.ok = true
}
