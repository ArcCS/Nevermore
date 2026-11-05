package objects

import (
	"strings"
	"testing"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/text"
)

const lookTestItemId = 999001

// lookTemplate registers a template for the test and returns a copy of it, the
// way a spawned item is.
func lookTemplate(t *testing.T, template *Item) *Item {
	t.Helper()
	template.ItemId = lookTestItemId
	if template.Flags == nil {
		template.Flags = map[string]bool{}
	}
	old, had := Items[lookTestItemId]
	Items[lookTestItemId] = template
	t.Cleanup(func() {
		if had {
			Items[lookTestItemId] = old
		} else {
			delete(Items, lookTestItemId)
		}
	})
	copied := *template
	copied.Flags = map[string]bool{}
	for k, v := range template.Flags {
		copied.Flags[k] = v
	}
	return &copied
}

func TestLookWeaponLines(t *testing.T) {
	// Tier 10, excellent: 1d19+38 averages 48 of real cap 59.
	item := lookTemplate(t, &Item{ItemType: 0, NumDice: 1, SidesDice: 19, PlusDice: 38, MaxUses: 100})

	look := item.Look(&Character{Class: config.BARBARIAN})
	if !strings.Contains(look, "This item appears to be "+text.LightPurple+"excellent quality"+text.Info+".\n") {
		t.Errorf("quality line missing: %q", look)
	}
	if !strings.Contains(look, "It can be used from tier 10.\n") {
		t.Errorf("tier line missing: %q", look)
	}
	if look = item.Look(&Character{Class: config.FIGHTER}); !strings.Contains(look, "It can be used from tier 9.\n") {
		t.Errorf("fighter tier line missing: %q", look)
	}
	if look = item.Look(nil); !strings.Contains(look, "It can be used from tier 10.\n") {
		t.Errorf("nil viewer tier line missing: %q", look)
	}
}

func TestLookIgnoresEnchanting(t *testing.T) {
	// A plain copy of an unremarkable weapon stays that way once enchanted,
	// even though a template with these flags would rate up.
	weapon := lookTemplate(t, &Item{ItemType: 0, NumDice: 1, SidesDice: 19, PlusDice: 36, MaxUses: 100})
	weapon.Flags["magic"] = true
	weapon.Flags["light"] = true
	weapon.Flags["permanent"] = true
	weapon.Adjustment = 30
	if look := weapon.Look(nil); !strings.Contains(look, text.White+"unremarkable") {
		t.Errorf("enchanted weapon changed rating: %q", look)
	}
	Items[lookTestItemId].Flags["no_shatter"] = true
	if look := weapon.Look(nil); !strings.Contains(look, text.LightBlue+"good quality") {
		t.Errorf("template flag gave no bonus: %q", look)
	}

	armor := lookTemplate(t, &Item{ItemType: 20, Armor: 55, MaxUses: 100})
	armor.Armor = 200
	look := armor.Look(nil)
	if !strings.Contains(look, text.LightBlue+"good quality") || !strings.Contains(look, "It can be used from tier 10.\n") {
		t.Errorf("enchanted armor changed rating: %q", look)
	}
}

func TestLookSkipsUnratedWeapons(t *testing.T) {
	for _, template := range []*Item{
		{ItemType: 0, MaxUses: 100},                                         // 0 dice
		{ItemType: 16, NumDice: 1, SidesDice: 19, PlusDice: 38, MaxUses: 1}, // instrument
	} {
		look := lookTemplate(t, template).Look(nil)
		if strings.Contains(look, "appears to be") || strings.Contains(look, "can be used from") {
			t.Errorf("type %d rated: %q", template.ItemType, look)
		}
	}
}

func TestLookPastEveryCapHasNoTier(t *testing.T) {
	item := lookTemplate(t, &Item{ItemType: 0, NumDice: 100, SidesDice: 100, PlusDice: 100, MaxUses: 100})
	look := item.Look(nil)
	if !strings.Contains(look, text.LightOrange+"exceptional quality") || strings.Contains(look, "can be used from") {
		t.Errorf("past-cap weapon: %q", look)
	}
}

func TestLookDeviceLines(t *testing.T) {
	device := lookTemplate(t, &Item{ItemType: 6, Spell: "Vigor", MaxUses: 20})
	device.MaxUses = 100 // melded
	look := device.Look(nil)
	if !strings.Contains(look, "This is a device with healing effects.\n") {
		t.Errorf("effect line missing: %q", look)
	}
	if !strings.Contains(look, "This item is swelling unnaturally.\n") {
		t.Errorf("charge line missing: %q", look)
	}

	device.Flags["infinite"] = true
	if look = device.Look(nil); strings.Contains(look, "swelling") {
		t.Errorf("infinite item shows charges: %q", look)
	}

	wand := lookTemplate(t, &Item{ItemType: 8, Spell: "fireball", MaxUses: 10})
	if look = wand.Look(nil); !strings.Contains(look, "This is a wand with damaging effects.\n") {
		t.Errorf("wand line missing: %q", look)
	}

	unknown := lookTemplate(t, &Item{ItemType: 6, Spell: "remove-blind", MaxUses: 10})
	if look = unknown.Look(nil); !strings.Contains(look, "This is a device.\n") {
		t.Errorf("unknown spell line wrong: %q", look)
	}
}

func TestSpellEffectKind(t *testing.T) {
	cases := map[string]string{
		"vigor":      "healing",
		"HEAL":       "healing",
		"fireball":   "damaging",
		"protection": "supportive",
		"teleport":   "utility",
		"light":      "utility",
		"no-such":    "",
	}
	for name, want := range cases {
		if got := SpellEffectKind(name); got != want {
			t.Errorf("%s: got %q, want %q", name, got, want)
		}
	}
	for name := range Spells {
		if SpellEffectKind(name) == "" {
			t.Errorf("%s has no effect word", name)
		}
	}
}
