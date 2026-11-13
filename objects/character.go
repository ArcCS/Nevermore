package objects

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/prompt"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
)

type Character struct {
	Object
	io.Writer
	PromptStyle
	CharId     int
	Equipment  *Equipment
	Inventory  *ItemInventory
	Permission permissions.Permissions

	// Invisible, Hidden, Resists, OOC, AFK
	Flags         map[string]bool
	FlagProviders map[string][]string
	Effects       map[string]*Effect
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

	Tier       int
	Class      int
	Race       int
	Gender     string
	Birthday   int
	Birthdate  int
	Birthmonth int

	// Cool Downs
	Timers map[string]time.Time

	// Extra
	MinutesPlayed int

	ClassProps map[string]int

	Spells            []string
	Skills            map[int]*Accumulator
	ElementalAffinity map[string]*Accumulator

	CharTicker       *time.Ticker
	CharTickerUnload chan bool
	MsgBuffer        chan []byte
	SongTicker       *time.Ticker
	SongTickerUnload chan bool
	Hooks            map[string]map[string]*Hook
	LastRefresh      time.Time
	LoginTime        time.Time
	//Party Stuff
	PartyFollow     string
	PartyFollowers  []string
	Victim          interface{}
	Resist          bool
	OOCSwap         int
	LastSave        time.Time
	Unloader        func()
	LastMessenger   string
	DeathInProgress bool
	Rerolls         int
	Disconnect      func()
	Pose            string
	claimed         int32
	// ResumeGraceUntil marks the window after a session resume in which a
	// death is downgraded to a penalty-free lag death. It ends early the
	// moment the player fights or moves - see SetTimer and the GO handler.
	ResumeGraceUntil time.Time
	// LastCombat is the last moment the character landed a hit, cast an
	// offensive spell, or was attacked. Regen is cut while it is within
	// CombatRegenWindowSeconds, and monk chi only bleeds away once it is
	// ChiDecayGraceSeconds in the past.
	LastCombat time.Time
	// FeintCharges counts the incoming attacks the feint dodge bonus still covers.
	FeintCharges int
	// TodWarnedTarget is the last mob the monk was told looked vulnerable to
	// a touch of death, so the warning fires once per target.
	TodWarnedTarget *Mob
	// RescuedBy names the paladin drawing this character's attackers until
	// RescueUntil. Never persisted.
	RescuedBy   string
	RescueUntil time.Time
	// CurrentSong is the song a bard is performing ("" when silent). Sing mix
	// changes it mid-song, and the song ticker reads it on every pulse.
	CurrentSong string
	// SongLength is the full length in seconds of the song being performed,
	// fixed when it starts so a mid-song tier gain can't skew the stop refund.
	SongLength int
}

