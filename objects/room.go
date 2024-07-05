package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jinzhu/copier"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Room contains the map of fields for a room in nexus
type Room struct {
	Object
	sync.Mutex

	RoomId  int //room_id from database
	Creator string
	Exits   map[string]*Exit
	Mobs    *MobInventory
	Chars   *CharInventory
	Items   *ItemInventory
	Flags   map[string]bool
	// This is a whole number percentage out of 100
	EncounterRate int
	// MobID mapped to an encounter percentage
	EncounterTable    map[int]int
	EvacuateTime      time.Time
	LastEffectTime    time.Time
	LastEncounterTime time.Time
	EncounterSpeed    int
	EncounterJigger   int
	StoreOwner        string
	StoreInventory    *ItemInventory
	Songs             map[string]string
	LockPriority      string
}

var (
	ActivateRoom func(roomId int)
)

func LoadRoom(roomData map[string]interface{}) (*Room, bool) {
	newRoom := &Room{
		Object{
			Name:        roomData["name"].(string),
			Description: roomData["description"].(string),
			Placement:   3,
			Commands:    DeserializeCommands(roomData["commands"].(string)),
		},
		sync.Mutex{},
		int(roomData["room_id"].(int64)),
		roomData["creator"].(string),
		make(map[string]*Exit),
		RestoreMobs(int(roomData["room_id"].(int64)), roomData["mobs"].(string)),
		NewCharInventory(int(roomData["room_id"].(int64))),
		RestoreInventory(roomData["inventory"].(string)),
		make(map[string]bool),
		int(roomData["encounter_rate"].(int64)),
		make(map[int]int),
		time.Time{},
		time.Time{},
		time.Time{},
		//int(roomData["encounter_speed"].(int64)),
		config.RoomDefaultEncounterSpeed,
		0,
		roomData["store_owner"].(string),
		RestoreInventory(roomData["store_inventory"].(string)),
		make(map[string]string),
		"",
	}

	for _, encounter := range roomData["encounters"].([]interface{}) {
		if encounter != nil {
			encData := encounter.(map[string]interface{})
			if encData["chance"] != nil {
				newRoom.EncounterTable[int(encData["mob_id"].(int64))] = int(encData["chance"].(int64))
			}
		}
	}

	for _, exit := range roomData["exits"].([]interface{}) {
		if exit != nil {
			exitData := exit.(map[string]interface{})
			if exitData["dest"] != nil {
				newRoom.Exits[strings.ToLower(exitData["direction"].(string))] = NewExit(int(roomData["room_id"].(int64)), exitData)
			}
		}
	}

	for k, v := range roomData["flags"].(map[string]interface{}) {
		if v == nil {
			newRoom.Flags[k] = false
		} else {
			newRoom.Flags[k] = int(v.(int64)) != 0
		}
	}
	newRoom.Mobs.ContinueEmpty = newRoom.ContinueEmpty
	return newRoom, true
}

// Crowded Evaluate if there are too many players in this rooms inventory
func (r *Room) Crowded() (crowded bool) {
	if r != nil {
		crowded = len(r.Chars.Contents) >= config.Inventory.CrowdSize
	}
	return
}

func (r *Room) LockRoom(lockRequester string, suppress bool) {
	if config.DebugVerbose && !suppress {
		log.Println("Locking Room: " + r.Name + " (" + strconv.Itoa(r.RoomId) + ")" + " by " + lockRequester)
	}
	r.Lock()
	if config.DebugVerbose && !suppress {
		log.Println("Success!! - Locking Room: " + r.Name + " (" + strconv.Itoa(r.RoomId) + ")" + " by " + lockRequester)
	}
}

