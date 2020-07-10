package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/prompt"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jinzhu/copier"
	"strconv"
	"strings"
	"time"
)

// Room contains the map of fields for a room in nexus
type Room struct {
	Object

	RoomId int  //room_id from database
	Creator string
	Exits map[string]*Exit
	Mobs *MobInventory
	Chars *CharInventory
	Items *ItemInventory
	Flags map[string]bool
	Commands map[string]prompt.MenuItem
	// This is a whole number percentage out of 100
	EncounterRate int
	// MobID mapped to an encounter percentage
	EncounterTable map[int]int
	roomTicker *time.Ticker
	roomTickerUnload chan bool
}

// Pop the room data
func LoadRoom(roomData map[string]interface{}) (*Room, bool){
	newRoom := &Room{
		Object{
			Name:        roomData["name"].(string),
			Description: roomData["description"].(string),
			Placement:   3,
		},
		int(roomData["room_id"].(int64)),
		roomData["creator"].(string),
		make(map[string]*Exit),
		NewMobInventory(int(roomData["room_id"].(int64))),
		NewCharInventory(int(roomData["room_id"].(int64))),
		NewItemInventory(),
		make(map[string]bool),
		make(map[string]prompt.MenuItem),
		int(roomData["encounter_rate"].(int64)),
		make(map[int]int),
		nil,
		make(chan bool),
	}

	// TODO: Load Permanent Items

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

	for k, v := range roomData["flags"].(map[string]interface{}){
		if v == nil{
			newRoom.Flags[k] = false
		}else {
			newRoom.Flags[k] = int(v.(int64)) != 0
		}
	}
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
func (r *Room) Look(gm bool) (buildText string) {
	invis := ""
	hidden := ""
	inactive := ""
	if !gm {
		buildText += r.Description + "\n"
		if len(r.Exits) > 0 {
			exitText := ""
			for _, exiti := range r.Exits {
				// Clean up just in case a delete didn't get cleaned up...
				if nextRoom, ok := Rooms[exiti.ToId]; !ok{
					delete(r.Exits, exiti.Name)
				}else {
					if exiti.Flags["invisible"] != true &&
						exiti.Flags["hidden"] != true &&
						nextRoom.Flags["active"] == true {
						exitText += exiti.Name + " "
					}
				}
			}
			if exitText != ""{
				buildText += "From here you can go: " + exitText
			}else{
				buildText += "You see no apparent exits."
			}
		} else {
			buildText += "You see no apparent exits."
		}
		return buildText
	}else{
		buildText = text.Cyan + r.Name + " [ID:" + strconv.Itoa(r.RoomId) +  "] (" + r.Creator +")\n" + text.Reset
		buildText += text.Yellow + r.Description + "\n"
		if len(r.Exits) > 0 {
			buildText += "From here you can go: "
			for _, exiti := range r.Exits {
				invis = ""
				hidden = ""
				inactive = ""
				if nextRoom, ok := Rooms[exiti.ToId]; !ok{
					delete(r.Exits, exiti.Name)
				}else {
					if exiti.Flags["invisible"] {
						invis = "[X]"
					}
					if exiti.Flags["hidden"] {
						hidden = "[-]"
					}
					if nextRoom.Flags["active"] == false {
						inactive = "[i]"
					}
					buildText += exiti.Name + " " + hidden + invis + inactive + "(" + strconv.Itoa(exiti.Placement) + ")[ID:" + strconv.Itoa(exiti.ToId) + "], "
				}
			}

		} else {
			buildText += "You see no apparent exits."
		}
		return buildText
	}
}

func (r *Room) CleanExits(){
	for _, exiti := range r.Exits {
		// Clean up just in case a delete didn't get cleaned up...
		if _, ok := Rooms[exiti.ToId]; !ok{
			delete(r.Exits, exiti.Name)
		}
	}
}

func (r *Room) FindExit(exitName string) *Exit {
	for k, v := range r.Exits {
		if strings.Contains(k, exitName){
			return v
		}
	}
	return nil
}

func (r *Room) ToggleFlag(flagName string) bool {
	if val, exists := r.Flags[flagName]; exists{
		r.Flags[flagName] = !val
		return true
	}else{
		return false
	}
}

// Actions to be taken on the first person entering a room
func (r *Room) FirstPerson() {
	// Construct and institute the ticker
	//*
	if r.Flags["encounters_on"] || r.Flags["fire"] || r.Flags["earth"] || r.Flags["wind"] || r.Flags["water"] {
		r.roomTicker = time.NewTicker(8 * time.Second)
		go func() {
			for {
				select {
				case <-r.roomTickerUnload:
					return
				case <-r.roomTicker.C:
					// Is the room crowded?
					if len(r.Mobs.Contents) < 10 {
						// Roll the dice and see if we get a mob here
						if utils.Roll(100, 1, 0) <= r.EncounterRate {
							// Successful roll:  Roll again to pick the mob
							mobCalc := 0
							mobPick := utils.Roll(100, 1, 0)
							for mob, chance := range r.EncounterTable {
								mobCalc += chance
								if mobPick <= mobCalc {
									// This is the mob!  Put it in the room!
									newMob := Mob{}
									copier.Copy(&newMob, Mobs[mob])
									newMob.Placement = 5
									r.Mobs.Add(&newMob)
									newMob.StartTicking()
									break
								}
							}
						}
					}
					//TODO: Do some elemental damage
				}
			}
		}()
	}
	// Resume permanent mob/item tickers
}

func (r *Room) LastPerson(){

	r.Mobs.RemoveNonPerms()

	r.Items.RemoveNonPerms()

	// Destruct the ticker
	if r.Flags["encounters_on"] || r.Flags["fire"] || r.Flags["earth"] || r.Flags["wind"] || r.Flags["water"] {
		r.roomTickerUnload<-true
		r.roomTicker.Stop()
	}

}

func (r *Room) MessageAll(Message string){
	// Message all the characters in this room
	for _, chara := range r.Chars.Contents{
		chara.Write([]byte(Message))
	}
}

func (r *Room) MessageVisible(Message string){
	// Message all the characters in this room
	for _, chara := range r.Chars.Contents{
		// Check invisible detection
		visDetect, err := chara.Flags["detect_invisible"]; if !err {
			continue
		}
		if visDetect || chara.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster){
			chara.Write([]byte(Message + "\n"))
		}
	}
}

