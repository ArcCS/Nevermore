package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
	"strings"
)

func init() {
	addHandler(destroy{}, "destroy", "delete", "del")
	addHelp("Usage:  destroy (room|mob|item|exit) ID/name) \n \n Delete the item entirely from the database.  If this is a builder account, you must be the creator of the item to delete.  If you delete a room, it will delete the exits to and from it; be mindful and grab id's or tunnel around prior to deletion.", 50, "destroy", "delete", "del")
}

type destroy cmd

func (destroy) process(s *state) {
	// Handle Permissions
	if s.actor.Class < 50 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
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
			if s.actor.Class > 50 || room.Creator == s.actor.Name {
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
			if s.actor.Class > 50 || s.where.Creator == s.actor.Name {
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