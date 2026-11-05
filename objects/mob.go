package objects

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jinzhu/copier"
)

// Mob implements a control object for mobs interacting with players and each other
type Mob struct {
	Object
	MobId         int
	Inventory     *ItemInventory
	ItemList      map[int]int
	Flags         map[string]bool
	FlagProviders map[string][]string
	Effects       map[string]*Effect
	Hooks         map[string]map[string]*Hook

	MobType    int // construct/celestial/daemonic/undead/humanoid/beast/draconic/elemental
	MobSubType int // Melee, Ranged, Caster, Support

	// ParentId is the room id for the room
	ParentId int
	Gold     int
	// ExpMultiplier scales config.MobBaseExperience for the mob's level;
	// see Experience.
	ExpMultiplier float64
	Level         int

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
	BreathWeapon    string

	//Threat table attacker -> damage
	TotalThreatDamage int
	ThreatTable       map[string]int
	CurrentTarget     string

	NumWander  int
	TicksAlive int
	WimpyValue int

	MobTickerUnload chan bool
	MobCommands     chan string
	MobTicker       *time.Ticker
	TickModifier    int
	// An int to hold a stun time.
	MobStunned int
	IsStunned  bool
	IsActive   bool
	IsThinking bool
	// StunnedUntil is when the current stun releases the mob's tick.
	StunnedUntil time.Time
	// NextTick is when the mob's ticker is due to fire next. Every reset of
	// MobTicker goes through scheduleTick so this stays accurate; Stun reads
	// it to decide whether a stun would actually delay the mob.
	NextTick time.Time

	// Attackers holds the names of characters who have taken an offensive
	// action against this mob, hit or miss. Unlike ThreatTable it ignores
	// threat from heals, rescues and the mob picking its own target.
	Attackers map[string]bool

	// SpawnedAt is when the mob started ticking. A mob can't follow anyone
	// who hasn't attacked it out of the room until
	// config.MobFollowSpawnLockout has passed.
	SpawnedAt time.Time
	// StoppedAt is when a permanent mob stopped ticking because its room
	// emptied; see Rest.
	StoppedAt time.Time

	// WeakenSteps counts the $WEAKEN offerings already baked into this mob's
	// template stats. It lives on the template so it survives death and
	// reboots, see Mob.Weaken.
	WeakenSteps int
	// WeakenOffered counts offerings made toward the next step when a
	// $WEAKEN asks for several per step, see Mob.Offer.
	WeakenOffered int

	// QuestEvent names the event whose spawn group made this copy, so it can
	// be removed when the event ends. Empty for ordinary mobs.
	QuestEvent string
	// Summoned marks an add a script called in; it carries no loot, event
	// drops included.
	Summoned bool

	// DamageBonusPct raises every hit by that percent, e.g. while a script
	// has the mob enraged.
	DamageBonusPct int

	// scripts is this copy's mob script bookkeeping; see mob_script.go.
	scripts *mobScriptState
}

