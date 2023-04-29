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
	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything to cast a spell!!")
		return
	}

	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What do you want to cast?")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	if s.actor.CheckFlag("singing") {
		s.msg.Actor.SendBad("You can't cast a spell while singing!")
		return
	}

	spellInstance, ok := objects.Spells[strings.ToLower(s.input[0])]
	if !ok {
		s.msg.Actor.SendBad("What spell do you want to cast?")
		return
	}
	cost := spellInstance.Cost
	if s.actor.Class == 8 {
		cost = cost / 4
		if cost < 1 {
			cost = 1
		}
	}

	if !utils.StringIn(spellInstance.Name, s.actor.Spells) && s.actor.Class != 100 {
		s.msg.Actor.SendBad("You do not have that spell in your spellbook.")
		return
	}

	if s.actor.GetStat("int") < config.IntMinCast {
		s.msg.Actor.SendBad("You simply do not have the mental capacity to cast spells.")
		return
	}

	if s.actor.GetStat("int") < config.IntNoFizzle {
		if utils.Roll(100, 1, 0) <= config.FizzleSave {
			s.msg.Actor.SendBad("You attempt to cast the spell, but it fizzles out.")
			s.actor.Mana.Current -= cost
			s.actor.SetTimer("combat", 8)
			return
		}
	}

	if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		if minLevel, ok := spellInstance.Classes[config.AvailableClasses[s.actor.Class]]; !ok {
			s.msg.Actor.SendBad("The comprehension of this spell is beyond you.")
			return
		} else if s.actor.Tier < minLevel {
			s.msg.Actor.SendBad("You are not high enough level to cast this spell.")
			return
		}
	}

	if s.actor.Mana.Current < cost && s.actor.Class != 100 {
		s.msg.Actor.SendBad("You do not have enough mana to cast this spell. ")
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
		// Try Mobs First
		var whatMob *objects.Mob
		whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
		// It was a mob!
		if whatMob != nil {
			s.actor.RunHook("combat")
			s.actor.Victim = whatMob
			msg = objects.Cast(s.actor, whatMob, spellInstance.Effect, spellInstance.Magnitude)
			s.actor.Mana.Subtract(cost)
			s.actor.SetTimer("combat", 8)
			// TODO: At level 15 wizards can change  the chant to a different action to invoke spells,
			// bards simply stop chanting and can invoke spells at will.
			if (s.actor.Class == 7 || s.actor.Class == 4) && s.actor.Tier <= 15 {
				s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
				s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
			}
			s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on " + whatMob.Name)
			s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + whatMob.Name)
			if strings.Contains(msg, "$CRIPT") {
				go Script(s.actor, strings.Replace(msg, "$CRIPT ", "", 1))
			} else if msg != "" {
				s.msg.Actor.SendGood(msg)
			}
			DeathCheck(s, whatMob)
			s.ok = true
			return
		}

		// Are we casting on a character
		var whatChar *objects.Character

		// Check if it's a remote spell
		if utils.StringIn(spellInstance.Name, objects.RemoteSpells) {
			whatChar = objects.ActiveCharacters.Find(name)
			// It was a person!
			if whatChar != nil {
				s.participant = whatChar
				s.actor.RunHook("combat")
				s.actor.Victim = whatChar
				if strings.Contains(spellInstance.Effect, "damage") {
					//TODO PVP flags etc.
					s.msg.Actor.SendBad("No PVP implemented yet. ")
					s.ok = true
					return
				}
				if spellInstance.Name == "heal" {
					if s.actor.ClassProps["heals"] <= 0 {
						s.msg.Actor.SendBad("You cannot cast heal anymore today.")
						s.ok = true
						return
					} else {
						s.actor.ClassProps["heals"]--
					}
				}
				if spellInstance.Name == "restore" {
					if s.actor.ClassProps["restores"] <= 0 {
						s.msg.Actor.SendBad("You cannot cast restore anymore today.")
						s.ok = true
						return
					} else {
						s.actor.ClassProps["restore"]--
					}
				}
				msg = objects.Cast(s.actor, whatChar, spellInstance.Effect, spellInstance.Magnitude)
				s.actor.Mana.Subtract(cost)
				s.actor.SetTimer("combat", config.CombatCooldown)
				s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
				s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
				s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on " + whatChar.Name)
				if strings.Contains(msg, "$CRIPT") {
					go Script(s.actor, strings.Replace(msg, "$CRIPT ", "", 1))
				} else if msg != "" {
					s.msg.Participant.SendGood(msg)
				}
				s.ok = true
				return
			}
		}

		whatChar = s.where.Chars.Search(name, s.actor)
		// It was a person!
		if whatChar != nil {
			s.participant = whatChar
			s.actor.RunHook("combat")
			s.actor.Victim = whatChar
			if strings.Contains(spellInstance.Effect, "damage") {
				//TODO PVP flags etc.
				s.msg.Actor.SendBad("No PVP implemented yet. ")
				s.ok = true
				return
			}
			msg = objects.Cast(s.actor, whatChar, spellInstance.Effect, spellInstance.Magnitude)
			s.actor.Mana.Subtract(cost)
			s.actor.SetTimer("combat", config.CombatCooldown)
			s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
			s.msg.Participant.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
			s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
			s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on " + whatChar.Name)
			s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + whatChar.Name)
			s.msg.Participant.SendInfo(s.actor.Name + " cast a " + spellInstance.Name + " spell on you")
			if strings.Contains(msg, "$CRIPT") {
				go Script(s.actor, strings.Replace(msg, "$CRIPT ", "", 1))
			} else if msg != "" {
				s.msg.Participant.Send(msg)
			}
			s.ok = true
			return
		}

	} else {

		if strings.Contains(spellInstance.Effect, "damage") {
			s.msg.Actor.SendBad("You cannot cast this on yourself.")
			s.ok = true
			return
		}

		s.actor.RunHook("combat")
		msg = objects.Cast(s.actor, s.actor, spellInstance.Effect, spellInstance.Magnitude)
		s.actor.SetTimer("combat", 8)
		s.actor.Mana.Subtract(cost)
		s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
		s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
		if utils.StringIn(spellInstance.Name, objects.AmbiguousTargets) {
			s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell.")
			s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell.")
		} else {
			s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on yourself")
			s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + config.TextDescPronoun[s.actor.Gender] + "self.")
		}
		if strings.Contains(msg, "$CRIPT") {
			go Script(s.actor, strings.Replace(msg, "$CRIPT ", "", 1))
		} else if msg != "" {
			s.msg.Actor.SendGood(msg)
		}
		s.ok = true
		return

	}

	s.msg.Actor.SendInfo("Cast on who?")
	s.ok = true
}
