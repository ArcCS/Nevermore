package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
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
		return
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

		spellInstance, ok := objects.Spells[strings.ToLower(s.input[0])]
		if !ok {
			s.msg.Actor.SendBad("What spell do you want to cast?")
			return
		}

		if !utils.StringIn(spellInstance.Name, s.actor.Spells) && s.actor.Class != 100 {
			s.msg.Actor.SendBad("You do not have that spell in your spellbook.")
			return
		}
		if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			if minLevel, ok := spellInstance.Classes[config.AvailableClasses[s.actor.Class]]; !ok {
				s.msg.Actor.SendBad("The comprehension of this spell is beyond you.")
				return
			} else if s.actor.Tier < minLevel {
				s.msg.Actor.SendBad("You are not high enough level to cast this spell.")
			}
		}

		s.actor.RunHook("combat")

		// Try Mobs First
		var whatMob *objects.Mob
		whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
		// It was a mob!
		if whatMob != nil {
			s.actor.Victim = whatMob
			cost := spellInstance.Cost
			if s.actor.Class == 8 {
				cost /= 4
			}
			if s.actor.Mana.Current >= cost || s.actor.Class == 100 {
				s.actor.Mana.Subtract(cost)
				msg = objects.Cast(s.actor, whatMob, spellInstance.Effect, spellInstance.Magnitude)
				s.actor.SetTimer("combat", 8)
				s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
				s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
				s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on " + whatMob.Name)
				s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + whatMob.Name)
				if strings.Contains(msg, "$CRIPT"){
					go Script(s.actor, strings.Replace(msg, "$CRIPT ", "",1))
				}else if msg != "" {
					s.msg.Actor.SendGood(msg)
				}
				go whatMob.DeathCheck(s.actor)
				s.ok=true
				return
			} else {
				s.msg.Actor.SendBad("You do not have enough mana to cast this spell.")
			}
		}

		// Are we casting on a character
		var whatChar *objects.Character
		whatChar = s.where.Chars.Search(name, s.actor)
		// It was a person!
		if whatChar != nil {
			s.actor.Victim = whatChar
			if strings.Contains(spellInstance.Effect, "damage") {
				//TODO PVP flags etc.
				s.msg.Actor.SendBad("No PVP implemented yet. ")
				s.ok=true
				return
			}
			msg = objects.Cast(s.actor, whatChar, spellInstance.Effect, spellInstance.Magnitude)
			s.actor.SetTimer("combat", config.CombatCooldown)
			s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
			s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
			s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on " + whatChar.Name)
			s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + whatChar.Name)
			s.participant = whatChar
			s.msg.Participant.SendInfo(s.actor.Name + " cast a " + spellInstance.Name + "spell on you")
			if strings.Contains(msg, "$CRIPT"){
				go Script(s.actor, strings.Replace(msg, "$CRIPT ", "",1))
			}else if msg != "" {
				s.msg.Actor.SendGood(msg)
			}
			s.ok=true
			return
		}
	} else {

		spellInstance, ok := objects.Spells[strings.ToLower(s.input[0])]
		if !ok {
			s.msg.Actor.SendBad("What spell do you want to cast?")
			s.ok=true
			return
		}

		if !utils.StringIn(spellInstance.Name, s.actor.Spells) && s.actor.Class != 100 {
			s.msg.Actor.SendBad("You do not have that spell in your spellbook.")
			s.ok=true
			return
		}

		if strings.Contains(spellInstance.Effect, "damage") {
			s.msg.Actor.SendBad("You cannot cast this on yourself.")
			s.ok=true
			return
		}

		msg = objects.Cast(s.actor, s.actor, spellInstance.Effect, spellInstance.Magnitude)
		s.actor.SetTimer("combat", 8)
		s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
		s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
		s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on yourself")
		s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + config.TextDescPronoun[s.actor.Gender] + "self.")
		if strings.Contains(msg, "$CRIPT"){
			go Script(s.actor, strings.Replace(msg, "$CRIPT ", "",1))
		}else if msg != "" {
			s.msg.Actor.SendGood(msg)
		}
		s.ok=true
		return

	}

	s.msg.Actor.SendInfo("Cast on who?")
	s.ok = true
}