func LoadCharacter(charName string, writer io.Writer, disconnect func()) (*Character, bool) {
	charData, err := data.LoadChar(charName)
	NewEquipment, ErroredEquip := RestoreEquipment(charData["equipment"].(string), int(charData["class"].(int64)))
	NewInventory := RestoreInventory(charData["inventory"].(string))
	for _, item := range ErroredEquip {
		NewInventory.Add(&item)
	}
	lastRefresh, _ := time.Parse(time.RFC3339, charData["lastrefresh"].(string))
	if err {
		return nil, true
	} else {
		FilledCharacter := &Character{
			Object{
				Name:        utils.Title(charData["name"].(string)),
				Description: charData["description"].(string),
				Placement:   3,
				Commands:    make(map[string]prompt.MenuItem),
			},
			writer,
			StyleNone,
			int(charData["character_id"].(int64)),
			NewEquipment,
			NewInventory,
			0,
			make(map[string]bool),
			make(map[string][]string),
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
			int(charData["birthdate"].(int64)),
			int(charData["birthmonth"].(int64)),
			map[string]time.Time{"global": time.Now(), "use": time.Now(), "combat": time.Now()},
			int(charData["played"].(int64)),
			make(map[string]int),
			[]string{},
			map[int]*Accumulator{0: {int(charData["sharpexp"].(int64))},
				1:  {int(charData["thrustexp"].(int64))},
				2:  {int(charData["bluntexp"].(int64))},
				3:  {int(charData["poleexp"].(int64))},
				4:  {int(charData["missileexp"].(int64))},
				5:  {int(charData["handexp"].(int64))},
				6:  {int(charData["fireexp"].(int64))},
				7:  {int(charData["airexp"].(int64))},
				8:  {int(charData["earthexp"].(int64))},
				9:  {int(charData["waterexp"].(int64))},
				10: {int(charData["divinity"].(int64))},
				11: {int(charData["stealthexp"].(int64))}},
			map[string]*Accumulator{
				"fire":  {0},
				"earth": {0},
				"water": {0},
				"air":   {0}},
			nil,
			make(chan bool),
			make(chan []byte, 100),
			nil,
			make(chan bool, 1),
			map[string]map[string]*Hook{
				"act":      make(map[string]*Hook),
				"combat":   make(map[string]*Hook),
				"peek":     make(map[string]*Hook),
				"gridmove": make(map[string]*Hook),
				"move":     make(map[string]*Hook),
				"say":      make(map[string]*Hook),
				"use":      make(map[string]*Hook),
				"attacked": make(map[string]*Hook),
				"ooc":      make(map[string]*Hook),
				"hide":     make(map[string]*Hook),
			},
			lastRefresh,
			time.Now(),
			"",
			[]string{},
			nil,
			true,
			int(charData["oocswap"].(int64)),
			time.Now(),
			nil,
			"",
			false,
			int(charData["rerolls"].(int64)),
			disconnect,
			"",
			0,
			time.Time{},
			time.Now(),
			0,
			nil,
			"",
			time.Time{},
			"",
			0,
		}

		// The pool formula may have changed since the character was last
		// saved, so rebuild max mana (chi for monks) rather than trusting the DB.
		FilledCharacter.Mana.Max = config.CalcMana(FilledCharacter.Tier, FilledCharacter.Int.Current, FilledCharacter.Pie.Current, FilledCharacter.Class)
		if FilledCharacter.Mana.Current > FilledCharacter.Mana.Max {
			FilledCharacter.Mana.Current = FilledCharacter.Mana.Max
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

		if FilledCharacter.Class == 4 || FilledCharacter.Class == 5 {
			FilledCharacter.ClassProps["enchants"] = int(charData["enchants"].(int64))
		}
		if FilledCharacter.Class == 5 || FilledCharacter.Class == 6 {
			FilledCharacter.ClassProps["heals"] = int(charData["heals"].(int64))
		}
		if FilledCharacter.Class == 7 || FilledCharacter.Class == 5 {
			// Clerics only gained restores on 2026-10-07; one with none stored starts with a full day's worth.
			if restores, ok := charData["restores"].(int64); ok {
				FilledCharacter.ClassProps["restores"] = int(restores)
			} else {
				FilledCharacter.ClassProps["restores"] = 5
			}
		}

		// GM Specifics:
		if FilledCharacter.Class >= 99 {
			FilledCharacter.Flags["hidden"] = true
			FilledCharacter.Flags["invisible"] = true
		}

		// Refresh or not to refresh on load?
		if RefreshDue(lastRefresh, time.Now()) {
			FilledCharacter.Refresh()
			FilledCharacter.LastRefresh = time.Now()
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

		FilledCharacter.SerialRestoreEffects(charData["effects"].(string))
		FilledCharacter.SerialRestoreTimers(charData["timers"].(string))

		FilledCharacter.Equipment.Owner = FilledCharacter.Name
		FilledCharacter.Equipment.FlagOn = FilledCharacter.FlagOn
		FilledCharacter.Equipment.FlagOff = FilledCharacter.FlagOff
		FilledCharacter.Equipment.CanEquip = FilledCharacter.CanEquip
		FilledCharacter.Equipment.ReturnToInventory = FilledCharacter.ReturnToInventory
		FilledCharacter.Equipment.CheckEquipment()
		FilledCharacter.Equipment.PostEquipmentLight()
		go FilledCharacter.EnableWrites()
		return FilledCharacter, false
	}
}

// GetCurrentWeight returns the current carrying weight of the character.
func (c *Character) GetCurrentWeight() int {
	return c.Inventory.GetTotalWeight() + c.Equipment.GetWeight() + c.Equipment.GetPreparedWeight()
}

// SuppressWrites Used when erroring a character out, or when we don't want to put anything on the write buffer
func (c *Character) SuppressWrites() {
	c.Flags["allow_writes"] = false
}

func (c *Character) EnableWrites() {
	c.Flags["allow_writes"] = true
}

// SetResumeGrace opens the post-resume window in which a death is downgraded
// to a penalty-free lag death. Only ResumeGame should call this.
func (c *Character) SetResumeGrace(seconds int) {
	c.ResumeGraceUntil = time.Now().Add(time.Duration(seconds) * time.Second)
}

// InResumeGrace reports whether the character is still inside that window.
func (c *Character) InResumeGrace() bool {
	return time.Now().Before(c.ResumeGraceUntil)
}

// ClearResumeGrace ends the window early: a character who has fought or moved
// since resuming is playing again and takes deaths at full price.
func (c *Character) ClearResumeGrace() {
	c.ResumeGraceUntil = time.Time{}
}

// MarkCombat notes that the character is fighting, which cuts regen and
// holds off chi decay.
func (c *Character) MarkCombat() {
	c.LastCombat = time.Now()
}

// InCombat reports whether the character fought recently enough for regen
// to be halved.
func (c *Character) InCombat() bool {
	return time.Since(c.LastCombat) < time.Duration(config.CombatRegenWindowSeconds)*time.Second
}

// GainChi awards chi to a monk and returns how much was actually added. Chi
// is not generated while flurrying unless force is set (leap strike), and
// meditating boosts every gain. No single award adds more than
// config.ChiGainCapPercent of the pool, applied after the meditate boost.
func (c *Character) GainChi(amount int, force bool) int {
	if c.Class != config.MONK || amount <= 0 {
		return 0
	}
	c.MarkCombat()
	if c.CheckFlag("flurry") && !force {
		return 0
	}
	if c.CheckFlag("meditate") {
		amount = int(math.Ceil(float64(amount) * config.MeditateChiMultiplier))
	}
	if cap := config.ChiGainCap(c.Mana.Max); amount > cap {
		amount = cap
	}
	before := c.Mana.Current
	c.Mana.Add(amount)
	return c.Mana.Current - before
}

// tickChi settles chi toward config.ChiDecayFloor once the monk has been out
// of combat long enough: a pool above it bleeds down to it, and one below it
// recovers up to it, at the same rate, so a rested monk always has enough
// chi to open a fight with a leap.
func (c *Character) tickChi() {
	floor := min(config.ChiDecayFloor(c.Tier), c.Mana.Max)
	if c.Mana.Current == floor || c.CheckFlag("meditate") {
		return
	}
	if time.Since(c.LastCombat) < time.Duration(config.ChiDecayGraceSeconds)*time.Second {
		return
	}
	amount := config.ChiDecayAmount(c.Mana.Max)
	if c.Mana.Current < floor {
		c.Mana.Add(min(amount, floor-c.Mana.Current))
		return
	}
	c.Mana.Subtract(min(amount, c.Mana.Current-floor))
}

// ConsumeFeintCharge spends one feint charge; the flag drops with the last
// one. Safe to call from the mob goroutine: it touches only Flags and an int.
func (c *Character) ConsumeFeintCharge() {
	c.FeintCharges--
	if c.FeintCharges <= 0 {
		c.FeintCharges = 0
		c.FlagOff("feint", "feint")
	}
}

// RecalcMaxes rebuilds the stamina, vitality, and mana pools from the
// character's tier, stats, and class.
func (c *Character) RecalcMaxes() {
	c.Stam.Max = config.CalcStamina(c.Tier, c.Con.Current, c.Class)
	c.Vit.Max = config.CalcHealth(c.Tier, c.Con.Current, c.Class)
	c.Mana.Max = config.CalcMana(c.Tier, c.Int.Current, c.Pie.Current, c.Class)
	if c.Stam.Current > c.Stam.Max {
		c.Stam.Current = c.Stam.Max
	}
	if c.Vit.Current > c.Vit.Max {
		c.Vit.Current = c.Vit.Max
	}
	if c.Mana.Current > c.Mana.Max {
		c.Mana.Current = c.Mana.Max
	}
}

func (c *Character) SetTimer(timer string, seconds int) {
	// Every combat command funnels through here with a "combat"-prefixed
	// timer name, so this is the one place that sees them all.
	if strings.HasPrefix(timer, "combat") {
		c.ClearResumeGrace()
	}
	// Looks like Sub on timers only accounts for seconds and not milliseconds, so folks are still blazing past a timer here add one second to increase to 8.
	seconds += 1
	if c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		return
	}
	if hasted, ok := c.Flags["haste"]; ok {
		if hasted {
			c.Timers[timer] = time.Now().Add(time.Duration(seconds-config.CalcHaste(c.Tier)) * time.Second)
			return
		}
	}
	c.Timers[timer] = time.Now().Add(time.Duration(seconds) * time.Second)
}

func (c *Character) ReturnToInventory(item *Item) {
	c.Inventory.Add(item)
}

// TimerReady reports whether the named timer and the global timer have both expired.
func (c *Character) TimerReady(timer string) (bool, string) {
	return c.timerReady(timer, true)
}

// TimerReadyIgnoreGlobal reports whether the named timer has expired, letting
// the action through even while the global timer (e.g. blocked movement)
// is still running. A stun (including paralysis) still blocks the action.
func (c *Character) TimerReadyIgnoreGlobal(timer string) (bool, string) {
	if ready, msg := c.timerReady("stun", false); !ready {
		return ready, msg
	}
	return c.timerReady(timer, false)
}

// TimerRemaining is how many seconds are left on the named timer, zero when
// it has expired or was never set.
func (c *Character) TimerRemaining(timer string) float64 {
	return math.Max(0, time.Until(c.Timers[timer]).Seconds())
}

// SetStanceTimer starts the cooldown between stance toggles. It is a fixed
// anti-spam delay, so unlike SetTimer haste does not shorten it.
func (c *Character) SetStanceTimer() {
	if c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		return
	}
	c.Timers["stance"] = time.Now().Add(time.Duration(config.StanceCooldown) * time.Second)
}

// Root stops the character moving under their own power for the given
// seconds: no squares, no exits, no teleporting away. Haste doesn't shorten
// it, and a shorter root never cuts a longer one short.
func (c *Character) Root(seconds int) {
	if c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		return
	}
	until := time.Now().Add(time.Duration(seconds) * time.Second)
	if until.After(c.Timers["rooted"]) {
		c.Timers["rooted"] = until
	}
}

