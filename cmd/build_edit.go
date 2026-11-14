package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"math"
	"strconv"
	"strings"
)

func init() {
	addHandler(edit{}, "Usage:  edit (room|mob|item|exit|char) (Property|Toggle) (SubjectName) (Values) \n \n Use this command and the sub commands to make modifications to world objects: \n"+
		"Rooms: you must be in the room you wish to edit name not required\n"+
		"Exits:  it must be already created in the room you are standing in \n"+
		"Items: You must edit the item in your inventory.\n"+
		"Change value:   edit room description Here is a new description \n\n"+
		"Toggle flag(s):   edit exit toggle north unpickable closed hidden\n\n"+
		"Leaving out the value empties the field, where the field allows it.\n"+
		"Item and mob edits change the template: only the edited field is saved, whatever state your copy is in.\n",
		permissions.Builder,
		"edit", "modify")
}

type edit cmd

func (edit) process(s *state) {
	// Check arguments
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Edit what, how?")
		return
	}

	target := strings.ToLower(s.words[0])
	property := strings.ToLower(s.words[1])
	// Rooms edit the room you're in, everything else names its subject first.
	valueIndex := 3
	if target == "room" {
		valueIndex = 2
	} else if len(s.words) < 3 {
		s.msg.Actor.SendInfo("Edit what, how?")
		return
	}

	// A value left out edits the field to empty, unless it can't be.
	wordArg := editArg(s.words, valueIndex)
	if wordArg == "" && editRequiresValue(target, property) {
		s.msg.Actor.SendBad(property + " cannot be empty.")
		return
	}
	if _, err := strconv.Atoi(wordArg); wordArg != "" && err != nil && editIntProperty(target, property) {
		s.msg.Actor.SendBad(property + " must be a whole number.")
		return
	}

	switch target {
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
				s.msg.Actor.SendGood("Room name changed.")
			case "storeowner":
				s.where.StoreOwner = editArg(s.input, 2)
				s.msg.Actor.SendGood("Store Owner changed.")
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
					intKey, _ := strconv.Atoi(wordArg)
					exit.KeyId = intKey
					s.msg.Actor.SendGood("Change Key Id")
				case "placement":
					intKey, _ := strconv.Atoi(wordArg)
					if intKey >= 1 && intKey <= 5 {
						exit.Placement = intKey
						s.msg.Actor.SendGood("Changed placement")
					} else {
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
			// The copy in hand has worn down since it spawned, so it is only
			// ever given the edit: what gets saved starts from the database.
			template, ok := loadItemTemplate(item.ItemId)
			if !ok {
				s.msg.Actor.SendBad("That item has no template to edit.")
				return
			}
			if property == "toggle" {
				for _, flag := range s.input[3:] {
					flagName := strings.ToLower(flag)
					if template.ToggleFlag(flagName) {
						item.Flags[flagName] = template.Flags[flagName]
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}

				// Set a variable
			} else {
				result, ok := applyItemEdit(template, property, s.input[3:])
				if !ok {
					s.msg.Actor.SendBad(result)
					return
				}
				applyItemEdit(item, property, s.input[3:])
				s.msg.Actor.SendGood(result)
			}
			template.Save()
			objects.Items[item.ItemId], _ = objects.LoadItem(data.LoadItem(item.ItemId))
		} else {
			s.msg.Actor.SendBad("Item not found.")
		}

		return

	// Handle Mobs
	case "mob":
		// Toggle Flags
		mobName := s.input[2]
		mob := s.where.Mobs.Search(mobName, 1, s.actor)

		if mob != nil {
			// Same as items: a live mob's vitals and stats aren't its template's.
			template, ok := loadMobTemplate(mob.MobId)
			if !ok {
				s.msg.Actor.SendBad("That mob has no template to edit.")
				return
			}
			if property == "toggle" {
				for _, flag := range s.input[3:] {
					flagName := strings.ToLower(flag)
					if template.EditToggleFlag(flagName) {
						mob.Flags[flagName] = template.Flags[flagName]
						s.msg.Actor.SendGood("Toggled " + flag)
						if flag == "permanent" {
							log.Println("Executing Save")
							s.where.Save()
						}
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}
				// Set a variable
			} else {
				result, ok := applyMobEdit(template, property, s.input[3:])
				if !ok {
					s.msg.Actor.SendBad(result)
					return
				}
				applyMobEdit(mob, property, s.input[3:])
				if property == "spells" {
					// Spells toggle, so the live mob follows the template's list.
					mob.Spells = append([]string{}, template.Spells...)
				}
				s.msg.Actor.SendGood(result)
			}
			template.Save()
			objects.Mobs[mob.MobId], _ = objects.LoadMob(data.LoadMob(mob.MobId))
		} else {
			s.msg.Actor.SendBad("Mob not found.")
		}

		return
		// Handle Mobs
	case "char":
		// Toggle Flags
		charName := s.input[2]

		character := objects.ActiveCharacters.Find(charName)

		if character != nil {
			objects.ActiveCharacters.Lock()
			defer objects.ActiveCharacters.Unlock()
			if strings.ToLower(s.input[1]) == "toggle" {
				for _, flag := range s.input[3:] {
					character.ToggleFlag(strings.ToLower(flag), "")
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
				case "title":
					character.Title = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Title changed.")
				case "tier":
					value, _ := strconv.Atoi(wordArg)
					character.Tier = value
					character.RecalcMaxes()
					s.msg.Actor.SendGood("Changed Tier")
				case "class":
					value, _ := strconv.Atoi(wordArg)
					character.Class = value
					character.ClassTitle = config.ClassTitle(value, character.Gender, character.Tier)
					s.msg.Actor.SendGood("Changed Class")
				case "race":
					value, _ := strconv.Atoi(wordArg)
					character.Race = value
					s.msg.Actor.SendGood("Changed Race")
				case "experience":
					types, _ := strconv.Atoi(wordArg)
					character.Experience.Value = types
					s.msg.Actor.SendGood("Changed amount of experience.")
				case "gold":
					value, _ := strconv.Atoi(wordArg)
					character.Gold.Value = value
					s.msg.Actor.SendGood("Changed amount of gold on character")
				case "bankgold":
					value, _ := strconv.Atoi(wordArg)
					character.BankGold.Value = value
					s.msg.Actor.SendGood("Changed amount of gold in bank.")
				case "passages":
					value, _ := strconv.Atoi(wordArg)
					character.Passages.Value = value
					s.msg.Actor.SendGood("Changed number of passages")
				case "bonuspoints":
					value, _ := strconv.Atoi(wordArg)
					character.BonusPoints.Value = value
					s.msg.Actor.SendGood("Changed number of bonus points")
				case "broadcasts":
					value, _ := strconv.Atoi(wordArg)
					character.Broadcasts = value
					s.msg.Actor.SendGood("Changed broadcasts.")
				case "evals":
					value, _ := strconv.Atoi(wordArg)
					character.Evals = value
					s.msg.Actor.SendGood("Changed evals")
				case "rerolls":
					value, _ := strconv.Atoi(wordArg)
					character.Rerolls = value
					s.msg.Actor.SendGood("Changed rerolls")
				case "concur":
					value, _ := strconv.Atoi(wordArg)
					character.Con.Current = value
					s.msg.Actor.SendGood("Changed constitution")
				case "intcur":
					value, _ := strconv.Atoi(wordArg)
					character.Int.Current = value
					s.msg.Actor.SendGood("Changed intelligence")
				case "strcur":
					value, _ := strconv.Atoi(wordArg)
					character.Str.Current = value
					s.msg.Actor.SendGood("Changed strength")
				case "dexcur":
					value, _ := strconv.Atoi(wordArg)
					character.Dex.Current = value
					s.msg.Actor.SendGood("Changed dexterity")
				case "piecur":
					value, _ := strconv.Atoi(wordArg)
					character.Pie.Current = value
					s.msg.Actor.SendGood("Changed piety")
				case "stamcur":
					value, _ := strconv.Atoi(wordArg)
					character.Stam.Add(value)
					s.msg.Actor.SendGood("Changed current stamina")
				case "stammax":
					value, _ := strconv.Atoi(wordArg)
					character.Stam.Max = value
					s.msg.Actor.SendGood("Changed stamina maximum")
				case "manamax":
					value, _ := strconv.Atoi(wordArg)
					character.Mana.Max = value
					s.msg.Actor.SendGood("Changed mana maximum")
				case "manacur":
					value, _ := strconv.Atoi(wordArg)
					character.Mana.Add(value)
					s.msg.Actor.SendGood("Changed mana current")
				case "vitmax":
					value, _ := strconv.Atoi(wordArg)
					character.Vit.Max = value
					s.msg.Actor.SendGood("Changed vitality max")
				case "vitcur":
					value, _ := strconv.Atoi(wordArg)
					character.Vit.Add(value)
					s.msg.Actor.SendGood("Changed vit current")
				case "sharpexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[0].Value = value
					s.msg.Actor.SendGood("Changed sharp exp")
				case "thrustexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[1].Value = value
					s.msg.Actor.SendGood("Changed thrust exp")
				case "bluntexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[2].Value = value
					s.msg.Actor.SendGood("Changed blunt exp")
				case "poleexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[3].Value = value
					s.msg.Actor.SendGood("Changed pole exp")
				case "missileexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[4].Value = value
					s.msg.Actor.SendGood("Changed missile exp")
				case "handexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[5].Value = value
					s.msg.Actor.SendGood("Changed hand to hand exp")
				case "fireexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[6].Value = value
					s.msg.Actor.SendGood("Changed fire affinity")
				case "airexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[7].Value = value
					s.msg.Actor.SendGood("Changed air affinity")
				case "earthexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[8].Value = value
					s.msg.Actor.SendGood("Changed earth affinity")
				case "waterexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[9].Value = value
					s.msg.Actor.SendGood("Changed water affinity")
				case "divinity":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[10].Value = value
					s.msg.Actor.SendGood("Changed divinity")
				case "stealthexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[11].Value = value
					s.msg.Actor.SendGood("Changed stealth exp")

				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
			character.Save()
		} else {
			if strings.ToLower(s.input[1]) == "toggle" {
				s.msg.Actor.SendBad("Cannot toggle offline character.")
			} else {
				switch strings.ToLower(s.input[1]) {
				case "description":
					data.SaveCharField(charName, "description", strings.Join(s.input[3:], " "))
					s.msg.Actor.SendGood("Description changed.")
				case "name":
					data.SaveCharField(charName, "name", strings.Join(s.input[3:], " "))
					s.msg.Actor.SendGood("Name changed.")
				case "title":
					data.SaveCharField(charName, "title", strings.Join(s.input[3:], " "))
					s.msg.Actor.SendGood("Title changed.")
				case "tier":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "tier", value)
					s.msg.Actor.SendGood("Changed Tier")
				case "experience":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "experience", value)
					s.msg.Actor.SendGood("Changed amount of experience.")
				case "gold":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "gold", value)
					s.msg.Actor.SendGood("Changed amount of gold on character")
				case "bankgold":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "bankgold", value)
					s.msg.Actor.SendGood("Changed amount of gold in bank.")
				case "passages":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "passages", value)
					s.msg.Actor.SendGood("Changed number of passages")
				case "bonuspoints":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "bonuspoints", value)
					s.msg.Actor.SendGood("Changed number of bonus points")
				case "broadcasts":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "broadcasts", value)
					s.msg.Actor.SendGood("Changed broadcasts.")
				case "evals":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "evals", value)
					s.msg.Actor.SendGood("Changed evals")
				case "rerolls":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "rerolls", value)
					s.msg.Actor.SendGood("Changed rerolls")
				case "concur":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "concur", value)
					s.msg.Actor.SendGood("Changed constitution")
				case "intcur":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "intcur", value)
					s.msg.Actor.SendGood("Changed intelligence")
				case "strcur":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "strcur", value)
					s.msg.Actor.SendGood("Changed strength")
				case "dexcur":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "dexcur", value)
					s.msg.Actor.SendGood("Changed dexterity")
				case "piecur":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "piecur", value)
					s.msg.Actor.SendGood("Changed piety")
				case "stamcur":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "curr_stam", value)
					s.msg.Actor.SendGood("Changed current stamina")
				case "stammax":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "max_stam", value)
					s.msg.Actor.SendGood("Changed stamina maximum")
				case "manamax":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "max_mana", value)
					s.msg.Actor.SendGood("Changed mana maximum")
				case "manacur":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "curr_Mana", value)
					s.msg.Actor.SendGood("Changed mana current")
				case "vitmax":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "max_vit", value)
					s.msg.Actor.SendGood("Changed vitality max")
				case "vitcur":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "curr_vit", value)
					s.msg.Actor.SendGood("Changed vit current")
				case "sharpexp":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "sharpexp", value)
					s.msg.Actor.SendGood("Changed sharp exp")
				case "thrustexp":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "thrustexp", value)
					s.msg.Actor.SendGood("Changed thrust exp")
				case "bluntexp":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "bluntexp", value)
					s.msg.Actor.SendGood("Changed blunt exp")
				case "poleexp":
					value, _ := strconv.Atoi(wordArg)
					data.SaveCharField(charName, "poleexp", value)
					s.msg.Actor.SendGood("Changed pole exp")
				case "missileexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[4].Value = value
					s.msg.Actor.SendGood("Changed missile exp")
				case "handexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[5].Value = value
					s.msg.Actor.SendGood("Changed hand to hand exp")
				case "fireexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[6].Value = value
					s.msg.Actor.SendGood("Changed fire affinity")
				case "airexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[7].Value = value
					s.msg.Actor.SendGood("Changed air affinity")
				case "earthexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[8].Value = value
					s.msg.Actor.SendGood("Changed earth affinity")
				case "waterexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[9].Value = value
					s.msg.Actor.SendGood("Changed water affinity")
				case "divinity":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[10].Value = value
					s.msg.Actor.SendGood("Changed divinity")
				case "stealthexp":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[11].Value = value
					s.msg.Actor.SendGood("Changed stealth exp")
				case "parent_id":
					value, _ := strconv.Atoi(wordArg)
					character.Skills[4].Value = value
					s.msg.Actor.SendGood("Changed room")
				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
		}

		return
	default:
		s.msg.Actor.SendBad("Not an object that can be edited.")
	}

	s.ok = true
	return
}

