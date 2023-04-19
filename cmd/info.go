package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"log"
	"strings"
	"text/template"
)

// Syntax: ( INFORMATION | INFO | INF | ME | STATS)
func init() {
	addHandler(information{},
		"Displays all of your current character information.",
		permissions.Player,
		"INF", "INFORMATION", "STATS", "INFO")
}

type information cmd

func (information) process(s *state) {
	// Do a save just because.
	s.actor.Save()

	berz, ok := s.actor.Flags["berserk"]
	if !ok {
		berz = false
	}
	monk := false
	if s.actor.Class == 8 {
		monk = true
	}

	age := (config.ImperialYearStart + objects.YearPlus) - s.actor.Birthyear

	char_template := "{{.Charname}}, the {{.Tier}} tier {{.Race}} {{.Title}}\n" +
		"----------------------------------------------------------------------\n" +
		"Str: {{.Str}}/{{.Max_str}}, Dex: {{.Dex}}/{{.Max_dex}}, Con: {{.Con}}/{{.Max_con}}, Int: {{.Int}}/{{.Max_int}}, Piety: {{.Pie}}/{{.Max_pie}}.\n" +
		"You have an armor resistance of {{.Armor_resistance}}.\n" +
		"{{if .God}} You bear the mark of a devotee of {{.God}}.\n{{end}}" +
		"{{if .Berz}}" + text.Red + "The red rage grips you!" + text.Good +
		"{{else}}You have {{.Stamina}}/{{.Max_stamina}} stamina, {{.Health}}/{{.Max_health}} health, and {{.Mana}}/{{.Max_mana}} {{if .Monk}}chi{{else}}mana{{end}} pts.{{end}}\n" +
		"You require {{.Next_level}} additional experience pts for your next tier.\n" +
		"You are carrying {{.Gold}} gold marks in your coin purse.\n" +
		"{{if .Dark_vision}} You can see in the dark naturally. \n{{end}}" +
		"You have {{.Broadcasts}} broadcasts remaining today.\n" +
		"You have {{.Evals}} evaluates remaining today.\n" +
		"You have logged {{.Hours}} hours and {{.Minutes}} minutes with this character.\n" +
		"You have {{.Bonus_points}} role-play bonus points.\n" +
		"You were born on {{.Day}}, the {{.Day_number}} of the month of {{.Month}}\n" +
		"in the year {{.GodsYear}} since the Godswar, and year {{.EmpYear}} of the Empire.\n" +
		"You are {{.Age}} years old.\n\n"

	data := struct {
		Charname         string
		Tier             string
		Race             string
		Title            string
		Str              int
		Max_str          int
		Dex              int
		Max_dex          int
		Con              int
		Max_con          int
		Int              int
		Max_int          int
		Pie              int
		Max_pie          int
		Armor_resistance int
		God              string
		Stamina          int
		Max_stamina      int
		Health           int
		Max_health       int
		Mana             int
		Max_mana         int
		Monk             bool
		Next_level       int
		Gold             int
		Dark_vision      bool
		Broadcasts       int
		Evals            int
		Hours            int
		Minutes          int
		Bonus_points     int
		Day              string
		Day_number       string
		Month            string
		Age              int
		Berz             bool
		GodsYear         int
		EmpYear          int
	}{
		s.actor.Name,
		config.TextTiers[s.actor.Tier],
		config.AvailableRaces[s.actor.Race],
		s.actor.ClassTitle,
		s.actor.GetStat("str"),
		s.actor.Str.Max,
		s.actor.GetStat("dex"),
		s.actor.Dex.Max,
		s.actor.GetStat("con"),
		s.actor.Con.Max,
		s.actor.GetStat("int"),
		s.actor.Int.Max,
		s.actor.GetStat("pie"),
		s.actor.Pie.Max,
		s.actor.GetStat("armor"),
		"",
		s.actor.Stam.Current,
		s.actor.Stam.Max,
		s.actor.Vit.Current,
		s.actor.Vit.Max,
		s.actor.Mana.Current,
		s.actor.Mana.Max,
		monk,
		config.TierExpLevels[s.actor.Tier+1] - s.actor.Experience.Value,
		s.actor.Gold.Value,
		s.actor.CheckFlag("darkvision"),
		s.actor.Broadcasts,
		s.actor.Evals,
		s.actor.MinutesPlayed / 60,
		s.actor.MinutesPlayed % 60,
		s.actor.BonusPoints.Value,
		strings.Title(config.Days[s.actor.Birthday]),
		config.PrintNumbers[s.actor.Birthdate],
		strings.Title(config.Months[s.actor.Birthmonth]["name"].(string)),
		age,
		berz,
		2705 - age,
		2228 - age,
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
