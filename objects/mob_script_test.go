package objects

import (
	"strings"
	"testing"
	"time"

	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/prompt"
	"github.com/jinzhu/copier"
)

// scriptFighter is a player who can be hit, stunned and moved by mob scripts.
func scriptFighter(name string) *Character {
	return &Character{
		Object:     Object{Name: name},
		Permission: permissions.Anyone | permissions.Player,
		Flags:      map[string]bool{},
		Modifiers:  map[string]int{},
		Timers:     map[string]time.Time{},
		Hooks:      map[string]map[string]*Hook{},
		Equipment:  &Equipment{},
		Stam:       Meter{Max: 100, Current: 100},
		Vit:        Meter{Max: 100, Current: 100},
	}
}

// scriptBoss puts a mob with the given @trigger scripts in a room with the
// given players, targeting the first of them.
func scriptBoss(t *testing.T, roomId int, scripts map[string]string, chars ...*Character) (*Mob, *Room) {
	m := &Mob{
		Object:      Object{Name: "boss", Commands: map[string]prompt.MenuItem{}},
		ParentId:    roomId,
		Flags:       map[string]bool{},
		Effects:     map[string]*Effect{},
		ThreatTable: map[string]int{},
		Stam:        Meter{Max: 1000, Current: 1000},
		NumDice:     1, SidesDice: 1, PlusDice: 99,
	}
	for trigger, script := range scripts {
		m.AddCommands(trigger, script)
	}
	r := &Room{
		RoomId: roomId, Flags: map[string]bool{},
		Chars: &CharInventory{ParentId: roomId, Contents: chars},
		Mobs:  &MobInventory{ParentId: roomId, Contents: []*Mob{m}},
		Items: NewItemInventory(),
	}
	for _, c := range chars {
		c.ParentId = roomId
	}
	if len(chars) > 0 {
		m.CurrentTarget = chars[0].Name
	}
	Rooms[roomId] = r
	t.Cleanup(func() { delete(Rooms, roomId) })
	return m, r
}

func TestValidMobScript(t *testing.T) {
	cases := []struct {
		script, bad string
		ok          bool
	}{
		{"$CHANCE 20 ; $COOLDOWN 30 ; $AOE 4d10+50 fire", "", true},
		{"$msay Fools! ; $enrage 30 50", "", true},
		{"$MSAY hi ; $GIVEGOLD 5", "$GIVEGOLD", false},
		{" ; ", " ; ", false},
		{"$CHANCE 20 ; $ONFAIL $MSAY Bah ; $NEXT ; $STUN second 24", "", true},
		{"$CHANCE 20 ; $ONFAIL $GIVEGOLD 5", "$GIVEGOLD", false},
		{"$CHANCE 20 ; $ONFAIL $NEXT", "$NEXT", false},
	}
	for _, c := range cases {
		if bad, ok := ValidMobScript(c.script); ok != c.ok || bad != c.bad {
			t.Errorf("ValidMobScript(%q) = %q, %v; want %q, %v", c.script, bad, ok, c.bad, c.ok)
		}
	}
}

func TestRollAmount(t *testing.T) {
	if v, ok := rollAmount("120"); !ok || v != 120 {
		t.Errorf("flat amount = %d, %v", v, ok)
	}
	for i := 0; i < 50; i++ {
		if v, ok := rollAmount("4d10+50"); !ok || v < 54 || v > 90 {
			t.Fatalf("4d10+50 rolled %d, %v", v, ok)
		}
	}
	for _, bad := range []string{"", "d10", "4d", "fire", "4d0", "0d6"} {
		if _, ok := rollAmount(bad); ok {
			t.Errorf("rollAmount(%q) accepted", bad)
		}
	}
}

