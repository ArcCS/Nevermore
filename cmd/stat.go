package cmd

import (
	"bytes"
	"log"
	"text/template"
)

func init() {
	addHandler(stat{}, "STAT", "STATS")
	addHelp("Usage:  stats \n \n Show your current stat line", 0, "stat")
}

type stat cmd

func (stat) process(s *state) {

	stat_template := "You have {{.Stamina}}/{{.Max_stamina}} stamina, {{.Health}}/{{.Max_health}} health, and {{.Mana}}/{{.Max_mana}} mana pts.\n"

	data := struct {
		Stamina int64
		Max_stamina int64
		Health int64
		Max_health int64
		Mana int64
		Max_mana int64
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
	}else{
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