// Rooted reports whether the character is held in place, and the line to
// show them if they try to move.
func (c *Character) Rooted() (bool, string) {
	remaining := time.Until(c.Timers["rooted"])
	if remaining <= 0 {
		return false, ""
	}
	return true, "You are held fast and can't move for " + strconv.Itoa(int(math.Ceil(remaining.Seconds()))) + " more seconds."
}

// StanceReady reports whether the stance cooldown has expired. A stance can
// be changed at any time, so nothing else is checked: not the global or
// combat timers, and not a stun.
func (c *Character) StanceReady() (bool, string) {
	remaining := time.Until(c.Timers["stance"])
	if remaining <= 0 {
		return true, ""
	}
	return false, text.Gray + "You have " + strconv.Itoa(int(math.Ceil(remaining.Seconds()))) + " seconds before you can perform this action. (stance)"
}

func (c *Character) timerReady(timer string, checkGlobal bool) (bool, string) {
	remaining := 0.0
	globalRemaining := 0.0
	if checkGlobal {
		globalRemaining = math.Ceil(float64(c.Timers["global"].Sub(time.Now()) / time.Second))
	}
	timeCheck := timer
	if curTimer, ok := c.Timers[timer]; ok {
		remaining = math.Ceil(float64(curTimer.Sub(time.Now()) / time.Second))
	}
	if remaining <= 0 && globalRemaining <= 0 {
		return true, ""
	} else if globalRemaining > remaining {
		timeCheck = "global"
		remaining = globalRemaining
	}

	return false, text.Gray + "You have " + strconv.Itoa(int(remaining)) + " seconds before you can perform this action. (" + timeCheck + ")"

}

// SingSong Bards get their own special ticker
// TODO(?): Maybe eventually this is a generalized aura?
func (c *Character) SingSong(song string, tickRate int) {
	c.CurrentSong = song
	c.FlagOn("singing", "sing")
	c.SongTicker = time.NewTicker(time.Duration(tickRate) * time.Second)
	go func() {
		for {
			select {
			case <-c.SongTickerUnload:
				c.SongTicker.Stop()
				c.FlagOffAndMsg("singing", "sing", "You stop singing.")
				return
			case <-c.SongTicker.C:
				room := Rooms[c.ParentId]
				room.LockRoom(c.Name+":SongTick", false)
				c.songPulse()
				room.UnlockRoom(c.Name+":SongTick", false)
			}
		}
	}()
}

// MixSong switches the song being performed without ending it, so the time
// left on it is unchanged, and lands the new song's first pulse at once. The
// caller holds the room lock.
func (c *Character) MixSong(song string) {
	c.CurrentSong = song
	if c.SongTicker != nil {
		// Restart the period so the next pulse is a full tick after this one.
		c.SongTicker.Reset(time.Duration(config.SongTickSeconds(c.Tier)) * time.Second)
	}
	c.songPulse()
}

// RefundTimer takes seconds off a running timer, never past now. A timer that
// has expired or was never set (staff have none) is left alone.
func (c *Character) RefundTimer(timer string, seconds float64) {
	until, ok := c.Timers[timer]
	now := time.Now()
	if !ok || seconds <= 0 || !until.After(now) {
		return
	}
	until = until.Add(-time.Duration(seconds * float64(time.Second)))
	if until.Before(now) {
		until = now
	}
	c.Timers[timer] = until
}

func (c *Character) Unload() {
	c.CharTicker.Stop()
	c.CharTickerUnload <- true
	close(c.MsgBuffer)
}

func (c *Character) IsClaimed() bool {
	return atomic.LoadInt32(&c.claimed) == 1
}

func (c *Character) ClearClaim() {
	atomic.StoreInt32(&c.claimed, 0)
}

// PrepareUnload saves and removes the character from the world. Returns false
// if the character was claimed by a reconnecting session and should not be
// destroyed.
func (c *Character) PrepareUnload() bool {
	c.Save()
	c.Unfollow()
	c.LoseParty()
	c.PurgeEffects()
	if !ActiveCharacters.Remove(c) {
		return false
	}
	Rooms[c.ParentId].Chars.Remove(c)
	return true
}

func (c *Character) ToggleFlag(flagName string, provider string) {

	if _, exists := c.Flags[flagName]; exists {
		if c.Flags[flagName] == true && utils.StringIn(provider, c.FlagProviders[flagName]) && len(c.FlagProviders[flagName]) > 1 {
			c.FlagProviders[flagName][utils.IndexOf(provider, c.FlagProviders[flagName])] = c.FlagProviders[flagName][len(c.FlagProviders[flagName])-1] // Copy last element to index i.
			c.FlagProviders[flagName][len(c.FlagProviders[flagName])-1] = ""                                                                            // Erase last element (write zero value).
			c.FlagProviders[flagName] = c.FlagProviders[flagName][:len(c.FlagProviders[flagName])-1]                                                    // Truncate slice.
		} else if c.Flags[flagName] == true && !utils.StringIn(provider, c.FlagProviders[flagName]) && len(c.FlagProviders[flagName]) >= 1 {
			c.FlagProviders[flagName] = append(c.FlagProviders[flagName], provider)
		} else if c.Flags[flagName] == true && len(c.FlagProviders[flagName]) == 1 {
			c.Flags[flagName] = false
			c.FlagProviders[flagName] = []string{}
		} else if c.Flags[flagName] == false && provider == "" {
			c.Flags[flagName] = true
		} else if c.Flags[flagName] == true && provider == "" {
			c.Flags[flagName] = false
			c.FlagProviders[flagName] = []string{}
		} else if c.Flags[flagName] == false && provider != "" {
			c.Flags[flagName] = true
			c.FlagProviders[flagName] = []string{provider}
		}
	} else {
		c.Flags[flagName] = true
		c.FlagProviders[flagName] = []string{provider}
	}
}

func (c *Character) FlagOn(flagName string, provider string) {

	if _, exists := c.Flags[flagName]; exists {
		if c.Flags[flagName] == true && !utils.StringIn(provider, c.FlagProviders[flagName]) && len(c.FlagProviders[flagName]) >= 1 {
			c.FlagProviders[flagName] = append(c.FlagProviders[flagName], provider)
		} else if c.Flags[flagName] == false {
			c.Flags[flagName] = true
			c.FlagProviders[flagName] = []string{provider}
		}
	} else {
		c.Flags[flagName] = true
		c.FlagProviders[flagName] = []string{provider}
	}
}

func (c *Character) FlagOff(flagName string, provider string) {
	if _, exists := c.Flags[flagName]; exists {
		if c.Flags[flagName] == true && utils.StringIn(provider, c.FlagProviders[flagName]) && len(c.FlagProviders[flagName]) > 1 {
			c.FlagProviders[flagName][utils.IndexOf(provider, c.FlagProviders[flagName])] = c.FlagProviders[flagName][len(c.FlagProviders[flagName])-1] // Copy last element to index i.
			c.FlagProviders[flagName][len(c.FlagProviders[flagName])-1] = ""                                                                            // Erase last element (write zero value).
			c.FlagProviders[flagName] = c.FlagProviders[flagName][:len(c.FlagProviders[flagName])-1]                                                    // Truncate slice.
		} else if c.Flags[flagName] == true && len(c.FlagProviders[flagName]) == 1 {
			c.Flags[flagName] = false
			c.FlagProviders[flagName] = []string{}
		}
	}
}

func (c *Character) FlagOnAndMsg(flagName string, provider string, msg string) {
	c.FlagOn(flagName, provider)
	if _, err := c.Write([]byte(msg)); err != nil {
		log.Println("Error writing to player: ", err)
	}
}

