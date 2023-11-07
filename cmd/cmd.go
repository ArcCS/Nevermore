// Copyright 2017 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

// init adds a handler for the empty command. See the process method for
// details.
func init() {
	addHandler(cmd{}, "", 0)
}

// cmd is the default type used to build commands.
type cmd struct{}

func (cmd) process(s *state) {
	s.ok = true
}
