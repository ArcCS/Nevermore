package objects

import (
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jinzhu/copier"
)

// Mob scripts are the @TRIGGER scripts builders attach to a mob with
// addCommand. They run as the mob, inside its tick or the hit that fired
// them, with the mob's room already locked - so unlike player scripts they
// never go through the command dispatcher, and their verbs live here.
//
// Triggers:
//
//	@SPAWN   once, when this copy first starts thinking
//	@AGGRO   when it goes from no target to a target
//	@COMBAT  every tick it has a target
//	@IDLE    every tick it has none
//	@HP75 @HP50 @HP25 @HP10  once each, the first tick after its stamina drops
//	         below that percentage - boss phases
//	@HIT     whenever a player damages it (attacker is the striker)
//
// (@DEATH also lives on mobs but runs as the killer, through the player
// script engine; see Mob.DeathCheck.)
//
// Control steps stop the rest of the script when they fail: $CHANCE pct,
// $COOLDOWN seconds, $ONCE, $IFSTAGE event stage. Put $COOLDOWN after any
// $CHANCE, or a failed chance roll still starts the cooldown. An action
// that can't happen (no one to hit, say) fails the same way.
//
// $NEXT splits a script into blocks tried in order. A failed step skips to
// the next block instead of ending the script, and so does a block that ran
// without taking an action; the first block that takes one ends the script,
// so the mob still does one thing a turn.
//
// $ONFAIL step runs step only if the step before it failed, then ends the
// block. Several in a row all run.
//
// $WINDUP [conditions] [text] splits a script over two ticks, for a move
// players can see coming and stop. The steps before it run now (the warning)
// and the mob spends its turn winding up; the steps after it run on the
// mob's next tick, against whoever was its target when it wound up, even if
// it has turned to someone else since. If the victim has left the room by
// then, nothing happens. A mob winds up one thing at a time.
//
// The conditions say what cancels it; any one of them does. With none
// given, a stun does:
//
//	stun       any stun that lands on the mob
//	damage N   the mob takes N or more damage in the window
//	hits N     the mob is hit N or more times in the window
//	rescue     the victim is under a rescue when it comes due
//	escape     the victim is no longer at the mob's position when it comes due
//	none       nothing cancels it
//
// Whatever follows the conditions is shown after the mob's name when it is
// cancelled, in place of a stock line - so the text can't begin with one of
// those words.
//
// %target%, %attacker% and %self% in a step are replaced with those names,
// and so are %second%, %last% and %threat2%, %threat3%... (see Targets). A
// step whose threat-list placeholder names nobody fails.
//
// An action step (cast, damage, stun, summon, heal, enrage, teleport) uses
// the mob's turn: on a tick, a script that took an action replaces the
// mob's normal move, breath, spell and attack. Talk, flags and stage steps
// are free, and so is anything a @HIT script does. $MTELEPORT does nothing
// from @HIT.
//
// Targets: target (current target), attacker (the striker on @HIT, else the
// target), threatN (Nth on its threat list: the target is threat1, then
// everyone else with threat on it, most first), second (threat2 - often the
// healer), last (least threat, never the target), random, all. Only players
// the mob can see; staff are never hit.
//
// Damage types: physical (armor), fire/air/earth/water (that resistance),
// true (no mitigation). Amounts are a number or dice: 120, 4d10, 4d10+50.
// 30% is 30% of each target's max stamina, dealt like any other hit; 30%v
// is 30% of their max vitality, dealt straight to vitality, and takes only
// physical (armor) or true.

// mobScriptState is one copy's script bookkeeping. copier leaves unexported
// fields alone, so every copy of a template starts with its own.
type mobScriptState struct {
	running    bool // a script is mid-run; see RunScript
	cooldowns  map[string]time.Time
	once       map[string]bool
	hpFired    map[int]bool
	lastTarget string
	spawned    bool
	windup     *mobWindup
}

// mobWindup is the rest of a script waiting for the mob's next tick.
type mobWindup struct {
	trigger  string
	step     int    // the step to resume at
	victim   string // the mob's target when it wound up
	breakMsg string // shown, after the mob's name, if it is cancelled

	// What cancels it, and how far along the cancelling is.
	stun, rescue, escape   bool
	damage, hits           int
	damageTaken, hitsTaken int
}

