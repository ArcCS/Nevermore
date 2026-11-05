package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(sweep{},
		"Usage:  sweep [target #] \n\n Strike low to knock a target off its feet, stunning it and drawing its attention onto you.  Throws extra strikes while flurrying, though a flurried sweep builds no chi.  With no target named, you sweep the foe you are already fighting.",
		permissions.Monk,
		"sweep")
}

type sweep cmd

func (sweep) process(s *state) {
	if !hasMobTarget(s) {
		s.msg.Actor.SendBad("Sweep what exactly?")
		return
	}

	if !s.requireSight() {
		return
	}

	if !s.requireStamina() {
		return
	}

	if !s.requireTier(config.MonkSweepTier) {
		return
	}

	ready, msg := s.actor.TimerReady("combat_sweep")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	whatMob := mobTarget(s)
	if whatMob == nil {
		s.msg.Actor.SendInfo("Sweep what?")
		s.ok = true
		return
	}

	if s.actor.Placement != whatMob.Placement {
		s.msg.Actor.SendBad("You are too far away to sweep them.")
		return
	}

	s.actor.Victim = whatMob
	s.actor.RunHook("combat")

	skillLevel := attackSkillLevel(s)
	attacks, flurried := monkAttackPlan(s, skillLevel)
	result := resolveHits(s, whatMob, attacks, skillLevel, "sweep")
	if result.hits == 0 {
		s.msg.Observers.SendBad(s.actor.Name + " fails to sweep " + whatMob.Name)
		// A missed sweep still draws the mob's attention, with no threat bump.
		whatMob.AddThreatDamage(1, s.actor)
		whatMob.CurrentTarget = s.actor.Name
		s.actor.SetTimer("combat", config.CombatCooldown)
		return
	}

	if !flurried {
		gainChiFromHits(s, result.chiHits())
	}

	// A landed sweep always stuns, like circle and bash.
	whatMob.StunBy(config.SweepStuns, "sweep-stun", s.actor)
	// A light taunt: the mob turns on the monk without out-holding the real tanks.
	whatMob.AddThreatDamage(config.ThreatPercent(whatMob.Stam.Max, config.SmallThreatPercent), s.actor)
	whatMob.CurrentTarget = s.actor.Name
	s.msg.Actor.SendGood("You sweep " + whatMob.Name + " off its feet, and it turns its attention to you!")
	s.msg.Observers.SendInfo(s.actor.Name + " sweeps " + whatMob.Name + " off its feet, drawing its attention!")

	DeathCheck(s, whatMob)
	s.actor.SetTimer("combat_sweep", config.SweepTimer)
	s.actor.SetTimer("combat", config.CombatCooldown)
}
