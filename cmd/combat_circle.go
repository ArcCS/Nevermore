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
	addHandler(circle{},
		"Usage:  circle target # \n\n Try to circle a mob to apply a short duration stun and generate threat.  Fighters stun for longer as their weapon skill rises.",
		permissions.Fighter|permissions.Barbarian,
		"circle", "cir")
}

type circle cmd

func (circle) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Circle what exactly?")
		return
	}

	if !s.requireSight() {
		return
	}

	if !s.requireStamina() {
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("combat_circle")
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

		s.actor.RunHook("combat")
		s.actor.Victim = whatMob
		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == nil {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		// Shortcut Missile weapon:
		if s.actor.Equipment.Main.ItemType == 4 {
			s.msg.Actor.SendBad("You cannot circle with a ranged weapon.")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to circle them.")
			return
		}

		// Check for a miss
		if utils.Roll(100, 1, 0) <= CircleMissChance(s, whatMob.Level-s.actor.Tier, whatMob.EngagedCount(s.actor)) {
			s.msg.Actor.SendBad("You missed!!")
			data.StoreCombatMetric("circle-stun-miss", 0, 0, 0, 0, 0, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			s.actor.SetTimer("combat_circle", config.CircleTimer)
			s.actor.SetTimer("combat", config.CombatCooldown)
			whatMob.AddThreatDamage(1, s.actor)
			whatMob.CurrentTarget = s.actor.Name
			s.msg.Observers.SendBad(s.actor.Name + " fails to circle " + whatMob.Name)
			return
		}

		whatMob.StunBy(config.CircleStunFor(s.actor.Class, attackSkillLevel(s)), "circle-stun", s.actor)
		whatMob.AddThreatDamage(config.ThreatPercent(whatMob.Stam.Max, config.LargeThreatPercent), s.actor)
		whatMob.CurrentTarget = s.actor.Name
		s.actor.SetTimer("combat_circle", config.CircleTimer)
		s.actor.SetTimer("combat", config.CombatCooldown)
		s.msg.Actor.SendInfo("You circled " + whatMob.Name)
		s.msg.Observers.SendInfo(s.actor.Name + " circles " + whatMob.Name)
		return

	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}

// CircleMissChance is the weapon miss chance without the reckless bonus.
// Circle never pays the reckless stance's stamina cost, so it does not get
// the stance's to-hit bonus either. The bonus is cancelled in the base so
// the 5-95 clamp still applies to the final number.
func CircleMissChance(s *state, lvlDiff int, engaged int) int {
	skill := 5
	if s.actor.Class != config.MONK {
		skill = s.actor.Equipment.Main.ItemType
	}
	base := config.WeaponMissChance(s.actor.Skills[skill].Value)
	if s.actor.CheckFlag("reckless") {
		base += config.RecklessMissReduction
	}
	return missChance(s, base, lvlDiff, engaged)
}
