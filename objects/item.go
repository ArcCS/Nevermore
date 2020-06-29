package objects

type Item struct {
	Object
	ParentItemId int
	ItemId       int
	ItemType         int
	Flags        map[string]bool
	Creator      string
	NumDice      int
	PlusDice     int
	SidesDice    int
	WeaponSpeed  int
	Armor        int
	MaxUses      int
	Value        int
	Spell        string

	Storage ItemInventory
	Weight  int
	LinkId int
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
		},
		0,
		int(itemData["item_id"].(int64)),
		int(itemData["type"].(int64)),
		make(map[string]bool),
		itemData["creator"].(string),
		int(itemData["ndice"].(int64)),
		int(itemData["pdice"].(int64)),
		int(itemData["sdice"].(int64)),
		int(itemData["weapon_speed"].(int64)),
		int(itemData["armor"].(int64)),
		int(itemData["max_uses"].(int64)),
		int(itemData["value"].(int64)),
		itemData["spell"].(string),
		ItemInventory{},
		int(itemData["weight"].(int64)),
		0,
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
	if i.ItemType == 9 && !i.Flags["weightless"] {
		return i.Weight + i.Storage.TotalWeight
	} else {
		return i.Weight
	}
}

func (i *Item) Look() string {
	return i.Description
}

func (i *Item) Use(parentId int, target int) {
	return
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
	// TODO: Invoke a static save as a new item
	return
}
