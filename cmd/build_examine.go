package cmd

import (
	"bytes"
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

	switch typeE {
	case "room":
		roomRef := s.where
		if len(s.words) > 1 {
			roomNumber, ok := strconv.Atoi(s.words[1])
			if ok == nil{
				roomRef = objects.Rooms[int64(roomNumber)]
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

		if exitRef, ok := s.where.Exits[strings.ToLower(s.words[1])]; ok {
			exitTemplate := `Examining exit...
	name={{.ExitName}}		Exit Name
	description={{.Description}}
	placement={{.Placement}}		Exit Placement in the room
	key_id={{.Key_id}}		Key Id that can open this door
	Toggle Flags:
	  closeable={{.Closeable}},		Can the door be closed
	  closed={{.Closed}},			Is the door closed on start
	  autoclose={{.Autoclose}},		Does this door close itself
	  lockable={{.Lockable}},           Can it be locked
	  unpickable={{.Unpickable}},         Can it be picked
	  locked={{.Locked}},             Is it locked on start
	  hidden={{.Hidden}},             Is the exit hidden
	  invisible={{.Invisible}},	Is the exit invisible
	  levitate={{.Levitate}},	Does the character have to leviate to access
	  day_only={{.Day_only}}	Only accessible during the day
	  night_only={{.Night_only}}		Only accessible during the night
      placement_dependent={{.Placement_dependent}}		The character must be in the same placement to use it`

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
				strconv.Itoa(int(exitRef.Placement)),
				strconv.Itoa(int(exitRef.KeyId)),
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
		if exitRef, ok := s.where.Exits[strings.ToLower(s.words[1])]; ok {
			exitTemplate := `Examining exit...
	name={{.ExitName}}		Exit Name
	description={{.Description}}
	placement={{.Placement}}		Exit Placement in the room
	key_id={{.Key_id}}		Key Id that can open this door
	Toggle Flags:
	  closeable={{.Closeable}},		Can the door be closed
	  closed={{.Closed}},			Is the door closed on start
	  autoclose={{.Autoclose}},		Does this door close itself
	  lockable={{.Lockable}},           Can it be locked
	  unpickable={{.Unpickable}},         Can it be picked
	  locked={{.Locked}},             Is it locked on start
	  hidden={{.Hidden}},             Is the exit hidden
	  invisible={{.Invisible}},	Is the exit invisible
	  levitate={{.Levitate}},	Does the character have to leviate to access
	  day_only={{.Day_only}}	Only accessible during the day
	  night_only={{.Night_only}}		Only accessible during the night
      placement_dependent={{.Placement_dependent}}		The character must be in the same placement to use it`

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
				strconv.Itoa(int(exitRef.Placement)),
				strconv.Itoa(int(exitRef.KeyId)),
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
				strconv.Itoa(int(exitRef.Placement)),
				strconv.Itoa(int(exitRef.KeyId)),
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