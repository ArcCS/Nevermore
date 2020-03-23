package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/config"
	"log"
	"text/template"
)

// Syntax: WHO
func init() {
	addHandler(skills{}, "skill", "skills")
	addHelp("Usage:  skill \n \n Display the current level of your various weapon skills", 0, "skills")
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
		Sharp string
		Thrust string
		Blunt string
		Pole string
		Missile string
	}{
		config.WeaponExpTitle(s.actor.SharpExperience.Value),
		config.WeaponExpTitle(s.actor.ThrustExperience.Value),
		config.WeaponExpTitle(s.actor.BluntExperience.Value),
		config.WeaponExpTitle(s.actor.PoleExperience.Value),
		config.WeaponExpTitle(s.actor.MissileExperience.Value),
	}

	tmpl, _ := template.New("stat_info").Parse(skill_template)
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	}else{
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
