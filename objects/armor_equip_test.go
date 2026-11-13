package objects

import (
	"testing"

	"github.com/ArcCS/Nevermore/config"
)

func TestEnchantedArmorMayExceedCap(t *testing.T) {
	const id = 987654
	Items[id] = &Item{Object: Object{Name: "vest"}, ItemId: id, ItemType: 5, Armor: 90, ArmorClass: 1, Flags: map[string]bool{}}
	defer delete(Items, id)
	fighter := &Character{Class: config.FIGHTER, Tier: 10}

	// Template at the t10-14 light cap (90), copy enchanted to 95.
	enchanted := &Item{Object: Object{Name: "vest"}, ItemId: id, ItemType: 5, Armor: 95, ArmorClass: 1, Flags: map[string]bool{}}
	if ok, msg := fighter.CanEquip(enchanted); !ok {
		t.Errorf("enchant bonus should be allowed past the cap: %s", msg)
	}

	// A template that is itself over the cap is still refused.
	Items[id].Armor = 95
	if ok, _ := fighter.CanEquip(enchanted); ok {
		t.Error("a piece whose own armor is over the cap should be refused")
	}
}

func TestLevelCapIgnoresGearCaps(t *testing.T) {
	const id = 987655
	Items[id] = &Item{Object: Object{Name: "plate"}, ItemId: id, ItemType: 5, Armor: 500, ArmorClass: 3, Flags: map[string]bool{}}
	defer delete(Items, id)
	plate := &Item{Object: Object{Name: "plate"}, ItemId: id, ItemType: 5, Armor: 500, ArmorClass: 3, Flags: map[string]bool{}}
	sword := &Item{Object: Object{Name: "sword"}, ItemType: 0, SidesDice: 100, NumDice: 2, Flags: map[string]bool{}}

	below := &Character{Class: config.FIGHTER, Tier: config.LevelCap - 1}
	if ok, _ := below.CanEquip(plate); ok {
		t.Error("over-cap armor should be refused below LevelCap")
	}
	if ok, _ := below.CanEquip(sword); ok {
		t.Error("over-cap weapon should be refused below LevelCap")
	}

	capped := &Character{Class: config.FIGHTER, Tier: config.LevelCap}
	if ok, msg := capped.CanEquip(plate); !ok {
		t.Errorf("over-cap armor should be allowed at LevelCap: %s", msg)
	}
	if ok, msg := capped.CanEquip(sword); !ok {
		t.Errorf("over-cap weapon should be allowed at LevelCap: %s", msg)
	}

	// Class rules still hold at LevelCap.
	mage := &Character{Class: config.MAGE, Tier: config.LevelCap}
	if ok, _ := mage.CanEquip(plate); ok {
		t.Error("armor class restriction should still apply at LevelCap")
	}
}

func TestVitalStrikeRespectsInertialBarrier(t *testing.T) {
	c := &Character{Class: config.MAGE, Tier: 20, Flags: map[string]bool{"inertial-barrier": true}, Modifiers: map[string]int{}}
	c.Int.Current = 30
	c.Vit.Current, c.Vit.Max = 1000, 1000
	c.Equipment = &Equipment{}
	// No armor, so only the barrier applies: 1000 - ceil(1000 * .1595)
	if got, _ := c.ReceiveVitalDamage(1000); got != 840 {
		t.Errorf("vital damage through barrier = %d, want 840", got)
	}
}

func TestCheckEquipmentReturnsDisallowedGear(t *testing.T) {
	helm := &Item{Object: Object{Name: "helm"}, ItemType: 25, Flags: map[string]bool{}}
	robe := &Item{Object: Object{Name: "robe"}, ItemType: 5, Flags: map[string]bool{}}
	e := &Equipment{Head: helm, Chest: robe}
	e.CanEquip = func(i *Item) (bool, string) { return i.ItemType != 25, "" }
	var returned []*Item
	e.ReturnToInventory = func(i *Item) { returned = append(returned, i) }
	e.CheckEquipment()
	if e.Head != (*Item)(nil) || e.Chest != robe {
		t.Error("helm should be removed and robe kept")
	}
	if len(returned) != 1 || returned[0] != helm {
		t.Errorf("helm not returned to inventory: %v", returned)
	}
}