func (c *Character) FlagOffAndMsg(flagName string, provider string, msg string) {
	c.FlagOff(flagName, provider)
	if _, err := c.Write([]byte(msg)); err != nil {
		log.Println("Error writing to player: ", err)
	}
}

func (c *Character) ToggleFlagAndMsg(flagName string, provider string, msg string) {
	c.ToggleFlag(flagName, provider)
	if _, err := c.Write([]byte(msg)); err != nil {
		log.Println("Error writing to player: ", err)
	}
}

func (c *Character) FindFlagProviders(flagName string) []string {
	if _, exists := c.Flags[flagName]; exists {
		return c.FlagProviders[flagName]
	}
	return []string{}
}

func (c *Character) CheckFlag(flagName string) bool {
	if flag, ok := c.Flags[flagName]; ok {
		return flag
	}
	return false
}

// SerialSaveEffects serializes all current user effects, removes them, and saves them to the database
func (c *Character) SerialSaveEffects() string {
	effectList := make(map[string]map[string]interface{})

	for efN, effect := range c.Effects {
		// Ignore any bard songs, we won't take them with us.
		if !strings.Contains(efN, "_song") && !strings.Contains(efN, "sing") {
			effectList[efN] = effect.ReturnEffectProps()
		}
	}

	dataJson, err := json.Marshal(effectList)
	if err != nil {
		return "[]"
	} else {
		return string(dataJson)
	}
}

func (c *Character) SerialRestoreEffects(effectsBlob string) {
	obj := make(map[string]map[string]interface{})
	err := json.Unmarshal([]byte(effectsBlob), &obj)
	if err != nil {
		return
	}
	for name, properties := range obj {
		if !strings.Contains(name, "_song") && !strings.Contains(name, "sing") {
			Effects[name](c, c, int(properties["magnitude"].(float64)))
			c.Effects[name].AlterTime(properties["timeRemaining"].(float64))
		}
	}
}

func (c *Character) PurgeEffects() {
	for _, effect := range c.Effects {
		effect.effectOff()
	}
}

// SerialSaveTimers serializes all current user timers
func (c *Character) SerialSaveTimers() string {
	timerList := make(map[string]float64)

	for efN, effect := range c.Timers {
		timerList[efN] = math.Ceil(effect.Sub(time.Now()).Seconds())
	}

	dataStr, err := json.Marshal(timerList)
	if err != nil {
		return "[]"
	} else {
		return string(dataStr)
	}
}

func (c *Character) SerialRestoreTimers(timerBlob string) {
	obj := make(map[string]float64)
	err := json.Unmarshal([]byte(timerBlob), &obj)
	if err != nil {
		return
	}
	for name, duration := range obj {
		c.SetTimer(name, int(duration))
	}
}

func (c *Character) GetModifier(modifier string) int {
	if mod, ok := c.Modifiers[modifier]; ok {
		return mod
	} else {
		return 0
	}
}

// HasShield reports whether a shield is equipped in the offhand.
func (c *Character) HasShield() bool {
	return c.Equipment.Off != (*Item)(nil) && c.Equipment.Off.ItemType == 23
}

// Rescuer returns the living paladin in this room still within their rescue
// window for this character, or nil.
func (c *Character) Rescuer() *Character {
	if c.RescuedBy == "" || time.Now().After(c.RescueUntil) {
		return nil
	}
	for _, ch := range Rooms[c.ParentId].Chars.Contents {
		if ch != c && ch.Name == c.RescuedBy && ch.Vit.Current > 0 {
			return ch
		}
	}
	return nil
}

func (c *Character) SetModifier(modifier string, value int) {
	if _, ok := c.Modifiers[modifier]; ok {
		c.Modifiers[modifier] += value
	} else {
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
		if c.Class == config.MONK {
			return c.Modifiers["armor"] + config.MonkNaturalArmor(c.Tier, c.GetStat("con"))
		}
		return c.Equipment.Armor + c.Modifiers["armor"] + int(float64(c.Con.Current)*config.ConArmorMod*float64(c.Equipment.Armor))
	default:
		return 0
	}
}

func (c *Character) Save() {
	c.MinutesPlayed += int(time.Now().Sub(c.LoginTime).Minutes())
	c.LoginTime = time.Now()
	charData := make(map[string]interface{})
	charData["title"] = c.Title
	charData["name"] = c.Name
	charData["tier"] = c.Tier
	charData["class"] = c.Class
	charData["character_id"] = c.CharId
	charData["experience"] = c.Experience.Value
	charData["spells"] = strings.Join(c.Spells, ",")
	charData["thrustexp"] = c.Skills[1].Value
	charData["bluntexp"] = c.Skills[2].Value
	charData["missileexp"] = c.Skills[4].Value
	charData["handexp"] = c.Skills[5].Value
	charData["poleexp"] = c.Skills[3].Value
	charData["sharpexp"] = c.Skills[0].Value
	charData["fireexp"] = c.Skills[6].Value
	charData["airexp"] = c.Skills[7].Value
	charData["earthexp"] = c.Skills[8].Value
	charData["waterexp"] = c.Skills[9].Value
	charData["divinity"] = c.Skills[10].Value
	charData["stealthexp"] = c.Skills[11].Value
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
	// A prepared weapon saves with the inventory, quickdraw pairings don't survive a logout.
	charData["inventory"] = c.Inventory.JsonifyWith(c.Equipment.Prepared)
	charData["effects"] = c.SerialSaveEffects()
	charData["timers"] = c.SerialSaveTimers()
	charData["lastrefresh"] = c.LastRefresh.Format(time.RFC3339)
	charData["oocswap"] = c.OOCSwap
	charData["bonuspoints"] = c.BonusPoints.Value
	charData["ooc"] = utils.Btoi(c.Flags["ooc"])
	charData["enchants"] = 0
	charData["heals"] = 0
	charData["restores"] = 0
	charData["rerolls"] = c.Rerolls
	if c.Class == 4 || c.Class == 5 {
		charData["enchants"] = c.ClassProps["enchants"]
	}
	if c.Class == 5 || c.Class == 6 {
		charData["heals"] = c.ClassProps["heals"]
	}
	if c.Class == 7 || c.Class == 5 {
		charData["restores"] = c.ClassProps["restores"]
	}
	data.SaveChar(charData)
}

func (c *Character) TickSaveWrapper() {
	// Look the room up once. The lock can be a wait - a death script holding the
	// room they died in, say - and by the time it is ours that script may have
	// moved them on. Unlocking Rooms[c.ParentId] afterwards would then release
	// the room they were moved to, which this never locked: either a fatal
	// "unlock of unlocked mutex" here, or it silently frees whatever command is
	// holding that room, and that command's own unlock is the one that crashes.
	room := Rooms[c.ParentId]
	room.LockRoom(c.Name+":SaveWrap", false)
	if _, err := c.Write([]byte(text.Info + "Saving...." + text.Reset)); err != nil {
		log.Println("Error writing to player: ", err)
	}
	c.Save()
	if _, err := c.Write([]byte(text.Info + "Saved!" + text.Reset)); err != nil {
		log.Println("Error writing to player: ", err)
	}
	room.UnlockRoom(c.Name+":SaveWrap", false)
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
		return []byte(text.Prompt + ">" + text.Reset + "\n")
	case StyleStat:
		if c.CheckFlag("berserk") {
			return []byte((text.Prompt + "(" + text.Red + "Berserk!!" + text.Prompt + "): " + text.Reset + "\n"))
		}
		return []byte((text.Prompt + "(" + text.Yellow +
			strconv.Itoa(c.Stam.Current) + "|" +
			text.Red + strconv.Itoa(c.Vit.Current) + "|" +
			text.Cyan + strconv.Itoa(c.Mana.Current) +
			text.Prompt + "): " + text.Reset + "\n"))
	default:
		return []byte{}
	}
}

