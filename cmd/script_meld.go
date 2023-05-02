package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(scriptMeld{},
		"",
		permissions.Player,
		"$MELD")
	addHandler(confirmMeld{},
		"",
		permissions.Player,
		"$CONFIRMMELD")
}

type scriptMeld cmd

func (scriptMeld) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendBad("Meld what into what?")
		return
	}

	// We have at least 2 items here so lets move forward with that
	argParse := 1
	targetStr := s.words[0]
	targetNum := 1

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
		argParse = 2
	}

	if argParse == 2 && len(s.words) <= 2 {
		s.msg.Actor.SendInfo("Meld it into what?")
		return
	}

	meldStr := s.words[argParse]
	meldNum := 1

	if len(s.words) >= argParse+2 {
		if val, err := strconv.Atoi(s.words[argParse+1]); err == nil {
			meldNum = val
		}
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)
	meld := s.actor.Inventory.Search(meldStr, meldNum)

	if target == nil || meld == nil {
		s.msg.Actor.SendBad("You have no " + targetStr + " to meld.")
		return
	} else {
		if utils.IntIn(meld.ItemType, []int{6, 15}) {
			if meld.ItemType != target.ItemType {
				s.msg.Actor.SendBad("These are not the same item type.")
				return
			}
			if meld.Spell != target.Spell {
				s.msg.Actor.SendBad("These items do not contain the same spell.")
				return
			}
			cost := int(float64(target.Value+target.Value/2) + ((float64(meld.MaxUses) / float64(objects.Items[meld.ItemId].MaxUses)) * float64(objects.Items[meld.ItemId].Value)))
			s.msg.Actor.SendInfo("The cost to meld this item will be " + strconv.Itoa(cost) + ".  Do you want to meld it? (Type yes to meld)")
			s.actor.AddCommands("yes", "$CONFIRMMELD "+targetStr+" "+strconv.Itoa(targetNum)+" "+meldStr+" "+strconv.Itoa(meldNum))
			s.actor.AddCommands("y", "$CONFIRMMELD "+targetStr+" "+strconv.Itoa(targetNum)+" "+meldStr+" "+strconv.Itoa(meldNum))
		} else {
			s.msg.Actor.SendBad("These are not meldable items")
			return
		}
	}
}

type confirmMeld cmd

func (confirmMeld) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendBad("Meld error")
		return
	}
	// We have at least 2 items here so lets move forward with that
	argParse := 1
	targetStr := s.words[0]
	targetNum := 1

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
		argParse = 2
	}

	if argParse == 2 && len(s.words) <= 2 {
		s.msg.Actor.SendInfo("Meld it into what?")
		return
	}

	meldStr := s.words[argParse]
	meldNum := 1

	if len(s.words) >= argParse+2 {
		if val, err := strconv.Atoi(s.words[argParse+1]); err == nil {
			meldNum = val
		}
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)
	meld := s.actor.Inventory.Search(meldStr, meldNum)

	if target == nil || meld == nil {
		s.msg.Actor.SendBad("Meld error")
		return
	} else {
		if utils.IntIn(target.ItemType, config.ArmorTypes) || utils.IntIn(meld.ItemType, config.WeaponTypes) || utils.IntIn(meld.ItemType, []int{6, 15}) {
			cost := int(float64(target.Value+target.Value/2) + ((float64(meld.MaxUses) / float64(objects.Items[meld.ItemId].MaxUses)) * float64(objects.Items[meld.ItemId].Value)))
			if s.actor.Gold.Value < cost {
				s.msg.Actor.SendBad("You do not have enough money to meld this item.")
				return
			} else {
				s.actor.Gold.Subtract(cost)
			}
			target.MaxUses += meld.MaxUses
			s.actor.Inventory.Remove(meld)
			s.msg.Actor.SendGood("Meld completed.")
			meld = nil
		} else {
			s.msg.Actor.SendBad("These are not meldable items")
			return
		}
	}
}
