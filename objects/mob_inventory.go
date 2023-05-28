package objects

import (
	"encoding/json"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/jinzhu/copier"
	"log"
	"strconv"
	"strings"
)

type MobInventory struct {
	ParentId      int
	Contents      []*Mob
	Flags         map[string]bool
	ContinueEmpty func() bool
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

// Add adds the specified Mob to the Contents.
func (i *MobInventory) AddWithMessage(o *Mob, message string, silent bool) {
	o.ParentId = i.ParentId
	i.Contents = append(i.Contents, o)
	if !silent {
		if o.Flags["invisible"] {
			Rooms[i.ParentId].MessageVisible(text.Magenta + message + text.Reset + "\n")
		} else if !o.Flags["hidden"] {
			Rooms[i.ParentId].MessageAll(text.Magenta + message + text.Reset + "\n")
		}
	}
}

// Pass mob as a pointer, compare and remove
func (i *MobInventory) Remove(o *Mob) {
	log.Println("Unloading mob from inventory: " + o.Name)
	go func() { o.MobTickerUnload <- true }()
	for c, p := range i.Contents {
		if p == o {
			copy(i.Contents[c:], i.Contents[c+1:])
			i.Contents[len(i.Contents)-1] = nil
			i.Contents = i.Contents[:len(i.Contents)-1]
		}
		p = nil
	}
	if len(i.Contents) == 0 {
		i.Contents = nil
		i.Contents = make([]*Mob, 0, 10)
	}
}

// Clear all non permanent
func (i *MobInventory) RemoveNonPerms() {
	var contentRef []*Mob
	for _, mob := range i.Contents {
		if mob.Flags["permanent"] != true {
			contentRef = append(contentRef, mob)
		} else {
			log.Println("Unload mob: " + mob.Name + " ticker, but do not delete")
			mob.MobTickerUnload <- true
		}
	}
	// Check if we should continue to empty, this is only relevant if mobs have been thinking and we have to back out of this loop entirely
	for i.ContinueEmpty() && len(contentRef) > 0 {
		for index, mob := range contentRef {
			if !mob.IsThinking {
				i.Remove(mob)
				mob = nil
			} else {
				contentRef = contentRef[index:]
				break
			}
		}
		contentRef = nil
	}
}

func (i *MobInventory) RestartPerms() {
	for _, mob := range i.Contents {
		if mob.Flags["permanent"] && !mob.IsActive {
			mob.StartTicking()
		}
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
						observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) {

					if c.Flags["invisible"] == false ||
						(c.Flags["invisible"] == true &&
							observer.Flags["detect_invisible"]) ||
						observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
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
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) {

			if c.Flags["invisible"] == false ||
				(c.Flags["invisible"] == true &&
					observer.Flags["detect_invisible"]) ||
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
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
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) {

			if c.Flags["invisible"] == false ||
				(c.Flags["invisible"] == true &&
					observer.Flags["detect_invisible"]) ||
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				items = append(items, c)
			}
		}
	}
	return items
}

// ListChars the items in this CharInventory
func (i *MobInventory) ListHiddenMobs(observer *Character) []*Mob {
	// Determine how many items we need if this is an all request.. and we have only one entry.  Return nothing
	items := make([]*Mob, 0)

	for _, m := range i.Contents {
		// List all
		if m.CheckFlag("hidden") {
			if !m.CheckFlag("invisible") || observer.CheckFlag("detect-invisible") {
				items = append(items, m)
			}
		}
	}
	return items
}

// List the items in this MobInventory
func (i *MobInventory) ListAttackers(observer *Character) string {
	items := ""
	victim := ""

	for _, o := range i.Contents {
		if o.CurrentTarget != "" {
			victim = o.CurrentTarget
			if o.CurrentTarget == observer.Name {
				victim = text.Bold + "you" + text.Reset
			}
			// List all
			if observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				items += o.Name + " #" + strconv.Itoa(i.GetNumber(o)) + " is attacking " + victim + "!\n"
				// List non-hiddens invis
			} else if observer.Flags["detect_invisible"] {
				if o.Flags["hidden"] != true {
					items += o.Name + " #" + strconv.Itoa(i.GetNumber(o)) + " is attacking " + victim + "!\n"
				}
				// List non-hiddens
			} else {
				if o.Flags["invisible"] != true && o.Flags["hidden"] != true {
					items += o.Name + " #" + strconv.Itoa(i.GetNumber(o)) + " is attacking " + victim + "!\n"
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
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) {

			if c.Flags["invisible"] == false ||
				(c.Flags["invisible"] == true &&
					observer.Flags["detect_invisible"]) ||
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
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
		copier.CopyWithOption(&newMob, Mobs[int(mob["mobId"].(float64))], copier.Option{DeepCopy: true})
		newMob.Placement = int(mob["placement"].(float64))
		newMob.Inventory = RestoreInventory(mob["inventory"].(string))
		NewInventory.Add(&newMob, true)

	}
	return NewInventory
}
