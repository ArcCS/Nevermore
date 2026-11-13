package cmd

import (
	"testing"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
)

func TestEditArgLeftOutIsEmpty(t *testing.T) {
	args := []string{"edit", "me"}
	if got := editArg(args, 1); got != "me" {
		t.Errorf("editArg present = %q", got)
	}
	if got := editArg(args, 2); got != "" {
		t.Errorf("editArg left out = %q, want empty", got)
	}
}

func TestEditRequiresValue(t *testing.T) {
	cases := []struct {
		target, property string
		want             bool
	}{
		{"room", "name", true},
		{"room", "description", false},
		{"room", "storeowner", false},
		{"exit", "name", true},
		{"exit", "placement", true},
		{"exit", "description", false},
		{"exit", "key_id", false},
		{"item", "name", true},
		{"item", "type", true},
		{"item", "description", false},
		{"item", "spell", false},
		{"item", "weight", false},
		{"item", "toggle", false},
		{"mob", "name", true},
		{"mob", "level", true},
		{"mob", "stam", true},
		{"mob", "spells", true},
		{"mob", "breathes", false},
		{"mob", "gold", false},
		{"char", "tier", true},
		{"char", "gold", true},
		{"char", "name", true},
		{"char", "description", false},
		{"char", "title", false},
		{"char", "toggle", false},
		{"nonsense", "name", false},
	}
	for _, c := range cases {
		if got := editRequiresValue(c.target, c.property); got != c.want {
			t.Errorf("editRequiresValue(%s, %s) = %v, want %v", c.target, c.property, got, c.want)
		}
	}
}

func TestEditIntProperty(t *testing.T) {
	cases := []struct {
		target, property string
		want             bool
	}{
		{"exit", "key_id", true},
		{"exit", "name", false},
		{"item", "weight", true},
		{"item", "max_uses", true},
		{"item", "spell", false},
		{"item", "charge_type", false},
		{"mob", "fire_resistance", true},
		{"mob", "exp_multiplier", false},
		{"mob", "breathes", false},
		{"char", "tier", true},
		{"char", "stealthexp", true},
		{"char", "title", false},
		{"room", "name", false},
	}
	for _, c := range cases {
		if got := editIntProperty(c.target, c.property); got != c.want {
			t.Errorf("editIntProperty(%s, %s) = %v, want %v", c.target, c.property, got, c.want)
		}
	}
}

func TestApplyItemEditEmptiesLeftOutValues(t *testing.T) {
	item := &objects.Item{Weight: 7, Spell: "fireball", ChargeType: config.ChargeTypeLiquid}
	item.Description = "old"

	if _, ok := applyItemEdit(item, "weight", nil); !ok || item.Weight != 0 {
		t.Errorf("weight left out: ok=%v weight=%d, want cleared", ok, item.Weight)
	}
	if _, ok := applyItemEdit(item, "spell", nil); !ok || item.Spell != "" {
		t.Errorf("spell left out: ok=%v spell=%q, want cleared", ok, item.Spell)
	}
	if _, ok := applyItemEdit(item, "charge_type", nil); !ok || item.ChargeType != config.ChargeTypeDefault {
		t.Errorf("charge_type left out: ok=%v type=%d, want default", ok, item.ChargeType)
	}
	if _, ok := applyItemEdit(item, "description", nil); !ok || item.Description != "" {
		t.Errorf("description left out: ok=%v description=%q, want cleared", ok, item.Description)
	}
}

func TestApplyItemEditRejectsBadValues(t *testing.T) {
	item := &objects.Item{ItemType: 5, Spell: "fireball", ChargeType: config.ChargeTypeLiquid}
	for _, c := range []struct {
		property string
		args     []string
	}{
		{"type", nil},
		{"type", []string{"sword"}},
		{"spell", []string{"not-a-spell"}},
		{"charge_type", []string{"gas"}},
		{"no_such_property", []string{"1"}},
	} {
		if _, ok := applyItemEdit(item, c.property, c.args); ok {
			t.Errorf("applyItemEdit(%s, %v) was accepted", c.property, c.args)
		}
	}
	if item.ItemType != 5 || item.Spell != "fireball" || item.ChargeType != config.ChargeTypeLiquid {
		t.Errorf("rejected edits changed the item: %+v", item)
	}
}

// An edit moves only the named field between a worn copy and its template.
func TestApplyItemEditLeavesTheCopysStateOutOfTheTemplate(t *testing.T) {
	template := &objects.Item{MaxUses: 20, Adjustment: 0, Armor: 100}
	worn := &objects.Item{MaxUses: 3, Adjustment: 5, Armor: 150}

	for _, item := range []*objects.Item{template, worn} {
		if _, ok := applyItemEdit(item, "description", []string{"A", "new", "look."}); !ok {
			t.Fatal("description edit rejected")
		}
	}
	if template.Description != "A new look." || worn.Description != "A new look." {
		t.Errorf("description not on both: %q, %q", template.Description, worn.Description)
	}
	if template.MaxUses != 20 || template.Adjustment != 0 || template.Armor != 100 {
		t.Errorf("template picked up the copy's state: %+v", template)
	}
	if worn.MaxUses != 3 || worn.Adjustment != 5 || worn.Armor != 150 {
		t.Errorf("copy lost its own state: %+v", worn)
	}
}

func TestApplyMobEdit(t *testing.T) {
	mob := &objects.Mob{Gold: 50, BreathWeapon: "fire"}
	mob.Stam = objects.Meter{Max: 100, Current: 40}

	if _, ok := applyMobEdit(mob, "gold", nil); !ok || mob.Gold != 0 {
		t.Errorf("gold left out: ok=%v gold=%d, want cleared", ok, mob.Gold)
	}
	if _, ok := applyMobEdit(mob, "breathes", nil); !ok || mob.BreathWeapon != "" {
		t.Errorf("breathes left out: ok=%v breath=%q, want unset", ok, mob.BreathWeapon)
	}
	if _, ok := applyMobEdit(mob, "breathes", []string{"lava"}); ok {
		t.Error("invalid breath weapon accepted")
	}
	if _, ok := applyMobEdit(mob, "spells", []string{"not-a-spell"}); ok {
		t.Error("unknown spell accepted")
	}
	if _, ok := applyMobEdit(mob, "no_such_property", []string{"1"}); ok {
		t.Error("unknown property accepted")
	}
	if _, ok := applyMobEdit(mob, "stam", []string{"200"}); !ok || mob.Stam.Max != 200 || mob.Stam.Current != 200 {
		t.Errorf("stam edit: ok=%v stam=%+v", ok, mob.Stam)
	}
}
