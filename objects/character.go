package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/prompt"
	"github.com/ArcCS/Nevermore/text"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
)

type Character struct {
	Object
	io.Writer
	PromptStyle
	Menu map[string]prompt.MenuItem
	CharId int64
	// Our stuff!
	Equipment Equipment
	Inventory ItemInventory
	Permission permissions.Permissions


	Flags map[string]bool
	Effects map[string]Effect
	HiddenEffects map[string]Effect
	//TODO ??? Modifiers map[string]int64

	// Should we count idle time based on last command entry and register the character as absent

	// ParentId is the room id for the room
	ParentId int64

	// Titles for all to see
	ClassTitle string
	Title string

	// Gold
	BankGold Accumulator
	Gold Accumulator

	// Exp
	Experience Accumulator
	BonusPoints Accumulator
	Passages Accumulator
	AttrMoves Accumulator
	Broadcasts int64
	Evals int64
	//Char Stats
	Stam Meter
	Vit Meter
	Mana Meter

	// Attributes
	Str Meter
	Dex Meter
	Con Meter
	Int Meter
	Pie Meter

	Tier int64
	Class int64
	Race int64
	Gender string
	Birthday int64

	// Cool Downs
	Global Cooldown
	Use Cooldown
	Action Cooldown

	// Extra
	MinutesPlayed int64

	//TODO: Class Properties Heals/Enchants
	ClassProps map[string]interface{}

	Spells []string

	SharpExperience Accumulator
	ThrustExperience Accumulator
	BluntExperience Accumulator
	PoleExperience Accumulator
	MissileExperience Accumulator

	CharTicker *time.Ticker
	CharTickerUnload chan bool
}

func LoadCharacter(charName string, writer io.Writer) (*Character, bool){
	charData, err := data.LoadChar(charName)
	if err {
		return nil, true
	}else{
		FilledCharacter := &Character{
			Object{
				Name:        charData["name"].(string),
				Description: charData["description"].(string),
				Placement:   3,
			},
			writer,
			StyleNone,
			make(map[string]prompt.MenuItem),
			charData["character_id"].(int64),
			Equipment{},
			ItemInventory{},
			0,
			make(map[string]bool),
			make(map[string]Effect),
			make(map[string]Effect),
			charData["parentid"].(int64),
			config.ClassTitle(
				charData["class"].(int64),
				charData["gender"].(string),
				charData["tier"].(int64)),
			charData["title"].(string),
			Accumulator{charData["bankgold"].(int64)},
			Accumulator{charData["gold"].(int64)},
			Accumulator{charData["experience"].(int64)},
			Accumulator{charData["bonuspoints"].(int64)},
			Accumulator{charData["passages"].(int64)},
			Accumulator{charData["attrmoves"].(int64)},
			charData["broadcasts"].(int64),
			charData["evals"].(int64),
			Meter{charData["stammax"].(int64), charData["stamcur"].(int64)},
			Meter{charData["vitmax"].(int64), charData["vitcur"].(int64)},
			Meter{charData["manamax"].(int64), charData["manacur"].(int64)},
			Meter{int64(config.RaceDefs[config.AvailableRaces[charData["race"].(int64)]].StrMax), charData["strcur"].(int64)},
			Meter{int64(config.RaceDefs[config.AvailableRaces[charData["race"].(int64)]].DexMax), charData["dexcur"].(int64)},
			Meter{int64(config.RaceDefs[config.AvailableRaces[charData["race"].(int64)]].ConMax), charData["concur"].(int64)},
			Meter{int64(config.RaceDefs[config.AvailableRaces[charData["race"].(int64)]].IntMax), charData["intcur"].(int64)},
			Meter{int64(config.RaceDefs[config.AvailableRaces[charData["race"].(int64)]].PieMax), charData["piecur"].(int64)},
			charData["tier"].(int64),
			charData["class"].(int64),
			charData["race"].(int64),
			charData["gender"].(string),
			charData["birthday"].(int64),
			Cooldown{},
			Cooldown{},
			Cooldown{},
			charData["played"].(int64),
			make(map[string]interface{}),
			strings.Split(charData["spells"].(string), ","),
			Accumulator{charData["sharpexp"].(int64)},
			Accumulator{charData["thrustexp"].(int64)},
			Accumulator{charData["bluntexp"].(int64)},
			Accumulator{charData["sharpexp"].(int64)},
			Accumulator{charData["missileexp"].(int64)},
			nil,
			make(chan bool),
		}

		for k, v := range charData["flags"].(map[string]interface{}){
			if v == nil{
				FilledCharacter.Flags[k] = false
			}else {
				FilledCharacter.Flags[k] = v.(int64) != 0
			}
		}

		// GM Specifics:
		if FilledCharacter.Class >= 99 {
			FilledCharacter.Flags["hidden"] = true
			FilledCharacter.Flags["invisible"] = true
		}

		FilledCharacter.CharTicker = time.NewTicker(8 * time.Second)
		go func() {
			for {
				select {
				case <-FilledCharacter.CharTickerUnload:
					return
				case <-FilledCharacter.CharTicker.C:
					FilledCharacter.Tick()
				}
			}
		}()

		return FilledCharacter, false
	}
}


