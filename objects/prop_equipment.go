package objects

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"strings"
)

type Equipment struct{
	// Status
	Armor int
	Weight int
	DamageIgnore int
	// TODO: Eventually create an effect system for equipment

	Head *Item
	Chest *Item
	Neck *Item
	Legs *Item
	Feet *Item
	Arms *Item
	Hands *Item
	Ring1 *Item
	Ring2 *Item

	Ammo *Item

	// Hands, can hold shield or weapon
	Main *Item
	Off *Item
}

func (e *Equipment) List() []*Item {
	equipList := make([]*Item, 0)

	if e.Head != (*Item)(nil) { equipList = append(equipList, e.Head) }
	if e.Chest != (*Item)(nil) { equipList = append(equipList, e.Chest) }
	if e.Neck != (*Item)(nil) { equipList = append(equipList, e.Neck) }
	if e.Legs != (*Item)(nil) { equipList = append(equipList, e.Legs) }
	if e.Feet != (*Item)(nil) { equipList = append(equipList, e.Feet) }
	if e.Arms != (*Item)(nil) { equipList = append(equipList, e.Arms) }
	if e.Hands != (*Item)(nil) { equipList = append(equipList, e.Hands) }
	if e.Ring1 != (*Item)(nil) { equipList = append(equipList, e.Ring1) }
	if e.Ring2 != (*Item)(nil) { equipList = append(equipList, e.Ring2) }
	if e.Main != (*Item)(nil) { equipList = append(equipList, e.Main) }
	if e.Off != (*Item)(nil) { equipList = append(equipList, e.Off) }

	return equipList
}

func (e *Equipment) GetText(ref string) string {
	if ref == "head" && e.Head != (*Item)(nil) { return e.Head.Name }
	if ref == "chest" && e.Chest != (*Item)(nil) { return e.Chest.Name}
	if ref == "neck" && e.Neck != (*Item)(nil) { return e.Neck.Name }
	if ref == "legs" && e.Legs != (*Item)(nil) { return e.Legs.Name }
	if ref == "feet" && e.Feet != (*Item)(nil) { return e.Feet.Name }
	if ref == "arms" && e.Arms != (*Item)(nil) { return e.Arms.Name }
	if ref == "hands" && e.Hands != (*Item)(nil) { return e.Hands.Name }
	if ref == "ring1" && e.Ring1 != (*Item)(nil) { return e.Ring1.Name }
	if ref == "ring2" &&  e.Ring2 != (*Item)(nil) { return e.Ring2.Name }
	if ref == "main" && e.Main != (*Item)(nil) { return e.Main.Name }
	if ref == "off" && e.Off != (*Item)(nil) { return e.Off.Name }
	return ""
}

// Search the ItemInventory to return a specific instance of something
func (e *Equipment) Search(alias string) *Item {
	for _, c := range e.List() {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)){
			return c
		}
	}

	return nil
}

func (e *Equipment) Equip(item *Item) (ok bool){
	ok = false
	if item.ItemType == 5 && e.Chest == (*Item)(nil) { e.Chest = item; ok = true}  //body",
	if item.ItemType == 6 && e.Off == (*Item)(nil) { e.Off = item; ok=true }  //device",
	if item.ItemType == 7 && e.Off == (*Item)(nil) { e.Off = item; ok=true }  //scroll",
	if item.ItemType == 8 && e.Off == (*Item)(nil) { e.Off = item; ok=true }  //wand",
	if item.ItemType == 15 && e.Ammo == (*Item)(nil) && e.Main != (*Item)(nil) { if e.Main.ItemType == 4 { e.Ammo = item; ok=true } }  //ammo",
	if item.ItemType == 16 && e.Off == (*Item)(nil) { e.Off = item; ok=true }  //instrument",
	if item.ItemType == 17 && e.Off == (*Item)(nil) { e.Off = item; ok=true }  //beverage",
	if item.ItemType == 19 && e.Feet == (*Item)(nil) { e.Feet = item; ok=true }  //feet",
	if item.ItemType == 20 && e.Legs == (*Item)(nil) { e.Legs = item; ok=true }  //legs",
	if item.ItemType == 21 && e.Arms == (*Item)(nil) { e.Arms = item; ok=true }  //arms",
	if item.ItemType == 22 && e.Neck == (*Item)(nil) { e.Neck = item; ok=true }  //neck",
	if item.ItemType == 23 && e.Off == (*Item)(nil) { e.Off = item; ok=true }  //shield",
	if item.ItemType == 24 && (e.Ring1 == (*Item)(nil) || e.Ring2 == (*Item)(nil)) { if e.Ring1 == (*Item)(nil) { e.Ring1 = item } else { e.Ring2 = item}; ok=true }  //finger",
	if item.ItemType == 25 && e.Head == (*Item)(nil) { e.Head = item; ok=true }  //head",
	if item.ItemType == 26 && e.Hands == (*Item)(nil) { e.Hands = item; ok=true }  //hands",
	if item.ItemType == 0 && e.Main == (*Item)(nil) { e.Main = item; ok=true }  //sharp",
	if item.ItemType == 1 && e.Main == (*Item)(nil) { e.Main = item; ok=true }  //thrust",
	if item.ItemType == 2 && e.Main == (*Item)(nil) { e.Main = item; ok=true }  //blunt",
	if item.ItemType == 3 && e.Main == (*Item)(nil) { e.Main = item; ok=true }  //pole",
	if item.ItemType == 4 && e.Main == (*Item)(nil) { e.Main = item; ok=true }  //range",

	// Update armor values
	if ok {
		e.Armor += item.Armor
		e.Weight += item.Weight
	}
	return ok
}

