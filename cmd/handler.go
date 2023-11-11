// Copyright 2017 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strings"
	"time"
)

// handler is the interface for command processing handlers.
type handler interface {
	process(*state)
}

type helpTextStruct struct {
	helpText string
	aliases  string
}

// handlers is a list of commands and their handlers. addHandler should be used
// to add new handlers. dispatchHandler then uses this list to look up the
// correct handler to invoke for a given command.
var handlers = map[string]handler{}
var handlerPermission = map[string]permissions.Permissions{}
var helpText = map[string]helpTextStruct{}
var oocCommands = []string{"SAY", "QUIT", "HELP", "WHO", "LOOK", "IC", "$POOF", "AFK", "GO", "ACT"}
var excludeFromLogs = []string{"SAYTO", "SAY", "TELL", "OSAY", "SEND", "R", "REPLY", "REP", "PARTYTELL", "PTELL"}
var reverseLookup = map[string]string{}
var emotes = []string{"ACT", "BLINK", "BLUSH", "BOW", "BURP", "CACKLE", "CHEER", "CHUCKLE", "CLAP", "CONFUSED", "COUGH", "CROSSARMS", "CROSSFINGERS", "CRY",
	"DANCE", "EMOTE", "FLEX", "FLINCH", "FROWN", "GASP", "GIGGLE", "GRIN", "GROAN", "HICCUP", "JUMP", "KNEEL", "LAUGH", "NOD", "PONDER", "SALUTE", "SHAKE", "SHIVER", "SHRUG",
	"SIGH", "SNEEZE", "SNAP", "SMILE", "SMIRK", "SNICKER", "SPIT", "STARE", "STRETCH", "TAP", "THUMBSDOWN", "THUMBSUP", "WAVE", "WHISTLE", "WINK", "YAWN",
	"BUG", "BOW", "HUG", "ANGRY", "GLARE", "STARE", "TICKLE", "POKE", "SLAP", "KICK", "WAVE", "WINK"}

// addHandler adds the given commands for the specified handler.
// It requires the command handler,  a help string to add to the help data, a bitmask permission, and the relative
// commands that will be each added to dispatch
func addHandler(h handler, helpString string, permission permissions.Permissions, cmds ...string) {
	primeString := ""
	if len(cmds) != 0 {
		primeString = strings.ToUpper(cmds[0])
	}
	if helpString != "" {
		helpText[strings.ToUpper(cmds[0])] = helpTextStruct{helpString, strings.Join(cmds[0:], ", ")}
	}
	for _, cmd := range cmds {
		if cmd != primeString {
			reverseLookup[strings.ToUpper(cmd)] = primeString
		}
		handlers[strings.ToUpper(cmd)] = h
		handlerPermission[strings.ToUpper(cmd)] = permission
	}
}

// dispatch handler takes the command sent and attempts to find it in a stack of command locations for execution
func dispatchHandler(s *state) {

	if len(s.cmd) > 0 {
		//if !s.scripting {
		if !utils.StringIn(strings.ToUpper(s.cmd), emotes) && !utils.StringIn(strings.ToUpper(s.cmd), excludeFromLogs) {
			log.Println(s.actor.Name + " sent " + s.cmd + " " + strings.Join(s.input, " "))
			s.actor.LastAction = time.Now()
		}
		//}

		if s.where.RoomId == config.OocRoom &&
			!s.actor.Permission.HasAnyFlags(permissions.Dungeonmaster, permissions.Gamemaster) &&
			!utils.StringIn(strings.ToUpper(s.cmd), oocCommands) {
			s.msg.Actor.SendBad("You must be IC to do that.")
			return
		}

		if s.cmd[0] == '$' && !s.scripting {
			s.msg.Actor.SendBad("Unknown command, type HELP to get a list of commands (2)")
			return
		}

		// Check the player stack for the command first
		completeCommand := s.cmd + " " + strings.Join(s.input, " ")
		if val, ok := s.actor.Commands[completeCommand]; ok {
			s.scriptActor(val.Command)
			return
		} else if val, ok := s.actor.Commands[s.cmd]; ok {
			s.scriptActor(val.Command, strings.Join(s.input, " "))
			return
		}
		s.actor.EmptyCommands()

		for _, i := range s.actor.Inventory.Contents {
			// Check the room stack for a command second:
			if val, ok := i.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := i.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Head != nil {
			if val, ok := s.actor.Equipment.Head.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Head.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Chest != nil {
			if val, ok := s.actor.Equipment.Chest.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Chest.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Neck != nil {
			if val, ok := s.actor.Equipment.Neck.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Neck.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Legs != nil {
			if val, ok := s.actor.Equipment.Legs.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Legs.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Feet != nil {
			if val, ok := s.actor.Equipment.Feet.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Feet.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Arms != nil {
			if val, ok := s.actor.Equipment.Arms.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Arms.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Hands != nil {
			if val, ok := s.actor.Equipment.Hands.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Hands.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Ring1 != nil {
			if val, ok := s.actor.Equipment.Ring1.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Ring1.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Ring2 != nil {
			if val, ok := s.actor.Equipment.Ring2.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Ring2.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Main != nil {
			if val, ok := s.actor.Equipment.Main.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Main.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		if s.actor.Equipment.Off != nil {
			if val, ok := s.actor.Equipment.Off.Commands[completeCommand]; ok {
				s.scriptAll(val.Command)
				return
			} else if val, ok := s.actor.Equipment.Off.Commands[s.cmd]; ok {
				s.scriptAll(val.Command, strings.Join(s.input, " "))
				return
			}
		}

		// Check the room stack for a command second:
		if val, ok := s.where.Commands[completeCommand]; ok {
			s.scriptAll(val.Command)
			return
		} else if val, ok := s.where.Commands[s.cmd]; ok {
			s.scriptAll(val.Command, strings.Join(s.input, " "))
			return
		}

		for _, i := range s.where.Items.Contents {
			// Check the room stack for a command second:
			if i.Flags["permanent"] {
				if i.Placement == s.actor.Placement {
					if val, ok := i.Commands[completeCommand]; ok {
						s.scriptAll(val.Command)
						return
					} else if val, ok := i.Commands[s.cmd]; ok {
						s.scriptAll(val.Command, strings.Join(s.input, " "))
						return
					}
				}
			}
		}

		for _, i := range s.where.Mobs.Contents {
			// Check the room stack for a command second:
			if i.Flags["permanent"] {
				if i.Placement == s.actor.Placement {
					if val, ok := i.Commands[completeCommand]; ok {
						s.scriptAll(val.Command)
						return
					} else if val, ok := i.Commands[s.cmd]; ok {
						s.scriptAll(val.Command, strings.Join(s.input, " "))
						return
					}
				}
			}
		}

		switch handler, valid := handlers[s.cmd]; {
		case valid:
			if s.actor.Permission.HasFlag(handlerPermission[s.cmd]) || s.actor.Permission.HasAnyFlags(permissions.Dungeonmaster, permissions.Gamemaster) {
				handler.process(s)
			} else {
				s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
			}
		default:
			s.msg.Actor.SendBad("Unknown command, type HELP to get a list of commands (3)")
		}

	} else {
		s.msg.Actor.SendBad("Unknown command, type HELP to get a list of commands (4)")
	}
}
