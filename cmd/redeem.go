package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"math"
	"strconv"
)

func init() {
	addHandler(redeem{},
		"Usage:  redeem gold|exp ## Redeem your bonus points for 1% gold or experience of your current required tier",
		permissions.Player,
		"redeem")
}

type redeem cmd

func (redeem) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Redeem how much?")
		return
	}

	redType := s.words[0]
	amt := 0
	var err error

	if redType != "GOLD" && redType != "EXP" {
		s.msg.Actor.SendBad("You can only redeem for 'gold' or 'exp'.")
		return
	}

	if amt, err = strconv.Atoi(s.words[1]); err != nil {
		s.msg.Actor.SendBad("Not a valid value to redeem.")
		return
	}

	if amt > s.actor.BonusPoints.Value {
		s.msg.Actor.SendBad("You do not have that many bonus points to redeem.")
		return
	}

	if redType == "GOLD" {
		s.actor.BonusPoints.Subtract(amt)
		totalGold := int(math.Floor(float64(config.GoldPerLevel[s.actor.Tier+1])*.01)) * amt
		s.actor.Gold.Add(totalGold)
		s.msg.Actor.SendGood("You have redeemed ", strconv.Itoa(amt), " bonus points for ", strconv.Itoa(totalGold), " gold.")
		return
	}

	if redType == "EXP" {
		s.actor.BonusPoints.Subtract(amt)
		experienceNeeded := config.TierExpLevels[s.actor.Tier+1] - config.TierExpLevels[s.actor.Tier]
		expAward := int(math.Floor(float64(experienceNeeded)*.01)) * amt
		s.actor.Experience.Add(expAward)
		s.msg.Actor.SendGood("You have redeemed ", strconv.Itoa(amt), " bonus points for ", strconv.Itoa(expAward), " experience.")
		return
	}

	s.ok = true
	return
}
