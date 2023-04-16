package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(equip{},
		"Usage:  equip item # \n\n Try to equip an item from your inventory",
		permissions.Player,
		"equip", "wield", "wear")
}

type equip cmd

func (equip) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendBad("What did you want to equip?")
		return
	}

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	what := s.actor.Inventory.Search(name, nameNum)
	if what != nil {
		s.actor.RunHook("combat")
		if s.actor.Class == 8 && utils.IntIn(what.ItemType, []int{0, 1, 2, 3, 4}) {
			s.msg.Actor.SendBad("You cannot wield weapons effectively.")
			return
		}
		if utils.IntIn(what.ItemType, []int{5, 19, 20, 21, 22, 23, 24, 25, 26}) {
			if !config.CheckArmor(what.ItemType, s.actor.Tier, what.Armor) {
				s.msg.Actor.SendBad("You are unsure of how to maximize the benefit of this armor and cannot wear it.")
				return
			}
			if !config.CheckMaxArmor("max", s.actor.Tier, what.Armor+s.actor.Equipment.Armor) {
				s.msg.Actor.SendBad("You cannot equip this item in addition to your current equipment, max armor value exceeded.")
				return
			}
			if what.ItemType == 23 {
				if !config.CheckMaxArmor("shield_block", s.actor.Tier, what.Armor+s.actor.Equipment.Armor) {
					s.msg.Actor.SendBad("You cannot equip this item in addition to your current equipment, max shield block exceeded.")
					return
				}
			}
		}
		if utils.IntIn(what.ItemType, []int{0, 1, 2, 3, 4}) &&
			!s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			if !config.CanWield(s.actor.Tier, s.actor.Class, utils.RollMax(what.SidesDice, what.NumDice, what.PlusDice)) {
				s.msg.Actor.SendBad("You are not well enough trained to wield " + what.Name)
				return
			}
		}
		s.actor.Inventory.Lock()
		if s.actor.Equipment.Equip(what) {
			s.msg.Actor.SendGood("You equip " + what.DisplayName())
			s.msg.Observers.SendInfo(s.actor.Name + " equips " + what.DisplayName())
			s.actor.Inventory.Remove(what)
		} else {
			s.msg.Actor.SendBad("You cannot equip that.")
		}

		s.actor.Inventory.Unlock()

		s.ok = true
		return
	}
	s.msg.Actor.SendInfo("What did you want to equip?")
	s.ok = true
}
