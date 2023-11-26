package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
	"strings"
	"text/template"
)

func init() {
	addHandler(look{},
		"Usage:  look [object|exit|character|mob] # \n \n Put your peepers on something. (Also can use short hand L",
		permissions.Player,
		"LOOK", "L")
}

type look cmd

func (look) process(s *state) {
	// Check to see if this person can see
	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything!")
		return
	}

	// Check if they have darkvision, a light source, or if they are a GM
	if !s.actor.CheckFlag("darkvision") && !s.actor.CheckFlag("light") && !objects.Rooms[s.actor.ParentId].Flags["light_always"] && !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		// Check if they are flagged for a light source
		if objects.Rooms[s.actor.ParentId].Flags["dark_always"] || (objects.Rooms[s.actor.ParentId].Flags["natural_light"] && !objects.DayTime) {
			s.msg.Actor.SendBad("It's too dark to see anything!")
			return
		}
	}

	var others []string
	var mobs string
	var poses []string
	var mobAttacking string
	var charAttacking string
	if len(s.input) == 0 {
		roomLook := objects.Rooms[s.actor.ParentId]
		s.msg.Actor.SendInfo(roomLook.Look(s.actor))
		others = objects.Rooms[s.actor.ParentId].Chars.List(s.actor)
		poses = objects.Rooms[s.actor.ParentId].Chars.ListPoses(s.actor)
		mobs = objects.Rooms[s.actor.ParentId].Mobs.ReducedList(s.actor)
		mobAttacking = objects.Rooms[s.actor.ParentId].Mobs.ListAttackers(s.actor)
		charAttacking = objects.Rooms[s.actor.ParentId].Chars.ListAttackers(s.actor)
		if len(others) == 1 {
			s.msg.Actor.SendInfo(strings.Join(others, ", "), " is also here.")
		} else if len(others) > 1 {
			s.msg.Actor.SendInfo(strings.Join(others, ", "), " are also here.")
		}
		if len(poses) > 0 {
			s.msg.Actor.SendInfo(strings.Join(poses, "\n"))
		}
		if len(mobs) > 0 {
			s.msg.Actor.SendInfo("You see " + mobs)
		}
		permItems := roomLook.Items.PermanentReducedList()
		if len(permItems) > 0 {
			s.msg.Actor.SendInfo("You see " + permItems)
		}
		items := roomLook.Items.RoomReducedList()
		if len(items) > 0 {
			s.msg.Actor.SendInfo("You see " + items)
		}
		if len(mobAttacking) > 0 {
			s.msg.Actor.Send(text.Red + mobAttacking + text.Reset)
		}
		if len(charAttacking) > 0 {
			s.msg.Actor.Send(text.Bold + text.Green + charAttacking + text.Reset)
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
		equipTemplate := "{{if .Chest}}\n{{.SubPronoun}} {{.Isare}} wearing {{.Chest}} about {{.PosPronoun}} body{{end}}" +
			" {{if .Neck}}\n{{.SubPronoun}} {{.Isare}} wearing a {{.Neck}} around {{.PosPronoun}} neck.{{end}}" +
			" {{if .Main}}\n{{.SubPronoun}} {{.Isare}} holding a {{.Main}} in {{.PosPronoun}} main hand.{{end}}" +
			" {{if .Offhand}}\n{{.SubPronoun}} {{.Isare}} holding a {{.Offhand}} in {{.PosPronoun}} off hand.{{end}}" +
			" {{if .Arms}}\n{{.SubPronoun}} {{.Isare}} wearing some {{.Arms}} on {{.PosPronoun}} arms.{{end}}" +
			" {{if .Finger1}}\n{{.SubPronoun}} {{.HasHave}} a {{.Finger1}} on {{.PosPronoun}} finger.{{end}}" +
			" {{if .Finger2}}\n{{.SubPronoun}} {{.HasHave}} a {{.Finger2}} on {{.PosPronoun}} finger.{{end}}" +
			" {{if .Legs}}\n{{.SubPronoun}} {{.HasHave}} {{.Legs}} on {{.PosPronoun}} legs.{{end}}" +
			" {{if .Hands}}\n{{.SubPronoun}} {{.HasHave}} {{.Hands}} on {{.PosPronoun}} hands.{{end}}" +
			" {{if .Feet}}\n{{.SubPronoun}} {{.HasHave}} {{.Feet}} on {{.PosPronoun}} feet.{{end}}" +
			" {{if .Head}}\n{{.SubPronoun}} {{.Isare}} wearing a {{.Head}}.{{end}}" +
			text.Reset

		data := struct {
			SubPronoun string
			PosPronoun string
			Isare      string
			HasHave    string
			Chest      string
			Neck       string
			Main       string
			Offhand    string
			Arms       string
			Finger1    string
			Finger2    string
			Legs       string
			Hands      string
			Feet       string
			Head       string
		}{
			utils.Title(config.TextSubPronoun[whatChar.Gender]),
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
			whatChar.Equipment.GetText("hands"),
			whatChar.Equipment.GetText("feet"),
			whatChar.Equipment.GetText("head"),
		}

		tmpl, _ := template.New("char_info").Parse(equipTemplate)

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
		s.msg.Actor.SendInfo(whatExit.Look())
		if whatExit.Flags["placement_dependent"] {
			s.msg.Actor.SendInfo("It is" + utils.WhereAt(whatExit.Placement, s.actor.Placement))
		} else {
			s.msg.Actor.SendInfo("It can be used from anywhere in the room.")
		}
		return
	}

	// Check mobs
	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	// It was a mob!
	if whatMob != nil {
		s.msg.Actor.SendInfo(whatMob.Look())
		s.msg.Actor.SendInfo("It is standing" + utils.WhereAt(whatMob.Placement, s.actor.Placement))
		s.msg.Actor.SendInfo("It" + whatMob.ReturnState() + ".")
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
		s.msg.Actor.SendInfo("It is" + utils.WhereAt(what.Placement, s.actor.Placement))
		return
	}

	what = s.actor.Inventory.Search(name, nameNum)

	// It was on you the whole time
	if what != nil {
		s.msg.Actor.SendInfo("You examine " + what.Name)
		s.msg.Actor.SendInfo(what.Look())
		return
	}

	what = s.actor.Equipment.Search(name, nameNum)

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
