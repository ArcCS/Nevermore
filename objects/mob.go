package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jinzhu/copier"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Mob implements a control object for mobs interacting with players and each other
type Mob struct {
	Object
	MobId int
	Inventory *ItemInventory
	ItemList map[int]int
	Flags map[string]bool
	Effects map[string]Effect

	// ParentId is the room id for the room
	ParentId int
	Gold int
	Experience int
	Level int

	Stam Meter
	Mana Meter

	// Attributes
	Str Meter
	Dex Meter
	Con Meter
	Int Meter
	Pie Meter
	Armor int

	// Dice
	NumDice int
	SidesDice int
	PlusDice int

	// Magic
	ChanceCast int
	Spells []string
	WaterResistance int
	AirResistance int
	FireResistance int
	EarthResistance int

	//Threat table attacker -> damage
	//TODO These are crappy short cuts,  we should use a target shim, or direct object pointers
	TotalThreatDamage int
	ThreatTable map[string]int
	CurrentTarget string

	NumWander int
	TicksAlive int
	WimpyValue int

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
		int(mobData["mob_id"].(int64)),
		NewItemInventory(),
		make(map[int]int),
		make(map[string]bool),
		make(map[string]Effect),
		-1,
		int(mobData["gold"].(int64)),
		int(mobData["experience"].(int64)),
		int(mobData["level"].(int64)),
		Meter{int(mobData["hpmax"].(int64)), int(mobData["hpcur"].(int64))},
		Meter{int(mobData["mpmax"].(int64)), int(mobData["mpcur"].(int64))},
		Meter{40, int(mobData["strength"].(int64))},
		Meter{40,int(mobData["dexterity"].(int64))},
		Meter{40,int(mobData["constitution"].(int64))},
		Meter{40,int(mobData["intelligence"].(int64))},
		Meter{40,int(mobData["piety"].(int64))},
		int(mobData["armor"].(int64)),
		int(mobData["ndice"].(int64)),
		int(mobData["sdice"].(int64)),
		int(mobData["pdice"].(int64)),
		int(mobData["casting_probability"].(int64)),
		strings.Split(mobData["spells"].(string), " "),
		int(mobData["water_resistance"].(int64)),
		int(mobData["air_resistance"].(int64)),
		int(mobData["fire_resistance"].(int64)),
		int(mobData["earth_resistance"].(int64)),
		0,
		nil,
		"",
		int(mobData["wimpyvalue"].(int64)),
		0,
		int(mobData["numwander"].(int64)),
		nil,
		nil,
		0,
	}

	for _, drop := range mobData["drops"].([]interface{}) {
		if drop != nil {
			dropData := drop.(map[string]interface{})
			if dropData["chance"] != nil {
				newMob.ItemList[int(dropData["item_id"].(int64))] = int(dropData["chance"].(int64))
			}
		}
	}

	for k, v := range mobData["flags"].(map[string]interface{}){
		if v == nil{
			newMob.Flags[k] = false
		}else {
			newMob.Flags[k] = int(v.(int64)) != 0
		}
	}
	return newMob, true
}

