package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
)

func init() {
	addHandler(remdrop{},
		"Usage:  remdrop mob_id item_id \n \n Remove a spawn from the encounter table \n" ,
		permissions.Builder,
		"remdrop")
}

type remdrop cmd

func (remdrop) process(s *state) {
	if len(s.words) < 1{
		s.msg.Actor.SendInfo("Remove what?")
		return
	}

	var mob_id, item_id int
	mob_id, err := strconv.Atoi(s.words[0])
	if err != nil {
		log.Println(err)
	}
	item_id, err2 := strconv.Atoi(s.words[1])
	if err2 != nil {
		log.Println(err2)
	}

	if mob, ok := objects.Mobs[mob_id]; ok {
		if _, ok := mob.ItemList[item_id]; ok {
			delete(mob.ItemList, item_id)
			data.DeleteDrop(mob_id, item_id)
			s.msg.Actor.SendGood("Item removed from this mobs encounter table.")
		}
	}else{
		s.msg.Actor.SendBad("That mob ID doesn't exist.")
		return
	}

	s.ok = true
	return
}
