// Copyright 2017 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

// handler is the interface for command processing handlers.
type handler interface {
	process(*state)
}

// handlers is a list of commands and their handlers. addHandler should be used
// to add new handlers. dispatchHandler then uses this list to lookup the
// correct handler to invoke for a given command.
var handlers = map[string]handler{}
var handlerPermission = map[string]permissions.Permissions{}

// addHandler adds the given commands for the specified handler. The commands
// will automatically be uppercased. Each command and its aliases should
// register its handler in its init function. For example:
//
//	func init() {
//		addHandler(Look{}, "L", "LOOK")
//	}
//
// In this example the LOOK command and it's alias 'L' register an instance of
// Look as their handler. If a handler is added for an existing command or
// alias the original handler will be replaced.
func addHandler(h handler, helpText string, permission permissions.Permissions, cmds ...string ) {
	for _, cmd := range cmds {
		handlers[strings.ToUpper(cmd)] = h
		handlerPermission[strings.ToUpper(cmd)] = permission
	}
	if helpText != "" {
		addHelp(helpText, cmds...)
	}
}

// dispatchHandler runs the registered handler for the current state command.
// If a handler cannot be found a message will be written to the actor's output
// buffer.
//
// dispatchHandler will only allow scripting specific commands to be executed
// if the state.scripting field is set to true.
func dispatchHandler(s *state) {

	if len(s.cmd) > 0 {
		//log.Println("received",  s.cmd, s.actor.Permission)
		if s.cmd[0] == '$' && !s.scripting {
			s.msg.Actor.SendBad("Unknown command, type HELP to get a list of commands (2)")
			return
		}

		// Check the player stack for the command first
		if val, ok := s.actor.Menu[strings.ToLower(s.cmd)]; ok {
			s.scriptActor(val.Command, strings.Join(s.input, " "))
			return
		}
		s.actor.EmptyMenu()

		// Check the room stack for a command second:
		if val, ok := s.where.Commands[s.cmd]; ok {
			s.scriptAll(val.Command, strings.Join(s.input, " "))
			return
		}

		switch handler, valid := handlers[s.cmd]; {
		case valid:
			//if s.actor.Permission.HasFlag(handlerPermission[s.cmd]) {
			//	s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
			//}else {
				handler.process(s)
			//}
		default:
			s.msg.Actor.SendBad("Unknown command, type HELP to get a list of commands (3)")
		}

	}else{
		s.msg.Actor.SendBad("Unknown command, type HELP to get a list of commands (4)")
	}
}
