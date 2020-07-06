package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/jinzhu/copier"
	"strings"
)

func init() {
	addHandler(create{},
	"Usage:  create name \n \n Create a brand new item with a name. \n Note:  Use the modify command to add modify traits of the object.",
	permissions.Builder,
	"create", "new")
}

type create cmd

func (create) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Create what?")
		return
	}

	itemId, err := data.CreateItem(map[string]interface{}{
		"name":    strings.Join(s.input, " "),
		"creator": s.actor.Name,
		"type":    13,
	})

	if err {
		s.msg.Actor.SendBad ("Failed to create item.")
	}else{
		objects.Items[itemId], _ = objects.LoadItem(data.LoadItem(itemId))
		newItem := objects.Item{}
		copier.Copy(&newItem, objects.Items[itemId])
		s.actor.Inventory.Add(&newItem)
		s.msg.Actor.SendGood(newItem.Name + " added to your inventory.")
	}

	s.ok = true
	return
}