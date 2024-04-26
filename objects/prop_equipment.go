package objects

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"log"
	"math/rand"
	"strings"
)

type Equipment struct {
	// Status
	Armor int

	Head  *Item
	Chest *Item
	Neck  *Item
	Legs  *Item
	Feet  *Item
	Arms  *Item
	Hands *Item
	Ring1 *Item
	Ring2 *Item

	Ammo *Item

	// Hands, can hold shield or weapon
	Main *Item
	Off  *Item

	FlagOn            func(flagName string, provider string)
	FlagOff           func(flagName string, provider string)
	CanEquip          func(item *Item) (bool, string)
	ReturnToInventory func(item *Item)
}

func (e *Equipment) GetWeight() (total int) {
	for _, item := range e.List() {
		total += item.GetWeight()
	}
	return total
}

func (e *Equipment) List() []*Item {
	equipList := make([]*Item, 0)

	if e.Head != (*Item)(nil) {
		equipList = append(equipList, e.Head)
	}
	if e.Chest != (*Item)(nil) {
		equipList = append(equipList, e.Chest)
	}
	if e.Neck != (*Item)(nil) {
		equipList = append(equipList, e.Neck)
	}
	if e.Legs != (*Item)(nil) {
		equipList = append(equipList, e.Legs)
	}
	if e.Feet != (*Item)(nil) {
		equipList = append(equipList, e.Feet)
	}
	if e.Arms != (*Item)(nil) {
		equipList = append(equipList, e.Arms)
	}
	if e.Hands != (*Item)(nil) {
		equipList = append(equipList, e.Hands)
	}
	if e.Ring1 != (*Item)(nil) {
		equipList = append(equipList, e.Ring1)
	}
	if e.Ring2 != (*Item)(nil) {
		equipList = append(equipList, e.Ring2)
	}
	if e.Main != (*Item)(nil) {
		equipList = append(equipList, e.Main)
	}
	if e.Off != (*Item)(nil) {
		equipList = append(equipList, e.Off)
	}

	return equipList
}

func (e *Equipment) GetText(ref string) string {
	if ref == "head" && e.Head != (*Item)(nil) {
		return e.Head.DisplayName()
	}
	if ref == "chest" && e.Chest != (*Item)(nil) {
		return e.Chest.DisplayName()
	}
	if ref == "neck" && e.Neck != (*Item)(nil) {
		return e.Neck.DisplayName()
	}
	if ref == "legs" && e.Legs != (*Item)(nil) {
		return e.Legs.DisplayName()
	}
	if ref == "feet" && e.Feet != (*Item)(nil) {
		return e.Feet.DisplayName()
	}
	if ref == "arms" && e.Arms != (*Item)(nil) {
		return e.Arms.DisplayName()
	}
	if ref == "hands" && e.Hands != (*Item)(nil) {
		return e.Hands.DisplayName()
	}
	if ref == "ring1" && e.Ring1 != (*Item)(nil) {
		return e.Ring1.DisplayName()
	}
	if ref == "ring2" && e.Ring2 != (*Item)(nil) {
		return e.Ring2.DisplayName()
	}
	if ref == "main" && e.Main != (*Item)(nil) {
		return e.Main.DisplayName()
	}
	if ref == "off" && e.Off != (*Item)(nil) {
		return e.Off.DisplayName()
	}
	return ""
}

