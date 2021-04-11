package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"text/template"
)

// Syntax: WHO
func init() {
	addHandler(skills{},
		"Usage:  skill \n \n Display the current level of your various weapon skills",
		permissions.Player,
		"skill")
}

type skills cmd

// TODO: Add experience to the display so users are aware of their current status

func (skills) process(s *state) {

	skill_template :=
		`Skill                Level of Mastery
-----------------------------------------------------------------
Sharp Weapons        {{.Sharp}}
Thrust Weapons       {{.Thrust}}
Blunt Weapons        {{.Blunt}}
Pole Weapons         {{.Pole}}
Missile Weapons      {{.Missile}}
`
	data := struct {
		Sharp   string
		Thrust  string
		Blunt   string
		Pole    string
		Missile string
	}{
		config.WeaponExpTitle(s.actor.Skills[0].Value),
		config.WeaponExpTitle(s.actor.Skills[1].Value),
		config.WeaponExpTitle(s.actor.Skills[2].Value),
		config.WeaponExpTitle(s.actor.Skills[3].Value),
		config.WeaponExpTitle(s.actor.Skills[4].Value),
	}

	tmpl, _ := template.New("stat_info").Parse(skill_template)
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	} else {
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
