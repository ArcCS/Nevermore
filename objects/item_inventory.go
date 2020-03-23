package objects

import (
	"strings"
	"sync"
)

type ItemInventory struct {
	Contents    []*Item
	sync.Mutex
	TotalWeight int64
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

// Search the ItemInventory to return a specific instance of something
func (i *ItemInventory) Search(alias string, num int) *Item {
	if i == nil {
		return nil
	}

	pass := 1
	for _, c := range i.Contents {
		if strings.Contains(c.Name, alias){
			if pass == num {
				return c
			}else{
				pass++
			}
		}
	}

	return nil
}

func (i *ItemInventory) Serialize() map[string]interface{}{
	return map[string]interface{}{}
}

// List the items in this ItemInventory
func (i *ItemInventory) List() []string {
	items := make([]string, len(i.Contents))

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