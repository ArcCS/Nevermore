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
	MobId     int
	Inventory *ItemInventory
	ItemList  map[int]int
	Flags     map[string]bool
	Effects   map[string]*Effect
	Hooks map[string]map[string]*Hook

	// ParentId is the room id for the room
	ParentId   int
	Gold       int
	Experience int
	Level      int

	Stam Meter
	Mana Meter

	// Attributes
	Str   Meter
	Dex   Meter
	Con   Meter
	Int   Meter
	Pie   Meter
	Armor int

	// Dice
	NumDice   int
	SidesDice int
	PlusDice  int

	// Magic
	ChanceCast      int
	Spells          []string
	WaterResistance int
	AirResistance   int
	FireResistance  int
	EarthResistance int
	BreathWeapon string

	//Threat table attacker -> damage
	TotalThreatDamage int
	ThreatTable       map[string]int
	CurrentTarget     string

	NumWander  int
	TicksAlive int
	WimpyValue int

	MobTickerUnload chan bool
	MobTicker       *time.Ticker
	// An int to hold a stun time.
	MobStunned int
}

// Pop the mob data
func LoadMob(mobData map[string]interface{}) (*Mob, bool) {
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
		make(map[string]*Effect),
		map[string]map[string]*Hook{
			"act": make(map[string]*Hook),
			"combat": make(map[string]*Hook),
			"peek": make(map[string]*Hook),
			"gridmove": make(map[string]*Hook),
			"move": make(map[string]*Hook),
			"say": make(map[string]*Hook),
		},
		-1,
		int(mobData["gold"].(int64)),
		int(mobData["experience"].(int64)),
		int(mobData["level"].(int64)),
		Meter{int(mobData["hpmax"].(int64)), int(mobData["hpcur"].(int64))},
		Meter{int(mobData["mpmax"].(int64)), int(mobData["mpcur"].(int64))},
		Meter{40, int(mobData["strength"].(int64))},
		Meter{40, int(mobData["dexterity"].(int64))},
		Meter{40, int(mobData["constitution"].(int64))},
		Meter{40, int(mobData["intelligence"].(int64))},
		Meter{40, int(mobData["piety"].(int64))},
		int(mobData["armor"].(int64)),
		int(mobData["ndice"].(int64)),
		int(mobData["sdice"].(int64)),
		int(mobData["pdice"].(int64)),
		int(mobData["casting_probability"].(int64)),
		[]string{},
		int(mobData["water_resistance"].(int64)),
		int(mobData["air_resistance"].(int64)),
		int(mobData["fire_resistance"].(int64)),
		int(mobData["earth_resistance"].(int64)),
		mobData["breathes"].(string),
		0,
		nil,
		"",
		int(mobData["numwander"].(int64)),
		0,
		int(mobData["wimpyvalue"].(int64)),
		nil,
		nil,
		0,
	}

	for _, spellN := range strings.Split(mobData["spells"].(string), ",") {
		if spellN != "" {
			newMob.Spells = append(newMob.Spells, spellN)
		}
	}

	for _, drop := range mobData["drops"].([]interface{}) {
		if drop != nil {
			dropData := drop.(map[string]interface{})
			if dropData["chance"] != nil {
				newMob.ItemList[int(dropData["item_id"].(int64))] = int(dropData["chance"].(int64))
			}
		}
	}

	for k, v := range mobData["flags"].(map[string]interface{}) {
		if v == nil {
			newMob.Flags[k] = false
		} else {
			newMob.Flags[k] = int(v.(int64)) != 0
		}
	}
	return newMob, true
}

