package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jinzhu/copier"
	"strconv"
	"strings"
)

func init() {
	addHandler(spawn{},
		"Usage:  spawn (mob|item) (name) \n \n Use this command to spawn a mob or item to be modified: \n"+
			"Items: Item will be added to your inventory\n"+
			"  -->  If you wish to save it as the template for that item, use the 'savetemplate item' command\n"+
			"Mob:  Mob will be spawned into your room. \n"+
			"  -->  If you wish to save it as the template for that mob, use the 'savetemplate mob' command\n\n",
		permissions.Builder,
		"spawn")
}

type spawn cmd

func (spawn) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Spawn what?")
		return
	}

	switch strings.ToLower(s.words[0]) {
	// Handle Rooms
	case "mob":
		//log.Println("Trying to do a spawn...")
		mobId, err := strconv.Atoi(s.words[1])
		if err != nil {
			s.msg.Actor.SendBad("What mob ID do you want to spawn?")
			return
		}
		//log.Println("Copying mob")
		newMob := objects.Mob{}
		copier.CopyWithOption(&newMob, objects.Mobs[mobId], copier.Option{DeepCopy: true})
		if newMob.Placement <= 0 {
			newMob.Placement = 5
		} else if newMob.Placement >= 6 {
			newMob.Placement = utils.Roll(5, 1, 0)
		}
		s.where.Mobs.Add(&newMob, false)
		newMob.StartTicking()
	case "item":
		itemId, err := strconv.Atoi(s.words[1])
		if err != nil {
			s.msg.Actor.SendBad("What item ID do you want to spawn?")
			return
		}
		newItem := objects.Item{}
		copier.CopyWithOption(&newItem, objects.Items[itemId], copier.Option{DeepCopy: true})
		s.actor.Inventory.Add(&newItem)
		s.msg.Actor.SendGood(newItem.Name + " added to your inventory.")
	default:
		s.msg.Actor.SendBad("Not an object that can be spawned")
	}

	s.ok = true
	return
}
