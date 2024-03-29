// Copyright 2023 Nevermore.

package intelligence

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
	"time"
)

var ActiveRooms []int

func init() {
	objects.ActivateRoom = ActivateRoom
}

func ActivateRoom(roomId int) {
	if !utils.IntIn(roomId, ActiveRooms) {
		ActiveRooms = append(ActiveRooms, roomId)
		go InitialRoom(roomId)
	}
}

func DeactivateRoom(roomId int) {
	for c, p := range ActiveRooms {
		if p == roomId {
			copy(ActiveRooms[c:], ActiveRooms[c+1:])
			ActiveRooms = ActiveRooms[:len(ActiveRooms)-1]
			break
		}
	}
}

func StartRoomAI() {
	log.Println("Starting room AI")
	go func() {
		for {
			LoopRooms()
			time.Sleep(1 * time.Second)
		}
	}()
}

func LoopRooms() {
	for _, r := range ActiveRooms {
		objects.Rooms[r].LockRoom(strconv.Itoa(r)+":RoomTick", true)
		if len(objects.Rooms[r].Chars.Contents) <= 0 {
			if !time.Time.IsZero(objects.Rooms[r].EvacuateTime) &&
				time.Now().Sub(objects.Rooms[r].EvacuateTime).Seconds() > float64(config.RoomClearTimer) {
				DeactivateRoom(r)
				objects.Rooms[r].CleanRoom()
			}
		} else {
			if objects.Rooms[r].Flags["encounters_on"] {
				if time.Now().Sub(objects.Rooms[r].LastEncounterTime).Seconds() > float64(objects.Rooms[r].EncounterSpeed) {
					// Apply jigger
					if time.Now().Sub(objects.Rooms[r].LastEncounterTime).Seconds() > float64(objects.Rooms[r].EncounterSpeed+objects.Rooms[r].EncounterJigger) {
						objects.Rooms[r].Encounter()
					}
				}
			}
			if (objects.Rooms[r].Flags["fire"] ||
				objects.Rooms[r].Flags["earth"] ||
				objects.Rooms[r].Flags["wind"] ||
				objects.Rooms[r].Flags["water"]) &&
				time.Now().Sub(objects.Rooms[r].LastEffectTime).Seconds() > float64(config.RoomEffectInvocation) {
				objects.Rooms[r].ElementalDamage()
			}
		}
		objects.Rooms[r].UnlockRoom(strconv.Itoa(r)+":RoomTick", true)
	}
}

func InitialRoom(r int) {
	objects.Rooms[r].LockRoom(strconv.Itoa(r)+":RoomInit", false)
	if len(objects.Rooms[r].Chars.Contents) >= 0 {
		if objects.Rooms[r].Flags["encounters_on"] {
			objects.Rooms[r].Encounter()
		}
	}
	if objects.Rooms[r].Flags["fire"] ||
		objects.Rooms[r].Flags["earth"] ||
		objects.Rooms[r].Flags["wind"] ||
		objects.Rooms[r].Flags["water"] {
		objects.Rooms[r].ElementalDamage()
	}
	objects.Rooms[r].UnlockRoom(strconv.Itoa(r)+":RoomInit", false)
}