// windupConditions are the words $WINDUP recognises before its text; the
// value says whether the word takes a number.
var windupConditions = map[string]bool{"stun": false, "damage": true, "hits": true, "rescue": false, "escape": false, "none": false}

// parseWindup reads $WINDUP's arguments into a wind-up.
func parseWindup(args []string) *mobWindup {
	w := &mobWindup{}
	given := false
	i := 0
	for i < len(args) {
		word := strings.ToLower(args[i])
		takesNumber, ok := windupConditions[word]
		if !ok {
			break
		}
		n := 0
		if takesNumber {
			if i+1 >= len(args) {
				break
			}
			if v, err := strconv.Atoi(args[i+1]); err == nil && v > 0 {
				n = v
			}
			i++
		}
		switch word {
		case "stun":
			w.stun = true
		case "damage":
			w.damage = n
		case "hits":
			w.hits = n
		case "rescue":
			w.rescue = true
		case "escape":
			w.escape = true
		}
		given = true
		i++
	}
	if !given {
		w.stun = true
	}
	w.breakMsg = strings.Join(args[i:], " ")
	return w
}

func (m *Mob) scriptState() *mobScriptState {
	if m.scripts == nil {
		m.scripts = &mobScriptState{
			cooldowns: map[string]time.Time{},
			once:      map[string]bool{},
			hpFired:   map[int]bool{},
		}
	}
	return m.scripts
}

// mobRun is one script running.
type mobRun struct {
	m        *Mob
	room     *Room
	trigger  string
	step     int
	attacker *Character
	victim   string // a resumed wind-up's target; overrides the current one
	usedTurn bool
	stopped  bool
}

type mobVerb struct {
	run      func(r *mobRun, args []string)
	usesTurn bool
	Usage    string
}

// MobVerbs lists every mob script verb; addCommand checks scripts against it.
var MobVerbs map[string]mobVerb

func init() {
	MobVerbs = map[string]mobVerb{
		"$NEXT":       {nil, false, "$NEXT - start the next block; a failed step skips to it, and the first block to act ends the script"},
		"$ONFAIL":     {nil, false, "$ONFAIL step - run step only if the step before it failed, then end the block"},
		"$WINDUP":     {mobWindUp, false, "$WINDUP [stun|damage N|hits N|rescue|escape|none ...] [text] - finish the script next tick unless cancelled; text shows if it is"},
		"$CHANCE":     {mobChance, false, "$CHANCE pct - stop unless a pct% roll succeeds"},
		"$COOLDOWN":   {mobCooldown, false, "$COOLDOWN seconds - stop if this step ran less than seconds ago"},
		"$ONCE":       {mobOnce, false, "$ONCE - stop if this script already got past here on this copy"},
		"$IFSTAGE":    {mobIfStage, false, "$IFSTAGE event stage - stop unless the event is live and at stage"},
		"$EVENTSTAGE": {mobEventStage, false, "$EVENTSTAGE event stage - advance a live event to stage"},
		"$MSAY":       {mobSay, false, "$MSAY text - the mob says text"},
		"$MEMOTE":     {mobEmote, false, "$MEMOTE text - the room sees '<mob> text'"},
		"$EXITTO":     {mobExitTo, false, "$EXITTO room_id exit to_room_id seconds - the exit leads elsewhere for a while (0 = put it back)"},
		"$FLAG":       {mobFlag, false, "$FLAG name on|off - set one of the mob's flags"},
		"$MCAST":      {mobCast, true, "$MCAST spell [target|attacker|second|threatN|last|random|all] - cast without mana"},
		"$DAMAGE":     {mobDamage, true, "$DAMAGE target|attacker|second|threatN|last|random|all amount [type] - hurt players; amount 120, 4d10+50, 30% (of max stamina) or 30%v (of max vitality, straight to it)"},
		"$AOE":        {mobAoe, true, "$AOE amount [type] - hurt everyone the mob can see; amounts as $DAMAGE"},
		"$STUN":       {mobStun, true, "$STUN target|attacker|second|threatN|last|random|all seconds - stun players"},
		"$HEAL":       {mobHeal, true, "$HEAL pct - heal pct% of max stamina"},
		"$SUMMON":     {mobSummon, true, "$SUMMON mob_id [count] - call up to 5 adds (no loot)"},
		"$ENRAGE":     {mobEnrage, true, "$ENRAGE seconds pct - deal pct% more damage for seconds"},
		"$MTELEPORT":  {mobTeleport, true, "$MTELEPORT target|attacker|second|threatN|last|random|all room_id - send players away"},
	}
}

