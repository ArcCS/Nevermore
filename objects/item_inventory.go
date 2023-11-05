package objects

import (
	"encoding/json"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jinzhu/copier"
	"sort"
	"strings"
)

type ItemInventory struct {
	Contents []*Item
	Flags    map[string]bool
}

// NewItemInventory New ItemInventory returns a new basic ItemInventory structure
func NewItemInventory(o ...*Item) *ItemInventory {
	i := &ItemInventory{
		Contents: make([]*Item, 0, len(o)),
	}

	for _, ob := range o {
		i.Add(ob)
	}

	return i
}

func (i *ItemInventory) GetTotalWeight() (total int) {
	for _, item := range i.Contents {
		total += item.GetWeight()
	}
	return total
}

// Add adds the specified object to the Contents.
func (i *ItemInventory) Add(o *Item) {
	index := sort.Search(len(i.Contents), func(num int) bool { return i.Contents[num].DisplayName() >= o.DisplayName() })
	if index < len(i.Contents) && i.Contents[index].DisplayName() == o.DisplayName() {
		index++
	}
	i.Contents = append(i.Contents[:index], append([]*Item{o}, i.Contents[index:]...)...)
}

// Remove Pass item as a pointer to be removed
func (i *ItemInventory) Remove(o *Item) (err error) {
	defer func() (err error) {
		if r := recover(); r != nil {
			return r.(error)
		}
		return nil
	}()
	for c, p := range i.Contents {
		if p == o {
			copy(i.Contents[c:], i.Contents[c+1:])
			i.Contents[len(i.Contents)-1] = nil
			i.Contents = i.Contents[:len(i.Contents)-1]
			break
		}
	}
	if len(i.Contents) == 0 {
		i.Contents = make([]*Item, 0, 0)
	}
	return nil
}

// RemoveNonPerms Clear all non-permanent
func (i *ItemInventory) RemoveNonPerms() {
	var newItems []*Item
	for _, item := range i.Contents {
		//log.Println("Checking item storage", strconv.Itoa(len(item.Storage.List())))
		if (strings.Contains(strings.ToLower(item.Name), "corpse of") && len(item.Storage.List()) != 0) ||
			(item.Flags["permanent"] && !strings.Contains(strings.ToLower(item.Name), "corpse of")) {
			continue
		} else {
			newItems = append(newItems, item)
		}
	}
	for _, item := range newItems {
		i.Remove(item)
	}
}

// Search the ItemInventory to return a specific instance of something
func (i *ItemInventory) Search(alias string, num int) *Item {
	if i == nil || alias == "" {
		return nil
	}

	pass := 1
	for _, c := range i.Contents {
		if strings.Contains(strings.ToLower(c.DisplayName()), strings.ToLower(alias)) {
			if pass == num {
				return c
			} else {
				pass++
			}
		}
	}

	return nil
}

// List the items in this ItemInventory
func (i *ItemInventory) List() []string {
	items := make([]string, 0)

	switch len(i.Contents) {
	case 0:
		return items
	}

	for _, o := range i.Contents {
		if strings.TrimSpace(o.DisplayName()) != "" {
			items = append(items, o.DisplayName())
		}
	}

	return items
}

// ListChars the items in this CharInventory
func (i *ItemInventory) ListHiddenItems(observer *Character) []*Item {
	// Determine how many items we need if this is an all request.. and we have only one entry.  Return nothing
	items := make([]*Item, 0)

	for _, item := range i.Contents {
		// List all
		if item.Flags["hidden"] {
			items = append(items, item)
		}
	}
	return items
}

// ListItems List the items in this ItemInventory
func (i *ItemInventory) ListItems() []*Item {
	items := make([]*Item, 0)

	switch len(i.Contents) {
	case 0:
		return items
	}

	for _, o := range i.Contents {
		if strings.TrimSpace(o.DisplayName()) != "" {
			items = append(items, o)
		}
	}

	return items
}

func (i *ItemInventory) Jsonify() string {
	itemList := make([]map[string]interface{}, 0)

	switch len(i.Contents) {
	case 0:
		return "[]"
	}

	for _, o := range i.Contents {
		itemList = append(itemList, ReturnItemInstanceProps(o))
	}

	data, err := json.Marshal(itemList)
	if err != nil {
		return "[]"
	} else {
		return string(data)
	}
}