func (e *Equipment) DamageRandomArmor() (retString string) {
	armorList := make([]string, 0)
	if e.Head != nil {
		if e.Head.Armor > 0 {
			armorList = append(armorList, "head")
		}
	}
	if e.Chest != nil {
		if e.Chest.Armor > 0 {
			armorList = append(armorList, "chest")
		}
	}
	if e.Neck != nil {
		if e.Neck.Armor > 0 {
			armorList = append(armorList, "neck")
		}
	}
	if e.Legs != nil {
		if e.Legs.Armor > 0 {
			armorList = append(armorList, "legs")
		}
	}
	if e.Feet != nil {
		if e.Feet.Armor > 0 {
			armorList = append(armorList, "feet")
		}
	}
	if e.Arms != nil {
		if e.Arms.Armor > 0 {
			armorList = append(armorList, "arms")
		}
	}
	if e.Hands != nil {
		if e.Hands.Armor > 0 {
			armorList = append(armorList, "hands")
		}
	}
	if e.Ring1 != nil {
		if e.Ring1.Armor > 0 {
			armorList = append(armorList, "ring1")
		}
	}
	if e.Ring2 != nil {
		if e.Ring2.Armor > 0 {
			armorList = append(armorList, "ring2")
		}
	}
	if e.Off != nil {
		if e.Off.ItemType == 23 || e.Off.Armor > 0 {
			armorList = append(armorList, "off")
		}
	}
	if e.Main != nil {
		if e.Main.Armor > 0 {
			armorList = append(armorList, "main")
		}
	}

	if len(armorList) > 0 {
		damageItem := armorList[rand.Intn(len(armorList))]
		if damageItem == "head" {
			e.Head.MaxUses -= 1
			if e.Head.MaxUses <= 0 {
				retString = "Your " + e.Head.DisplayName() + " falls apart."
				e.UnequipSpecific("head")
				return
			}
			return ""
		} else if damageItem == "chest" {
			e.Chest.MaxUses -= 1
			if e.Chest.MaxUses <= 0 {
				retString = "Your " + e.Chest.DisplayName() + " falls apart."
				e.UnequipSpecific("chest")
				return
			}
			return ""
		} else if damageItem == "neck" {
			e.Neck.MaxUses -= 1
			if e.Neck.MaxUses <= 0 {
				retString = "Your " + e.Neck.DisplayName() + " falls apart."
				e.UnequipSpecific("neck")
				return
			}
			return ""
		} else if damageItem == "legs" {
			e.Legs.MaxUses -= 1
			if e.Legs.MaxUses <= 0 {
				retString = "Your " + e.Legs.DisplayName() + " falls apart."
				e.UnequipSpecific("legs")
				return
			}
			return ""
		} else if damageItem == "feet" {
			e.Feet.MaxUses -= 1
			if e.Feet.MaxUses <= 0 {
				retString = "Your " + e.Feet.DisplayName() + " falls apart."
				e.UnequipSpecific("feet")
				return
			}
			return ""
		} else if damageItem == "arms" {
			e.Arms.MaxUses -= 1
			if e.Arms.MaxUses <= 0 {
				retString = "Your " + e.Arms.DisplayName() + " falls apart."
				e.UnequipSpecific("arms")
				return
			}
			return ""
		} else if damageItem == "hands" {
			e.Hands.MaxUses -= 1
			if e.Hands.MaxUses <= 0 {
				retString = "Your " + e.Hands.DisplayName() + " falls apart."
				e.UnequipSpecific("hands")
				return
			}
			return ""
		} else if damageItem == "ring1" {
			e.Ring1.MaxUses -= 1
			if e.Ring1.MaxUses <= 0 {
				retString = "Your " + e.Ring1.DisplayName() + " falls apart."
				e.UnequipSpecific("ring1")
				return
			}
			return ""
		} else if damageItem == "ring2" {
			e.Ring2.MaxUses -= 1
			if e.Ring2.MaxUses <= 0 {
				retString = "Your " + e.Ring2.DisplayName() + " falls apart."
				e.UnequipSpecific("ring2")
				return
			}
			return ""
		} else if damageItem == "off" {
			e.Off.MaxUses -= 1
			if e.Off.MaxUses <= 0 {
				retString = "Your " + e.Off.DisplayName() + " falls apart."
				e.UnequipSpecific("off")
				return
			}
			return ""
		} else if damageItem == "main" {
			e.Main.MaxUses -= 1
			if e.Main.MaxUses <= 0 {
				retString = "Your " + e.Main.DisplayName() + " falls apart."
				e.UnequipSpecific("main")
				return
			}
			return ""
		}

	}

	return ""
}

