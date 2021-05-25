package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/spells"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
	"strings"
)

func init() {
	addHandler(cast{},
		"Usage:  cast spell_name target # \n\n Attempts to cast a known spell from your spellbook",
		permissions.Player,
		"cast")
}

type cast cmd

func (cast) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What do you want to cast?")
		return
	}

	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
	}
	if len(s.words) > 1 {
		name := s.input[1]
		nameNum := 1

		if len(s.words) > 2 {
			// Try to snag a number off the list
			if val, err := strconv.Atoi(s.words[2]); err == nil {
				nameNum = val
			}
		}

		spellInstance, ok := spells.Spells[strings.ToLower(s.input[0])]
		if !ok {
			s.msg.Actor.SendBad("What spell do you want to cast?")
			return
		}

		if !utils.StringIn(spellInstance.Name, s.actor.Spells) && s.actor.Class != 100 {
			s.msg.Actor.SendBad("You do not have that spell in your spellbook.")
			return
		}

		s.actor.RunHook("combat")

		// Try Mobs First
		var whatMob *objects.Mob
		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			whatMob = s.where.Mobs.Search(name, nameNum, true)
		} else {
			whatMob = s.where.Mobs.Search(name, nameNum, false)
		}
		// It was a mob!
		if whatMob != nil {
			if s.actor.Mana.Current > spellInstance.Cost || s.actor.Class == 100 {
				msg = spells.PlayerCast(s.actor, whatMob, spellInstance.Effect, map[string]interface{}{"magnitude": spellInstance.Magnitude})
				s.actor.SetTimer("combat", 8)
				s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
				s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
				s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on " + whatMob.Name)
				s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + whatMob.Name)
				s.msg.Actor.SendGood(msg)
				if whatMob.Stam.Current <= 0 {
					s.msg.Actor.SendInfo("You killed " + whatMob.Name + text.Reset)
					s.msg.Observers.SendInfo(s.actor.Name + " killed " + whatMob.Name + text.Reset)
					//TODO Calculate experience
					stringExp := strconv.Itoa(whatMob.Experience)
					for k := range whatMob.ThreatTable {
						s.where.Chars.Search(k, true).Write([]byte(text.Cyan + "You earn " + stringExp + " exp for the defeat of the " + whatMob.Name + "\n" + text.Reset))
						s.where.Chars.Search(k, true).Experience.Add(whatMob.Experience)
					}
					s.msg.Observers.SendInfo(whatMob.Name + " dies.")
					s.msg.Actor.SendInfo(whatMob.DropInventory())
					objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
					whatMob = nil
				}
				return
			} else {
				s.msg.Actor.SendBad("You do not have enough mana to cast this spell.")
			}
		}

		// Are we casting on a character
		var whatChar *objects.Character
		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			whatChar = s.where.Chars.Search(name, true)
		} else {
			whatChar = s.where.Chars.Search(name, false)
		}
		// It was a person!
		if whatChar != nil {
			if strings.Contains(spellInstance.Effect, "damage") {
				//TODO PVP flags etc.
				s.msg.Actor.SendBad("No PVP implemented yet. ")
				return
			}
			msg = spells.PlayerCast(s.actor, whatChar, spellInstance.Effect, map[string]interface{}{"magnitude": spellInstance.Magnitude})
			s.actor.SetTimer("combat", 8)
			s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
			s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
			s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on " + whatChar.Name)
			s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + whatChar.Name)
			s.participant = whatChar
			s.msg.Participant.SendInfo(s.actor.Name + " cast a " + spellInstance.Name + "spell on you")
			s.msg.Actor.SendGood(msg)
			return
		}
	} else {

		spellInstance, ok := spells.Spells[strings.ToLower(s.input[0])]
		if !ok {
			s.msg.Actor.SendBad("What spell do you want to cast?")
			return
		}

		if !utils.StringIn(spellInstance.Name, s.actor.Spells) && s.actor.Class != 100 {
			s.msg.Actor.SendBad("You do not have that spell in your spellbook.")
			return
		}

		if strings.Contains(spellInstance.Effect, "damage") {
			s.msg.Actor.SendBad("You cannot cast this on yourself.")
			return
		}

		msg = spells.PlayerCast(s.actor, s.actor, spellInstance.Effect, map[string]interface{}{"magnitude": spellInstance.Magnitude})
		s.actor.SetTimer("combat", 8)
		s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
		s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
		s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on yourself")
		s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + config.TextDescPronoun[s.actor.Gender] + "self.")
		s.msg.Actor.SendGood(msg)
		return
	}

	s.msg.Actor.SendInfo("Cast on who?")
	s.ok = true
}
