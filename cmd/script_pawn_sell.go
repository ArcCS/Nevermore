package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(sell{},
	"",
	permissions.Player,
	"$SELL")
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
		/*finalValue := (.5*(target.Value)) +
			((45/s.actor.Int.Current) * (.25*target.Value)) +
			((utils.Roll(10, 1)/10) * (.25*target.Value))*/
		finalValue := int(.65*float64(target.Value))

		s.msg.Actor.SendGood("The pawn broker offers you ", strconv.Itoa(finalValue), " for ", target.Name, ".")
		s.msg.Actor.SendGood("Accept offer? (y, yes to confirm)")
		s.actor.AddCommands("yes", "$SELL_CONFIRM " + targetStr + " " + strconv.Itoa(targetNum))
		s.actor.AddCommands("y", "$SELL_CONFIRM " + targetStr + " " + strconv.Itoa(targetNum))
	}else{
		s.msg.Actor.SendInfo("What're you trying to sell??")
		return
	}

}
