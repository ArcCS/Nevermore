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

	if e.Head != nil { equipList = append(equipList, e.Head) }
	if e.Chest != nil { equipList = append(equipList, e.Chest) }
	if e.Neck != nil { equipList = append(equipList, e.Neck) }
	if e.Legs != nil { equipList = append(equipList, e.Legs) }
	if e.Feet != nil { equipList = append(equipList, e.Feet) }
	if e.Arms != nil { equipList = append(equipList, e.Arms) }
	if e.Hands != nil { equipList = append(equipList, e.Hands) }
	if e.Ring1 != nil { equipList = append(equipList, e.Ring1) }
	if e.Ring2 != nil { equipList = append(equipList, e.Ring2) }
	if e.Main != nil { equipList = append(equipList, e.Main) }
	if e.Off != nil { equipList = append(equipList, e.Off) }

	return equipList
}

func (e *Equipment) GetText(ref string) string {
	if ref == "head" && e.Head != nil { return e.Head.Name }
	if ref == "chest" && e.Chest != nil { return e.Chest.Name}
	if ref == "neck" && e.Neck != nil { return e.Neck.Name }
	if ref == "legs" && e.Legs != nil { return e.Legs.Name }
	if ref == "feet" && e.Feet != nil { return e.Feet.Name }
	if ref == "arms" && e.Arms != nil { return e.Arms.Name }
	if ref == "hands" && e.Hands != nil { return e.Hands.Name }
	if ref == "ring1" && e.Ring1 != nil { return e.Ring1.Name }
	if ref == "ring2" &&  e.Ring2 != nil { return e.Ring2.Name }
	if ref == "main" && e.Main != nil { return e.Main.Name }
	if ref == "off" && e.Off != nil { return e.Off.Name }
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
	if item.ItemType == 5 && e.Chest == nil { e.Chest = item; ok = true}  //body",
	if item.ItemType == 6 && e.Off == nil { e.Off = item; ok=true }  //device",
	if item.ItemType == 7 && e.Off == nil { e.Off = item; ok=true }  //scroll",
	if item.ItemType == 8 && e.Off == nil { e.Off = item; ok=true }  //wand",
	if item.ItemType == 15 && e.Ammo == nil && e.Main != nil { if e.Main.ItemType == 4 { e.Ammo = item; ok=true } }  //ammo",
	if item.ItemType == 16 && e.Off == nil { e.Off = item; ok=true }  //instrument",
	if item.ItemType == 17 && e.Off == nil { e.Off = item; ok=true }  //beverage",
	if item.ItemType == 19 && e.Feet == nil { e.Feet = item; ok=true }  //feet",
	if item.ItemType == 20 && e.Legs == nil { e.Legs = item; ok=true }  //legs",
	if item.ItemType == 21 && e.Arms == nil { e.Arms = item; ok=true }  //arms",
	if item.ItemType == 22 && e.Neck == nil { e.Neck = item; ok=true }  //neck",
	if item.ItemType == 23 && e.Off == nil { e.Off = item; ok=true }  //shield",
	if item.ItemType == 24 && (e.Ring1 == nil || e.Ring2 == nil) { if e.Ring1 == nil { e.Ring1 = item } else { e.Ring2 = item}; ok=true }  //finger",
	if item.ItemType == 25 && e.Head == nil { e.Head = item; ok=true }  //head",
	if item.ItemType == 0 && e.Main == nil { e.Main = item; ok=true }  //sharp",
	if item.ItemType == 1 && e.Main == nil { e.Main = item; ok=true }  //thrust",
	if item.ItemType == 2 && e.Main == nil { e.Main = item; ok=true }  //blunt",
	if item.ItemType == 3 && e.Main == nil { e.Main = item; ok=true }  //pole",
	if item.ItemType == 4 && e.Main == nil { e.Main = item; ok=true }  //range",

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

	if e.Head != nil { if strings.Contains(strings.ToLower(e.Head.Name), strings.ToLower(alias)) { item = e.Head; ok=true }}
	if e.Chest != nil { if strings.Contains(strings.ToLower(e.Chest.Name), strings.ToLower(alias)){  item = e.Chest; ok=true}}
	if e.Neck != nil { if strings.Contains(strings.ToLower(e.Neck.Name), strings.ToLower(alias)){  item = e.Neck; ok=true}}
	if e.Legs != nil { if strings.Contains(strings.ToLower(e.Legs.Name), strings.ToLower(alias)){  item = e.Legs; ok=true}}
	if e.Feet != nil { if strings.Contains(strings.ToLower(e.Feet.Name), strings.ToLower(alias)){  item = e.Feet; ok=true}}
	if e.Arms != nil { if strings.Contains(strings.ToLower(e.Arms.Name), strings.ToLower(alias)){  item = e.Arms; ok=true}}
	if e.Hands != nil { if strings.Contains(strings.ToLower(e.Hands.Name), strings.ToLower(alias)){  item = e.Hands; ok=true}}
	if e.Ring1 != nil { if strings.Contains(strings.ToLower(e.Ring1.Name), strings.ToLower(alias)){  item = e.Ring1; ok=true}}
	if e.Ring2 != nil { if strings.Contains(strings.ToLower(e.Ring2.Name), strings.ToLower(alias)){  item = e.Ring2; ok=true}}
	if e.Main != nil { if strings.Contains(strings.ToLower(e.Main.Name), strings.ToLower(alias)){  item = e.Main; ok=true}}
	if e.Off != nil { if strings.Contains(strings.ToLower(e.Off.Name), strings.ToLower(alias)){  item = e.Off; ok=true}}

	// Update armor values
	if ok && item != nil {
		e.Armor -= item.Armor
		e.Weight -= item.Weight
	}
	return ok, item
}

func (e *Equipment) Jsonify() string {
	itemList := make([]map[string]interface{}, 0)

	if e.Head != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Head)) }
	if e.Chest != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Chest)) }
	if e.Neck != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Neck)) }
	if e.Legs != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Legs)) }
	if e.Feet != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Feet)) }
	if e.Arms != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Arms)) }
	if e.Hands != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Hands)) }
	if e.Ring1 != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Ring1)) }
	if e.Ring2 != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Ring2)) }
	if e.Main != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Main)) }
	if e.Off != nil { itemList = append(itemList, ReturnItemInstanceProps(e.Off)) }

	data, err := json.Marshal(itemList)
	if err != nil {
		return "[]"
	} else {
		return string(data)
	}
}

func RestoreEquipment(jsonString string) *Equipment{
	var obj interface{}
	NewEquipment := &Equipment{}
	err := json.Unmarshal([]byte(jsonString), &obj)
	if err != nil {
		return NewEquipment
	}
	for _, item := range obj.([]map[string]interface{}) {
		newItem := Item{}
		copier.Copy(&newItem, Items[item["itemId"].(int)])
		newItem.Name = item["name"].(string)
		newItem.MaxUses	= item["uses"].(int)
		newItem.Flags["magic"] = item["magic"].(int) != 0
		newItem.Spell = item["spell"].(string)
		newItem.Armor =	item["armor"].(int)
		NewEquipment.Equip(&newItem)
	}
	return NewEquipment
}
