package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"math"
	"math/rand"
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
	//TODO These are crappy short cuts,  we should use a target shim, or direct object pointers
	TotalThreatDamage int
	ThreatTable map[string]int
	CurrentTarget string

	NumWander int64
	TicksAlive int64
	WimpyValue int64

	MobTickerUnload chan bool
	MobTicker *time.Ticker
	// An int to hold a stun time.
	MobStunned int
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
		0,
		make(map[string]int),
		"",
		mobData["wimpyvalue"].(int64),
		0,
		mobData["numwander"].(int64),
		nil,
		nil,
		0,
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
	if m.MobStunned > 0 {
		m.MobStunned -= 1
	}else {
		// We're kind of managing our own state...  set all the locks
		Rooms[m.ParentId].Chars.Lock()
		Rooms[m.ParentId].Mobs.Lock()
		Rooms[m.ParentId].Items.Lock()

		m.TicksAlive++
		if m.TicksAlive >= m.NumWander && m.CurrentTarget == "" {
			Rooms[m.ParentId].WanderMob(m)
		}

		// Make sure the current target is still in the room and didn't flee
		if !utils.StringIn(m.CurrentTarget, Rooms[m.ParentId].Chars.List(m.Flags["detect_invisible"], true,  "", false)){
			potentials := Rooms[m.ParentId].Chars.List(m.Flags["detect_invisible"], false,"", false)
			if len(potentials) > 0 {
				rand.Seed(time.Now().Unix())
				m.CurrentTarget = potentials[rand.Intn(len(potentials))]
				Rooms[m.ParentId].MessageAll(m.Name + " turns to " + m.CurrentTarget)
			}
		}

		// Am I hostile?  Should I pick a target?
		if m.CurrentTarget == "" && m.Flags["hostile"] {
			potentials := Rooms[m.ParentId].Chars.List(m.Flags["detect_invisible"], false,"", false)
			if len(potentials) > 0 {
				rand.Seed(time.Now().Unix())
				m.CurrentTarget = potentials[rand.Intn(len(potentials))]
				Rooms[m.ParentId].MessageAll(m.Name + " attacks " + m.CurrentTarget)
			}
		}

		// Do I want to chang targets? 33% chance if the current target isn't the highest on the threat table
		if len(m.ThreatTable) > 1 {
			rankedThreats := utils.RankMapStringInt(m.ThreatTable)
			if m.CurrentTarget != rankedThreats[0] {
				if utils.Roll(3, 1, 0) == 1 {
					m.CurrentTarget = rankedThreats[0]
					Rooms[m.ParentId].MessageAll(m.Name + " turns to " + m.CurrentTarget)
				}
			}
		}

		// TODO: Do I pick stuff up off the ground?

		// I have no target and want to move
		if (m.CurrentTarget == "" && m.Placement != 3) ||
			(m.CurrentTarget != "" && !m.Flags["ranged"] &&
				m.Placement != Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement) ||
			(m.CurrentTarget != "" && m.Flags["ranged"] &&
				(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement)) > 1)) {
			// We aren't fighting, we don't want to fight, and we aren't in the middle of the room.  Lets get there.
			oldPlacement := m.Placement
			if m.Placement > 3 {
				m.Placement--
			} else {
				m.Placement++
			}
			if !m.Flags["hidden"] {
				whichNumber := Rooms[m.ParentId].Mobs.GetNumber(m)
				Rooms[m.ParentId].MessageMovement(oldPlacement, m.Placement, m.Name+" #"+strconv.Itoa(whichNumber))
			}
		// Next to attack
		} else if m.CurrentTarget != "" && !m.Flags["ranged"] &&
			m.Placement == Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement {
			// Am I against a fighter and they succeed in a parry roll?
			target := Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false)
			if target.Class == 0 && config.RollParry(int(target.Skills[int(target.Equipment.Main.Type)].Value)) {
				if target.Tier >= 10 {
					// It's a riposte
					actualDamage := m.ReceiveDamage(int64(math.Ceil(float64(target.InflictDamage()))))
					target.Write([]byte(text.Green + "You parry and riposte the attack from " + m.Name + " for " + strconv.Itoa(actualDamage) + " damage!"))
					if m.Stam.Current <= 0 {
						Rooms[m.ParentId].MessageAll(text.Green + target.Name + " killed " + m.Name)
						m.Died()
						Rooms[m.ParentId].ClearMob(m)
					}
					m.MobStunned = config.ParryStuns
				}else{
					target.Write([]byte(text.Green + "You parry the attack from " + m.Name))
					m.MobStunned = config.ParryStuns
				}
			}else{
				actualDamage := target.ReceiveDamage(int(math.Ceil(float64(m.InflictDamage()))))
				target.Write([]byte(text.Red + m.Name + " attacks you for " + strconv.Itoa(actualDamage) + " damage!"))
				if target.Vit.Current == 0 {
					target.Died()
				}
			}
		}else if m.CurrentTarget != "" && m.Flags["ranged"] &&
				(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement)) > 1){
			target := Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false)
			actualDamage := target.ReceiveDamage(int(math.Ceil(float64(m.InflictDamage()))))
			target.Write([]byte(text.Red + "Thwwip!! " + m.Name + " attacks you for " + strconv.Itoa(actualDamage) + " damage!"))
			if target.Vit.Current == 0 {
				target.Died()
			}
		}

		// TODO: Can I cast spells.

		Rooms[m.ParentId].Chars.Unlock()
		Rooms[m.ParentId].Mobs.Unlock()
		Rooms[m.ParentId].Items.Unlock()
	}
}

