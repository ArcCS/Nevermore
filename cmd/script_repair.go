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
	addHandler(scriptRepair{},
		"",
		permissions.Player,
		"$REPAIR")
	addHandler(confirmRepair{},
		"",
		permissions.Player,
		"$CONFIRMREPAIR")
}

type scriptRepair cmd

func (scriptRepair) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Repair what?")
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
			// Calculate a cost of repair
			// TODO: Replace this with a formula based on int, later with crafting/blacksmithing

			cost := int(math.Round(float64(3) * (float64(what.Value)) * (float64(objects.Items[what.ItemId].MaxUses-what.MaxUses) / float64(objects.Items[what.ItemId].MaxUses))))
			s.msg.Actor.SendInfo("The cost to repair this item will be " + strconv.Itoa(cost) + ".  Do you want to repair it? (Type yes to repair)")
			s.actor.AddCommands("yes", "$CONFIRMREPAIR "+targetStr+" "+strconv.Itoa(targetNum))
		} else {
			s.msg.Actor.SendBad("This is not a repairable item.")
			return
		}
	} else {
		s.msg.Actor.SendBad("You don't have anything like that in your inventory to repair.")
		return
	}
}

type confirmRepair cmd

func (confirmRepair) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Repair what?")
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
			cost := int(math.Round(float64(3) * (float64(what.Value)) * (float64(objects.Items[what.ItemId].MaxUses-what.MaxUses) / float64(objects.Items[what.ItemId].MaxUses))))
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