// Damage types go through the matching mitigation; staff are never hit.
func TestMobScriptDamageTargets(t *testing.T) {
	tank, mage := scriptFighter("Tank"), scriptFighter("Mage")
	dm := scriptFighter("Dm")
	dm.Permission = permissions.Anyone | permissions.Dungeonmaster
	m, _ := scriptBoss(t, -91, map[string]string{
		"@COMBAT": "$AOE 30 true",
		"@HP50":   "$DAMAGE target 20 true",
	}, tank, mage, dm)

	if !m.RunScript("COMBAT", nil) {
		t.Error("an $AOE didn't use the mob's turn")
	}
	if tank.Stam.Current != 70 || mage.Stam.Current != 70 {
		t.Errorf("true AoE left tank %d mage %d, want 70 each", tank.Stam.Current, mage.Stam.Current)
	}
	if dm.Stam.Current != 100 {
		t.Error("an AoE hit staff")
	}
	m.RunScript("HP50", nil)
	if tank.Stam.Current != 50 || mage.Stam.Current != 70 {
		t.Errorf("$DAMAGE target hit tank %d mage %d, want 50 and 70", tank.Stam.Current, mage.Stam.Current)
	}
}

func TestMobScriptControlSteps(t *testing.T) {
	tank := scriptFighter("Tank")
	m, _ := scriptBoss(t, -92, map[string]string{
		"@COMBAT": "$COOLDOWN 60 ; $DAMAGE target 10 true",
		"@IDLE":   "$CHANCE 0 ; $DAMAGE target 10 true",
		"@AGGRO":  "$ONCE ; $DAMAGE target 5 true",
		"@SPAWN":  "$MSAY You dare? ; $FLAG no_stun on",
	}, tank)

	if !m.RunScript("COMBAT", nil) || tank.Stam.Current != 90 {
		t.Fatalf("first run didn't act: stamina %d", tank.Stam.Current)
	}
	if m.RunScript("COMBAT", nil) || tank.Stam.Current != 90 {
		t.Error("the script ran again inside its cooldown")
	}
	if m.RunScript("IDLE", nil) || tank.Stam.Current != 90 {
		t.Error("a 0% chance passed")
	}
	m.RunScript("AGGRO", nil)
	m.RunScript("AGGRO", nil)
	if tank.Stam.Current != 85 {
		t.Errorf("$ONCE script ran more or less than once: stamina %d", tank.Stam.Current)
	}
	if m.RunScript("SPAWN", nil) {
		t.Error("talking and flags used the mob's turn")
	}
	if !m.Flags["no_stun"] {
		t.Error("$FLAG didn't set the flag")
	}
}

// 30% hits for 30% of each target's max stamina and spills into vitality;
// 30%v goes straight to vitality and won't take an element.
func TestMobScriptPercentDamage(t *testing.T) {
	tank, mage := scriptFighter("Tank"), scriptFighter("Mage")
	tank.Stam = Meter{Max: 1000, Current: 1000}
	mage.Stam = Meter{Max: 200, Current: 50}
	m, _ := scriptBoss(t, -96, map[string]string{
		"@COMBAT": "$AOE 30% true",
		"@HP50":   "$DAMAGE target 25%v true",
		"@HP25":   "$DAMAGE target 25%v fire ; $ONFAIL $MSAY Bah",
	}, tank, mage)

	m.RunScript("COMBAT", nil)
	if tank.Stam.Current != 700 || tank.Vit.Current != 100 {
		t.Errorf("tank %d/%d, want 300 off stamina only", tank.Stam.Current, tank.Vit.Current)
	}
	if mage.Stam.Current != 0 || mage.Vit.Current != 90 {
		t.Errorf("mage %d/%d, want 60 taking the last 50 stamina and 10 vitality", mage.Stam.Current, mage.Vit.Current)
	}
	m.RunScript("HP50", nil)
	if tank.Stam.Current != 700 || tank.Vit.Current != 75 {
		t.Errorf("vitality hit left tank %d/%d, want 700/75", tank.Stam.Current, tank.Vit.Current)
	}
	if m.RunScript("HP25", nil) || tank.Vit.Current != 75 {
		t.Errorf("fire vitality hit landed: tank vitality %d", tank.Vit.Current)
	}
}