// editArg is the argument at index, empty when it was left out.
func editArg(args []string, index int) string {
	if index < len(args) {
		return args[index]
	}
	return ""
}

// editRequiresValue reports whether a property can't be left empty.  A player's
// stats are all required, a typo there would zero what they earned.
func editRequiresValue(target string, property string) bool {
	switch target {
	case "room":
		return property == "name"
	case "exit":
		return utils.StringIn(property, []string{"name", "placement"})
	case "item":
		return utils.StringIn(property, []string{"name", "type"})
	case "mob":
		return utils.StringIn(property, []string{"name", "level", "stam", "spells", "placement", "exp_multiplier"})
	case "char":
		return !utils.StringIn(property, []string{"toggle", "description", "title"})
	}
	return false
}

// editIntProperty reports whether a property is stored as an integer, the
// int64 its loader reads, so anything else typed in is refused, not saved as 0.
func editIntProperty(target string, property string) bool {
	switch target {
	case "exit":
		return utils.StringIn(property, []string{"key_id", "placement"})
	case "item":
		return utils.StringIn(property, []string{"weight", "type", "value", "ndice", "armor", "pdice",
			"adjustment", "sdice", "max_uses", "armor_class"})
	case "mob":
		return utils.StringIn(property, []string{"level", "gold", "con", "int", "str", "dex", "pie", "mana",
			"stam", "ndice", "armor", "pdice", "sdice", "chancecast", "numwander", "wimpyvalue",
			"air_resistance", "earth_resistance", "fire_resistance", "water_resistance", "placement"})
	case "char":
		return !utils.StringIn(property, []string{"toggle", "description", "name", "title"})
	}
	return false
}

