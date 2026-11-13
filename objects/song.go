package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
)

// Songs are keyed by the name a bard sings. "effect" names the SongEffects
// entry that does the work, and "tune" is how others hear it ("Lala (singing
// a joyous tune)"), so it hints at what the song does.
var Songs = map[string]map[string]string{
	"celebration-night": {
		"desc":   "Hastens stamina and vitality regeneration on players in the area.",
		"verse":  " let your spirits soar on the wings of eagles, let music wash the sleep from your eyes.",
		"effect": "celebration-night",
		"tune":   "a joyous tune",
	},
	"curious-canticle": {
		"desc":   "prevents creatures from following anyone in the area.",
		"verse":  " now curious, get it in your head, if you follow me, then you'll be dead.",
		"effect": "curious-canticle",
		"tune":   "a bewildering canticle",
	},
	"sweet-comfort": {
		"desc":   "comforts creatures, lulling them so they will not flee.",
		"verse":  " relax, my child, sleep and dream, for things aren't as bad as they seem.",
		"effect": "sweet-comfort",
		"tune":   "a soothing lullaby",
	},
	"run-run-away": {
		"desc":   "lulls creatures so that they will not block escape.",
		"verse":  "see chameleon, lying there in the sun, all things to everyone, run, run away.",
		"effect": "run-away",
		"tune":   "a quick-stepping waltz",
	},
	"draens-tale": {
		"desc":   "increases the the mana regeneration of everyone present in the area.",
		"verse":  "draen, throw your spells anew, know the weave will answer you.",
		"effect": "draens-tale",
		"tune":   "a mystical ballad",
	},
	/*
		"victors-chorus": {
			"desc":   "decreases the damage output of all creatures in the area.",
			"verse":  "we fight to victory, we fight to win, we fight for glory our enemies will never see",
			"effect": "victors-chorus",
		},
			"warriors-threnody": {
				"desc": "instead of death, will restore another character completely (max one per 2 minutes)",
				"verse": "",
				"effect": "warriors-threnody",
			},
			}*/
	"champions-anthem": {
		"desc":   "increases the melee damage output of all players in the area.",
		"verse":  "feel thy strength swell, with every swing, feel the power of my hymn",
		"effect": "champions-anthem",
		"tune":   "a rousing anthem",
	},
	"banshees-lament": {
		"desc":   "causes damage periodically to all creatures in the area.",
		"verse":  "hear my song, hear my sorrow, feel my will tear you apart",
		"effect": "banshees-lament",
		"tune":   "a shrill dirge",
	},
}

// songPulse carries what is decided once per pulse rather than per target.
type songPulse struct {
	// roll is one instrument roll, adjustment included.
	roll int
}

// A song effect applies one pulse to one target and returns the amount it
// applied (mana, healing, damage or bonus), which the pulse logs.
type song struct {
	target string
	effect func(target interface{}, singer *Character, pulse songPulse) int
}

var SongEffects = map[string]song{
	"draens-tale":       {"players", DraensTale},
	"run-away":          {"mobs", RunAway},
	"sweet-comfort":     {"mobs", SweetComfort},
	"curious-canticle":  {"mobs", CuriousCanticle},
	"celebration-night": {"players", CelebrationNight},
	"champions-anthem":  {"players", ChampionsAnthem},
	"banshees-lament":   {"mobs", BansheesLament},
}

// SongEffectFor returns the effect behind a sung song name, looked up through
// the song's "effect" key: the two names differ for run-run-away.
func SongEffectFor(songName string) (string, song, bool) {
	effectName := Songs[songName]["effect"]
	effect, ok := SongEffects[effectName]
	return effectName, effect, ok
}

// instrumentRoll is one roll of the singer's instrument plus its adjustment,
// or 0 when they hold none.
func instrumentRoll(singer *Character) int {
	off := singer.Equipment.Off
	if off == (*Item)(nil) {
		return 0
	}
	return utils.Roll(off.SidesDice, off.NumDice, off.PlusDice) + off.Adjustment
}

// songPulse applies one pulse of the character's current song to the room
// they are in, and logs it. The caller holds the room lock.
func (c *Character) songPulse() {
	effectName, effect, ok := SongEffectFor(c.CurrentSong)
	if !ok {
		return
	}
	room, ok := Rooms[c.ParentId]
	if !ok {
		return
	}
	pulse := songPulse{roll: instrumentRoll(c)}
	total, targets, victimType := 0, 0, 0
	switch effect.target {
	case "mobs":
		victimType = 1
		for _, mob := range room.Mobs.Contents {
			if mob.CheckFlag("hostile") {
				total += effect.effect(mob, c, pulse)
				targets++
			}
		}
	case "players":
		for _, player := range room.Chars.Contents {
			total += effect.effect(player, c, pulse)
			targets++
		}
	}
	// Banshee's lament already logs each mob it hurts.
	if effectName != "banshees-lament" {
		data.StoreCombatMetric("song_"+effectName, 0, 0, total, 0, 0, 0, c.CharId, c.Tier, victimType, targets)
	}
}