func (r *Room) UnlockRoom(lockRequester string, suppress bool) {
	if config.DebugVerbose && !suppress {
		log.Println("Unlocking Room: " + r.Name + " (" + strconv.Itoa(r.RoomId) + ")" + " by " + lockRequester)
	}
	r.Unlock()
	if config.DebugVerbose && !suppress {
		log.Println("Success!! - Unlocking Room: " + r.Name + " (" + strconv.Itoa(r.RoomId) + ")" + " by " + lockRequester)
	}
}

// Look Drop out the description of this room
func (r *Room) Look(looker *Character) (buildText string) {
	invis := ""
	hidden := ""
	inactive := ""
	if !looker.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		buildText += r.Description + "\n" + text.Turquoise
		if len(r.Exits) > 0 {
			exitText := make([]string, 0)
			longExit := make([]string, 0)
			for _, exiti := range r.Exits {
				// Clean up just in case delete didn't get cleaned up...
				if nextRoom, ok := Rooms[exiti.ToId]; !ok {
					delete(r.Exits, exiti.Name)
				} else {
					if exiti.Flags["invisible"] != true &&
						exiti.Flags["hidden"] != true &&
						nextRoom.Flags["active"] == true {
						if len(strings.Split(exiti.Name, " ")) >= 3 {
							longExit = append(longExit, "You see "+exiti.Name)
						} else {
							exitText = append(exitText, exiti.Name)
						}
					}
				}
			}
			sort.Strings(exitText)
			sort.Strings(longExit)
			if len(exitText) > 0 {
				buildText += "Obvious exits are " + strings.Join(exitText, ", ")
			}
			if len(longExit) > 0 {
				if len(exitText) > 0 {
					buildText += "\n"
				}
				buildText += strings.Join(longExit, "\n")
			}
			if len(longExit) == 0 && len(exitText) == 0 {
				buildText += "You see no apparent exits."
			}
		} else {
			buildText += "You see no apparent exits."
		}
		return buildText + text.Reset
	} else {
		buildText = text.Cyan + r.Name + " [ID:" + strconv.Itoa(r.RoomId) + "] (" + r.Creator + ")\n" + text.Reset
		buildText += text.Yellow + r.Description + "\n" + text.Turquoise
		exitText := make([]string, 0)
		longExit := make([]string, 0)
		if len(r.Exits) > 0 {
			for _, exiti := range r.Exits {
				invis = ""
				hidden = ""
				inactive = ""
				if nextRoom, ok := Rooms[exiti.ToId]; !ok {
					delete(r.Exits, exiti.Name)
				} else {
					if exiti.Flags["invisible"] {
						invis = "[X]"
					}
					if exiti.Flags["hidden"] {
						hidden = "[-]"
					}
					if nextRoom.Flags["active"] == false {
						inactive = "[i]"
					}
					if len(strings.Split(exiti.Name, " ")) >= 3 {
						longExit = append(longExit, "You see "+exiti.Name+" "+hidden+invis+inactive+"("+strconv.Itoa(exiti.Placement)+")[ID:"+strconv.Itoa(exiti.ToId)+"]")
					} else {
						exitText = append(exitText, exiti.Name+" "+hidden+invis+inactive+"("+strconv.Itoa(exiti.Placement)+")[ID:"+strconv.Itoa(exiti.ToId)+"]")
					}
				}
			}

			sort.Strings(exitText)
			sort.Strings(longExit)

			if len(exitText) > 0 {
				buildText += "Obvious exits are " + strings.Join(exitText, ", ")
			}
			if len(longExit) > 0 {
				if len(exitText) > 0 {
					buildText += "\n"
				}
				buildText += strings.Join(longExit, "\n")
			}
			if len(longExit) == 0 && len(exitText) == 0 {
				buildText += "You see no apparent exits."
			}
		} else {
			buildText += "You see no apparent exits."
		}
		return buildText + text.Reset
	}
}

func (r *Room) ContinueEmpty() bool {
	if len(r.Chars.Contents) == 0 {
		return true
	}
	return false
}

func (r *Room) CleanExits() {
	for _, exiti := range r.Exits {
		if _, ok := Rooms[exiti.ToId]; !ok {
			delete(r.Exits, exiti.Name)
		}
	}
}