func LoadMob(mobData map[string]interface{}) (*Mob, bool) {
	description := ""
	var ok bool
	if description, ok = mobData["description"].(string); !ok {
		description = "A mob...  yup"
	}
	commands, _ := mobData["commands"].(string)
	newMob := &Mob{
		Object{
			Name:        mobData["name"].(string),
			Description: description,
			Placement:   int(mobData["placement"].(int64)),
			// Mobs saved before Save wrote commands have a null here.
			Commands: DeserializeCommands(commands),
		},
		int(mobData["mob_id"].(int64)),
		NewItemInventory(),
		make(map[int]int),
		make(map[string]bool),
		make(map[string][]string),
		make(map[string]*Effect),
		map[string]map[string]*Hook{
			"act":      make(map[string]*Hook),
			"combat":   make(map[string]*Hook),
			"peek":     make(map[string]*Hook),
			"gridmove": make(map[string]*Hook),
			"move":     make(map[string]*Hook),
			"say":      make(map[string]*Hook),
			"attacked": make(map[string]*Hook),
		},
		0,
		0,
		-1,
		int(mobData["gold"].(int64)),
		expMultiplier(mobData["exp_multiplier"]),
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
		nil,
		0,
		0,
		false,
		false,
		false,
		time.Time{},
		time.Time{},
		nil,
		time.Time{},
		time.Time{},
		0,
		0,
		"",
		false,
		0,
		nil,
	}
	if steps, ok := mobData["weaken_steps"].(int64); ok {
		newMob.WeakenSteps = int(steps)
	}
	if offered, ok := mobData["weaken_offered"].(int64); ok {
		newMob.WeakenOffered = int(offered)
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

// Rest heals a permanent mob for the time its room sat empty, at
// config.PermMobRegenPctPerMinute, and re-arms the HP triggers it is back
// above. Its target and any windup it was holding are dropped, so it picks
// whoever walks in as soon as it wakes and swings on its first tick.
func (m *Mob) Rest() {
	if m.StoppedAt.IsZero() {
		return
	}
	m.CurrentTarget = ""
	st := m.scriptState()
	st.windup = nil
	st.lastTarget = ""
	if m.Stam.Max <= 0 {
		m.StoppedAt = time.Time{}
		return
	}
	minutes := time.Since(m.StoppedAt).Minutes()
	m.StoppedAt = time.Time{}
	heal := int(minutes * float64(m.Stam.Max*config.PermMobRegenPctPerMinute) / 100)
	m.Stam.Current = min(m.Stam.Max, m.Stam.Current+heal)

	pct := m.Stam.Current * 100 / m.Stam.Max
	for threshold := range st.hpFired {
		if pct >= threshold {
			delete(st.hpFired, threshold)
		}
	}
}

func (m *Mob) StartTicking() {
	if m.IsActive {
		return
	}
	log.Println(m.Name + " not active starting ticking")
	m.IsActive = true
	m.SpawnedAt = time.Now()
	m.CalculateInventory()
	m.ThreatTable = make(map[string]int)
	m.Attackers = make(map[string]bool)
	m.MobTickerUnload = make(chan bool)
	m.MobCommands = make(chan string)
	m.TickModifier = 0
	if fastMoving, ok := m.Flags["fast_moving"]; ok {
		if fastMoving {
			m.TickModifier = 2
		}
	}
	m.MobTicker = time.NewTicker(m.tickPeriod())
	m.NextTick = time.Now().Add(m.tickPeriod())
	/* This allows a mob to cast beneficial spells on itself on load, but it also means that any mob with a heal in its spell list will heal on load, which is not desired.
	for _, spell := range m.Spells {
		if utils.StringIn(spell, MobSupportSpells) {
			Cast(m, m, spell, 0)
		}
	}
	*/
	// Execute Immediately - Do not wrap locks, this is called from an existing lock
	go func() {
		Rooms[m.ParentId].LockRoom(m.Name+"MobInitTick", false)
		m.runSpawnScript()
		m.PickTarget()
		Rooms[m.ParentId].UnlockRoom(m.Name+"MobInitTick", false)
	}()
	go func() {
		for {
			select {
			case <-m.MobTickerUnload:
				m.MobTicker.Stop()
				m.IsActive = false
				return
			case <-m.MobTicker.C:
				Rooms[m.ParentId].LockRoom(m.Name+" MobTick", false)
				// The ticker keeps its current period unless Tick resets it.
				m.NextTick = time.Now().Add(m.tickPeriod())
				m.Tick()
				Rooms[m.ParentId].UnlockRoom(m.Name+" MobTick", false)
			}
		}
	}()
	go func() {
		for {
			select {
			case msg := <-m.MobCommands:
				// This function call will immediately call a command off the stack and run it, ideally to decouple state
				if msg == "" {
					return
				}
				var params = strings.Split(msg, " ")
				go m.ProcessCommand(params[0], params[1:])
			}
		}
	}()

}

func (m *Mob) GetSpellMultiplier() float32 {
	return 1
}

func (m *Mob) CheckThreatTable(charName string) bool {
	if _, ok := m.ThreatTable[charName]; ok {
		return true
	}
	return false
}

// Tick The mob brain is this ticker
func (m *Mob) Tick() {
	// It appears sometimes that mobs get a second wind round if they try to lock at the sametime they die, lets prevent that.
	if m.Stam.Current <= 0 {
		return
	}
	if m.IsStunned {
		m.scheduleTick(m.tickPeriod())
		m.IsStunned = false
		m.MobStunned = 0
	}
	m.TicksAlive++
	// Permanent mobs never wander off, and the wander count doesn't apply to
	// them at all: however long they have stood idle they keep ticking in
	// full, so a hostile one still picks up whoever walks in.
	if m.TicksAlive >= m.NumWander && m.CurrentTarget == "" && !m.Flags["permanent"] {
		go Rooms[m.ParentId].WanderMob(m)
		return
	} else {
		// Picking up treasure
		if m.Flags["take_treasure"] {
			// Roll to see if mob is picking it up
			if utils.Roll(100, 1, 0) <= config.MobTakeChance {
				// Loop inventory, and take the first thing they find
				for _, item := range Rooms[m.ParentId].Items.Contents {
					if m.Placement == item.Placement && !item.Flags["hidden"] && !item.Flags["permanent"] {
						if err := Rooms[m.ParentId].Items.Remove(item); err != nil {
							log.Println("Error mob removing item", err)
						}
						m.Flags["no_steal"] = true
						m.Inventory.Add(item)
						Rooms[m.ParentId].MessageAll(m.Name + " picks up " + item.DisplayName() + text.Reset + "\n")
						break
					}
				}
			}
		}

		if m.CurrentTarget != "" {
			if Rooms[m.ParentId].Chars.SearchAll(m.CurrentTarget) == nil {
				m.CurrentTarget = ""
			}
		}

		// A mob that only now picks a target spends the tick on it: it may
		// close in, but breathes, casts and attacks from the next tick.
		hadTarget := m.CurrentTarget != ""
		m.PickTarget()
		picked := !hadTarget && m.CurrentTarget != ""

		// Do I want to change targets? MobThreatSwitchChance percent per tick
		// if the current target isn't the highest on the threat table.
		if len(m.ThreatTable) > 1 {
			rankedThreats := utils.RankMapStringInt(m.ThreatTable)
			if m.CurrentTarget != rankedThreats[0] {
				if utils.Roll(100, 1, 0) <= config.MobThreatSwitchChance {
					if utils.StringIn(rankedThreats[0], Rooms[m.ParentId].Chars.MobList(m)) {
						m.CurrentTarget = rankedThreats[0]
						Rooms[m.ParentId].MessageAll(m.Name + " turns to " + m.CurrentTarget + "\n" + text.Reset)
					}
				}
			}
		}

		// Scripts go before the mob's own moves; one that took an action
		// was the mob's turn.
		if m.runTickScripts() {
			return
		}

		if m.CurrentTarget == "" && m.Placement != 3 && !m.CheckFlag("immobile") {
			oldPlacement := m.Placement
			if m.Placement > 3 {
				m.Placement--
			} else {
				m.Placement++
			}
			if !m.Flags["hidden"] {
				whichNumber := Rooms[m.ParentId].Mobs.GetNumber(m)
				if len(Rooms[m.ParentId].Mobs.Contents) > 1 && whichNumber > 1 {
					Rooms[m.ParentId].MessageMovement(oldPlacement, m.Placement, m.Name+" #"+strconv.Itoa(whichNumber))
				} else {
					Rooms[m.ParentId].MessageMovement(oldPlacement, m.Placement, m.Name)
				}
			}
			return
		}

		if !picked && m.CurrentTarget != "" && m.BreathWeapon != "" &&
			(math.Abs(float64(m.Placement-Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m).Placement)) == 1) {

			// Roll to see if we're going to breathe
			if utils.Roll(100, 1, 0) <= 30 {
				target := Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m)
				var targets []*Character
				for _, character := range Rooms[m.ParentId].Chars.Contents {
					if character.Placement == target.Placement && !character.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
						log.Println("Adding target: ", character.Name, " to breath list")
						targets = append(targets, character)
					}
				}

				Rooms[m.ParentId].MessageAll("The " + m.Name + " breathes " + m.BreathWeapon + " at " + target.Name + "\n")
				damageTotal := config.BreatheDamage(m.Level)
				for _, t := range targets {
					if utils.StringIn(m.BreathWeapon, []string{"fire", "air", "earth", "water"}) {
						t.RunHook("attacked")
						breathDamage := int(math.Ceil(float64(damageTotal) * m.rollMonkBreathDodge(t)))
						stamDam, vitDam, resisted := t.ReceiveMagicDamage(breathDamage, m.BreathWeapon, 0)
						if _, err := t.Write([]byte(text.Bad + m.Name + "'s breath  struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality. You resisted " + strconv.Itoa(resisted) + "damage." + text.Reset + "\n")); err != nil {
							log.Println("Error writing to player:", err)
						}
						// Once one reflection kills the mob, a second DeathCheck would pay out the kill again.
						if t.CheckFlag("reflection") && m.Stam.Current > 0 {
							reflectDamage := int(float64(damageTotal) * (float64(t.GetStat("int")) * config.ReflectDamagePerInt))
							m.ReceiveDamage(reflectDamage)
							if _, err := t.Write([]byte(text.Cyan + "You reflect " + strconv.Itoa(reflectDamage) + " damage back to " + m.Name + "!\n" + text.Reset)); err != nil {
								log.Println("Error writing to player:", err)
							}
							m.DeathCheck(t)
						}
						t.DeathCheck("was slain by a " + m.Name + ".")
					} else if m.BreathWeapon == "paralytic" {
						if _, err := t.Write([]byte(text.Gray + m.Name + " breathes paralytic gas on to you.\n")); err != nil {
							log.Println("Error writing to player:", err)
						}
						if t.MeditativeSave("paralysis", m.Level) {
							continue
						}
						// Paralysis is a stun: it blocks healing and support casts too.
						t.SetTimer("global", 16)
						t.SetTimer("stun", 16)
						data.StoreCombatMetric("breath-stun", 0, 1, 16, 0, 0, 1, m.MobId, m.Level, 0, t.CharId)
					} else if m.BreathWeapon == "pestilence" {
						if _, err := t.Write([]byte(text.Gray + m.Name + " breathes infectious gas on to you.\n")); err != nil {
							log.Println("Error writing to player:", err)
						}
						Effects["disease"](m, t, m.Level)
					}
				}
				return
			}
		}

		// only hostile mobs or mobs with a current target will try to cast
		if !picked && (m.CheckFlag("hostile") || m.CurrentTarget != "") {
			if m.ChanceCast > 0 && len(m.Spells) > 0 {
				log.Println("Trying to cast a spell")
				if utils.Roll(100, 1, 0) <= m.ChanceCast {
					spellIndices := rand.Perm(len(m.Spells)) // Shuffle spell list
					for _, idx := range spellIndices {
						selectSpell := m.Spells[idx]
						if selectSpell == "" {
							continue
						}
						spellInstance, ok := Spells[selectSpell]
						if !ok || m.Mana.Current < spellInstance.Cost {
							continue
						}
						if utils.StringIn(selectSpell, MobSupportSpells) {
							var mobTargets []*Mob
							for _, mob := range Rooms[m.ParentId].Mobs.Contents {
								// Healing spells: only target mobs not at full health
								//log.Println("current health:", mob.Stam.Current, "max health:", mob.Stam.Max)
								if utils.StringIn(selectSpell, MobHealingSpells) {
									if !mob.CheckFlag(spellInstance.Effect) && mob.Stam.Current < mob.Stam.Max {
										mobTargets = append(mobTargets, mob)
									}
								} else {
									// Non-healing beneficial spells
									if !mob.CheckFlag(spellInstance.Effect) {
										mobTargets = append(mobTargets, mob)
									}
								}
							}
							if len(mobTargets) == 0 {
								continue // Try next spell
							}
							targetMob := mobTargets[rand.Intn(len(mobTargets))]
							Rooms[m.ParentId].MessageAll(m.Name + " casts a " + spellInstance.Name + " spell on " + targetMob.Name + ".\n")
							m.Mana.Subtract(spellInstance.Cost)
							result := Cast(m, targetMob, spellInstance.Effect, spellInstance.Magnitude)
							if strings.Contains(result, "$SCRIPT") {
								m.MobScript(result)
							}
							break // Spell cast, exit loop
						} else if utils.StringIn(selectSpell, OffensiveSpells) && m.CurrentTarget != "" {
							// Offensive spell: find a character target
							var target *Character
							for target == nil && len(Rooms[m.ParentId].Chars.MobList(m)) > 0 {
								target = Rooms[m.ParentId].Chars.MobSearch(utils.RandMapKeySelection(m.ThreatTable), m)
							}
							if target == nil {
								continue // Try next spell
							}
							Rooms[m.ParentId].MessageAll(m.Name + " casts a " + spellInstance.Name + " spell on " + target.Name + "\n")
							target.RunHook("attacked")
							m.Mana.Subtract(spellInstance.Cost)
							result := Cast(m, target, spellInstance.Effect, spellInstance.Magnitude)
							if strings.Contains(result, "$SCRIPT") {
								m.MobScript(result)
							}
							target.DeathCheck("was slain by a " + m.Name + ".")
							return //Offensive spell cast, end turn
						}
					}
				}
			}
		}

		// Physical attack: ranged if flagged and not adjacent, otherwise close
		// in and melee. Everything after the parry check is shared.
		if m.CurrentTarget == "" {
			return
		}
		target := Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m)
		if target == nil {
			return
		}
		distance := math.Abs(float64(m.Placement - target.Placement))
		ranged := m.Flags["ranged_attack"] && distance >= 1

		// A paladin rescuing the target intercepts this attack: a melee mob is
		// pulled onto the paladin and charges them, a ranged mob keeps its
		// target but this shot lands on the paladin instead. A ranged mob
		// already standing on the paladin's square is pulled like a melee one.
		if rescuer := target.Rescuer(); rescuer != nil {
			if !ranged || m.Placement == rescuer.Placement {
				m.Charge(rescuer, target)
				return
			}
			target = rescuer
		}

		if !ranged {
			if (!m.CheckFlag("immobile") && m.Placement != target.Placement) || distance > 1 {
				// Close the gap; attack next tick.
				oldPlacement := m.Placement
				if m.Placement > target.Placement {
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
			if m.Placement != target.Placement {
				// Immobile and out of reach.
				return
			}
		}
		if picked {
			return
		}

		// Am I against a fighter, and they succeed in a parry roll? A shield
		// widens the parry window, but a parry that only the shield caught is
		// a block and never ripostes.
		parried, shieldBlock := false, false
		if target.Class == config.FIGHTER && target.Equipment.Main != nil {
			parried, shieldBlock = config.RollFighterParry(config.WeaponLevel(target.Skills[target.Equipment.Main.ItemType].Value, target.Class, target.Equipment.Main.ItemType), target.Tier, target.HasShield())
		}
		if parried {
			if ranged {
				if shieldBlock {
					target.writeCombat(text.Green + "You catch the attack from " + m.Name + " on your shield\n" + text.Reset)
				} else {
					target.writeCombat(text.Green + "You deflect the attack from " + m.Name + "\n" + text.Reset)
				}
				data.StoreCombatMetric("range-parry", 0, 1, 0, 0, 0, 1, m.MobId, m.Level, 0, target.CharId)
				return
			} else if shieldBlock {
				target.writeCombat(text.Green + "You block the attack from " + m.Name + " with your shield\n" + text.Reset)
				data.StoreCombatMetric("melee_player_shield_block", 0, 1, 0, 0, 0, 1, m.MobId, m.Level, 0, target.CharId)
				// A block only stops the hit; unlike a parry it does not stun.
				return
			} else if target.Tier >= config.SpecialAbilityTier {
				// It's a riposte
				target.RunHook("attacked")
				actualDamage, _, resisted := m.ReceiveDamage(int(math.Ceil(float64(target.InflictDamage()))))
				data.StoreCombatMetric("melee_player_riposte", 0, 1, actualDamage+resisted, resisted, actualDamage, 0, target.CharId, target.Tier, 1, m.MobId)
				target.AdvanceSkillExp((float64(actualDamage) / float64(m.Stam.Max) * float64(m.Experience())))
				target.writeCombat(text.Green + "You parry and riposte the attack from " + m.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + "\n" + text.Reset)
				if m.DeathCheck(target) {
					return
				}
				m.StunBy(config.ParryStuns, "parry-stun", target)
				return
			} else {
				target.writeCombat(text.Green + "You parry the attack from " + m.Name + "\n" + text.Reset)
				m.StunBy(config.ParryStuns, "parry-stun", target)
				return
			}
		}

		metric, hitPrefix := "melee", ""
		if ranged {
			metric, hitPrefix = "range", "Thwwip!! "
		}
		if m.RollMiss(target, metric) {
			return
		}
		target.RunHook("attacked")
		if !ranged {
			// Touch effects (blind/disease/poison/acid) only land in melee.
			m.CheckForExtraAttack(target)
		}
		style, mult := m.RollSpecial(target)
		if !Rooms[m.ParentId].InQuestMode() {
			mult *= config.OutlevelDamageMult(m.Level-target.Tier, m.EngagedCount(target))
		}
		m.ApplyStrike(target, m.InflictDamage(), style, mult, StrikeOpts{
			Metric:    metric,
			Mode:      1,
			HitPrefix: hitPrefix,
			DeathMsg:  "was slain by a " + m.Name + ".",
		})
	}
}

