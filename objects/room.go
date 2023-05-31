package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jinzhu/copier"
	"log"
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
	EncounterTable     map[int]int
	roomTicker         *time.Ticker
	roomTickerUnload   chan bool
	effectTicker       *time.Ticker
	effectTickerUnload chan bool
	StoreOwner         string
	StoreInventory     *ItemInventory
	Songs              map[string]string
	LockPriority       string
}

// Pop the room data
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
		nil,
		make(chan bool),
		nil,
		make(chan bool),
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

// Evaluate if there are too many players in this rooms inventory
func (r *Room) Crowded() (crowded bool) {
	if r != nil {
		crowded = len(r.Chars.Contents) >= config.Inventory.CrowdSize
	}
	return
}

// Drop out the description of this room
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
				// Clean up just in case a delete didn't get cleaned up...
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
		// Clean up just in case a delete didn't get cleaned up...
		if _, ok := Rooms[exiti.ToId]; !ok {
			delete(r.Exits, exiti.Name)
		}
	}
}

func (r *Room) FindExit(exitName string) *Exit {
	for k, v := range r.Exits {
		if strings.Contains(strings.ToLower(k), strings.ToLower(exitName)) {
			return v
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

// Actions to be taken on the first person entering a room
func (r *Room) FirstPerson() {
	// Construct and institute the ticker
	//*
	if r.Flags["encounters_on"] {
		r.roomTicker = time.NewTicker(10 * time.Second)
		go func() {
			for {
				select {
				case <-r.roomTickerUnload:
					return
				case <-r.roomTicker.C:
					r.Encounter()
				}
			}
		}()
	}
	if r.Flags["fire"] || r.Flags["earth"] || r.Flags["wind"] || r.Flags["water"] {
		r.effectTicker = time.NewTicker(45 * time.Second)
		go func() {
			for {
				select {
				case <-r.effectTickerUnload:
					return
				case <-r.effectTicker.C:

					if r.Flags["earth"] || r.Flags["fire"] || r.Flags["air"] || r.Flags["water"] {
						for _, c := range r.Chars.Contents {
							if r.Flags["earth"] {
								if !c.Flags["resist_earth"] {
									c.Write([]byte(text.Brown + "The earth swells up around you." + "\n"))
									c.ReceiveMagicDamage(20, "earth")
									c.DeathCheck("was swallowed by the earth.")
								} else {
									c.Write([]byte(text.Brown + "Your earth resistance protects you from the environment." + "\n"))
								}
							} else if r.Flags["fire"] {
								if !c.Flags["resist_fire"] {
									c.Write([]byte(text.Brown + "Burning flames overwhelm you." + "\n"))
									c.ReceiveMagicDamage(20, "fire")
									c.DeathCheck("was burned alived.")
								} else {
									c.Write([]byte(text.Brown + "Your fire resistance protects you from the environment." + "\n"))
								}
							} else if r.Flags["water"] {
								if !c.Flags["resist_water"] {
									c.Write([]byte(text.Brown + "The water overwhelms you, choking you." + "\n"))
									c.DeathCheck("drowned.")
									c.ReceiveMagicDamage(20, "water")
								} else {
									c.Write([]byte(text.Brown + "Your water resistance protects you from the environment." + "\n"))
								}
							} else if r.Flags["air"] {
								if !c.Flags["resist_air"] {
									c.Write([]byte(text.Brown + "The icy air buffets you." + "\n"))
									c.DeathCheck("was frozen solid.")
									c.ReceiveMagicDamage(20, "air")
								} else {
									c.Write([]byte(text.Brown + "Your air protection protects you from the icy winds." + "\n"))
								}
							}
						}
					}
				}
			}
		}()
	}
	r.Mobs.RestartPerms()
}

func (r *Room) Encounter() {
	// Check if encounters are off, a GM can change this live.
	if r.Flags["encounters_on"] {
		log.Println("Room# " + strconv.Itoa(r.RoomId) + " Run the encounter function!")
		if len(r.Mobs.Contents) < 10 {
			// Augment the encounter based on the number of players in the room
			aug := len(r.Chars.Contents)
			if aug <= 1 {
				aug = 0
			}
			// Roll the dice and see if we get a mob here
			if utils.Roll(100, 1, 0) <= r.EncounterRate+(aug*config.MobAugmentPerCharacter) {
				// Successful roll:  Roll again to pick the mob
				mobCalc := 0
				mobPick := utils.Roll(100, 1, 0)
				for mob, chance := range r.EncounterTable {
					if (DayTime && !Mobs[mob].Flags["night_only"]) || (!DayTime && !Mobs[mob].Flags["day_only"]) {
						mobCalc += chance
						if mobPick <= mobCalc {
							// This is the mob!  Put it in the room!
							newMob := Mob{}
							copier.CopyWithOption(&newMob, Mobs[mob], copier.Option{DeepCopy: true})
							newMob.Placement = 5
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

func (r *Room) LastPerson() {
	// Verify that no one else is in here after a follow mob invocation
	if len(r.Chars.Contents) != 0 {
		return
	}

	// This was the last character, invoke the whole cleaning routine now.
	log.Println("Clearing Room: " + r.Name + " (" + strconv.Itoa(r.RoomId) + ")")
	r.Items.RemoveNonPerms()
	go r.Mobs.RemoveNonPerms()

	for _, exit := range r.Exits {
		if exit.Flags["autoclose"] {
			exit.Close()
		}
	}

	// Destruct the ticker
	if r.Flags["encounters_on"] {
		r.roomTickerUnload <- true
		r.roomTicker.Stop()
	}

	// Destruct the ticker
	if r.Flags["fire"] || r.Flags["earth"] || r.Flags["wind"] || r.Flags["water"] {
		r.effectTickerUnload <- true
		r.effectTicker.Stop()
	}

	// Relock all the exits.
	for _, exit := range r.Exits {
		if exit.Flags["lockable"] {
			exit.Flags["locked"] = true
		}
	}

	go r.Save()

}

func (r *Room) MessageAll(Message string) {
	// Message all the characters in this room
	for _, chara := range r.Chars.Contents {
		chara.Write([]byte(Message))
	}
}

func (r *Room) MessageVisible(Message string) {
	// Message all the characters in this room
	for _, chara := range r.Chars.Contents {
		// Check invisible detection
		visDetect, err := chara.Flags["detect_invisible"]
		if !err {
			continue
		}
		if visDetect || chara.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			chara.Write([]byte(Message + "\n"))
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
	r.Lock()
	if o.Flags["invisible"] {
		r.MessageVisible(o.Name + " wanders away. \n" + text.Reset)
	} else if !o.Flags["hidden"] {
		r.MessageAll(o.Name + " wanders away. \n" + text.Reset)
	}
	r.Mobs.Remove(o)
	o = nil
	r.Unlock()
}

func (r *Room) FleeMob(o *Mob) {
	r.Lock()
	if o.Flags["invisible"] {
		r.MessageVisible(o.Name + " flees!! \n" + text.Reset)
	} else if !o.Flags["hidden"] {
		r.MessageAll(o.Name + " flees!! \n" + text.Reset)
	}
	r.Mobs.Remove(o)
	o = nil
	r.Unlock()
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