func (m *Mob) StartTicking() {
	m.CalculateInventory()
	m.ThreatTable = make(map[string]int)
	m.MobTickerUnload = make(chan bool)
	tickModifier := 0
	if fastMoving, ok := m.Flags["fast_moving"]; ok {
		if fastMoving {
			tickModifier = 2
		}
	}
	m.MobTicker = time.NewTicker(time.Duration(8-tickModifier) * time.Second)
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

func (m *Mob) GetSpellMultiplier() float32 {
	return 1
}

// The mob brain is this ticker
func (m *Mob) Tick() {
	if m.MobStunned > 0 {
		m.MobStunned -= 1
	} else {
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
						if utils.StringIn(rankedThreats[0], Rooms[m.ParentId].Chars.MobList(m.Flags["detect_invisible"], true)) {
							m.CurrentTarget = rankedThreats[0]
							Rooms[m.ParentId].MessageAll(m.Name + " turns to " + m.CurrentTarget + "\n" + text.Reset)
						}
					}
				}
			}


			if m.CurrentTarget == "" && m.Placement != 3 {
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
				return
			}

			if m.CurrentTarget != "" && m.ChanceCast > 0 &&
				(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement)) >= 1) {
				// Try to cast a spell first
				target := Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false)
				spellSelected := false
				selectSpell := ""
				if utils.Roll(100, 1, 0) <= m.ChanceCast {
					for _ = range m.Spells {
						rand.Seed(time.Now().Unix())
						selectSpell = m.Spells[rand.Intn(len(m.Spells))]
						if selectSpell != "" {
							if utils.StringIn(selectSpell, OffensiveSpells) {
								spellSelected = true
							}
						}
					}

					if spellSelected {
						spellInstance, ok := Spells[selectSpell]
						if !ok {
							spellSelected = false
						}
						Rooms[m.ParentId].MessageAll(m.Name + " chants: " + spellInstance.Chant + "\n")
						Rooms[m.ParentId].MessageAll(m.Name + " cast a " + spellInstance.Name + " spell on " + target.Name + "\n")
						MobCast(m, target, spellInstance.Effect, map[string]interface{}{"magnitude": spellInstance.Magnitude})
						if target.Vit.Current == 0 {
							target.Died()
						}
						return
					}
				}
			}

			if m.CurrentTarget != "" && m.Flags["ranged_attack"] &&
				(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement)) >= 1) {
				target := Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false)
				// If we made it here, default out and do a range hit.
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
				return
			}

			if (m.CurrentTarget != "" &&
				m.Placement != Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement) ||
				(m.CurrentTarget != "" &&
					(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement)) > 1)) {
				oldPlacement := m.Placement
				if m.Placement > Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement {
					m.Placement--
				} else {
					m.Placement++
				}
				if !m.Flags["hidden"] {
					whichNumber := Rooms[m.ParentId].Mobs.GetNumber(m)
					Rooms[m.ParentId].MessageMovement(oldPlacement, m.Placement, m.Name+" #"+strconv.Itoa(whichNumber))
				}
				// Next to attack
			} else if m.CurrentTarget != "" &&
				m.Placement == Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false).Placement {
				// Am I against a fighter and they succeed in a parry roll?
				target := Rooms[m.ParentId].Chars.Search(m.CurrentTarget, false)
				if target.Class == 0 && target.Equipment.Main != nil && config.RollParry(target.Skills[target.Equipment.Main.ItemType].Value) {
					if target.Tier >= 10 {
						// It's a riposte
						actualDamage, _ := m.ReceiveDamage(int(math.Ceil(float64(target.InflictDamage()))))
						target.Write([]byte(text.Green + "You parry and riposte the attack from " + m.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + "\n" + text.Reset))
						if m.Stam.Current <= 0 {
							Rooms[m.ParentId].MessageAll(text.Green + target.Name + " killed " + m.Name)
							stringExp := strconv.Itoa(m.Experience)
							for k := range m.ThreatTable {
								Rooms[m.ParentId].Chars.Search(k, true).Write([]byte(text.Cyan + "You earn " + stringExp + "exp for the defeat of the " + m.Name + "\n" + text.Reset))
								Rooms[m.ParentId].Chars.Search(k, true).Experience.Add(m.Experience)
							}
							Rooms[m.ParentId].MessageAll(m.Name + " dies.")
							target.Write([]byte(text.White + m.DropInventory()))
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
			}

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
				// Successful roll!  Add this item to the inventory!
				newItem := Item{}
				copier.Copy(&newItem, Items[k])
				m.Inventory.Add(&newItem)
			}
		}
	}
}

