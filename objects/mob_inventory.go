package objects

import (
	"encoding/json"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/jinzhu/copier"
	"strconv"
	"strings"
	"sync"
)

type MobInventory struct {
	ParentId int
	Contents []*Mob
	sync.Mutex
	Flags map[string]bool
}

// NewMobInventory returns a new basic MobInventory structure
func NewMobInventory(ParentID int, o ...*Mob) *MobInventory {
	i := &MobInventory{
		ParentId: ParentID,
		Contents: make([]*Mob, 0, len(o)),
	}

	for _, ob := range o {
		i.Add(ob, false)
	}

	return i
}

// Add adds the specified Mob to the Contents.
func (i *MobInventory) Add(o *Mob, silent bool) {
	o.ParentId = i.ParentId
	i.Contents = append(i.Contents, o)
	if !silent {
		if o.Flags["invisible"] {
			Rooms[i.ParentId].MessageVisible(text.Magenta + "You encounter: " + o.Name + text.Reset + "\n")
		} else if !o.Flags["hidden"] {
			Rooms[i.ParentId].MessageAll(text.Magenta + "You encounter: " + o.Name + text.Reset + "\n")
		}
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
			mob.MobTickerUnload <- true
		} else {
			mob.MobTickerUnload <- true
			mob = nil
		}
	}
	i.Contents = newContents
}

func (i *MobInventory) RestartPerms() {
	for _, mob := range i.Contents {
		mob.StartTicking()
	}
}

// Search the MobInventory to return a specific instance of something
func (i *MobInventory) Search(alias string, num int, observer *Character) *Mob {
	if i == nil {
		return nil
	}

	pass := 1
	for _, c := range i.Contents {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)) {
			if pass == num {
				if c.Flags["hidden"] == false ||
					(c.Flags["hidden"] == true &&
						observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)){

					if c.Flags["invisible"] == false ||
						(c.Flags["invisible"] == true &&
							observer.Flags["detect_invisible"]) ||
						observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster){
						return c
					}
				}
			} else {
				pass++
			}
		}
	}

	return nil
}

// GetNumber Search the MobInventory to return a specific instance of something
func (i *MobInventory) GetNumber(o *Mob) int {
	pass := 1
	for _, c := range i.Contents {
		if c == o {
			return pass
		} else if c.Name == o.Name {
			pass++
		}
	}
	return pass
}

// List the items in this MobInventory
func (i *MobInventory) List(observer *Character) []string {
	items := make([]string, 0)

	for _, c := range i.Contents {
		if c.Flags["hidden"] == false ||
			(c.Flags["hidden"] == true &&
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)){

			if c.Flags["invisible"] == false ||
				(c.Flags["invisible"] == true &&
					observer.Flags["detect_invisible"]) ||
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster){
				items = append(items, c.Name)
			}
		}
	}
	return items
}

// List the items in this MobInventory
func (i *MobInventory) ListMobs(observer *Character) []*Mob {
	items := make([]*Mob, 0)

	for _, c := range i.Contents {
		if c.Flags["hidden"] == false ||
			(c.Flags["hidden"] == true &&
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)){

			if c.Flags["invisible"] == false ||
				(c.Flags["invisible"] == true &&
					observer.Flags["detect_invisible"]) ||
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster){
				items = append(items, c)
			}
		}
	}
	return items
}

// List the items in this MobInventory
func (i *MobInventory) ListAttackers(observer *Character) string {
	items := ""

	for _, o := range i.Contents {
		if o.CurrentTarget != "" {
			// List all
			if observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				items += o.Name + " #" + strconv.Itoa(i.GetNumber(o)) + " is attacking " + o.CurrentTarget + "!\n"
				// List non-hiddens invis
			} else if observer.Flags["detect_invisible"] {
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
func (i *MobInventory) ReducedList(observer *Character) string {
	items := make(map[string]int, 0)

	for _, c := range i.Contents {
		_, inMap := items[c.Name]
		if c.Flags["hidden"] == false ||
			(c.Flags["hidden"] == true &&
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)){

			if c.Flags["invisible"] == false ||
				(c.Flags["invisible"] == true &&
					observer.Flags["detect_invisible"]) ||
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster){
				if inMap {
					items[c.Name]++
				} else {
					items[c.Name] = 1
				}
			}
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

func (i *MobInventory) Jsonify() string {
	mobList := make([]map[string]interface{}, 0)

	if len(i.Contents) == 0 {
		return "[]"
	}

	for _, o := range i.Contents {
		if o.Flags["permanent"] {
			mobList = append(mobList, ReturnMobInstanceProps(o))
		}
	}

	data, err := json.Marshal(mobList)
	if err != nil {
		return "[]"
	} else {
		return string(data)
	}
}

func RestoreMobs(ParentID int, jsonString string) *MobInventory {
	NewInventory := &MobInventory{
		ParentId: ParentID,
		Contents: make([]*Mob, 0, 0),
	}
	obj := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(jsonString), &obj)
	if err != nil {
		return NewInventory
	}
	for _, mob := range obj {
		newMob := Mob{}
		copier.Copy(&newMob, Mobs[int(mob["mobId"].(float64))])
		newMob.Stam.Current = int(mob["health"].(float64))
		newMob.Mana.Current = int(mob["mana"].(float64))
		newMob.Placement = int(mob["placement"].(float64))
		newMob.Inventory = RestoreInventory(mob["inventory"].(string))
		NewInventory.Add(&newMob, true)

	}
	return NewInventory
}