// Charge is the melee half of a rescue: the mob abandons rescued for the
// rescuer, becomes fixed on them as its top threat, closes to their placement
// and lands a hit that ignores armor. An immobile mob that cannot reach the
// rescuer is retargeted but gets no hit.
func (m *Mob) Charge(rescuer *Character, rescued *Character) {
	m.CurrentTarget = rescuer.Name
	top := 0
	for _, threat := range m.ThreatTable {
		if threat > top {
			top = threat
		}
	}
	if want := config.RescueThreat(top); m.ThreatTable[rescuer.Name] < want {
		m.addThreat(want-m.ThreatTable[rescuer.Name], rescuer)
	}
	if m.Placement != rescuer.Placement {
		if m.CheckFlag("immobile") {
			Rooms[m.ParentId].MessageAll(m.Name + " turns to " + rescuer.Name + "\n" + text.Reset)
			return
		}
		oldPlacement := m.Placement
		m.Placement = rescuer.Placement
		if !m.Flags["hidden"] {
			whichNumber := Rooms[m.ParentId].Mobs.GetNumber(m)
			Rooms[m.ParentId].MessageMovement(oldPlacement, m.Placement, m.Name+" #"+strconv.Itoa(whichNumber))
		}
	}
	Rooms[m.ParentId].MessageAll(text.Red + m.Name + " charges at " + rescuer.Name + "!\n" + text.Reset)
	rescuer.RunHook("attacked")
	damage := m.InflictDamage()
	stamDamage, vitDamage := rescuer.ReceiveDamageNoArmor(damage)
	data.StoreCombatMetric("melee_rescue", 0, 1, damage, 0, stamDamage+vitDamage, 1, m.MobId, m.Level, 0, rescuer.CharId)
	rescuer.writeCombat(text.Red + m.Name + " slams into you for " + strconv.Itoa(stamDamage+vitDamage) + " points of damage!\n" + text.Reset)
	rescuer.DeathCheck("was slain by a " + m.Name + " while rescuing " + rescued.Name + ".")
}