// TODO:  A hooking system
// Extend the anon scripts to bind from items and add hooks to characters
// Rooms should take the hook system as well and invoke onActions.
/*static_str
static_text
num_ranges
num_vals
hi_numeric // Hidden numerics
hi_string  // Hidden string
hi_text  // Hidden text
hook [list]
 onaction
 onmove
 onattack
 onget
 onreset
 oncleanup
veto*/

func (c *Character) Unload(){
	c.CharTicker.Stop()
	c.CharTickerUnload<-true
}

func (c *Character) OnAction(act string){
	//TODO: Loop the actions based on the act sent
	// Invoke functions tied to the act
	return
}

func (c *Character) ToggleFlag(flagName string) bool {
	if val, exists := c.Flags[flagName]; exists{
		c.Flags[flagName] = !val
		return true
	}else{
		return false
	}
}

func (c *Character) Save(){
	charData := make(map[string]interface{})
	charData["title"] = c.Title
	charData["name"] = c.Name
	charData["tier"] = c.Tier
	charData["character_id"] = c.CharId
	charData["experience"] = c.Experience.Value
	charData["spells"] = strings.Join(c.Spells, ",")
	charData["thrustexp"] = c.ThrustExperience.Value
	charData["bluntexp"] = c.BluntExperience.Value
	charData["missileexp"] = c.MissileExperience.Value
	charData["poleexp"] = c.PoleExperience.Value
	charData["sharpexp"] = c.SharpExperience.Value
	charData["bankgold"] = c.BankGold.Value
	charData["gold"] = c.Gold.Value
	charData["evals"] = c.Evals
	charData["broadcasts"] = c.Broadcasts
	charData["played"] = c.MinutesPlayed
	charData["description"] = c.Description
	charData["parent_id"] = c.ParentId
	charData["str"] = c.Str.Current
	charData["con"] = c.Con.Current
	charData["dex"] = c.Dex.Current
	charData["pie"] = c.Pie.Current
	charData["intel"] = c.Int.Current
	charData["manacur"] = c.Mana.Current
	charData["vitcurr"] = c.Vit.Current
	charData["stamcurr"] = c.Stam.Current
	charData["manamax"] = c.Mana.Max
	charData["vitmax"] = c.Vit.Max
	charData["stammax"] = c.Stam.Max
	//TODO Flags?
	data.SaveChar(charData)

	//TODO Process Equipment
	//TODO Process Effects
}

func (c *Character) SetPromptStyle(new PromptStyle) (old PromptStyle) {
	old, c.PromptStyle = c.PromptStyle, new
	return
}

// buildPrompt creates a prompt appropriate for the current PromptStyle. This
// is mostly useful for dynamic prompts that show Character statistics.
func (c *Character) buildPrompt() []byte {
	switch c.PromptStyle {
	case StyleNone:
		return []byte(text.Prompt + " > ")
	case StyleStat:
		return []byte(text.Prompt +
				strconv.Itoa(int(c.Stam.Current)) + "|" +
				strconv.Itoa(int(c.Vit.Current)) + "|" +
				strconv.Itoa(int(c.Mana.Current)) +
			" > ")
	default:
		return []byte{}
	}
}

// Write writes the specified byte slice to the associated client.
func (c *Character) Write(b []byte) (n int, err error) {
	if c == nil {
		return
	}

	b = append(b, c.buildPrompt()...)
	if c != nil {
		n, err = c.Writer.Write(b)
	}
	return
}