// MobTriggers are the triggers that run as the mob.
var MobTriggers = []string{"@SPAWN", "@AGGRO", "@COMBAT", "@IDLE", "@HP75", "@HP50", "@HP25", "@HP10", "@HIT"}

// ValidMobScript checks every step of a mob script starts with a mob verb
// and returns the first one that doesn't.
func ValidMobScript(script string) (string, bool) {
	steps := splitSteps(script)
	if len(steps) == 0 {
		return script, false
	}
	for _, step := range steps {
		fields := strings.Fields(step)
		verb := strings.ToUpper(fields[0])
		if _, ok := MobVerbs[verb]; !ok {
			return verb, false
		}
		if verb == "$ONFAIL" && len(fields) > 1 {
			body := strings.ToUpper(fields[1])
			if v, ok := MobVerbs[body]; !ok || v.run == nil {
				return body, false
			}
		}
	}
	return "", true
}

func splitSteps(script string) []string {
	var steps []string
	for _, step := range strings.Split(script, ";") {
		if step = strings.TrimSpace(step); step != "" {
			steps = append(steps, step)
		}
	}
	return steps
}

// RunScript runs the mob's @trigger script, if it has one, and reports
// whether it took an action. Call it with the mob's room locked.
func (m *Mob) RunScript(trigger string, attacker *Character) bool {
	return m.runScriptFrom(trigger, attacker, 0, "")
}

// runScriptFrom runs a script from the given step. A victim names the
// target a resumed wind-up is aimed at.
func (m *Mob) runScriptFrom(trigger string, attacker *Character, start int, victim string) bool {
	val, ok := m.Commands["@"+trigger]
	if !ok || m.Stam.Current <= 0 {
		return false
	}
	room, ok := Rooms[m.ParentId]
	if !ok {
		return false
	}
	// One script at a time per mob: an action that ends up damaging the mob
	// (a reflected spell, say) must not start its @HIT script inside the
	// script that caused it.
	st := m.scriptState()
	if st.running {
		return false
	}
	st.running = true
	defer func() { st.running = false }()

	r := &mobRun{m: m, room: room, trigger: trigger, attacker: attacker, victim: victim}
	steps := splitSteps(val.Command)
	skipping := false // a step in this block failed; wait for the next $NEXT
	failed := false   // the step just before failed, so $ONFAIL steps run
	for i := start; i < len(steps); i++ {
		filled, named := r.fillNames(steps[i])
		fields := strings.Fields(filled)
		name := strings.ToUpper(fields[0])
		onFail := name == "$ONFAIL"
		switch {
		case name == "$NEXT":
			// A resumed wind-up finishes its own block and no more.
			if start > 0 || r.usedTurn {
				return r.usedTurn
			}
			skipping, failed = false, false
			continue
		case onFail:
			if !failed || len(fields) < 2 {
				continue
			}
			fields = fields[1:]
			name = strings.ToUpper(fields[0])
		case skipping:
			failed = false
			continue
		}
		verb, ok := MobVerbs[name]
		if !ok || verb.run == nil {
			continue
		}
		r.step = i
		if !named {
			r.stopped = true
		} else if verb.run(r, fields[1:]); name == "$WINDUP" {
			// The rest waits for the next tick; see runTickScripts.
			break
		}
		if r.stopped {
			r.stopped = false
			if !onFail {
				skipping, failed = true, true
			}
			continue
		}
		r.usedTurn = r.usedTurn || verb.usesTurn
	}
	return r.usedTurn
}

// runTickScripts runs the scripts due on this tick: health phases, then
// aggro, then combat or idle. It reports whether any took the mob's turn.
func (m *Mob) runTickScripts() bool {
	st := m.scriptState()
	// A wind-up that nobody interrupted comes due before anything else.
	if w := st.windup; w != nil {
		victim := m.windupVictim(w)
		switch {
		case victim == nil:
			st.windup = nil
		case w.rescue && victim.Rescuer() != nil, w.escape && victim.Placement != m.Placement:
			m.cancelWindup(w)
		default:
			st.windup = nil
			m.runScriptFrom(w.trigger, nil, w.step, w.victim)
			st.lastTarget = m.CurrentTarget
			return true
		}
	}
	used := false
	if m.Stam.Max > 0 {
		pct := m.Stam.Current * 100 / m.Stam.Max
		for _, threshold := range []int{75, 50, 25, 10} {
			if pct < threshold && !st.hpFired[threshold] {
				st.hpFired[threshold] = true
				used = m.RunScript("HP"+strconv.Itoa(threshold), nil) || used
			}
		}
	}
	if st.lastTarget == "" && m.CurrentTarget != "" && !used {
		used = m.RunScript("AGGRO", nil)
	}
	st.lastTarget = m.CurrentTarget
	if used {
		return true
	}
	if m.CurrentTarget != "" {
		return m.RunScript("COMBAT", nil)
	}
	return m.RunScript("IDLE", nil)
}