// RunDeathScript has the killer run the mob's @DEATH script. Kills happen
// with the room locked, so it waits its turn on a goroutine of its own.
func (m *Mob) RunDeathScript(killer *Character) {
	if val, ok := m.Commands["@DEATH"]; ok && killer != nil && Script != nil {
		go Script(killer, "$RUN "+val.Command)
	}
}

func (m *Mob) DeathCheck(target *Character) bool {
	if m.Stam.Current <= 0 {
		m.RunDeathScript(target)
		Rooms[m.ParentId].MessageAll(target.Name + " killed " + m.Name)
		partyLead := target
		if target.PartyFollow != "" {
			partyLead = ActiveCharacters.Find(target.PartyFollow)
		}
		partyMembers := append(partyLead.PartyFollowers, partyLead.Name)
		highestTier := 0
		expReduce := len(Rooms[m.ParentId].Chars.Contents)
		//reusing this because we are already running through everyone in the room
		for _, gm := range Rooms[m.ParentId].Chars.Contents {
			if gm.Tier > highestTier && gm.Class <= 8 {
				highestTier = gm.Tier
			}
			if gm.Class > 8 {
				expReduce -= 1
			}
		}

		if expReduce > 5 {
			expReduce = 5
		}
		experienceAwarded := m.KillExperience(Rooms[m.ParentId], expReduce)
		for _, member := range Rooms[m.ParentId].Chars.Contents {
			buildActorString := ""
			charClean := Rooms[m.ParentId].Chars.SearchAll(member.Name)
			if charClean != nil {
				partyCheck := false
				for _, name := range partyMembers {
					if charClean.Name == name {
						partyCheck = true
					}
				}
				if Rooms[m.ParentId].InQuestMode() {
					buildActorString += text.Cyan + "You earn " + strconv.Itoa(experienceAwarded) + " experience for the defeat of the " + m.Name + "\n"
					charClean.GainExperience(experienceAwarded)
				} else if partyCheck || m.CheckThreatTable(charClean.Name) {
					if int(math.Ceil((float64(charClean.Tier+1))*1.2)) < highestTier {
						buildActorString += text.Cyan + "You learn nothing for the defeat of the " + m.Name + "\n"
					} else {
						buildActorString += text.Cyan + "You earn " + strconv.Itoa(experienceAwarded) + " experience for the defeat of the " + m.Name + "\n"
						charClean.GainExperience(experienceAwarded)
					}
				}
				if charClean == target {
					buildActorString += text.Green + m.DropInventory() + "\n"
					if _, err := target.Write([]byte(buildActorString)); err != nil {
						log.Println("Error writing to player:", err)
					}
				} else {
					go func() {
						if _, err := charClean.Write([]byte(buildActorString + "\n" + text.Reset)); err != nil {
							log.Println("Error writing to player:", err)
						}
					}()
				}
				if charClean.Victim == m {
					charClean.Victim = nil
				}
			}
		}

		go Rooms[m.ParentId].ClearMob(m)
		return true
	}
	return false
}

// Experience is the mob's full exp value: the per-level table scaled by its
// multiplier. Kill awards, weapon skill and elemental skill exp all start
// from it.
func (m *Mob) Experience() int {
	return int(math.Round(float64(config.MobBaseExperience(m.Level)) * m.ExpMultiplier))
}

// expMultiplier reads exp_multiplier from a mob node. Neo4j hands back an
// int64 when a whole number was written, and nil on a mob that predates the
// field, which gets the plain table value.
func expMultiplier(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case int64:
		return float64(v)
	}
	return 1
}

// KillExperience is the experience each member earns for killing the mob in
// room, shared by melee kills and spell/callback kills. expReduce is the
// number of players in the room, capped at 5.
func (m *Mob) KillExperience(room *Room, expReduce int) int {
	experienceAwarded := 0
	baseExp := m.Experience()
	if room.InQuestMode() {
		experienceAwarded = baseExp
	} else if m.CheckFlag("hostile") {
		experienceAwarded = int(float64(baseExp) * (config.ExperienceReduction[expReduce] + (float64(utils.Roll(10, 1, 0)) / 100)))
	} else {
		experienceAwarded = baseExp / 10
	}

	for flag, modifier := range config.XP_Modifiers {
		if m.CheckFlag(flag) {
			experienceAwarded += int(float64(baseExp) * modifier)
		}
	}

	// Check for elemental damage in the room and add bonus experience
	if room.FlagOn("earth") || room.FlagOn("fire") || room.FlagOn("water") || room.FlagOn("air") {
		experienceAwarded = int(float64(experienceAwarded) * 1.06) // 6% bonus for elemental rooms
	}
	return experienceAwarded
}

func (m *Mob) CheckForExtraAttack(target *Character) {

	if m.Flags["blinds"] {
		if utils.Roll(100, 1, 0) > 90 {
			if _, err := target.Write([]byte(text.Red + m.Name + " blinds you!" + "\n" + text.Reset)); err != nil {
				log.Println("Error writing to player:", err)
			}
			Effects["blind"](m, target, 0)
			return
		}
	}

	if m.Flags["diseases"] {
		if utils.Roll(100, 1, 0) > 50 {
			if _, err := target.Write([]byte(text.Red + m.Name + " tries to spread disease on to you!" + "\n" + text.Reset)); err != nil {
				log.Println("Error writing to player:", err)
			}
			magnitude := m.Level
			if magnitude < 4 {
				magnitude = 4
			}
			Effects["disease"](m, target, magnitude)
			return
		}
	}
	if m.Flags["poisons"] {
		if utils.Roll(100, 1, 0) > 50 {
			if _, err := target.Write([]byte(text.Red + m.Name + " injects you with venom!" + "\n" + text.Reset)); err != nil {
				log.Println("Error writing to player:", err)
			}
			magnitude := m.Level
			if magnitude < 4 {
				magnitude = 4
			}
			Effects["poison"](m, target, m.Level)
			return
		}
	}
	if m.Flags["spits_acid"] {
		if target.CheckFlag("resilient-aura") {
			if _, err := target.Write([]byte(text.Red + m.Name + " spits acid on you, but your aura protects your gear!" + "\n" + text.Reset)); err != nil {
				log.Println("Error writing to player:", err)
			}
			return
		}
		if utils.Roll(100, 1, 0) > 50 {
			if _, err := target.Write([]byte(text.Red + m.Name + " spits acid on you, damaging your armor !" + "\n" + text.Reset)); err != nil {
				log.Println("Error writing to player:", err)
			}
			msg := target.Equipment.DamageRandomArmor()
			if msg != "" {
				if _, err := target.Write([]byte(text.Info + msg + "\n" + text.Reset)); err != nil {
					log.Println("Error writing to player:", err)
				}
			}
			return
		}
	}
	return
}

func (m *Mob) PickTarget() {
	// Am I hostile?  Should I pick a target?
	if m.CurrentTarget == "" && m.Flags["hostile"] && len(Rooms[m.ParentId].Chars.MobList(m)) > 0 {
		for m.CurrentTarget == "" && len(Rooms[m.ParentId].Chars.MobList(m)) > 0 {
			for i := 0; i <= 4; i++ {
				if m.CurrentTarget != "" {
					break
				}
				potentials := Rooms[m.ParentId].Chars.MobListAt(m, i)
				if len(potentials) > 0 {
					if utils.Roll(100, 1, 0) <= config.ProximityChance-(i*config.ProximityStep) {
						potential := utils.RandListSelection(potentials)
						m.addThreat(1, Rooms[m.ParentId].Chars.MobSearch(potential, m))
						Rooms[m.ParentId].MessageAll(m.Name + " attacks " + m.CurrentTarget + text.Reset + "\n")
						break
					}
				}
			}
		}
	}
}