// Free makes sure references are nil'ed when the Character attribute is freed.
func (c *Character) Free() {
	c.EmptyMenu()
	//TODO: Clear Effect Timers

	if c != nil {
		c.Writer = nil
	}
}

type PromptStyle int

const (
	StyleNone = iota
	StyleStat
)

func (c *Character) Tick(){
	//_,_ = c.Write([]byte(text.Good + "Your Ticker Executed Here!"))
	// The tick is affected by all things around the character and any currently applied effects
	if Rooms[c.ParentId].Flags["heal_fast"] {
		c.Stam.Add(c.Con.Current * 2)
		c.Vit.Add(c.Con.Current * 2)
		c.Mana.Add(c.Pie.Current * 2)
	} else {
		c.Stam.Add(c.Con.Current)
		c.Mana.Add(c.Pie.Current)
	}

	// Loop the currently applied effects, drop them if needed, or execute their functions as necessary
	for name, effect := range c.Effects {
		// Process Removing the effect
		if effect.TimeRemaining() <= 0 {
			//TODO: Execute the spell to turn this off, but for now, just toggles
			if effect.effectOff == "toggle"{
				c.ToggleFlag(name)
			}
			delete(c.Effects, name)
			continue
		}

		//TODO:  Process an interval execution of the effect
	}


}

func (c *Character) Died() {

}

// Drop out the description of this character
func (c *Character) Look() (buildText string) {
	buildText = "You see " + c.Name + ", the young, " + config.TextGender[c.Gender]  + ", " + config.AvailableRaces[c.Race] + " " + c.ClassTitle  +"."
	return buildText
}

func (c *Character) EmptyMenu () {
	for k := range c.Menu {
		delete(c.Menu, k)
	}
}

func (c *Character) AddMenu(menuItem string, menuCmd string) {
	c.Menu[menuItem] = prompt.MenuItem{
		Command: menuCmd,
	}
}

func (c *Character) Info() (buildText string ) {
	buildText = ""

	return
}

func (c *Character) ApplyEffect(){
	return
}

func (c *Character) RemoteEffect(effect string){
	return
}


func (c *Character) ReceiveDamage(damage int){
	return
}

func (c *Character) ReceiveVitalDamage(damage int){

}

func (c *Character) Heal(damage int){
	return
}

func (c *Character) HealVital(damage int){

}

func (c *Character) RestoreMana(damage int){

}

func (c *Character) InflictDamage() (damage int){
	return 0
}

func (c *Character) CastSpell(spell string) bool {
	return false
}

func (c *Character) MaxWeight() int64 {
	return config.MaxWeight(c.Str.Current)
}

func (c *Character) WriteMovement(previous int64, new int64, subject string) {
	mvAmnt := math.Abs(float64(previous - new))
	color := text.Yellow
	// Moving backwards
	if (previous > new) && (mvAmnt == 1) && (new > c.Placement) {
		c.Write([]byte(color + subject + " moves backwards, towards you." + text.Reset))
	}else if (previous > new) && (mvAmnt == 1) && (new == c.Placement) {
		c.Write([]byte(color + subject + " moves backwards, next to you." + text.Reset))
	}else if (previous > new) && (mvAmnt == 2) && (new > c.Placement) {
		c.Write([]byte(color + subject + " sprints backwards, towards you." + text.Reset))
	}else if (previous > new) && (mvAmnt == 2) && (new == c.Placement) {
		c.Write([]byte(color + subject + " sprints backwards, next to you." + text.Reset))
	// Moving forwards
	}else if (previous < new) && (mvAmnt == 1) && (new < c.Placement) {
		c.Write([]byte(color + subject + " moves forwards, towards you." + text.Reset))
	}else if (previous < new) && (mvAmnt == 1) && (new == c.Placement) {
		c.Write([]byte(color + subject + " moves forwards, next to you." + text.Reset))
	}else if (previous < new) && (mvAmnt == 2) && (new < c.Placement) {
		c.Write([]byte(color + subject + " sprints forwards, towards you." + text.Reset))
	}else if (previous < new) && (mvAmnt == 2) && (new == c.Placement) {
		c.Write([]byte(color + subject + " sprints forwards, next to you." + text.Reset))
	}

}
