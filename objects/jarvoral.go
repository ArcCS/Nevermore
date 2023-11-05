package objects

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strconv"
	"strings"
)

// Variables used for command line parameters
var (
	DiscordSession *discordgo.Session
	err            error
)

func StartJarvoral() {
	// Create a new Discord session using the provided bot token.
	if os.Getenv("DISCORDTOKEN") == "" {
		fmt.Println("No Discord Token found")
		return
	}

	DiscordSession, err = discordgo.New("Bot " + os.Getenv("DISCORDTOKEN"))
	if err != nil {
		fmt.Println("Discord session was not created: \n,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	DiscordSession.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	DiscordSession.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = DiscordSession.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Discord Session Initiated")

}

func StopJarvoral() {
	if DiscordSession != nil {
		DiscordSession.Close()
	}
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
		players := ActiveCharacters.List()

		if len(players) == 0 {
			//log.Println(m.ChannelID)
			s.ChannelMessageSend(m.ChannelID, "There is currently no one visibly playing.")
			return
		}

		if len(players) > 1 {
			message = "There are currently " + strconv.Itoa(len(players)) + " players in the realms. \n"
		} else {
			message = "There is currently " + strconv.Itoa(len(players)) + " player in the realms. \n"
		}

		for _, player := range players {
			message += "	" + player + " \n"
		}

		s.ChannelMessageSend(m.ChannelID, message)
	}
}