// MobsAllowedIn reports whether a mob may be relocated into this room. The
// resurrection room, the entry room and the OOC lounge have to stay clear: a
// player arriving there has just died, just logged in, or is out of character,
// and none of them can be expected to deal with a mob that chased them in.
func MobsAllowedIn(r *Room) bool {
	switch r.RoomId {
	case config.HealingHand, config.StartingRoom, config.OocRoom:
		return false
	}
	return r.Flags["active"]
}

// FollowChar relocates a hunting mob from one room into another behind the
// character it is chasing, and reports whether it moved.
//
// The caller must already hold the locks on both rooms. That is the point of
// the signature: this is called from the character's own move, which holds
// them, so nothing is locked here. The old version ran on the mob's goroutine
// with nothing held, spun for lock priority on two rooms, and had to unload the
// mob to move it - Mobs.Remove stops the ticker and closes MobCommands, so what
// arrived was a rebuilt mob with an empty threat table, a recalculated
// inventory and a closed channel that any later send would panic on. This moves
// the same mob and leaves its goroutines alone; they pick up the new ParentId
// on the next tick.
//
// The destination is passed in rather than looked up from the character, so a
// mob can only ever arrive in the room on the other side of the exit the
// character used. A mob can no longer end up wherever the character happens to
// be by the time it gets around to chasing them - the healing hand, most of
// all, which is where they are if they died on the way out.
func (m *Mob) FollowChar(target *Character, from *Room, to *Room) bool {
	if m.ParentId != from.RoomId || !MobsAllowedIn(to) {
		return false
	}

	if !m.CheckFlag("follows") || m.CheckFlag("curious_canticle") || m.CheckFlag("run_away") {
		return false
	}

	// Permanent mobs belong to their room - they survive a room clear, are the
	// ones Jsonify writes back out, and are restarted in place by RestartPerms.
	// One that wandered off behind a player would be saved into the wrong room
	// and leave a hole in the one it was built for.
	if m.CheckFlag("permanent") {
		return false
	}

	// Still angry at this one in particular, and in any shape to give chase?
	if m.CurrentTarget != target.Name || m.Stunned() || m.FollowLocked(target.Name) || m.Stam.Current <= 0 {
		return false
	}

	curChance := config.MobFollow - ((target.Tier - m.Level) * config.MobFollowPerLevel)
	if curChance > 85 {
		curChance = 85
	}
	if m.CheckFlag("crippled") {
		curChance = config.CrippledChance(curChance)
	}
	if utils.Roll(100, 1, 0) > curChance {
		return false
	}

	from.MessageAll(m.Name + " follows " + target.Name + "\n\n")
	from.Mobs.Release(m)
	// Add sets ParentId to the new room for us.
	to.Mobs.AddWithMessage(m, m.Name+" follows "+target.Name+" into the area.", false)
	m.Placement = 3
	return true
}

func (m *Mob) Flee(params []string) {
	log.Println("I'm gonna try to flee", params)
	if m.Stam.Current > 0 && !m.CheckFlag("sweet_comfort") {
		if utils.Roll(100, 1, 0) <= 50 {
			go Rooms[m.ParentId].FleeMob(m)
		}
	}
	return
}

func (m *Mob) ProcessCommand(cmdStr string, params []string) {
	// StateCommands is a list of stack functions that cause the mob to do something that affects state and causes actions
	log.Println("Processing command " + cmdStr)
	var StateCommands = map[string]func(params []string){
		"attack":   nil,
		"cast":     nil,
		"move":     nil,
		"pickup":   nil,
		"teleport": nil,
		"flee":     m.Flee,
	}
	// Following is no longer driven from here - it is done inline by the move
	// that triggers it, see Mob.FollowChar. Everything else in this map is a
	// placeholder, so check before calling rather than panicking on a nil.
	if command, ok := StateCommands[cmdStr]; ok && command != nil {
		command(params)
	}
}

func (m *Mob) MobScript(inputStr string) {
	input := strings.Split(inputStr, " ")
	switch input[0] {
	case "$TELEPORT":
		m.Teleport(strings.Join(input[1:], " "))
	}

}

// tickPeriod is how long the mob normally waits between ticks.
func (m *Mob) tickPeriod() time.Duration {
	return time.Duration(8-m.TickModifier) * time.Second
}

// scheduleTick makes the mob's next tick fire d from now and records when
// that is. Always use this instead of resetting MobTicker directly.
func (m *Mob) scheduleTick(d time.Duration) {
	m.MobTicker.Reset(d)
	m.NextTick = time.Now().Add(d)
}

// Stun holds the mob's next tick off until amt seconds from now. A stun is
// a floor on the mob's next action, not an addition to it: if the mob was
// already going to wait longer than amt (its normal tick is further out, or
// a longer stun is running) the call is a no-op, and otherwise the next tick
// moves out to exactly now+amt. Stuns can therefore never be chained to hold
// a mob indefinitely; the longest single stun wins.
func (m *Mob) Stun(amt int) {
	if m.Flags["no_stun"] || amt <= 0 {
		return
	}
	// Any stun that lands breaks a scripted wind-up that stuns cancel, even
	// one too short to move the mob's next tick.
	m.windupStunned()
	until := time.Now().Add(time.Duration(amt) * time.Second)
	if !until.After(m.NextTick) {
		return
	}
	m.IsStunned = true
	m.MobStunned = amt
	m.StunnedUntil = until
	m.scheduleTick(time.Until(until))
}

// StunBy is Stun with a combat metric saying who did it. The seconds go in
// total_damage and final_damage stays 0, so damage totals aren't touched. If
// the stun broke a wind-up, the action gets _broke_windup on the end.
func (m *Mob) StunBy(amt int, action string, by *Character) {
	if m.Flags["no_stun"] || amt <= 0 {
		return
	}
	pending := m.pendingWindup()
	m.Stun(amt)
	if pending != nil && m.pendingWindup() == nil {
		action += "_broke_windup"
	}
	data.StoreCombatMetric(action, 0, 0, amt, 0, 0, 0, by.CharId, by.Tier, 1, m.MobId)
}

// Stunned reports whether the mob is currently held by a stun. A stunned mob
// cannot act on its tick, block an exit, or follow a character out of the room.
func (m *Mob) Stunned() bool {
	return m.IsStunned && time.Now().Before(m.StunnedUntil)
}

// FollowLocked reports whether the mob spawned too recently to follow the
// named character out of the room. A character who has attacked the mob
// gets no lockout.
func (m *Mob) FollowLocked(name string) bool {
	return !m.AttackedBy(name) && time.Since(m.SpawnedAt) < config.MobFollowSpawnLockout
}

// Teleport Special handler for handling a mobs cast of a teleport spell
func (m *Mob) Teleport(target string) {
	newRoom := RandomTeleportRoom(m.ParentId)
	targetName := strings.Split(target, " ")
	workingRoom := Rooms[m.ParentId]
	targetChar := workingRoom.Chars.SearchAll(targetName[0])
	if targetChar != nil {
		if targetChar.Resist {
			// For every 5 points of int over the target there's an extra 10% chance to teleport
			diff := (m.Level - targetChar.Tier) * 5
			chance := 10 + diff
			if utils.Roll(100, 1, 0) > chance {
				if _, err := targetChar.Write([]byte(m.Name + " failed to teleport you.\n")); err != nil {
					log.Println("Error writing to player:", err)
				}
				return
			}
		}
		if _, err := targetChar.Write([]byte(m.Name + " teleported you.\n")); err != nil {
			log.Println("Error writing to player:", err)
		}
		newRoom.LockRoom(m.Name+":teleport:enterR-"+strconv.Itoa(m.ParentId)+":"+targetChar.Name, false)
		workingRoom.Chars.Remove(targetChar)
		newRoom.Chars.Add(targetChar)
		targetChar.ParentId = newRoom.RoomId
		if _, err := targetChar.Write([]byte(newRoom.Look(targetChar))); err != nil {
			log.Println("Error writing to player:", err)
		}
		newRoom.UnlockRoom(m.Name+":teleport:enterR-"+strconv.Itoa(m.ParentId)+":"+targetChar.Name, false)
	}
}