// SongDescriptor is how others hear the character's current song, such as
// "a joyous tune".
func (c *Character) SongDescriptor() string {
	if tune := Songs[c.CurrentSong]["tune"]; tune != "" {
		return tune
	}
	return "a song"
}

func DraensTale(target interface{}, singer *Character, pulse songPulse) int {
	switch target := target.(type) {
	case *Character:
		before := target.Mana.Current
		target.RestoreMana(config.DraensTaleMana(pulse.roll, singer.GetStat("pie")))
		return target.Mana.Current - before
	}
	return 0
}

func RunAway(target interface{}, singer *Character, _ songPulse) int {
	switch target := target.(type) {
	case *Mob:
		target.ApplyEffect(singer.Name+"_run_away_song", "16", 0, 0,
			func(triggers int) {
				target.FlagOn("run_away", singer.Name+"_run_away")
			},
			func() {
				target.FlagOff("run_away", singer.Name+"_run_away")
			})
	}
	return 0
}

func SweetComfort(target interface{}, singer *Character, _ songPulse) int {
	switch target := target.(type) {
	case *Mob:
		target.ApplyEffect(singer.Name+"_sweet_comfort_song", "16", 0, 0,
			func(triggers int) {
				target.FlagOn("sweet_comfort", singer.Name+"_sweet_comfort")
			},
			func() {
				target.FlagOff("sweet_comfort", singer.Name+"_sweet_comfort")
			})
	}
	return 0
}

func CuriousCanticle(target interface{}, singer *Character, _ songPulse) int {
	switch target := target.(type) {
	case *Mob:
		target.ApplyEffect(singer.Name+"_curious_canticle_song", "16", 0, 0,
			func(triggers int) {
				target.FlagOn("curious_canticle", singer.Name+"_curious_canticle")
			},
			func() {
				target.FlagOff("curious_canticle", singer.Name+"_curious_canticle")
			})
	}
	return 0
}

// CelebrationNight rolls the instrument for each player, so their heals vary.
func CelebrationNight(target interface{}, singer *Character, _ songPulse) int {
	switch target := target.(type) {
	case *Character:
		heal := config.CelebrationHeal(instrumentRoll(singer), len(Rooms[target.ParentId].Chars.Contents))
		if _, err := singer.Write([]byte(text.Red + "Your song healed " + strconv.Itoa(heal) + " damage to " + target.Name + ".\n" + text.Reset)); err != nil {
			log.Println("Error writing to player:", err)
		}
		before := target.Stam.Current + target.Vit.Current
		target.Heal(heal)
		return target.Stam.Current + target.Vit.Current - before
	}
	return 0
}

func ChampionsAnthem(target interface{}, singer *Character, _ songPulse) int {
	switch target := target.(type) {
	case *Character:
		// The bonus is computed once when the anthem lands and the same amount
		// is removed when it ends, so a piety change on the singer mid-song
		// (pray starting or expiring) cannot leave a residue on the modifier.
		bonus := 0
		target.ApplyEffect(singer.Name+"_champions_anthem_song", "16", 0, 0,
			func(triggers int) {
				bonus = singer.GetStat("pie") * config.ScalePerPiety
				target.SetModifier("base_damage", bonus)
			},
			func() {
				target.SetModifier("base_damage", -bonus)
			})
		return singer.GetStat("pie") * config.ScalePerPiety
	}
	return 0
}

func BansheesLament(target interface{}, singer *Character, _ songPulse) int {
	switch target := target.(type) {
	case *Mob:
		damage := (instrumentRoll(singer) / 2) / len(Rooms[target.ParentId].Mobs.Contents)
		if _, err := singer.Write([]byte(text.Red + "Your song caused " + strconv.Itoa(damage) + " damage to " + target.Name + ".\n" + text.Reset)); err != nil {
			log.Println("Error writing to player:", err)
		}
		actualDamage, resisted := target.ReceiveDamageNoArmor(damage)
		data.StoreCombatMetric("banshees_lament", 0, 0, actualDamage+resisted, resisted, actualDamage, 0, singer.CharId, singer.Tier, 1, target.MobId)
		target.AddThreatDamage(damage, singer)
		return actualDamage
	}
	return 0
}
