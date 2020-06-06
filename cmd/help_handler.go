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
	Text string
}

var helpText= map[string]HelpCmd{}

func addHelp(Text string, cmd ...string) {
	for _, cmd := range cmd {
		helpText[strings.ToUpper(cmd)] = HelpCmd{Text}
	}
}