// A failed step skips to the next $NEXT block, a block that only talks lets
// the next one try, and the first block to act ends the script.
func TestMobScriptNext(t *testing.T) {
	tank := scriptFighter("Tank")
	m, _ := scriptBoss(t, -97, map[string]string{
		"@COMBAT": "$CHANCE 0 ; $DAMAGE target 50 true ; $NEXT ; $MSAY Hah ; $NEXT ; " +
			"$COOLDOWN 60 ; $DAMAGE target 10 true ; $NEXT ; $DAMAGE target 1 true",
	}, tank)

	if !m.RunScript("COMBAT", nil) || tank.Stam.Current != 90 {
		t.Fatalf("first run: stamina %d, want only the cooldown block's 10", tank.Stam.Current)
	}
	if !m.RunScript("COMBAT", nil) || tank.Stam.Current != 89 {
		t.Errorf("second run: stamina %d, want the last block's 1 while the cooldown runs", tank.Stam.Current)
	}
}

// $ONFAIL runs only after a failed step, every one in a row runs, and the
// block ends after them.
func TestMobScriptOnFail(t *testing.T) {
	tank := scriptFighter("Tank")
	m, _ := scriptBoss(t, -98, map[string]string{
		"@COMBAT": "$CHANCE 0 ; $ONFAIL $DAMAGE target 5 true ; $ONFAIL $DAMAGE target 3 true ; $DAMAGE target 50 true",
		"@IDLE":   "$CHANCE 100 ; $ONFAIL $DAMAGE target 5 true ; $DAMAGE target 2 true",
		"@AGGRO":  "$CHANCE 0 ; $ONFAIL $MSAY Bah ; $NEXT ; $DAMAGE target 1 true",
	}, tank)

	if !m.RunScript("COMBAT", nil) || tank.Stam.Current != 92 {
		t.Errorf("failed chance: stamina %d, want both $ONFAIL steps (92)", tank.Stam.Current)
	}
	if m.RunScript("IDLE", nil); tank.Stam.Current != 90 {
		t.Errorf("passed chance: stamina %d, want the $ONFAIL skipped (90)", tank.Stam.Current)
	}
	if !m.RunScript("AGGRO", nil) || tank.Stam.Current != 89 {
		t.Errorf("talking $ONFAIL: stamina %d, want the next block to act (89)", tank.Stam.Current)
	}
}

// second is the most threat after the current target, and nobody when the
// target is alone.
func TestMobScriptSecondTarget(t *testing.T) {
	tank, healer, rogue := scriptFighter("Tank"), scriptFighter("Healer"), scriptFighter("Rogue")
	m, _ := scriptBoss(t, -99, map[string]string{
		"@COMBAT": "$STUN second 24 ; $ONFAIL $DAMAGE target 7 true",
	}, tank, healer, rogue)
	m.ThreatTable = map[string]int{"Tank": 500, "Healer": 300, "Rogue": 100}

	m.RunScript("COMBAT", nil)
	if ready, _ := healer.TimerReady("stun"); ready {
		t.Error("the healer, second on threat, wasn't stunned")
	}
	if ready, _ := rogue.TimerReady("stun"); !ready {
		t.Error("the rogue was stunned")
	}
	if ready, _ := tank.TimerReady("stun"); !ready {
		t.Error("the tank was stunned")
	}

	m.ThreatTable = map[string]int{"Tank": 500}
	m.RunScript("COMBAT", nil)
	if tank.Stam.Current != 93 {
		t.Errorf("with no second, stamina %d, want the $ONFAIL hit (93)", tank.Stam.Current)
	}
}