// loadItemTemplate reads an item's template fresh from the database.
func loadItemTemplate(itemId int) (*objects.Item, bool) {
	if _, ok := objects.Items[itemId]; !ok {
		return nil, false
	}
	itemData := data.LoadItem(itemId)
	if itemData == nil {
		return nil, false
	}
	return objects.LoadItem(itemData)
}

// loadMobTemplate reads a mob's template fresh from the database.
func loadMobTemplate(mobId int) (*objects.Mob, bool) {
	if _, ok := objects.Mobs[mobId]; !ok {
		return nil, false
	}
	mobData := data.LoadMob(mobId)
	if mobData == nil {
		return nil, false
	}
	return objects.LoadMob(mobData)
}

// applyItemEdit sets one property from args, the words after the item's name.
// It returns the message for the builder and whether the edit was accepted.
func applyItemEdit(item *objects.Item, property string, args []string) (string, bool) {
	arg := editArg(args, 0)
	value, _ := strconv.Atoi(arg)
	switch property {
	case "description":
		item.Description = strings.Join(args, " ")
		return "Description changed.", true
	case "name":
		item.Name = strings.Join(args, " ")
		return "Name changed.", true
	case "spell":
		if arg != "" {
			_, isSpell := objects.Spells[arg]
			_, isSong := objects.Songs[arg]
			if !isSpell && !isSong {
				return "Spell not found.", false
			}
		}
		item.Spell = arg
		return "Spell changed.", true
	case "weight":
		item.Weight = value
		return "Change weight", true
	case "type":
		types, err := strconv.Atoi(arg)
		if err != nil {
			return "Type must be an integer, use command 'types' to print types.", false
		}
		item.ItemType = types
		return "Changed item type to " + config.ItemTypes[types], true
	case "value":
		item.Value = value
		return "Changed value", true
	case "ndice":
		item.NumDice = value
		return "Changed number of dice", true
	case "armor":
		item.Armor = value
		return "Changed armor value", true
	case "pdice":
		item.PlusDice = value
		return "Changed plus dice", true
	case "adjustment":
		item.Adjustment = value
		return "Changed adjustment to damage", true
	case "sdice":
		item.SidesDice = value
		return "Changed sides of dice", true
	case "max_uses":
		item.MaxUses = value
		return "Changed max_uses", true
	case "armor_class":
		item.ArmorClass = value
		return "Changed armor class", true
	case "charge_type":
		chargeType := config.ChargeTypeDefault
		if arg != "" {
			parsed, ok := config.ParseChargeType(arg)
			if !ok {
				return "Charge type must be default, liquid, food or plural.", false
			}
			chargeType = parsed
		}
		item.ChargeType = chargeType
		return "Changed charge type to " + config.ChargeTypes[chargeType], true
	}
	return "Property not found.", false
}

