package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"strconv"
	"strings"
)

func init() {
	addHandler(examine{},
		"Usage:  examine (room|mob|object|exit) (name|#####) \n\n  Examine will display the item and all of it's modifiable properties",
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
			{"T", "repair", strconv.FormatBool(roomRef.Flags["repair"]), "Can repair items here."},
			{"T", "mana_drain", strconv.FormatBool(roomRef.Flags["mana_drain"]), "Drains mana per tick"},
			{"T", "no_summon", strconv.FormatBool(roomRef.Flags["no_summon"]), "Blocks summon"},
			{"T", "heal_fast", strconv.FormatBool(roomRef.Flags["heal_fast"]), "Regenerate quickly"},
			{"T", "no_teleport", strconv.FormatBool(roomRef.Flags["no_teleport"]), "Blocks teleporting out"},
			{"T", "no_word_of_recall", strconv.FormatBool(roomRef.Flags["no_word_of_recall"]), "Blocks WoD"},
			{"T", "no_scry", strconv.FormatBool(roomRef.Flags["no_scry"]), "Blocks clairvoyance"},
			{"T", "no_magic", strconv.FormatBool(roomRef.Flags["no_magic"]), "Cannot cast spells"},
			{"T", "dark_always", strconv.FormatBool(roomRef.Flags["dark_always"]), "Always dark"},
			{"T", "light_always", strconv.FormatBool(roomRef.Flags["light_always"]), "Always lit"},
			{"T", "natural_light", strconv.FormatBool(roomRef.Flags["natural_light"]), "Follows day/night cycle"},
			{"T", "indoors", strconv.FormatBool(roomRef.Flags["indoors"]), "Room is indoors"},
			{"T", "urban", strconv.FormatBool(roomRef.Flags["urban"]), "Room is part of a city"},
			{"T", "underground", strconv.FormatBool(roomRef.Flags["underground"]), "Room is underground"},
			{"T", "encounters_on",strconv.FormatBool(roomRef.Flags["encounters_on"]), "Monsters can spawn"},
			{"T", "fire", strconv.FormatBool(roomRef.Flags["fire"]), "Room causes and amplifies fire."},
			{"T", "water", strconv.FormatBool(roomRef.Flags["water"]), "Room causes and amplifies water."},
			{"T", "earth", strconv.FormatBool(roomRef.Flags["earth"]), "Room causes and amplifies earth."},
			{"T", "wind", strconv.FormatBool(roomRef.Flags["wind"]), "Room causes and amplifies wind."},
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
		mobRef := s.where.Mobs.Search(s.words[1], nameNum, true)
		if mobRef != nil {

			t := table.NewWriter()
			t.SetAllowedRowLength(rowLength)
			t.Style().Options.SeparateRows = true
			t.AppendHeader(table.Row{"Type", "Variable Name", "Value", "Description"})
			t.AppendRows([]table.Row{
				{"X", "mobId", strconv.Itoa(mobRef.MobId), "DB id for mob"},
				{"V", "name", mobRef.Name, "Mob Name"},
				{"V", "description",  text.WrapSoft(mobRef.Description, rowLength/5), "Mob description"},
				{"V", "level", strconv.Itoa(mobRef.Level), "Mob Level"},
				{"V", "experience", strconv.Itoa(mobRef.Experience), "Awarded Experience"},
				{"V", "gold", strconv.Itoa(mobRef.Gold), "Gold dropped"},
				{"V", "con", strconv.Itoa(mobRef.Con.Max), "Amt of Con"},
				{"V", "str", strconv.Itoa(mobRef.Str.Max), "Amt of Str"},
				{"V", "int", strconv.Itoa(mobRef.Int.Max), "Amt of Int"},
				{"V", "dex", strconv.Itoa(mobRef.Dex.Max), "Amt of Dex"},
				{"V", "pie", strconv.Itoa(mobRef.Pie.Max), "Amt of Pie"},
				{"V", "mana", strconv.Itoa(mobRef.Mana.Max), "Available Mana"},
				{"V", "stam", strconv.Itoa(mobRef.Stam.Max), "Total HP"},
				{"V", "ndice", strconv.Itoa(mobRef.NumDice), "Number of Dice"},
				{"V", "sdice", strconv.Itoa(mobRef.SidesDice), "Sides of Dice"},
				{"V", "pdice", strconv.Itoa(mobRef.PlusDice), "Plus modifier to dice"},
				{"V", "chancecast", strconv.Itoa(mobRef.ChanceCast), "Chance to cast spell"},
				{"V", "armor", strconv.Itoa(mobRef.Armor), "Amt of Armor"},
				{"V", "numwander", strconv.Itoa(mobRef.NumWander), "Number of ticks to wander"},
				{"V", "wimpyvalue", strconv.Itoa(mobRef.WimpyValue), "Amt of damage that causes a flee chance"},
				{"T", "hide_encounter", strconv.FormatBool(mobRef.Flags["hide_encounter"]), "Hide when encounter"},
				{"T", "invisible", strconv.FormatBool(mobRef.Flags["invisible"]), "Invisible Mob"},
				{"T", "permanent", strconv.FormatBool(mobRef.Flags["permanent"]), "Does not despawn from room"},
				{"T", "hostile", strconv.FormatBool(mobRef.Flags["hostile"]), "Mob is hostile"},
			})
			t.SetCaption("X = Cannot Modify,  T=Toggle to Edit, V=Edit by value name\nSee 'help edit' for more.")
			s.msg.Actor.SendGood(t.Render())

		} else {
			s.msg.Actor.SendBad("Couldn't find the object in the current room.")
		}

	case "item":
		nameNum := 1
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
				{"V", "armor", strconv.Itoa(objRef.Armor	), "Armor value from this item"},
				{"V", "ndice", strconv.Itoa(objRef.NumDice), "Number of Dice to Roll"},
				{"V", "sdice", strconv.Itoa(objRef.SidesDice), "Sides of Dice Being Rolled"},
				{"V", "pdice", strconv.Itoa(objRef.PlusDice), "Additional value to add to roll."},
				{"V", "spell", objRef.Spell, "Spell cast when used."},
				{"T", "always_crit", strconv.FormatBool(objRef.Flags["always_crit"]), "Always criticals when used"},
				{"T", "permanent", strconv.FormatBool(objRef.Flags["permanent"]), "Does not despawn on ground"},
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
		if exitRef, ok := s.where.Exits[strings.ToLower(s.words[1])]; ok {
			t := table.NewWriter()
			t.SetAllowedRowLength(rowLength)
			t.Style().Options.SeparateRows = true
			t.AppendHeader(table.Row{"Type", "Variable Name", "Value", "Description"})
			t.AppendRows([]table.Row{
				{"V", "name", exitRef.Name, "Name of the exit"},
				{"V", "description", exitRef.Description, "Peers into next room if empty."},
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
	default:
		s.msg.Actor.SendBad("Couldn't figure out what to examine.")
	}

	s.ok = true
	return
}
