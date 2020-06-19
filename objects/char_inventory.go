package objects

import (
	"strings"
	"sync"
)

type CharInventory struct {
	ParentId int64
	Contents    []*Character
	sync.Mutex
	Flags map[string]bool
}


// New CharInventory returns a new basic CharInventory structure
func NewCharInventory(roomID int64, o ...*Character) *CharInventory {
	i := &CharInventory{
		ParentId: roomID,
		Contents:  make([]*Character, 0, len(o)),
	}

	for _, ob := range o {
		i.Add(ob)
	}

	return i
}

// Add adds the specified object to the contents.
func (i *CharInventory) Add(o *Character) {
	if len(i.Contents) == 0 {
		Rooms[i.ParentId].FirstPerson()
	}
	i.Contents = append(i.Contents, o)
}

// Pass character as a pointer, compare and remove
func (i *CharInventory) Remove(o *Character) {
	for c, p := range i.Contents {
		if p == o {
			copy(i.Contents[c:], i.Contents[c+1:])
			i.Contents[len(i.Contents)-1] = nil
			i.Contents = i.Contents[:len(i.Contents)-1]
			break
		}
	}
	if len(i.Contents) == 0{
		Rooms[i.ParentId].LastPerson()
	}
	if len(i.Contents) == 0 {
		i.Contents = make([]*Character, 0, 10)
	}
}

// Search the CharInventory to return a specific instance of something
func (i *CharInventory) Search(alias string, gm bool) *Character {
	if i == nil {
		return nil
	}

	for _, c := range i.Contents {
		if c.Flags["invisible"] == false || gm {
			if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)) {
				return c
			}
		}
	}

	return nil
}

// List the items in this CharInventory
func (i *CharInventory) List(seeInvisible bool, exclude string, gm bool) []string {
	// Determine how many items we need if this is an all request.. and we have only one entry.  Return nothing
	items := make([]string, 0)

	for _, o := range i.Contents {
		// List all
		if o.Name != exclude {
			if seeInvisible && gm{
				items = append(items, o.Name)
			// List non-hiddens
			} else if seeInvisible && !gm {
				if o.Flags["hidden"] != true {
					items = append(items, o.Name)
				}
			} else {
				if o.Flags["invisible"] != true && o.Flags["hidden"] != true {
					items = append(items, o.Name)
				}
			}
		}
	}

	return items
}



// Free recursively calls Free on all of it's content when the CharInventory
// attribute is freed.
func (i *CharInventory) Free() {
	if i == nil {
		return
	}
	for x, t := range i.Contents {
		i.Contents[x] = nil
		t.Free()
	}
}