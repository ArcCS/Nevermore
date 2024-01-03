package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"log"
	"strconv"
	"strings"
)

func init() {
	addHandler(examine{},
		"Usage:  examine (room|mob|item|exit|char) (name|#####) \n\n  Examine will display the subject and all of it's modifiable properties",
		permissions.Builder,
		"examine")
}

type examine cmd

func (examine) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What do you want to examine?")
		return
	}

	rowLength := 120
	typeE := strings.ToLower(s.words[0])

	nameNum := 1

	switch typeE {
	case "room":
		roomRef := s.where
		if len(s.words) > 1 {
			roomNumber, ok := strconv.Atoi(s.words[1])
			if ok == nil {
				roomRef = objects.Rooms[roomNumber]
			} else {
				roomRef = s.where
			}
		}

		t := table.NewWriter()
		t.SetAllowedRowLength(rowLength)
		t.Style().Options.SeparateRows = true
		t.AppendHeader(table.Row{"Type", "Variable Name", "Value", "Description"})
		t.AppendRows([]table.Row{
			{"V", "name", roomRef.Name, "Title/Name of the room"},
			{"V", "storeowner", roomRef.StoreOwner, "Owner that can add items."},
			{"T", "repair", strconv.FormatBool(roomRef.Flags["repair"]), "Can repair items here."},
			{"T", "mana_drain", strconv.FormatBool(roomRef.Flags["mana_drain"]), "Drains mana per tick"},
			{"T", "no_summon", strconv.FormatBool(roomRef.Flags["no_summon"]), "Blocks summon"},
			{"T", "heal_fast", strconv.FormatBool(roomRef.Flags["heal_fast"]), "Regenerate quickly"},
			{"T", "no_teleport", strconv.FormatBool(roomRef.Flags["no_teleport"]), "Blocks teleporting out"},
			{"T", "no_word_of_recall", strconv.FormatBool(roomRef.Flags["no_word_of_recall"]), "Blocks WoD"},
			{"T", "no_magic", strconv.FormatBool(roomRef.Flags["no_magic"]), "Cannot cast spells"},
			{"T", "dark_always", strconv.FormatBool(roomRef.Flags["dark_always"]), "Always dark"},
			{"T", "light_always", strconv.FormatBool(roomRef.Flags["light_always"]), "Always lit"},
			{"T", "natural_light", strconv.FormatBool(roomRef.Flags["natural_light"]), "Follows day/night cycle"},
			{"T", "indoors", strconv.FormatBool(roomRef.Flags["indoors"]), "Room is indoors"},
			{"T", "urban", strconv.FormatBool(roomRef.Flags["urban"]), "Room is part of a city"},
			{"T", "underground", strconv.FormatBool(roomRef.Flags["underground"]), "Room is underground"},
			{"T", "encounters_on", strconv.FormatBool(roomRef.Flags["encounters_on"]), "Monsters can spawn"},
			{"T", "fire", strconv.FormatBool(roomRef.Flags["fire"]), "Room causes and amplifies fire."},
			{"T", "water", strconv.FormatBool(roomRef.Flags["water"]), "Room causes and amplifies water."},
			{"T", "earth", strconv.FormatBool(roomRef.Flags["earth"]), "Room causes and amplifies earth."},
			{"T", "wind", strconv.FormatBool(roomRef.Flags["wind"]), "Room causes and amplifies wind."},
			{"T", "active", strconv.FormatBool(roomRef.Flags["active"]), "The room is activated."},
			{"T", "train", strconv.FormatBool(roomRef.Flags["train"]), "Characters can train here."},
		})
		t.SetCaption("Light is evaluated with dark_always first, then light_always, then natural light.\nX = Cannot Modify,  T=Toggle to Edit, V=Edit by value name\nSee 'help edit' for more.")
		s.msg.Actor.SendGood(t.Render())

	case "mob":
		if len(s.words) > 2 {
			// Try to snag a number off the list
			if val, err := strconv.Atoi(s.words[2]); err == nil {
				nameNum = val
			}
		}
		mobRef := s.where.Mobs.Search(s.words[1], nameNum, s.actor)
		if mobRef != nil {

			t := table.NewWriter()
			t.SetAllowedRowLength(rowLength)
			t.Style().Options.SeparateRows = true
			t.AppendHeader(table.Row{"Type", "Variable Name", "Value", "Description"})
			t.AppendRows([]table.Row{
				{"X", "mobId", strconv.Itoa(mobRef.MobId), "DB id for mob"},
				{"V", "name", mobRef.Name, "Mob Name"},
				{"V", "description", text.WrapSoft(mobRef.Description, rowLength/5), "Mob description"},
				{"V", "level", strconv.Itoa(mobRef.Level), "Mob Level"},
				{"V", "experience", strconv.Itoa(mobRef.Experience), "Awarded Experience"},
				{"V", "gold", strconv.Itoa(mobRef.Gold), "Gold dropped"},
				{"V", "con", strconv.Itoa(mobRef.Con.Current), "Amt of Con"},
				{"V", "str", strconv.Itoa(mobRef.Str.Current), "Amt of Str"},
				{"V", "int", strconv.Itoa(mobRef.Int.Current), "Amt of Int"},
				{"V", "dex", strconv.Itoa(mobRef.Dex.Current), "Amt of Dex"},
				{"V", "pie", strconv.Itoa(mobRef.Pie.Current), "Amt of Pie"},
				{"V", "mana", strconv.Itoa(mobRef.Mana.Current), "Available Mana"},
				{"V", "stam", strconv.Itoa(mobRef.Stam.Current), "Total HP"},
				{"V", "ndice", strconv.Itoa(mobRef.NumDice), "Number of Dice"},
				{"V", "sdice", strconv.Itoa(mobRef.SidesDice), "Sides of Dice"},
				{"V", "pdice", strconv.Itoa(mobRef.PlusDice), "Plus modifier to dice"},
				{"V", "chancecast", strconv.Itoa(mobRef.ChanceCast), "Chance to cast spell"},
				{"V", "armor", strconv.Itoa(mobRef.Armor), "Amt of Armor"},
				{"V", "numwander", strconv.Itoa(mobRef.NumWander), "Number of ticks to wander"},
				{"V", "wimpyvalue", strconv.Itoa(mobRef.WimpyValue), "Flee @ Damage"},
				{"V", "air_resistance", mobRef.AirResistance, "Resists air damage %."},
				{"V", "earth_resistance", mobRef.EarthResistance, "Resists earth damage %."},
				{"V", "fire_resistance", mobRef.FireResistance, "Resists fire damage %."},
				{"V", "water_resistance", mobRef.WaterResistance, "Resists water damage %."},
				{"V", "breathes", mobRef.BreathWeapon, text.WrapSoft("Element breath, earth air fire water paralytic pestilence", rowLength/5)},
				{"V", "spells", text.WrapSoft(strings.Join(mobRef.Spells, ", "), rowLength/5), "Available spells"},
				{"V", "placement", strconv.Itoa(mobRef.Placement), "Mob spawn location"},
				{"T", "no_specials", strconv.FormatBool(mobRef.Flags["no_specials"]), "Whether the mob crits/doubles or not"},
				{"T", "fast_moving", strconv.FormatBool(mobRef.Flags["fast_moving"]), "Mob is moves quickly"},
				{"T", "guard_treasure", strconv.FormatBool(mobRef.Flags["guard_treasure"]), "Mob guards treasure."},
				{"T", "take_treasure", strconv.FormatBool(mobRef.Flags["take_treasure"]), "Mob takes treasure."},
				{"T", "steals", strconv.FormatBool(mobRef.Flags["steals"]), "Mob will steal from target."},
				{"T", "block_exit", strconv.FormatBool(mobRef.Flags["block_exit"]), "Mob will block exits"},
				{"T", "follows", strconv.FormatBool(mobRef.Flags["follows"]), "Mob will follow exiting target."},
				{"T", "no_steal", strconv.FormatBool(mobRef.Flags["no_steal"]), "Cannot be stolen from."},
				{"T", "detect_invisible", strconv.FormatBool(mobRef.Flags["detect_invisible"]), "Can see past invisibility."},
				{"T", "no_stun", strconv.FormatBool(mobRef.Flags["no_stun"]), "Mob cannot be stunned."},
				{"T", "diseases", strconv.FormatBool(mobRef.Flags["diseases"]), "Mob can disease targets."},
				{"T", "poisons", strconv.FormatBool(mobRef.Flags["poisons"]), "Mob can poison targets."},
				{"T", "spits_acid", strconv.FormatBool(mobRef.Flags["spits_acid"]), "Mob can spit acid."},
				{"T", "ranged_attack", strconv.FormatBool(mobRef.Flags["ranged_attack"]), "Mob can attack from afar."},
				{"T", "flees", strconv.FormatBool(mobRef.Flags["flees"]), "Flee @ Wimpy Value"},
				{"T", "blinds", strconv.FormatBool(mobRef.Flags["blinds"]), "Mob can blind targets."},
				{"T", "hide_encounter", strconv.FormatBool(mobRef.Flags["hide_encounter"]), "Hide when encounter"},
				{"T", "invisible", strconv.FormatBool(mobRef.Flags["invisible"]), "Invisible Mob"},
				{"T", "permanent", strconv.FormatBool(mobRef.Flags["permanent"]), "Does not despawn from room"},
				{"T", "immobile", strconv.FormatBool(mobRef.Flags["immobile"]), "Mob cannot move"},
				{"T", "hostile", strconv.FormatBool(mobRef.Flags["hostile"]), "Mob is hostile"},
				{"T", "undead", strconv.FormatBool(mobRef.Flags["undead"]), "Mob is undead"},
				{"T", "day_only", strconv.FormatBool(mobRef.Flags["day_only"]), "Day Spawn Only."},
				{"T", "night_only", strconv.FormatBool(mobRef.Flags["night_only"]), "Night spawn only"},
			})
			t.SetCaption("X = Cannot Modify,  T=Toggle to Edit, V=Edit by value name\nSee 'help edit' for more.")
			s.msg.Actor.SendGood(t.Render())

		} else {
			s.msg.Actor.SendBad("Couldn't find the object in the current room.")
		}

	case "item":
		nameNum := 1
		if len(s.words) < 2 {
			s.msg.Actor.SendBad("What item do you want to examine?")
			return
		}
		if len(s.words) > 2 {
			// Try to snag a number off the list
			if val, err := strconv.Atoi(s.words[2]); err == nil {
				nameNum = val
			}
		}
		objRef := s.actor.Inventory.Search(s.words[1], nameNum)
		if objRef != nil {
			t := table.NewWriter()
			t.SetAllowedRowLength(rowLength)
			t.Style().Options.SeparateRows = true
			t.AppendHeader(table.Row{"Type", "Variable Name", "Value", "Description"})
			t.AppendRows([]table.Row{
				{"X", "itemId", strconv.Itoa(objRef.ItemId), "Database value for items"},
				{"V", "name", objRef.Name, "Name of the item"},
				{"X", "creator", objRef.Creator, "Last person to modify item."},
				{"V", "description", text.WrapSoft(objRef.Description, rowLength/5), "Look output"},
				{"V", "weight", strconv.Itoa(objRef.Weight), "Inventory weight"},
				{"V", "type", config.ItemTypes[objRef.ItemType], "Type of Item"},
				{"V", "value", strconv.Itoa(objRef.Value), "Pawning Value"},
				{"V", "max_uses", strconv.Itoa(objRef.MaxUses), "Number of uses before breakage"},
				{"V", "armor", strconv.Itoa(objRef.Armor), "Armor value from this item"},
				{"V", "armor_class", config.ArmorClass[objRef.ArmorClass], "Armor class of this item"},
				{"V", "ndice", strconv.Itoa(objRef.NumDice), "Number of Dice to Roll"},
				{"V", "sdice", strconv.Itoa(objRef.SidesDice), "Sides of Dice Being Rolled"},
				{"V", "pdice", strconv.Itoa(objRef.PlusDice), "Additional value to add to roll."},
				{"V", "adjustment", strconv.Itoa(objRef.Adjustment), "Adjustment to final roll damage"},
				{"V", "spell", objRef.Spell, "Spell/Song learned/cast when used."},
				{"T", "always_crit", strconv.FormatBool(objRef.Flags["always_crit"]), "Always criticals when used"},
				{"T", "permanent", strconv.FormatBool(objRef.Flags["permanent"]), "Does not despawn"},
				{"T", "magic", strconv.FormatBool(objRef.Flags["magic"]), "Magical item"},
				{"T", "no_take", strconv.FormatBool(objRef.Flags["no_take"]), "Cannot be picked up."},
				{"T", "light", strconv.FormatBool(objRef.Flags["light"]), "Provides user illumination."},
				{"T", "weightless_chest", strconv.FormatBool(objRef.Flags["weightless_chest"]), "Holds items weightlessly"},
			})
			t.SetCaption("X = Cannot Modify,  T=Toggle to Edit, V=Edit by value name\nSee 'help edit' for more.")
			s.msg.Actor.SendGood(t.Render())

		} else {
			s.msg.Actor.SendBad("Couldn't find the object in the current room.")
		}

	case "exit":
		if len(s.words) < 2 {
			s.msg.Actor.SendBad("What item do you want to examine?")
			return
		}
		exitName := s.input[1]
		objectRef := strings.ToLower(exitName)
		if !utils.StringIn(strings.ToUpper(objectRef), directionals) {
			for txtE := range s.where.Exits {
				if strings.Contains(txtE, objectRef) {
					objectRef = txtE
				}
			}
		}
		if exitRef, ok := s.where.Exits[strings.ToLower(objectRef)]; ok {
			t := table.NewWriter()
			t.SetAllowedRowLength(rowLength)
			t.Style().Options.SeparateRows = true
			t.AppendHeader(table.Row{"Type", "Variable Name", "Value", "Description"})
			t.AppendRows([]table.Row{
				{"V", "name", exitRef.Name, "Name of the exit"},
				{"V", "description", text.WrapSoft(exitRef.Description, rowLength/5), "Peers into next room if empty."},
				{"V", "placement", strconv.Itoa(exitRef.Placement), "Where the exit is in the room. \n5 is front, 1 is back."},
				{"V", "key_id", strconv.Itoa(exitRef.KeyId), "Id of key that can unlock exit."},
				{"T", "closeable", strconv.FormatBool(exitRef.Flags["closeable"]), "The exit can be closed."},
				{"T", "closed", strconv.FormatBool(exitRef.Flags["closed"]), "Exit is closed."},
				{"T", "autoclose", strconv.FormatBool(exitRef.Flags["autoclose"]), "Exit autocloses after 1 minute"},
				{"T", "lockable", strconv.FormatBool(exitRef.Flags["lockable"]), "Exit can be locked."},
				{"T", "unpickable", strconv.FormatBool(exitRef.Flags["unpickable"]), "Exit cannot be picked."},
				{"T", "locked", strconv.FormatBool(exitRef.Flags["locked"]), "Exit is locked when closed"},
				{"T", "hidden", strconv.FormatBool(exitRef.Flags["hidden"]), "Exit is hidden"},
				{"T", "invisible", strconv.FormatBool(exitRef.Flags["invisible"]), "Exit is invisible"},
				{"T", "levitate", strconv.FormatBool(exitRef.Flags["levitate"]), "Must be levitating to use exit"},
				{"T", "day_only", strconv.FormatBool(exitRef.Flags["day_only"]), "Can only use during day"},
				{"T", "night_only", strconv.FormatBool(exitRef.Flags["night_only"]), "Can only use at night"},
				{"T", "placement_dependent", strconv.FormatBool(exitRef.Flags["placement_dependent"]), "Location of exit matters."},
			})
			t.SetCaption("X = Cannot Modify,  T=Toggle to Edit, V=Edit by value name\nSee 'help edit' for more.")
			s.msg.Actor.SendGood(t.Render())

		} else {
			s.msg.Actor.SendBad("Couldn't find the exit in the current room.")
		}

	case "char":
		log.Println("Starting search...")
		charName := s.words[1]
		character := objects.ActiveCharacters.Find(charName)
		t := table.NewWriter()
		t.SetAllowedRowLength(rowLength)
		t.Style().Options.SeparateRows = true
		t.AppendHeader(table.Row{"Type", "Variable Name", "Value", "Description"})
		if character == nil {
			charData, err := data.LoadChar(charName)
			if err || charData == nil {
				s.msg.Actor.SendBad("Could not load the character from the database.")
				return
			}
			t.AppendRows([]table.Row{
				{"V", "name", charData["name"].(string), "Characters Name"},
				{"V", "description", text.WrapSoft(charData["description"].(string), rowLength/5), "Description"},
				{"X", "character_id", charData["character_id"].(int64), "DB Char ID"},
				{"V", "class", config.AvailableClasses[charData["class"].(int64)], "Character Class"},
				{"V", "race", config.AvailableRaces[charData["race"].(int64)], "Character Race"},
				{"V", "parentid", charData["parentid"].(int64), "Room ID that the character is in."},
				{"V", "title", charData["title"].(string), "Character Titles"},
				{"V", "bankgold", charData["bankgold"].(int64), "Amount of gold in the bank"},
				{"V", "gold", charData["gold"].(int64), "Amount of gold on character"},
				{"V", "experience", charData["experience"].(int64), "Character amount of experience."},
				{"V", "bonuspoints", charData["bonuspoints"].(int64), "Bonus points for character"},
				{"V", "passages", charData["passages"].(int64), "Passages"},
				{"V", "broadcasts", charData["broadcasts"].(int64), "Broadcasts (refreshes daily)"},
				{"V", "evals", charData["evals"].(int64), "Evaluates (refreshes daily)"},
				{"V", "stamcur", charData["stammax"].(int64), "Max Stamina"},
				{"V", "stamcur", charData["stamcur"].(int64), "Current Stamina"},
				{"V", "vitmax", charData["vitmax"].(int64), "Maximum Vitality"},
				{"V", "vitcur", charData["vitcur"].(int64), "Current Vitality"},
				{"V", "manamax", charData["manamax"].(int64), "Maximum Mana"},
				{"V", "manacur", charData["manacur"].(int64), "Current Mana"},
				{"V", "strcur", charData["strcur"].(int64), "Current Strength"},
				{"V", "dexcur", charData["dexcur"].(int64), "Current Dex"},
				{"V", "concur", charData["concur"].(int64), "Current Con"},
				{"V", "intcur", charData["intcur"].(int64), "Current Int"},
				{"V", "piecur", charData["piecur"].(int64), "Current Piety"},
				{"V", "tier", charData["tier"].(int64), "Character Tier"},
				{"V", "sharpexp", charData["sharpexp"].(int64), "Sharp Skill Experience"},
				{"V", "thrustexp", charData["thrustexp"].(int64), "Thrust Skill Experience"},
				{"V", "bluntexp", charData["bluntexp"].(int64), "Blunt Skill Experience"},
				{"V", "poleexp", charData["poleexp"].(int64), "Pole Skill Experience"},
				{"V", "missileexp", charData["missileexp"].(int64), "Missile Skill Experience"},
				{"V", "handexp", charData["handexp"].(int64), "Hand to Hand Experience"},
				{"V", "fireexp", charData["fireexp"].(int64), "Fire Affinity"},
				{"V", "airexp", charData["airexp"].(int64), "Air Affinity"},
				{"V", "earthexp", charData["earthexp"].(int64), "Earth Affinity"},
				{"V", "waterexp", charData["waterexp"].(int64), "Water Affinity"},
				{"V", "divinity", charData["divinity"].(int64), "Divinity"},
				{"V", "stealthexp", charData["stealthexp"].(int64), "Stealth Experience"},
				{"T", "darkvision", strconv.FormatBool(charData["flags"].(map[string]interface{})["darkvision"].(int64) != 0), "Permanent Dark Vision"},
			})
			t.SetCaption("X = Cannot Modify,  T=Toggle to Edit, V=Edit by value name\nSee 'help edit' for more.")
			s.msg.Actor.SendGood(t.Render())
		} else {
			objects.ActiveCharacters.Lock()
			t.AppendRows([]table.Row{
				{"V", "name", character.Name, "Characters Name"},
				{"V", "description", text.WrapSoft(character.Description, rowLength/5), "Description"},
				{"X", "character_id", character.CharId, "DB Char ID"},
				{"V", "parentid", character.ParentId, "Room ID that the character is in."},
				{"V", "title", character.Title, "Character Titles"},
				{"V", "bankgold", character.BankGold.Value, "Amount of gold in the bank"},
				{"V", "gold", character.Gold.Value, "Amount of gold on character"},
				{"V", "spells", strings.Join(character.Spells, ", "), "Available spells"},
				{"V", "experience", character.Experience.Value, "Character amount of experience."},
				{"V", "bonuspoints", character.BonusPoints.Value, "Bonus points for character"},
				{"V", "passages", character.Passages.Value, "Passages"},
				{"V", "broadcasts", character.Broadcasts, "Broadcasts (refreshes daily)"},
				{"V", "evals", character.Evals, "Evaluates (refreshes daily)"},
				{"V", "stammax", character.Stam.Max, "Max Stamina"},
				{"V", "stamcur", character.Stam.Current, "Current Stamina"},
				{"V", "vitmax", character.Vit.Max, "Maximum Vitality"},
				{"V", "vitcur", character.Vit.Current, "Current Vitality"},
				{"V", "manamax", character.Mana.Max, "Maximum Mana"},
				{"V", "manacur", character.Mana.Current, "Current Mana"},
				{"V", "strcur", character.Str.Current, "Current Strength"},
				{"V", "dexcur", character.Dex.Current, "Current Dex"},
				{"V", "concur", character.Con.Current, "Current Con"},
				{"V", "intcur", character.Int.Current, "Current Int"},
				{"V", "piecur", character.Pie.Current, "Current Piety"},
				{"V", "tier", character.Tier, "Character Tier"},
				{"V", "sharpexp", character.Skills[0].Value, "Sharp Skill Experience"},
				{"V", "thrustexp", character.Skills[1].Value, "Thrust Skill Experience"},
				{"V", "bluntexp", character.Skills[2].Value, "Blunt Skill Experience"},
				{"V", "poleexp", character.Skills[3].Value, "Pole Skill Experience"},
				{"V", "missileexp", character.Skills[4].Value, "Missile Skill Experience"},
				{"V", "handexp", character.Skills[5].Value, "Hand to Hand Experience"},
				{"V", "fireexp", character.Skills[6].Value, "Fire Affinity"},
				{"V", "airexp", character.Skills[7].Value, "Air Affinity"},
				{"V", "earthexp", character.Skills[8].Value, "Earth Affinity"},
				{"V", "waterexp", character.Skills[9].Value, "Water Affinity"},
				{"V", "divinity", character.Skills[10].Value, "Divinity"},
				{"V", "stealthexp", character.Skills[11].Value, "Stealth Experience"},
				{"T", "invisible", strconv.FormatBool(character.Flags["invisible"]), "Character is currently invisible"},
				{"T", "darkvision", strconv.FormatBool(character.Flags["darkvision"]), "Permanent Dark Vision"},
			})
			t.SetCaption("X = Cannot Modify,  T=Toggle to Edit, V=Edit by value name\nSee 'help edit' for more.")
			s.msg.Actor.SendGood(t.Render())
			objects.ActiveCharacters.Unlock()
			return
		}
		return
	default:
		s.msg.Actor.SendBad("Couldn't figure out what to examine.")
	}

	s.ok = true
	return
}
