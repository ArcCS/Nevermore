package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
	"text/template"
)

// Syntax: WHO
func init() {
	addHandler(skills{},
		"Usage:  skill \n \n Display the current level of your various weapon skills or elemental affinities",
		permissions.Player,
		"skill", "sk", "skills", "ski")
}

type skills cmd

func (skills) process(s *state) {
	skill_header :=
		`Weapon Skills   |    Level of Mastery (Experience/NextLevel)
----------------------------------------------------------------
`
	standard_skills :=
		`Sharp Weapons        {{.Sharp}} ({{.SharpTotal}}/{{.SharpNext}})
Thrust Weapons       {{.Thrust}} ({{.ThrustTotal}}/{{.ThrustNext}})
Blunt Weapons        {{.Blunt}} ({{.BluntTotal}}/{{.BluntNext}})
Pole Weapons         {{.Pole}} ({{.PoleTotal}}/{{.PoleNext}})
Missile Weapons      {{.Missile}} ({{.MissileTotal}}/{{.MissileNext}})
`

	monk_skills :=
		`Hand-to-Hand combat  {{.Unarmed}}    {{.UnarmedTotal}}/{{.UnarmedNext}}`

	mage_skills :=
		`
Elemental Affinity  |   Level of Attunement (Experience/NextLevel)
-----------------------------------------------------------------
Fire Affinity           {{.Fire}} ({{.FireTotal}}/{{.FireNext}})
Air Affinity            {{.Air}} ({{.AirTotal}}/{{.AirNext}})
Earth Affinity          {{.Earth}} ({{.EarthTotal}}/{{.EarthNext}})
Water Affinity          {{.Water}} ({{.WaterTotal}}/{{.WaterNext}})
`

	data := struct {
		Sharp        string
		Thrust       string
		Blunt        string
		Pole         string
		Missile      string
		Unarmed      string
		Fire         string
		Air          string
		Earth        string
		Water        string
		SharpTotal   string
		SharpNext    string
		ThrustTotal  string
		ThrustNext   string
		BluntTotal   string
		BluntNext    string
		PoleTotal    string
		PoleNext     string
		MissileTotal string
		MissileNext  string
		UnarmedTotal string
		UnarmedNext  string
		FireTotal    string
		FireNext     string
		AirTotal     string
		AirNext      string
		EarthTotal   string
		EarthNext    string
		WaterTotal   string
		WaterNext    string
	}{
		config.WeaponExpTitle(s.actor.Skills[0].Value, s.actor.Class),
		config.WeaponExpTitle(s.actor.Skills[1].Value, s.actor.Class),
		config.WeaponExpTitle(s.actor.Skills[2].Value, s.actor.Class),
		config.WeaponExpTitle(s.actor.Skills[3].Value, s.actor.Class),
		config.WeaponExpTitle(s.actor.Skills[4].Value, s.actor.Class),
		config.WeaponExpTitle(s.actor.Skills[5].Value, s.actor.Class),
		config.AffinityExpTitle(s.actor.Skills[6].Value),
		config.AffinityExpTitle(s.actor.Skills[7].Value),
		config.AffinityExpTitle(s.actor.Skills[8].Value),
		config.AffinityExpTitle(s.actor.Skills[9].Value),
		strconv.Itoa(s.actor.Skills[0].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[0].Value, s.actor.Class)),
		strconv.Itoa(s.actor.Skills[1].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[1].Value, s.actor.Class)),
		strconv.Itoa(s.actor.Skills[2].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[2].Value, s.actor.Class)),
		strconv.Itoa(s.actor.Skills[3].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[3].Value, s.actor.Class)),
		strconv.Itoa(s.actor.Skills[4].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[4].Value, s.actor.Class)),
		strconv.Itoa(s.actor.Skills[5].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[5].Value, s.actor.Class)),
		strconv.Itoa(s.actor.Skills[6].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[6].Value, s.actor.Class)),
		strconv.Itoa(s.actor.Skills[7].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[7].Value, s.actor.Class)),
		strconv.Itoa(s.actor.Skills[8].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[8].Value, s.actor.Class)),
		strconv.Itoa(s.actor.Skills[9].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[9].Value, s.actor.Class)),
	}

	if s.actor.Class == 8 {
		skill_header += monk_skills
	} else {
		skill_header += standard_skills
	}
	if s.actor.Class == 4 {
		skill_header += mage_skills
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
