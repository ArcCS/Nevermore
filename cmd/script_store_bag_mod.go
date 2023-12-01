package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(modbag{},
		"",
		permissions.Player,
		"$MODBAG")
	addHandler(modconfirm{},
		"",
		permissions.Player,
		"$MODCONFIRM")
}

type modbag cmd

type modconfirm cmd

func (modbag) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendBad("What bag do you want to modify? And how (weight, capacity, weightless)")
		return
	}

	targetStr := s.words[0]

	target := s.actor.Inventory.Search(targetStr, 1)

	actionStr := s.words[1]

	if target != nil {
		if target.ItemType != 9 {
			s.msg.Actor.SendBad("That's not a bag.")
			return
		}

		switch actionStr {
		case "WEIGHT":
			if newWeight, err := strconv.Atoi(s.words[2]); err == nil {
				if newWeight > target.Weight {
					s.msg.Actor.SendBad("You can't make the bag heavier.")
					return
				}
				newCost := bagWeight[newWeight] - bagWeight[target.Weight]
				s.msg.Actor.SendGood("It will cost you ", strconv.Itoa(newCost), " gold to make the bag ", strconv.Itoa(newWeight), " pounds.")
				s.msg.Actor.SendGood("Accept offer? (y, yes to confirm)")
				s.actor.AddCommands("yes", "$MODCONFIRM "+targetStr+" "+actionStr+" "+strconv.Itoa(newWeight))
				s.actor.AddCommands("y", "$MODCONFIRM "+targetStr+" "+actionStr+" "+strconv.Itoa(newWeight))

			} else {
				s.msg.Actor.SendBad("What weight do you want to set?")
				return
			}
		case "WEIGHTLESS":
			if target.Flags["weightless_chest"] {
				s.msg.Actor.SendBad("The bag is already weightless.")
				return
			}

			for _, item := range s.actor.Inventory.ListItems() {
				if item.Flags["weightless_chest"] {
					s.msg.Actor.SendBad("You already have a weightless bag.")
					return
				}
			}

			s.msg.Actor.SendGood("It will cost you " + strconv.Itoa(weightLess) + " gold to make the bag weightless.")
			s.msg.Actor.SendGood("Accept offer? (y, yes to confirm)")
			s.actor.AddCommands("yes", "$MODCONFIRM "+targetStr+" "+actionStr)
			s.actor.AddCommands("y", "$MODCONFIRM "+targetStr+" "+actionStr)
		case "CAPACITY":
			if newCapacity, err := strconv.Atoi(s.words[2]); err == nil {
				if newCapacity < target.MaxUses && newCapacity <= 30 {
					s.msg.Actor.SendBad("You can't modify the bag to hold less than it already does.")
					return
				}
				newCost := newCapacity * bagCapacity
				newCost -= target.MaxUses * bagCapacity
				s.msg.Actor.SendGood("It will cost you " + strconv.Itoa(newCost) + " gold to make the bag hold " + strconv.Itoa(newCapacity) + " items.")
				s.msg.Actor.SendGood("Accept offer? (y, yes to confirm)")
				s.actor.AddCommands("yes", "$MODCONFIRM "+targetStr+" "+actionStr+" "+strconv.Itoa(newCapacity))
				s.actor.AddCommands("y", "$MODCONFIRM "+targetStr+" "+actionStr+" "+strconv.Itoa(newCapacity))

			} else {
				s.msg.Actor.SendBad("Capacity could not be parsed")
				return
			}

		default:
			s.msg.Actor.SendBad("What do you want to modify about the bag?")
		}
	}
}

func (modconfirm) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendBad("Bag modification error")
		return
	}

	targetStr := s.words[0]

	target := s.actor.Inventory.Search(targetStr, 1)

	actionStr := s.words[1]

	if target != nil {
		if target.ItemType != 9 {
			s.msg.Actor.SendBad("That's not a bag.")
			return
		}

		switch actionStr {
		case "WEIGHT":
			if newWeight, err := strconv.Atoi(s.words[2]); err == nil {
				if newWeight > target.Weight {
					s.msg.Actor.SendBad("You can't make the bag heavier.")
					return
				}
				newCost := bagWeight[newWeight] - bagWeight[target.Weight]

				if s.actor.Gold.Value < newCost {
					s.msg.Actor.SendBad("You don't have enough gold to make the bag that heavy.")
					return
				} else {
					s.actor.Gold.Value -= newCost
					target.Weight = newWeight
					target.Save()

					s.msg.Actor.SendInfo("You pay ", strconv.Itoa(newCost), " gold to make the bag ", strconv.Itoa(newWeight), " pounds.")
				}

			} else {
				s.msg.Actor.SendBad("What weight do you want to set?")
				return
			}
		case "WEIGHTLESS":
			if target.Flags["weightless_chest"] {
				s.msg.Actor.SendBad("The bag is already weightless.")
				return
			}

			if s.actor.Gold.Value < weightLess {
				s.msg.Actor.SendBad("You don't have enough gold to make the bag weightless.")
				return
			} else {
				s.actor.Gold.Value -= weightLess
				target.Flags["weightless_chest"] = true
				target.Save()

				s.msg.Actor.SendInfo("You pay ", strconv.Itoa(weightLess), " gold to make the bag weightless.")
			}
		case "CAPACITY":
			if newCapacity, err := strconv.Atoi(s.words[2]); err == nil {
				if newCapacity < target.MaxUses {
					s.msg.Actor.SendBad("You can't modify the bag to hold less than it already does.")
					return
				}
				newCost := newCapacity * bagCapacity
				newCost -= target.MaxUses * bagCapacity

				if s.actor.Gold.Value < newCost {
					s.msg.Actor.SendBad("You don't have enough gold to make the bag that heavy.")
					return
				} else {
					s.actor.Gold.Value -= newCost
					target.MaxUses = newCapacity
					target.Save()

					s.msg.Actor.SendInfo("You pay ", strconv.Itoa(newCost), " gold to make the bag hold ", strconv.Itoa(newCapacity), " items.")
				}

			} else {
				s.msg.Actor.SendBad("Capacity could not be parsed")
				return
			}
		default:
			s.msg.Actor.SendBad("Error processing, inform a GM.")
		}

	}
}
