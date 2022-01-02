package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(withdraw{},
	"",
	permissions.Player,
	"$WITHDRAW")
}

type withdraw cmd

func (withdraw) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("What value would you like to withdraw?")
		return
	}

	value := s.words[0]

	if amount64, err := strconv.Atoi(value); err == nil {
		amount := amount64
		if s.actor.BankGold.CanSubtract(amount) {
			s.actor.RunHook("act")
			s.actor.BankGold.SubIfCan(amount)
			s.actor.Gold.Add(amount)
			s.msg.Actor.SendGood("You withdraw ", value, " from your bank account. \n ====")
			s.msg.Actor.SendGood("Your bank account currently has " + strconv.Itoa(s.actor.BankGold.Value) + " gold marks.")
		} else {
			s.msg.Actor.SendInfo("You don't have that much gold in your bank account.")
			return
		}
	}else{
		s.msg.Actor.SendBad("That's not a valid number.")
		return
	}
}
