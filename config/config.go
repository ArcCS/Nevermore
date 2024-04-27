package config

import (
	"fmt"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/spf13/viper"
	"log"
	"time"
)

// Server default configuration
var Server = struct {
	Host              string  // Host for server to listen on
	NEOUname          string  // The Username for the neo4j instance
	NEOPword          string  // The Password for the neo4j instance
	NEOAddress        string  // The address of the neo4j server
	PGUname           string  // The Username for the neo4j instance
	PGPword           string  // The Password for the neo4j instance
	PGAddress         string  // The address of the neo4j server
	PGPort            int     // The port of the postgres server
	Port              string  // Port for server to listen on
	IdleTimeout       float64 // Idle connection disconnect time in seconds
	AFKTimeout        float64
	OOCTimeout        float64
	MaxPlayers        int    // Max number of players at once
	DataDir           string // Main data directory
	MaxCharacters     int    // Maximum number of characters
	SearchResults     int    // Max search results
	Running           bool
	CreateChars       bool
	PermissionDefault int
	RestToken         string
}{
	Host:              "127.0.0.1",
	NEOUname:          "USERNAME",
	NEOPword:          "PASSWORD",
	NEOAddress:        "127.0.0.1",
	Port:              "4001",
	PGUname:           "USERNAME",
	PGPword:           "PASSWORD",
	PGAddress:         "127.0.0.1",
	PGPort:            5432,
	IdleTimeout:       15, // Minutes
	AFKTimeout:        30,
	OOCTimeout:        20,
	MaxPlayers:        1024,
	DataDir:           ".",
	MaxCharacters:     100,
	SearchResults:     15,
	Running:           true,
	CreateChars:       true,
	PermissionDefault: 2,
	RestToken:         "",
}

var DragonAscii = text.Red + `
                                      #**###                                                        
                                     ####**##                                                       
                         ###         #####***##                             #                       
                       ################**###*####                     #######                       
                     ###***##***#****##**##*****##              #########                           
                   ###*****##*****#******##********#      #####**###*####                           
                 ####******##******######***##******#     #####*##**#####                           
                ###*********#*****#########**********#      ##      ##*##                           
               ####*********#**################***#***#             ##**#                           
              #####*********###############****###**##*#          ##****                            
             ##**##******###########**#*#### ##*##***#############*****                             
            ###***#****# #########******###   *##**#*#####**##*******                               
           ###****#**** ########*****    # ####**************#**#*####                              
          ###*****###  ###*##**#*##    ######***#***********##############                          
          ##******##   ##******#*#  #######****************########  #####                          
         ##******###  *#****##  ## #####******************##                                        
         ##******  # ##*****      #####*******########*****                                         
         ##*****     ##***#        ####***** #########**##  *                                       
         #****#       ##****#      #####******#  #####*#######**********                            
         ##**#          ###*####*#   ###*****************************#*******#*                     
         ##*                            #########*########***####         ##****##                  
         ###                                                                   #**#                 
           ###          ` + text.White + `Welcome to Aalynor's Nexus!` + text.Red + `                                                 
             ###                                                                                    
                                                                                                    
                                                                                                    
` + text.White

// Flip to false to turn down verbose logs
var DebugVerbose = false

var MaxPlayItemNameLength = 26

var JarvoralChannel = "815306102627106836"
var BroadcastChannel = "854733320474329088"
var AppealChannel = "854733587018416138"
var BugChannel = "729467777416691712"

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
	CrowdSize: 25,
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

var QuestMode = false
var StartingRoom = 3
var OocRoom = 2
var ServerShutdown = make(chan bool)

// Load loads the configuration file and overrides the default configuration
// values with any values found.
func init() {

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")
	viper.AddConfigPath(".") // path to look for the config file in
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	Server.Host = viper.GetString("host")
	Server.NEOUname = viper.GetString("neouname")
	Server.NEOPword = viper.GetString("neopword")
	Server.NEOAddress = viper.GetString("neoaddress")
	Server.PGUname = viper.GetString("pguname")
	Server.PGPword = viper.GetString("pgpword")
	Server.PGAddress = viper.GetString("pgaddress")
	Server.Port = viper.GetString("port")
	Server.PGPort = viper.GetInt("pgport")
	Server.RestToken = viper.GetString("resttoken")
	// Setup global logging format
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)

	// Seed default random source

	if !Debug.LongLog {
		log.SetFlags(log.Ldate | log.Ltime)
		log.Printf("Switching to short log format.")
	}
}

var BlockedNames, _ = utils.ReadLines("restricted_names")
