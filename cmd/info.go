package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"log"
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
	singing, singOk := s.actor.Flags["singing"]
	if !singOk {
		singing = false
	}
	disprerolls := false
	if s.actor.Rerolls > 0 {
		disprerolls = true
	}

	showEnchants := false
	enchants := 0
	showHeals := false
	heals := 0
	showRestores := false
	restores := 0

	if s.actor.Class == 4 || s.actor.Class == 6 {
		showEnchants = true
		enchants = s.actor.ClassProps["enchants"]
	}
	if (s.actor.Class == 5 && s.actor.Tier >= 8) || (s.actor.Class == 6 && s.actor.Tier >= 12) {
		showHeals = true
		heals = s.actor.ClassProps["heals"]
	}
	if (s.actor.Class == 7 && s.actor.Tier >= 14) || (s.actor.Class == 6 && s.actor.Tier >= 13) {
		showRestores = true
		restores = s.actor.ClassProps["restores"]
	}

	age := (config.ImperialYearStart + objects.YearPlus) - s.actor.Birthyear

	char_template := "{{.Charname}}, the {{.Tier}} tier {{.Race}} {{.Title}}\n" +
		"----------------------------------------------------------------------\n" +
		"Str: {{.Str}}/{{.Max_str}}, Dex: {{.Dex}}/{{.Max_dex}}, Con: {{.Con}}/{{.Max_con}}, Int: {{.Int}}/{{.Max_int}}, Piety: {{.Pie}}/{{.Max_pie}}.\n" +
		"You have an armor resistance of {{.Armor_resistance}}.\n" +
		"{{if .God}} You bear the mark of a devotee of {{.God}}.\n{{end}}" +
		"{{if .Singing}}" + text.Cyan + "You are currently performing a song!\n{{end}}" + text.Good +
		"{{if .Berz}}" + text.Red + "The red rage grips you!\n" + text.Good +
		"{{else}}You have {{.Stamina}}/{{.Max_stamina}} stamina, {{.Health}}/{{.Max_health}} health, and {{.Mana}}/{{.Max_mana}} {{if .Monk}}chi{{else}}mana{{end}} pts.{{end}}\n" +
		"You require {{.Next_level}} additional experience pts for your next tier.\n" +
		"You are carrying {{.Gold}} gold marks in your coin purse.\n" +
		"{{if .Poisoned}}" + text.Red + "You have poison coursing through your veins.\n{{end}}" + text.Good +
		"{{if .Diseased}}" + text.Brown + "You are suffering from affliction.\n{{end}}" + text.Good +
		"{{if .Blind}}" + text.Blue + "You have been blinded!!\n{{end}}" + text.Good +
		"{{if .Dark_vision}}You can see in the dark naturally. \n{{end}}" +
		"You have {{.Broadcasts}} broadcasts remaining today.\n" +
		"You have {{.Evals}} evaluates remaining today.\n" +
		"{{if .ShowEnchants}}You can enchant {{.Enchants}} more items today.\n{{end}}" +
		"{{if .ShowHeals}}You can cast the heal spell {{.Heals}} more times today.\n{{end}}" +
		"{{if .ShowRestores}}You can cast the restore spell {{.Restores}} more times today.\n{{end}}" +
		"You have logged {{.Hours}} hours and {{.Minutes}} minutes with this character.\n" +
		"You have {{.Bonus_points}} role-play bonus points.\n" +
		"{{if .DispRerolls}}You can reroll your character {{.Rerolls}} more times.\n{{end}}" +
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
		Poisoned         bool
		Diseased         bool
		Blind            bool
		Dark_vision      bool
		ShowEnchants     bool
		ShowHeals        bool
		ShowRestores     bool
		Enchants         int
		Heals            int
		Restores         int
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
		Singing          bool
		GodsYear         int
		EmpYear          int
		DispRerolls      bool
		Rerolls          int
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
		s.actor.CheckFlag("poisoned"),
		s.actor.CheckFlag("diseased"),
		s.actor.CheckFlag("blind"),
		s.actor.CheckFlag("darkvision"),
		showEnchants,
		showHeals,
		showRestores,
		enchants,
		heals,
		restores,
		s.actor.Broadcasts,
		s.actor.Evals,
		s.actor.MinutesPlayed / 60,
		s.actor.MinutesPlayed % 60,
		s.actor.BonusPoints.Value,
		utils.Title(config.Days[s.actor.Birthday]),
		config.PrintNumbers[s.actor.Birthdate],
		utils.Title(config.Months[s.actor.Birthmonth]["name"].(string)),
		age,
		berz,
		singing,
		2705 - age,
		2228 - age,
		disprerolls,
		s.actor.Rerolls,
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
