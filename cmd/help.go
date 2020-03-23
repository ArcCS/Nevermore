package cmd

import (
	"sort"
	"strings"
)

// Syntax: COMANDS
func init() {
	addHandler(help{},  "HELP", "COMMANDS", "CMDS")
}

// Width of gutter between columns
const gutter = 2

type help cmd

func (help) process(s *state) {

	if len(s.words) < 1 {
		s.msg.Actor.SendGood("Available commands:\n" +
			"(Type help cmd_name for more details)\n" +
			"=====================================================")
		cmds := make([]string, len(helpText), len(helpText))

		// Extract keys from handler map
		pos, ommit := 0, 0
		for cmd := range helpText {

			// Ommit empty handler if installed and special commands starting with '#'
			// and scripting commands starting with '$'
			if len(cmd) == 0 || cmd[0] == '#' || cmd[0] == '$' || helpText[cmd].permission >= int(s.actor.Class) {
				ommit++
				continue
			}

			cmds[pos] = cmd
			pos++
		}

		// Reslice to remove omitted slots
		cmds = cmds[0 : len(cmds)-ommit]

		// Find longest key extracted
		maxWidth := 0
		for _, cmd := range cmds {
			if l := len(cmd); l > maxWidth {
				maxWidth = l
			}
		}

		sort.Strings(cmds)

		var (
			columnWidth = maxWidth + gutter
			columnCount = 80 / columnWidth
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
			line := []byte{}
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
	}else{
		// Here we return the help text
		subject := s.words[0]
		if s.actor.Class >= int64(helpText[subject].permission) {
			s.msg.Actor.SendGood("Command: ", subject, "\n\n", helpText[subject].helptext)
		}
	}
}
