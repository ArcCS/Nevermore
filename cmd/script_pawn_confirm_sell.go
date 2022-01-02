package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(sell_confirm{},
	"",
	permissions.Player,
	"$SELL_CONFIRM")
}

type sell_confirm cmd

func (sell_confirm) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Issue completing sell.")
		return
	}

	targetStr := s.words[0]
	targetNum := 1

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target != nil {
		/*finalValue := (.5*(target.Value)) +
		((45/s.actor.Int.Current) * (.25*target.Value)) +
		((utils.Roll(10, 1)/10) * (.25*target.Value))*/
		finalValue := int(.65*float64(target.Value))
		s.actor.Inventory.Remove(target)
		s.actor.Gold.Add(finalValue)
		s.msg.Actor.SendGood("The pawn broker gives you ", strconv.Itoa(finalValue), " for ", target.Name, ".")

	}else{
		s.msg.Actor.SendInfo("What're you trying to sell??")
		return
	}

}