// CalculateInventory On copy to a room calculate the inventory
func (m *Mob) CalculateInventory() {
	m.Inventory = nil
	m.Inventory = NewItemInventory()
	if len(m.ItemList) > 0 {
		for k, v := range m.ItemList {
			if utils.Roll(100, 1, 0) <= v {
				// Successful roll!  Add this item to the inventory!
				newItem := Item{}
				if err := copier.CopyWithOption(&newItem, Items[k], copier.Option{DeepCopy: true}); err != nil {
					log.Println("Error copying item:", err)
				}
				m.Inventory.Add(&newItem)
			}
		}
	}
	if m.Summoned {
		return
	}
	for k, v := range QuestMobDrops(m) {
		if _, ok := Items[k]; ok && utils.Roll(100, 1, 0) <= v {
			newItem := Item{}
			if err := copier.CopyWithOption(&newItem, Items[k], copier.Option{DeepCopy: true}); err != nil {
				log.Println("Error copying event drop:", err)
				continue
			}
			m.Inventory.Add(&newItem)
		}
	}
}

func (m *Mob) ReturnState() string {
	stamStatus := text.Green + "healthy" + text.Info
	if m.Stam.Current < (m.Stam.Max - int(.90*float32(m.Stam.Max))) {
		stamStatus = text.BrightRed + "near death" + text.Info
	} else if m.Stam.Current < (m.Stam.Max - int(.75*float32(m.Stam.Max))) {
		stamStatus = text.Red + "badly injured" + text.Info
	} else if m.Stam.Current < (m.Stam.Max - int(.5*float32(m.Stam.Max))) {
		stamStatus = text.DarkYellow + "injured" + text.Info
	} else if m.Stam.Current < (m.Stam.Max - int(.25*float32(m.Stam.Max))) {
		stamStatus = text.DarkGreen + "slightly injured" + text.Info
	}
	return " looks " + stamStatus
}

func (m *Mob) DropInventory() string {
	var drops []string
	var tempStore []*Item
	for _, item := range m.Inventory.Contents {
		tempStore = append(tempStore, item)
	}
	if len(tempStore) > 0 {
		for _, item := range tempStore {
			if item != nil {
				if err := m.Inventory.Remove(item); err == nil {
					if len(Rooms[m.ParentId].Items.Contents) < 15 {
						item.Placement = m.Placement
						Rooms[m.ParentId].Items.Add(item)
						drops = append(drops, item.Name)
					} else {
						Rooms[m.ParentId].MessageAll(item.Name + " fall on top of other items and rolls away.\n" + text.Reset)
					}
				}
			}
		}
	}
	if m.Gold > 0 {
		newGold := Item{}
		if err := copier.CopyWithOption(&newGold, Items[3456], copier.Option{DeepCopy: true}); err != nil {
			log.Println("Error copying item:", err)
		}
		newGold.Name = strconv.Itoa(m.Gold) + " gold marks"
		newGold.Value = m.Gold
		newGold.Placement = m.Placement
		Rooms[m.ParentId].Items.Add(&newGold)
		drops = append(drops, newGold.Name)
	}
	if len(drops) == 0 {
		return "The " + m.Name + " was carrying:\n Nothing"
	} else {
		return "The " + m.Name + " was carrying:\n " + strings.Join(drops, ", ")
	}
}

// AddThreatDamage adds threat for an offensive action and records the
// attacker as having attacked this mob. The seal of courage raises it by
// SealCourageThreatPercent.
func (m *Mob) AddThreatDamage(damage int, attacker *Character) {
	m.MarkAttackedBy(attacker)
	threat := damage
	if damage > 0 && attacker.CheckFlag("seal-courage") {
		threat += config.ThreatPercent(damage, config.SealCourageThreatPercent)
	}
	m.addThreat(threat, attacker)
	if damage > 0 {
		m.RunScript("HIT", attacker)
	}
}

// MarkAttackedBy records that the character has attacked this mob.
func (m *Mob) MarkAttackedBy(attacker *Character) {
	if m.Attackers == nil {
		m.Attackers = make(map[string]bool)
	}
	m.Attackers[attacker.Name] = true
}

// AttackedBy reports whether the named character has attacked this mob.
func (m *Mob) AttackedBy(name string) bool {
	return m.Attackers[name]
}

// EngagedCount is how many characters in the mob's room are fighting it
// alongside c: c itself, c's party, and anyone on the threat table. Each is
// counted once and staff are ignored. Threat entries for characters who have
// left the room do not count. It is never less than 1.
func (m *Mob) EngagedCount(c *Character) int {
	names := map[string]bool{c.Name: true}
	lead := c
	if c.PartyFollow != "" {
		if l := ActiveCharacters.Find(c.PartyFollow); l != nil {
			lead = l
		}
	}
	names[lead.Name] = true
	for _, name := range lead.PartyFollowers {
		names[name] = true
	}
	for name := range m.ThreatTable {
		names[name] = true
	}
	room, ok := Rooms[m.ParentId]
	if !ok || room.Chars == nil {
		return 1
	}
	count := 0
	for _, char := range room.Chars.Contents {
		if char == c || (names[char.Name] && !char.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster)) {
			count++
		}
	}
	if count < 1 {
		return 1
	}
	return count
}

// addThreat adds threat without marking the character as an attacker, for
// heals, rescues and the mob choosing its own target.
func (m *Mob) addThreat(damage int, attacker *Character) {
	attacker.MarkCombat()
	if !attacker.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		m.ThreatTable[attacker.Name] += damage
		if m.CurrentTarget == "" {
			m.CurrentTarget = attacker.Name
		}
	}
}

func (m *Mob) ApplyEffect(effectName string, length string, interval int, magnitude int, effect func(triggers int), effectOff func()) {
	if effectInstance, ok := m.Effects[effectName]; ok {
		effectInstance.Reset(length)
		return
	}
	m.startEffect(effectName, NewEffect(length, interval, magnitude, effect, effectOff))
}

// ReplaceEffect is ApplyEffect for effects whose callbacks carry state, like a
// burn that remembers its caster: a running effect of the same name is dropped,
// without its effectOff, and the new one starts fresh.
func (m *Mob) ReplaceEffect(effectName string, length string, interval int, magnitude int, effect func(triggers int), effectOff func()) {
	delete(m.Effects, effectName)
	m.startEffect(effectName, NewEffect(length, interval, magnitude, effect, effectOff))
}

func (m *Mob) startEffect(effectName string, e *Effect) {
	m.Effects[effectName] = e
	e.RunEffect()
	if e.interval > 0 {
		m.scheduleEffectTick(effectName, e)
	}
}

