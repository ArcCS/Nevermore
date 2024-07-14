// Copyright 2016 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package frontend

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"sort"
	"strings"
)

// Start menu embeds a Frontend instance adding fields and methods specific to
// the main menu.
type Start struct {
	*Frontend
	optionEnd      int
	powerCharacter string
	characters     []string
	puppets        []string
	pwChange       [16]byte
}

func NewStart(f *Frontend) (m *Start) {
	m = &Start{Frontend: f}
	m.startDisplay()
	return
}

// menuDisplay shows the main menu of options available once a player is logged
// into the system.
func (m *Start) startDisplay() {
	// Load Characters
	for _, name := range data.ListChars(m.account) {
		if !utils.StringIn(name, m.characters) {
			m.characters = append(m.characters, name)
		}
	}
	sort.Strings(m.characters)
	m.powerCharacter, _ = data.ListPowerChar(m.account)

	if m.permissions.HasAnyFlags(permissions.Dungeonmaster, permissions.Gamemaster) {
		m.puppets = data.ListPuppets()
		sort.Strings(m.puppets)
	}

	var output strings.Builder
	m.optionEnd = 2
	charOption := ` X. Make a new character `
	if config.Server.CreateChars {
		charOption = ` 1. Make a new character `
	}

	output.WriteString(text.Good + `
 
=========
 Choose an action:
 ---------
 0. Quit` + "\n" + charOption + "\n" +
		` 2. Change account password
`)
	if m.permissions.HasFlag(permissions.Gamemaster) {
		if m.powerCharacter == "" {
			output.WriteString(" 3. Create a gamemaster account.\r\n")
		} else {
			output.WriteString(" 3. Load gamemaster account:" + m.powerCharacter + "\r\n")
		}
	} else if m.permissions.HasFlag(permissions.Builder) {
		if m.powerCharacter == "" {
			output.WriteString(" 3. Create a builder account.\r\n")
		} else {
			output.WriteString(" 3. Load builder account:" + m.powerCharacter + "\r\n")
		}
		m.optionEnd = 3
	} else if m.permissions.HasFlag(permissions.Dungeonmaster) {
		if m.powerCharacter == "" {
			output.WriteString(" 3. Create a dungeon master account.\r\n")
		} else {
			output.WriteString(" 3. Load dungeonmaster account:" + m.powerCharacter + "\r\n")
		}
		m.optionEnd = 3
	} else {
		m.optionEnd = 2
	}
	output.WriteString("\n==== Your Character List ====\n")
	output.WriteString("(Enter the name of the character you wish to play)\n")
	output.WriteString(m.characterList())

	m.buf.Send(output.String())
	m.nextFunc = m.startProcess
}

// menuProcess takes the current input and processes it as a menu option.
func (m *Start) startProcess() {
	switch string(m.input) {
	case "":
		return
	case "1":
		if config.Server.CreateChars {
			CreateNewChar(m.Frontend)
		} else {
			m.buf.Send(text.Bad, "New character creation is disabled at this time.", text.Reset)
		}
	case "0":
		// Say goodbye to client
		if _, err := m.Write([]byte(text.Info + "\nBye bye...\n\n")); err != nil {
			log.Println("Error writing to player:", err)
		}
		// Revert to default colors
		if _, err := m.Write([]byte(text.Reset)); err != nil {
			log.Println("Error writing to player:", err)
		}
		m.Close()
	case "2":
		m.buf.Send(text.Info, "Enter your new password:")
		m.nextFunc = m.verifyPw
	default:
		if m.permissions.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) && string(m.input) == "3" {
			if m.powerCharacter == "" {
				CreateNewPChar(m.Frontend)
				return
			} else {
				if objects.ActiveCharacters.Find(m.powerCharacter) == nil {
					StartGame(m.Frontend, m.powerCharacter)
				} else {
					m.buf.Send(text.Bad, "You're already in the game.  You cannot rejoin.", text.Reset)
					return
				}

			}
		} else if utils.StringInLower(string(m.input), m.characters) || utils.StringInLower(string(m.input), m.puppets) {

			if objects.ActiveCharacters.Find(string(m.input)) == nil {
				if m.accountAllowed() {
					StartGame(m.Frontend, string(m.input))
				} else {
					m.buf.Send(text.Bad, "Account already logged in.\n", text.Reset)
					return
				}
			} else {
				// Attempt to resume
				m.buf.Send(text.Good, "Resuming session...\n", text.Reset)
				var character = objects.ActiveCharacters.Find(string(m.input))
				log.Println("Disconnect old lease")
				character.Disconnect()
				log.Println("Resume game")
				ResumeGame(m.Frontend, character)
			}
		}
		m.buf.Send(text.Bad, "Invalid option selected. Please try again.", text.Reset)
	}
}

func (m *Start) accountAllowed() bool {
	if m.permissions < 16 {
		if _, ok := accounts.inuse[m.account]; ok {
			return false
		}
	}
	return true
}

func (m *Start) verifyPw() {
	switch l := len(m.input); {
	case l == 0:
		m.buf.Send(text.Bad, "No text sent, returning to menu.", text.Reset)
		m.startDisplay()
	default:
		m.pwChange = md5.Sum(m.input)
		m.buf.Send(text.Good, "Please type your password again to change.")
		m.nextFunc = m.changePw
	}
	return
}

func (m *Start) changePw() {
	switch l := len(m.input); {
	case l == 0:
		m.buf.Send(text.Info, "Password change cancelled.\n", text.Reset)
		m.startDisplay()
	default:
		if md5.Sum(m.input) != m.pwChange {
			m.buf.Send(text.Bad, "Passwords do not match, please try again. Or enter nothing to cancel\n", text.Reset)
			m.nextFunc = m.changePw
		} else {
			data.UpdatePassword(m.account, hex.EncodeToString(m.pwChange[:]))
			m.buf.Send(text.Good, "Password changed!\n")
			m.startDisplay()
		}

	}
	return
}

func (m *Start) characterList() string {
	var charList strings.Builder
	charList.WriteString("-----------------------------\n  ")
	charList.Write([]byte(strings.Join(m.characters, ", ")))
	if m.permissions.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		charList.WriteString("\n\n==== Puppet List ====\n")
		charList.WriteString("(Enter the name of the puppet you wish to play)\n")
		charList.WriteString("-----------------------------\n  ")
		charList.Write([]byte(strings.Join(m.puppets, ", ")))
	}
	return charList.String()

}