func (m *Mob) StartTicking(){
	m.CalculateInventory()
	m.ThreatTable = make(map[string]int)
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
		m.TicksAlive++
		if m.TicksAlive >= m.NumWander && m.CurrentTarget == "" {
			go Rooms[m.ParentId].WanderMob(m)
			return
		} else {
			Rooms[m.ParentId].Chars.Lock()
			Rooms[m.ParentId].Mobs.Lock()
			Rooms[m.ParentId].Items.Lock()
			// Am I hostile?  Should I pick a target?
			if m.CurrentTarget == "" && m.Flags["hostile"] {
				potentials := Rooms[m.ParentId].Chars.MobList(m.Flags["detect_invisible"], false)
				if len(potentials) > 0 {
					rand.Seed(time.Now().Unix())
					m.CurrentTarget = potentials[rand.Intn(len(potentials))]
					m.AddThreatDamage(1, m.CurrentTarget)
					Rooms[m.ParentId].MessageAll(m.Name + " attacks " + m.CurrentTarget + text.Reset + "\n")
				}
			}

			if m.CurrentTarget != "" {
				if !utils.StringIn(m.CurrentTarget, Rooms[m.ParentId].Chars.MobList(m.Flags["detect_invisible"], true)) {
					m.CurrentTarget = ""
				}
			}

			// Do I want to change targets? 33% chance if the current target isn't the highest on the threat table
			if len(m.ThreatTable) > 1 {
				rankedThreats := utils.RankMapStringInt(m.ThreatTable)
				if m.CurrentTarget != rankedThreats[0] {
					if utils.Roll(3, 1, 0) == 1 {
						if !utils.StringIn(rankedThreats[0], Rooms[m.ParentId].Chars.MobList(m.Flags["detect_invisible"], true)) {
							m.CurrentTarget = rankedThreats[0]
							Rooms[m.ParentId].MessageAll(m.Name + " turns to " + m.CurrentTarget + "\n" + text.Reset)
						}
					}
				}
			}

			// TODO: Do I pick stuff up off the ground?
			log.Println(m.Name + "My target is: :" + m.CurrentTarget )
			if (m.CurrentTarget == "" && m.Placement != 3) ||
				(m.CurrentTarget != "" && !m.Flags["ranged"] &&
					m.Placement != Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement) ||
				(m.CurrentTarget != "" && m.Flags["ranged"] &&
					(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement)) > 1)) {
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
				if target.Class == 0 && target.Equipment.Main != nil && config.RollParry(target.Skills[target.Equipment.Main.ItemType].Value) {
					if target.Tier >= 10 {
						// It's a riposte
						actualDamage,_ := m.ReceiveDamage(int(math.Ceil(float64(target.InflictDamage()))))
						target.Write([]byte(text.Green + "You parry and riposte the attack from " + m.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + "\n" + text.Reset))
						if m.Stam.Current <= 0 {
							Rooms[m.ParentId].MessageAll(text.Green + target.Name + " killed " + m.Name)
							stringExp := strconv.Itoa(m.Experience)
							for k := range m.ThreatTable {
								Rooms[m.ParentId].Chars.Search(k, true).Write([]byte(text.Cyan + "You earn " + stringExp + "exp for the defeat of the " + m.Name + "\n" + text.Reset))
								Rooms[m.ParentId].Chars.Search(k, true).Experience.Add(m.Experience)
							}
							Rooms[m.ParentId].MessageAll(m.Name + " dies.")
							m.DropInventory()
							go Rooms[m.ParentId].ClearMob(m)
							return
						}
						m.MobStunned = config.ParryStuns
					} else {
						target.Write([]byte(text.Green + "You parry the attack from " + m.Name + "\n" + text.Reset))
						m.MobStunned = config.ParryStuns
					}
				} else {
					stamDamage, vitDamage := target.ReceiveDamage(int(math.Ceil(float64(m.InflictDamage()))))
					buildString := ""
					if stamDamage != 0 {
						buildString += strconv.Itoa(stamDamage) + " stamina"
					}
					if stamDamage != 0 && vitDamage != 0 {
						buildString += " and "
					}
					if vitDamage != 0 {
						buildString += strconv.Itoa(vitDamage) + " vitality"
					}
					if stamDamage == 0 && vitDamage == 0 {
						target.Write([]byte(text.Red + m.Name + " attacks bounces off of you for no damage!" + "\n" + text.Reset))
					} else {
						target.Write([]byte(text.Red + m.Name + " attacks you for " + buildString + " points of damage!" + "\n" + text.Reset))
					}
					if target.Vit.Current == 0 {
						target.Died()
					}
				}
			} else if m.CurrentTarget != "" && m.Flags["ranged"] &&
				(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement)) > 1) {
				target := Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false)
				stamDamage, vitDamage := target.ReceiveDamage(int(math.Ceil(float64(m.InflictDamage()))))
				buildString := ""
				if stamDamage != 0 {
					buildString += strconv.Itoa(stamDamage) + " stamina"
				}
				if stamDamage != 0 && vitDamage != 0 {
					buildString += " and "
				}
				if vitDamage != 0 {
					buildString += strconv.Itoa(vitDamage) + " vitality"
				}
				target.Write([]byte(text.Red + "Thwwip!! " + m.Name + " attacks you for " + buildString + " points of damage!" + "\n" + text.Reset))
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
}

// On copy to a room calculate the inventory
func (m *Mob) CalculateInventory() {
	//log.Println("Attempting to add some inventory...")
	if len(m.ItemList) > 0 {
		for k, v := range m.ItemList {
			if utils.Roll(100, 1, 0) <= v {
				log.Println("Adding inventory!!")
				// Successful roll!  Add this item to the inventory!
				newItem := Item{}
				copier.Copy(&newItem, Items[k])
				m.Inventory.Add(&newItem)
			}
		}
	}
}

func (m *Mob) DropInventory(){
	for _, item := range m.Inventory.Contents{
		m.Inventory.Remove(item)
		Rooms[m.ParentId].Items.Add(item)
	}
}

// On death calculate and distribute experience
func (m *Mob) CalculateExperience(attackerName string) {
	return
}

func (m *Mob) AddThreatDamage(damage int, attackerName string){
	m.ThreatTable[attackerName] += damage
	if m.CurrentTarget == "" {
		m.CurrentTarget = attackerName
	}
}

func (m *Mob) ApplyEffect(effect string){
	return
}

func (m *Mob) RemoveEffect(effect string){
	return
}

func (m *Mob) ToggleFlag(flagName string) bool {
	if val, exists := m.Flags[flagName]; exists {
		m.Flags[flagName] = !val
		return true
	} else {
		return false
	}
}

func (m *Mob) ReceiveDamage(damage int) (int, int) {
	//TODO: Review the numbers for armor here
	finalDamage := math.Ceil(float64(damage) * (1 - (float64(m.Armor/config.MobArmorReductionPoints)*config.MobArmorReduction)))
	m.Stam.Subtract(int(finalDamage))
	return int(finalDamage), 0
}

func (m *Mob) ReceiveVitalDamage(damage int) int{
	damageMod, _ := m.ReceiveDamage(damage)
	return damageMod
}

func (m *Mob) Heal(damage int) (int, int){
	m.Stam.Add(damage)
	return damage, 0
}

func (m *Mob) HealStam(damage int){
	m.Stam.Add(damage)
}

func (m *Mob) HealVital(damage int){
	m.Heal(damage)
}

func (m *Mob) RestoreMana(damage int){
	m.Mana.Add(damage)
}

func (m *Mob) InflictDamage() int {
	damage := 0
	if m.NumDice > 0 && m.SidesDice > 0 {
		damage = utils.Roll(m.SidesDice, m.NumDice, m.PlusDice)
	}
	return damage
}

func (m *Mob) CastSpell(spell string) bool {
	return true
}

func (m *Mob) Look() string {
	buildText := "You see a " + m.Name + ", " + config.TextTiers[m.Level] + " level. \n"
	buildText += m.Description
	if m.Flags["hostile"] {
		buildText += "\n It looks hostile!"
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