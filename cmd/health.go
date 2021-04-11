package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/permissions"
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

	char_template := "You have {{.Stamina}}/{{.Max_stamina}} stamina, {{.Health}}/{{.Max_health}} health, and {{.Mana}}/{{.Max_mana}} mana pts.\n"

	data := struct {
		Stamina     int
		Max_stamina int
		Health      int
		Max_health  int
		Mana        int
		Max_mana    int
	}{
		s.actor.Stam.Current,
		s.actor.Stam.Max,
		s.actor.Vit.Current,
		s.actor.Vit.Max,
		s.actor.Mana.Current,
		s.actor.Mana.Max,
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
