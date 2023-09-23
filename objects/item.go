package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

type Item struct {
	Object
	ParentItemId int
	ItemId       int
	ItemType     int
	Flags        map[string]bool
	Creator      string
	NumDice      int
	PlusDice     int
	SidesDice    int
	Armor        int
	Armor_Class  int
	MaxUses      int
	Value        int
	Spell        string
	StorePrice   int
	Adjustment   int

	Storage *ItemInventory
	Weight  int
}

// Pop the room data
func LoadItem(itemData map[string]interface{}) (*Item, bool) {
	description := ""
	var ok bool
	if description, ok = itemData["description"].(string); !ok {
		description = "This one was missing a description.. maybe fix that?"
	}
	newItem := &Item{
		Object{
			Name:        itemData["name"].(string),
			Description: description,
			Placement:   3,
			Commands:    DeserializeCommands(itemData["commands"].(string)),
		},
		0,
		int(itemData["item_id"].(int64)),
		int(itemData["type"].(int64)),
		make(map[string]bool),
		itemData["creator"].(string),
		int(itemData["ndice"].(int64)),
		int(itemData["pdice"].(int64)),
		int(itemData["sdice"].(int64)),
		int(itemData["armor"].(int64)),
		int(itemData["armor_class"].(int64)),
		int(itemData["max_uses"].(int64)),
		int(itemData["value"].(int64)),
		itemData["spell"].(string),
		0,
		int(itemData["adjustment"].(int64)),
		&ItemInventory{},
		int(itemData["weight"].(int64)),
	}
	for k, v := range itemData["flags"].(map[string]interface{}) {
		if v == nil {
			newItem.Flags[k] = false
		} else {
			newItem.Flags[k] = v.(int64) != 0
		}
	}
	return newItem, true
}

func (i *Item) GetWeight() int {
	if i.ItemType == 9 && !i.Flags["weightless_chest"] {
		return i.Weight + i.Storage.GetTotalWeight()
	} else {
		return i.Weight
	}
}

func (i *Item) Look() string {
	resString := i.Description + "\n\n"
	if utils.IntIn(i.ItemType, config.WeaponTypes) {
		resString = "It is a " + config.ItemTypes[i.ItemType] + " weapon, and it" + i.ReturnState() + "\n" + resString
	}
	if utils.IntIn(i.ItemType, config.ArmorTypes) {
		resString = "It is a " + config.ArmorClass[i.Armor_Class] + " " + config.ItemTypes[i.ItemType] + " armor, and it" + i.ReturnState() + "\n" + resString
	}
	if i.ItemType == 9 {
		items := i.Storage.ReducedList()
		if len(items) > 0 {
			resString += "The " + i.Name + " contains " + strconv.Itoa(len(i.Storage.Contents)) + " items: \n" + items
		}
	}
	return resString
}

// DisplayName Return a display name with numerics
func (i *Item) DisplayName() string {
	typeReturn := 0
	preName := ""
	if i.Flags["magic"] {
		preName += "magic "
	}
	// Mapping value definitions
	switch i.ItemType {
	case 0:
		typeReturn = 1
	case 1:
		typeReturn = 1
	case 2:
		typeReturn = 1
	case 3:
		typeReturn = 1
	case 4:
		typeReturn = 1
	case 5:
		typeReturn = 2
	case 15:
		typeReturn = 1
	case 16:
		typeReturn = 1
	case 19:
		typeReturn = 2
	case 20:
		typeReturn = 2
	case 21:
		typeReturn = 2
	case 22:
		typeReturn = 2
	case 23:
		typeReturn = 2
	case 24:
		typeReturn = 2
	case 25:
		typeReturn = 2
	case 26:
		typeReturn = 2
	}
	if typeReturn == 1 {
		return preName + i.Name + " (" + strconv.Itoa(i.Adjustment) + ")"
	} else if typeReturn == 2 {
		if i.Armor > 0 {
			return preName + i.Name + " (" + strconv.Itoa(i.Armor) + ")"
		} else {
			return preName + i.Name
		}
	}
	return preName + i.Name
}

func (i *Item) ToggleFlag(flagName string) bool {
	if val, exists := i.Flags[flagName]; exists {
		i.Flags[flagName] = !val
		return true
	} else {
		return false
	}
}

func (i *Item) Save() {
	itemData := make(map[string]interface{})
	itemData["item_id"] = i.ItemId
	itemData["ndice"] = i.NumDice
	itemData["weight"] = i.Weight
	itemData["description"] = i.Description
	itemData["type"] = i.ItemType
	itemData["pdice"] = i.PlusDice
	itemData["armor"] = i.Armor
	itemData["armor_class"] = i.Armor_Class
	itemData["max_uses"] = i.MaxUses
	itemData["name"] = i.Name
	itemData["sdice"] = i.SidesDice
	itemData["value"] = i.Value
	itemData["spell"] = i.Spell
	itemData["always_crit"] = utils.Btoi(i.Flags["always_crit"])
	itemData["permanent"] = utils.Btoi(i.Flags["permanent"])
	itemData["magic"] = utils.Btoi(i.Flags["magic"])
	itemData["light"] = utils.Btoi(i.Flags["light"])
	itemData["no_take"] = utils.Btoi(i.Flags["no_take"])
	itemData["weightless_chest"] = utils.Btoi(i.Flags["weightless_chest"])
	itemData["adjustment"] = i.Adjustment
	itemData["commands"] = i.SerializeCommands()
	data.UpdateItem(itemData)
	return
}

