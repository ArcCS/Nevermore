package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(scriptBind{},
		"",
		permissions.Player,
		"$SOULBIND")
	addHandler(confirmBind{},
		"",
		permissions.Player,
		"$CONFIRMBIND")
}

type scriptBind cmd

func (scriptBind) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Bind what????")
		return
	}

	targetStr := s.words[0]
	targetNum := 1

	if len(s.words) == 2 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		}
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target != nil {
		if target.Flags["permanent"] {
			s.msg.Actor.SendBad("That item is already bound to your soul.")
			return
		}

		s.msg.Actor.SendGood("The price to bind " + target.Name + " to your soul is " + strconv.Itoa(config.BindCost) + " gold marks.")
		s.msg.Actor.SendGood("Accept offer? (y, yes to confirm)")
		s.actor.AddCommands("yes", "$CONFIRMBIND "+targetStr+" "+strconv.Itoa(targetNum))
		s.actor.AddCommands("y", "$CONFIRMBIND "+targetStr+" "+strconv.Itoa(targetNum))
	} else {
		s.msg.Actor.SendInfo("What're you trying to bind")
		return
	}
}

type confirmBind cmd

func (confirmBind) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Bind error")
		return
	}

	targetStr := s.words[0]
	targetNum := 1

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target != nil {
		if !s.actor.Gold.CanSubtract(config.BindCost) {
			s.msg.Actor.SendBad("You do not have enough gold to bind this item to you.")
			return
		}
		s.actor.Gold.Subtract(config.BindCost)
		target.Flags["permanent"] = true
		s.msg.Actor.SendGood(target.Name + " is now bound to you and cannot be dropped or lost through death.")

	} else {
		s.msg.Actor.SendInfo("What're you trying to bind??")
		return
	}
}
