package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/stats"
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

	//log.Println("Trying to edit: " + strings.ToLower(s.words[0]))
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
		objectRef := strings.ToLower(exitName)
		if !utils.StringIn(strings.ToUpper(objectRef), directionals) {
			for txtE := range s.where.Exits {
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
					exit.KeyId = intKey
					s.msg.Actor.SendGood("Change Key Id")
				case "placement":
					intKey, _ :=  strconv.Atoi(s.words[3])
					if intKey >= 1 && intKey <= 5 {
						exit.KeyId = intKey
						s.msg.Actor.SendGood("Changed placement")
					}else{
						s.msg.Actor.SendBad("Placement Id not valid. ")
					}
				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
			exit.Save()
		} else {
			s.msg.Actor.SendBad("Exit not found.")
		}

		return

	// Handle Items
	case "item":
		// Toggle Flags
		itemName := s.input[2]
		item := s.actor.Inventory.Search(itemName, 1)

		if item != nil {
			if strings.ToLower(s.input[1]) == "toggle" {
				for _, flag := range s.input[3:] {
					if item.ToggleFlag(strings.ToLower(flag)) {
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}

				// Set a variable
			} else {
				switch strings.ToLower(s.input[1]) {
				case "description":
					item.Description = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Description changed.")
				case "name":
					item.Name = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Name changed.")
				case "spell":
					item.Name = s.input[3]
					s.msg.Actor.SendGood("Spell changed.")
				case "weight":
					weight, _ :=  strconv.Atoi(s.words[3])
					item.Weight = weight
					s.msg.Actor.SendGood("Change weight")
				case "type":
					types, err :=  strconv.Atoi(s.words[3])
					if err != nil {
						s.msg.Actor.SendBad("Type must be an integer, use command 'types' to print types.")
						return
					}
					item.ItemType = types
					s.msg.Actor.SendGood("Changed types.")
				case "value":
					value, _ :=  strconv.Atoi(s.words[3])
					item.Value = value
					s.msg.Actor.SendGood("Changed value")
				case "ndice":
					value, _ :=  strconv.Atoi(s.words[3])
					item.NumDice = value
					s.msg.Actor.SendGood("Changed number of dice")
				case "armor":
					value, _ :=  strconv.Atoi(s.words[3])
					item.Armor = value
					s.msg.Actor.SendGood("Changed armor value")
				case "pdice":
					value, _ :=  strconv.Atoi(s.words[3])
					item.PlusDice = value
					s.msg.Actor.SendGood("Changed plus dice")
				case "sdice":
					value, _ :=  strconv.Atoi(s.words[3])
					item.SidesDice	 = value
					s.msg.Actor.SendGood("Changed sides of dice")
				case "max_uses":
					value, _ :=  strconv.Atoi(s.words[3])
					item.MaxUses = value
					s.msg.Actor.SendGood("Changed max_uses")
				/*case "placement":
					intKey, _ :=  strconv.Atoi(s.words[3])
					if intKey >= 1 && intKey <= 5 {
						exit.KeyId = intKey
						s.msg.Actor.SendGood("Changed placement")
					}else{
						s.msg.Actor.SendBad("Placement Id not valid. ")
					}*/
				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
			item.Save()
		} else {
			s.msg.Actor.SendBad("Exit not found.")
		}

		return

	// Handle Mobs
	case "mob":
		// Toggle Flags
		mobName := s.input[2]
		mob := s.where.Mobs.Search(mobName, 1, true)

		if mob != nil {
			if strings.ToLower(s.input[1]) == "toggle" {
				for _, flag := range s.input[3:] {
					if mob.ToggleFlag(strings.ToLower(flag)) {
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}

				// Set a variable
			} else {
				switch strings.ToLower(s.input[1]) {
				case "description":
					mob.Description = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Description changed.")
				case "name":
					mob.Name = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Name changed.")
				case "level":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.Level = value
					s.msg.Actor.SendGood("Change level")
				case "experience":
					types, _ :=  strconv.Atoi(s.words[3])
					mob.Experience = types
					s.msg.Actor.SendGood("Changed experience value.")
				case "gold":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.Gold = value
					s.msg.Actor.SendGood("Changed amount of gold dropped.")
				case "con":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.Con.Current = value
					s.msg.Actor.SendGood("Changed constitution")
				case "int":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.Int.Current = value
					s.msg.Actor.SendGood("Changed intelligence")
				case "str":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.Str.Current = value
					s.msg.Actor.SendGood("Changed strength")
				case "dex":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.Dex.Current = value
					s.msg.Actor.SendGood("Changed dexterity")
				case "pie":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.Pie.Current = value
					s.msg.Actor.SendGood("Changed piety")
				case "mana":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.Mana.Max = value
					s.msg.Actor.SendGood("Changed mana")
				case "stam":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.Mana.Max = value
					s.msg.Actor.SendGood("Changed stam")
				case "ndice":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.NumDice = value
					s.msg.Actor.SendGood("Changed number of dice")
				case "armor":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.Armor = value
					s.msg.Actor.SendGood("Changed armor value")
				case "pdice":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.PlusDice = value
					s.msg.Actor.SendGood("Changed plus dice")
				case "sdice":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.SidesDice	 = value
					s.msg.Actor.SendGood("Changed sides of dice")
				case "chancecast":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.ChanceCast = value
					s.msg.Actor.SendGood("Changed chance to cast")
				case "numwander":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.NumWander = value
					s.msg.Actor.SendGood("Changed amount of ticks to wander")
				case "wimpyvalue":
					value, _ :=  strconv.Atoi(s.words[3])
					mob.WimpyValue = value
					s.msg.Actor.SendGood("Changed value that mob tries to flee")
				/*case "placement":
				intKey, _ :=  strconv.Atoi(s.words[3])
				if intKey >= 1 && intKey <= 5 {
					exit.KeyId = intKey
					s.msg.Actor.SendGood("Changed placement")
				}else{
					s.msg.Actor.SendBad("Placement Id not valid. ")
				}*/
				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
			mob.Save()
		} else {
			s.msg.Actor.SendBad("Exit not found.")
		}

		return
		// Handle Mobs
	case "char":
		// Toggle Flags
		charName := s.input[2]

		stats.ActiveCharacters.Lock()
		character := stats.ActiveCharacters.Find(charName)
		
		if character != nil {
			if strings.ToLower(s.input[1]) == "toggle" {
				for _, flag := range s.input[3:] {
					if character.ToggleFlag(strings.ToLower(flag)) {
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}

				// Set a variable
			} else {
				switch strings.ToLower(s.input[1]) {
				case "description":
					character.Description = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Description changed.")
				case "name":
					character.Name = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Name changed.")
				case "level":
					value, _ :=  strconv.Atoi(s.words[3])
					character.Tier	 = value
					s.msg.Actor.SendGood("Changed Tier")
				case "experience":
					types, _ :=  strconv.Atoi(s.words[3])
					character.Experience = types
					s.msg.Actor.SendGood("Changed amount of experience.")
				case "gold":
					value, _ :=  strconv.Atoi(s.words[3])
					character.Gold.Value = value
					s.msg.Actor.SendGood("Changed amount of gold on character")
				case "bank gold":
					value, _ :=  strconv.Atoi(s.words[3])
					character.BankGold.Value = value
					s.msg.Actor.SendGood("Changed amount of gold in bank.")
				case "con":
					value, _ :=  strconv.Atoi(s.words[3])
					character.Con.Current = value
					s.msg.Actor.SendGood("Changed constitution")
				case "int":
					value, _ :=  strconv.Atoi(s.words[3])
					character.Int.Current = value
					s.msg.Actor.SendGood("Changed intelligence")
				case "str":
					value, _ :=  strconv.Atoi(s.words[3])
					character.Str.Current = value
					s.msg.Actor.SendGood("Changed strength")
				case "dex":
					value, _ :=  strconv.Atoi(s.words[3])
					character.Dex.Current = value
					s.msg.Actor.SendGood("Changed dexterity")
				case "pie":
					value, _ :=  strconv.Atoi(s.words[3])
					character.Pie.Current = value
					s.msg.Actor.SendGood("Changed piety")
				case "mana":
					value, _ :=  strconv.Atoi(s.words[3])
					character.Mana.Max = value
					s.msg.Actor.SendGood("Changed mana")
				case "stam":
					value, _ :=  strconv.Atoi(s.words[3])
					character.Mana.Max = value
					s.msg.Actor.SendGood("Changed stam")
				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
			mob.Save()
		} else {
			s.msg.Actor.SendBad("Exit not found.")
		}
		stats.ActiveCharacters.Unlock()
		return
	default:
		s.msg.Actor.SendBad("Not an object that can be edited, or WIP")
	}

	s.ok = true
	return
}
