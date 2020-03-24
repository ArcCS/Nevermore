package config

import (
	"log"
	"math/rand"
	"time"
)

// Server default configuration
var Server = struct {
	Host           string        // Host for server to listen on
	DBUname		   string		// The Username for the neo4j instance
	DBPword		   string		// The Password for the neo4j instance
	Port           string        // Port for server to listen on
	Greeting       []byte        // Connection greeting
	Motd		   string		// MOTD when logging in
	IdleTimeout    time.Duration // Idle connection disconnect time
	MaxPlayers     int           // Max number of players allowed to login at once
	DataDir        string        // Main data directory
	MaxCharacters  int			// Maximum number of characters
	SearchResults int		// Max search results
	Running	bool
}{
	Host:           "127.0.0.1",
	DBUname:		"USERNAME",
	DBPword:		"PASSWORD",
	Port:           "4001",
	Greeting:       []byte("Welcome to Aalynor's Nexus."),
	Motd:			"",
	IdleTimeout:    10 * time.Minute,
	MaxPlayers:     1024,
	DataDir:        ".",
	MaxCharacters:  20,
	SearchResults: 15,
	Running: true,
}

// Stats default configuration
var Stats = struct {
	Rate time.Duration // Stats collection and display rate
	GC   bool          // Run garbage collection before stat collection
}{
	Rate: 10 * time.Second,
	GC:   false,
}

// Inventory default configuration
var Inventory = struct {
	Compact   int // only compact if cap - len*2 > compact
	CrowdSize int // If inventory has more player than this it's a crowd
}{
	Compact:   4,
	CrowdSize: 10,
}

// Login default configuration
var Login = struct {
	AccountLength  int
	PasswordLength int
	SaltLength     int
}{
	AccountLength:  3,
	PasswordLength: 8,
	SaltLength:     32,
}

// Debugging configuration
var Debug = struct {
	LongLog    bool // Long log with microseconds & filename?
	Panic      bool // Let goroutines panic and stop server?
	AllowDump  bool // Allow use of #DUMP command?
	AllowDebug bool // Allow use of #DEBUG command?
	Events     bool // Log events? - this can make the log quite noisy
	Things     bool // Log additional information for Thing?
}{
	LongLog:    false,
	Panic:      false,
	AllowDump:  false,
	AllowDebug: false,
	Events:     false,
	Things:     false,
}

var StartingRoom = 3
var OocRoom = 2

// Load loads the configuration file and overrides the default configuration
// values with any values found.
func init() {

	// Setup global logging format
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)

	// Seed default random source
	rand.Seed(time.Now().UnixNano())

	if !Debug.LongLog {
		log.SetFlags(log.Ldate | log.Ltime)
		log.Printf("Switching to short log format.")
	}
}