// On copy to a room calculate the inventory
func (m *Mob) CalculateInventory(){
	return
}

// On death calculate and distribute experience
func (m *Mob) CalculateExperience(attackerName string){
	return
}

func (m *Mob) AddThreatDamage(damage int, attackerName string){
	m.ThreatTable[attackerName] += damage
}

func (m *Mob) ApplyEffect(){
	return
}

func (m *Mob) RemoveEffect(effect string){
	return
}


func (m *Mob) ReceiveDamage(damage int64) int {
	//TODO: Review the numbers for armor here
	finalDamage := math.Ceil(float64(damage) * (1 - (float64(int(m.Armor)/config.MobArmorReductionPoints)*config.MobArmorReduction)))
	m.Stam.Subtract(int64(finalDamage))
	return int(finalDamage)
}

func (m *Mob) ReceiveVitalDamage(damage int64){
	m.ReceiveDamage(damage)
}

func (m *Mob) Heal(damage int){
	m.Stam.Add(int64(damage))
}

func (m *Mob) HealVital(damage int){
	m.Heal(damage)
}

func (m *Mob) RestoreMana(damage int){
	m.Mana.Add(int64(damage))
}

func (m *Mob) InflictDamage() int {
	return utils.Roll(int(m.SidesDice), int(m.NumDice), int(m.PlusDice))
}

func (m *Mob) CastSpell(spell string) bool {
	return true
}

func (m *Mob) Died() {
	Rooms[m.ParentId].MessageAll(m.Name + "dies.")
	stringExp := strconv.Itoa(int(m.Experience))
	for k, _ := range m.ThreatTable {
		Rooms[m.ParentId].Chars.Search(k, true).Write([]byte(text.Blue + "You earn " + stringExp + " for the defeat of the " + m.Name))
		Rooms[m.ParentId].Chars.Search(k, true).Experience.Add(m.Experience)
	}
}

func (m *Mob) Look() string {
	buildText := "You see a " + m.Name + ", " + config.TextTiers[m.Level] + " level. \n"
	buildText += m.Description + "\n"
	if m.Flags["hostile"] {
		buildText += "It looks hostile!"
	}
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