func (e *Equipment) DamageWeapon(whichHand string, damage int) string {
	if whichHand == "main" {
		e.Main.MaxUses -= damage
		if e.Main.MaxUses <= 0 {
			e.ReturnToInventory(e.Main)
			e.Main = (*Item)(nil)
			return "Your weapon breaks!"
		}
	} else if whichHand == "off" {
		e.Off.MaxUses -= damage
		if e.Off.MaxUses <= 0 {
			e.ReturnToInventory(e.Off)
			e.Off = (*Item)(nil)
			return "Your instrument breaks!"
		}
	}
	return ""
}

// Search the ItemInventory to return a specific instance of something
func (e *Equipment) Search(alias string, nameNum int) *Item {
	passes := 1
	for _, c := range e.List() {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)) {
			if passes == nameNum {
				return c
			} else {
				passes++
			}
		}
	}

	return nil
}

func (e *Equipment) Equip(item *Item, charClass int) (ok bool) {
	ok = false
	itemSlot := ""

	if item.ItemType == 5 && e.Chest == (*Item)(nil) {
		e.Chest = item
		itemSlot = "chest"
		ok = true
	} //body
	if item.ItemType == 6 && e.Off == (*Item)(nil) {
		e.Off = item
		itemSlot = "off"
		ok = true
	} //device
	if item.ItemType == 7 && e.Off == (*Item)(nil) {
		e.Off = item
		itemSlot = "off"
		ok = true
	} //scroll
	if item.ItemType == 8 && e.Off == (*Item)(nil) {
		e.Off = item
		itemSlot = "off"
		ok = true
	} //wand
	if item.ItemType == 12 && e.Off == (*Item)(nil) {
		e.Off = item
		itemSlot = "off"
		ok = true
	} //light source
	if item.ItemType == 13 && e.Off == (*Item)(nil) {
		e.Off = item
		itemSlot = "off"
		ok = true
	} //just random crap to hold I guess.
	/* TODO: Lets look at ammo later
	if item.ItemType == 15 && e.Ammo == (*Item)(nil) && e.Main != (*Item)(nil) {
		if e.Main.ItemType == 4 {
			e.Ammo = item
			ok = true
		}
	} //ammo", */
	if item.ItemType == 16 && e.Off == (*Item)(nil) {
		e.Off = item
		itemSlot = "off"
		ok = true
	} //instrument
	if item.ItemType == 17 && e.Off == (*Item)(nil) {
		e.Off = item
		itemSlot = "off"
		ok = true
	} //beverage
	if item.ItemType == 19 && e.Feet == (*Item)(nil) {
		e.Feet = item
		itemSlot = "feet"
		ok = true
	} //feet
	if item.ItemType == 20 && e.Legs == (*Item)(nil) {
		e.Legs = item
		itemSlot = "legs"
		ok = true
	} //legs
	if item.ItemType == 21 && e.Arms == (*Item)(nil) {
		e.Arms = item
		itemSlot = "arms"
		ok = true
	} //arms
	if item.ItemType == 22 && e.Neck == (*Item)(nil) {
		e.Neck = item
		itemSlot = "neck"
		ok = true
	} //neck
	if item.ItemType == 23 && e.Off == (*Item)(nil) {
		e.Off = item
		itemSlot = "off"
		ok = true
	} //shield
	if item.ItemType == 24 && (e.Ring1 == (*Item)(nil) || e.Ring2 == (*Item)(nil)) {
		if e.Ring1 == (*Item)(nil) {
			e.Ring1 = item
			itemSlot = "ring1"
		} else {
			e.Ring2 = item
			itemSlot = "ring2"
		}
		ok = true
	} //finger
	if item.ItemType == 25 && e.Head == (*Item)(nil) {
		e.Head = item
		itemSlot = "head"
		ok = true
	} //head
	if item.ItemType == 26 && e.Hands == (*Item)(nil) {
		e.Hands = item
		itemSlot = "hands"
		ok = true
	} //hands
	if item.ItemType >= 0 && item.ItemType <= 3 { // 0: sharp, 1: thrust, 2: blunt, 3: pole
		if e.Main == (*Item)(nil) {
			log.Printf("hi")
			e.Main = item
			itemSlot = "main"
			ok = true
		}
	}
	if item.ItemType == 4 && e.Main == (*Item)(nil) {
		e.Main = item
		itemSlot = "main"
		ok = true
	} //range

	// Update armor values
	if ok {
		if e.FlagOn != nil && (item.Flags["light"] || item.ItemType == 12) {
			e.FlagOn("light", itemSlot)
		}
		e.Armor += item.Armor
	}
	return ok
}