// Each health phase fires once, in order, and a script that acts takes the
// tick; @AGGRO fires when a target is first picked up.
func TestRunTickScripts(t *testing.T) {
	tank := scriptFighter("Tank")
	m, _ := scriptBoss(t, -93, map[string]string{
		"@HP75":   "$DAMAGE target 1 true",
		"@HP25":   "$DAMAGE target 2 true",
		"@HP10":   "$DAMAGE target 16 true",
		"@AGGRO":  "$DAMAGE target 4 true",
		"@COMBAT": "$DAMAGE target 8 true",
	}, tank)
	lost := func() int { return 100 - tank.Stam.Current }

	if !m.runTickScripts() || lost() != 4 {
		t.Fatalf("first tick with a target: lost %d, want the aggro script's 4", lost())
	}
	if m.runTickScripts(); lost() != 12 {
		t.Errorf("second tick: lost %d, want aggro once then combat (12)", lost())
	}
	m.Stam.Current = 200 // 20%: both 75 and 25 are due, and 50 has no script
	if m.runTickScripts(); lost() != 15 {
		t.Errorf("phase tick: lost %d, want both phases and no combat script (15)", lost())
	}
	if m.runTickScripts(); lost() != 23 {
		t.Errorf("tick after phases: lost %d, want only the combat script (23)", lost())
	}
	m.Stam.Current = 50 // 5%: the last phase
	if m.runTickScripts(); lost() != 39 {
		t.Errorf("10%% phase tick: lost %d, want the @HP10 script in place of combat (39)", lost())
	}
}

func TestMobScriptStunHealEnrage(t *testing.T) {
	tank := scriptFighter("Tank")
	m, _ := scriptBoss(t, -94, map[string]string{
		"@COMBAT": "$STUN target 10",
		"@HP50":   "$HEAL 25",
		"@HP25":   "$ENRAGE 30 50",
	}, tank)

	m.RunScript("COMBAT", nil)
	if ready, _ := tank.TimerReady("stun"); ready {
		t.Error("$STUN didn't stun the target")
	}

	if m.RunScript("HP50", nil) {
		t.Error("a heal at full health used the turn")
	}
	m.Stam.Current = 400
	if !m.RunScript("HP50", nil) || m.Stam.Current != 650 {
		t.Errorf("$HEAL 25 from 400/1000 gave %d, want 650", m.Stam.Current)
	}

	if base := m.InflictDamage(); base != 100 {
		t.Fatalf("base damage = %d, want 100", base)
	}
	m.RunScript("HP25", nil)
	if got := m.InflictDamage(); got != 150 {
		t.Errorf("enraged damage = %d, want 150", got)
	}
	m.RunScript("HP25", nil)
	if got := m.InflictDamage(); got != 150 {
		t.Errorf("a second enrage stacked: damage %d", got)
	}
	m.expireEffect("enrage", m.Effects["enrage"])
	if got := m.InflictDamage(); got != 100 {
		t.Errorf("damage after the rage ended = %d, want 100", got)
	}
}

func TestMobScriptHitTrigger(t *testing.T) {
	tank, rogue := scriptFighter("Tank"), scriptFighter("Rogue")
	m, _ := scriptBoss(t, -95, map[string]string{
		"@HIT": "$DAMAGE attacker 7 true ; $MTELEPORT attacker -96",
	}, tank, rogue)
	Rooms[-96] = &Room{RoomId: -96, Flags: map[string]bool{}, Chars: &CharInventory{ParentId: -96}, Mobs: &MobInventory{}}
	defer delete(Rooms, -96)

	m.AddThreatDamage(50, rogue)

	if rogue.Stam.Current != 93 || tank.Stam.Current != 100 {
		t.Errorf("thorns hit rogue %d tank %d, want only the attacker", rogue.Stam.Current, tank.Stam.Current)
	}
	if rogue.ParentId != -95 {
		t.Error("a @HIT script teleported the attacker out from under their own command")
	}
	m.AddThreatDamage(0, rogue)
	if rogue.Stam.Current != 93 {
		t.Error("@HIT fired on a hit that did no damage")
	}
}

