// Copyright 2021 Nevermore.

package main

import (
	"github.com/ArcCS/Nevermore/comms"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/intelligence"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/stats"
	"io"
	"log"
	"os"
	"time"
)

var RoomSyncTicker *time.Ticker

func main() {
	logFile, err := os.OpenFile("log_"+time.Now().Format("01-02-2006")+".txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Println("Error closing log file:", err)
		}
	}()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	stats.Start()
	go objects.StartJarvoral()
	objects.Load()
	log.Println("Starting time...")
	StartTime()
	intelligence.StartRoomAI()
	StartSync()
	comms.Listen(config.Server.Host, config.Server.Port)
}

func StartSync() {
	RoomSyncTicker = time.NewTicker(3 * time.Minute)
	go func() {
		for {
			select {
			case <-RoomSyncTicker.C:
				objects.FlushRoomUpdates()
				data.FlushChatLogs()
				data.FlushItemSales()
				data.FlushCombatMetrics()
			}
		}
	}()
}

func StartTime() {
	SyncTime()
	objects.WorldTicker = time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-objects.WorldTickerUnload:
				return
			case <-objects.WorldTicker.C:
				SyncTime()
			}
		}
	}()
}

func SyncTime() {
	// Sync Time
	diffHours, diffDays, diffMonths, diffYears := config.SyncCurrentTime()
	objects.YearPlus = diffYears
	objects.CurrentDay = diffDays % 9
	if diffDays%30 == 0 {
		objects.DayOfMonth = 30
	} else {
		objects.DayOfMonth = diffDays % 30
	}
	objects.CurrentMonth = diffMonths % 12
	if objects.CurrentHour != diffHours%24 && diffHours%24 == config.Months[objects.CurrentMonth]["sunrise"] {
		objects.ActiveCharacters.MessageAll("### The suns rise over the mountains to the east.")
	} else if objects.CurrentHour != diffHours%24 && diffHours%24 == config.Months[objects.CurrentMonth]["sunset"] {
		objects.ActiveCharacters.MessageAll("### The suns dip below the horizon to the west.")
	}
	objects.CurrentHour = diffHours % 24
	if objects.CurrentHour >= config.Months[objects.CurrentMonth]["sunrise"].(int) && objects.CurrentHour < config.Months[objects.CurrentMonth]["sunset"].(int) {
		objects.DayTime = true
	} else {
		objects.DayTime = false
	}
}