// UnequipSpecific removes a specific slot rather than searching for a name
func (e *Equipment) UnequipSpecific(alias string) (ok bool) {
	ok = true
	iArmor := 0
	lightBearing := false

	if alias == "head" {
		if e.Head != (*Item)(nil) {
			iArmor = e.Head.Armor
			lightBearing = e.Head.Flags["light"]
			e.Head = (*Item)(nil)
		} else {
			return false
		}
	} else if alias == "chest" {
		if e.Chest != (*Item)(nil) {
			iArmor = e.Chest.Armor
			lightBearing = e.Chest.Flags["light"]
			e.Chest = (*Item)(nil)
		} else {
			return false
		}
	} else if alias == "neck" {
		if e.Neck != (*Item)(nil) {
			iArmor = e.Neck.Armor
			lightBearing = e.Neck.Flags["light"]
			e.Neck = (*Item)(nil)
		} else {
			return false
		}
	} else if alias == "legs" {
		if e.Legs != (*Item)(nil) {
			iArmor = e.Legs.Armor
			lightBearing = e.Legs.Flags["light"]
			e.Legs = (*Item)(nil)
		} else {
			return false
		}
	} else if alias == "feet" {
		if e.Feet != (*Item)(nil) {
			iArmor = e.Feet.Armor
			lightBearing = e.Feet.Flags["light"]
			e.Feet = (*Item)(nil)
		} else {
			return false
		}
	} else if alias == "arms" {
		if e.Arms != (*Item)(nil) {
			iArmor = e.Arms.Armor
			lightBearing = e.Arms.Flags["light"]
			e.Arms = (*Item)(nil)
		} else {
			return false
		}
	} else if alias == "hands" {
		if e.Hands != (*Item)(nil) {
			iArmor = e.Hands.Armor
			lightBearing = e.Hands.Flags["light"]
			e.Hands = (*Item)(nil)
		} else {
			return false
		}
	} else if alias == "ring1" {
		if e.Ring1 != (*Item)(nil) {
			iArmor = e.Ring1.Armor
			lightBearing = e.Ring1.Flags["light"]
			e.Ring1 = (*Item)(nil)
		} else {
			return false
		}
	} else if alias == "ring2" {
		if e.Ring2 != (*Item)(nil) {
			iArmor = e.Ring2.Armor
			lightBearing = e.Ring2.Flags["light"]
			e.Ring2 = (*Item)(nil)
		} else {
			return false
		}
	} else if alias == "off" {
		if e.Off != (*Item)(nil) {
			iArmor = e.Off.Armor
			lightBearing = e.Off.Flags["light"]
			e.Off = (*Item)(nil)
		} else {
			return false
		}
	} else if alias == "main" {
		if e.Main != (*Item)(nil) {
			iArmor = e.Main.Armor
			lightBearing = e.Main.Flags["light"]
			e.Main = (*Item)(nil)
		} else {
			return false
		}
	} else {
		return false
	}

	if lightBearing {
		e.FlagOff("light", alias)
	}
	e.Armor -= iArmor
	return ok
}

