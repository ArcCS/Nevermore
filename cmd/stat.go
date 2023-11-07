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

	statTemplate := "You have {{.Stamina}}/{{.MaxStamina}} stamina, {{.Health}}/{{.MaxHealth}} health, and {{.Mana}}/{{.MaxMana}} mana pts.\n"

	data := struct {
		Stamina    int
		MaxStamina int
		Health     int
		MaxHealth  int
		Mana       int
		MaxMana    int
	}{
		s.actor.Stam.Current,
		s.actor.Stam.Max,
		s.actor.Vit.Current,
		s.actor.Vit.Max,
		s.actor.Mana.Current,
		s.actor.Mana.Max,
	}

	tmpl, _ := template.New("stat_info").Parse(statTemplate)
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	} else {
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
