package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"sort"
	"strconv"
	"strings"
)

// Syntax: COMANDS
func init() {
	addHandler(help{},
		"",
		permissions.Player,
		"HELP", "COMMANDS", "CMDS")
}

// Width of gutter between columns
const gutter = 2

type help cmd

func (help) process(s *state) {

	if len(s.words) < 1 {
		s.msg.Actor.SendGood("Available commands:\n" +
			"(Type help cmd_name for more details, or emotes for a list of emotes, or spell_list for a list of spells, or song_list for a list of songs)\n" +
			"=====================================================")
		cmds := make([]string, len(helpText), len(helpText))

		// Extract keys from handler map
		pos, ommit := 0, 0
		for cmd := range helpText {

			// Ommit empty handler if installed and special commands starting with '#'
			// and scripting commands starting with '$'
			if len(cmd) == 0 || cmd[0] == '#' || cmd[0] == '$' || !s.actor.Permission.HasFlag(handlerPermission[strings.ToUpper(cmd)]) {
				ommit++
				continue
			}

			cmds[pos] = cmd
			pos++
		}

		// Reslice to remove omitted slots
		cmds = cmds[0 : len(cmds)-ommit]

		maxWidth := 0
		for _, cmd := range cmds {
			if l := len(cmd); l > maxWidth {
				maxWidth = l
			}
		}

		sort.Strings(cmds)

		var (
			columnWidth = maxWidth + gutter
			columnCount = 100 / columnWidth
			rowCount    = len(cmds) / columnCount
		)

		// If we have a partial row we need to account for it
		if len(cmds) > rowCount*columnCount {
			rowCount++
		}

		// NOTE: cell is not (row * columnCount) + column as we are pivoting the
		// table so that the commands are alphabetical DOWN the rows not across the
		// columns.
		for row := 0; row < rowCount; row++ {
			var line []byte
			for column := 0; column < columnCount; column++ {
				cell := (column * rowCount) + row
				if cell < len(cmds) {
					line = append(line, cmds[cell]...)
					line = append(line, strings.Repeat(" ", columnWidth-len(cmds[cell]))...)
				}
			}
			s.msg.Actor.Send(string(line))
		}

		s.ok = true
	} else {
		// Here we return the help text
		subject := s.words[0]

		// shortcut the help text to add emotes, spells and song lists
		if subject == "EMOTES" {
			s.msg.Actor.SendGood("All Emotes:\n" +
				"(Type help emote for more details)\n" +
				"=====================================================")
			cmds := make([]string, len(emotes), len(emotes))

			// Extract keys from handler map
			pos := 0
			for _, emote := range emotes {
				cmds[pos] = emote
				pos++
			}

			maxWidth := 0
			for _, spell := range cmds {
				if l := len(spell); l > maxWidth {
					maxWidth = l
				}
			}

			sort.Strings(cmds)

			var (
				columnWidth = maxWidth + gutter
				columnCount = 80 / columnWidth
				rowCount    = len(cmds) / columnCount
			)

			if len(cmds) > rowCount*columnCount {
				rowCount++
			}

			for row := 0; row < rowCount; row++ {
				var line []byte
				for column := 0; column < columnCount; column++ {
					cell := (column * rowCount) + row
					if cell < len(cmds) {
						line = append(line, cmds[cell]...)
						line = append(line, strings.Repeat(" ", columnWidth-len(cmds[cell]))...)
					}
				}
				s.msg.Actor.Send(string(line))
			}

			s.ok = true
			return
		}

		if subject == "SPELL_LIST" {
			s.msg.Actor.SendGood("All Spells:\n" +
				"(Type help spell_name for more details)\n" +
				"=====================================================")
			cmds := make([]string, len(objects.Spells), len(objects.Spells))

			// Extract keys from handler map
			pos := 0
			for spell := range objects.Spells {
				cmds[pos] = spell
				pos++
			}

			maxWidth := 0
			for _, spell := range cmds {
				if l := len(spell); l > maxWidth {
					maxWidth = l
				}
			}

			sort.Strings(cmds)

			var (
				columnWidth = maxWidth + gutter
				columnCount = 80 / columnWidth
				rowCount    = len(cmds) / columnCount
			)

			if len(cmds) > rowCount*columnCount {
				rowCount++
			}

			for row := 0; row < rowCount; row++ {
				var line []byte
				for column := 0; column < columnCount; column++ {
					cell := (column * rowCount) + row
					if cell < len(cmds) {
						line = append(line, cmds[cell]...)
						line = append(line, strings.Repeat(" ", columnWidth-len(cmds[cell]))...)
					}
				}
				s.msg.Actor.Send(string(line))
			}

			s.ok = true
			return
		}

		if subject == "SONG_LIST" {
			s.msg.Actor.SendGood("All Songs:\n" +
				"(Type help song_name for more details)\n" +
				"=====================================================")
			cmds := make([]string, len(objects.Songs), len(objects.Songs))

			// Extract keys from handler map
			pos := 0
			for song := range objects.Songs {
				cmds[pos] = song
				pos++
			}

			maxWidth := 0
			for _, song := range cmds {
				if l := len(song); l > maxWidth {
					maxWidth = l
				}
			}

			sort.Strings(cmds)

			var (
				columnWidth = maxWidth + gutter
				columnCount = 80 / columnWidth
				rowCount    = len(cmds) / columnCount
			)

			if len(cmds) > rowCount*columnCount {
				rowCount++
			}

			for row := 0; row < rowCount; row++ {
				var line []byte
				for column := 0; column < columnCount; column++ {
					cell := (column * rowCount) + row
					if cell < len(cmds) {
						line = append(line, cmds[cell]...)
						line = append(line, strings.Repeat(" ", columnWidth-len(cmds[cell]))...)
					}
				}
				s.msg.Actor.Send(string(line))
			}

			s.ok = true
			return
		}

		// Search Spells:
		if spell, ok := objects.Spells[strings.ToLower(subject)]; ok {
			s.msg.Actor.SendGood("Spell: ", subject, "\n", spell.Description, "\n")
			s.msg.Actor.SendGood("Chant: ", spell.Chant, "\n")
			s.msg.Actor.SendGood("Mana Cost: ", strconv.Itoa(spell.Cost), "\n")
			s.msg.Actor.SendGood("Castable by: ")
			for class, val := range spell.Classes {
				s.msg.Actor.SendGood("	", class, ": Tier ", strconv.Itoa(val))
			}

			return
		}

		if song, ok := objects.Songs[strings.ToLower(subject)]; ok {
			s.msg.Actor.SendGood("Song: ", subject, "\n\n", song["desc"])
			return
		}

		if _, ok := reverseLookup[strings.ToUpper(subject)]; ok {
			subject = reverseLookup[strings.ToUpper(subject)]
		}
		if s.actor.Permission.HasFlag(handlerPermission[strings.ToUpper(subject)]) {
			s.msg.Actor.SendGood("Command: ", subject, "\n\n", helpText[subject].helpText, "\n\n", "Aliases:", helpText[subject].aliases)
		} else {
			s.msg.Actor.SendBad("Not a command available to you.")
		}
	}
}