func (m *Mob) DropInventory() string {
	var drops []string
	if len(m.Inventory.Contents) > 0 {
		for _, item := range m.Inventory.Contents {
			if item != nil {
				if err := m.Inventory.Remove(item); err == nil {
					if len(Rooms[m.ParentId].Items.Contents) < 15 {
						item.Placement = m.Placement
						Rooms[m.ParentId].Items.Add(item)
						drops = append(drops, item.Name)
					}else{
						Rooms[m.ParentId].MessageAll(item.Name + " falls on top of other items and rolls away.")
					}
				}
			}
		}
	}
	if m.Gold > 0 {
		newGold := Item{}
		copier.Copy(&newGold, Items[3456])
		newGold.Name = strconv.Itoa(m.Gold) + " gold pieces"
		newGold.Value = m.Gold
		Rooms[m.ParentId].Items.Add(&newGold)
		drops = append(drops, newGold.Name)
	}
	if len(drops) == 0 {
		return "The " + m.Name + " was carrying:\n Nothing"
	} else {
		return "The " + m.Name + " was carrying:\n " + strings.Join(drops, ", ")
	}
}

// On death calculate and distribute experience
func (m *Mob) CalculateExperience(attackerName string) {
	return
}

func (m *Mob) AddThreatDamage(damage int, attackerName string) {
	m.ThreatTable[attackerName] += damage
	if m.CurrentTarget == "" {
		m.CurrentTarget = attackerName
	}
}

func (m *Mob) ApplyEffect(effectName string, length string, interval string, effect func(), effectOff func()) {
	m.Effects[effectName] = NewEffect(length, interval, effect, effectOff)
	effect()
}

func (m *Mob) RemoveEffect(effectName string) {
	delete(m.Effects, effectName)
}


func (m *Mob) ApplyHook(hook string, hookName string, executions int, length string, interval int, effect func(), effectOff func()) {
	m.Hooks[hook][hookName] = NewHook(executions, length, interval, effect, effectOff)
}

func (m *Mob) RemoveHook(hook string, hookName string) {
	m.Hooks[hook][hookName].effectOff()
	valPresent := false
	for k, _ := range m.Hooks{
		valPresent = false
		for hName, _ := range m.Hooks[k] {
			if hName == hookName {
				valPresent = true
			}
		}
		if valPresent {
			delete(m.Hooks[k], hookName)
		}
	}
}

func (m *Mob) RunHook(hook string){
	for name, hookInstance := range m.Hooks[hook] {
		// Process Removing the hook
		if hookInstance.TimeRemaining() == 0 {
			m.RemoveHook(hook, name)
			continue
		}
		if hookInstance.interval > 0 {
			log.Println("Executing Hook", hook)
			log.Println(hookInstance.LastTriggerInterval())
			if hookInstance.LastTriggerInterval() <= 0 {
				hookInstance.RunHook()
			}
		}else if hookInstance.interval == -1 {
			log.Println("Executing Hook", hook)
			hookInstance.RunHook()
		}
	}
	return
}

func (m *Mob) GetInt() int {
	return m.Int.Current
}

func (m *Mob) ToggleFlag(flagName string) bool {
	if val, exists := m.Flags[flagName]; exists {
		m.Flags[flagName] = !val
		return true
	} else {
		return false
	}
}

func (m *Mob) ToggleFlagAndMsg(flagName string, msg string) {
	if val, exists := m.Flags[flagName]; exists {
		m.Flags[flagName] = !val
	} else {
		m.Flags[flagName] = true
	}
	log.Println(m.Name, " informed: ", msg)
}

