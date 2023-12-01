package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/jinzhu/copier"
	"strings"
)

func init() {
	addHandler(buybag{},
		"",
		permissions.Player,
		"$BUYBAG")
}

type buybag cmd

var bagItem = 3978

func (buybag) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendBad("Your bag needs a name and a description.")
		return
	}

	if s.actor.Gold.Value < baseBag {
		s.msg.Actor.SendBad("You don't have enough gold to buy a bag.")
		return
	}

	name := s.input[0]
	desc := strings.Join(s.input[1:], " ")

	newItemId, succ := data.CopyItem(bagItem)
	if succ {
		objects.Items[newItemId], _ = objects.LoadItem(data.LoadItem(newItemId))
		objects.Items[newItemId].Name = name
		objects.Items[newItemId].Description = desc

		newItem := objects.Item{}
		if err := copier.CopyWithOption(&newItem, objects.Items[newItemId], copier.Option{DeepCopy: true}); err != nil {
			s.msg.Actor.SendBad("Bag could not be created.")
			return
		}
		newItem.Save()
		s.actor.Inventory.Add(&newItem)
		s.actor.Gold.Value -= baseBag
		s.msg.Actor.SendGood(newItem.Name + " added to your inventory.")
		return
	} else {
		s.msg.Actor.SendBad("Failed to copy primary bag item.")
		return
	}
}