// runSpawnScript runs @SPAWN the first time this copy starts thinking; a
// permanent mob restarting when its room wakes up again doesn't re-run it.
func (m *Mob) runSpawnScript() {
	st := m.scriptState()
	if st.spawned {
		return
	}
	st.spawned = true
	m.RunScript("SPAWN", nil)
}

// windupVictim is the wound-up move's target, if still in the room.
func (m *Mob) windupVictim(w *mobWindup) *Character {
	room, ok := Rooms[m.ParentId]
	if !ok {
		return nil
	}
	for _, c := range room.Chars.Contents {
		if c.Name == w.victim {
			return c
		}
	}
	return nil
}

// windupStunned is told when a stun lands on the mob.
func (m *Mob) windupStunned() {
	if w := m.pendingWindup(); w != nil && w.stun {
		m.cancelWindup(w)
	}
}

// windupDamaged is told how much every hit took off the mob.
func (m *Mob) windupDamaged(amount int) {
	w := m.pendingWindup()
	if w == nil || amount <= 0 {
		return
	}
	w.damageTaken += amount
	w.hitsTaken++
	if (w.damage > 0 && w.damageTaken >= w.damage) || (w.hits > 0 && w.hitsTaken >= w.hits) {
		m.cancelWindup(w)
	}
}

func (m *Mob) pendingWindup() *mobWindup {
	if m.scripts == nil {
		return nil
	}
	return m.scripts.windup
}

// cancelWindup drops a pending wind-up and tells the room.
func (m *Mob) cancelWindup(w *mobWindup) {
	if m.scripts == nil || m.scripts.windup != w {
		return
	}
	m.scripts.windup = nil
	msg := "is interrupted!"
	if w.breakMsg != "" {
		msg = w.breakMsg
	}
	if room, ok := Rooms[m.ParentId]; ok {
		room.MessageAll(text.Green + m.Name + " " + msg + text.Reset + "\n")
	}
}

// threatName matches the threat-list placeholders: %second%, %last%,
// %threat2% and so on.
var threatName = regexp.MustCompile(`(?i)%(second|last|threat[0-9]+)%`)

// fillNames replaces %target%, %attacker%, %self% and the threat-list
// placeholders in a step. It reports false when a threat-list placeholder
// names nobody, so the step fails rather than talking about no one.
func (r *mobRun) fillNames(step string) (string, bool) {
	if !strings.Contains(step, "%") {
		return step, true
	}
	found := true
	step = threatName.ReplaceAllStringFunc(step, func(match string) string {
		who := r.targets(strings.Trim(match, "%"))
		if len(who) == 0 {
			found = false
			return match
		}
		return who[0].Name
	})
	target := r.m.CurrentTarget
	if r.victim != "" {
		target = r.victim
	}
	attacker := target
	if r.attacker != nil {
		attacker = r.attacker.Name
	}
	return strings.NewReplacer(
		"%target%", target, "%TARGET%", target,
		"%attacker%", attacker, "%ATTACKER%", attacker,
		"%self%", r.m.Name, "%SELF%", r.m.Name,
	).Replace(step), found
}

// --- control ---

func mobWindUp(r *mobRun, args []string) {
	st := r.m.scriptState()
	targets := r.targets("target")
	r.stopped = true
	if st.windup != nil || len(targets) == 0 {
		return
	}
	w := parseWindup(args)
	w.trigger, w.step, w.victim = r.trigger, r.step+1, targets[0].Name
	st.windup = w
	r.usedTurn = true
}

func (r *mobRun) key() string { return r.trigger + "#" + strconv.Itoa(r.step) }

func mobChance(r *mobRun, args []string) {
	if pct, ok := intArg(args, 0); !ok || utils.Roll(100, 1, 0) > pct {
		r.stopped = true
	}
}

