package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
	"text/template"
)

func init() {
	addHandler(evaluate{},
		"Usage:  evaluate target\n\n  Examine a monster or item to find it's properties. ",
		permissions.Anyone,
		"evaluate", "eval")
}

type evaluate cmd

func (evaluate) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What do you want to evaluate?")
		return
	}

	if s.actor.Evals <= 0 {
		s.msg.Actor.SendBad("You cannot perform anymore evaluations today.")
		return
	}

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	/*
		// Check mobs
		var whatMob *objects.Mob
		whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
		// It was a mob!
		if whatMob != nil {
			s.actor.Evals -= 1
			s.msg.Actor.SendInfo(whatMob.Eval())
			return
		}

	*/

	// Check items
	whatItem := s.where.Items.Search(name, nameNum)

	// Item in the room?
	if whatItem != nil {
		s.actor.Evals -= 1
		s.msg.Actor.SendInfo(whatItem.Eval())
		return
	}

	whatItem = s.actor.Inventory.Search(name, nameNum)

	// It was on you the whole time
	if whatItem != nil {
		s.actor.Evals -= 1
		s.msg.Actor.SendInfo(whatItem.Eval())
		return
	}

	whatItem = s.actor.Equipment.Search(name, nameNum)

	// Check your equipment
	if whatItem != nil {
		s.actor.Evals -= 1
		s.msg.Actor.SendInfo(whatItem.Eval())
		return
	}

	// Check people

	whatChar := s.where.Chars.Search(name, s.actor)

	if whatChar != nil {
		s.actor.Evals -= 1
		berz, ok := s.actor.Flags["berserk"]
		if !ok {
			berz = false
		}
		monk := false
		if s.actor.Class == 8 {
			monk = true
		}

		age := (config.ImperialYearStart + objects.YearPlus) - s.actor.Birthyear

		charTemplate := "{{.Charname}}, the {{.Tier}} tier {{.Race}} {{.Title}}\n" +
			"----------------------------------------------------------------------\n" +
			"Str: {{.Str}}/{{.MaxStr}}, Dex: {{.Dex}}/{{.MaxDex}}, Con: {{.Con}}/{{.MaxCon}}, Int: {{.Int}}/{{.MaxInt}}, Piety: {{.Pie}}/{{.MaxPie}}.\n" +
			"They have an armor resistance of {{.ArmorResistance}}.\n" +
			"{{if .God}} You bear the mark of a devotee of {{.God}}.\n{{end}}" +
			"{{if .Berz}}" + text.Red + "They are in the throes of the red rage!\n" + text.Good +
			"{{else}}They have {{.Stamina}}/{{.MaxStamina}} stamina, {{.Health}}/{{.MaxHealth}} health, and {{.Mana}}/{{.MaxMana}} {{if .Monk}}chi{{else}}mana{{end}} pts.{{end}}\n" +
			"They require {{.NextLevel}} additional experience pts for their next tier.\n" +
			"They are carrying {{.Gold}} gold marks in their coin purse.\n" +
			"{{if .Poisoned}}" + text.Red + "They have poison coursing through their veins.\n{{end}}" + text.Good +
			"{{if .Diseased}}" + text.Brown + "They are suffering from affliction.\n{{end}}" + text.Good +
			"{{if .Blind}}" + text.Blue + "They have been blinded!!\n{{end}}" + text.Good +
			"{{if .DarkVision}}They can see in the dark naturally. \n{{end}}" +
			"They were born on {{.Day}}, the {{.DayNumber}} of the month of {{.Month}}\n" +
			"in the year {{.GodsYear}} since the Godswar, and year {{.EmpYear}} of the Empire.\n"

		data := struct {
			Charname        string
			Tier            string
			Race            string
			Title           string
			Str             int
			MaxStr          int
			Dex             int
			MaxDex          int
			Con             int
			MaxCon          int
			Int             int
			MaxInt          int
			Pie             int
			MaxPie          int
			ArmorResistance int
			God             string
			Stamina         int
			MaxStamina      int
			Health          int
			MaxHealth       int
			Mana            int
			MaxMana         int
			Monk            bool
			NextLevel       int
			Gold            int
			Poisoned        bool
			Diseased        bool
			Blind           bool
			DarkVision      bool
			Day             string
			DayNumber       string
			Month           string
			Age             int
			GodsYear        int
			EmpYear         int
			Berz            bool
		}{
			whatChar.Name,
			config.TextTiers[whatChar.Tier],
			config.AvailableRaces[whatChar.Race],
			whatChar.ClassTitle,
			whatChar.GetStat("str"),
			whatChar.Str.Max,
			whatChar.GetStat("dex"),
			whatChar.Dex.Max,
			whatChar.GetStat("con"),
			whatChar.Con.Max,
			whatChar.GetStat("int"),
			whatChar.Int.Max,
			whatChar.GetStat("pie"),
			whatChar.Pie.Max,
			whatChar.GetStat("armor"),
			"",
			whatChar.Stam.Current,
			whatChar.Stam.Max,
			whatChar.Vit.Current,
			whatChar.Vit.Max,
			whatChar.Mana.Current,
			whatChar.Mana.Max,
			monk,
			config.TierExpLevels[whatChar.Tier+1] - whatChar.Experience.Value,
			whatChar.Gold.Value,
			whatChar.CheckFlag("poisoned"),
			whatChar.CheckFlag("diseased"),
			whatChar.CheckFlag("blind"),
			whatChar.CheckFlag("darkvision"),
			utils.Title(config.Days[whatChar.Birthday]),
			config.PrintNumbers[whatChar.Birthdate],
			utils.Title(config.Months[whatChar.Birthmonth]["name"].(string)),
			age,
			2705 - age,
			2228 - age,
			berz,
		}

		tmpl, _ := template.New("char_info").Parse(charTemplate)
		var output bytes.Buffer
		err := tmpl.Execute(&output, data)
		if err != nil {
			log.Println(err)
		} else {
			s.msg.Actor.SendGood(output.String())
		}

		return
	}

	s.ok = true
	s.msg.Actor.SendInfo("Could not find anything to evaluate based on your input.")
	return

}
