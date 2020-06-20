package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
	"strings"
	"time"
)

// Mob implements a control object for mobs interacting with players and each other
type Mob struct {
	Object
	MobId int64
	Inventory *ItemInventory
	ItemList map[int64]int64
	Flags map[string]bool
	Effects map[string]Effect

	// ParentId is the room id for the room
	ParentId int64
	Gold int64
	Experience int64
	Level int64

	Stam Meter
	Mana Meter

	// Attributes
	Str Meter
	Dex Meter
	Con Meter
	Int Meter
	Pie Meter
	Armor int64

	// Dice
	NumDice int64
	SidesDice int64
	PlusDice int64

	// Magic
	ChanceCast int64
	Spells []string
	WaterResistance int64
	AirResistance int64
	FireResistance int64
	EarthResistance int64

	//Threat table attacker -> damage
	TotalThreatDamage Accumulator
	ThreatTable map[string]Accumulator
	CurrentTarget string

	NumWander int64
	TicksAlive int64
	WimpyValue int64

	MobTickerUnload chan bool
	MobTicker *time.Ticker
}

// Pop the mob data
func LoadMob(mobData map[string]interface{}) (*Mob, bool){
	description := ""
	var ok bool
	if description, ok = mobData["description"].(string); !ok {
		description = "A mob...  yup"
	}
	newMob := &Mob{
		Object{
			Name:        mobData["name"].(string),
			Description: description,
			Placement:   5,
		},
		mobData["mob_id"].(int64),
		NewItemInventory(),
		make(map[int64]int64),
		make(map[string]bool),
		make(map[string]Effect),
		-1,
		mobData["gold"].(int64),
		mobData["experience"].(int64),
		mobData["level"].(int64),
		Meter{mobData["hpmax"].(int64), mobData["hpcur"].(int64)},
		Meter{mobData["mpmax"].(int64), mobData["mpcur"].(int64)},
		Meter{40, mobData["strength"].(int64)},
		Meter{40,mobData["dexterity"].(int64)},
		Meter{40,mobData["constitution"].(int64)},
		Meter{40,mobData["intelligence"].(int64)},
		Meter{40,mobData["piety"].(int64)},
		mobData["armor"].(int64),
		mobData["ndice"].(int64),
		mobData["sdice"].(int64),
		mobData["pdice"].(int64),
		mobData["casting_probability"].(int64),
		strings.Split(mobData["spells"].(string), " "),
		mobData["water_resistance"].(int64),
		mobData["air_resistance"].(int64),
		mobData["fire_resistance"].(int64),
		mobData["earth_resistance"].(int64),
		Accumulator{0},
		make(map[string]Accumulator),
		"",
		mobData["wimpyvalue"].(int64),
		0,
		mobData["numwander"].(int64),
		nil,
		nil,
	}

	for _, drop := range mobData["drops"].([]interface{}) {
		if drop != nil {
			dropData := drop.(map[string]interface{})
			if dropData["dest"] != nil {
				newMob.ItemList[dropData["item_id"].(int64)] = dropData["chance"].(int64)
			}
		}
	}

	for k, v := range mobData["flags"].(map[string]interface{}){
		if v == nil{
			newMob.Flags[k] = false
		}else {
			newMob.Flags[k] = v.(int64) != 0
		}
	}
	return newMob, true
}

func (m *Mob) StartTicking(){
	m.MobTickerUnload = make(chan bool)
	//TODO Modify this ticker if the mob is especially fast moving
	m.MobTicker = time.NewTicker(8 * time.Second)
	go func() {
		for {
			select {
			case <-m.MobTickerUnload:
				return
			case <-m.MobTicker.C:
				m.Tick()
			}
		}
	}()
}

