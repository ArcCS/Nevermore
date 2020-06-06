package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(dig{},
           "Usage:  dig exit_name exit_back (room_name)  \n Required: exit_name room_name \n \n Use dig with single world exit names (can be changed later) and a room name to create the new room, the exit to, and exit back.   Temporary exit can be deleted if you don't want an exit back to the previous room. \n",
           permissions.Builder,
           "dig")
}

type dig cmd

func (dig) process(s *state) {
	if len(s.words) < 3 {
		s.msg.Actor.SendInfo("Not enough parameters to dig")
		return
	}

	// Check that we aren't duping an exit
	if !data.ExitExists(strings.ToLower(s.words[0]), s.where.RoomId){
	// First create the room
	roomId, rErr := data.CreateRoom(strings.Join(s.input[2:], " "), s.actor.Name)
	if !rErr {
		createTo := data.CreateExit(map[string]interface{}{
			"name": strings.ToLower(s.words[0]),
			"fromId": s.where.RoomId,
			"toId": roomId,
		})
		if createTo {
			s.msg.Actor.SendBad("To exit creation failed.")
			return
		}
		createBack := data.CreateExit(map[string]interface{}{
			"name": strings.ToLower(s.words[1]),
			"fromId": roomId,
			"toId": s.where.RoomId,
		})
		if createBack {
			s.msg.Actor.SendBad("To exit creation failed.")
			return
		}
		objects.Rooms[roomId], _ = objects.LoadRoom(data.LoadRoom(roomId))
		exitData := data.LoadExit(strings.ToLower(s.words[0]), s.where.RoomId, roomId)
		s.where.Exits[strings.ToLower(s.words[0])] = objects.NewExit(s.where.RoomId, exitData)
		s.msg.Actor.SendGood("Exits and room created and loaded into game.")
	}
	}else{
		s.msg.Actor.SendBad("An exit with that name  exists in your current room.")
	}

	s.ok = true
	return
}