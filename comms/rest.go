package comms

import (
	"encoding/json"
	"fmt"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Token struct {
	Token string `json:"token"`
}

func ListenRest() {
	log.Println("Starting REST server on port 1234")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/cleanshutdown", CleanShutDown).Methods("POST")

	log.Fatal(http.ListenAndServe(config.Server.Host+":1234", router))
}

func CleanShutDown(w http.ResponseWriter, r *http.Request) {
	var t Token
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Token != config.Server.RestToken {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

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

	log.Println("Received Rest Call for a Clean Shutdown")
	fmt.Fprint(w, "Valid token received")
}
