package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
	"strings"
)

func init() {
	addHandler(link{},
		"Usage:  link exit_name (room_id) [exit_back] \n Required: exit_name room_name \n \n Dig creates a new exit with the exit name, and links it to the room ID specified. If you specify a name back, the exit back will be automatically generated.  \n Optionals: exit_back will create the exit name back to current room \n",
		permissions.Builder,
		"link", "tunnel")
}

type link cmd

func (link) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Link where?")
		return
	}

	roomRef, _ := strconv.Atoi(s.input[1])
	// Check that the room exists.
	if room, ok := objects.Rooms[roomRef]; !ok {
		s.msg.Actor.SendBad("That room ID doesn't appear to exist.")
	} else {
		// Check that we aren't duping an exit
		if !data.ExitExists(strings.ToLower(s.words[0]), s.where.RoomId) {
			createTo := data.CreateExit(map[string]interface{}{
				"name":   strings.ToLower(s.words[0]),
				"fromId": s.where.RoomId,
				"toId":   roomRef,
			})
			if createTo {
				s.msg.Actor.SendBad("To exit creation failed.")
				return
			}
			exitData := data.LoadExit(strings.ToLower(s.words[0]), s.where.RoomId, roomRef)
			s.where.Exits[strings.ToLower(s.words[0])] = objects.NewExit(s.where.RoomId, exitData)
			if len(s.input) == 3 {
				createBack := data.CreateExit(map[string]interface{}{
					"name":   strings.ToLower(s.words[2]),
					"fromId": roomRef,
					"toId":   s.where.RoomId,
				})
				if createBack {
					s.msg.Actor.SendBad("Exit back creation failed.")
					return
				}
				exitData := data.LoadExit(strings.ToLower(s.words[0]), roomRef, s.where.RoomId)
				if exitData == nil {
					s.msg.Actor.SendBad("Exit back creation failed.")
					return
				}
				room.Exits[strings.ToLower(s.words[0])] = objects.NewExit(room.RoomId, exitData)
			}
			s.msg.Actor.SendGood("Exits created and loaded into game.")
		} else {
			s.msg.Actor.SendBad("An exit with that name  exists in your current room.")
		}
	}

	s.ok = true
	return
}