func mobCooldown(r *mobRun, args []string) {
	secs, ok := intArg(args, 0)
	st := r.m.scriptState()
	if !ok || time.Now().Before(st.cooldowns[r.key()]) {
		r.stopped = true
		return
	}
	st.cooldowns[r.key()] = time.Now().Add(time.Duration(secs) * time.Second)
}

func mobOnce(r *mobRun, _ []string) {
	st := r.m.scriptState()
	if st.once[r.key()] {
		r.stopped = true
		return
	}
	st.once[r.key()] = true
}

func mobIfStage(r *mobRun, args []string) {
	stage, ok := intArg(args, 1)
	if !ok {
		r.stopped = true
		return
	}
	if current, live := QuestEventStageIn(args[0], r.room); !live || current < stage {
		r.stopped = true
	}
}

func mobEventStage(r *mobRun, args []string) {
	stage, ok := intArg(args, 1)
	if !ok {
		return
	}
	if current, live := QuestEventStageIn(args[0], r.room); live && stage > current {
		_, _ = AdvanceQuestEvent(args[0], stage, nil, true)
	}
}

// --- free actions ---

func mobSay(r *mobRun, args []string) {
	r.room.MessageAll(text.Yellow + r.m.Name + " says, \"" + strings.Join(args, " ") + "\"" + text.Reset + "\n")
}

func mobEmote(r *mobRun, args []string) {
	r.room.MessageAll(text.Yellow + r.m.Name + " " + strings.Join(args, " ") + text.Reset + "\n")
}

func mobFlag(r *mobRun, args []string) {
	if len(args) < 2 {
		return
	}
	name := strings.ToLower(args[0])
	switch strings.ToUpper(args[1]) {
	case "ON":
		r.m.Flags[name] = true
	case "OFF":
		r.m.Flags[name] = false
	}
}

func mobExitTo(r *mobRun, args []string) {
	roomId, exit, to, d, ok := ParseExitTo(args)
	if !ok {
		r.stopped = true
		return
	}
	RedirectExitFor(roomId, exit, to, d)
}

// --- actions ---

// targets resolves a target word to the players it means.
func (r *mobRun) targets(word string) []*Character {
	var visible []*Character
	for _, c := range r.room.Chars.Contents {
		if isStaff(c) || c.Flags["hidden"] || (c.Flags["invisible"] && !r.m.Flags["detect-invisible"]) {
			continue
		}
		visible = append(visible, c)
	}
	pick := func(name string) []*Character {
		for _, c := range visible {
			if c.Name == name {
				return []*Character{c}
			}
		}
		return nil
	}
	word = strings.ToUpper(word)
	if word == "SECOND" {
		word = "THREAT2"
	}
	if n, err := strconv.Atoi(strings.TrimPrefix(word, "THREAT")); err == nil && strings.HasPrefix(word, "THREAT") {
		ranked := r.threatOrder(visible)
		if n < 1 || n > len(ranked) {
			return nil
		}
		return ranked[n-1 : n]
	}
	switch word {
	case "ALL":
		return visible
	case "RANDOM":
		if len(visible) == 0 {
			return nil
		}
		return []*Character{visible[rand.Intn(len(visible))]}
	case "LAST":
		// Never the target itself, so a lone tank isn't "last" too.
		if ranked := r.threatOrder(visible); len(ranked) > 1 {
			return ranked[len(ranked)-1:]
		}
		return nil
	case "ATTACKER":
		if r.attacker != nil {
			if c := pick(r.attacker.Name); c != nil {
				return c
			}
		}
		return pick(r.m.CurrentTarget)
	default: // TARGET
		if r.victim != "" {
			return pick(r.victim)
		}
		return pick(r.m.CurrentTarget)
	}
}

// threatOrder is the mob's threat list among visible: its current target
// first, then everyone else with threat on it, most first.
func (r *mobRun) threatOrder(visible []*Character) []*Character {
	var target *Character
	var rest []*Character
	for _, c := range visible {
		switch {
		case c.Name == r.m.CurrentTarget:
			target = c
		case r.m.ThreatTable[c.Name] > 0:
			rest = append(rest, c)
		}
	}
	sort.SliceStable(rest, func(i, j int) bool {
		return r.m.ThreatTable[rest[i].Name] > r.m.ThreatTable[rest[j].Name]
	})
	if target == nil {
		return rest
	}
	return append([]*Character{target}, rest...)
}

