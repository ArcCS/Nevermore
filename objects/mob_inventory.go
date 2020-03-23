package objects

import (
	"strings"
	"sync"
)

type MobInventory struct {
	Contents    []*Mob
	sync.Mutex
	Flags map[string]bool
}


// New MobInventory returns a new basic MobInventory structure
func NewMobInventory(o ...*Mob) *MobInventory {
	i := &MobInventory{
		Contents:  make([]*Mob, 0, len(o)),
	}

	for _, ob := range o {
		i.Add(ob)
	}

	return i
}

// Add adds the specified Mob to the Contents.
func (i *MobInventory) Add(o *Mob) {
	//log.Println("Adding mob" + o.Name)
	i.Contents = append(i.Contents, o)
}

// Pass character as a pointer, compare and remove
func (i *MobInventory) Remove(o *Mob) {
	for c, p := range i.Contents {
		if p == o {
			copy(i.Contents[c:], i.Contents[c+1:])
			i.Contents[len(i.Contents)-1] = nil
			i.Contents = i.Contents[:len(i.Contents)-1]
			break
		}
	}
	if len(i.Contents) == 0 {
		i.Contents = make([]*Mob, 0, 10)
	}
}

// Search the MobInventory to return a specific instance of something
func (i *MobInventory) Search(alias string, num int64, gm bool) *Mob {
	if i == nil {
		return nil
	}

	pass := 1
	for _, c := range i.Contents {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)){
			if pass == int(num) {
				if i.Flags["invisible"] == false || gm {
					return c
				}
			}else{
				pass++
			}
		}
	}

	return nil
}

// List the items in this MobInventory
func (i *MobInventory) List(seeInvisible bool, gm bool) []string {
	items := make([]string, 0)

	for _, o := range i.Contents {
		// List all
		if seeInvisible && gm{
				items = append(items, o.Name)
			// List non-hiddens invis
		}else if seeInvisible && !gm {
			if o.Flags["hidden"] != true {
				items = append(items, o.Name)
			}
			// List non-hiddens
		} else {
			if o.Flags["invisible"] != true && o.Flags["hidden"] != true {
				items = append(items, o.Name)
			}
		}
	}

	return items
}



// Free recursively calls Free on all of it's content when the MobInventory
// attribute is freed.
func (i *MobInventory) Free() {
	if i == nil {
		return
	}
	for x, t := range i.Contents {
		i.Contents[x] = nil
		t.Free()
	}
}