func TestMobScriptTeleport(t *testing.T) {
	tank := scriptFighter("Tank")
	m, r := scriptBoss(t, -97, map[string]string{"@COMBAT": "$MTELEPORT target -98"}, tank)
	dest := &Room{RoomId: -98, Flags: map[string]bool{}, Chars: &CharInventory{ParentId: -98, Contents: []*Character{scriptFighter("Other")}}, Mobs: &MobInventory{}}
	Rooms[-98] = dest
	defer delete(Rooms, -98)

	dest.Lock()
	if m.RunScript("COMBAT", nil) || tank.ParentId != -97 {
		t.Error("teleported into a room that was locked by someone else")
	}
	dest.Unlock()

	if !m.RunScript("COMBAT", nil) {
		t.Fatal("the teleport didn't happen")
	}
	if tank.ParentId != -98 || len(r.Chars.Contents) != 0 || len(dest.Chars.Contents) != 2 {
		t.Errorf("after teleport: parent %d, %d left behind, %d at destination", tank.ParentId, len(r.Chars.Contents), len(dest.Chars.Contents))
	}
}

// A spawned copy carries its template's scripts but none of another copy's
// cooldowns or fired phases.
func TestMobCopyGetsScriptsNotState(t *testing.T) {
	tank := scriptFighter("Tank")
	template, _ := scriptBoss(t, -99, map[string]string{"@COMBAT": "$ONCE ; $DAMAGE target 10 true"}, tank)
	template.RunScript("COMBAT", nil)
	if tank.Stam.Current != 90 {
		t.Fatalf("setup: stamina %d", tank.Stam.Current)
	}

	spawned := Mob{}
	if err := copier.CopyWithOption(&spawned, template, copier.Option{DeepCopy: true}); err != nil {
		t.Fatal(err)
	}
	if spawned.scripts != nil {
		t.Fatal("the copy shares or inherited the original's script state")
	}
	if !spawned.RunScript("COMBAT", nil) || tank.Stam.Current != 80 {
		t.Errorf("the copy's $ONCE script didn't run fresh: stamina %d", tank.Stam.Current)
	}
}

// A script can't start another of the mob's scripts from inside itself.
func TestMobScriptNotReentrant(t *testing.T) {
	tank := scriptFighter("Tank")
	m, _ := scriptBoss(t, -100, map[string]string{"@HIT": "$DAMAGE attacker 5 true"}, tank)
	m.scriptState().running = true
	m.AddThreatDamage(10, tank)
	if tank.Stam.Current != 100 {
		t.Error("@HIT ran while another of the mob's scripts was running")
	}
}

// A permanent mob past its wander count still gets its @IDLE script.
func TestPermanentMobIdleScriptStillTicks(t *testing.T) {
	m, r := scriptBoss(t, -101, map[string]string{"@IDLE": "$FLAG woke on"})
	r.Chars.Contents = []*Character{}
	m.Flags["permanent"] = true
	m.NumWander, m.TicksAlive = 5, 50

	m.Tick()

	if !m.Flags["woke"] {
		t.Error("@IDLE didn't run for a long-standing permanent mob")
	}
}

// The wander count doesn't apply to permanent mobs: a hostile one that has
// stood idle far longer still picks up a player and gets on with its turn.
func TestPermanentHostileMobAggrosPastWanderCount(t *testing.T) {
	tank := scriptFighter("Tank")
	m, _ := scriptBoss(t, -102, map[string]string{"@AGGRO": "$DAMAGE target 5 true"}, tank)
	m.CurrentTarget = ""
	m.Flags["permanent"], m.Flags["hostile"] = true, true
	m.NumWander, m.TicksAlive = 5, 500

	m.Tick()

	if m.CurrentTarget != "Tank" {
		t.Fatalf("target = %q, want the player in the room", m.CurrentTarget)
	}
	if tank.Stam.Current != 95 {
		t.Errorf("@AGGRO didn't fire on the pickup: stamina %d", tank.Stam.Current)
	}
}