func mobCast(r *mobRun, args []string) {
	if len(args) == 0 {
		r.stopped = true
		return
	}
	spell, ok := Spells[strings.ToLower(args[0])]
	if !ok {
		r.stopped = true
		return
	}
	if utils.StringIn(strings.ToLower(args[0]), MobSupportSpells) {
		r.room.MessageAll(r.m.Name + " casts a " + spell.Name + " spell on itself.\n")
		Cast(r.m, r.m, spell.Effect, spell.Magnitude)
		return
	}
	who := "target"
	if len(args) > 1 {
		who = args[1]
	}
	targets := r.targets(who)
	if len(targets) == 0 {
		r.stopped = true
		return
	}
	for _, t := range targets {
		r.room.MessageAll(r.m.Name + " casts a " + spell.Name + " spell on " + t.Name + ".\n")
		t.RunHook("attacked")
		Cast(r.m, t, spell.Effect, spell.Magnitude)
		t.DeathCheck("was slain by a " + r.m.Name + ".")
	}
}

func mobDamage(r *mobRun, args []string) {
	if len(args) < 2 {
		r.stopped = true
		return
	}
	r.hurt(r.targets(args[0]), args[1:])
}

func mobAoe(r *mobRun, args []string) {
	if len(args) < 1 {
		r.stopped = true
		return
	}
	r.hurt(r.targets("all"), args)
}

// hurt deals amount [type] to each target through the normal damage paths.
func (r *mobRun) hurt(targets []*Character, args []string) {
	if len(targets) == 0 {
		r.stopped = true
		return
	}
	kind := "physical"
	if len(args) > 1 {
		kind = strings.ToLower(args[1])
	}
	for _, t := range targets {
		amount, vital, ok := amountFor(args[0], t)
		// Only armor stands between a blow and vitality; no element does.
		if !ok || (vital && kind != "physical" && kind != "true") {
			r.stopped = true
			return
		}
		t.RunHook("attacked")
		var stam, vit int
		switch {
		case vital && kind == "true":
			vit = t.ReceiveVitalDamageNoArmor(amount)
		case vital:
			vit, _ = t.ReceiveVitalDamage(amount)
		case kind == "true":
			stam, vit = t.ReceiveDamageNoArmor(amount)
		case kind == "fire" || kind == "air" || kind == "earth" || kind == "water":
			stam, vit, _ = t.ReceiveMagicDamage(amount, kind, 0)
		default:
			stam, vit, _ = t.ReceiveDamage(amount)
		}
		writeChar(t, text.Bad+r.m.Name+" hits you for "+strconv.Itoa(stam)+" stamina and "+strconv.Itoa(vit)+" vitality damage."+text.Reset)
		t.DeathCheck("was slain by a " + r.m.Name + ".")
	}
}

func mobStun(r *mobRun, args []string) {
	if len(args) < 2 {
		r.stopped = true
		return
	}
	secs, ok := intArg(args, 1)
	targets := r.targets(args[0])
	if !ok || len(targets) == 0 {
		r.stopped = true
		return
	}
	for _, t := range targets {
		t.RunHook("attacked")
		t.SetTimer("global", secs)
		t.SetTimer("stun", secs)
		writeChar(t, text.Bad+r.m.Name+" stuns you!"+text.Reset)
	}
}

func mobHeal(r *mobRun, args []string) {
	pct, ok := intArg(args, 0)
	if !ok || r.m.Stam.Current >= r.m.Stam.Max {
		r.stopped = true
		return
	}
	r.m.Stam.Current = min(r.m.Stam.Max, r.m.Stam.Current+r.m.Stam.Max*pct/100)
	r.room.MessageAll(text.Green + r.m.Name + " looks revitalized." + text.Reset + "\n")
}

// maxSummon caps one $SUMMON, under the room's own cap of 10 mobs.
const maxSummon = 5

