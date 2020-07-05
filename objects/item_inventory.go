package objects

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"strings"
	"sync"
)

type ItemInventory struct {
	Contents    []*Item
	sync.Mutex
	TotalWeight int
	Flags map[string]bool
}


// New ItemInventory returns a new basic ItemInventory structure
func NewItemInventory(o ...*Item) *ItemInventory {
	i := &ItemInventory{
		Contents:  make([]*Item, 0, len(o)),
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

// Pass character as a pointer, compare and remove
func (i *ItemInventory) Remove(o *Item) {
	for c, p := range i.Contents {
		if p == o {
			copy(i.Contents[c:], i.Contents[c+1:])
			i.Contents[len(i.Contents)-1] = nil
			i.Contents = i.Contents[:len(i.Contents)-1]
			break
		}
	}
	if len(i.Contents) == 0 {
		i.Contents = make([]*Item, 0, 10)
	}
	i.TotalWeight -= o.GetWeight()
}

// Clear all non permanent
func (i *ItemInventory) RemoveNonPerms() {
	newContents := make([]*Item, 0, 0)
	newWeight := 0
	for _, item := range i.Contents {
		if item.Flags["permanent"] == true {
			newContents = append(newContents, item)
			newWeight += item.GetWeight()
		}else{
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
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)){
			if pass == num {
				return c
			}else{
				pass++
			}
		}
	}

	return nil
}

// List the items in this ItemInventory
func (i *ItemInventory) List() []string {
	items := make([]string, 0)

	switch len(i.Contents){
	case 0:
		return items
	}

	for _, o := range i.Contents {
		items = append(items, o.Name)
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

	switch len(i.Contents){
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

func RestoreInventory(jsonString string) *ItemInventory {
	obj := make([]map[string]interface{}, 0)
	NewInventory := &ItemInventory{}
	err := json.Unmarshal([]byte(jsonString), &obj)
	if err != nil {
		return NewInventory
	}
	for _, item := range obj {
		newItem := Item{}
		copier.Copy(&newItem, Items[int(item["itemId"].(float64))])
		newItem.Name = item["name"].(string)
		newItem.MaxUses	= int(item["uses"].(float64))
		newItem.Flags["magic"] = int(item["magic"].(float64)) != 0
		newItem.Spell = item["spell"].(string)
		newItem.Armor =	int(item["armor"].(float64))
		if newItem.ItemType == 9 {
			newItem.Storage = RestoreInventory(item["contents"].(string))
		}
		NewInventory.Add(&newItem)
	}
	return NewInventory
}