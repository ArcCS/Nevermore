package cmd

import (
	"testing"
	"time"

	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
)

// A rooted player can't take an exit or change squares until the root ends.
func TestRootedPlayerCantMove(t *testing.T) {
	const from, to = -133, -134
	room := func(id int) *objects.Room {
		r := &objects.Room{
			RoomId: id,
			Flags:  map[string]bool{"active": true},
			Exits:  map[string]*objects.Exit{},
			Chars:  &objects.CharInventory{ParentId: id},
			Mobs:   &objects.MobInventory{ParentId: id},
			Items:  objects.NewItemInventory(),
		}
		objects.Rooms[id] = r
		t.Cleanup(func() { delete(objects.Rooms, id) })
		return r
	}
	start, end := room(from), room(to)
	end.Chars.Contents = []*objects.Character{{Object: objects.Object{Name: "Bystander"}, Flags: map[string]bool{}, ParentId: to}}
	start.Exits["gate"] = &objects.Exit{Object: objects.Object{Name: "gate", Placement: 3}, ParentId: from, ToId: to, Flags: map[string]bool{}}

	c := &objects.Character{
		Object:     objects.Object{Name: "Walker", Placement: 3},
		Permission: permissions.Anyone | permissions.Player,
		ParentId:   from,
		Flags:      map[string]bool{},
		Modifiers:  map[string]int{},
		Timers:     map[string]time.Time{},
		Hooks:      map[string]map[string]*objects.Hook{},
		Stam:       objects.Meter{Max: 10, Current: 10},
		Str:        objects.Meter{Max: 20, Current: 20},
		Inventory:  objects.NewItemInventory(),
		Equipment:  &objects.Equipment{},
	}
	// Someone stays behind, so leaving never empties the room and sets off
	// its last-person bookkeeping.
	start.Chars.Contents = []*objects.Character{c,
		{Object: objects.Object{Name: "Stayer"}, Flags: map[string]bool{}, ParentId: from}}

	c.Root(60)
	for _, input := range []string{"go gate", "f", "b", "move forward", "sprint forward"} {
		Parse(c, input)
		if c.ParentId != from || c.Placement != 3 {
			t.Errorf("%q moved a rooted player: room %d, square %d", input, c.ParentId, c.Placement)
		}
	}

	c.Timers["rooted"] = time.Now().Add(-time.Second)
	Parse(c, "f")
	if c.Placement != 4 {
		t.Errorf("square %d after the root ended, want 4", c.Placement)
	}
	c.Timers["global"] = time.Time{}
	c.Placement = 3
	Parse(c, "go gate")
	if c.ParentId != to {
		t.Errorf("room %d after the root ended, want %d", c.ParentId, to)
	}
}
