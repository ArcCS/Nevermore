package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(deposit{},
	"",
	permissions.Player,
	"$DEPOSIT")
}

type deposit cmd

func (deposit) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("What value would you like to deposit?")
		return
	}

	value := s.words[0]

	if amount64, err := strconv.Atoi(value); err == nil {
		amount := amount64
		if s.actor.Gold.CanSubtract(amount) {
			s.actor.RunHook("act")
			s.actor.Gold.SubIfCan(amount)
			s.actor.BankGold.Add(amount)
			s.msg.Actor.SendGood("You deposit ", value, " into you bank account. \n ====")
			s.msg.Actor.SendGood("Your bank account currently has " + strconv.Itoa(s.actor.BankGold.Value) + " gold marks.")
		} else {
			s.msg.Actor.SendInfo("You don't have that much gold.")
			return
		}
	}else{
		s.msg.Actor.SendBad("That's not a valid number.")
		return
	}
}
