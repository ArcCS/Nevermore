// Copyright 2017 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

import (
	"strings"
)

// handlers is a list of commands and their handlers. addHandler should be used
// to add new handlers. dispatchHandler then uses this list to lookup the
// correct handler to invoke for a given command.
type HelpCmd struct {
	permission int
	helptext string
}

var helpText= map[string]HelpCmd{}

func addHelp(helptext string, permission int, cmd ...string) {
	for _, cmd := range cmd {
		helpText[strings.ToUpper(cmd)] = HelpCmd{permission, helptext}
	}
}

