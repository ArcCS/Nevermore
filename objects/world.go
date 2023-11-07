package objects

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"runtime"
	"time"
)

// Rooms contains all the world rooms tagged with their room_id
// This makes it very simple to move people by room_id and retain their connecting exit
var (
	Rooms = map[int]*Room{}
	Mobs  = map[int]*Mob{}
	Items = map[int]*Item{}

	CurrentDay   = 0
	YearPlus     = 0
	CurrentMonth = 0
	DayOfMonth   = 0
	CurrentHour  = 0

	WorldTicker       *time.Ticker
	WorldTickerUnload = make(chan bool)

	DayTime = false
)

// Create a bind variable
var (
	Script             func(o *Character, input string) string
	RoomsPendingUpdate []int
)

func AddRoomUpdate(roomId int) {
	if !utils.IntIn(roomId, RoomsPendingUpdate) {
		RoomsPendingUpdate = append(RoomsPendingUpdate, roomId)
	}
}

func FlushRoomUpdates() {
	for _, roomId := range RoomsPendingUpdate {
		if len(Rooms[roomId].Chars.Contents) <= 0 {
			Rooms[roomId].Save()
		}
	}
	RoomsPendingUpdate = nil
}

// Load fills the world with love.
func Load() {
	log.Printf("Loading mobs")
	preparse := data.LoadMobs()
	for _, mob := range preparse {
		if mob != nil {
			mobData := mob.(map[string]interface{})
			Mobs[int(mobData["mob_id"].(int64))], _ = LoadMob(mobData)
		}

	}
	log.Printf("Finished loading %d mobs.", len(Mobs))
	preparse = nil

	log.Printf("Loading items")
	preparse = data.LoadItems()
	for _, item := range preparse {
		if item != nil {
			itemData := item.(map[string]interface{})
			Items[int(itemData["item_id"].(int64))], _ = LoadItem(itemData)
		}

	}
	log.Printf("Finished loading %d items.", len(Items))
	preparse = nil

	log.Printf("Loading rooms")
	preparse = data.LoadRooms()
	for _, room := range preparse {
		if room != nil {
			roomData := room.(map[string]interface{})
			Rooms[int(roomData["room_id"].(int64))], _ = LoadRoom(roomData)
		}

	}

	log.Printf("Finished loading %d rooms.", len(Rooms))

	preparse = nil

	runtime.GC()

}
