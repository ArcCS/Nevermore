package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/stats"
	"strconv"
	"time"
)

func init() {
	addHandler(shutdown{},
           "Usage:  shutdown \n \n Safely shutdown the game",
           permissions.Gamemaster,
           "shutdown")
}

type shutdown cmd

func (shutdown) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Shutting down in 5 minutes")
		tickerShutdown := time.NewTicker(60 * time.Second)
		stats.ActiveCharacters.MessageAll("The server will shut down in 5 minutes.  Please save your character and exit the game.")
		countDown := 5
		countCapture := 0
		go func() {
			for {
				select {
				case <-tickerShutdown.C:
					countCapture += 1
					if countCapture == countDown {
						tickerShutdown.Stop()
						config.Server.Running = false
						config.ServerShutdown<-true
					}else {
						stats.ActiveCharacters.MessageAll("The server will shut down in " + strconv.Itoa(countDown - countCapture) + " minutes.  Please save your character and exit the game.")
					}
				}
			}
		}()
	}else if s.words[0] == "NOW" {
		s.msg.Actor.SendInfo("Shutting down now!")
		stats.ActiveCharacters.MessageAll("GM initiated immediate server shut down.")
		config.Server.Running = false
		config.ServerShutdown<-true
	} else if minutes, err := strconv.Atoi(s.words[0]); err==nil {
		s.msg.Actor.SendInfo("Shutting down in " + strconv.Itoa(minutes) + " minutes")
		tickerShutdown := time.NewTicker(60 * time.Second)
		stats.ActiveCharacters.MessageAll("The server will shut down in " + strconv.Itoa(minutes) + " minutes minutes.  Please save your character and exit the game.")
		countCapture := 0
		go func() {
			for {
				select {
				case <-tickerShutdown.C:
					countCapture += 1
					if countCapture == minutes {
						tickerShutdown.Stop()
						config.Server.Running = false
						config.ServerShutdown<-true
					}else {
						stats.ActiveCharacters.MessageAll("The server will shut down in " + strconv.Itoa(minutes - countCapture) + " minutes.  Please save your character and exit the game.")
					}
				}
			}
		}()
	}else {
		s.msg.Actor.SendInfo("Unrecognized shutdown input")
	}
	s.ok = true
	return
}