// Unequip Attempt to unequip by name, or type
func (e *Equipment) Unequip(alias string) (ok bool, item *Item) {
	ok = false
	itemSlot := ""
	if e.Head != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Head.Name), strings.ToLower(alias)) {
			item = e.Head
			e.Head = (*Item)(nil)
			itemSlot = "head"
			ok = true
		}
	}
	if e.Chest != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Chest.Name), strings.ToLower(alias)) {
			item = e.Chest
			e.Chest = (*Item)(nil)
			itemSlot = "chest"
			ok = true
		}
	}
	if e.Neck != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Neck.Name), strings.ToLower(alias)) {
			item = e.Neck
			e.Neck = (*Item)(nil)
			itemSlot = "neck"
			ok = true
		}
	}
	if e.Legs != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Legs.Name), strings.ToLower(alias)) {
			item = e.Legs
			e.Legs = (*Item)(nil)
			itemSlot = "legs"
			ok = true
		}
	}
	if e.Feet != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Feet.Name), strings.ToLower(alias)) {
			item = e.Feet
			e.Feet = (*Item)(nil)
			itemSlot = "feet"
			ok = true
		}
	}
	if e.Arms != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Arms.Name), strings.ToLower(alias)) {
			item = e.Arms
			e.Arms = (*Item)(nil)
			itemSlot = "arms"
			ok = true
		}
	}
	if e.Hands != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Hands.Name), strings.ToLower(alias)) {
			item = e.Hands
			e.Hands = (*Item)(nil)
			itemSlot = "hands"
			ok = true
		}
	}
	if e.Ring1 != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Ring1.Name), strings.ToLower(alias)) {
			item = e.Ring1
			e.Ring1 = (*Item)(nil)
			itemSlot = "ring1"
			ok = true
		}
	}
	if e.Ring2 != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Ring2.Name), strings.ToLower(alias)) {
			item = e.Ring2
			e.Ring2 = (*Item)(nil)
			itemSlot = "ring2"
			ok = true
		}
	}
	if e.Main != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Main.Name), strings.ToLower(alias)) {
			item = e.Main
			e.Main = (*Item)(nil)
			itemSlot = "main"
			ok = true
		}
	}
	if e.Off != (*Item)(nil) && ok == false {
		if strings.Contains(strings.ToLower(e.Off.Name), strings.ToLower(alias)) {
			item = e.Off
			e.Off = (*Item)(nil)
			itemSlot = "off"
			ok = true
		}
	}

	// Update armor values
	if ok && item != (*Item)(nil) {
		if item.Flags["light"] || item.ItemType == 12 {
			e.FlagOff("light", itemSlot)
		}
		e.Armor -= item.Armor
	}
	return ok, item
}

// FindLocation Attempt to find an item by name, return location
func (e *Equipment) FindLocation(alias string) (slot string) {
	itemSlot := ""
	if e.Head != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Head.Name), strings.ToLower(alias)) {
			itemSlot = "head"
		}
	}
	if e.Chest != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Chest.Name), strings.ToLower(alias)) {
			itemSlot = "chest"
		}
	}
	if e.Neck != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Neck.Name), strings.ToLower(alias)) {
			itemSlot = "neck"
		}
	}
	if e.Legs != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Legs.Name), strings.ToLower(alias)) {
			itemSlot = "legs"
		}
	}
	if e.Feet != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Feet.Name), strings.ToLower(alias)) {
			itemSlot = "feet"
		}
	}
	if e.Arms != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Arms.Name), strings.ToLower(alias)) {
			itemSlot = "arms"
		}
	}
	if e.Hands != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Hands.Name), strings.ToLower(alias)) {
			itemSlot = "hands"
		}
	}
	if e.Ring1 != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Ring1.Name), strings.ToLower(alias)) {
			itemSlot = "ring1"
		}
	}
	if e.Ring2 != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Ring2.Name), strings.ToLower(alias)) {
			itemSlot = "ring2"
		}
	}
	if e.Main != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Main.Name), strings.ToLower(alias)) {
			itemSlot = "main"
		}
	}
	if e.Off != (*Item)(nil) {
		if strings.Contains(strings.ToLower(e.Off.Name), strings.ToLower(alias)) {
			itemSlot = "off"
		}
	}

	return itemSlot
}

