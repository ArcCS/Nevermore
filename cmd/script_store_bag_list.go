package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/jedib0t/go-pretty/table"
	"strconv"
)

func init() {
	addHandler(listbag{},
		"",
		permissions.Player,
		"$BAGLIST")
}

var baseBag = 15500
var bagCapacity = 1500
var weightLess = 150000
var bagWeight = map[int]int{
	19: 1000,
	18: 2000,
	17: 3000,
	16: 4000,
	15: 5000,
	14: 8000,
	13: 12000,
	12: 16000,
	11: 18000,
	10: 20000,
	9:  30000,
	8:  45000,
	7:  75000,
	6:  95000,
	5:  125000,
}
var bagWeightOList = []int{
	20,
	19,
	18,
	17,
	16,
	15,
	14,
	13,
	12,
	11,
	10,
	9,
	8,
	7,
	6,
	5,
}

type listbag cmd

func (listbag) process(s *state) {
	if len(s.where.StoreInventory.Contents) != 0 {
		rowLength := 120
		t := table.NewWriter()
		t.SetAllowedRowLength(rowLength)
		t.Style().Options.SeparateRows = false
		t.AppendHeader(table.Row{"Bag Modification", "Cost"})
		t.AppendRows([]table.Row{
			{"Base Bag", strconv.Itoa(baseBag)}})
		t.AppendRows([]table.Row{
			{"Weightless Holding", strconv.Itoa(weightLess)}})
		t.AppendRows([]table.Row{
			{"Per Item Capacity", bagCapacity}})
		t.AppendRows([]table.Row{
			{"Weights:", ""}})
		for _, bagW := range bagWeightOList {
			t.AppendRows([]table.Row{
				{strconv.Itoa(bagW) + " Weight ", strconv.Itoa(bagWeight[bagW])}})
		}

		s.msg.Actor.SendGood(t.Render())
	} else {
		s.msg.Actor.SendInfo("There is nothing for sale here.")
	}

}