// The mob brain is this ticker
func (m *Mob) Tick(){
	//log.Println(m.Name + " did a tick!")
	m.TicksAlive++
	if m.TicksAlive >= m.NumWander && m.CurrentTarget == "" {
		Rooms[m.ParentId].WanderMob(m)
	}
	// Am I hostile?  Should I pick a target?
	if m.CurrentTarget == "" && !m.Flags["hostile"] {
		// Randomly choose someone from this room to be hostile towards

	}
	// todo: Threat table management

	// Perform movement/action/combat/stealing
	// I have no target and want to move
	if m.CurrentTarget == "" && m.Placement != 3 {
		// We aren't fighting, we don't want to fight, and we aren't in the middle of the room.  Lets get there.
		oldPlacement := m.Placement
		if m.Placement > 3 {
			m.Placement--
		}else{
			m.Placement++
		}
		if !m.Flags["hidden"] {
			whichNumber := Rooms[m.ParentId].Mobs.GetNumber(m)
			Rooms[m.ParentId].MessageMovement(oldPlacement, m.Placement, m.Name + " #" + strconv.Itoa(whichNumber))
		}
	}
	// Am I mad at something?
	//if m.CurrentTarget != "" {
		// Is that target near enough to me?
	//	target := Rooms[m.ParentId].Chars.Search()
	//}
	  // Can I cast spells or ranged hit?
	 // Am I in the center of the room?

	// Combat
	// Am I changing targets because of threat table generation?
	// Am I in range to attack?
	  // Spells
	  // Ranged
	  // Melee

	// Do I pick stuff up off the ground?
	  // Yoink!

	// Do I want to wander away?

}

// On copy to a room calculate the inventory
func (m *Mob) CalculateInventory(){
	return
}

// On death calculate and distribute experience
func (m *Mob) CalculateExperience(attackerName string){
	return
}

func (m *Mob) AddThreatDamage(damage string, attackerName string){
	return
}

func (m *Mob) ApplyEffect(){
	return
}

func (m *Mob) RemoveEffect(effect string){
	return
}


func (m *Mob) ReceiveDamage(damage int){
	return
}

func (m *Mob) ReceiveVitalDamage(damage int){

}

func (m *Mob) Heal(damage int){
	return
}

func (m *Mob) HealVital(damage int){

}

func (m *Mob) RestoreMana(damage int){

}

func (m *Mob) InflictDamage() (damage int){
	return
}

func (m *Mob) CastSpell(spell string) bool {
	return true
}

func (m *Mob) Died() {

}

func (m *Mob) Look() string {
	buildText := "You see a " + m.Name + ", " + config.TextTiers[m.Level] + " level. \n"
	buildText += m.Description + "\n"
	/*
	TODO: Location He is standing a couple steps in front of you.
	TODO: Hostile He looks hostile!
	TODO: ThreatTable He looks very angry at you.
	TODO: Who attacking He is attacking you.

	 */
	return buildText
}

func (m *Mob) Save() {
	mobData := make(map[string]interface{})
	mobData["mob_id"] = m.MobId 
	mobData["name"] = m.Name 
	mobData["description"]= m.Description 
	mobData["experience"]= m.Experience 
	mobData["level"]= m.Level 
	mobData["gold"]= m.Gold 
	mobData["constitution"]= m.Con.Current 
	mobData["strength"]= m.Str.Current 
	mobData["intelligence"]= m.Int.Current 
	mobData["dexterity"]= m.Dex.Current 
	mobData["piety"]= m.Pie.Current 
	mobData["mpmax"]= m.Mana.Max 
	mobData["mpcur"]= m.Mana.Current 
	mobData["hpcur"]= m.Stam.Current 
	mobData["hpmax"]= m.Stam.Max 
	mobData["sdice"]= m.SidesDice 
	mobData["ndice"]= m.NumDice 
	mobData["pdice"]= m.PlusDice 
	mobData["spells"] = strings.Join(m.Spells, ",")
	mobData["casting_probability"] = m.ChanceCast 
	mobData["armor"]= m.Armor
	mobData["numwander"]= m.NumWander
	mobData["wimpyvalue"]= m.WimpyValue
	mobData["air_resistance"]= m.AirResistance
	mobData["fire_resistance"]= m.FireResistance
	mobData["earth_resistance"]= m.EarthResistance
	mobData["water_resistance"]= m.WaterResistance
	mobData["hide_encounter"]= utils.Btoi(m.Flags["hide_encounter"])
	mobData["invisible"]= utils.Btoi(m.Flags["invisible"])
	mobData["permanent"]= utils.Btoi(m.Flags["permanent"])
	mobData["hostile"]= utils.Btoi(m.Flags["hostile"])
	data.UpdateMob(mobData)
}