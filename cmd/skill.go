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



	skill_header :=
`Skill                Level of Mastery      Experience
-----------------------------------------------------------------
`
standard_skills :=
`Sharp Weapons        {{.Sharp}}
Thrust Weapons       {{.Thrust}}
Blunt Weapons        {{.Blunt}}
Pole Weapons         {{.Pole}}
Missile Weapons      {{.Missile}}
`

monk_skills :=
`hand-to-hand combat  {{.Unarmed}}`

	data := struct {
		Sharp   string
		Thrust  string
		Blunt   string
		Pole    string
		Missile string
		Unarmed string
	}{
		config.WeaponExpTitle(s.actor.Skills[0].Value, s.actor.Class),
		config.WeaponExpTitle(s.actor.Skills[1].Value, s.actor.Class),
		config.WeaponExpTitle(s.actor.Skills[2].Value, s.actor.Class),
		config.WeaponExpTitle(s.actor.Skills[3].Value, s.actor.Class),
		config.WeaponExpTitle(s.actor.Skills[4].Value, s.actor.Class),
		config.WeaponExpTitle(s.actor.Skills[5].Value, s.actor.Class),
	}

	if s.actor.Class == 8 {
		skill_header += monk_skills
	}else{
		skill_header += standard_skills
	}
	tmpl, _ := template.New("stat_info").Parse(skill_header)
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	} else {
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
