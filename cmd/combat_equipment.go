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
           "equipment", "gear")
}

type equipment cmd

func (equipment) process(s *state) {

	equip_template := " You take a look at your equipment..." +
	" {{if .Chest}}\n{{.Sub_pronoun}} {{.Isare}} wearing {{.Chest}} about {{.Pos_pronoun}} body{{end}}" +
	" {{if .Neck}}\n{{.Sub_pronoun}} {{.Isare}} a {{.Neck}} around {{.Pos_pronoun}} neck.{{end}}" +
	" {{if .Main}}\n{{.Sub_pronoun}} {{.Isare}} holding a {{.Main}} in {{.Pos_pronoun}} main hand.{{end}}" +
	" {{if .Offhand}}\n{{.Sub_pronoun}} {{.Isare}} holding a {{.Offhand}} in {{.Pos_pronoun}} off hand.{{end}}" +
	" {{if .Arms}}\n{{.Sub_pronoun}} {{.Isare}} wearing some {{.Arms}} on {{.Pos_pronoun}} arms.{{end}}" +
	" {{if .Finger1}}\n{{.Sub_pronoun}} {{.HasHave}} a {{.Finger1}} on {{.Pos_pronoun}} finger.{{end}}" +
	" {{if .Finger2}}\n{{.Sub_pronoun}} {{.HasHave}} a {{.Finger2}} on {{.Pos_pronoun}} finger.{{end}}" +
	" {{if .Legs}}\n{{.Sub_pronoun}} {{.HasHave}} {{.Legs}} on {{.Pos_pronoun}} legs.{{end}}" +
	" {{if .Feet}}\n{{.Sub_pronoun}} {{.HasHave}} {{.Feet}} on {{.Pos_pronoun}} feet.{{end}}" +
 " {{if .Head}}\n{{.Sub_pronoun}} {{.Isare}} wearing {{.Head}}.{{end}}" +
	text.Reset

	data := struct {
		Sub_pronoun string
		Pos_pronoun string
		Isare string
		HasHave string
		Chest string
		Neck string
		Main string
		Offhand	string
		Arms string
		Finger1 string
		Finger2 string
		Legs string
		Feet string
		Head string
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
		s.actor.Equipment.GetText("feet"),
		s.actor.Equipment.GetText("head"),
	}

	tmpl, _ := template.New("char_info").Parse(equip_template)

	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	}else{
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