func (r *Room) FindExit(exitName string, observer *Character) *Exit {
	for k, v := range r.Exits {
		if strings.Contains(strings.ToLower(k), strings.ToLower(exitName)) {
			if _, ok := Rooms[v.ToId]; ok {
				if Rooms[v.ToId].Flags["active"] || observer.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					return v
				}
			}
		}
	}
	return nil
}

func (r *Room) ToggleFlag(flagName string) bool {
	if val, exists := r.Flags[flagName]; exists {
		r.Flags[flagName] = !val
		return true
	} else {
		return false
	}
}

// FirstPerson Actions to be taken on the first person entering a room
func (r *Room) FirstPerson() {
	// Construct and institute the ticker
	ActivateRoom(r.RoomId)
	r.EvacuateTime = time.Time{}
	r.LastEffectTime = time.Now()
	r.LastEncounterTime = time.Now()
	r.Mobs.RestartPerms()
}

func (r *Room) Encounter() {
	// Check if encounters are off, a GM can change this live.
	if r.Flags["encounters_on"] {
		log.Println("Room# " + strconv.Itoa(r.RoomId) + " Run the encounter function!")
		r.LastEncounterTime = time.Now()
		r.EncounterJigger = utils.Roll(config.RoomMaxJigger, 1, 0)
		if len(r.Mobs.Contents) < 10 {
			// Augment the encounter based on the number of players in the room
			aug := len(r.Chars.Contents)
			if aug <= 1 {
				aug = 0
			}
			// Roll the dice and see if we get a mob here
			if utils.Roll(100, 1, 0) <= r.EncounterRate+(aug*config.MobAugmentPerCharacter) {
				// Successful roll:  Roll again to pick the mob
				multMob := 1
				doubleChance := 0
				if len(r.Chars.Contents) >= 5 {
					doubleChance = 20
				} else if len(r.Chars.Contents) == 4 {
					doubleChance = 15
				} else if len(r.Chars.Contents) == 3 {
					doubleChance = 10
				}
				if utils.Roll(100, 1, 0) <= doubleChance && len(r.Mobs.ListHostile()) <= config.RoomEncNoDoubles {
					multMob = 2
				}
				for i := 0; i < multMob; i++ {
					mobCalc := 0
					mobPick := utils.Roll(100, 1, 0)
					for mob, chance := range r.EncounterTable {
						if (DayTime && !Mobs[mob].Flags["night_only"]) || (!DayTime && !Mobs[mob].Flags["day_only"]) {
							mobCalc += chance
							if mobPick <= mobCalc {
								// This is the mob!  Put it in the room!
								newMob := Mob{}
								if err := copier.CopyWithOption(&newMob, Mobs[mob], copier.Option{DeepCopy: true}); err != nil {
									log.Println("Error copying mob during encounter: ", err)
								}
								if newMob.Placement <= 0 {
									newMob.Placement = 5
								} else if newMob.Placement >= 6 {
									newMob.Placement = utils.Roll(5, 1, 0)
								}
								r.Mobs.Add(&newMob, false)
								newMob.StartTicking()
								break
							}
						}
					}
				}
			}
		}
	}
}

