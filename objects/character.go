package objects

import "C"
import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/prompt"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type Character struct {
	Object
	io.Writer
	PromptStyle
	CharId int
	// Our stuff!
	Equipment  *Equipment
	Inventory  *ItemInventory
	Permission permissions.Permissions

	// Invisible, Hidden, Resists, OOC, AFK
	Flags         map[string]bool
	FlagProviders map[string][]string
	Effects       map[string]*Effect
	HiddenEffects map[string]*Effect
	Modifiers     map[string]int

	// ParentId is the room id for the room
	ParentId int

	// Titles for all to see
	ClassTitle string
	Title      string

	// Gold
	BankGold Accumulator
	Gold     Accumulator

	// Exp
	Experience  Accumulator
	BonusPoints Accumulator
	Passages    Accumulator
	Broadcasts  int
	Evals       int
	//Char Stats
	Stam Meter
	Vit  Meter
	Mana Meter

	// Attributes
	Str Meter
	Dex Meter
	Con Meter
	Int Meter
	Pie Meter

	Tier     int
	Class    int
	Race     int
	Gender   string
	Birthday int

	// Cool Downs
	Timers map[string]time.Time

	// Extra
	MinutesPlayed int

	//TODO: Class Properties Heals/Enchants
	ClassProps map[string]interface{}

	Spells []string
	Skills map[int]*Accumulator

	CharTicker       *time.Ticker
	CharTickerUnload chan bool
	Hooks map[string]map[string]*Hook
	LastAction time.Time
}