// The overboard throw: a warning tick, then on the next tick the victim is
// stunned and sent away - unless the mob was stunned in between.
func TestMobScriptWindup(t *testing.T) {
	const deck, sea = -140, -141
	script := map[string]string{
		"@COMBAT": "$MEMOTE grapples %target%! ; $WINDUP loses its grip on %target%! ; $STUN target 8 ; $MTELEPORT target -141",
	}
	setup := func() (*Mob, *Room, *Character, *Character) {
		tank, rogue := scriptFighter("Tank"), scriptFighter("Rogue")
		m, r := scriptBoss(t, deck, script, tank, rogue)
		Rooms[sea] = &Room{RoomId: sea, Flags: map[string]bool{}, Mobs: &MobInventory{},
			Chars: &CharInventory{ParentId: sea, Contents: []*Character{scriptFighter("Swimmer")}}}
		t.Cleanup(func() { delete(Rooms, sea) })
		m.MobTicker = time.NewTicker(time.Hour)
		m.NextTick = time.Now().Add(8 * time.Second)
		return m, r, tank, rogue
	}
	thrown := func(c *Character) bool {
		ready, _ := c.TimerReady("stun")
		return c.ParentId == sea && !ready
	}

	// Nobody stops it: the tank goes over, even though the mob has since
	// turned to the rogue.
	m, r, tank, rogue := setup()
	if !m.runTickScripts() || m.scriptState().windup == nil {
		t.Fatal("the warning tick didn't wind up, or didn't use the turn")
	}
	if tank.ParentId != deck {
		t.Fatal("the tank was thrown on the warning tick")
	}
	m.CurrentTarget = rogue.Name
	if !m.runTickScripts() {
		t.Error("the follow-up tick didn't count as the mob's turn")
	}
	if !thrown(tank) || rogue.ParentId != deck {
		t.Errorf("after the follow-up: tank in %d, rogue in %d; want the tank thrown and stunned", tank.ParentId, rogue.ParentId)
	}

	// A stun in the window - even one too short to delay the tick - breaks it.
	m, _, tank, _ = setup()
	m.runTickScripts()
	m.Stun(1)
	if m.scriptState().windup != nil {
		t.Fatal("a landed stun didn't cancel the wind-up")
	}
	m.runTickScripts()
	if tank.ParentId != deck {
		t.Error("the tank was thrown after the mob had been stunned")
	}

	// A no_stun mob can't be interrupted.
	m, _, tank, _ = setup()
	m.Flags["no_stun"] = true
	m.runTickScripts()
	m.Stun(10)
	if m.scriptState().windup == nil {
		t.Error("a stun that can't land cancelled the wind-up")
	}

	// The victim got away first: nothing happens and the mob acts normally.
	m, r, tank, _ = setup()
	_ = r
	m.runTickScripts()
	r.Chars.Remove(tank)
	if m.runTickScripts() {
		t.Error("a wind-up with no victim still took the turn")
	}
	if ready, _ := tank.TimerReady("stun"); !ready {
		t.Error("a victim who had left was still stunned")
	}
}

func TestMobScriptFillNames(t *testing.T) {
	m := &Mob{Object: Object{Name: "kraken"}, CurrentTarget: "Tank"}
	r := &mobRun{m: m, attacker: scriptFighter("Rogue")}
	got := r.fillNames("$MEMOTE %self% grapples %target% while %ATTACKER% watches")
	if got != "$MEMOTE kraken grapples Tank while Rogue watches" {
		t.Errorf("fillNames = %q", got)
	}
	r.victim = "Mage"
	if got := r.fillNames("%target%"); got != "Mage" {
		t.Errorf("a wind-up's victim wasn't used: %q", got)
	}
}

