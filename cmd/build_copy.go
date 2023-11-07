package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
	"strings"
)

func init() {
	addHandler(buildCopy{}, "Usage: copy (mob|item) (SubjectID) \n \n Use this to copy an existing item in the database \n",
		permissions.Builder,
		"copy", "duplicate")
}

type buildCopy cmd

func (buildCopy) process(s *state) {

	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Spawn what?")
		return
	}

	switch strings.ToLower(s.words[0]) {
	// Handle Rooms
	case "mob":
		mobId, err := strconv.Atoi(s.words[1])
		if err != nil {
			s.msg.Actor.SendBad("What mob ID do you want to spawn?")
			return
		}
		newMobId, succ := data.CopyMob(mobId)
		if succ {
			objects.Mobs[newMobId], _ = objects.LoadMob(data.LoadMob(newMobId))
			s.msg.Actor.SendGood("Created New mob with id: " + strconv.Itoa(newMobId))
			return
		} else {
			s.msg.Actor.SendBad("Failed to copy the mob.")
			return
		}

	case "item":
		itemId, err := strconv.Atoi(s.words[1])
		if err != nil {
			s.msg.Actor.SendBad("What item ID do you want to spawn?")
			return
		}
		newItemId, succ := data.CopyItem(itemId)
		if succ {
			objects.Items[newItemId], _ = objects.LoadItem(data.LoadItem(newItemId))
			s.msg.Actor.SendGood("Created New mob with id: " + strconv.Itoa(newItemId))
			return
		} else {
			s.msg.Actor.SendBad("Failed to copy the item.")
			return
		}
	default:
		s.msg.Actor.SendBad("Not an object that can be spawned")
	}

	s.ok = true
	return
}
