package cmd

import (
	"log"
	"math"
	"strconv"

	"github.com/ArcCS/Nevermore/data"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(kill{},
		"Usage:  kill target # \n\n Try to attack something! Can also use attack, or shorthand k",
		permissions.Player,
		"kill", "k")
}

type kill cmd

func (kill) process(s *state) {
	if len(s.input) < 1 && s.actor.Victim == nil {
		s.msg.Actor.SendBad("Attack what exactly?")
		return
	}

	if !s.requireSight() {
		return
	}

	if !s.requireStamina() {
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	name := ""
	nameNum := 1
	if len(s.words) < 1 && s.actor.Victim != nil {
		switch s.actor.Victim.(type) {
		case *objects.Character:
			name = s.actor.Victim.(*objects.Character).Name
		case *objects.Mob:
			name = s.actor.Victim.(*objects.Mob).Name
			nameNum = s.where.Mobs.GetNumber(s.actor.Victim.(*objects.Mob))
		}
	} else {
		name = s.input[0]
	}

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	if whatMob != nil {
		s.actor.Victim = whatMob
		s.actor.MarkCombat()

		// This is an override for a GM to delete a mob
		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			s.msg.Actor.SendInfo("You smashed ", whatMob.Name, " out of existence.")
			s.actor.Victim = nil
			objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
			whatMob = nil
			return
		}

		s.actor.RunHook("combat")

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == (*objects.Item)(nil) && s.actor.Class != config.MONK {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		if _, err := whatMob.ThreatTable[s.actor.Name]; !err {
			s.msg.Actor.Send(text.White + "You engaged " + whatMob.Name + " #" + strconv.Itoa(s.where.Mobs.GetNumber(whatMob)) + " in combat.")
			s.msg.Observers.Send(text.White + s.actor.Name + " attacks " + whatMob.Name)
			whatMob.AddThreatDamage(0, s.actor)
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		weapon := s.actor.Equipment.Main
		if s.actor.Class == config.MONK {
			weapon = (*objects.Item)(nil)
		}
		if inRange, rangeMsg := weaponInRange(s, weapon, whatMob); !inRange {
			s.msg.Actor.SendBad(rangeMsg)
			return
		}

		performAttack(s, whatMob)
		return

	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}

// weaponInRange reports whether the mob is at a distance the weapon can reach, and
// why not when it isn't.  A nil weapon is an unarmed attack.
func weaponInRange(s *state, weapon *objects.Item, whatMob *objects.Mob) (bool, string) {
	if weapon == (*objects.Item)(nil) {
		if s.actor.Placement != whatMob.Placement {
			return false, "You are too far away to attack."
		}
		return true, ""
	}

	if (weapon.ItemType != 4 && weapon.ItemType != 3) && (s.actor.Placement != whatMob.Placement) {
		return false, "You are too far away to attack."
	} else if weapon.ItemType == 4 && (s.actor.Placement == whatMob.Placement) {
		return false, "You are too close to attack."
	} else if weapon.ItemType == 3 && (s.actor.Placement == whatMob.Placement) {
		return false, "You are too close to attack."
	} else if weapon.ItemType == 3 && (int(math.Abs(float64(s.actor.Placement-whatMob.Placement))) > 1) {
		return false, "You are too far away to attack."
	}

	return true, ""
}

// performAttack runs a standard weapon attack against the mob with whatever is in the
// main hand: the fighter lethal and ranger snipe rolls, the multi attack chain, mob
// reflection, the death check, weapon wear and the attack timer.
//
// The multi attack resolves in two passes.  Every swing rolls to hit first, with the
// first swing at the normal miss chance and follow-ups at a penalised one.  The
// damage multipliers are then handed out to the landed hits in order, so whichever
// swing connects first always takes the full damage slot and any critical or double
// roll rides on it.
func performAttack(s *state, whatMob *objects.Mob) {
	recklessAttack(s)
	// use a list of attacks,  so we can expand this later if other classes get multi style attacks
	attacks := []float64{
		1.0,
	}

	skillLevel := attackSkillLevel(s)
	// Kill is really the fighters realm for specialty.
	if s.actor.Permission.HasAnyFlags(permissions.Fighter) && s.actor.Equipment.Main.ItemType != 4 {
		// mob lethal? A tier 15 expert executes wounded mobs more often.
		execute := config.ExecuteMultiplier(s.actor.Tier, skillLevel, float64(whatMob.Stam.Current)/float64(whatMob.Stam.Max))
		lethal := config.RollLethalScaled(skillLevel, execute)
		if lethal && whatMob.CheckFlag("no_instakill") {
			// Too strong to fell in one blow: the lethal takes half of what is
			// left, and the swing still lands as usual below.
			damage, _ := whatMob.ReceiveDamageNoArmor(whatMob.Stam.Current / 2)
			s.msg.Actor.SendInfo("You landed a lethal blow, but " + whatMob.Name + " endures, torn for " + strconv.Itoa(damage) + " damage.")
			s.msg.Observers.SendInfo(s.actor.Name + " landed a lethal blow on " + whatMob.Name + ", and it reels.")
			data.StoreCombatMetric("lethal_halved", 0, 0, damage, 0, damage, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			whatMob.AddThreatDamage(damage, s.actor)
			s.actor.AdvanceSkillExp(float64(damage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience()))
		} else if lethal {
			// Sure did.  Kill this fool and bail.
			s.msg.Actor.SendInfo("You landed a lethal blow on the " + whatMob.Name)
			s.msg.Observers.SendInfo(s.actor.Name + " landed a lethal blow on " + whatMob.Name)
			if weapMsg := s.actor.Equipment.DamageWeapon("main", 1); weapMsg != "" {
				s.msg.Actor.SendInfo(weapMsg)
			}
			data.StoreCombatMetric("lethal", 0, 0, whatMob.Stam.Current, 0, 0, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			whatMob.AddThreatDamage(whatMob.Stam.Current, s.actor)
			lethalFraction := math.Max(float64(whatMob.Stam.Current)/float64(whatMob.Stam.Max), config.LethalSkillFloor)
			s.actor.AdvanceSkillExp(lethalFraction * float64(whatMob.Experience()))
			whatMob.Stam.Current = 0
			DeathCheck(s, whatMob)
			s.actor.SetTimer("combat", config.CombatCooldown)
			return
		}
	}

	if s.actor.Permission.HasAnyFlags(permissions.Ranger) && s.actor.Equipment.Main.ItemType == 4 && s.actor.Tier > 7 {
		// Sniper
		if utils.Roll(1000, 1, 0) == 1 {
			// Throw a snipe
			whatMob.AddThreatDamage(whatMob.Stam.Max/10, s.actor)
			actualDamage, _, resisted := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * config.CombatModifiers["snipe"])))
			data.StoreCombatMetric("snipe", 0, 0, actualDamage+resisted, resisted, actualDamage, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			s.msg.Actor.SendInfo("You sniped the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
			s.actor.AdvanceSkillExp((float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience())))
			s.msg.Observers.SendInfo(s.actor.Name + " snipes " + whatMob.Name)
			if whatMob.CheckFlag("reflection") {
				reflectDamage := int(float64(actualDamage) * config.ReflectDamageFromMob)
				stamDamage, vitDamage, resisted := s.actor.ReceiveDamage(reflectDamage)
				data.StoreCombatMetric("snipe_mob_reflect", 0, 0, stamDamage+vitDamage+resisted, resisted, stamDamage+vitDamage, 1, whatMob.MobId, whatMob.Level, 0, s.actor.CharId)
				s.msg.Actor.Send("The " + whatMob.Name + " reflects " + strconv.Itoa(reflectDamage) + " damage back at you!")
				s.actor.DeathCheck(" was killed by reflection!")
			}
			DeathCheck(s, whatMob)
			s.actor.SetTimer("combat", config.CombatCooldown)
			return
		}
	}

	flurried := false
	if s.actor.Permission.HasAnyFlags(permissions.Fighter) || (s.actor.Permission.HasAnyFlags(permissions.Ranger) && s.actor.Equipment.Main.ItemType == 4) {
		if mults, ok := config.MultiAttackMultipliers[skillLevel]; ok {
			attacks = mults
		}
	} else if s.actor.Class == config.MONK {
		attacks, flurried = monkAttackPlan(s, skillLevel)
	}

	result := resolveHits(s, whatMob, attacks, skillLevel, "kill")
	if s.actor.Class == config.MONK && result.hits > 0 && !flurried {
		gainChiFromHits(s, result.chiHits())
	}
	DeathCheck(s, whatMob)
	if s.actor.Class != config.MONK && s.actor.Equipment.Main != nil {
		weapMsg := s.actor.Equipment.DamageWeapon("main", result.weaponDamage)
		if weapMsg != "" {
			s.msg.Actor.SendInfo(weapMsg)
		}
	}

	s.actor.SetTimer("combat", config.CombatCooldown)
}

// canCrush reports whether the actor's regular attacks can land a crushing
// blow: barbarians wielding a blunt or two-handed melee weapon.
func canCrush(s *state) bool {
	main := s.actor.Equipment.Main
	if s.actor.Class != config.BARBARIAN || main == nil || main.ItemType == 4 {
		return false
	}
	return main.ItemType == 2 || main.IsTwoHanded()
}

// attackSkillLevel is the weapon skill level behind the actor's attacks:
// hand-to-hand for monks, otherwise the wielded weapon's skill.
func attackSkillLevel(s *state) int {
	if s.actor.Class == config.MONK {
		return config.WeaponLevel(s.actor.Skills[5].Value, s.actor.Class, 5)
	}
	return config.WeaponLevel(s.actor.Skills[s.actor.Equipment.Main.ItemType].Value, s.actor.Class, s.actor.Equipment.Main.ItemType)
}

// monkAttackPlan returns the swing multipliers for a monk's attack round and
// whether the flurry stance paid for them. Flurry costs chi every round; when
// the monk can't cover it the stance drops and the round is a single swing.
func monkAttackPlan(s *state, skillLevel int) ([]float64, bool) {
	single := []float64{1.0}
	if !s.actor.CheckFlag("flurry") {
		return single, false
	}
	cost := config.FlurryChiCost(skillLevel)
	if s.actor.Mana.Current < cost {
		s.actor.RemoveEffect("flurry")
		s.msg.Actor.SendBad("You are too drained to keep up the flurry.")
		return single, false
	}
	s.actor.Mana.Subtract(cost)
	return config.MonkFlurryFor(skillLevel), true
}

// gainChiFromHits awards a monk chi for landed hits and tells them about it.
func gainChiFromHits(s *state, hits int) {
	gained := s.actor.GainChi(config.ChiPerHitFor(s.actor.GetStat("pie"))*hits, false)
	if gained > 0 {
		s.msg.Actor.Send(text.Cyan + "Your chi rises by " + strconv.Itoa(gained) + "." + text.Reset)
	}
}

// hitResult is what resolveHits reports back: how many swings landed, the
// damage dealt, how much wear the weapon took (10 on a critical), and whether
// the round landed a critical strike.
type hitResult struct {
	hits         int
	totalDamage  int
	weaponDamage int
	critical     bool
}

// chiHits is the number of hits' worth of chi the round builds: a critical
// strike counts config.CritChiMultiplier times.
func (r hitResult) chiHits() int {
	if r.critical {
		return r.hits + config.CritChiMultiplier - 1
	}
	return r.hits
}

// resolveHits runs the swings of one attack round against the mob and reports
// the outcome. metric names the combat metric family ("kill", "leap", ...).
//
// It resolves in two passes. Every swing rolls to hit first, with the first
// swing at the normal miss chance and follow-ups (multi attack or flurry) at
// FollowUpMissChance. The damage
// multipliers are then handed out to the landed hits in order, so whichever
// swing connects first always takes the full damage slot and any critical or
// double roll rides on it. Damage and reflection are totalled and reported once.
func resolveHits(s *state, whatMob *objects.Mob, attacks []float64, skillLevel int, metric string) hitResult {
	result := hitResult{weaponDamage: 1}
	lvlDiff, engaged := whatMob.Level-s.actor.Tier, whatMob.EngagedCount(s.actor)
	baseMiss := DetermineMissChance(s, lvlDiff, engaged)
	followUpMiss := FollowUpMissChance(s, lvlDiff, engaged, skillLevel)
	for swing := range attacks {
		missChance := baseMiss
		if swing > 0 {
			missChance = followUpMiss
		}
		if utils.Roll(100, 1, 0) <= missChance {
			data.StoreCombatMetric(metric+"-miss", 0, 0, 0, 0, 0, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
			whatMob.AddThreatDamage(1, s.actor)
			continue
		}
		result.hits++
	}
	if result.hits == 0 {
		s.msg.Actor.SendBad("You missed!!")
		return result
	}

	alwaysCrit := false
	if s.actor.Class != config.MONK {
		alwaysCrit = s.actor.Equipment.Main.Flags["always_crit"]
	}
	crushing := false
	shatterRisk := false
	totalReflect := 0
	for hit := 0; hit < result.hits; hit++ {
		mult := attacks[hit]
		action := metric
		if hit == 0 {
			if canCrush(s) && config.RollCrushing(skillLevel) {
				crushing = true
				shatterRisk = true
				mult *= config.CombatModifiers["crushing"]
				s.msg.Actor.SendGood("Craaackk!! A crushing blow!")
				action = metric + "-crushing"
			} else if config.RollCritical(skillLevel) || alwaysCrit {
				mult *= config.CombatModifiers["critical"]
				s.msg.Actor.SendGood("Critical Strike!")
				result.weaponDamage = 10
				result.critical = true
				shatterRisk = !alwaysCrit
				action = metric + "-critical"
			} else if config.RollDouble(skillLevel) {
				mult *= config.CombatModifiers["double"]
				s.msg.Actor.SendGood("Double Damage!")
				action = metric + "-double"
			}
		}

		actualDamage, _, resisted := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * mult)))
		data.StoreCombatMetric(action, 0, 0, actualDamage+resisted, resisted, actualDamage, 0, s.actor.CharId, s.actor.Tier, 1, whatMob.MobId)
		whatMob.AddThreatDamage(actualDamage, s.actor)
		s.actor.AdvanceSkillExp((float64(actualDamage) / float64(whatMob.Stam.Max) * float64(whatMob.Experience())))
		result.totalDamage += actualDamage
		if whatMob.CheckFlag("reflection") {
			reflectDamage := int(float64(actualDamage) * config.ReflectDamageFromMob)
			stamDamage, vitDamage, resisted := s.actor.ReceiveDamage(reflectDamage)
			data.StoreCombatMetric(metric+"_mob_reflect", 0, 0, stamDamage+vitDamage+resisted, resisted, stamDamage+vitDamage, 1, whatMob.MobId, whatMob.Level, 0, s.actor.CharId)
			totalReflect += reflectDamage
		}
	}
	if result.hits == 1 {
		s.msg.Actor.SendInfo("You hit the " + whatMob.Name + " for " + strconv.Itoa(result.totalDamage) + " damage!" + text.Reset)
	} else {
		s.msg.Actor.SendInfo("You hit the " + whatMob.Name + " " + strconv.Itoa(result.hits) + " times for " + strconv.Itoa(result.totalDamage) + " damage!" + text.Reset)
	}
	if crushing {
		whatMob.StunBy(config.CrushingStuns, "crushing-stun", s.actor)
		s.msg.Observers.SendInfo(s.actor.Name + " lands a crushing blow on " + whatMob.Name + "!")
	}
	if totalReflect > 0 {
		s.msg.Actor.Send("The " + whatMob.Name + " reflects " + strconv.Itoa(totalReflect) + " damage back at you!")
		s.actor.DeathCheck(" was killed by reflection!")
	}
	if shatterRisk {
		shatterWeapon(s)
	}
	warnTouchVulnerable(s, whatMob)
	return result
}

// shatterWeapon rolls the chance that a critical strike, crushing blow or
// Thunk destroyed the weapon that landed it. A no_shatter weapon never rolls,
// nor does an always_crit one, or its own crits would grind it away.
func shatterWeapon(s *state) {
	main := s.actor.Equipment.Main
	if s.actor.Class == config.MONK || main == nil || main.Flags["no_shatter"] || main.Flags["always_crit"] || !config.RollShatter() {
		return
	}
	s.msg.Actor.SendBad(s.actor.Equipment.ShatterWeapon())
	s.msg.Observers.SendInfo(s.actor.Name + "'s weapon shatters!")
}

// DeathCheck Universal death check for mobs on whatever the current state is
func DeathCheck(s *state, m *objects.Mob) {
	if m.Stam.Current <= 0 {
		m.RunDeathScript(s.actor)
		s.msg.Actor.SendGood("You killed " + m.Name)
		s.msg.Observers.SendGood(s.actor.Name + " killed " + m.Name)
		bloodlust(s, m.Name)
		partyLead := s.actor
		if s.actor.PartyFollow != "" {
			partyLead = objects.ActiveCharacters.Find(s.actor.PartyFollow)
		}
		partyMembers := append(partyLead.PartyFollowers, partyLead.Name)
		highestTier := 0
		expReduce := len(s.where.Chars.Contents)
		//reusing this because we are already running through everyone in the room
		for _, gm := range s.where.Chars.Contents {
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
		//debuging stuff
		//s.msg.Actor.SendGood("Highest Tier: " + strconv.Itoa(highestTier))
		//s.msg.Actor.SendGood(strconv.Itoa(tierLimit))
		experienceAwarded := m.KillExperience(s.where, expReduce)

		for _, member := range s.where.Chars.Contents {
			buildActorString := ""
			charClean := s.where.Chars.SearchAll(member.Name)
			if charClean != nil {
				partyCheck := false
				if !s.where.InQuestMode() {
					for _, name := range partyMembers {
						if charClean.Name == name {
							partyCheck = true
						}
					}
				}
				if s.where.InQuestMode() {
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
				if charClean == s.actor {
					buildActorString += text.Green + m.DropInventory() + "\n"
					s.msg.Actor.Send(buildActorString)
				} else {
					go func() {
						if _, err := charClean.Write([]byte(buildActorString + "\n" + text.Reset)); err != nil {
							log.Println("Error writing to player: ", err)
						}
					}()
				}
				if charClean.Victim == m {
					charClean.Victim = nil
				}
			}
		}

		s.where.Mobs.Remove(m)
	}
}

// DetermineMissChance is the miss chance of a weapon attack: the weapon skill
// sets the base and the shared modifiers do the rest.
func DetermineMissChance(s *state, lvlDiff int, engaged int) int {
	return missChance(s, weaponMiss(s), lvlDiff, engaged)
}

// FollowUpMissChance is the miss chance of every swing after the first in a
// multi attack or flurry: the weapon roll plus the follow-up penalty, added
// before the clamp so dex works it off, floored at the skill level's
// follow-up floor instead of 5. A monk's flurry penalty never exceeds
// config.FlurryMaxMissPenalty.
func FollowUpMissChance(s *state, lvlDiff int, engaged int, skillLevel int) int {
	penalty := config.MultiAttackMissPenaltyFor(skillLevel)
	if s.actor.Class == config.MONK {
		penalty = min(penalty, config.FlurryMaxMissPenalty)
	}
	return missChanceFloor(s, weaponMiss(s)+penalty, lvlDiff, engaged, config.FollowUpMissFloorFor(skillLevel))
}

// weaponMiss is the skill-derived base miss of a regular swing: the
// hand-to-hand skill for monks, the wielded weapon's skill for everyone else.
func weaponMiss(s *state) int {
	skill := config.HandSkill
	if s.actor.Class != config.MONK {
		skill = s.actor.Equipment.Main.ItemType
	}
	return config.WeaponMissChance(s.actor.Skills[skill].Value)
}

// BackstabMissChance is the miss chance of a backstab: the regular weapon
// roll plus a penalty that the stealth skill works off.
func BackstabMissChance(s *state, lvlDiff int, engaged int) int {
	stealth := config.StealthLevel(s.actor.Skills[11].Value)
	weapon := config.WeaponMissChance(s.actor.Skills[s.actor.Equipment.Main.ItemType].Value)
	return missChance(s, weapon+config.BackstabMissPenaltyFor(stealth), lvlDiff, engaged)
}

// TodMissChance is the miss chance of a touch of death: the regular
// hand-to-hand roll plus a penalty that tier and piety work off.
func TodMissChance(s *state, lvlDiff int, engaged int) int {
	weapon := config.WeaponMissChance(s.actor.Skills[config.HandSkill].Value)
	return missChance(s, weapon+config.TodMissPenaltyFor(s.actor.Tier, s.actor.GetStat("pie")), lvlDiff, engaged)
}

// missChance builds a to-hit roll's miss chance from a skill-derived base.
// Only the base differs between kinds of attack; the level penalty (which
// loosens with the number of characters engaged, off in quest mode), dex,
// combat flags such as bless and reckless, and the 5-95 clamp are shared.
func missChance(s *state, baseMiss int, lvlDiff int, engaged int) int {
	return missChanceFloor(s, baseMiss, lvlDiff, engaged, 5)
}

// missChanceFloor is missChance with a different lower clamp.
func missChanceFloor(s *state, baseMiss int, lvlDiff int, engaged int, floor int) int {
	missChance := baseMiss
	if !s.where.InQuestMode() {
		missChance += config.OutlevelMissPenalty(lvlDiff, engaged)
	}
	missChance -= s.actor.GetStat("dex") * config.HitPerDex
	if s.actor.CheckFlag("reckless") {
		missChance -= config.RecklessMissReduction
	}
	// The seal of justice carries the to-hit half of bless for paladins, who
	// cannot cast bless themselves. The two do not stack.
	if s.actor.CheckFlag("bless") || s.actor.CheckFlag("seal-justice") {
		missChance -= config.BlessHitBonus
	}
	if missChance >= 100 {
		missChance = 95
	}
	if missChance < floor {
		missChance = floor
	}
	return missChance
}