// WriteBuffer writes the specified byte slice to the associated client.
func (c *Character) Write(b []byte) (n int, err error) {

	if !c.CheckFlag("allow_writes") {
		return 0, nil
	}
	if c == nil {
		return
	}

	b = append(b, c.buildPrompt()...)
	if c != nil {
		if n, err := c.Writer.Write(b); err != nil {
			log.Println(err)
			return n, err
		}
	}
	return
}

func (c *Character) ReturnVictim() string {
	switch c.Victim.(type) {
	case *Character:
		target := c.Victim.(*Character)
		return target.Name + target.ReturnStateMasked(c.CheckFlag("berserk")) + "," + utils.WhereAt(target.Placement, c.Placement)
	case *Mob:
		target := c.Victim.(*Mob)
		return target.Name + target.ReturnState() + "," + utils.WhereAt(target.Placement, c.Placement)
	default:
		return "No victim."
	}
}

func (c *Character) LookVictim() *Mob {
	switch c.Victim.(type) {
	case *Mob:
		target := c.Victim.(*Mob)
		if target.ParentId == c.ParentId {
			return target
		}
		return nil
	default:
		return nil
	}
}

// ReturnState describes the character's full visible condition: stamina, vitality, and effects.
func (c *Character) ReturnState() string {
	return c.ReturnStateMasked(false)
}

// ReturnStateMasked describes the character's condition. When hideStamAndDrunk is true (a berserk
// viewer), the stamina word and drunkenness are omitted; vitality and other effects still show.
// A berserk character always reads as "appears berserk", followed by their injury level if they
// are hurt and any other visible statuses; the rage hides their stamina.
func (c *Character) ReturnStateMasked(hideStamAndDrunk bool) string {
	stamStatus := text.Green + "energetic" + text.Info
	vitStatus := text.Green + "healthy" + text.Info
	vitHurt := c.Vit.Current < (c.Vit.Max - int(.10*float32(c.Vit.Max)))
	effectStatus := ""
	if c.Stam.Current < (c.Stam.Max - int(.75*float32(c.Stam.Max))) {
		stamStatus = text.Red + "exhausted" + text.Info
	} else if c.Stam.Current < (c.Stam.Max - int(.5*float32(c.Stam.Max))) {
		stamStatus = text.LightYellow + "fatigued" + text.Info
	} else if c.Stam.Current < (c.Stam.Max - int(.25*float32(c.Stam.Max))) {
		stamStatus = text.LightGreen + "slightly fatigued" + text.Info
	}

	if c.Vit.Current < (c.Vit.Max - int(.75*float32(c.Vit.Max))) {
		vitStatus = text.BrightRed + "mortally wounded" + text.Info
	} else if c.Vit.Current < (c.Vit.Max - int(.5*float32(c.Vit.Max))) {
		vitStatus = text.Red + "badly injured" + text.Info
	} else if c.Vit.Current < (c.Vit.Max - int(.25*float32(c.Vit.Max))) {
		vitStatus = text.LightYellow + "injured" + text.Info
	} else if c.Vit.Current < (c.Vit.Max - int(.10*float32(c.Vit.Max))) {
		vitStatus = text.LightGreen + "slightly injured" + text.Info
	}

	if c.CheckFlag("poisoned") {
		effectStatus = effectStatus + " and " + text.LightRed + "poisoned" + text.Info
	}
	if c.CheckFlag("disease") {
		effectStatus = effectStatus + " and " + text.LightRed + "diseased" + text.Info
	}
	if c.CheckFlag("blind") {
		effectStatus = effectStatus + " and " + text.LightRed + "blinded" + text.Info
	}
	if !hideStamAndDrunk {
		if state := DrunkState(c.DrunkLevel()); state != "" {
			effectStatus = effectStatus + " and " + text.LightYellow + state + text.Info
		}
	}

	if c.CheckFlag("berserk") {
		state := " appears " + text.Red + "berserk" + text.Info
		if vitHurt {
			state += " and " + vitStatus
		}
		return state + effectStatus
	}

	if hideStamAndDrunk {
		if c.Vit.Current == c.Vit.Max && effectStatus == "" {
			return " appears untouched"
		}
		return " looks " + vitStatus + effectStatus
	}

	if c.Stam.Current == c.Stam.Max && c.Vit.Current == c.Vit.Max && effectStatus == "" {
		return " appears untouched"
	}

	return " looks " + stamStatus + " and " + vitStatus + effectStatus

}

type PromptStyle int

const (
	StyleNone = iota
	StyleStat
)

