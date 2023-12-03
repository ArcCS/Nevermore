package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(mobWander{},
		"",
		permissions.None,
		"$MWANDER")
}

type mobWander cmd

func (mobWander) process(s *state) {
	/*
			if m.Flags["takes_treasure"] {
			// Roll to see if mob is picking it up
			if utils.Roll(100, 1, 0) <= config.MobTakeChance {
				// Loop inventory, and take the first thing they find
				for _, item := range Rooms[m.ParentId].Items.Contents {
					if m.Placement == item.Placement && !item.Flags["hidden"] {
						if err := Rooms[m.ParentId].Items.Remove(item); err != nil {
							log.Println("Error mob removing item", err)
						}
						m.Inventory.Add(item)
						Rooms[m.ParentId].MessageAll(m.Name + " picks up " + item.DisplayName() + text.Reset + "\n")
						break
					}
				}
			}
		}
	*/
	return
}
