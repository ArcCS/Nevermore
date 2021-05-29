package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"log"
	"math"
	"strconv"
	"strings"
	"text/template"
)

// Overloaded Look object for all of your looking pleasure
// Syntax: ( LOOK | L ) has.Thing
func init() {
	addHandler(look{},
		"Usage:  look [object|exit|character|mob] # \n \n Put your peepers on something. (Also can use short hand L",
		permissions.Player,
		"LOOK", "L")
}

type look cmd

func (look) process(s *state) {
	var others []string
	var mobs string
	var mobAttacking string
	if len(s.input) == 0 {
		roomLook := objects.Rooms[s.actor.ParentId]
		s.msg.Actor.SendInfo(roomLook.Look(s.actor))
		others = objects.Rooms[s.actor.ParentId].Chars.List(s.actor)
		mobs = objects.Rooms[s.actor.ParentId].Mobs.ReducedList(s.actor)
		mobAttacking = objects.Rooms[s.actor.ParentId].Mobs.ListAttackers(s.actor)
		if len(others) == 1 {
			s.msg.Actor.SendInfo(strings.Join(others, ", "), " is also here.")
		} else if len(others) > 1 {
			s.msg.Actor.SendInfo(strings.Join(others, ", "), " are also here.")
		}
		if len(mobs) > 0 {
			s.msg.Actor.SendInfo("You see " + mobs)
		}
		items := roomLook.Items.ReducedList()
		if len(items) > 0 {
			s.msg.Actor.SendInfo("On the ground you see " + items)
		}
		if len(mobAttacking) > 0 {
			s.msg.Actor.SendInfo(mobAttacking)
		}
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

	var whatChar *objects.Character
	// Check characters in the room first.
	whatChar = s.where.Chars.Search(name, s.actor)
	// It was a person!
	if whatChar != nil {
		s.msg.Actor.SendInfo(whatChar.Look())
		s.msg.Actor.SendInfo(text.Gray + whatChar.Description + "\n")
		equip_template := "{{if .Chest}}\n{{.Sub_pronoun}} {{.Isare}} wearing {{.Chest}} about {{.Pos_pronoun}} body{{end}}" +
			" {{if .Neck}}\n{{.Sub_pronoun}} {{.Isare}} a {{.Neck}} around {{.Pos_pronoun}} neck.{{end}}" +
			" {{if .Main}}\n{{.Sub_pronoun}} {{.Isare}} holding a {{.Main}} in {{.Pos_pronoun}} main hand.{{end}}" +
			" {{if .Offhand}}\n{{.Sub_pronoun}} {{.Isare}} holding a {{.Offhand}} in {{.Pos_pronoun}} off hand.{{end}}" +
			" {{if .Arms}}\n{{.Sub_pronoun}} {{.Isare}} wearing some {{.Arms}} on {{.Pos_pronoun}} arms.{{end}}" +
			" {{if .Finger1}}\n{{.Sub_pronoun}} {{.HasHave}} a {{.Finger1}} on {{.Pos_pronoun}} finger.{{end}}" +
			" {{if .Finger2}}\n{{.Sub_pronoun}} {{.HasHave}} a {{.Finger2}} on {{.Pos_pronoun}} finger.{{end}}" +
			" {{if .Legs}}\n{{.Sub_pronoun}} {{.HasHave}} {{.Legs}} on {{.Pos_pronoun}} legs.{{end}}" +
			" {{if .Feet}}\n{{.Sub_pronoun}} {{.HasHave}} {{.Feet}} on {{.Pos_pronoun}} feet.{{end}}" +
			" {{if .Head}}\n{{.Sub_pronoun}} {{.Isare}} wearing a {{.Head}}.{{end}}" +
			text.Reset

		data := struct {
			Sub_pronoun string
			Pos_pronoun string
			Isare       string
			HasHave     string
			Chest       string
			Neck        string
			Main        string
			Offhand     string
			Arms        string
			Finger1     string
			Finger2     string
			Legs        string
			Feet        string
			Head        string
		}{
			strings.Title(config.TextSubPronoun[whatChar.Gender]),
			config.TextPosPronoun[whatChar.Gender],
			"is",
			"has",
			whatChar.Equipment.GetText("chest"),
			whatChar.Equipment.GetText("neck"),
			whatChar.Equipment.GetText("main"),
			whatChar.Equipment.GetText("off"),
			whatChar.Equipment.GetText("arms"),
			whatChar.Equipment.GetText("ring1"),
			whatChar.Equipment.GetText("ring2"),
			whatChar.Equipment.GetText("legs"),
			whatChar.Equipment.GetText("feet"),
			whatChar.Equipment.GetText("head"),
		}

		tmpl, _ := template.New("char_info").Parse(equip_template)

		var output bytes.Buffer
		err := tmpl.Execute(&output, data)
		if err != nil {
			log.Println(err)
		} else {
			s.msg.Actor.SendGood(output.String())
		}

		s.ok = true

		return
	}

	// Check exits
	whatExit := s.where.FindExit(strings.ToLower(name))

	// Nice, looking at an exit.
	if whatExit != nil {
		if whatExit.Description == "" {
			s.msg.Actor.SendInfo(objects.Rooms[whatExit.ToId].Look(s.actor))
		} else {
			s.msg.Actor.SendInfo(whatExit.Look())
		}
		return
	}

	// Check mobs
	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	// It was a mob!
	if whatMob != nil {
		s.msg.Actor.SendInfo(whatMob.Look())
		s.msg.Actor.SendInfo("It is standing" + WhereAt(whatMob.Placement, s.actor.Placement))
		//log.Println(whatMob.ThreatTable)
		_, ok := whatMob.ThreatTable[s.actor.Name]
		if !ok {
			s.msg.Actor.SendInfo("It isn't paying attention to you.")
		} else {
			s.msg.Actor.SendInfo("It appears very angry at you!")
		}
		if whatMob.CurrentTarget == s.actor.Name {
			s.msg.Actor.SendInfo("It is attacking you!")
		} else if whatMob.CurrentTarget != "" {
			s.msg.Actor.SendInfo("it is attacking " + whatMob.CurrentTarget)
		}
		return
	}

	// Check items
	what := s.where.Items.Search(name, nameNum)

	// Item in the room?
	if what != nil {
		s.msg.Actor.SendInfo(what.Look())
		s.msg.Actor.SendInfo("It is on the ground " + WhereAt(what.Placement, s.actor.Placement))
		return
	}

	what = s.actor.Inventory.Search(name, nameNum)

	// It was on you the whole time
	if what != nil {
		s.msg.Actor.SendInfo("You examine " + what.Name)
		s.msg.Actor.SendInfo(what.Look())
		return
	}

	what = s.actor.Equipment.Search(name)

	// Check your equipment
	if what != nil {
		s.msg.Actor.SendInfo("You examine your equipped " + what.Name)
		s.msg.Actor.SendInfo(what.Look())
		return
	} else {
		s.msg.Actor.SendBad("You see no '", name, "' to examine.")
		return
	}
}

func WhereAt(subLoc int, charLoc int) string {
	// Moving backwards
	if subLoc == charLoc {
		return " next to you."
	}
	diff := math.Abs(float64(subLoc - charLoc))
	steps := ""
	direction := ""
	if subLoc > charLoc {
		direction = "in front of you."
	} else {
		direction = "behind you."
	}
	if diff == 1 {
		steps = " a couple steps "
	} else if diff == 2 {
		steps = " a dozen steps "
	} else if diff == 3 {
		steps = " many steps "
	} else if diff == 4 {
		steps = " at the other end of the room "
	}
	return steps + direction
}
