package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(turn{},
		"Usage:  turn target # \n\n Channel the power of your faith into an undead target, instilling"+
			"fear and potentially destroying them completely.",
		permissions.Cleric|permissions.Paladin,
		"turn")
}

type turn cmd

func (turn) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Turn what exactly?")
		return
	}
	if s.actor.Tier < 5 {
		s.msg.Actor.SendBad("You aren't high enough level to perform that skill.")
		return
	}
	// Check some timers
	ready, msg := s.actor.TimerReady("combat_turn")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	if whatMob != nil {

		if whatMob.Flags["undead"] != true {
			s.msg.Actor.SendBad("Your target isn't undead!")
		}

		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to turn them.")
			return
		}

		s.actor.RunHook("combat")
		s.actor.SetTimer("combat_turn", config.TurnTimer)
		s.actor.SetTimer("combat", config.CombatCooldown)
		// base chance is 15% to hide
		curChance := config.TurnMax
		if whatMob.Level > s.actor.Tier {
			curChance -= config.TurnScaleDown * (whatMob.Level - s.actor.Tier)
		} else if s.actor.Tier > whatMob.Level {
			curChance += config.TurnScaleDown * (s.actor.Tier - whatMob.Level)
		}

		turnRoll := utils.Roll(100, 1, 0)
		if turnRoll <= config.DisintegrateChance {
			s.msg.Actor.SendInfo("Your faith overwhelms the " + whatMob.Name + " and utterly demolishes them.")
			s.msg.Observers.SendInfo(s.actor.Name + " disintegrates " + whatMob.Name)
			whatMob.Stam.Current = 0
			//TODO Calculate experience
			stringExp := strconv.Itoa(whatMob.Experience)
			for k := range whatMob.ThreatTable {
				s.where.Chars.Search(k, s.actor).Write([]byte(text.Cyan + "You earn " + stringExp + " exp for the defeat of the " + whatMob.Name + "\n" + text.Reset))
				s.where.Chars.Search(k, s.actor).Experience.Add(whatMob.Experience)
			}
			s.msg.Actor.SendInfo(whatMob.DropInventory())
			objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
			whatMob = nil
		}else if curChance >= 100 || turnRoll <= curChance {
			s.msg.Actor.SendInfo("Your faith pours into " + whatMob.Name + " and damages them!")
			s.msg.Observers.SendInfo(s.actor.Name + " turned " + whatMob.Name)
			whatMob.AddThreatDamage(whatMob.Stam.Current/2, s.actor.Name)
			whatMob.Stam.Subtract(whatMob.Stam.Current/2)
		}else{
			s.msg.Actor.SendBad("You fail to turn the " + whatMob.Name + ".  They charge you!")
			whatMob.CurrentTarget = s.actor.Name
			whatMob.AddThreatDamage(whatMob.Stam.Current, s.actor.Name)
			s.actor.ReceiveDamage(s.actor.Stam.Max)
			s.msg.Observers.SendInfo(s.actor.Name + " turn attempt fails and enrages " + whatMob.Name )
		}
		return
	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}