func (r *Room) ElementalDamage() {
	r.LastEffectTime = time.Now()
	for _, c := range r.Chars.Contents {
		if !c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			envDamage := math.Ceil(.20 * float64(c.Stam.Max+c.Vit.Max))
			if r.Flags["earth"] {
				if !c.Flags["resist-earth"] {
					if _, err := c.Write([]byte(text.Brown + "The earth swells up around you." + "\n")); err != nil {
						log.Println("Error writing to player:", err)
					}
					c.ReceiveEnvironmentalDamage(int(envDamage), "earth")
					c.DeathCheck("was swallowed by the earth.")
				} else {
					if _, err := c.Write([]byte(text.Brown + "Your earth resistance protects you from the environment." + "\n")); err != nil {
						log.Println("Error writing to player:", err)
					}
				}
			} else if r.Flags["fire"] {
				if !c.Flags["resist-fire"] {
					if _, err := c.Write([]byte(text.Brown + "Burning flames overwhelm you." + "\n")); err != nil {
						log.Println("Error writing to player:", err)
					}
					c.ReceiveEnvironmentalDamage(int(envDamage), "fire")
					c.DeathCheck("was burned alive.")
				} else {
					if _, err := c.Write([]byte(text.Brown + "Your fire resistance protects you from the environment." + "\n")); err != nil {
						log.Println("Error writing to player:", err)
					}
				}
			} else if r.Flags["water"] {
				if !c.Flags["resist-water"] {
					if _, err := c.Write([]byte(text.Brown + "The water overwhelms you, choking you." + "\n")); err != nil {
						log.Println("Error writing to player:", err)
					}
					c.ReceiveEnvironmentalDamage(int(envDamage), "water")
					c.DeathCheck("drowned.")
				} else {
					if _, err := c.Write([]byte(text.Brown + "Your water resistance protects you from the environment." + "\n")); err != nil {
						log.Println("Error writing to player:", err)
					}
				}
			} else if r.Flags["air"] {
				if !c.Flags["resist-air"] {
					if _, err := c.Write([]byte(text.Brown + "The icy air buffets you." + "\n")); err != nil {
						log.Println("Error writing to player:", err)
					}
					c.ReceiveEnvironmentalDamage(int(envDamage), "air")
					c.DeathCheck("was frozen solid.")
				} else {
					if _, err := c.Write([]byte(text.Brown + "Your air protection protects you from the icy winds." + "\n")); err != nil {
						log.Println("Error writing to player:", err)
					}
				}
			}
		}
	}
}

func (r *Room) LastPerson() {
	// Set the last person time to now
	r.EvacuateTime = time.Now()
}

func (r *Room) CleanRoom() {
	// Verify that no one else is in here after a follow mob invocation

	// This was the last character, invoke the whole cleaning routine now.
	//log.Println("Clearing Room: " + r.Name + " (" + strconv.Itoa(r.RoomId) + ")")
	r.Items.RemoveNonPerms()
	go r.Mobs.RemoveNonPerms()

	for _, exit := range r.Exits {
		if exit.Flags["autoclose"] {
			exit.Close()
		}
	}

	// Relock all the exits.
	for _, exit := range r.Exits {
		if exit.Flags["lockable"] {
			exit.Flags["locked"] = true
		}
	}

	AddRoomUpdate(r.RoomId)

}

func (r *Room) MessageAll(Message string) {
	// Message all the characters in this room
	for _, chara := range r.Chars.Contents {
		if _, err := chara.Write([]byte(Message)); err != nil {
			log.Println("Error writing to player:", err)
		}
	}
}

func (r *Room) MessageVisible(Message string) {
	// Message all the characters in this room
	for _, chara := range r.Chars.Contents {
		// Check invisible detection
		visDetect, err := chara.Flags["detect-invisible"]
		if !err {
			continue
		}
		if visDetect || chara.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			if _, err := chara.Write([]byte(Message + "\n")); err != nil {
				log.Println("Error writing to player:", err)
			}
		}
	}
}

func (r *Room) MessageMovement(previous int, new int, subject string) {
	// Message all the characters in this room
	for _, chara := range r.Chars.Contents {
		chara.WriteMovement(previous, new, subject)
	}
}

func (r *Room) WanderMob(o *Mob) {
	r.LockRoom(strconv.Itoa(r.RoomId)+":WanderMob:"+o.Name, false)
	if o.Flags["invisible"] {
		r.MessageVisible(o.Name + " wanders away. \n" + text.Reset)
	} else if !o.Flags["hidden"] {
		r.MessageAll(o.Name + " wanders away. \n" + text.Reset)
	}
	r.Mobs.Remove(o)
	o = nil
	r.UnlockRoom(strconv.Itoa(r.RoomId)+":WanderMob", false)
}

