package cmd

import "testing"

func TestDetermineWeaponEnchant(t *testing.T) {
	for tier, want := range map[int]int{1: 2, 5: 3, 10: 4, 15: 5, 20: 6, 24: 6, 25: 7} {
		if got := DetermineWeaponEnchant(tier); got != want {
			t.Errorf("tier %d weapon enchant = %d, want %d", tier, got, want)
		}
	}
}

func TestDetermineArmorEnchant(t *testing.T) {
	for tier, want := range map[int]int{1: 1, 5: 2, 10: 3, 15: 4, 20: 5, 24: 5, 25: 6} {
		if got := DetermineArmorEnchant(tier); got != want {
			t.Errorf("tier %d armor enchant = %d, want %d", tier, got, want)
		}
	}
}
