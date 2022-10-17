package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"math"
	"strconv"
)

func init() {
	addHandler(scriptMeld{},
		"",
		permissions.Player,
		"$REPAIR")
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
		s.msg.Actor.SendInfo("Put it where?")
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
	} else if meld == nil {
		s.msg.Actor.SendBad("You have no " + meldStr + " to meld.")
		return
	} else {
		if utils.IntIn(target.ItemType, config.ArmorTypes) || utils.IntIn(meld.ItemType, config.WeaponTypes) || utils.IntIn(meld.ItemType, []int{6, 15}) {
			// Calculate a cost of repair
			// TODO: Replace this with a formula based on int, later with crafting/blacksmithing
			base_cost := target.Value + target.Value/2
			secondary_cost := (meld.MaxUses / objects.Items[meld.ParentItemId].MaxUses) * objects.Items[meld.ParentItemId].Value
			s.msg.Actor.SendInfo("The cost to repair this item will be " + strconv.Itoa(cost) + ".  Do you want to repair it? (Type yes to repair)")
			s.actor.AddCommands("yes", "$CONFIRMMELD "+targetStr+" "+strconv.Itoa(targetNum)+" "+meldStr+" "+strconv.Itoa(meldNum))
		} else {
			s.msg.Actor.SendBad("This is not a repairable item.")
			return
		}
	} else {
		s.msg.Actor.SendBad("You cannot meld these items together")
		return
	}
}

type confirmMeld cmd

func (confirmMeld) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Meld error")
		return
	}

	targetStr := s.words[0]
	targetNum := 1

	if len(s.words) > 1 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		}
	}

	what := s.actor.Inventory.Search(targetStr, targetNum)

	if what != nil {
		if utils.IntIn(what.ItemType, config.ArmorTypes) || utils.IntIn(what.ItemType, config.WeaponTypes) {
			// TODO: Replace this with a formula based on int, later with crafting/blacksmithing
			cost := int(math.Round(3 * (float64(what.Value) * float64(objects.Items[what.ItemId].MaxUses/(objects.Items[what.ItemId].MaxUses-what.MaxUses)))))
			if s.actor.Gold.CanSubtract(cost) {
				s.actor.Gold.Subtract(cost)
				what.MaxUses = objects.Items[what.ItemId].MaxUses
				s.msg.Actor.SendInfo("Your item was repaired.")
			} else {
				s.msg.Actor.SendBad("You don't have enough gold to repair this item.")
				return
			}

		} else {
			s.msg.Actor.SendBad("This is not an item that can be repaired")
			return
		}
	} else {
		s.msg.Actor.SendBad("You don't have anything like that in your inventory to repair.")
		return
	}
}