// scheduleEffectTick runs a ticking effect every interval seconds until it
// expires. It keeps its own timer instead of riding the mob's tick, so a stun
// does not hold it off.
func (m *Mob) scheduleEffectTick(effectName string, e *Effect) {
	time.AfterFunc(time.Duration(e.interval)*time.Second, func() {
		room := m.lockCurrentRoom(effectName + " effect tick")
		if room == nil {
			return
		}
		defer room.UnlockRoom(m.Name+" "+effectName+" effect tick", false)

		// Stop if the mob died, left the world, or the effect was removed or replaced.
		if m.Effects[effectName] != e || m.Stam.Current <= 0 || !room.Mobs.Contains(m) {
			return
		}
		e.RunEffect()
		if m.Effects[effectName] != e || m.Stam.Current <= 0 {
			return
		}
		// Expire unless another full interval fits (half a second of slack for
		// timer drift), so a duration that isn't a multiple of the interval
		// ends without a late tick.
		if e.TimeRemaining() < float64(e.interval)-0.5 {
			m.expireEffect(effectName, e)
			return
		}
		m.scheduleEffectTick(effectName, e)
	})
}

// lockCurrentRoom locks the room the mob is in, re-checking afterwards in case
// it moved while we waited for the lock.
func (m *Mob) lockCurrentRoom(requester string) *Room {
	for {
		room, ok := Rooms[m.ParentId]
		if !ok {
			return nil
		}
		room.LockRoom(m.Name+" "+requester, false)
		if m.ParentId == room.RoomId {
			return room
		}
		room.UnlockRoom(m.Name+" "+requester, false)
	}
}

func (m *Mob) expireEffect(effectName string, e *Effect) {
	if e.effectOff != nil {
		e.effectOff()
	}
	m.RemoveEffect(effectName)
}

func (m *Mob) RemoveEffect(effectName string) {
	delete(m.Effects, effectName)
	m.FlagOff(effectName, "")
}

func (m *Mob) ApplyHook(hook string, hookName string, executions int, length string, interval int, effect func(), effectOff func()) {
	m.Hooks[hook][hookName] = NewHook(executions, length, interval, effect, effectOff)
}

func (m *Mob) RemoveHook(hook string, hookName string) {
	m.Hooks[hook][hookName].effectOff()
	valPresent := false
	for k := range m.Hooks {
		valPresent = false
		for hName := range m.Hooks[k] {
			if hName == hookName {
				valPresent = true
			}
		}
		if valPresent {
			delete(m.Hooks[k], hookName)
		}
	}
}

func (m *Mob) RunHook(hook string) {
	for name, hookInstance := range m.Hooks[hook] {
		// Process Removing the hook
		if hookInstance.TimeRemaining() == 0 {
			m.RemoveHook(hook, name)
			continue
		}
		if hookInstance.interval > 0 {
			//log.Println(hookInstance.LastTriggerInterval())
			if hookInstance.LastTriggerInterval() <= 0 {
				hookInstance.RunHook()
			}
		} else if hookInstance.interval == -1 {
			hookInstance.RunHook()
		}
	}
	return
}

func (m *Mob) GetInt() int {
	return m.Int.Current
}

func (m *Mob) EditToggleFlag(flagName string) bool {
	if val, exists := m.Flags[flagName]; exists {
		m.Flags[flagName] = !val
		return true
	} else {
		return false
	}
}

func (m *Mob) ToggleFlag(flagName string, provider string) {

	if _, exists := m.Flags[flagName]; exists {
		if m.Flags[flagName] == true && utils.StringIn(provider, m.FlagProviders[flagName]) && len(m.FlagProviders[flagName]) > 1 {
			m.FlagProviders[flagName][utils.IndexOf(provider, m.FlagProviders[flagName])] = m.FlagProviders[flagName][len(m.FlagProviders[flagName])-1] // Copy last element to index i.
			m.FlagProviders[flagName][len(m.FlagProviders[flagName])-1] = ""                                                                            // Erase last element (write zero value).
			m.FlagProviders[flagName] = m.FlagProviders[flagName][:len(m.FlagProviders[flagName])-1]                                                    // Truncate slice.
		} else if m.Flags[flagName] == true && !utils.StringIn(provider, m.FlagProviders[flagName]) && len(m.FlagProviders[flagName]) >= 1 {
			m.FlagProviders[flagName] = append(m.FlagProviders[flagName], provider)
		} else if m.Flags[flagName] == true && len(m.FlagProviders[flagName]) == 1 {
			m.Flags[flagName] = false
			m.FlagProviders[flagName] = []string{}
		} else if m.Flags[flagName] == false && provider == "" {
			m.Flags[flagName] = true
		} else if m.Flags[flagName] == true && provider == "" {
			m.Flags[flagName] = false
			m.FlagProviders[flagName] = []string{}
		} else if m.Flags[flagName] == false && provider != "" {
			m.Flags[flagName] = true
			m.FlagProviders[flagName] = []string{provider}
		}
	} else {
		m.Flags[flagName] = true
		m.FlagProviders[flagName] = []string{provider}
	}
}

func (m *Mob) ToggleFlagAndMsg(flagName string, provider string, msg string) {
	m.ToggleFlag(flagName, provider)
	log.Println(m.Name, " informed: ", msg)
}

func (m *Mob) FlagOn(flagName string, provider string) {

	if _, exists := m.Flags[flagName]; exists {
		if m.Flags[flagName] == true && !utils.StringIn(provider, m.FlagProviders[flagName]) && len(m.FlagProviders[flagName]) >= 1 {
			m.FlagProviders[flagName] = append(m.FlagProviders[flagName], provider)
		} else if m.Flags[flagName] == false {
			m.Flags[flagName] = true
			m.FlagProviders[flagName] = []string{provider}
		}
	} else {
		m.Flags[flagName] = true
		m.FlagProviders[flagName] = []string{provider}
	}
}

func (m *Mob) FlagOff(flagName string, provider string) {
	if _, exists := m.Flags[flagName]; exists {
		if m.Flags[flagName] == true && utils.StringIn(provider, m.FlagProviders[flagName]) && len(m.FlagProviders[flagName]) > 1 {
			m.FlagProviders[flagName][utils.IndexOf(provider, m.FlagProviders[flagName])] = m.FlagProviders[flagName][len(m.FlagProviders[flagName])-1] // Copy last element to index i.
			m.FlagProviders[flagName][len(m.FlagProviders[flagName])-1] = ""                                                                            // Erase last element (write zero value).
			m.FlagProviders[flagName] = m.FlagProviders[flagName][:len(m.FlagProviders[flagName])-1]                                                    // Truncate slice.
		} else if m.Flags[flagName] == true && len(m.FlagProviders[flagName]) == 1 {
			m.Flags[flagName] = false
			m.FlagProviders[flagName] = []string{}
		}
	}
}

func (m *Mob) CheckFlag(flagName string) bool {
	if _, ok := m.Flags[flagName]; ok {
		return m.Flags[flagName]
	}
	return false
}

func (m *Mob) ReceiveDamage(damage int) (int, int, int) {
	resist := min(int(((float64(m.Armor)/100)*config.MobArmorReduction)*float64(damage)), damage)
	finalDamage := damage - resist
	if m.CheckFlag("inertial-barrier") {
		finalDamage -= int(math.Ceil(float64(damage) * config.InertialDamageIgnore))
	}
	// Armor of 200 or more resists the whole hit. A negative remainder would
	// add stamina through Meter.Subtract, so a hit never does less than 0.
	finalDamage = max(finalDamage, 0)
	if m.Stam.Current < finalDamage {
		finalDamage = m.Stam.Current
	}
	m.Stam.Subtract(finalDamage)
	m.windupDamaged(finalDamage)
	if finalDamage > m.WimpyValue && m.CheckFlag("flees") {
		m.MobCommands <- "flee"
	}
	return finalDamage, 0, resist
}

func (m *Mob) ReceiveDamageNoArmor(damage int) (int, int) {
	// Resistances over 100% hand in a negative amount; it must not heal.
	finalDamage := max(damage, 0)
	if m.Stam.Current < finalDamage {
		finalDamage = m.Stam.Current
	}
	m.Stam.Subtract(finalDamage)
	m.windupDamaged(finalDamage)
	if finalDamage > m.WimpyValue && m.CheckFlag("flees") {
		m.MobCommands <- "flee"
	}
	return finalDamage, 0
}

