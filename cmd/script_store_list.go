package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
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
			t.AppendRows([]table.Row{
				{item.Name, item.StorePrice, config.ItemTypes[item.ItemType]},
			})
		}
		s.msg.Actor.SendGood(t.Render())
	}else{
		s.msg.Actor.SendInfo("There is nothing for sale here.")
	}

}