package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"math"
	"strconv"
	"strings"
)

func init() {
	addHandler(cast{},
		"Usage:  cast spell_name target # \n\n Attempts to cast a known spell from your spellbook",
		permissions.Mage|permissions.Bard|permissions.Cleric|permissions.Ranger|permissions.Paladin|permissions.Thief|permissions.Monk|permissions.Gamemaster|permissions.Dungeonmaster|permissions.Builder|permissions.Barbarian|permissions.Fighter,
		"cast", "ca", "c")
}

type cast cmd

// dailyCharges maps spells limited per day to the ClassProps counter that gates them.
var dailyCharges = map[string]string{
	"heal":    "heals",
	"restore": "restores",
}

// spendDailyCharge uses one daily charge if the spell is limited, telling the
// caster and returning false when none are left. Staff are never limited.
func spendDailyCharge(s *state, spell string) bool {
	prop, limited := dailyCharges[spell]
	if !limited || s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		return true
	}
	if s.actor.ClassProps[prop] <= 0 {
		s.msg.Actor.SendBad("You cannot cast " + spell + " anymore today.")
		return false
	}
	s.actor.ClassProps[prop]--
	return true
}

func (cast) process(s *state) {
	/* Allowing blind casting
	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything to cast a spell!!")
		return
	}
	*/

	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What do you want to cast?")
		return
	}

	if !s.requireStamina() {
		return
	}

	// Only bards can sing, and they may cast through it; see noteSingingCast.
	if s.actor.CheckFlag("singing") && s.actor.Class != config.BARD {
		s.msg.Actor.SendBad("You can't cast a spell while singing!")
		return
	}

	spellInstance, ok := objects.ResolveSpell(s.input[0])
	if !ok {
		s.msg.Actor.SendBad("What spell do you want to cast?")
		return
	}

	// Healing and status-removal spells only wait on the cast timer; everything
	// else also respects the global timer (e.g. blocked movement). Stun and paralysis block all casting.
	var ready bool
	var msg string
	if objects.IsRestorativeSpell(spellInstance.Name) {
		ready, msg = s.actor.TimerReadyIgnoreGlobal("cast")
	} else {
		ready, msg = s.actor.TimerReady("cast")
	}
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	cost := spellInstance.Cost
	if s.actor.Class == config.MONK {
		cost = cost / config.MonkSpellCostDivisor
		if cost < 1 {
			cost = 1
		}
	}

	if !utils.StringIn(spellInstance.Name, s.actor.Spells) && s.actor.Class != 100 {
		s.msg.Actor.SendBad("You do not have that spell in your spellbook.")
		return
	}

	// Clerics and paladins can lean on piety to steady divinity-based spells.
	castStat := s.actor.GetStat("int")
	if (s.actor.Class == config.CLERIC || s.actor.Class == config.PALADIN) && utils.StringIn(spellInstance.Name, objects.HealingSpells) {
		castStat = int(math.Max(float64(castStat), float64((castStat+s.actor.GetStat("pie"))/2)))
	}
	if castStat < config.IntMajorPenalty {
		var fizzlePenalty = int(math.Min(float64(config.FizzleSave*(config.IntMajorPenalty-castStat)), float64(75)))
		var roll = utils.Roll(100, 1, 0)
		if roll <= fizzlePenalty {
			s.msg.Actor.SendBad("You attempt to cast the spell, but it fizzles out.")
			s.actor.Mana.Current -= cost
			s.actor.SetTimer("cast", 8)
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

	if (s.actor.Class == 5 || s.actor.Class == 4 || s.actor.Class == 7) && utils.StringIn(spellInstance.Name, objects.OffensiveSpells) {
		// Make sure we check that the combat timer is ready as well if this spell is offensive.
		ready, msg := s.actor.TimerReady("combat")
		if !ready {
			s.msg.Actor.SendBad(msg)
			return
		}
	}

	if s.actor.Mana.Current < cost && s.actor.Class != 100 {
		s.msg.Actor.SendBad("You do not have enough mana to cast this spell. ")
		return
	}

	if objects.Rooms[s.actor.ParentId].Flags["no_magic"] {
		s.msg.Actor.SendBad("An oppressive anti-magic aura prevents you from casting magic here.")
		return
	}

	name := ""
	if len(s.words) > 1 {
		name = s.input[1]
	}
	nameNum := 1

	if len(s.words) > 2 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[2]); err == nil {
			nameNum = val
		}
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
			if strings.Contains(spellInstance.Effect, "damage") {
				s.msg.Actor.SendBad("Who are you trying to cast on?")
				s.ok = true
				return
			}
			if !spendDailyCharge(s, spellInstance.Name) {
				s.ok = true
				return
			}
			noteSingingCast(s)
			s.actor.FlagOn("casting", "cast")
			msg = objects.Cast(s.actor, whatChar, spellInstance.Effect, spellInstance.Magnitude)
			s.actor.FlagOff("casting", "cast")
			s.actor.Mana.Subtract(cost)
			if (s.actor.Class == 5 || s.actor.Class == 4 || s.actor.Class == 7) && utils.StringIn(spellInstance.Name, objects.OffensiveSpells) {
				s.actor.SetTimer("combat", config.CombatCooldown)
			}
			s.actor.SetTimer("cast", config.CombatCooldown)
			s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
			s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
			s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on " + whatChar.Name)
			if strings.Contains(msg, "$CRIPT") {
				go Script(whatChar, strings.Replace(msg, "$CRIPT ", "", 1))
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
		log.Println("Selected Character: ", whatChar.Name)
		if whatChar == s.actor && spellInstance.Name == "restore" {
			s.msg.Actor.SendBad("You can only cast this spell on others.")
			return
		}
		s.participant = whatChar
		s.actor.RunHook("combat")
		if strings.Contains(spellInstance.Effect, "damage") {
			//TODO PVP flags etc.
			s.msg.Actor.SendBad("Who are you trying to cast on?")
			s.ok = true
			return
		}
		if !spendDailyCharge(s, spellInstance.Name) {
			s.ok = true
			return
		}
		noteSingingCast(s)
		s.actor.FlagOn("casting", "cast")
		msg = objects.Cast(s.actor, whatChar, spellInstance.Effect, spellInstance.Magnitude)
		s.actor.FlagOff("casting", "cast")
		s.actor.Mana.Subtract(cost)
		if (s.actor.Class == 5 || s.actor.Class == 4 || s.actor.Class == 7) && utils.StringIn(spellInstance.Name, objects.OffensiveSpells) {
			s.actor.SetTimer("combat", config.CombatCooldown)
		}
		s.actor.SetTimer("cast", config.CombatCooldown)
		s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
		s.msg.Participant.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
		s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
		s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on " + whatChar.Name)
		s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + whatChar.Name)
		s.msg.Participant.SendInfo(s.actor.Name + " cast a " + spellInstance.Name + " spell on you")
		if strings.Contains(msg, "$CRIPT") {
			if whatChar == s.actor {
				go Script(s.actor, strings.Replace(msg, "$CRIPT ", "", 1))
			} else {
				go Script(whatChar, strings.Replace(msg, "$CRIPT ", "", 1))
			}
		} else if msg != "" {
			s.msg.Participant.Send(msg)
		}
		s.ok = true
		return
	}

	// Try Mobs Last
	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)

	if whatMob == nil && name == "" {
		if s.actor.Victim != nil && utils.StringIn(spellInstance.Name, objects.OffensiveSpells) {
			whatMob = s.actor.Victim.(*objects.Mob)
		}
	}

	// It was a mob!
	if whatMob != nil {
		s.actor.RunHook("combat")
		if utils.StringIn(spellInstance.Name, objects.OffensiveSpells) {
			s.actor.Victim = whatMob
			whatMob.MarkAttackedBy(s.actor)
		}
		if !spendDailyCharge(s, spellInstance.Name) {
			s.ok = true
			return
		}
		noteSingingCast(s)
		s.actor.FlagOn("casting", "cast")
		msg = objects.Cast(s.actor, whatMob, spellInstance.Effect, spellInstance.Magnitude)
		objects.SpellHitCallbacks(s.actor, spellInstance, whatMob)
		s.actor.FlagOff("casting", "cast")
		s.actor.Mana.Subtract(cost)
		if (s.actor.Class == 5 || s.actor.Class == 4 || s.actor.Class == 7) && utils.StringIn(spellInstance.Name, objects.OffensiveSpells) {
			s.actor.SetTimer("combat", 8)
		}
		s.actor.SetTimer("cast", 8)
		s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
		s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
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

	// A named target that matched nothing is an error, never a self cast.
	if name != "" {
		if utils.StringIn(spellInstance.Name, objects.RemoteSpells) {
			s.msg.Actor.SendBad("You cannot sense anyone named '" + name + "'.")
		} else {
			s.msg.Actor.SendBad("You don't see '" + name + "' here.")
		}
		return
	}

	if utils.StringIn(spellInstance.Name, objects.OffensiveSpells) {
		s.msg.Actor.SendBad("Who are you trying to cast on?")
		return
	}

	if spellInstance.Name == "restore" {
		s.msg.Actor.SendBad("You can only cast this spell on others.")
		return
	}

	log.Println("Casting on self")
	s.actor.RunHook("combat")
	if !spendDailyCharge(s, spellInstance.Name) {
		s.ok = true
		return
	}
	noteSingingCast(s)
	s.actor.FlagOn("casting", "cast")
	msg = objects.Cast(s.actor, s.actor, spellInstance.Effect, spellInstance.Magnitude)
	s.actor.FlagOff("casting", "cast")
	if (s.actor.Class == 5 || s.actor.Class == 4 || s.actor.Class == 7) && utils.StringIn(spellInstance.Name, objects.OffensiveSpells) {
		s.actor.SetTimer("combat", config.CombatCooldown)
	}
	s.actor.SetTimer("cast", config.CombatCooldown)
	s.actor.Mana.Subtract(cost)
	s.msg.Actor.SendGood("You chant: \"" + spellInstance.Chant + "\"")
	s.msg.Observers.SendGood(s.actor.Name + " chants: \"" + spellInstance.Chant + "\"")
	if utils.StringIn(spellInstance.Name, objects.AmbiguousTargets) {
		s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell.")
		s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell.")
	} else {
		s.msg.Actor.SendGood("You cast a " + spellInstance.Name + " spell on yourself")
		s.msg.Observers.SendGood(s.actor.Name + " cast a " + spellInstance.Name + " spell on " + config.DescPronoun(s.actor.Gender) + "self.")
	}
	if strings.Contains(msg, "$CRIPT") {
		go Script(s.actor, strings.Replace(msg, "$CRIPT ", "", 1))
	} else if msg != "" {
		s.msg.Actor.SendGood(msg)
	}
	s.ok = true
	return

}

// noteSingingCast warns a low-tier bard that their song is weakening the
// spell they are about to cast.
func noteSingingCast(s *state) {
	if s.actor.CheckFlag("singing") && s.actor.SingCastMod() < 1 {
		s.msg.Actor.SendInfo("Your song strains your concentration, weakening the spell.")
	}
}
