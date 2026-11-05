package cmd

import (
	"math"
	"strconv"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
)

func init() {
	addHandler(leap{},
		"Usage:  leap target # \n\n Close the distance to a target in a single bound and strike.  A short leap costs nothing; a seasoned monk can spend chi to cross the whole room.  A leap is always a single strike and always builds chi.",
		permissions.Monk,
		"leap", "leapstrike", "ls")
}

type leap cmd

func (leap) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Leap at what exactly?")
		return
	}

	if !s.requireSight() {
		return
	}

	if !s.requireStamina() {
		return
	}

	if !s.requireTier(config.MonkLeapTier) {
		return
	}

	if rooted, msg := s.actor.Rooted(); rooted {
		s.msg.Actor.SendBad(msg)
		return
	}

	if config.LeapTimer > 0 {
		ready, msg := s.actor.TimerReady("combat_leap")
		if !ready {
			s.msg.Actor.SendBad(msg)
			return
		}
	}
	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	name := s.input[0]
	nameNum := 1
	if len(s.words) > 1 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	whatMob := s.where.Mobs.Search(name, nameNum, s.actor)
	if whatMob == nil {
		s.msg.Actor.SendInfo("Leap at what?")
		s.ok = true
		return
	}

	distance := int(math.Abs(float64(s.actor.Placement - whatMob.Placement)))
	if distance == 0 {
		s.msg.Actor.SendBad("You are too close to leap at them.")
		return
	}
	if distance > config.LeapFreeDistance {
		if s.actor.Tier < config.MonkLongLeapTier {
			s.msg.Actor.SendBad("They are too far away to reach in a single leap.")
			return
		}
		if s.actor.Mana.Current < config.LongLeapChiCost {
			s.msg.Actor.SendBad("You need " + strconv.Itoa(config.LongLeapChiCost) + " chi to leap that far.")
			return
		}
	}

	s.actor.Victim = whatMob
	s.actor.RunHook("combat")

	if _, engaged := whatMob.ThreatTable[s.actor.Name]; !engaged {
		s.msg.Actor.Send(text.White + "You engaged " + whatMob.Name + " #" + strconv.Itoa(s.where.Mobs.GetNumber(whatMob)) + " in combat.")
		whatMob.AddThreatDamage(0, s.actor)
	}

	if distance > config.LeapFreeDistance {
		s.actor.Mana.Subtract(config.LongLeapChiCost)
	}

	s.actor.Placement = whatMob.Placement
	s.msg.Actor.SendInfo("You leap at " + whatMob.Name + "!")
	s.msg.Observers.SendInfo(s.actor.Name + " leaps at " + whatMob.Name + "!")

	result := resolveHits(s, whatMob, []float64{1.0}, attackSkillLevel(s), "leap")
	if result.hits > 0 {
		gained := s.actor.GainChi(config.ChiPerHitFor(s.actor.GetStat("pie"))*result.chiHits(), true)
		if gained > 0 {
			s.msg.Actor.Send(text.Cyan + "Your chi rises by " + strconv.Itoa(gained) + "." + text.Reset)
		}
	}
	DeathCheck(s, whatMob)

	if config.LeapTimer > 0 {
		s.actor.SetTimer("combat_leap", config.LeapTimer)
	}
	s.actor.SetTimer("combat", config.CombatCooldown)
}
