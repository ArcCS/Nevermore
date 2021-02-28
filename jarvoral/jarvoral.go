package jarvoral

import (
	"fmt"
	"github.com/ArcCS/Nevermore/stats"
	"github.com/bwmarrin/discordgo"
	"os"
	"strconv"
	"strings"
)

// Variables used for command line parameters
var (
	dg *discordgo.Session
	err error
)


func StartJarvoral() {
	// Create a new Discord session using the provided bot token.
	dg, err = discordgo.New("Bot " + os.Getenv("DISCORDTOKEN"))
	if err != nil {
		fmt.Println("Discord session was not created: \n,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Discord Session Initiated")

}

func StopJarvoral(){
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if strings.ToLower(m.Content) == "who" {
		var message string
		players := stats.ActiveCharacters.List()

		if len(players) == 0 {
			s.ChannelMessageSend(m.ChannelID, "There is currently no one visibly playing.")
			return
		}

		if len(players) > 1 {
			message  = "There are currently " + strconv.Itoa(len(players)) + " players in the realms. \n"
		}else {
			message = "There is currently " + strconv.Itoa(len(players)) + "player in the realms. \n"
		}

		//TODO: Add some more information to the output list AFK? Currently hunting?  RP only?
		for _, player := range players {
			message += "	" + player + " \n"
		}

		s.ChannelMessageSend(m.ChannelID, message)
	}
}
