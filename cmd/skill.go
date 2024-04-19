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
	skillHeader :=
		`Weapon Skills   |    Level of Mastery (Experience/NextLevel)
----------------------------------------------------------------
`
	standardSkills :=
		`Sharp Weapons        {{.Sharp}} ({{.SharpTotal}}/{{.SharpNext}})
Thrust Weapons       {{.Thrust}} ({{.ThrustTotal}}/{{.ThrustNext}})
Blunt Weapons        {{.Blunt}} ({{.BluntTotal}}/{{.BluntNext}})
Pole Weapons         {{.Pole}} ({{.PoleTotal}}/{{.PoleNext}})
Missile Weapons      {{.Missile}} ({{.MissileTotal}}/{{.MissileNext}})
`

	monkSkills :=
		`Hand-to-Hand combat  {{.Unarmed}}    {{.UnarmedTotal}}/{{.UnarmedNext}}`

	mageSkills :=
		`
Elemental Affinity  |   Level of Attunement (Experience/NextLevel)
-----------------------------------------------------------------
Fire Affinity           {{.Fire}} ({{.FireTotal}}/{{.FireNext}})
Air Affinity            {{.Air}} ({{.AirTotal}}/{{.AirNext}})
Earth Affinity          {{.Earth}} ({{.EarthTotal}}/{{.EarthNext}})
Water Affinity          {{.Water}} ({{.WaterTotal}}/{{.WaterNext}})
`

	healerSkills :=
		`
Divinity        |    Level of Devotion (Experience/NextLevel)
-----------------------------------------------------------------
Curative Arts        {{.Divinity}} ({{.DivinityTotal}}/{{.DivinityNext}})

`

	thiefSkills :=
		`
Clandestine Skills    |    Level of Skill (Experience/NextLevel)
-----------------------------------------------------------------
Covert Arts              {{.Stealth}} ({{.StealthTotal}}/{{.StealthNext}})		
`

	data := struct {
		Sharp         string
		Thrust        string
		Blunt         string
		Pole          string
		Missile       string
		Unarmed       string
		Fire          string
		Air           string
		Earth         string
		Water         string
		Divinity      string
		Stealth       string
		SharpTotal    string
		SharpNext     string
		ThrustTotal   string
		ThrustNext    string
		BluntTotal    string
		BluntNext     string
		PoleTotal     string
		PoleNext      string
		MissileTotal  string
		MissileNext   string
		UnarmedTotal  string
		UnarmedNext   string
		FireTotal     string
		FireNext      string
		AirTotal      string
		AirNext       string
		EarthTotal    string
		EarthNext     string
		WaterTotal    string
		WaterNext     string
		DivinityTotal string
		DivinityNext  string
		StealthTotal  string
		StealthNext   string
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
		config.DivinityExpTitle(s.actor.Skills[10].Value),
		config.StealthExpTitle(s.actor.Skills[11].Value),
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
		strconv.Itoa(s.actor.Skills[10].Value),
		strconv.Itoa(config.WeaponExpNext(s.actor.Skills[10].Value, s.actor.Class)),
		strconv.Itoa(s.actor.Skills[11].Value),
		strconv.Itoa(config.StealthExpNext(s.actor.Skills[11].Value)),
	}

	if s.actor.Class == 8 {
		skillHeader += monkSkills
	} else {
		skillHeader += standardSkills
	}
	if s.actor.Class == 4 {
		skillHeader += mageSkills
	}
	if s.actor.Class == 2 {
		skillHeader += thiefSkills
	}
	if s.actor.Class == 5 || s.actor.Class == 6 {
		skillHeader += healerSkills
	}
	tmpl, _ := template.New("stat_info").Parse(skillHeader)
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	} else {
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
