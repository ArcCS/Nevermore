package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
	"strings"
)

func init() {
	addHandler(edit{},"Usage:  edit (room|mob|item|exit) (Name) (commands) \n \n Use this command and the sub commands to make modifications to world objects: \n" +
		"Rooms: you must be in the room you wish to edit name not required\n" +
		"Exits:  it must be already created in the room you are standing in \n" +
		"Items: You must edit the item in your inventory.\n" +
		"  -->  If you wish to save it as the template for that item, use the 'savetemplate item' command\n" +
		"Mob: You must summon the mob into the room you are in to edit it. \n" +
		"  -->  If you wish to save it as the template for that item, use the 'savetemplate mob' command\n\n" +
		"Change value:   edit room description Here is a new description \n\n" +
		"Toggle flag(s):   edit exit north toggle unpickable closed hidden\n",
		permissions.Builder,
		"edit", "modify")
}

type edit cmd

func (edit) process(s *state) {
	// Check arguements
	if len(s.words) < 3 {
		s.msg.Actor.SendInfo("Edit what, how?")
		return
	}

	log.Println("Trying to edit: " + strings.ToLower(s.words[0]))
	switch strings.ToLower(s.words[0]) {
	// Handle Rooms
	case "room":
		// Toggle Flags
		if strings.ToLower(s.words[1]) == "toggle" {
			for _, flag := range s.input[2:] {
				if (s.actor.Permission.HasFlags(permissions.Builder, permissions.Dungeonmaster)) || flag != "active" {
					if s.where.ToggleFlag(strings.ToLower(flag)) {
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}
			}

			// Set a variable
		} else {
			switch strings.ToLower(s.words[1]) {
			case "description":
				s.where.Description = strings.Join(s.input[2:], " ")
				s.msg.Actor.SendGood("Description changed.")
			case "name":
				s.where.Name = strings.Join(s.input[2:], " ")
				s.msg.Actor.SendGood("Name changed.")
			default:
				s.msg.Actor.SendBad("Property not found.")
			}
		}
		s.where.Save()
		return

	// Handle Exits
	case "exit":
		// Toggle Flags
		exitName := s.input[2]
		log.Println("Attempting to edit ", exitName)
		//if len(s.words) > 0 {
		//	exitName = strings.Join(s.input[1:], " ")
		//}
		objectRef := strings.ToLower(exitName)
		if !utils.StringIn(strings.ToUpper(objectRef), directionals) {
			for txtE, _ := range s.where.Exits {
				if strings.Contains(txtE, objectRef) {
					objectRef = txtE
				}
			}
		}
		if exit, exists := s.where.Exits[objectRef]; exists {
			if strings.ToLower(s.input[1]) == "toggle" {
				for _, flag := range s.input[3:] {
					if exit.ToggleFlag(strings.ToLower(flag)) {
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}

			// Set a variable
			} else {
				switch strings.ToLower(s.input[1]) {
				case "description":
					exit.Description = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Description changed.")
				case "name":
					oldName := exit.Name
					exit.Name = strings.Join(s.input[3:], " ")
					s.where.Exits[strings.ToLower(strings.Join(s.input[3:], " "))] = exit
					delete(s.where.Exits, oldName)
					data.RenameExit(exit.Name, oldName, exit.ParentId, exit.ToId)
					s.msg.Actor.SendGood("Name changed.")
				case "key_id":
					intKey, _ :=  strconv.Atoi(s.words[3])
					exit.KeyId = int64(intKey)
					s.msg.Actor.SendGood("Change Key Id")
				case "placement":
					intKey, _ :=  strconv.Atoi(s.words[3])
					exit.KeyId = int64(intKey)
					s.msg.Actor.SendGood("Changed placement")
				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
			exit.Save()
		} else {
			s.msg.Actor.SendBad("Exit not found.")
		}

		return
	default:
		s.msg.Actor.SendBad("Not an object that can be edited, or WIP")
	}

	s.ok = true
	return
}
