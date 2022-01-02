package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(balance{},
	"Usage:  balance \n \n Displays the ",
	permissions.Player,
	"$BALANCE")
}

type balance cmd

func (balance) process(s *state) {
	s.msg.Actor.SendGood("Your bank account currently has " + strconv.Itoa(s.actor.BankGold.Value) + " gold marks.")
}
