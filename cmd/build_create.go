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
		"Usage:  create (mob|item) name \n \n Create a brand new item with a name. \n Note:  Use the modify command to add modify traits of the object.",
		permissions.Builder,
		"create", "new")
}

type create cmd

func (create) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Create what?")
		return
	}

	switch strings.ToLower(s.input[0]) {
	case "item":
		itemId, err := data.CreateItem(map[string]interface{}{
			"name":    strings.Join(s.input[1:], " "),
			"creator": s.actor.Name,
			"type":    13,
		})

		if err {
			s.msg.Actor.SendBad("Failed to create item.")
		} else {
			objects.Items[itemId], _ = objects.LoadItem(data.LoadItem(itemId))
			newItem := objects.Item{}
			if err := copier.CopyWithOption(&newItem, objects.Items[itemId], copier.Option{DeepCopy: true}); err != nil {
				s.msg.Actor.SendBad("Failed to copy item.")
				return
			}
			s.actor.Inventory.Add(&newItem)
			s.msg.Actor.SendGood(newItem.Name + " added to your inventory.")
		}
	case "mob":
		mobId, err := data.CreateMob(strings.Join(s.input[1:], " "), s.actor.Name)

		if err {
			s.msg.Actor.SendBad("Failed to create mob.")
		} else {
			objects.Mobs[mobId], _ = objects.LoadMob(data.LoadMob(mobId))
			newMob := objects.Mob{}
			if err := copier.CopyWithOption(&newMob, objects.Mobs[mobId], copier.Option{DeepCopy: true}); err != nil {
				s.msg.Actor.SendBad("Failed to copy mob.")
				return
			}
			s.where.Mobs.Add(&newMob, false)
			newMob.StartTicking()
			s.msg.Actor.SendGood(newMob.Name + " added to the room.")
		}
	default:
		s.msg.Actor.SendBad("What do you want to create?")
	}

	s.ok = true
	return
}