func mobSummon(r *mobRun, args []string) {
	mobId, ok := intArg(args, 0)
	template, exists := Mobs[mobId]
	if !ok || !exists {
		r.stopped = true
		return
	}
	count := 1
	if n, ok := intArg(args, 1); ok {
		count = max(1, min(n, maxSummon))
	}
	summoned := 0
	for i := 0; i < count && len(r.room.Mobs.Contents) < 10; i++ {
		add := Mob{}
		if err := copier.CopyWithOption(&add, template, copier.Option{DeepCopy: true}); err != nil {
			continue
		}
		// Summoned adds carry no loot or gold, so a boss can't be farmed
		// for them, and leave with the summoner's event.
		add.ItemList = map[int]int{}
		add.Gold = 0
		add.QuestEvent = r.m.QuestEvent
		add.Summoned = true
		add.Placement = r.m.Placement
		r.room.Mobs.Add(&add, true)
		add.StartTicking()
		summoned++
	}
	if summoned == 0 {
		r.stopped = true
		return
	}
	r.room.MessageAll(text.Magenta + r.m.Name + " summons " + template.Name + "!" + text.Reset + "\n")
}

func mobEnrage(r *mobRun, args []string) {
	secs, ok1 := intArg(args, 0)
	pct, ok2 := intArg(args, 1)
	if !ok1 || !ok2 || secs <= 0 {
		r.stopped = true
		return
	}
	m := r.m
	if _, raging := m.Effects["enrage"]; !raging {
		m.DamageBonusPct += pct
		r.room.MessageAll(text.Red + m.Name + " flies into a rage!" + text.Reset + "\n")
	}
	m.ApplyEffect("enrage", strconv.Itoa(secs), secs, 0, func(int) {}, func() {
		m.DamageBonusPct -= pct
		if room, ok := Rooms[m.ParentId]; ok {
			room.MessageAll(m.Name + " calms down.\n")
		}
	})
}

func mobTeleport(r *mobRun, args []string) {
	if len(args) < 2 {
		r.stopped = true
		return
	}
	roomId, ok := intArg(args, 1)
	dest, exists := Rooms[roomId]
	targets := r.targets(args[0])
	// A @HIT script runs inside the attacker's own command, which would carry
	// on in the room they had just been taken out of.
	if r.trigger == "HIT" || !ok || !exists || dest == r.room || len(targets) == 0 {
		r.stopped = true
		return
	}
	// The mob's room is already locked. Waiting for a second room here could
	// deadlock against a command holding it and wanting this one, so if it is
	// busy the teleport just doesn't happen this time.
	if !dest.TryLock() {
		r.stopped = true
		return
	}
	defer dest.Unlock()
	for _, t := range targets {
		r.room.Chars.Remove(t)
		dest.Chars.Add(t)
		t.ParentId = dest.RoomId
		writeChar(t, text.Magenta+r.m.Name+" sends you elsewhere!"+text.Reset)
		writeChar(t, dest.Look(t))
		r.room.MessageAll(t.Name + " vanishes!\n")
	}
}

// --- helpers ---

func intArg(args []string, i int) (int, bool) {
	if i >= len(args) {
		return 0, false
	}
	v, err := strconv.Atoi(args[i])
	return v, err == nil
}

// amountFor reads a damage amount for one target: anything rollAmount
// reads, 30% for 30% of their max stamina, or 30%v for 30% of their max
// vitality, which is dealt straight to vitality (vital).
func amountFor(s string, t *Character) (amount int, vital, ok bool) {
	s = strings.ToLower(s)
	if pct, found := strings.CutSuffix(s, "%v"); found {
		n, err := strconv.Atoi(pct)
		return t.Vit.Max * n / 100, true, err == nil && n > 0
	}
	if pct, found := strings.CutSuffix(s, "%"); found {
		n, err := strconv.Atoi(pct)
		return t.Stam.Max * n / 100, false, err == nil && n > 0
	}
	amount, ok = rollAmount(s)
	return amount, false, ok
}

// rollAmount reads 120, 4d10 or 4d10+50.
func rollAmount(s string) (int, bool) {
	s = strings.ToLower(s)
	if v, err := strconv.Atoi(s); err == nil {
		return max(v, 0), true
	}
	num, rest, found := strings.Cut(s, "d")
	if !found {
		return 0, false
	}
	sides, plus := rest, "0"
	if a, b, ok := strings.Cut(rest, "+"); ok {
		sides, plus = a, b
	}
	n, err1 := strconv.Atoi(num)
	sd, err2 := strconv.Atoi(sides)
	p, err3 := strconv.Atoi(plus)
	if err1 != nil || err2 != nil || err3 != nil || n <= 0 || sd <= 0 {
		return 0, false
	}
	return utils.Roll(sd, n, p), true
}

func writeChar(c *Character, msg string) {
	_, _ = c.Write([]byte(msg + "\n"))
}
