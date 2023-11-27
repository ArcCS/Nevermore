package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
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
		objects.ActiveCharacters.MessageAll("The server will shut down in 5 minutes.  Please save your character and exit the game.", config.JarvoralChannel)
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
						config.ServerShutdown <- true
					} else {
						objects.ActiveCharacters.MessageAll("The server will shut down in "+strconv.Itoa(countDown-countCapture)+" minutes.  Please save your character and exit the game.", config.JarvoralChannel)
					}
				}
			}
		}()
	} else if s.words[0] == "NOW" {
		s.msg.Actor.SendInfo("Shutting down now!")
		objects.ActiveCharacters.MessageAll("GM initiated immediate server shut down.", config.JarvoralChannel)
		config.Server.Running = false
		config.ServerShutdown <- true
	} else if minutes, err := strconv.Atoi(s.words[0]); err == nil {
		s.msg.Actor.SendInfo("Shutting down in " + strconv.Itoa(minutes) + " minutes")
		tickerShutdown := time.NewTicker(60 * time.Second)
		objects.ActiveCharacters.MessageAll("The server will shut down in "+strconv.Itoa(minutes)+" minutes.  Please save your character and exit the game.", config.JarvoralChannel)
		countCapture := 0
		go func() {
			for {
				select {
				case <-tickerShutdown.C:
					countCapture += 1
					if countCapture == minutes {
						tickerShutdown.Stop()
						config.Server.Running = false
						config.ServerShutdown <- true
					} else {
						objects.ActiveCharacters.MessageAll("The server will shut down in "+strconv.Itoa(minutes-countCapture)+" minutes.  Please save your character and exit the game.", config.JarvoralChannel)
					}
				}
			}
		}()
	} else {
		s.msg.Actor.SendInfo("Unrecognized shutdown input")
	}
	s.ok = true
	return
}
