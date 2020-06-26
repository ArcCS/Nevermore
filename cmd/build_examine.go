package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
	"strings"
	"text/template"
)

func init() {
	addHandler(examine{},
           "Usage:  examine (room|mob|object|exit) (name|#####) \n\n  Examine will display the item and all of it's modifiable properties",
           permissions.Builder,
           "examine")
}

type examine cmd

func (examine) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("What do you want to examine?")
		return
	}

	typeE := strings.ToLower(s.words[0])

	nameNum := 1

	switch typeE {
	case "room":
		roomRef := s.where
		if len(s.words) > 1 {
			roomNumber, ok := strconv.Atoi(s.words[1])
			if ok == nil{
				roomRef = objects.Rooms[roomNumber]
			}
		}

		roomTemplate := `name:	{{.RoomName}}
description:	{{.Description}}

	Toggle flags
	 mana_drain:		{{.Mana_drain}}
	 no_summon:		{{.No_summon}}
	 heal_fast:		{{.Heal_fast}}
	 no_teleport:		{{.No_teleport}}
	 no_scry:		{{.No_scry}}
	 dark_always:		{{.Dark_always}}
	 light_always:		{{.Light_always}}
	 natural_light:		{{.Natural_light}}
	 indoors:		{{.Indoors}}
	 fire:			{{.Fire}}
	 encounters_on:	{{.Encounters_on}}
	 no_word_of_recall:	{{.No_word_of_recall}}
	 water:			{{.Water}}
	 no_magic:		{{.No_magic}}
	 urban:			{{.Urban}}
	 underground:		{{.Underground}}
	 earth:			{{.Earth}}
	 wind:			{{.Wind}}`

		data := struct {
			RoomName          string
			Description       string
			Repair            string
			Mana_drain        string
			No_summon         string
			Heal_fast         string
			No_teleport       string
			No_scry           string
			Shielded          string
			Dark_always       string
			Light_always      string
			Natural_light     string
			Indoors           string
			Fire              string
			Encounters_on    string
			No_word_of_recall string
			Water             string
			No_magic          string
			Urban             string
			Underground       string
			Earth             string
			Wind              string
		}{
			roomRef.Name,
			roomRef.Description,
			strconv.FormatBool(roomRef.Flags["repair"]),
			strconv.FormatBool(roomRef.Flags["mana_drain"]),
			strconv.FormatBool(roomRef.Flags["no_summon"]),
			strconv.FormatBool(roomRef.Flags["heal_fast"]),
			strconv.FormatBool(roomRef.Flags["no_teleport"]),
			strconv.FormatBool(roomRef.Flags["no_scry"]),
			strconv.FormatBool(roomRef.Flags["shielded"]),
			strconv.FormatBool(roomRef.Flags["dark_always"]),
			strconv.FormatBool(roomRef.Flags["light_always"]),
			strconv.FormatBool(roomRef.Flags["natural_light"]),
			strconv.FormatBool(roomRef.Flags["indoors"]),
			strconv.FormatBool(roomRef.Flags["fire"]),
			strconv.FormatBool(roomRef.Flags["encounters_on"]),
			strconv.FormatBool(roomRef.Flags["no_word_of_recall"]),
			strconv.FormatBool(roomRef.Flags["water"]),
			strconv.FormatBool(roomRef.Flags["no_magic"]),
			strconv.FormatBool(roomRef.Flags["urban"]),
			strconv.FormatBool(roomRef.Flags["underground"]),
			strconv.FormatBool(roomRef.Flags["earth"]),
			strconv.FormatBool(roomRef.Flags["wind"]),
		}

		tmpl, _ := template.New("room_info").Parse(roomTemplate)
		var output bytes.Buffer
		err := tmpl.Execute(&output, data)
		if err != nil {
			log.Println(err)
		} else {
			s.msg.Actor.SendGood(output.String())
		}


	case "mob":
		if len(s.words) > 2 {
			// Try to snag a number off the list
			if val, err := strconv.Atoi(s.words[2]); err == nil {
				nameNum = val
			}
		}
		mobRef := s.where.Mobs.Search(s.words[1], nameNum, true)
		if mobRef != nil {
			exitTemplate := `Examining mob...
    mob_id={{.MobId}}		Database Mob Id
	name={{.MobName}}		Mob Name
	description={{.Description}}
	level={{.Level}}
	experience={{.Experience}}
	gold={{.Gold}}

	## Stats ##
	Con: {{.Con}}
	Str: {{.Str}}
	Int: {{.Int}}
	Dex: {{.Dex}}
	Pie: {{.Pie}}
	Mana:        {{.ManaMax}}
    Hit Points:  {{.HPMax}}

	Combat:
	ndice={{.Ndice}}		Number of Damage Dice
	sdice={{.Sdice}}		Number of Sides on Damage Dice
	pdice={{.Pdice}}		Addition to to dice roll value. 
	casting_probability={{.CastingProb}}  % Chance of casting a spell
	armor={{.Armor}}	Amount of armor

	Wander={{.NumWander}}   How many ticks before they wander away while not in combat
	Wimpy={{.WimpyValue}}  Amount of damage in a single hit before it tries to flee

	Toggle Flags:
	  hide_encounter={{.Hide_Encounter}},		Is it hidden when it shows up?
	  invisible={{.Invisible}},			Invisible
	  permament={{.Permanent}},		Does not despawn
	  hostile={{.Hostile}},           Hostile to players`

			data := struct {
				MobId string
				MobName string
				Description string
				Level string
				Experience string
				Gold string
				Con string
				Str string
				Int string
				Dex string
				Pie string
				ManaMax string
				HPMax string
				Ndice string
				Sdice string
				Pdice string
				CastingProb string
				Armor string
				NumWander string
				WimpyValue string
				Hide_Encounter string
				Invisible string
				Permanent string
				Hostile string
			}{
				strconv.Itoa(mobRef.MobId),
				mobRef.Name,
				mobRef.Description,
				strconv.Itoa(mobRef.Level),
				strconv.Itoa(mobRef.Experience),
				strconv.Itoa(mobRef.Experience),
				strconv.Itoa(mobRef.Con.Max),
				strconv.Itoa(mobRef.Str.Max),
				strconv.Itoa(mobRef.Int.Max),
				strconv.Itoa(mobRef.Dex.Max),
				strconv.Itoa(mobRef.Pie.Max),
				strconv.Itoa(mobRef.Mana.Max),
				strconv.Itoa(mobRef.Stam.Max),
				strconv.Itoa(mobRef.NumDice),
				strconv.Itoa(mobRef.SidesDice),
				strconv.Itoa(mobRef.PlusDice),
				strconv.Itoa(mobRef.ChanceCast),
				strconv.Itoa(mobRef.Armor),
				strconv.Itoa(mobRef.NumWander),
				strconv.Itoa(mobRef.WimpyValue),
				strconv.FormatBool(mobRef.Flags["hide_encounter"]),
				strconv.FormatBool(mobRef.Flags["invisible"]),
				strconv.FormatBool(mobRef.Flags["permanent"]),
				strconv.FormatBool(mobRef.Flags["hostile"]),

			}

			tmpl, _ := template.New("char_info").Parse(exitTemplate)
			var output bytes.Buffer
			err := tmpl.Execute(&output, data)
			if err != nil {
				log.Println(err)
			} else {
				s.msg.Actor.SendGood(output.String())
			}
		}else{
			s.msg.Actor.SendBad("Couldn't find the object in the current room.")
		}

	case "object":
		nameNum := 1
		if len(s.words) > 2 {
			// Try to snag a number off the list
			if val, err := strconv.Atoi(s.words[2]); err == nil {
				nameNum = val
			}
		}
		objRef := s.actor.Inventory.Search(s.words[1],  nameNum)
		if objRef != nil {
			exitTemplate := `Examining object...
	itemId={{.ItemId}}			Database ID of the object
	name={{.ObjectName}}		Object Name
    creator={{.Creator}}		Who made this (Or last renamed)
	description={{.Description}}
	weight={{.Weight}}		Item weight
	type={{.Type}			Type of item
	value={{.Value}}		What will this pawn for?
	max_uses={{.Uses}}      The number of times this can be used before breaking
	Combat Item Stats:
	ndice={{.Ndice}}		Number of Damage Dice
	sdice={{.Sdice}}		Number of Sides on Damage Dice
	pdice={{.Pdice}}		Addition to to dice roll value. 
	
	Combat Values: 
	Toggle Flags:
	  always_crit={{.AlwaysCrit}}  If yes, criticals every hit. 
	  permanent={{.Permanent}},		If yes, never respawns. 
	  magic={{.Magic}},			Is it magic?
	  no_take={{.NoTake}},		Prevent player from taking
	  light={{.Light}},           Does it shed light?
	  weightless_chest={{.WeightLessChest}},         Weight less holding of items`

			data := struct {
				ItemId string
				ObjectName string
				Creator string
				Description string
				Weight string
				Type string
				Value string
				Uses string
				Ndice string
				Sdice string
				Pdice string
				AlwaysCrit string
				Permanent string
				Magic string
				NoTake string
				Light string
				WeightLessChest string
			}{
				strconv.Itoa(objRef.ItemId),
				objRef.Name,
				objRef.Creator,
				objRef.Description,
				strconv.Itoa(objRef.Weight),
				config.ItemTypes[objRef.Type],
				strconv.Itoa(objRef.Value),
				strconv.Itoa(objRef.MaxUses),
				strconv.Itoa(objRef.NumDice),
				strconv.Itoa(objRef.SidesDice),
				strconv.Itoa(objRef.PlusDice),
				strconv.FormatBool(objRef.Flags["always_crit"]),
				strconv.FormatBool(objRef.Flags["permanent"]),
				strconv.FormatBool(objRef.Flags["magic"]),
				strconv.FormatBool(objRef.Flags["no_take"]),
				strconv.FormatBool(objRef.Flags["light"]),
				strconv.FormatBool(objRef.Flags["weightless_chest"]),
			}

			tmpl, _ := template.New("object_info").Parse(exitTemplate)
			var output bytes.Buffer
			err := tmpl.Execute(&output, data)
			if err != nil {
				log.Println(err)
			} else {
				s.msg.Actor.SendGood(output.String())
			}
		}else{
			s.msg.Actor.SendBad("Couldn't find the object in the current room.")
		}

		//noinspection ALL
	case "exit":
		if exitRef, ok := s.where.Exits[strings.ToLower(s.words[1])]; ok {
			exitTemplate := `Examining exit...
	name:		{{.ExitName}}
	description:	{{.Description}}
	placement:		{{.Placement}}	#Exit Placement in the room
	key_id:		{{.Key_id}}	#Key Id that can open this door
	Toggle Flags:
	 closeable:		{{.Closeable}}
	 closed:		{{.Closed}}
	 autoclose:		{{.Autoclose}}
	 lockable:		{{.Lockable}}
	 unpickable:		{{.Unpickable}}
	 locked:		{{.Locked}}
	 hidden:		{{.Hidden}}
	 invisible:		{{.Invisible}}
	 levitate:		{{.Levitate}}
	 day_only:		{{.Day_only}}
	 night_only:		{{.Night_only}}
	 placement_dependent:	{{.Placement_dependent}}		#The character must be in the same placement to use it`

			data := struct {
				ExitName string
				Description string
				Placement string
				Key_id string
				Closeable string
				Closed string
				Autoclose string
				Lockable string
				Unpickable string
				Locked string
				Hidden string
				Invisible string
				Levitate string
				Day_only string
				Night_only string
				Placement_dependent string
			}{
				exitRef.Name,
				exitRef.Description,
				strconv.Itoa(exitRef.Placement),
				strconv.Itoa(exitRef.KeyId),
				strconv.FormatBool(exitRef.Flags["closeable"]),
				strconv.FormatBool(exitRef.Flags["closed"]),
				strconv.FormatBool(exitRef.Flags["autoclose"]),
				strconv.FormatBool(exitRef.Flags["lockable"]),
				strconv.FormatBool(exitRef.Flags["unpickable"]),
				strconv.FormatBool(exitRef.Flags["locked"]),
				strconv.FormatBool(exitRef.Flags["hidden"]),
				strconv.FormatBool(exitRef.Flags["invisible"]),
				strconv.FormatBool(exitRef.Flags["levitate"]),
				strconv.FormatBool(exitRef.Flags["day_only"]),
				strconv.FormatBool(exitRef.Flags["night_only"]),
				strconv.FormatBool(exitRef.Flags["placement_dependent"]),

			}

			tmpl, _ := template.New("exit_info").Parse(exitTemplate)
			var output bytes.Buffer
			err := tmpl.Execute(&output, data)
			if err != nil {
				log.Println(err)
			} else {
				s.msg.Actor.SendGood(output.String())
			}
		}else{
			s.msg.Actor.SendBad("Couldn't find the exit in the current room.")
		}
	default:
		s.msg.Actor.SendBad("Couldn't figure out what to examine.")
	}

	s.ok = true
	return
}