func (i *Item) ReturnState() string {
	stamStatus := "slightly used"
	if i.MaxUses == Items[i.ItemId].MaxUses {
		stamStatus = "pristine"
	} else if i.MaxUses < (Items[i.ItemId].MaxUses - int(.90*float32(Items[i.ItemId].MaxUses))) {
		stamStatus = "about to break"
	} else if i.MaxUses < (Items[i.ItemId].MaxUses - int(.75*float32(Items[i.ItemId].MaxUses))) {
		stamStatus = "badly damaged"
	} else if i.MaxUses < (Items[i.ItemId].MaxUses - int(.5*float32(Items[i.ItemId].MaxUses))) {
		stamStatus = "well used"
	} else if i.MaxUses < (Items[i.ItemId].MaxUses - int(.25*float32(Items[i.ItemId].MaxUses))) {
		stamStatus = "used"
	}
	return " looks " + stamStatus
}

// Function to return only the modifiable properties
func ReturnItemInstanceProps(item *Item) map[string]interface{} {
	serialList := map[string]interface{}{
		"itemId":     item.ItemId,
		"name":       item.Name,
		"uses":       item.MaxUses,
		"adjustment": item.Adjustment, //  Adjustable by Mages
		"magic":      utils.Btoi(item.Flags["magic"]),
		"spell":      item.Spell,
		"light":      utils.Btoi(item.Flags["light"]),
		"armor":      item.Armor, // Adjustable by Paladins
	}
	if _, ok := item.Flags["infinite"]; ok {
		serialList["infinite"] = utils.Btoi(item.Flags["infinite"])
	}
	if item.StorePrice != 0 {
		serialList["store_price"] = item.StorePrice
	}
	if item.ItemType == 9 {
		serialList["contents"] = item.Storage.Jsonify()
	}
	return serialList
}

func (i *Item) Eval() string {

	stringOut := "You study the " + i.Name + " in your minds eye.... \n\n"

	if utils.IntIn(i.ItemType, []int{0, 1, 2, 3, 4}) { // Weapons
		stringOut += "It is a " + config.ItemTypes[i.ItemId] + " weapon. \n" +
			"It deals between " + strconv.Itoa(utils.RollMin(i.NumDice, i.PlusDice)+i.Adjustment) + " and " + strconv.Itoa(utils.RollMax(i.SidesDice, i.NumDice, i.PlusDice)+i.Adjustment) + " damage. \n" +
			"It has " + strconv.Itoa(i.MaxUses) + " uses before it breaks \n."
	} else if utils.IntIn(i.ItemType, []int{5, 26, 25, 24, 23, 22, 21, 20, 19}) { // Armor
		stringOut += "It is a " + config.ArmorClass[i.Armor_Class] + " " + config.ItemTypes[i.ItemId] + " armor. \n" +
			"It has " + strconv.Itoa(i.MaxUses) + " uses before it breaks. \n"
	} else if i.ItemType == 17 { // Beverage
		stringOut += "It is a beverage. \n" +
			"It has " + strconv.Itoa(i.MaxUses) + " sips remaining. \n"
	} else if i.ItemType == 18 { // Music
		stringOut += "It is sheet music. \n" +
			"It contains " + i.Spell + ".\n"
	} else if i.ItemType == 16 { // Instrument
		stringOut += "It is an instrument. \n"
	} else if i.ItemType == 15 { // Ammo  //TODO Implement if we do ammo.
		stringOut += "It is ammunition, and not implemented. \n"
	} else if i.ItemType == 6 || i.ItemType == 8 { //device/wand
		stringOut += "It is a " + config.ItemTypes[i.ItemType] + ". \n" +
			"It is charged with " + i.Spell + ".\n" +
			"It has " + strconv.Itoa(i.MaxUses) + " uses remaining. \n"
	} else if i.ItemType == 7 {
		stringOut += "It is a scroll. \n" +
			"It contains " + i.Spell + ".\n"
	} else if i.ItemType == 9 { //chest
		stringOut += "It is a container. \n" +
			"It can hold " + strconv.Itoa(i.MaxUses) + " items. \n"
		if i.Flags["weightless_chest"] {
			stringOut += "It holds it's contents weightlessly. \n"
		}
	} else if i.ItemType == 10 { //gold
		stringOut += "It is gold. \n"
	} else if i.ItemType == 11 { //key
		stringOut += "It is a key. \n"
	} else if i.ItemType == 12 { //light
		stringOut += "It is a light source. \n"
	} else if i.ItemType == 13 { //object
		stringOut += "It's just an object. \n"
	}

	if i.Flags["permanent"] {
		stringOut += "It is permanent. \n"
	}
	if i.Flags["no_take"] {
		stringOut += "It cannot be picked up \n"
	}

	stringOut += "\n You determine its weight to be " + strconv.Itoa(i.Weight) + "lbs. \nYou judge its value to be " + strconv.Itoa(i.Value) + " gold marks."
	return stringOut
}