func (m *Mob) ReceiveDamage(damage int) (int, int) {
	finalDamage := math.Ceil(float64(damage) * (1 - (float64(m.Armor/config.MobArmorReductionPoints) * config.MobArmorReduction)))
	m.Stam.Subtract(int(finalDamage))
	return int(finalDamage), 0
}

func (m *Mob) ReceiveVitalDamage(damage int) int {
	damageMod, _ := m.ReceiveDamage(damage)
	return damageMod
}

func (m *Mob) Heal(damage int) (int, int) {
	m.Stam.Add(damage)
	return damage, 0
}

func (m *Mob) HealStam(damage int) {
	m.Stam.Add(damage)
}

func (m *Mob) HealVital(damage int) {
	m.Heal(damage)
}

func (m *Mob) RestoreMana(damage int) {
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
	mobData["description"] = m.Description
	mobData["experience"] = m.Experience
	mobData["level"] = m.Level
	mobData["gold"] = m.Gold
	mobData["constitution"] = m.Con.Current
	mobData["strength"] = m.Str.Current
	mobData["intelligence"] = m.Int.Current
	mobData["dexterity"] = m.Dex.Current
	mobData["piety"] = m.Pie.Current
	mobData["mpmax"] = m.Mana.Max
	mobData["mpcur"] = m.Mana.Current
	mobData["hpcur"] = m.Stam.Current
	mobData["hpmax"] = m.Stam.Max
	mobData["sdice"] = m.SidesDice
	mobData["ndice"] = m.NumDice
	mobData["pdice"] = m.PlusDice
	mobData["spells"] = strings.Join(m.Spells, ",")
	mobData["casting_probability"] = m.ChanceCast
	mobData["armor"] = m.Armor
	mobData["numwander"] = m.NumWander
	mobData["wimpyvalue"] = m.WimpyValue
	mobData["air_resistance"] = m.AirResistance
	mobData["fire_resistance"] = m.FireResistance
	mobData["earth_resistance"] = m.EarthResistance
	mobData["water_resistance"] = m.WaterResistance
	mobData["hide_encounter"] = utils.Btoi(m.Flags["hide_encounter"])
	mobData["invisible"] = utils.Btoi(m.Flags["invisible"])
	mobData["permanent"] = utils.Btoi(m.Flags["permanent"])
	mobData["hostile"] = utils.Btoi(m.Flags["hostile"])
	mobData["undead"] = utils.Btoi(m.Flags["undead"])
	mobData["breathes"] = m.BreathWeapon
	mobData["fast_moving"] = utils.Btoi(m.Flags["fast_moving"])
	mobData["guard_treasure"] = utils.Btoi(m.Flags["guard_treasure"])
	mobData["take_treasure"] = utils.Btoi(m.Flags["take_treasure"])
	mobData["steals"] = utils.Btoi(m.Flags["steals"])
	mobData["block_exit"] = utils.Btoi(m.Flags["block_exit"])
	mobData["follows"] = utils.Btoi(m.Flags["block_exit"])
	mobData["no_steal"] = utils.Btoi(m.Flags["no_steal"])
	mobData["detect_invisible"] = utils.Btoi(m.Flags["detect_invisible"])
	mobData["no_stun"] = utils.Btoi(m.Flags["no_stun"])
	mobData["diseases"] = utils.Btoi(m.Flags["diseases"])
	mobData["poisons"] = utils.Btoi(m.Flags["poisons"])
	mobData["spits_acid"] = utils.Btoi(m.Flags["spits_acid"])
	mobData["ranged_attack"] = utils.Btoi(m.Flags["ranged_attack"])
	mobData["flees"] = utils.Btoi(m.Flags["flees"])
	mobData["blinds"] = utils.Btoi(m.Flags["blinds"])
	data.UpdateMob(mobData)
}
