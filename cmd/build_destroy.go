package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
	"strings"
)

func init() {
	addHandler(destroy{},
	"Usage:  destroy (room|mob|item|exit) ID/name) \n \n Delete the item entirely from the database.  If this is a builder account, you must be the creator of the item to delete.  If you delete a room, it will delete the exits to and from it; be mindful and grab id's or tunnel around prior to deletion.",
	permissions.Builder,
	"destroy", "delete", "del")
}

type destroy cmd

func (destroy) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Delete what?")
		return
	}

	switch strings.ToLower(s.input[0]){
	case "room":
		objectRef, _ := strconv.Atoi(s.input[1])
		if int64(objectRef) == s.where.RoomId {
			s.msg.Actor.SendBad("Don't delete a room while you're standing in it...")
			return
		}
		room, rErr := objects.Rooms[int64(objectRef)]
		if rErr {
			if s.actor.Permission.HasFlag(permissions.Builder) || room.Creator == s.actor.Name {
				data.DeleteRoom(int64(objectRef))
				delete(objects.Rooms, int64(objectRef))
				s.where.CleanExits()
				s.msg.Actor.SendGood("Deleted room successfully.")
			}else{
				s.msg.Actor.SendBad("No permissions to modify your current location. ")
			}
		}else{
			s.msg.Actor.SendBad("Couldn't find room.")
		}
	case "exit":
		exitName := s.input[1]
		if len(s.words) > 0 {
			exitName = strings.Join(s.input[1:], " ")
		}
		objectRef := strings.ToLower(exitName)
		if !utils.StringIn(strings.ToUpper(objectRef), directionals) {
			for txtE, _ := range s.where.Exits {
				if strings.Contains(txtE, objectRef) {
					objectRef = txtE
				}
			}
		}
		exit, rErr := s.where.Exits[objectRef]
		if rErr {
			if (s.actor.Permission.HasFlags(permissions.Builder, permissions.Dungeonmaster)) || objects.Rooms[exit.ToId].Creator == s.actor.Name {
				data.DeleteExit(exit.Name, s.where.RoomId )
				delete(s.where.Exits, objectRef)
				s.where.CleanExits()
				s.msg.Actor.SendGood("Deleted exit successfully.")
			}else{
				s.msg.Actor.SendBad("No permissions to modify your current location. ")
			}

		}else{
			s.msg.Actor.SendBad("Couldn't find exit.")
		}
	default:
		s.msg.Actor.SendBad("Unknown world object")
	}

	s.ok = true
	return
}