package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(evaluate{},
		"Usage:  evaluate target\n\n  Examine a monster or item to find it's properties. ",
		permissions.Anyone,
		"evaluate", "eval")
}

type evaluate cmd

func (evaluate) process(s *state) {
	s.msg.Actor.SendInfo("Not implemented yet")
	return
}
	/*
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What do you want to evaluate?")
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

	// Check mobs
	var whatMob *objects.Mob
	if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		whatMob = s.where.Mobs.Search(name, nameNum, true)
	} else {
		whatMob = s.where.Mobs.Search(name, nameNum, false)
	}
	// It was a mob!
	if whatMob != nil {
		mob_template := "You study the {{.Name}} in your minds eye....\n" +
			"" +
			"It currently has {{.HP}} hits points remaining." +
			"It is worth {{.Exp}}experience points." +
			"{{if .Quick}} It is quick reacting \n{{end}}" +
			"{{if .Permanent}} It is permanent \n{{end}}" +
			"{{if .Hostile}} It is hostile \n{{end}}" +
			"{{if .Guards}} It guards treasure. \n{{end}}" +
			"{{if .PickUp}} It picks up treasure. \n{{end}}" +
			"{{if .Spells}} It casts spells \n{{end}}" +
		return
	}

	// Check items
	what := s.where.Items.Search(name, nameNum)

	// Item in the room?
	if what != nil {

		return
	}

	what = s.actor.Inventory.Search(name, nameNum)

	// It was on you the whole time
	if what != nil {

		return
	}

	what = s.actor.Equipment.Search(name)

	// Check your equipment
	if what != nil {

		return
	} else {
		s.msg.Actor.SendBad("You see no '", name, "' to examine.")
		return
	}

/*
TODO: There are more mob flags to be implemented


   You study the Cloud Giant Leader in your minds eye....

   It currently has 3,503 hits points remaining.
   It is worth 67,000 experience points.
   It is quick reacting.
   It breathes Lightning.
   It is permanent.
   It is hostile.
   It guards treasure.
   It block exits.
   It cannot be stolen from.
   It can cast spells.
   It can see invisibles.
   It can only be harmed by magical weapons.
   It resists magic.
   It cannot be stunned.
   It has a long ranged missile attack.
   It is carrying a Storm Cloud

   You study the Bullfrog in your minds eye....

   It currently has 1,000 hits points remaining.
   It is worth 2,000 experience points.
   It is quick reacting.
   It guards treasure.
   It flees.
   It picks up treasure.
   It causes disease.
   It has a long ranged missile attack.
   You are unable to discern what its carrying.

   You see an Orcish Battle-Savant, large, twenty-second level.
   A blood-soaked orcish warrior with his entire face tattooed wielding a
     massive claymore.
   He is standing a couple steps in front of you.
   He appears to be magical!
   He looks hostile!

   You study the Orcish Battle-Savant in your minds eye....

   It currently has 2,800 hits points remaining.
   It is worth 24,444 experience points.
   It is quick reacting.
   It is hostile.
   It block exits.
   It can only be harmed by magical weapons.
   It has a corrosive touch.
   It causes blindness.
   You are unable to discern what its carrying.


You study the Eye of the Lich in your minds eye....

It is a device.
It is charged with clairvoyance.
It has 15 charges remaining.
You determine its weight to be 1 lbs.
You judge its value to be 3,118 gold marks.
 */