func (r *Room) FleeMob(o *Mob) {
	r.LockRoom(strconv.Itoa(r.RoomId)+":FleeMob:"+o.Name, false)
	if o.Flags["invisible"] {
		r.MessageVisible(o.Name + " flees!! \n" + text.Reset)
	} else if !o.Flags["hidden"] {
		r.MessageAll(o.Name + " flees!! \n" + text.Reset)
	}
	r.Mobs.Remove(o)
	o = nil
	r.UnlockRoom(strconv.Itoa(r.RoomId)+":FleeMob", false)
}

func (r *Room) ClearMob(o *Mob) {
	r.Mobs.Remove(o)
	o = nil
}

func (r *Room) AddStoreItem(item *Item, price int, infinite bool) {
	if infinite {
		item.Flags["infinite"] = true
	}
	item.StorePrice = price
	r.StoreInventory.Add(item)
}

func (r *Room) SongPlaying(songName string) bool {
	songPlaying := false
	if _, ok := Songs[songName]; ok {
		for _, song := range r.Songs {
			if song == songName {
				songPlaying = true
			}
		}
	}
	return songPlaying
}

func (r *Room) Save() {
	r.Mobs.Jsonify()
	roomData := make(map[string]interface{})
	roomData["room_id"] = r.RoomId
	roomData["name"] = r.Name
	roomData["description"] = r.Description
	roomData["repair"] = utils.Btoi(r.Flags["repair"])
	roomData["encounter_rate"] = r.EncounterRate
	roomData["mana_drain"] = utils.Btoi(r.Flags["mana_drain"])
	roomData["no_summon"] = utils.Btoi(r.Flags["no_summon"])
	roomData["heal_fast"] = utils.Btoi(r.Flags["heal_fast"])
	roomData["no_teleport"] = utils.Btoi(r.Flags["no_teleport"])
	roomData["lo_level"] = utils.Btoi(r.Flags["lo_level"])
	roomData["no_scry"] = utils.Btoi(r.Flags["no_scry"])
	roomData["shielded"] = utils.Btoi(r.Flags["shielded"])
	roomData["dark_always"] = utils.Btoi(r.Flags["dark_always"])
	roomData["light_always"] = utils.Btoi(r.Flags["light_always"])
	roomData["natural_light"] = utils.Btoi(r.Flags["natural_light"])
	roomData["indoors"] = utils.Btoi(r.Flags["indoors"])
	roomData["fire"] = utils.Btoi(r.Flags["fire"])
	roomData["encounters_on"] = utils.Btoi(r.Flags["encounters_on"])
	roomData["no_word_of_recall"] = utils.Btoi(r.Flags["no_word_of_recall"])
	roomData["water"] = utils.Btoi(r.Flags["water"])
	roomData["no_magic"] = utils.Btoi(r.Flags["no_magic"])
	roomData["urban"] = utils.Btoi(r.Flags["urban"])
	roomData["underground"] = utils.Btoi(r.Flags["underground"])
	roomData["hilevel"] = utils.Btoi(r.Flags["hilevel"])
	roomData["earth"] = utils.Btoi(r.Flags["earth"])
	roomData["wind"] = utils.Btoi(r.Flags["wind"])
	roomData["active"] = utils.Btoi(r.Flags["active"])
	roomData["train"] = utils.Btoi(r.Flags["train"])
	roomData["mobs"] = r.Mobs.JsonRepr
	roomData["inventory"] = r.Items.Jsonify()
	roomData["commands"] = r.SerializeCommands()
	roomData["store_owner"] = r.StoreOwner
	roomData["store_inventory"] = r.StoreInventory.Jsonify()
	data.UpdateRoom(roomData)
}
