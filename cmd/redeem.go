package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(redeem{},
		"Usage:  redeem gold|exp ## Redeem your bonus points for 5k gold or 5% of your curent required experiene",
		permissions.Dungeonmaster,
		"redeem")
}

type redeem cmd

func (redeem) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Bonus how much?")
		return
	}

	if amt, err := strconv.Atoi(s.words[0]); err == nil {
		if len(s.words) >= 2 {
			targetChar := s.where.Chars.Search(s.words[1], s.actor)
			if targetChar != nil {
				targetChar.BonusPoints.Add(amt)
				s.participant = targetChar
				s.msg.Participant.SendGood("You've been awarded " + s.words[0] + " bonus points!")
			}
		} else {
			for _, actor := range s.where.Chars.Contents {
				actor.BonusPoints.Add(amt)
			}
			s.msg.Observers.SendGood("You've been awarded " + s.words[0] + " bonus points!")
			return
		}
	} else {
		s.msg.Actor.SendBad("Not an appropriate value to bonus with.")
		return
	}

	s.ok = true
	return
}
