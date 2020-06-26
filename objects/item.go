package objects

type Item struct {
	Object
	ParentItemId int64
	ItemId int64
	Type int64
	Flags map[string]bool
	Creator string
	NumDice int64
	PlusDice int64
	SidesDice int64
	WeaponSpeed int64
	Armor int64
	MaxUses int64
	Value int64
	Spell string

	Storage ItemInventory
	Weight int64
}

// Pop the room data
func LoadItem(itemData map[string]interface{}) (*Item, bool){
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
		itemData["item_id"].(int64),
		itemData["type"].(int64),
		make(map[string]bool),
		itemData["creator"].(string),
		itemData["ndice"].(int64),
		itemData["pdice"].(int64),
		itemData["sdice"].(int64),
		itemData["weapon_speed"].(int64),
		itemData["armor"].(int64),
		itemData["max_uses"].(int64),
		itemData["value"].(int64),
		itemData["spell"].(string),
		ItemInventory{},
		itemData["weight"].(int64),
	}
	for k, v := range itemData["flags"].(map[string]interface{}){
		if v == nil{
			newItem.Flags[k] = false
		}else {
			newItem.Flags[k] = v.(int64) != 0
		}
	}
	return newItem, true
}

func (i *Item) GetWeight() int64{
	if i.Type == 9 && !i.Flags["weightless"] {
		return i.Weight + i.Storage.TotalWeight
	}else{
		return i.Weight
	}
}

func (i *Item) Look() string {
	return i.Description
}

func (i *Item) Use(parentId int64, target int64){
	return
}

func (i *Item) ToggleFlag(flagName string) bool {
	if val, exists := i.Flags[flagName]; exists{
		i.Flags[flagName] = !val
		return true
	}else{
		return false
	}
}

func (i *Item) Save()  {
	// TODO: Invoke a static save as a new item
return
}