func (m *Mob) ReceiveMagicDamage(damage int, element string) (int, int, int) {
	var resisting float64

	switch element {
	case "fire":
		resisting = float64(m.FireResistance) / 100
		if m.CheckFlag("resist-fire") {
			resisting += float64(config.ElementResistPercent) / 100
		}
	case "air":
		resisting = float64(m.AirResistance) / 100
		if m.CheckFlag("resist-air") {
			resisting += float64(config.ElementResistPercent) / 100
		}
	case "earth":
		resisting = float64(m.EarthResistance) / 100
		if m.CheckFlag("resist-earth") {
			resisting += float64(config.ElementResistPercent) / 100
		}
	case "water":
		resisting = float64(m.WaterResistance) / 100
		if m.CheckFlag("resist-water") {
			resisting += float64(config.ElementResistPercent) / 100
		}
		/*if resisting > 0 {
			resisting = (float64(m.Int.Current) / 30) * resisting
		}*/
	}
	if m.CheckFlag("resist-magic") {
		resisting += .10
	}
	resisted := int(math.Ceil(float64(damage) * resisting))
	stamDam, vitDam := m.ReceiveDamageNoArmor(damage - resisted)
	return stamDam, vitDam, resisted
}

func (m *Mob) ReceiveVitalDamage(damage int) (int, int) {
	damageMod, _, resisted := m.ReceiveDamage(damage)
	return damageMod, resisted
}

func (m *Mob) Heal(damage int) (int, int) {
	m.Stam.Add(damage)
	return damage, 0

}

func (m *Mob) HealStam(damage int) {
	m.Stam.Add(damage)
	return
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
	if m.DamageBonusPct != 0 {
		damage += damage * m.DamageBonusPct / 100
	}
	return damage
}

func (m *Mob) CastSpell(spell string) bool {
	log.Println("Casting spell: ", spell)
	return true
}

func (m *Mob) Look() string {
	buildText := "You see a " + m.Name + ", " + config.TextTiers[m.Level] + " level. \n"
	buildText += m.Description
	if m.Flags["hostile"] {
		buildText += text.Bad + "\nIt looks hostile!" + text.Info
	}
	return buildText
}

// ReturnMobInstanceProps Function to return only the modifiable properties
func ReturnMobInstanceProps(mob *Mob) map[string]interface{} {
	serialList := map[string]interface{}{
		"mobId":     mob.MobId,
		"health":    mob.Stam.Current,
		"mana":      mob.Mana.Current,
		"placement": mob.Placement,
		"inventory": mob.Inventory.Jsonify(),
	}
	return serialList
}

func (m *Mob) Save() {
	mobData := make(map[string]interface{})
	mobData["mob_id"] = m.MobId
	mobData["name"] = m.Name
	mobData["description"] = m.Description
	mobData["exp_multiplier"] = m.ExpMultiplier
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
	mobData["night_only"] = utils.Btoi(m.Flags["night_only"])
	mobData["day_only"] = utils.Btoi(m.Flags["day_only"])
	mobData["guard_treasure"] = utils.Btoi(m.Flags["guard_treasure"])
	mobData["take_treasure"] = utils.Btoi(m.Flags["take_treasure"])
	mobData["steals"] = utils.Btoi(m.Flags["steals"])
	mobData["block_exit"] = utils.Btoi(m.Flags["block_exit"])
	mobData["follows"] = utils.Btoi(m.Flags["follows"])
	mobData["no_steal"] = utils.Btoi(m.Flags["no_steal"])
	mobData["detect_invisible"] = utils.Btoi(m.Flags["detect_invisible"])
	mobData["no_stun"] = utils.Btoi(m.Flags["no_stun"])
	mobData["diseases"] = utils.Btoi(m.Flags["diseases"])
	mobData["poisons"] = utils.Btoi(m.Flags["poisons"])
	mobData["spits_acid"] = utils.Btoi(m.Flags["spits_acid"])
	mobData["ranged_attack"] = utils.Btoi(m.Flags["ranged_attack"])
	mobData["flees"] = utils.Btoi(m.Flags["flees"])
	mobData["blinds"] = utils.Btoi(m.Flags["blinds"])
	mobData["no_specials"] = utils.Btoi(m.Flags["no_specials"])
	mobData["no_touch"] = utils.Btoi(m.Flags["no_touch"])
	mobData["no_turn"] = utils.Btoi(m.Flags["no_turn"])
	mobData["no_instakill"] = utils.Btoi(m.Flags["no_instakill"])
	mobData["placement"] = m.Placement
	mobData["immobile"] = utils.Btoi(m.Flags["immobile"])
	mobData["commands"] = m.SerializeCommands()
	mobData["weaken_steps"] = m.WeakenSteps
	mobData["weaken_offered"] = m.WeakenOffered
	data.UpdateMob(mobData)
}

func (m *Mob) Eval() string {
	var descriptions []string

	descriptions = append(descriptions, "You study the "+m.Name+" in your mind's eye....")

	if m.Level > 0 {
		descriptions = append(descriptions, "It is level "+strconv.Itoa(m.Level)+".")
	}

	if m.Stam.Current > 0 {
		descriptions = append(descriptions, "It currently has "+strconv.Itoa(m.Stam.Current)+" hit points remaining.")
	}

	if exp := m.Experience(); exp > 0 {
		descriptions = append(descriptions, "It is worth "+strconv.Itoa(exp)+" experience points.")
	}

	if m.CheckFlag("poisons") {
		descriptions = append(descriptions, "It can poison you.")
	}

	if m.CheckFlag("undead") {
		descriptions = append(descriptions, "It is an undead creature.")
	}
	if m.BreathWeapon != "" {
		descriptions = append(descriptions, "It can breathe "+m.BreathWeapon+" on you.")
	}
	if m.CheckFlag("fast_moving") {
		descriptions = append(descriptions, "It is a fast-moving creature.")
	}
	if m.CheckFlag("block_exit") {
		descriptions = append(descriptions, "It can block your way.")
	}
	if m.CheckFlag("follows") {
		descriptions = append(descriptions, "It will try to follow you.")
	}
	if m.CheckFlag("no_stun") {
		descriptions = append(descriptions, "It is immune to stun.")
	}
	if m.CheckFlag("crippled") {
		descriptions = append(descriptions, "It is crippled.")
	}
	if m.CheckFlag("no_touch") {
		descriptions = append(descriptions, "Its life force cannot be touched.")
	}
	if m.CheckFlag("no_turn") {
		descriptions = append(descriptions, "It is bound too strongly to be turned.")
	}
	if m.CheckFlag("no_instakill") {
		descriptions = append(descriptions, "It cannot be felled in a single blow.")
	}
	if m.CheckFlag("diseases") {
		descriptions = append(descriptions, "It can spread diseases.")
	}
	if m.CheckFlag("spits_acid") {
		descriptions = append(descriptions, "It can spit acid.")
	}
	if m.CheckFlag("blinds") {
		descriptions = append(descriptions, "It has the ability to blind your vision.")
	}
	if len(m.Spells) > 0 {
		descriptions = append(descriptions, "It can cast spells.")
	}
	if m.CheckFlag("no_steal") {
		descriptions = append(descriptions, "It cannot be stolen from.")
	}
	if m.CheckFlag("steals") {
		descriptions = append(descriptions, "It will try to steal your things!")
	}
	if m.CheckFlag("flees") {
		descriptions = append(descriptions, "It could try to flee from combat.")
	}
	if m.CheckFlag("hostile") {
		descriptions = append(descriptions, "It is  hostile towards you.")
	}
	if m.CheckFlag("ranged_attack") {
		descriptions = append(descriptions, "It has a ranged attack.")
	}
	return strings.Join(descriptions, "\n") + "\n"
}
