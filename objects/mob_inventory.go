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
	JsonRepr      string
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

// AddWithMessage Add adds the specified Mob to the Contents.
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

// Remove Pass mob as a pointer, compare and remove
func (i *MobInventory) Remove(o *Mob) {
	go func() {
		o.MobTickerUnload <- true
		close(o.MobCommands)
	}()
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

// RemoveNonPerms Clear all non-permanent
func (i *MobInventory) RemoveNonPerms() {
	var contentRef []*Mob
	for _, mob := range i.Contents {
		if mob.Flags["permanent"] != true {
			contentRef = append(contentRef, mob)
		} else {
			mob.MobTickerUnload <- true
			close(mob.MobCommands)
		}
	}
	// Check if we should continue to empty, this is only relevant if mobs have been thinking, and we have to back out of this loop entirely
	//Check if we should empty the room
	for i.ContinueEmpty() && len(contentRef) > 0 {
		for index, mob := range contentRef {
			if !mob.IsThinking {
				i.Remove(mob)
				mob = nil
				if len(contentRef) == index+1 {
					contentRef = nil
				}
			} else {
				contentRef = contentRef[index:]
				break
			}
		}
	}
	i.Jsonify()
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
							observer.Flags["detect-invisible"]) ||
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
					observer.Flags["detect-invisible"]) ||
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				items = append(items, c.Name)
			}
		}
	}
	return items
}

// ListHostile List the items in this MobInventory
func (i *MobInventory) ListHostile() []string {
	items := make([]string, 0)

	for _, c := range i.Contents {
		if c.CheckFlag("hostile") {
			items = append(items, c.Name)
		}
	}
	return items
}

// ListMobs List the items in this MobInventory
func (i *MobInventory) ListMobs(observer *Character) []*Mob {
	items := make([]*Mob, 0)

	for _, c := range i.Contents {
		if c.Flags["hidden"] == false ||
			(c.Flags["hidden"] == true &&
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) {

			if c.Flags["invisible"] == false ||
				(c.Flags["invisible"] == true &&
					observer.Flags["detect-invisible"]) ||
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				items = append(items, c)
			}
		}
	}
	return items
}

// ListHiddenMobs ListChars the items in this CharInventory
func (i *MobInventory) ListHiddenMobs(observer *Character) []*Mob {
	// Determine how many items we need if this is an all request. and we have only one entry.  Return nothing
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

// ListAttackers List the items in this MobInventory
func (i *MobInventory) ListAttackers(observer *Character) string {
	items := ""
	victim := ""

	for _, o := range i.Contents {
		if o.CurrentTarget != "" {
			victim = o.CurrentTarget
			if o.CurrentTarget == observer.Name {
				victim = text.Bold + "you" + text.Reset + text.Info
			}
			// List all
			if observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				items += o.Name + " #" + strconv.Itoa(i.GetNumber(o)) + " is attacking " + victim + "!\n"
				// List non-hiddens invis
			} else if observer.Flags["detect-invisible"] {
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

// ReducedList List the items in this MobInventory
func (i *MobInventory) ReducedList(observer *Character) string {
	items := make(map[string]int)

	for _, c := range i.Contents {
		_, inMap := items[c.Name]
		if c.Flags["hidden"] == false ||
			(c.Flags["hidden"] == true &&
				observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) {

			if c.Flags["invisible"] == false ||
				(c.Flags["invisible"] == true &&
					observer.Flags["detect-invisible"]) ||
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

func (i *MobInventory) Jsonify() {
	mobList := make([]map[string]interface{}, 0)

	if len(i.Contents) == 0 {
		i.JsonRepr = "[]"
		return
	}

	for _, o := range i.Contents {
		if isPerm, ok := o.Flags["permanent"]; ok && isPerm {
			mobList = append(mobList, ReturnMobInstanceProps(o))
		}
	}

	data, err := json.Marshal(mobList)
	if err != nil {
		i.JsonRepr = "[]"
	} else {
		i.JsonRepr = string(data)
	}
	return

}

func RestoreMobs(ParentID int, jsonString string) *MobInventory {
	NewInventory := &MobInventory{
		ParentId: ParentID,
		Contents: make([]*Mob, 0),
	}
	obj := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(jsonString), &obj)
	if err != nil {
		return NewInventory
	}
	for _, mob := range obj {
		newMob := Mob{}
		if err := copier.CopyWithOption(&newMob, Mobs[int(mob["mobId"].(float64))], copier.Option{DeepCopy: true}); err != nil {
			log.Println("Error copying mob during restore: ", err)
		}
		newMob.Placement = int(mob["placement"].(float64))
		newMob.Inventory = RestoreInventory(mob["inventory"].(string))
		NewInventory.Add(&newMob, true)

	}
	return NewInventory
}