// UnequipAll Remove all equipment
func (e *Equipment) UnequipAll() (items []*Item) {

	if e.Head != (*Item)(nil) {
		items = append(items, e.Head)
		if e.Head.Flags["light"] {
			e.FlagOff("light", "head")
		}
		e.Head = (*Item)(nil)
	}
	if e.Chest != (*Item)(nil) {
		items = append(items, e.Chest)
		if e.Chest.Flags["light"] {
			e.FlagOff("light", "chest")
		}
		e.Chest = (*Item)(nil)
	}
	if e.Neck != (*Item)(nil) {
		items = append(items, e.Neck)
		if e.Neck.Flags["light"] {
			e.FlagOff("light", "neck")
		}
		e.Neck = (*Item)(nil)
	}
	if e.Legs != (*Item)(nil) {
		items = append(items, e.Legs)
		if e.Legs.Flags["light"] {
			e.FlagOff("light", "legs")
		}
		e.Legs = (*Item)(nil)
	}
	if e.Feet != (*Item)(nil) {
		items = append(items, e.Feet)
		if e.Feet.Flags["light"] {
			e.FlagOff("light", "feet")
		}
		e.Feet = (*Item)(nil)
	}
	if e.Arms != (*Item)(nil) {
		items = append(items, e.Arms)
		if e.Arms.Flags["light"] {
			e.FlagOff("light", "arms")
		}
		e.Arms = (*Item)(nil)
	}
	if e.Hands != (*Item)(nil) {
		items = append(items, e.Hands)
		if e.Hands.Flags["light"] {
			e.FlagOff("light", "hands")
		}
		e.Hands = (*Item)(nil)
	}
	if e.Ring1 != (*Item)(nil) {
		items = append(items, e.Ring1)
		if e.Ring1.Flags["light"] {
			e.FlagOff("light", "ring1")
		}
		e.Ring1 = (*Item)(nil)
	}
	if e.Ring2 != (*Item)(nil) {
		items = append(items, e.Ring2)
		if e.Ring2.Flags["light"] {
			e.FlagOff("light", "ring2")
		}
		e.Ring2 = (*Item)(nil)
	}
	if e.Main != (*Item)(nil) {
		items = append(items, e.Main)
		if e.Main.Flags["light"] {
			e.FlagOff("light", "main")
		}
		e.Main = (*Item)(nil)
	}
	if e.Off != (*Item)(nil) {
		items = append(items, e.Off)
		if e.Off.Flags["light"] || e.Off.ItemType == 12 {
			e.FlagOff("light", "off")
		}
		e.Off = (*Item)(nil)
	}

	e.Armor = 0

	return items
}

func (e *Equipment) Jsonify() string {
	itemList := make([]map[string]interface{}, 0)

	if e.Head != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Head))
	}
	if e.Chest != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Chest))
	}
	if e.Neck != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Neck))
	}
	if e.Legs != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Legs))
	}
	if e.Feet != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Feet))
	}
	if e.Arms != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Arms))
	}
	if e.Hands != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Hands))
	}
	if e.Ring1 != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Ring1))
	}
	if e.Ring2 != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Ring2))
	}
	if e.Main != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Main))
	}
	if e.Off != (*Item)(nil) {
		itemList = append(itemList, ReturnItemInstanceProps(e.Off))
	}

	data, err := json.Marshal(itemList)
	if err != nil {
		return "[]"
	} else {
		return string(data)
	}
}

