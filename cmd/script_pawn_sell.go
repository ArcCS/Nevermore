package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
)

func init() {
	addHandler(sell{},
		"",
		permissions.Player,
		"$SELL")
	addHandler(sell_confirm{},
		"",
		permissions.Player,
		"$SELL_CONFIRM")
}

type sell cmd

func (sell) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Sell what????")
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
		finalValue := 0
		// real dumb
		if s.actor.GetStat("int") < 5 {
			finalValue = int(.10 * float64(target.Value))
		} else if s.actor.GetStat("int") < 10 {
			finalValue = int((.25 * float64(target.Value)) +
				(float64(utils.Roll(10, 1, 0))/float64(10))*(.25*float64(target.Value)))
		} else if s.actor.GetStat("int") >= 10 {
			finalValue = int((.5 * float64(target.Value)) +
				((float64(s.actor.Int.Current) / 45) * (.25 * float64(target.Value))) +
				(float64(utils.Roll(10, 1, 0))/float64(10))*(.25*float64(target.Value)))
		}

		s.msg.Actor.SendGood("The pawn broker offers you ", strconv.Itoa(finalValue), " for ", target.Name, ".")
		s.msg.Actor.SendGood("Accept offer? (y, yes to confirm)")
		s.actor.AddCommands("yes", "$SELL_CONFIRM "+targetStr+" "+strconv.Itoa(targetNum)+" "+strconv.Itoa(finalValue))
		s.actor.AddCommands("y", "$SELL_CONFIRM "+targetStr+" "+strconv.Itoa(targetNum)+" "+strconv.Itoa(finalValue))
	} else {
		s.msg.Actor.SendInfo("What're you trying to sell??")
		return
	}

}

type sell_confirm cmd

func (sell_confirm) process(s *state) {
	if len(s.words) < 3 {
		s.msg.Actor.SendInfo("Issue completing sell.")
		return
	}

	targetStr := s.words[0]
	targetNum := 1
	targetPrice := 0

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if val, err := strconv.Atoi(s.words[2]); err == nil {
		targetPrice = val
	} else {
		s.msg.Actor.SendInfo("Issue completing sell.")
		log.Println("Error converting target price to int: ", err)
		return
	}

	if target != nil {
		if ok := s.actor.Inventory.Remove(target); ok == nil {
			s.actor.Gold.Add(targetPrice)
			s.msg.Actor.SendGood("The pawn broker gives you ", strconv.Itoa(targetPrice), " for ", target.Name, ".")
		} else {
			s.msg.Actor.SendBad("Issue completing sell.")
			log.Println("Error removing item from inventory: ", ok)
		}

	} else {
		s.msg.Actor.SendInfo("What're you trying to sell??")
		return
	}

}
