package config

import "testing"

func TestPronouns(t *testing.T) {
	cases := []struct {
		gender, sub, pos, desc string
	}{
		{"m", "he", "his", "him"},
		{"f", "she", "her", "her"},
		{"", "it", "its", "it"},
		{"x", "it", "its", "it"},
	}
	for _, c := range cases {
		if got := SubPronoun(c.gender); got != c.sub {
			t.Errorf("SubPronoun(%q) = %q, want %q", c.gender, got, c.sub)
		}
		if got := PosPronoun(c.gender); got != c.pos {
			t.Errorf("PosPronoun(%q) = %q, want %q", c.gender, got, c.pos)
		}
		if got := DescPronoun(c.gender); got != c.desc {
			t.Errorf("DescPronoun(%q) = %q, want %q", c.gender, got, c.desc)
		}
	}
}