func (r *Room) MessageMovement(previous int, new int, subject string){
	// Message all the characters in this room
	for _, chara := range r.Chars.Contents{
		chara.WriteMovement(previous, new, subject)
	}
}

func (r *Room) WanderMob(o *Mob) {
	if o.Flags["invisible"] {
		r.MessageVisible(o.Name + " wanders away.")
	}else if !o.Flags["hidden"] {
		r.MessageAll(o.Name + " wanders away.")
	}
	r.Mobs.Remove(o)
	o.MobTickerUnload <- true
	o = nil
}

func (r *Room) ClearMob(o *Mob) {
	r.Mobs.Remove(o)
	o.MobTickerUnload <- true
	o = nil
}

func (r *Room) Save(){
	roomData := make(map[string]interface{})
	roomData["room_id"] = r.RoomId
	roomData["name"] = r.Name
	roomData["description"] = r.Description
	roomData["repair"] = utils.Btoi(r.Flags["repair"])
	roomData["encounter_rate"] = r.EncounterRate
	roomData["mana_drain"] =  utils.Btoi(r.Flags["mana_drain"])
	roomData["no_summon"] =  utils.Btoi(r.Flags["no_summon"])
	roomData["heal_fast"] =  utils.Btoi(r.Flags["heal_fast"])
	roomData["no_teleport"] =  utils.Btoi(r.Flags["no_teleport"])
	roomData["lo_level"] =  utils.Btoi(r.Flags["lo_level"])
	roomData["no_scry"] =  utils.Btoi(r.Flags["no_scry"])
	roomData["shielded"] =  utils.Btoi(r.Flags["shielded"])
	roomData["dark_always"] =  utils.Btoi(r.Flags["dark_always"])
	roomData["light_always"] =  utils.Btoi(r.Flags["light_always"])
	roomData["natural_light"] =  utils.Btoi(r.Flags["natural_light"])
	roomData["indoors"] =  utils.Btoi(r.Flags["indoors"])
	roomData["fire"] =  utils.Btoi(r.Flags["fire"])
	roomData["encounters_on"] =  utils.Btoi(r.Flags["encounters_on"])
	roomData["no_word_of_recall"] =  utils.Btoi(r.Flags["no_word_of_recall"])
	roomData["water"] =  utils.Btoi(r.Flags["water"])
	roomData["no_magic"] =  utils.Btoi(r.Flags["no_magic"])
	roomData["urban"] =  utils.Btoi(r.Flags["urban"])
	roomData["underground"] =  utils.Btoi(r.Flags["underground"])
	roomData["hilevel"] =  utils.Btoi(r.Flags["hilevel"])
	roomData["earth"] =  utils.Btoi(r.Flags["earth"])
	roomData["wind"] =  utils.Btoi(r.Flags["wind"])
	roomData["active"] = utils.Btoi(r.Flags["active"])
	roomData["train"] = utils.Btoi(r.Flags["train"])
	data.UpdateRoom(roomData)

	// TODO: Save permanent items
	// TODO: save permanent mobs
}