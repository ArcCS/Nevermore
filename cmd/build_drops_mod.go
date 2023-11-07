package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
)

func init() {
	addHandler(moddrop{},
		"Usage:  moddrop mob_id item_id drop_rate \n Modify an EXISTING itme drop from a mob.  In order: ## mob_id, ## item_id, ## new percent chance ",
		permissions.Builder,
		"moddrop")
}

type moddrop cmd

func (moddrop) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Mod which item and it's spawn rate?")
		return
	}

	var mobId, itemId, dropRate int
	mobId, err := strconv.Atoi(s.words[0])
	if err != nil {
		log.Println(err)
	}

	itemId, err2 := strconv.Atoi(s.words[1])
	if err2 != nil {
		log.Println(err2)
	}

	dropRate, err3 := strconv.Atoi(s.words[2])
	if err2 != nil {
		log.Println(err3)
	}

	if mob, ok := objects.Mobs[mobId]; ok {
		if _, ok := objects.Items[itemId]; ok {
			if dropRate > 100 {
				s.msg.Actor.SendBad("The sum of the drop rates is more than 100% with the new value")
			} else {
				mob.ItemList[itemId] = dropRate
				data.UpdateDrop(map[string]interface{}{
					"mob_id":  mobId,
					"item_id": itemId,
					"chance":  dropRate})
				s.msg.Actor.SendGood("Drop rate updated")
			}
		} else {
			s.msg.Actor.SendBad("Item not found.")
		}

	} else {
		s.msg.Actor.SendBad("That mob ID doesn't exist")
		return
	}

	s.ok = true
	return
}
