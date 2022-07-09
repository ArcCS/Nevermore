package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jedib0t/go-pretty/table"
)

func init() {
	addHandler(list{},
		"",
		permissions.Player,
		"$LIST")
}

type list cmd

func (list) process(s *state) {
	if len(s.where.StoreInventory.Contents) != 0 {
		rowLength := 120
		t := table.NewWriter()
		t.SetAllowedRowLength(rowLength)
		t.Style().Options.SeparateRows = false
		t.AppendHeader(table.Row{"Item", "Cost", "Type"})
		for _, item := range s.where.StoreInventory.Contents {
			if utils.IntIn(item.ItemType, []int{5, 19, 20, 21, 22, 23, 24, 25, 26}) {
				t.AppendRows([]table.Row{
					{item.DisplayName(), item.StorePrice, config.ItemTypes[item.ItemType]},
				})
			} else if utils.IntIn(item.ItemType, []int{0, 1, 2, 3, 4}) {
				t.AppendRows([]table.Row{
					{item.DisplayName(), item.StorePrice, config.ItemTypes[item.ItemType]},
				})
			} else if item.ItemType == 6 {
				t.AppendRows([]table.Row{
					{item.Name, item.StorePrice, "(D) " + objects.Spells[item.Spell].Name},
				})
			} else if item.ItemType == 8 {
				t.AppendRows([]table.Row{
					{item.Name, item.StorePrice, "(W) " + objects.Spells[item.Spell].Name},
				})

			} else {
				t.AppendRows([]table.Row{
					{item.Name, item.StorePrice, config.ItemTypes[item.ItemType]},
				})
			}
		}
		s.msg.Actor.SendGood(t.Render())
	} else {
		s.msg.Actor.SendInfo("There is nothing for sale here.")
	}

}