func TestParseWindup(t *testing.T) {
	w := parseWindup(strings.Fields("stun damage 400 hits 3 rescue escape loses its grip on %target%!"))
	if !w.stun || w.damage != 400 || w.hits != 3 || !w.rescue || !w.escape || w.breakMsg != "loses its grip on %target%!" {
		t.Errorf("parsed %+v", w)
	}
	w = parseWindup(strings.Fields("is interrupted by nothing"))
	if !w.stun || w.damage != 0 || w.breakMsg != "is interrupted by nothing" {
		t.Errorf("default should be stun only: %+v", w)
	}
	w = parseWindup(strings.Fields("none shrugs it off"))
	if w.stun || w.damage != 0 || w.hits != 0 || w.rescue || w.escape || w.breakMsg != "shrugs it off" {
		t.Errorf("none: %+v", w)
	}
	if w = parseWindup(nil); !w.stun || w.breakMsg != "" {
		t.Errorf("no args: %+v", w)
	}
}

// Each cancel condition stops the throw on its own, and none disables them.
func TestMobScriptWindupConditions(t *testing.T) {
	const deck, sea = -150, -151
	run := func(conds string, between func(m *Mob, tank *Character)) (thrown, pending bool) {
		tank := scriptFighter("Tank")
		m, _ := scriptBoss(t, deck, map[string]string{
			"@COMBAT": "$WINDUP " + conds + " ; $STUN target 8 ; $MTELEPORT target -151",
		}, tank)
		Rooms[sea] = &Room{RoomId: sea, Flags: map[string]bool{}, Mobs: &MobInventory{},
			Chars: &CharInventory{ParentId: sea, Contents: []*Character{scriptFighter("Swimmer")}}}
		t.Cleanup(func() { delete(Rooms, sea) })
		m.MobTicker = time.NewTicker(time.Hour)
		m.NextTick = time.Now().Add(8 * time.Second)
		m.runTickScripts()
		between(m, tank)
		pending = m.scriptState().windup != nil
		m.runTickScripts()
		return tank.ParentId == sea, pending
	}
	nothing := func(*Mob, *Character) {}

	if thrown, _ := run("stun", nothing); !thrown {
		t.Error("uncontested wind-up didn't land")
	}
	if thrown, _ := run("damage 300", func(m *Mob, _ *Character) { m.ReceiveDamageNoArmor(150); m.ReceiveDamageNoArmor(150) }); thrown {
		t.Error("300 damage in the window didn't cancel a 'damage 300' wind-up")
	}
	if thrown, _ := run("damage 300", func(m *Mob, _ *Character) { m.ReceiveDamageNoArmor(299) }); !thrown {
		t.Error("299 damage cancelled a 'damage 300' wind-up")
	}
	if thrown, _ := run("hits 2", func(m *Mob, _ *Character) { m.ReceiveDamage(1); m.ReceiveDamage(1) }); thrown {
		t.Error("two hits didn't cancel a 'hits 2' wind-up")
	}
	if thrown, _ := run("damage 300", func(m *Mob, _ *Character) { m.Stun(10) }); !thrown {
		t.Error("a stun cancelled a wind-up that only damage should")
	}
	if thrown, pending := run("none", func(m *Mob, _ *Character) { m.Stun(10); m.ReceiveDamageNoArmor(900) }); !thrown || !pending {
		t.Error("'none' was cancelled")
	}
	if thrown, _ := run("escape", func(m *Mob, tank *Character) { tank.Placement = m.Placement + 2 }); thrown {
		t.Error("the victim moving away didn't cancel an 'escape' wind-up")
	}
	if thrown, _ := run("rescue", func(_ *Mob, tank *Character) {
		Rooms[deck].Chars.Contents = append(Rooms[deck].Chars.Contents, scriptFighter("Paladin"))
		tank.RescuedBy, tank.RescueUntil = "Paladin", time.Now().Add(time.Minute)
	}); thrown {
		t.Error("a rescue didn't cancel a 'rescue' wind-up")
	}
	if thrown, _ := run("stun rescue", func(m *Mob, _ *Character) { m.Stun(10) }); thrown {
		t.Error("with two conditions, the first didn't cancel")
	}
}
