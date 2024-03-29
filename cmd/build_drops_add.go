package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
)

func init() {
	addHandler(adddrop{},
		"Usage:  adddrop mob_id item_id drop_rate \n In order: add a drop to MOB ##, ITEM ##, WHOLE NUMBER PERCENT CHANCE ##  \n",
		permissions.Builder,
		"adddrop")
}

type adddrop cmd

func (adddrop) process(s *state) {
	if len(s.words) < 3 {
		s.msg.Actor.SendInfo("Missing some arguments?")
		return
	}

	var mobId, itemId, dropRate int
	mobId, err := strconv.Atoi(s.words[0])
	if err != nil {
		log.Println(err)
	}

	itemId, err2 := strconv.Atoi(s.words[1])
	if err != nil {
		log.Println(err2)
	}

	dropRate, err3 := strconv.Atoi(s.words[2])
	if err != nil {
		log.Println(err3)
	}

	if _, ok := objects.Mobs[mobId]; ok {
		if _, ok := objects.Items[itemId]; ok {
			if len(objects.Mobs[mobId].ItemList) < 10 {
				if dropRate <= 100 {
					data.CreateDrop(map[string]interface{}{
						"mobId":  mobId,
						"itemId": itemId,
						"chance": dropRate})
					objects.Mobs[mobId].ItemList[itemId] = dropRate
					s.msg.Actor.SendGood("Drop added to mob drops")
				} else {
					s.msg.Actor.SendBad("You can't set a drop to more than 100%")
				}
			} else {
				s.msg.Actor.SendBad("There are already 10 items in this mobs drop list.")
			}

		} else {
			s.msg.Actor.SendBad("That item ID doesn't exist.")
		}
	} else {
		s.msg.Actor.SendBad("That mob ID doesn't exist")
		return
	}

	s.ok = true
	return
}