func RestoreEquipment(jsonString string, charClass int) (*Equipment, []Item) {
	obj := make([]map[string]interface{}, 0)
	NewEquipment := &Equipment{}
	var ErroredEquipment []Item
	err := json.Unmarshal([]byte(jsonString), &obj)
	if err != nil {
		return NewEquipment, ErroredEquipment
	}
	for _, item := range obj {
		newItem := Item{}
		if err := copier.CopyWithOption(&newItem, Items[int(item["itemId"].(float64))], copier.Option{DeepCopy: true}); err != nil {
			log.Println("Error copying item during restore: ", err)
		}
		newItem.Name = item["name"].(string)
		newItem.MaxUses = int(item["uses"].(float64))
		newItem.Flags["magic"] = int(item["magic"].(float64)) != 0
		newItem.Spell = item["spell"].(string)
		newItem.Armor = int(item["armor"].(float64))
		if _, ok := item["light"]; ok {
			newItem.Flags["light"] = int(item["light"].(float64)) != 0
		}
		if _, ok := item["adjustment"]; ok {
			newItem.Adjustment = int(item["adjustment"].(float64))
		}
		ok := NewEquipment.Equip(&newItem, charClass)
		if !ok {
			ErroredEquipment = append(ErroredEquipment, newItem)
		}
	}
	return NewEquipment, ErroredEquipment
}

func (e *Equipment) CheckEquipment() {
	if e.Head != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Head); !ok {
			e.UnequipSpecific("head")
		}
	}
	if e.Chest != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Chest); !ok {
			e.UnequipSpecific("chest")
		}
	}
	if e.Neck != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Neck); !ok {
			e.UnequipSpecific("neck")
		}
	}
	if e.Legs != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Legs); !ok {
			e.UnequipSpecific("legs")
		}
	}
	if e.Feet != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Feet); !ok {
			e.UnequipSpecific("feet")
		}
	}
	if e.Arms != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Arms); !ok {
			e.UnequipSpecific("arms")
		}
	}
	if e.Hands != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Hands); !ok {
			e.UnequipSpecific("hands")
		}
	}
	if e.Ring1 != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Ring1); !ok {
			e.UnequipSpecific("ring1")
		}
	}
	if e.Ring2 != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Ring2); !ok {
			e.UnequipSpecific("ring2")
		}
	}
	if e.Main != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Main); !ok {
			e.UnequipSpecific("main")
		}
	}
	if e.Off != (*Item)(nil) {
		if ok, _ := e.CanEquip(e.Off); !ok {
			e.UnequipSpecific("off")
		}
	}
}

func (e *Equipment) PostEquipmentLight() {
	if e.Head != (*Item)(nil) {
		if e.Head.Flags["light"] {
			e.FlagOn("light", "head")
		}
	}
	if e.Chest != (*Item)(nil) {
		if e.Chest.Flags["light"] {
			e.FlagOn("light", "chest")
		}
	}
	if e.Neck != (*Item)(nil) {
		if e.Neck.Flags["light"] {
			e.FlagOn("light", "neck")
		}
	}
	if e.Legs != (*Item)(nil) {
		if e.Legs.Flags["light"] {
			e.FlagOn("light", "legs")
		}
	}
	if e.Feet != (*Item)(nil) {
		if e.Feet.Flags["light"] {
			e.FlagOn("light", "feet")
		}
	}
	if e.Arms != (*Item)(nil) {
		if e.Arms.Flags["light"] {
			e.FlagOn("light", "arms")
		}
	}
	if e.Hands != (*Item)(nil) {
		if e.Hands.Flags["light"] {
			e.FlagOn("light", "hands")
		}
	}
	if e.Ring1 != (*Item)(nil) {
		if e.Ring1.Flags["light"] {
			e.FlagOn("light", "ring1")
		}
	}
	if e.Ring2 != (*Item)(nil) {
		if e.Ring2.Flags["light"] {
			e.FlagOn("light", "ring2")
		}
	}
	if e.Main != (*Item)(nil) {
		if e.Main.Flags["light"] {
			e.FlagOn("light", "main")
		}
	}
	if e.Off != (*Item)(nil) {
		if e.Off.Flags["light"] || e.Off.ItemType == 12 {
			e.FlagOn("light", "off")
		}
	}
}
