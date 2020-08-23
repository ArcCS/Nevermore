package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/text"
	"strconv"
	"strings"
	"sync"
)

type MobInventory struct {
	ParentId int
	Contents    []*Mob
	sync.Mutex
	Flags map[string]bool
}


// New MobInventory returns a new basic MobInventory structure
func NewMobInventory(ParentID int, o ...*Mob) *MobInventory {
	i := &MobInventory{
		ParentId: ParentID,
		Contents:  make([]*Mob, 0, len(o)),
	}

	for _, ob := range o {
		i.Add(ob)
	}

	return i
}

// Add adds the specified Mob to the Contents.
func (i *MobInventory) Add(o *Mob) {
	o.ParentId = i.ParentId
	i.Contents = append(i.Contents, o)
	if o.Flags["invisible"] {
		Rooms[i.ParentId].MessageVisible(text.Magenta + "You encounter: " + o.Name + text.Reset + "\n")
	}else if !o.Flags["hidden"] {
		Rooms[i.ParentId].MessageAll(text.Magenta + "You encounter: " + o.Name + text.Reset + "\n")
	}
}

// Pass mob as a pointer, compare and remove
func (i *MobInventory) Remove(o *Mob) {
	o.MobTickerUnload <- true
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

// Clear all non permanent
func (i *MobInventory) RemoveNonPerms() {
	newContents := make([]*Mob, 0, 0)
	for _, mob := range i.Contents {
		if mob.Flags["permanent"] == true {
			newContents = append(newContents, mob)
		}else{
			mob.MobTickerUnload <- true
			mob = nil
		}
	}
	i.Contents = newContents
}

// Search the MobInventory to return a specific instance of something
func (i *MobInventory) Search(alias string, num int, gm bool) *Mob {
	if i == nil {
		return nil
	}

	pass := 1
	for _, c := range i.Contents {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)){
			//log.Println("Searching for mob on pass " + strconv.Itoa(pass) + " looking for " + strconv.Itoa(num))
			if pass == num {
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

// Search the MobInventory to return a specific instance of something
func (i *MobInventory) GetNumber(o *Mob) int {
	pass := 1
	for _, c := range i.Contents {
		if c == o {
			return pass
		}else if c.Name == o.Name {
			pass++
		}
	}
	return pass
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

// List the items in this MobInventory
func (i *MobInventory) ListAttackers(seeInvisible bool, gm bool) string {
	items := ""

	for _, o := range i.Contents {
		if o.CurrentTarget != "" {
			// List all
			if seeInvisible && gm {
				items += o.Name + " #" + strconv.Itoa(i.GetNumber(o)) + " is attacking " + o.CurrentTarget + "!\n"
				// List non-hiddens invis
			} else if seeInvisible && !gm {
				if o.Flags["hidden"] != true {
					items += o.Name + " #" + strconv.Itoa(i.GetNumber(o)) + " is attacking " + o.CurrentTarget + "!\n"
				}
				// List non-hiddens
			} else {
				if o.Flags["invisible"] != true && o.Flags["hidden"] != true {
					items += o.Name + " #" + strconv.Itoa(i.GetNumber(o)) + " is attacking " + o.CurrentTarget + "!\n"
				}
			}
		}
	}

	return items
}

// List the items in this MobInventory
func (i *MobInventory) ReducedList(seeInvisible bool, gm bool) string {
	items := make(map[string]int, 0)

	for _, o := range i.Contents {
		// List all
		_, inMap := items[o.Name]
		if seeInvisible && gm{
			if inMap {
				items[o.Name]++
			}else {
				items[o.Name] = 1
			}
			// List non-hiddens invis
		}else if seeInvisible && !gm {
			if o.Flags["hidden"] != true {
				if inMap {
					items[o.Name]++
				}else {
					items[o.Name] = 1
				}
			}
			// List non-hiddens
		} else {
			if o.Flags["invisible"] != true && o.Flags["hidden"] != true {
				if inMap {
					items[o.Name]++
				}else {
					items[o.Name] = 1
				}
			}
		}
	}

	stringify := make([]string, 0)
	for k, v := range items {
		if v == 1 {
			stringify = append(stringify, "a "+k)
		}else {
			stringify = append(stringify, config.TextNumbers[v]+" "+k+"s")
		}
	}

	return strings.Join(stringify, ", ")
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