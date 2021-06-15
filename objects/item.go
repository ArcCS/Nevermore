package objects

import (
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
	MaxUses      int
	Value        int
	Spell        string
	StorePrice 	int
	Adjustment int

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
			Commands: DeserializeCommands(itemData["commands"].(string)),
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
		return i.Weight + i.Storage.TotalWeight
	} else {
		return i.Weight
	}
}

func (i *Item) Look() string {
	resString := i.Description + "\n\n"
	if i.ItemType==9 {
		items := i.Storage.ReducedList()
		if len(items) > 0 {
			resString += "The " + i.Name + " contains " + strconv.Itoa(len(i.Storage.Contents)) + " items: \n" + items
		}
	}
	return resString
}

// DisplayName Return a display name with numerics
func (i *Item) DisplayName() string{
	typeReturn := 0
	// Mapping value definitions
	switch i.ItemType{
		case 0:  typeReturn = 1
		case 1:  typeReturn = 1
		case 2:  typeReturn = 1
		case 3:  typeReturn = 1
		case 4:  typeReturn = 1
		case 5:  typeReturn = 2
		case 15: typeReturn = 1
		case 16: typeReturn = 1
		case 19: typeReturn = 2
		case 20: typeReturn = 2
		case 21: typeReturn = 2
		case 22: typeReturn = 2
		case 23: typeReturn = 2
		case 24: typeReturn = 2
		case 25: typeReturn = 2
		case 26: typeReturn = 2
	}
	if typeReturn == 1 {
		return i.Name + " (" + strconv.Itoa(i.Adjustment) + ")"
	}else if typeReturn == 2 {
		if i.Armor > 0 {
			return i.Name + " (" + strconv.Itoa(i.Armor) + ")"
		}else{
			return i.Name
		}
	}
	return i.Name
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

// Function to return only the modifiable properties
func ReturnItemInstanceProps(item *Item) map[string]interface{} {
	serialList := map[string]interface{}{
		"itemId": item.ItemId,
		"name":   item.Name,
		"uses":   item.MaxUses,
		"magic":  utils.Btoi(item.Flags["magic"]),
		"spell":  item.Spell,
		"armor":  item.Armor,
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
