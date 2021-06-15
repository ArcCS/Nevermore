package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(bonus{},
		"Usage:  bonus ## player Bonus a player or the whole room, leave player name off to bonus all present",
		permissions.Dungeonmaster,
		"bonus")
}

type bonus cmd

func (bonus) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Bonus how much?")
		return
	}

	if amt, err := strconv.Atoi(s.words[0]); err == nil {
		if len(s.words) >= 2 {
			targetChar := s.where.Chars.Search(s.words[1], s.actor)
			if targetChar != nil {
				targetChar.BonusPoints.Add(amt)
				s.msg.Actor.SendGood("You've been awarded " + s.words[0] + " bonus points!")
			}
		} else {
			for _, actor := range s.where.Chars.Contents {
				actor.BonusPoints.Add(amt)
			}
			s.msg.Observers.SendGood("You've been awarded " + s.words[0] + " bonus points!")
			return
		}
	}else{
		s.msg.Actor.SendBad("Not an appropriate value to bonus with.")
		return
	}

	s.ok = true
	return
}