// Attempt to unequip by name, or type
func (e *Equipment) Unequip(alias string) (ok bool, item *Item){
	ok = false

	if e.Head != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Head.Name), strings.ToLower(alias)) { item = e.Head; e.Head = (*Item)(nil); ok=true }}
	if e.Chest != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Chest.Name), strings.ToLower(alias)){  item = e.Chest; e.Chest = (*Item)(nil); ok=true}}
	if e.Neck != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Neck.Name), strings.ToLower(alias)){  item = e.Neck; e.Neck = (*Item)(nil); ok=true}}
	if e.Legs != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Legs.Name), strings.ToLower(alias)){  item = e.Legs; e.Legs = (*Item)(nil); ok=true}}
	if e.Feet != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Feet.Name), strings.ToLower(alias)){  item = e.Feet; e.Feet = (*Item)(nil); ok=true}}
	if e.Arms != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Arms.Name), strings.ToLower(alias)){  item = e.Arms; e.Arms = (*Item)(nil); ok=true}}
	if e.Hands != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Hands.Name), strings.ToLower(alias)){  item = e.Hands; e.Hands = (*Item)(nil); ok=true}}
	if e.Ring1 != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Ring1.Name), strings.ToLower(alias)){  item = e.Ring1; e.Ring1 = (*Item)(nil); ok=true}}
	if e.Ring2 != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Ring2.Name), strings.ToLower(alias)){  item = e.Ring2; e.Ring2 = (*Item)(nil); ok=true}}
	if e.Main != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Main.Name), strings.ToLower(alias)){  item = e.Main; e.Main = (*Item)(nil); ok=true}}
	if e.Off != (*Item)(nil) { if strings.Contains(strings.ToLower(e.Off.Name), strings.ToLower(alias)){  item = e.Off; e.Off = (*Item)(nil); ok=true}}

	// Update armor values
	if ok && item != (*Item)(nil) {
		e.Armor -= item.Armor
		e.Weight -= item.Weight
	}
	return ok, item
}

func (e *Equipment) Jsonify() string {
	itemList := make([]map[string]interface{}, 0)

	if e.Head != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Head)) }
	if e.Chest != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Chest)) }
	if e.Neck != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Neck)) }
	if e.Legs != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Legs)) }
	if e.Feet != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Feet)) }
	if e.Arms != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Arms)) }
	if e.Hands != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Hands)) }
	if e.Ring1 != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Ring1)) }
	if e.Ring2 != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Ring2)) }
	if e.Main != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Main)) }
	if e.Off != (*Item)(nil) { itemList = append(itemList, ReturnItemInstanceProps(e.Off)) }

	data, err := json.Marshal(itemList)
	if err != nil {
		return "[]"
	} else {
		return string(data)
	}
}

func RestoreEquipment(jsonString string) *Equipment{
	obj := make([]map[string]interface{}, 0)
	NewEquipment := &Equipment{}
	err := json.Unmarshal([]byte(jsonString), &obj)
	if err != nil {
		return NewEquipment
	}
	for _, item := range obj {
		newItem := Item{}
		copier.Copy(&newItem, Items[int(item["itemId"].(float64))])
		newItem.Name = item["name"].(string)
		newItem.MaxUses	= int(item["uses"].(float64))
		newItem.Flags["magic"] = int(item["magic"].(float64)) != 0
		newItem.Spell = item["spell"].(string)
		newItem.Armor =	int(item["armor"].(float64))
		NewEquipment.Equip(&newItem)
	}
	return NewEquipment
}
