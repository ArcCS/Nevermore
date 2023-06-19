package objects

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

type CharInventory struct {
	ParentId int
	Contents []*Character
	Flags    map[string]bool
}

// New CharInventory returns a new basic CharInventory structure
func NewCharInventory(roomID int, o ...*Character) *CharInventory {
	i := &CharInventory{
		ParentId: roomID,
		Contents: make([]*Character, 0, len(o)),
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
	if len(i.Contents) == 0 {
		go Rooms[i.ParentId].LastPerson()
	}
	if len(i.Contents) == 0 {
		i.Contents = make([]*Character, 0, 10)
	}
}

// Search the CharInventory to return a specific instance of something
func (i *CharInventory) SearchAll(alias string) *Character {
	if i == nil {
		return nil
	}

	for _, c := range i.Contents {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)) {
			return c
		}
	}

	return nil
}

// Search the CharInventory to return a specific instance of something
func (i *CharInventory) Search(alias string, observer *Character) *Character {
	if i == nil {
		return nil
	}

	for _, c := range i.Contents {
		if c.Flags["invisible"] == false ||
			(c.Flags["invisible"] == true &&
				observer.Flags["detect_invisible"] &&
				!c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) ||
			observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)) {
				return c
			}
		}
	}

	return nil
}

// Search the CharInventory to return a specific instance of something
func (i *CharInventory) MobSearch(alias string, observer *Mob) *Character {
	if i == nil {
		return nil
	}

	for _, c := range i.Contents {
		if c.Flags["invisible"] == false ||
			(c.Flags["invisible"] == true &&
				observer.Flags["detect_invisible"] &&
				!c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) {
			if strings.Contains(strings.ToLower(c.Name), strings.ToLower(alias)) {
				return c
			}
		}
	}

	return nil
}

// List the items in this CharInventory
func (i *CharInventory) List(observer *Character) []string {
	// Determine how many items we need if this is an all request.. and we have only one entry.  Return nothing
	items := make([]string, 0)

	for _, c := range i.Contents {
		// List all
		if strings.ToLower(c.Name) != strings.ToLower(observer.Name) {
			if c.Flags["hidden"] == false ||
				(c.Flags["hidden"] == true &&
					observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) {

				if c.Flags["invisible"] == false ||
					(c.Flags["invisible"] == true &&
						observer.Flags["detect_invisible"] &&
						!c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) ||
					observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					if c.CheckFlag("singing") {
						items = append(items, c.Name+" (singing)")
					} else {
						items = append(items, c.Name)

					}
				}
			}
		}
	}
	return items
}

// ListChars the items in this CharInventory
func (i *CharInventory) ListChars(observer *Character) []*Character {
	// Determine how many items we need if this is an all request.. and we have only one entry.  Return nothing
	items := make([]*Character, 0)

	for _, c := range i.Contents {
		// List all
		if strings.ToLower(c.Name) != strings.ToLower(observer.Name) {
			if c.Flags["hidden"] == false ||
				(c.Flags["hidden"] == true &&
					observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) {

				if c.Flags["invisible"] == false ||
					(c.Flags["invisible"] == true &&
						observer.Flags["detect_invisible"] &&
						!c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) ||
					observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					items = append(items, c)
				}
			}
		}
	}
	return items
}

// ListChars the items in this CharInventory
func (i *CharInventory) ListHiddenChars(observer *Character) []*Character {
	// Determine how many items we need if this is an all request.. and we have only one entry.  Return nothing
	items := make([]*Character, 0)

	for _, c := range i.Contents {
		// List all
		if strings.ToLower(c.Name) != strings.ToLower(observer.Name) {
			if c.CheckFlag("hidden") && !c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				if !c.CheckFlag("invisible") || observer.CheckFlag("detect-invisible") {
					items = append(items, c)
				}
			}
		}
	}
	return items
}

// MobList lists characters for a mobs point of view
func (i *CharInventory) MobList(observer *Mob) []string {
	// Determine how many items we need if this is an all request.. and we have only one entry.  Return nothing
	items := make([]string, 0)

	// List all
	for _, c := range i.Contents {
		if c.Flags["hidden"] == false {
			if c.Flags["invisible"] == false ||
				(c.Flags["invisible"] == true &&
					observer.Flags["detect_invisible"]) {
				items = append(items, c.Name)
			}
		}
	}
	return items
}
