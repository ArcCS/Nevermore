package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/config"
	"log"
	"text/template"
)

// Syntax: ( INFORMATION | INFO | INF | ME | STATS)
func init() {
	addHandler(information{}, "INF", "INFORMATION", "STATS", "INFO")
}

type information cmd

func (information) process(s *state) {

	char_template := "{{.Charname}}, the {{.Tier}} tier {{.Race}} {{.Title}}\n" +
"----------------------------------------------------------------------\n" +
"Str: {{.Str}}/{{.Max_str}}, Dex: {{.Dex}}/{{.Max_dex}}, Con: {{.Con}}/{{.Max_con}}, Int: {{.Int}}/{{.Max_int}}, Piety: {{.Pie}}/{{.Max_pie}}.\n"+
"You have an armor resistance of {{.Armor_resistance}} with a damage ignore of {{.Damage_ignore}}.\n"+
"{{if .God}} You bear the mark of a devotee of {{.God}}.\n{{end}}"+
	"You have {{.Stamina}}/{{.Max_stamina}} stamina, {{.Health}}/{{.Max_health}} health, and {{.Mana}}/{{.Max_mana}} mana pts.\n"+
	"You require {{.Next_level}} additional experience pts for your next tier.\n"+
	"You are carrying {{.Gold}} gold marks in your coin purse.\n"+
	"{{if .Dark_vision}} You can see in the dark naturally. \n{{end}}"+
	"You have {{.Broadcasts}} broadcasts remaining today.\n"+
	"You have {{.Evals}} evaluates remaining today.\n"+
	"You have logged {{.Hours}} hours with this character.\n"+
	"You have {{.Bonus_points}} role-play bonus points.\n"+
	"You may move {{.Attr_moves}} attribute point.\n"+
	"You were born on {{.Day}}, the {{.Day_number}} of the month of {{.Month}}\n"+
	"You are {{.Age}} years old.\n\n"


	data := struct {
		Charname string
		Tier string
		Race string
		Title string
		Str int64
		Max_str int64
		Dex	int64
		Max_dex int64
		Con int64
		Max_con int64
		Int int64
		Max_int int64
		Pie int64
		Max_pie int64
		Armor_resistance int64
		Damage_ignore int64
		God string
		Stamina int64
		Max_stamina int64
		Health int64
		Max_health int64
		Mana int64
		Max_mana int64
		Next_level int64
		Gold int64
		Dark_vision bool
		Broadcasts int64
		Evals int64
		Hours int64
		Bonus_points int64
		Attr_moves int64
		Day string
		Day_number int64
		Month string
		Age int64
	}{
		s.actor.Name,
		config.TextTiers[s.actor.Tier],
		config.AvailableRaces[s.actor.Race],
		s.actor.ClassTitle,
		s.actor.Str.Current,
		s.actor.Str.Max,
		s.actor.Dex.Current,
		s.actor.Dex.Max,
		s.actor.Con.Current,
		s.actor.Con.Max,
		s.actor.Int.Current,
		s.actor.Int.Max,
		s.actor.Pie.Current,
		s.actor.Pie.Max,
		0,
		0,
		"",
		s.actor.Stam.Current,
		s.actor.Stam.Max,
		s.actor.Vit.Current,
		s.actor.Vit.Max,
		s.actor.Mana.Current,
		s.actor.Mana.Max,
		config.TierExpLevels[s.actor.Tier + 1] - s.actor.Experience.Value,
		s.actor.Gold.Value,
		false,
		s.actor.Broadcasts,
		s.actor.Evals,
		0,
		s.actor.BonusPoints.Value,
		s.actor.AttrMoves.Value,
		"Day",
		0,
		"Month",
		0,
	}

	tmpl, _ := template.New("char_info").Parse(char_template)
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	}else{
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
