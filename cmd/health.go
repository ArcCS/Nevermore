package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"log"
	"text/template"
)

func init() {
	addHandler(health{},
		"Displays quick view of health/stam mana",
		permissions.Player,
		"STAT", "HEALTH", "HEA", "VIT", "STAM", "MANA")
}

type health cmd

func (health) process(s *state) {

	if s.actor.CheckFlag("berserk") {
		s.msg.Actor.SendGood("The red rage has you, there's no time for that.\n")
		return
	}

	char_template := "You have {{.Stamina}}/{{.Max_stamina}} stamina, {{.Health}}/{{.Max_health}} health, and {{.Mana}}/{{.Max_mana}} mana pts.\n" +
		"{{if .Poisoned}}" + text.Red + "You have poison coursing through your veins.\n{{end}}" + text.Good +
		"{{if .Diseased}}" + text.Brown + "You are suffering from affliction.\n{{end}}" + text.Good

	data := struct {
		Stamina     int
		Max_stamina int
		Health      int
		Max_health  int
		Mana        int
		Max_mana    int
		Poisoned    bool
		Diseased    bool
	}{
		s.actor.Stam.Current,
		s.actor.Stam.Max,
		s.actor.Vit.Current,
		s.actor.Vit.Max,
		s.actor.Mana.Current,
		s.actor.Mana.Max,
		s.actor.CheckFlag("poisoned"),
		s.actor.CheckFlag("disease"),
	}

	tmpl, _ := template.New("char_info").Parse(char_template)
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	} else {
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