func LoadCharacter(charName string, writer io.Writer) (*Character, bool) {
	charData, err := data.LoadChar(charName)
	if err {
		return nil, true
	} else {
		FilledCharacter := &Character{
			Object{
				Name:        strings.Title(charData["name"].(string)),
				Description: charData["description"].(string),
				Placement:   3,
				Commands: make(map[string]prompt.MenuItem),
			},
			writer,
			StyleNone,
			int(charData["character_id"].(int64)),
			RestoreEquipment(charData["equipment"].(string)),
			RestoreInventory(charData["inventory"].(string)),
			0,
			make(map[string]bool),
			make(map[string][]string),
			make(map[string]*Effect),
			make(map[string]*Effect),
			make(map[string]int),
			int(charData["parentid"].(int64)),
			config.ClassTitle(
				int(charData["class"].(int64)),
				charData["gender"].(string),
				int(charData["tier"].(int64))),
			charData["title"].(string),
			Accumulator{int(charData["bankgold"].(int64))},
			Accumulator{int(charData["gold"].(int64))},
			Accumulator{int(charData["experience"].(int64))},
			Accumulator{int(charData["bonuspoints"].(int64))},
			Accumulator{int(charData["passages"].(int64))},
			int(charData["broadcasts"].(int64)),
			int(charData["evals"].(int64)),
			Meter{int(charData["stammax"].(int64)), int(charData["stamcur"].(int64))},
			Meter{int(charData["vitmax"].(int64)), int(charData["vitcur"].(int64))},
			Meter{int(charData["manamax"].(int64)), int(charData["manacur"].(int64))},
			Meter{config.RaceDefs[config.AvailableRaces[int(charData["race"].(int64))]].StrMax, int(charData["strcur"].(int64))},
			Meter{config.RaceDefs[config.AvailableRaces[int(charData["race"].(int64))]].DexMax, int(charData["dexcur"].(int64))},
			Meter{config.RaceDefs[config.AvailableRaces[int(charData["race"].(int64))]].ConMax, int(charData["concur"].(int64))},
			Meter{config.RaceDefs[config.AvailableRaces[int(charData["race"].(int64))]].IntMax, int(charData["intcur"].(int64))},
			Meter{config.RaceDefs[config.AvailableRaces[int(charData["race"].(int64))]].PieMax, int(charData["piecur"].(int64))},
			int(charData["tier"].(int64)),
			int(charData["class"].(int64)),
			int(charData["race"].(int64)),
			charData["gender"].(string),
			int(charData["birthday"].(int64)),
			map[string]time.Time{"global": time.Now(), "use": time.Now(), "combat": time.Now()},
			int(charData["played"].(int64)),
			make(map[string]interface{}),
			[]string{},
			map[int]*Accumulator{0: {int(charData["sharpexp"].(int64))},
				1: {int(charData["thrustexp"].(int64))},
				2: {int(charData["bluntexp"].(int64))},
				3: {int(charData["poleexp"].(int64))},
				4: {int(charData["missileexp"].(int64))}},
			nil,
			make(chan bool),
			map[string]map[string]*Hook{
				"act": make(map[string]*Hook),
				"combat": make(map[string]*Hook),
				"peek": make(map[string]*Hook),
				"gridmove": make(map[string]*Hook),
				"move": make(map[string]*Hook),
				"say": make(map[string]*Hook),
			},
			time.Now(),
		}

		for _, spellN := range strings.Split(charData["spells"].(string), ",") {
			if spellN != "" {
				FilledCharacter.Spells = append(FilledCharacter.Spells, spellN)
			}
		}

		for k, v := range charData["flags"].(map[string]interface{}) {
			if v == nil {
				FilledCharacter.Flags[k] = false
			} else {
				FilledCharacter.Flags[k] = int(v.(int64)) != 0
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

func (c *Character) SetTimer(timer string, seconds int) {
	if timer == "combat" {
		if hasted, ok := c.Flags["haste"]; ok {
			if hasted {
				c.Timers[timer] = time.Now().Add(time.Duration(seconds - config.CalcHaste(c.Tier)) * time.Second)
				return
			}
		}
	}
	c.Timers[timer] = time.Now().Add(time.Duration(seconds) * time.Second)
}

func (c *Character) TimerReady(timer string) (bool, string) {
	// Always check Global:
	remaining := int(c.Timers["global"].Sub(time.Now()) / time.Second)
	if remaining <= 0 {
		if curTimer, ok := c.Timers[timer]; ok {
			remaining = int(curTimer.Sub(time.Now()) / time.Second)
			if remaining <= 0 {
				return true, ""
			}
		} else {
			return true, ""
		}

	}
	return false, "You have " + strconv.Itoa(remaining) + " seconds before you can perform this action."

}

func (c *Character) Unload() {
	c.CharTicker.Stop()
	c.CharTickerUnload <- true
}

func (c *Character) ToggleFlag(flagName string, provider string) {
	if _, exists := c.Flags[flagName]; exists {
		if c.Flags[flagName] == true && len(c.FlagProviders[flagName]) > 1{
			c.FlagProviders[flagName][utils.IndexOf(provider, c.FlagProviders[flagName])] = c.FlagProviders[flagName][len(c.FlagProviders[flagName])-1] // Copy last element to index i.
			c.FlagProviders[flagName][len(c.FlagProviders[flagName])-1] = ""   // Erase last element (write zero value).
			c.FlagProviders[flagName] = c.FlagProviders[flagName][:len(c.FlagProviders[flagName])-1]   // Truncate slice.
		}else if c.Flags[flagName] == true && len(c.FlagProviders[flagName]) == 1 {
			c.Flags[flagName] = false
			c.FlagProviders[flagName] = []string{}
		}else if c.Flags[flagName] == false && provider == ""{
			c.Flags[flagName] = true
		}else if c.Flags[flagName] == true && provider == "" {
			c.Flags[flagName] = false
			c.FlagProviders[flagName] = []string{provider}
		}else if c.Flags[flagName] == false && provider != "" {
			c.Flags[flagName] = true
			c.FlagProviders[flagName] = []string{provider}
		}
	}else{
		c.Flags[flagName] = true
		c.FlagProviders[flagName] = []string{provider}
	}
}

func (c *Character) ToggleFlagAndMsg(flagName string, provider string, msg string) {
	c.ToggleFlag(flagName, provider)
	c.Write([]byte(msg))
}

// SerialSaveEffects serializes all current user effects, removes them, and saves them to the database
func (c *Character) SerialSaveEffects(){

}

func (c *Character) SerialRestoreEffects(effectsBlob string){

}

func (c *Character) GetModifier(modifier string) int{
	if mod, ok := c.Modifiers[modifier]; ok {
		return mod
	}else{
		return 0
	}
}

func (c *Character) SetModifier(modifier string, value int) {
	if _, ok := c.Modifiers[modifier]; ok {
		c.Modifiers[modifier] += value
	}else{
		c.Modifiers[modifier] = value
	}
}

func (c *Character) GetStat(stat string) int {
	switch stat {
	case "int":
		return c.Int.Current + c.Modifiers["int"]
	case "str":
		return c.Str.Current + c.Modifiers["str"]
	case "dex":
		return c.Dex.Current + c.Modifiers["dex"]
	case "pie":
		return c.Pie.Current + c.Modifiers["pie"]
	case "con":
		return c.Con.Current + c.Modifiers["con"]
	case "armor":
		return c.Equipment.Armor + c.Modifiers["armor"]
	default:
		return 0
	}
}

func (c *Character) Save() {
	charData := make(map[string]interface{})
	charData["title"] = c.Title
	charData["name"] = c.Name
	charData["tier"] = c.Tier
	charData["character_id"] = c.CharId
	charData["experience"] = c.Experience.Value
	charData["spells"] = strings.Join(c.Spells, ",")
	charData["thrustexp"] = c.Skills[1].Value
	charData["bluntexp"] = c.Skills[2].Value
	charData["missileexp"] = c.Skills[4].Value
	charData["poleexp"] = c.Skills[3].Value
	charData["sharpexp"] = c.Skills[0].Value
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
	charData["equipment"] = c.Equipment.Jsonify()
	charData["inventory"] = c.Inventory.Jsonify()

	berz, ok := c.Flags["berserk"]
	if ok {
		if berz {
			charData["str"] = c.Str.Current - 5
		}
	}
	data.SaveChar(charData)

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
			strconv.Itoa(c.Stam.Current) + "|" +
			strconv.Itoa(c.Vit.Current) + "|" +
			strconv.Itoa(c.Mana.Current) +
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

type PromptStyle int

const (
	StyleNone = iota
	StyleStat
)

func (c *Character) Tick() {
	// TODO: Fix Tick, The tick is affected by all things around the character and any currently applied effects
	if Rooms[c.ParentId].Flags["heal_fast"] {
		c.Stam.Add(c.Con.Current * int(math.Round(config.ConHealRegenMod*2)))
		c.Vit.Add(c.Con.Current * int(math.Round(config.ConHealRegenMod*2)))
		c.Mana.Add(c.Pie.Current * int(math.Round(config.PieRegenMod*2)))
	} else {
		c.Stam.Add(c.Con.Current * int(math.Round(config.ConHealRegenMod)))
		c.Vit.Add(c.Con.Current * int(math.Round(config.ConHealRegenMod)))
		c.Mana.Add(c.Pie.Current * int(math.Round(config.PieRegenMod)))
	}
	// Loop the currently applied effects, drop them if needed, or execute their functions as necessary
	for name, effect := range c.Effects {
		// Process Removing the effect
		if effect.interval > 0 {
			if effect.LastTriggerInterval() <= 0 {
				effect.RunEffect()
			}
		}
		if effect.TimeRemaining() <= 0 {
			c.RemoveEffect(name)
			continue
		}
	}

}

func (c *Character) Died() {
	c.Write([]byte(text.Red + "#### OH GODS! YOU DIED!!!! Hahaha just kidding, this is beta. Here's a restore...\n" + text.Reset))
	c.Stam.Current = c.Stam.Max
	c.Vit.Current = c.Vit.Max
	c.Mana.Current = c.Mana.Max
	c.Write([]byte(text.Cyan + "## Your vitality, stamina, and mana were restored to max." + text.Reset + "\n"))
}

// Look Drop out the description of this character
func (c *Character) Look() (buildText string) {
	buildText = "You see " + c.Name + ", the young, " + config.TextGender[c.Gender] + ", " + config.AvailableRaces[c.Race] + " " + c.ClassTitle + "."
	return buildText
}

func (c *Character) ApplyEffect(effectName string, length string, interval string, effect func(), effectOff func()) {
	c.Effects[effectName] = NewEffect(length, interval, effect, effectOff)
	effect()
}

func (c *Character) RemoveEffect(effectName string) {
	if _, ok := c.Effects[effectName]; ok {
		c.Effects[effectName].effectOff()
		delete(c.Effects, effectName)
	}
}

func (c *Character) ApplyHook(hook string, hookName string, executions int, length string, interval int, effect func(), effectOff func()) {
	c.Hooks[hook][hookName] = NewHook(executions, length, interval, effect, effectOff)
}

func (c *Character) RemoveHook(hook string, hookName string) {
	c.Hooks[hook][hookName].effectOff()
	valPresent := false
	for k, _ := range c.Hooks{
		valPresent = false
		for hName, _ := range c.Hooks[k] {
			if hName == hookName {
				valPresent = true
			}
		}
		if valPresent {
			delete(c.Hooks[k], hookName)
		}
	}
}

func (c *Character) RunHook(hook string){
	for name, hookInstance := range c.Hooks[hook] {
		// Process Removing the hook
		if hookInstance.TimeRemaining() == 0 {
			c.RemoveHook(hook, name)
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

// ReceiveDamage Return stam and vital damage
func (c *Character) ReceiveDamage(damage int) (int, int) {
	stamDamage, vitalDamage := 0, 0
	resist := int(math.Ceil(float64(damage) * (float64(c.GetStat("armor")/config.ArmorReductionPoints) * config.ArmorReduction)))
	finalDamage := damage - resist
	if finalDamage > c.Stam.Current {
		stamDamage = c.Stam.Current
		vitalDamage = finalDamage - stamDamage
		c.Stam.Current = 0
		if vitalDamage > c.Vit.Current {
			vitalDamage = c.Vit.Current
			c.Vit.Current = 0
		} else {
			c.Vit.Subtract(vitalDamage)
		}
	} else {
		c.Stam.Subtract(finalDamage)
		stamDamage = finalDamage
		vitalDamage = 0
	}
	return stamDamage, vitalDamage
}

func (c *Character) ReceiveVitalDamage(damage int) int {
	finalDamage := int(math.Ceil(float64(damage) * (1 - (float64(c.GetStat("armor")/config.ArmorReductionPoints) * config.ArmorReduction))))
	if finalDamage > c.Vit.Current {
		finalDamage = c.Vit.Current
		c.Vit.Current = 0
	} else {
		c.Vit.Subtract(finalDamage)
	}
	return finalDamage
}

func (c *Character) ReceiveMagicDamage(damage int) (int, int) {
	//TODO Calculate some magic resistance damage
	return c.ReceiveDamage(damage)
}

func (c *Character) Heal(damage int) (int, int) {
	stamHeal, vitalHeal := 0, 0
	if damage > (c.Vit.Max - c.Vit.Current) {
		vitalHeal = c.Vit.Max - c.Vit.Current
		c.Vit.Current = c.Vit.Max
		stamHeal = damage - vitalHeal
		if stamHeal > (c.Stam.Max - c.Stam.Current) {
			stamHeal = c.Stam.Max - c.Stam.Current
			c.Stam.Current = c.Stam.Max
		} else {
			c.Stam.Add(stamHeal)
		}
	} else {
		c.Vit.Add(damage)
	}
	return stamHeal, vitalHeal
}

func (c *Character) HealVital(damage int) {
	c.Vit.Add(damage)
}

func (c *Character) HealStam(damage int) {
	c.Stam.Add(damage)
}

func (c *Character) RestoreMana(damage int) {
	c.Mana.Add(damage)
}

func (c *Character) GetSpellMultiplier() int {
	if c.Class == 4 {
		if c.Tier >= 15 {
			return 2
		} else if c.Tier >= 20 {
			return 4
		} else {
			return 1
		}
	}else{
		return 1
	}
}

func (c *Character) InflictDamage() (damage int) {
	//TODO: Monks need to not worry about weapons
	damage = utils.Roll(c.Equipment.Main.SidesDice,
		c.Equipment.Main.NumDice,
		c.Equipment.Main.PlusDice)

	damage += int(math.Ceil(float64(damage) * (config.StrDamageMod * float64(c.GetStat("str")))))
	// Add any modified base damage
	baseDamage, ok := c.Modifiers["base_damage"]
	if !ok {
		baseDamage = 0
	}
	damage += baseDamage
	if damage < 0 {
		damage = 0
	}
	return damage
}

func (c *Character) CastSpell(spell string) bool {
	return false
}

func (c *Character) MaxWeight() int {
	return config.MaxWeight(c.Str.Current)
}

func (c *Character) WriteMovement(previous int, new int, subject string) {
	mvAmnt := math.Abs(float64(previous - new))
	color := text.Yellow
	// Moving backwards
	if (previous > new) && (mvAmnt == 1) && (new > c.Placement) {
		c.Write([]byte(color + subject + " moves backwards, towards you." + text.Reset + "\n"))
	} else if (previous > new) && (mvAmnt == 1) && (new < c.Placement) {
		c.Write([]byte(color + subject + " moves backwards, away from you." + text.Reset + "\n"))
	} else if (previous > new) && (mvAmnt == 1) && (new == c.Placement) {
		c.Write([]byte(color + subject + " moves backwards, next to you." + text.Reset + "\n"))
	} else if (previous > new) && (mvAmnt == 2) && (new > c.Placement) {
		c.Write([]byte(color + subject + " sprints backwards, towards you." + text.Reset + "\n"))
	} else if (previous > new) && (mvAmnt == 2) && (new < c.Placement) {
		c.Write([]byte(color + subject + " sprints backwards, away from you." + text.Reset + "\n"))
	} else if (previous > new) && (mvAmnt == 2) && (new == c.Placement) {
		c.Write([]byte(color + subject + " sprints backwards, next to you." + text.Reset + "\n"))
		// Moving forwards
	} else if (previous < new) && (mvAmnt == 1) && (new < c.Placement) {
		c.Write([]byte(color + subject + " moves forwards, towards you." + text.Reset + "\n"))
	} else if (previous < new) && (mvAmnt == 1) && (new > c.Placement) {
		c.Write([]byte(color + subject + " moves forwards, away from you." + text.Reset + "\n"))
	} else if (previous < new) && (mvAmnt == 1) && (new == c.Placement) {
		c.Write([]byte(color + subject + " moves forwards, next to you." + text.Reset + "\n"))
	} else if (previous < new) && (mvAmnt == 2) && (new < c.Placement) {
		c.Write([]byte(color + subject + " sprints forwards, towards you." + text.Reset + "\n"))
	} else if (previous < new) && (mvAmnt == 2) && (new > c.Placement) {
		c.Write([]byte(color + subject + " sprints forwards, away from you." + text.Reset + "\n"))
	} else if (previous < new) && (mvAmnt == 2) && (new == c.Placement) {
		c.Write([]byte(color + subject + " sprints forwards, next to you." + text.Reset + "\n"))
	}
}
