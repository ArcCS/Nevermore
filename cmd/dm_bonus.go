package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
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

			// Search Active Characters
			if targetChar == nil {
				targetChar = objects.ActiveCharacters.Find(s.words[1])
			}

			if targetChar != nil {
				if targetChar.BonusPoints.Value <= 50 {
					targetChar.BonusPoints.Add(amt)
					s.participant = targetChar
					s.msg.Participant.SendGood("You've been awarded " + s.words[0] + " bonus points!")
					s.msg.Actor.SendGood("You've awarded " + s.words[0] + " bonus points to " + targetChar.Name + "!")
				} else {
					s.msg.Actor.SendBad("That player has too many bonus points already!")
					return
				}

			}
		} else {
			for _, actor := range s.where.Chars.Contents {
				if actor.BonusPoints.Value <= 50 {
					actor.BonusPoints.Add(amt)
				} else {
					s.msg.Actor.SendBad(actor.Name + " has too many bonus points already!")
					return
				}
			}
			s.msg.Observers.SendGood("You've been awarded " + s.words[0] + " bonus points!")
			s.msg.Actor.SendGood("You've awarded " + s.words[0] + " bonus points to everyone in the room!")
			return
		}
	} else {
		s.msg.Actor.SendBad("Not an appropriate value to bonus with.")
		return
	}

	s.ok = true
	return
}