// applyMobEdit sets one property from args, the words after the mob's name.
// It returns the message for the builder and whether the edit was accepted.
func applyMobEdit(mob *objects.Mob, property string, args []string) (string, bool) {
	arg := editArg(args, 0)
	value, _ := strconv.Atoi(arg)
	switch property {
	case "description":
		mob.Description = strings.Join(args, " ")
		return "Description changed.", true
	case "name":
		mob.Name = strings.Join(args, " ")
		return "Name changed.", true
	case "level":
		mob.Level = value
		return "Change level", true
	case "exp_multiplier":
		multiplier, err := strconv.ParseFloat(arg, 64)
		if err != nil || multiplier < 0 {
			return "The experience multiplier must be a number of 0 or more.", false
		}
		mob.ExpMultiplier = math.Round(multiplier*100) / 100
		return "Changed experience multiplier; it is now worth " + strconv.Itoa(mob.Experience()) + " experience.", true
	case "gold":
		mob.Gold = value
		return "Changed amount of gold dropped.", true
	case "con":
		mob.Con.Current = value
		return "Changed constitution", true
	case "int":
		mob.Int.Current = value
		return "Changed intelligence", true
	case "str":
		mob.Str.Current = value
		return "Changed strength", true
	case "dex":
		mob.Dex.Current = value
		return "Changed dexterity", true
	case "pie":
		mob.Pie.Current = value
		return "Changed piety", true
	case "mana":
		mob.Mana.Max = value
		mob.Mana.Current = value
		return "Changed mana", true
	case "stam":
		mob.Stam.Max = value
		mob.Stam.Current = value
		return "Changed stam", true
	case "ndice":
		mob.NumDice = value
		return "Changed number of dice", true
	case "armor":
		mob.Armor = value
		return "Changed armor value", true
	case "pdice":
		mob.PlusDice = value
		return "Changed plus dice", true
	case "sdice":
		mob.SidesDice = value
		return "Changed sides of dice", true
	case "chancecast":
		mob.ChanceCast = value
		return "Changed chance to cast", true
	case "numwander":
		mob.NumWander = value
		return "Changed amount of ticks to wander", true
	case "wimpyvalue":
		mob.WimpyValue = value
		return "Changed value that mob tries to flee", true
	case "air_resistance":
		mob.AirResistance = value
		return "Changed air resistance", true
	case "earth_resistance":
		mob.EarthResistance = value
		return "Changed earth resistance", true
	case "fire_resistance":
		mob.FireResistance = value
		return "Changed fire resistance", true
	case "water_resistance":
		mob.WaterResistance = value
		return "Changed water resistance", true
	case "spells":
		spellName := strings.ToLower(arg)
		if _, ok := objects.Spells[spellName]; !ok {
			return spellName + " not found in spell definitions.", false
		}
		if utils.StringIn(spellName, mob.Spells) {
			// Deleting
			mob.Spells[utils.IndexOf(spellName, mob.Spells)] = mob.Spells[len(mob.Spells)-1]
			mob.Spells = mob.Spells[:len(mob.Spells)-1]
			return spellName + " deleted from mob spellbook", true
		}
		// Adding
		mob.Spells = append(mob.Spells, spellName)
		return spellName + " add to mob spellbook", true
	case "breathes":
		spellName := strings.ToLower(arg)
		if spellName == "" {
			mob.BreathWeapon = ""
			return "Mob BreathWeapon unset", true
		}
		if !utils.StringIn(spellName, []string{"air", "fire", "water", "earth", "paralytic", "pestilence"}) {
			return "Not a valid BreathWeapon.", false
		}
		mob.BreathWeapon = spellName
		return "Mob BreathWeapon set to " + arg, true
	case "placement":
		placement, err := strconv.Atoi(arg)
		if err != nil {
			return "Placement value not valid", false
		}
		mob.Placement = placement
		return "Changed placement", true
	}
	return "Property not found.", false
}
