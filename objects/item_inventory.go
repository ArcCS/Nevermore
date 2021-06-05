package objects

import (
	"encoding/json"
	"github.com/ArcCS/Nevermore/config"
	"github.com/jinzhu/copier"
	"log"
	"strings"
	"sync"
)

type ItemInventory struct {
	Contents []*Item
	sync.Mutex
	TotalWeight int
	Flags       map[string]bool
}

// New ItemInventory returns a new basic ItemInventory structure
func NewItemInventory(o ...*Item) *ItemInventory {
	i := &ItemInventory{
		Contents: make([]*Item, 0, len(o)),
	}

	for _, ob := range o {
		i.Add(ob)
	}

	return i
}

// Add adds the specified object to the Contents.
func (i *ItemInventory) Add(o *Item) {
	i.Contents = append(i.Contents, o)
	i.TotalWeight += o.GetWeight()
}

// Pass item as a pointer to be removed
func (i *ItemInventory) Remove(o *Item) (err error) {
	defer func() (err error) {
		if r := recover(); r != nil {
			log.Println("Item inventory removal recovery, failed to process", r)
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
	i.TotalWeight -= o.GetWeight()
	return nil
}

// Clear all non permanent
func (i *ItemInventory) RemoveNonPerms() {
	newContents := make([]*Item, 0, 0)
	newWeight := 0
	for _, item := range i.Contents {
		if item.Flags["permanent"] == true {
			newContents = append(newContents, item)
			newWeight += item.GetWeight()
		} else {
			item = nil
		}
	}
	i.Contents = newContents
	i.TotalWeight = newWeight
}

// Search the ItemInventory to return a specific instance of something
func (i *ItemInventory) Search(alias string, num int) *Item {
	if i == nil {
		return nil
	}

	pass := 1
	for _, c := range i.Contents {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)) {
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
		if strings.TrimSpace(o.Name) != "" {
			items = append(items, o.Name)
		}
	}

	return items
}

// Free recursively calls Free on all of it's content when the ItemInventory
// attribute is freed.
func (i *ItemInventory) Free() {
	if i == nil {
		return
	}
	for x, t := range i.Contents {
		i.Contents[x] = nil
		t.Free()
	}
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

// List the items in this inentory
func (i *ItemInventory) ReducedList() string {
	items := make(map[string]int, 0)

	for _, o := range i.Contents {
		// List all
		_, inMap := items[o.Name]
		if inMap {
			items[o.Name]++
		} else {
			items[o.Name] = 1
		}
	}

	stringify := make([]string, 0)
	for k, v := range items {
		if v == 1 {
			stringify = append(stringify, "a "+k)
		} else {
			stringify = append(stringify, config.TextNumbers[v]+" "+k+"s")
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
		err = copier.Copy(&newItem, Items[int(item["itemId"].(float64))])
		if err == nil {
			newItem.Name = item["name"].(string)
			newItem.MaxUses = int(item["uses"].(float64))
			newItem.Flags["magic"] = int(item["magic"].(float64)) != 0
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
