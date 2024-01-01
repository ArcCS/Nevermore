// Stats and global listing of characters.

package objects

import (
	"fmt"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/message"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Currently active characters
type characterStats struct {
	sync.Mutex
	list []*Character
}

var LastActivity = map[string]time.Time{}
var ActiveCharacters = &characterStats{}
var IpMap = make(map[string]string)

func GetLastActivity(name string) time.Time {
	if val, ok := LastActivity[name]; ok {
		return val
	}
	// If there's an error return over an hour to ensure that everyone hits an AFK
	return time.Now()
}

// Add adds the specified character to the list of characters.
func (c *characterStats) Add(character *Character, address string) {
	if character.Flags["invisible"] || character.Permission.HasAnyFlags(permissions.Builder, permissions.Gamemaster, permissions.Dungeonmaster, permissions.God) {
		c.MessageGM("###: " + character.Name + "[" + address + "] joins the realm.")
	} else {
		c.MessageAll("###: "+character.Name+" joins the realm.", config.JarvoralChannel)
	}
	c.Lock()
	c.list = append(c.list, character)
	IpMap[character.Name] = address
	c.Unlock()
}

// Remove Pass character as a pointer, compare and remove
func (c *characterStats) Remove(character *Character) {
	c.Lock()
	if character.Flags["invisible"] || character.Permission.HasAnyFlags(permissions.God, permissions.Builder, permissions.Gamemaster, permissions.Dungeonmaster) {
		c.MessageGMExcept("###:"+character.Name+" departs the realm.", character)
	} else {
		c.MessageExcept("###: "+character.Name+" departs the realm.", character, config.JarvoralChannel)
	}

	for i, p := range c.list {
		if p == character {
			copy(c.list[i:], c.list[i+1:])
			c.list[len(c.list)-1] = nil
			c.list = c.list[:len(c.list)-1]
			delete(IpMap, character.Name)
			break
		}
	}

	delete(LastActivity, character.Name)

	if len(c.list) == 0 {
		c.list = make([]*Character, 0, 10)
	}

	c.Unlock()
	log.Println("Completed Character removal from stats containers")
}

func (c *characterStats) Find(name string) *Character {
	c.Lock()
	for _, p := range c.list {
		if strings.ToLower(p.Name) == strings.ToLower(name) {
			c.Unlock()
			return p
		}
	}
	c.Unlock()
	return nil
}

// List returns the names of all characters in the character list. The omit parameter
// may be used to specify a character that should be omitted from the list.
func (c *characterStats) List() []string {
	c.Lock()

	list := make([]string, 0, len(c.list))

	for _, character := range c.list {
		if character.Flags["invisible"] == true {
			continue
		}

		var calc time.Duration
		if !LastActivity[character.Name].IsZero() {
			calc = time.Now().Sub(LastActivity[character.Name])
		}
		charState := ""
		if calc.Minutes() > 2 {
			charState = fmt.Sprintf("[idle: %s]", strconv.Itoa(int(calc.Minutes())))
		}
		if character.Flags["ooc"] {
			charState += " [OOC] "
		}
		if character.Flags["afk"] {
			charState += " [AFK]"
		}
		if charState != "" {
			charState = "/" + charState
		}
		if character.Title != "" {
			list = append(list, fmt.Sprintf("%s(%s), the %s, %s, %s", character.Name, strconv.Itoa(character.Tier), character.ClassTitle, character.Title, charState))
		} else {
			list = append(list, fmt.Sprintf("%s(%s), the %s, %s", character.Name, strconv.Itoa(character.Tier), character.ClassTitle, charState))
		}
	}

	c.Unlock()
	return list
}

// GMList returns the names of all characters in the character list. The omit parameter
// may be used to specify a character that should be omitted from the list.
func (c *characterStats) GMList() []string {
	c.Lock()

	list := make([]string, 0, len(c.list))

	for _, character := range c.list {
		calc := time.Now().Sub(LastActivity[character.Name])
		charState := ""
		if character.Flags["ooc"] {
			charState += " [OOC] "
		}
		if character.Flags["afk"] {
			charState += " [AFK]"
		}
		if charState != "" {
			charState = "/" + charState + fmt.Sprintf("[Activity: %ss]", strconv.Itoa(int(calc.Seconds())))
		} else {
			charState = fmt.Sprintf("[Activity: %ss]", strconv.Itoa(int(calc.Seconds())))
		}

		if character.Title != "" {
			list = append(list, fmt.Sprintf("(Room: %s) (%s) %s(%s), %s, %s %s", strconv.Itoa(character.ParentId), IpMap[character.Name], character.Name, strconv.Itoa(character.Tier), character.ClassTitle, character.Title, charState))
		} else {
			list = append(list, fmt.Sprintf("(Room: %s) (%s) %s(%s), %s %s", strconv.Itoa(character.ParentId), IpMap[character.Name], character.Name, strconv.Itoa(character.Tier), character.ClassTitle, charState))
		}
	}

	c.Unlock()
	return list
}

func (c *characterStats) MessageExcept(msg string, except *Character, channel string) {
	if DiscordSession != nil {
		if _, err := DiscordSession.ChannelMessageSend(channel, msg); err != nil {
			log.Println("Error sending message to discord:", err)
		}
	}
	// Setup buffer
	msgbuf := message.AcquireBuffer()
	msgbuf.Send(text.White, msg)
	var players []io.Writer
	for _, p := range c.list {
		if p != except {
			players = append(players, p)
		}
	}
	msgbuf.Deliver(players...)

	return
}

func (c *characterStats) MessageAll(msg string, channel string) {
	if DiscordSession != nil {
		if _, err := DiscordSession.ChannelMessageSend(channel, msg); err != nil {
			log.Println("Error sending message to discord:", err)
		}
	}
	c.Lock()

	// Setup buffer
	msgbuf := message.AcquireBuffer()
	msgbuf.Send(text.White, msg)
	var players []io.Writer
	for _, p := range c.list {
		players = append(players, p)
	}
	msgbuf.Deliver(players...)

	c.Unlock()
	return
}

func (c *characterStats) MessageGMExcept(msg string, except *Character) {
	// Setup buffer
	msgbuf := message.AcquireBuffer()
	msgbuf.Send(text.White, "[GM] ", msg)
	var players []io.Writer
	for _, p := range c.list {
		if p != except && p.Permission.HasAnyFlags(permissions.God, permissions.NPC, permissions.Dungeonmaster, permissions.Gamemaster, permissions.Builder) {
			players = append(players, p)
		}
	}
	msgbuf.Deliver(players...)

	return
}

func (c *characterStats) MessageGM(msg string) {
	c.Lock()

	// Setup buffer
	msgbuf := message.AcquireBuffer()
	msgbuf.Send(text.White, "[GM] ", msg)
	var players []io.Writer
	for _, p := range c.list {
		if p.Permission.HasAnyFlags(permissions.God, permissions.NPC, permissions.Dungeonmaster, permissions.Gamemaster, permissions.Builder) {
			players = append(players, p)
		}
	}
	msgbuf.Deliver(players...)

	c.Unlock()
	return
}

// Len returns the length of the character list.
func (c *characterStats) Len() (l int) {
	c.Lock()
	l = len(c.list)
	c.Unlock()
	return
}
