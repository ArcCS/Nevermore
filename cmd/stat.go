package cmd

import (
	"bytes"
	"log"
	"text/template"
)

func init() {
	addHandler(stat{},
		"Usage:  stats \n \n Show your current stat line",
		0,
		"STAT")
}

type stat cmd

func (stat) process(s *state) {
	berz, ok := s.actor.Flags["berserk"]
	if ok {
		if berz {
			s.msg.Actor.SendBad("You are within the grip of the red rage!")
			return
		}
	}

	stat_template := "You have {{.Stamina}}/{{.Max_stamina}} stamina, {{.Health}}/{{.Max_health}} health, and {{.Mana}}/{{.Max_mana}} mana pts.\n"

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

	tmpl, _ := template.New("stat_info").Parse(stat_template)
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	} else {
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