func (c *Character) Tick() {
	if time.Now().Sub(c.LastSave) > 5*time.Minute {
		c.LastSave = time.Now()
		c.TickSaveWrapper()
	}
	roomMod := 1.0
	if Rooms[c.ParentId].Flags["heal_fast"] {
		roomMod = config.HealFastRoomRegenMod
	}
	inCombat := c.InCombat()
	// The seal of faith carries the regen half of bless for paladins, who
	// cannot cast bless themselves.
	blessed := c.CheckFlag("bless") || c.CheckFlag("seal-faith")
	c.Heal(config.HealthRegen(c.GetStat("con"), roomMod, inCombat, blessed))
	if c.Class == config.MONK {
		// Chi regenerates only while meditating; otherwise it comes from
		// fighting and bleeds away between fights.
		if c.CheckFlag("meditate") {
			regen := config.ManaRegen(c.GetStat("pie"), c.GetStat("int"), c.Tier, roomMod, inCombat, blessed)
			c.RestoreMana(int(math.Ceil(float64(regen) * config.MeditateChiRegenMod)))
		} else {
			c.tickChi()
		}
	} else {
		c.RestoreMana(config.ManaRegen(c.GetStat("pie"), c.GetStat("int"), c.Tier, roomMod, inCombat, blessed))
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

// Look Drop out the description of this character
func (c *Character) Look() (buildText string) {
	buildText = "You see " + c.Name + ", the " + config.TextGender[c.Gender] + ", " + config.AvailableRaces[c.Race] + " " + c.ClassTitle + "."
	return buildText
}

func (c *Character) ApplyEffect(effectName string, length string, interval int, magnitude int, effect func(triggers int), effectOff func()) {
	if effectInstance, ok := c.Effects[effectName]; ok {
		effectInstance.Reset(length)
		return
	}
	c.Effects[effectName] = NewEffect(length, interval, magnitude, effect, effectOff)
	c.Effects[effectName].RunEffect()
}

// ExtendEffect lengthens a running effect by seconds, capped at maxLength
// seconds in total, and reports the seconds actually added. A missing effect
// adds nothing.
func (c *Character) ExtendEffect(effectName string, seconds int, maxLength int) int {
	if effectInstance, ok := c.Effects[effectName]; ok {
		return effectInstance.Extend(seconds, maxLength)
	}
	return 0
}

func (c *Character) RemoveEffect(effectName string) {
	if _, ok := c.Effects[effectName]; ok {
		c.Effects[effectName].effectOff()
		delete(c.Effects, effectName)
	}
}

func (c *Character) GainExperience(value int) {
	if c.Experience.Value < config.TierExpLevels[config.LevelCap] {
		c.Experience.Add(value)
	}
}

func (c *Character) CanEquip(item *Item) (bool, string) {
	if c.Class >= 99 {
		return true, ""
	}
	if c.Class == config.MONK {
		//check if weapon
		if utils.IntIn(item.ItemType, []int{0, 1, 2, 3, 4}) {
			return false, "You cannot wield weapons effectively."
		}
		//Check if armor and has value greater than 0
		if utils.IntIn(item.ItemType, []int{5, 19, 20, 21, 22, 23, 24, 25, 26}) && item.Armor > 0 {
			return false, "This armor would disrupt the flow of your chi"
		}
	}
	if utils.IntIn(item.ItemType, []int{5, 19, 20, 21, 22, 23, 24, 25, 26}) {
		if !utils.IntIn(item.ArmorClass, config.ArmorRestrictions[c.Class]) {
			return false, "You cannot equip armor of this class."
		}
		if !config.IgnoresGearCaps(c.Tier) && !config.CheckArmor(item.ItemType, c.Tier, item.BaseArmor()) {
			return false, "You are unsure of how to maximize the benefit of this armor and cannot wear it."
		}

	}
	if utils.IntIn(item.ItemType, []int{0, 1, 2, 3, 4, 16}) &&
		!c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		if item.IsTwoHanded() && !config.CanWieldTwoHanded(c.Class) {
			return false, "Only fighters, barbarians and paladins can handle a two-handed weapon."
		}
		if !config.IgnoresGearCaps(c.Tier) && !config.CanWield(c.Tier, c.Class, utils.RollMax(item.SidesDice, item.NumDice, item.PlusDice), item.Flags["two_handed"]) {
			return false, "You are not well enough trained to wield " + item.Name
		}
	}

	return true, ""
}

// DrunkLevel returns the character's current drunk number, 0 if sober
func (c *Character) DrunkLevel() int {
	if effectInstance, ok := c.Effects["drunk"]; ok {
		return effectInstance.magnitude
	}
	return 0
}

func (c *Character) HasEffect(effectName string) bool {
	if _, ok := c.Effects[effectName]; ok {
		return true
	}
	return false
}

// Afflictions are the negative statuses meditate clears and the meditative
// save guards against. Each is an effect name on the character.
var Afflictions = []string{"poison", "disease", "blind", "drunk"}

// ClearAfflictions removes every active affliction and reports which ones
// were cleared. Each effect's own effectOff tells the player it lifted.
func (c *Character) ClearAfflictions() []string {
	var cleared []string
	for _, name := range Afflictions {
		if c.HasEffect(name) {
			c.RemoveEffect(name)
			cleared = append(cleared, name)
		}
	}
	return cleared
}

// MeditativeSave reports whether a meditating monk shrugs off an incoming
// affliction from a source of the given level. It only ever consults flags
// and plain fields, so it is safe to call from the mob goroutine. On a save
// the monk is told what they resisted.
func (c *Character) MeditativeSave(affliction string, attackerLevel int) bool {
	if c.Class != config.MONK || c.Tier < config.MonkMeditateSaveTier || !c.CheckFlag("meditate") {
		return false
	}
	if utils.Roll(100, 1, 0) > config.MeditateSaveChance(c.Tier, attackerLevel, c.GetStat("con")) {
		return false
	}
	if _, err := c.Write([]byte(text.Cyan + "Your trance holds; the " + affliction + " washes over you and finds no purchase.\n")); err != nil {
		log.Println("Error writing to player:", err)
	}
	return true
}

func (c *Character) ApplyHook(hook string, hookName string, executions int, length string, interval int, effect func(), effectOff func()) {
	c.Hooks[hook][hookName] = NewHook(executions, length, interval, effect, effectOff)
}

func (c *Character) RemoveHook(hook string, hookName string) {
	c.Hooks[hook][hookName].effectOff()
	valPresent := false
	for k := range c.Hooks {
		valPresent = false
		for hName := range c.Hooks[k] {
			if hName == hookName {
				valPresent = true
			}
		}
		if valPresent {
			delete(c.Hooks[k], hookName)
		}
	}
}

func (c *Character) RunHook(hook string) {
	if hook == "attacked" {
		c.MarkCombat()
	}
	for name, hookInstance := range c.Hooks[hook] {
		// Process Removing the hook
		if hookInstance.TimeRemaining() == 0 {
			c.RemoveHook(hook, name)
			continue
		}
		if hookInstance.interval > 0 {
			if hookInstance.LastTriggerInterval() <= 0 {
				hookInstance.RunHook()
			}
		} else if hookInstance.interval == -1 {
			hookInstance.RunHook()
		}
	}
	return
}

// AdvanceSkillExp credits weapon skill from a raw exp amount: hand-to-hand
// for monks, otherwise the slot matching the wielded weapon's type. The
// class's advancement multiplier for that slot is applied here.
func (c *Character) AdvanceSkillExp(rawAmount float64) {
	slot := -1
	if c.Class == config.MONK {
		slot = 5
	} else if c.Equipment.Main != nil && utils.IntIn(c.Equipment.Main.ItemType, []int{0, 1, 2, 3, 4}) {
		slot = c.Equipment.Main.ItemType
	}
	if slot < 0 {
		return
	}
	if c.Skills[slot].Value < config.SkillCap {
		c.Skills[slot].Add(int(rawAmount * config.WeaponAdvancementFor(c.Class, slot)))
	}
}

func (c *Character) AdvanceElementalExp(amount int, element string, class int) {
	if class == 4 {
		if c.Skills[magicSkillMap[element]].Value < config.SkillCap {
			c.Skills[magicSkillMap[element]].Add(amount)
		}
	}
	return
}

func (c *Character) AdvanceDivinity(amount int, class int) {
	if class == 5 || class == 6 {
		if c.Skills[10].Value < config.SkillCap {
			c.Skills[10].Add(amount)
		}
	}
	return
}

func (c *Character) AdvanceStealthExp(amount int) {
	if c.Skills[11].Value < config.SkillCap {
		c.Skills[11].Add(amount)
		return
	}
}

// ReceiveDamage Return stam and vital damage
func (c *Character) ReceiveDamage(damage int) (int, int, int) {
	if c.CheckFlag("surge") {
		damage += int(math.Ceil(float64(damage) * config.SurgeExtraDamage))
	}
	if c.CheckFlag("inertial-barrier") {
		damage -= int(math.Ceil(float64(damage) * config.InertialBarrierIgnore(c.Tier, c.Int.Current)))
	}
	stamDamage, vitalDamage := 0, 0
	mitigation := float64(config.ArmorReductionConstant) / (float64(config.ArmorReductionConstant) + float64(c.GetStat("armor")))

	finalDamage := int(math.Ceil(float64(damage) * mitigation))
	resist := damage - finalDamage
	msg := c.Equipment.DamageRandomArmor()
	if msg != "" {
		if _, err := c.Write([]byte(text.Info + msg + "\n" + text.Reset)); err != nil {
			log.Println("Error writing to player: ", err)
		}
	}
	if finalDamage < 0 {
		finalDamage = 0
	}
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
	return stamDamage, vitalDamage, resist
}

func (c *Character) ReceiveDamageNoArmor(damage int) (int, int) {
	stamDamage, vitalDamage := 0, 0
	finalDamage := damage
	if finalDamage < 0 {
		finalDamage = 0
	}
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

func (c *Character) ReceiveVitalDamage(damage int) (int, int) {
	if c.CheckFlag("inertial-barrier") {
		damage -= int(math.Ceil(float64(damage) * config.InertialBarrierIgnore(c.Tier, c.Int.Current)))
	}
	msg := c.Equipment.DamageRandomArmor()
	if msg != "" {
		if _, err := c.Write([]byte(text.Info + msg + "\n" + text.Reset)); err != nil {
			log.Println("Error writing to player: ", err)
		}
	}
	mitigation := float64(config.ArmorReductionConstant) / (float64(config.ArmorReductionConstant) + float64(c.GetStat("armor")))

	finalDamage := int(math.Ceil(float64(damage) * mitigation))
	resist := damage - finalDamage
	if finalDamage < 0 {
		finalDamage = 0
	}
	if finalDamage > c.Vit.Current {
		finalDamage = c.Vit.Current
		c.Vit.Current = 0
	} else {
		c.Vit.Subtract(finalDamage)
	}
	return finalDamage, resist
}

func (c *Character) ReceiveVitalDamageNoArmor(damage int) int {
	msg := c.Equipment.DamageRandomArmor()
	if msg != "" {
		if _, err := c.Write([]byte(text.Info + msg + "\n" + text.Reset)); err != nil {
			log.Println("Error writing to player: ", err)
		}
	}
	if damage < 0 {
		damage = 0
	}
	if damage > c.Vit.Current {
		damage = c.Vit.Current
		c.Vit.Current = 0
	} else {
		c.Vit.Subtract(damage)
	}
	return damage
}

// ReceiveMagicDamage applies a magic hit of the given spell tier (0 for
// breath, damage over time and scripts; see config.SpellTier).
func (c *Character) ReceiveMagicDamage(damage int, element string, tier int) (int, int, int) {
	resisting := 0.0
	// dodge spell
	if c.CheckFlag("dodge") {
		// Did they dodge completely?
		if utils.Roll(100, 1, 0) <= (int(float64(c.GetStat("dex")) * config.FullDodgeChancePerDex)) {
			if _, err := c.Write([]byte(text.Info + "You dodge the spell completely!\n" + text.Reset)); err != nil {
				log.Println("Error writing to player: ", err)
			}
			return 0, 0, 0
		} else {
			if _, err := c.Write([]byte(text.Info + "Your magically quickened reflexes allow you to lessen the effect of the magic!\n" + text.Reset)); err != nil {
				log.Println("Error writing to player: ", err)
			}
			resisting = float64(c.GetStat("dex")) * config.DodgeDamagePerDex
		}
	}

	switch element {
	case "fire":
		if c.CheckFlag("resist-fire") {
			resisting += float64(config.ElementResistPercent) / 100
		}
	case "air":
		if c.CheckFlag("resist-air") {
			resisting += float64(config.ElementResistPercent) / 100
		}
	case "earth":
		if c.CheckFlag("resist-earth") {
			resisting += float64(config.ElementResistPercent) / 100
		}
	case "water":
		if c.CheckFlag("resist-water") {
			resisting += float64(config.ElementResistPercent) / 100
		}
	}

	if c.CheckFlag("resist-magic") {
		resisting += .10
	}

	resisted := int(math.Ceil(float64(damage) * resisting))
	damage -= config.IntMagicResist(c.GetStat("int"), tier)
	damage -= resisted
	// A monk's chi ward takes its share of whatever the other resists let through.
	if c.Class == config.MONK && c.Tier >= config.MonkMagicResistTier {
		absorbed, cost := config.ChiWard(damage, config.MonkMagicResistPercent(c.GetStat("con")), c.Mana.Current)
		if absorbed > 0 {
			c.Mana.Subtract(cost)
			damage -= absorbed
			resisted += absorbed
			if _, err := c.Write([]byte(text.Cyan + "Your chi turns aside " + strconv.Itoa(absorbed) + " damage." + text.Reset + "\n")); err != nil {
				log.Println("Error writing to player: ", err)
			}
		}
	}
	stamDam, vitDam := c.ReceiveDamageNoArmor(damage)
	return stamDam, vitDam, resisted
}

func (c *Character) ReceiveEnvironmentalDamage(damage int, element string) (int, int, int) {
	stamDam, vitDam := c.ReceiveDamageNoArmor(damage)
	return stamDam, vitDam, 0
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

func (c *Character) HealVital(damage int) int {
	if damage > c.Vit.Max-c.Vit.Current {
		damage = c.Vit.Max - c.Vit.Current
		c.Vit.Current = c.Vit.Max
		return damage
	} else {
		c.Vit.Add(damage)
		return damage
	}

}

func (c *Character) HealStam(damage int) int {
	if damage > c.Stam.Max-c.Stam.Current {
		damage = c.Stam.Max - c.Stam.Current
		c.Stam.Current = c.Stam.Max
		return damage
	} else {
		c.Stam.Add(damage)
		return damage
	}
}

func (c *Character) RestoreMana(damage int) {
	c.Mana.Add(damage)
}

// HealPiety is the piety value healing spells scale from. Bards may lean on
// intelligence: max(pie, (int+pie)/2).
func (c *Character) HealPiety() int {
	pie := c.GetStat("pie")
	if c.Class == config.BARD {
		return int(math.Max(float64(pie), float64((c.GetStat("int")+pie)/2)))
	}
	return pie
}

// DivinityBonus is the percentage healing bonus the caster's divinity skill
// grants. Paladins receive only a share of it unless the seal of faith is
// active.
func (c *Character) DivinityBonus() float64 {
	bonus := float64(config.HealingSkill[config.WeaponLevel(c.Skills[config.DivinitySkill].Value, c.Class, config.DivinitySkill)])
	if c.Class == config.PALADIN && !c.CheckFlag("seal-faith") {
		bonus *= config.PaladinDivinityMod
	}
	return bonus
}

// SingCastMod is the damage and healing multiplier for a spell this
// character casts; only a low-tier singing bard pays a penalty.
func (c *Character) SingCastMod() float64 {
	return config.BardSingCastMod(c.Class, c.Tier, c.CheckFlag("singing"))
}

func (c *Character) CalcHealPenalty(damage int) int {
	if c.GetStat("pie") <= config.PieMajorPenalty {
		damage -= int(float64(damage) * (.10 * float64(6-c.GetStat("pie"))))
	}
	return damage
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
	} else {
		return 1
	}
}

func (c *Character) InflictDamage() (damage int) {
	return c.inflictDamage(0)
}

// InflictBackstabDamage is a weapon hit whose dice bonus counts dex on top of
// the weapon's usual stat, before the backstab multiplier.
func (c *Character) InflictBackstabDamage() int {
	return c.inflictDamage(c.GetStat("dex"))
}

// inflictDamage rolls a hit; extraStat points join the weapon stat in the
// percentage bonus to the dice roll.
func (c *Character) inflictDamage(extraStat int) (damage int) {
	if c.Class != config.MONK {
		damage = utils.Roll(c.Equipment.Main.SidesDice,
			c.Equipment.Main.NumDice,
			c.Equipment.Main.PlusDice)
		// Fighters lose their tier-ahead cap at LevelCap, so they get it back as
		// damage on the roll itself.
		if c.Class == config.FIGHTER && config.IgnoresGearCaps(c.Tier) {
			damage += config.MaxTierFighterBonus
		}

		stat := c.GetStat("str")
		if c.Equipment.Main.ItemType == 4 {
			stat = c.GetStat("dex")
		}
		damage += int(math.Ceil(float64(damage) * (config.StatDamageMod * float64(stat+extraStat))))
		damage += c.Equipment.Main.Adjustment
	} else {
		// Monks land the base for their tier, plus a strength share of it
		// (max str is 45), plus a bell-curved roll of a few dice, less a flat
		// cut that bites hardest at low tier.
		baseMonkDamage := config.MonkUnarmedBase(c.Tier)
		strDamage := int(math.Ceil(float64(c.GetStat("str")) / float64(45) * float64(baseMonkDamage)))
		rngDamage := utils.Roll(config.MonkUnarmedRollSides(c.Tier), config.MonkDamageDice, 0)
		damage = baseMonkDamage + strDamage + rngDamage - config.MonkDamageCut(c.Tier)
		if damage < config.MonkDamageFloor {
			damage = config.MonkDamageFloor
		}
	}

	if c.CheckFlag("surge") {
		damage += int(math.Ceil(float64(damage) * config.SurgeDamageBonus))
	}
	if c.CheckFlag("seal-justice") {
		damage += int(math.Ceil(float64(damage) * c.SealJusticeBonus()))
	}

	if c.CheckFlag("reckless") {
		damage += int(math.Ceil(float64(damage) * float64(config.RecklessDamagePercent) / 100))
	}

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

// SealJusticeBonus is the fraction of extra damage the seal of justice adds:
// a flat base plus shares of tier, piety and the wielded weapon's skill level.
func (c *Character) SealJusticeBonus() float64 {
	weaponLevel := 0
	if c.Equipment.Main != (*Item)(nil) {
		weaponLevel = config.WeaponLevel(c.Skills[c.Equipment.Main.ItemType].Value, c.Class, c.Equipment.Main.ItemType)
	}
	tierPart := math.Max(0, float64(c.Tier-config.SealJusticeTierStart)) / config.SealJusticeTierDiv
	percent := tierPart + float64(c.GetStat("pie"))/config.SealJusticePieDiv + float64(weaponLevel)/config.SealJusticeWeaponDiv
	return percent / 100
}

func (c *Character) MaxWeight() int {
	return config.MaxWeight(c.Str.Current)
}

// RefreshDue reports whether the last daily refresh happened before 0:00 UTC today.
func RefreshDue(last, now time.Time) bool {
	y, m, d := now.UTC().Date()
	return last.Before(time.Date(y, m, d, 0, 0, 0, 0, time.UTC))
}

func (c *Character) Refresh() {
	c.Broadcasts = config.BaseBroads + (c.GetStat("int") * config.IntBroad)
	c.Evals = config.BaseEvals + (int(math.Ceil(float64(c.GetStat("int")) / float64(config.IntEvalDivInt))))
	if c.Class == 4 || c.Class == 5 {
		c.ClassProps["enchants"] = 3
	}
	if c.Class == 5 || c.Class == 6 {
		c.ClassProps["heals"] = 5
	}
	if c.Class == 7 || c.Class == 5 {
		c.ClassProps["restores"] = 5
	}
}

func (c *Character) WriteMovement(previous int, new int, subject string) {
	mvAmnt := math.Abs(float64(previous - new))
	color := text.Yellow
	message := ""
	// Moving backwards
	if (previous > new) && (mvAmnt == 1) && (new > c.Placement) {
		message = " moves backwards, towards you." + text.Reset + "\n"
	} else if (previous > new) && (mvAmnt == 1) && (new < c.Placement) {
		message = " moves backwards, away from you." + text.Reset + "\n"
	} else if (previous > new) && (mvAmnt == 1) && (new == c.Placement) {
		message = " moves backwards, next to you." + text.Reset + "\n"
	} else if (previous > new) && (mvAmnt == 2) && (new > c.Placement) {
		message = " sprints backwards, towards you." + text.Reset + "\n"
	} else if (previous > new) && (mvAmnt == 2) && (new < c.Placement) {
		message = " sprints backwards, away from you." + text.Reset + "\n"
	} else if (previous > new) && (mvAmnt == 2) && (new == c.Placement) {
		message = " sprints backwards, next to you." + text.Reset + "\n"
		// Moving forwards
	} else if (previous < new) && (mvAmnt == 1) && (new < c.Placement) {
		message = " moves forwards, towards you." + text.Reset + "\n"
	} else if (previous < new) && (mvAmnt == 1) && (new > c.Placement) {
		message = " moves forwards, away from you." + text.Reset + "\n"
	} else if (previous < new) && (mvAmnt == 1) && (new == c.Placement) {
		message = " moves forwards, next to you." + text.Reset + "\n"
	} else if (previous < new) && (mvAmnt == 2) && (new < c.Placement) {
		message = " sprints forwards, towards you." + text.Reset + "\n"
	} else if (previous < new) && (mvAmnt == 2) && (new > c.Placement) {
		message = " sprints forwards, away from you." + text.Reset + "\n"
	} else if (previous < new) && (mvAmnt == 2) && (new == c.Placement) {
		message = " sprints forwards, next to you." + text.Reset + "\n"
	}
	if _, err := c.Write([]byte(color + subject + message)); err != nil {
		log.Println("Error writing to player: ", err)
	}
}

func (c *Character) LoseParty() {
	if len(c.PartyFollowers) > 0 {
		for _, player := range c.PartyFollowers {
			char := ActiveCharacters.Find(player)
			if char != nil {
				char.PartyFollow = ""
				if _, err := char.Write([]byte(text.Info + c.Name + " loses you." + text.Reset + "\n")); err != nil {
					log.Println("Error writing to player: ", err)
				}
			}

		}
		c.PartyFollowers = []string{}
	}
	return
}

func (c *Character) PartyFollowersNoGM() (followers []string) {
	for i, char := range c.PartyFollowers {
		if !ActiveCharacters.Find(char).Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			followers = append(followers, c.PartyFollowers[i])
		}
	}
	return followers
}

func (c *Character) Unfollow() {
	if c.PartyFollow != "" {
		leadChar := ActiveCharacters.Find(c.PartyFollow)
		if leadChar != nil {
			for i, char := range leadChar.PartyFollowers {
				if char == c.Name {
					leadChar.PartyFollowers = append(leadChar.PartyFollowers[:i], leadChar.PartyFollowers[i+1:]...)
					if !c.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
						if utils.StringIn(leadChar.Name, ActiveCharacters.List()) {
							if _, err := leadChar.Write([]byte(text.Info + c.Name + " stops following you." + text.Reset + "\n")); err != nil {
								log.Println("Error writing to player: ", err)
							}
						}
					}
					break
				}
			}

		}
		if _, err := c.Write([]byte(text.Info + "You stop following " + c.PartyFollow + text.Reset + "\n")); err != nil {
			log.Println("Error writing to player: ", err)
		}
		c.PartyFollow = ""
	}
}

func (c *Character) MessageParty(msg string, exclude *Character) {
	if len(c.PartyFollowers) > 0 {
		for _, findChar := range c.PartyFollowers {
			char := ActiveCharacters.Find(findChar)
			if char != nil && char != exclude {
				if _, err := char.Write([]byte(msg)); err != nil {
					log.Println("Error writing to player: ", err)
				}
			}
		}
	}
}

func (c *Character) DeathCheck(how string) {
	if c.DeathInProgress {
		return
	} else {
		c.DeathInProgress = true
	}
	if c.Vit.Current <= 0 {
		go Script(c, "$DEATH "+how)
	} else {
		c.DeathInProgress = false
	}
	return
}

func (c *Character) DeathCheckBool(how string) bool {
	if c.DeathInProgress {
		return true
	} else {
		c.DeathInProgress = true
	}
	if c.Vit.Current <= 0 {
		go Script(c, "$DEATH "+how)
		return true
	} else {
		c.DeathInProgress = false
		return false
	}

}