// PermanentReducedList the items in this inventory
func (i *ItemInventory) PermanentReducedList() string {
	items := make(map[string]int, 0)

	for _, o := range i.Contents {
		if o.Flags["permanent"] {
			// List all
			_, inMap := items[o.DisplayName()]
			if inMap {
				items[o.DisplayName()]++
			} else {
				items[o.DisplayName()] = 1
			}
		}
	}

	stringify := make([]string, 0)
	for _, v := range i.Contents {
		if v.Flags["permanent"] {
			if items[v.DisplayName()] == 1 {
				stringify = append(stringify, "a "+v.DisplayName())
			} else {
				if !utils.StringIn(config.TextNumbers[items[v.DisplayName()]]+" "+v.DisplayName()+"s", stringify) {
					stringify = append(stringify, config.TextNumbers[items[v.DisplayName()]]+" "+v.DisplayName()+"s")
				}
			}
		}
	}

	return strings.Join(stringify, ", ")
}

// RoomReducedList the items in this inventory
func (i *ItemInventory) RoomReducedList() string {
	items := make(map[string]int, 0)

	for _, o := range i.Contents {
		if !o.Flags["permanent"] {
			// List all
			_, inMap := items[o.DisplayName()]
			if inMap {
				items[o.DisplayName()]++
			} else {
				items[o.DisplayName()] = 1
			}
		}
	}

	stringify := make([]string, 0)
	for _, v := range i.Contents {
		if !v.Flags["permanent"] {
			if items[v.DisplayName()] == 1 {
				stringify = append(stringify, "a "+v.DisplayName())
			} else {
				if !utils.StringIn(config.TextNumbers[items[v.DisplayName()]]+" "+v.DisplayName()+"s", stringify) {
					stringify = append(stringify, config.TextNumbers[items[v.DisplayName()]]+" "+v.DisplayName()+"s")
				}
			}
		}
	}

	return strings.Join(stringify, ", ")
}

// ReducedList the items in this inventory
func (i *ItemInventory) ReducedList() string {
	items := make(map[string]int, 0)

	for _, o := range i.Contents {
		// List all
		_, inMap := items[o.DisplayName()]
		if inMap {
			items[o.DisplayName()]++
		} else {
			items[o.DisplayName()] = 1
		}
	}

	stringify := make([]string, 0)
	for _, v := range i.Contents {
		if items[v.DisplayName()] == 1 {
			stringify = append(stringify, "a "+v.DisplayName())
		} else {
			if !utils.StringIn(config.TextNumbers[items[v.DisplayName()]]+" "+v.DisplayName()+"s", stringify) {
				stringify = append(stringify, config.TextNumbers[items[v.DisplayName()]]+" "+v.DisplayName()+"s")
			}
		}
	}

	return strings.Join(stringify, ", ")
}

func RestoreInventory(jsonString string) *ItemInventory {
	obj := make([]map[string]interface{}, 0)
	NewInventory := &ItemInventory{}
	err := json.Unmarshal([]byte(jsonString), &obj)
	if err != nil {
		return NewInventory
	}
	for _, item := range obj {
		newItem := Item{}
		err = copier.CopyWithOption(&newItem, Items[int(item["itemId"].(float64))], copier.Option{DeepCopy: true})
		if err == nil {
			newItem.Name = item["name"].(string)
			newItem.MaxUses = int(item["uses"].(float64))
			newItem.Flags["magic"] = int(item["magic"].(float64)) != 0
			if _, ok := item["light"]; ok {
				newItem.Flags["light"] = int(item["light"].(float64)) != 0
			}
			if _, ok := item["adjustment"]; ok {
				newItem.Adjustment = int(item["adjustment"].(float64))
			}
			if _, ok := item["infinite"]; ok {
				newItem.Flags["infinite"] = int(item["infinite"].(float64)) != 0
			}
			if _, ok := item["store_price"]; ok {
				newItem.StorePrice = int(item["store_price"].(float64))
			}
			newItem.Spell = item["spell"].(string)
			newItem.Armor = int(item["armor"].(float64))
			if newItem.ItemType == 9 {
				newItem.Storage = RestoreInventory(item["contents"].(string))
			}
			NewInventory.Add(&newItem)
		}
	}
	return NewInventory
}
