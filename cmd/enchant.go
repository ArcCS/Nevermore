package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

// Syntax: DROP item
func init() {
	addHandler(enchant{},
		"Usage:  enchant itemName # \n \n Allows mages to imbue damaging magics into weapons, and paladins to imbue protective magic into armor",
		permissions.Mage|permissions.Paladin,
		"enchant")
}

type enchant cmd

func (enchant) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to enchant")
		return
	}

	// We have at least 2 items here so lets move forward with that
	targetStr := s.words[0]
	targetNum := 1

	if len(s.words) < 1 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		} else {
			s.msg.Actor.SendInfo("What did you want to enchant?")
			return
		}
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target == nil {
		s.msg.Actor.SendInfo("What're you trying to enchant??")
		return
	}

	if target.Flags["magic"] {
		s.msg.Actor.SendBad("That item is already enchanted!")
		return
	}

	if s.actor.Class == 4 && !utils.IntIn(target.ItemType, config.WeaponTypes) {
		s.msg.Actor.SendBad("You can only enchant weapons!")
		return
	}

	if s.actor.Class == 6 && !utils.IntIn(target.ItemType, config.ArmorTypes) {
		s.msg.Actor.SendBad("You can only enchant armor!")
		return
	}

	s.actor.RunHook("act")
	s.msg.Actor.SendGood("You chant: \"I inject my magicks into thee!\"")
	s.msg.Observers.SendGood(s.actor.Name + " chants: \"I inject my magicks into thee!\"")
	s.actor.ClassProps["enchants"]--
	target.Flags["magic"] = true
	if utils.IntIn(target.ItemType, config.ArmorTypes) {
		if target.ItemType == 24 {
			target.Flags["light"] = true
			s.msg.Actor.SendGood("The ", target.Name, " glows brightly.")
		} else {
			target.Armor += DetermineArmorEnchant(s.actor.Tier)
		}
	}
	if utils.IntIn(target.ItemType, config.WeaponTypes) {
		target.Adjustment += DetermineWeaponEnchant(s.actor.Tier)
	}
	s.msg.Actor.SendGood("You enchanted ", target.Name, ".")
	s.ok = true
}

func DetermineWeaponEnchant(tier int) int {
	switch {
	case tier < 5:
		return 2
	case tier < 10:
		return 3
	case tier < 15:
		return 4
	case tier < 20:
		return 5
	case tier < 25:
		return 6
	default:
		return 2
	}
}

func DetermineArmorEnchant(tier int) int {
	switch {
	case tier < 5:
		return 1
	case tier < 10:
		return 2
	case tier < 15:
		return 3
	case tier < 20:
		return 4
	case tier < 25:
		return 5
	default:
		return 1
	}
}
