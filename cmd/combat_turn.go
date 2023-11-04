package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
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
	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything!")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}
	if s.actor.Tier < config.SpecialAbilityTier {
		s.msg.Actor.SendBad("You must be at least tier " + strconv.Itoa(config.SpecialAbilityTier) + " to use this skill.")
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
			return
		}

		if s.actor.Class == 6 && (s.actor.Placement != whatMob.Placement) {
			s.msg.Actor.SendBad("You are too far away to turn them.")
			return
		}
		s.actor.Victim = whatMob

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
			whatMob.AddThreatDamage(whatMob.Stam.Current, s.actor)
			whatMob.Stam.Current = 0
			data.StoreCombatMetric("turn_whole", 0, 0, whatMob.Stam.Max, 0, whatMob.Stam.Max, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			DeathCheck(s, whatMob)
			whatMob = nil
		} else if curChance >= 100 || turnRoll <= curChance {
			s.msg.Actor.SendInfo("Your faith pours into " + whatMob.Name + " and causes " + strconv.Itoa(whatMob.Stam.Current/2) + " damage.")
			s.msg.Observers.SendInfo(s.actor.Name + " turned " + whatMob.Name)
			whatMob.AddThreatDamage(whatMob.Stam.Current/2, s.actor)
			whatMob.Stam.Subtract(whatMob.Stam.Current / 2)
			data.StoreCombatMetric("turn_half", 0, 0, whatMob.Stam.Current, 0, whatMob.Stam.Current, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
		} else {
			s.msg.Actor.SendBad("You fail to turn the " + whatMob.Name + ".  They charge you!")
			whatMob.CurrentTarget = s.actor.Name
			whatMob.Placement = s.actor.Placement
			whatMob.AddThreatDamage(whatMob.Stam.Current, s.actor)
			stamDamage, vitDamage, resisted := s.actor.ReceiveDamage(s.actor.Stam.Max / 2)
			data.StoreCombatMetric("turn_fail_retaliate", 0, 0, stamDamage+vitDamage+resisted, resisted, stamDamage+vitDamage, 1, whatMob.MobId, whatMob.Level, 0, s.actor.CharId)
			s.msg.Observers.SendInfo(s.actor.Name + " turn attempt fails and enrages " + whatMob.Name)
		}
		return
	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}
