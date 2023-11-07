package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"log"
	"text/template"
)

func init() {
	addHandler(equipment{},
		"Usage:  equipment \n\n Display currently equipped gear",
		permissions.Player,
		"equipment", "gear", "eq")
}

type equipment cmd

func (equipment) process(s *state) {

	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything!")
		return
	}

	equipTemplate := " You take a look at your equipment..." +
		" {{if .Chest}}\n{{.SubPronoun}} {{.Isare}} wearing {{.Chest}} about {{.PosPronoun}} body{{end}}" +
		" {{if .Neck}}\n{{.SubPronoun}} {{.Isare}} wearing a {{.Neck}} around {{.PosPronoun}} neck.{{end}}" +
		" {{if .Main}}\n{{.SubPronoun}} {{.Isare}} holding a {{.Main}} in {{.PosPronoun}} main hand.{{end}}" +
		" {{if .Offhand}}\n{{.SubPronoun}} {{.Isare}} holding a {{.Offhand}} in {{.PosPronoun}} off hand.{{end}}" +
		" {{if .Arms}}\n{{.SubPronoun}} {{.Isare}} wearing some {{.Arms}} on {{.PosPronoun}} arms.{{end}}" +
		" {{if .Finger1}}\n{{.SubPronoun}} {{.HasHave}} a {{.Finger1}} on {{.PosPronoun}} finger.{{end}}" +
		" {{if .Finger2}}\n{{.SubPronoun}} {{.HasHave}} a {{.Finger2}} on {{.PosPronoun}} finger.{{end}}" +
		" {{if .Legs}}\n{{.SubPronoun}} {{.HasHave}} {{.Legs}} on {{.PosPronoun}} legs.{{end}}" +
		" {{if .Hands}}\n{{.SubPronoun}} {{.HasHave}} {{.Hands}} on {{.PosPronoun}} hands.{{end}}" +
		" {{if .Feet}}\n{{.SubPronoun}} {{.HasHave}} {{.Feet}} on {{.PosPronoun}} feet.{{end}}" +
		" {{if .Head}}\n{{.SubPronoun}} {{.Isare}} wearing {{.Head}}.{{end}}" +
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
		"You",
		"your",
		"are",
		"have",
		s.actor.Equipment.GetText("chest"),
		s.actor.Equipment.GetText("neck"),
		s.actor.Equipment.GetText("main"),
		s.actor.Equipment.GetText("off"),
		s.actor.Equipment.GetText("arms"),
		s.actor.Equipment.GetText("ring1"),
		s.actor.Equipment.GetText("ring2"),
		s.actor.Equipment.GetText("legs"),
		s.actor.Equipment.GetText("hands"),
		s.actor.Equipment.GetText("feet"),
		s.actor.Equipment.GetText("head"),
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
}
