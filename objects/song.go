package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

var Songs = map[string]map[string]string{
	"celebration-night": {
		"desc":   "Hastens stamina and vitality regeneration on players in the area.",
		"verse":  " let your spirits soar on the wings of eagles, let music wash the sleep from your eyes.",
		"effect": "celebration-night",
	},
	"curious-canticle": {
		"desc":   "prevents creatures from following anyone in the area.",
		"verse":  " now curious, get it in your head, if you follow me, then you'll be dead.",
		"effect": "curious-canticle",
	},
	"sweet-comfort": {
		"desc":   "comforts creatures, lulling them so they will not flee.",
		"verse":  " relax, my child, sleep and dream, for things aren't as bad as they seem.",
		"effect": "sweet-comfort",
	},
	"run-run-away": {
		"desc":   "lulls creatures so that they will not block escape.",
		"verse":  "see chameleon, lying there in the sun, all things to everyone, run, run away.",
		"effect": "run-away",
	},
	"draens-tale": {
		"desc":   "increases the the mana regeneration of everyone present in the area.",
		"verse":  "draen, throw your spells anew, know the weave will answer you.",
		"effect": "draens-tale",
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
	},
	"banshees-lament": {
		"desc":   "causes damage periodically to all creatures in the area.",
		"verse":  "hear my song, hear my sorrow, feel my will tear you apart",
		"effect": "banshees-lament",
	},
}

type song struct {
	target string
	effect func(target interface{}, singer *Character)
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

func DraensTale(target interface{}, singer *Character) {
	switch target := target.(type) {
	case *Character:
		target.RestoreMana(singer.GetStat("pie") * config.ScalePerPiety)
	}
}

func RunAway(target interface{}, singer *Character) {
	switch target := target.(type) {
	case *Mob:
		target.ApplyEffect(singer.Name+"_run_away", "16", 0, 0,
			func(triggers int) {
				target.ToggleFlag("run_away")
			},
			func() {
				target.ToggleFlag("run_away")
			})
	}
}

func SweetComfort(target interface{}, singer *Character) {
	switch target := target.(type) {
	case *Mob:
		target.ApplyEffect(singer.Name+"_sweet_comfort", "16", 0, 0,
			func(triggers int) {
				target.ToggleFlag("sweet_comfort")
			},
			func() {
				target.ToggleFlag("sweet_comfort")
			})
	}
}

func CuriousCanticle(target interface{}, singer *Character) {
	switch target := target.(type) {
	case *Mob:
		target.ApplyEffect(singer.Name+"_curious_canticle", "16", 0, 0,
			func(triggers int) {
				target.ToggleFlag("curious_canticle")
			},
			func() {
				target.ToggleFlag("curious_canticle")
			})
	}
}

func CelebrationNight(target interface{}, singer *Character) {
	switch target := target.(type) {
	case *Character:
		damage := (utils.Roll(singer.Equipment.Off.SidesDice, singer.Equipment.Off.NumDice, singer.Equipment.Off.PlusDice) + singer.Equipment.Off.Adjustment) / len(Rooms[target.ParentId].Chars.Contents)
		singer.Write([]byte(text.Red + "Your song healed " + strconv.Itoa(damage) + " damage to " + target.Name + ".\n" + text.Reset))
		target.Heal(damage)
	}
}

func ChampionsAnthem(target interface{}, singer *Character) {
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect(singer.Name+"_champions_anthem", "16", 0, 0,
			func(triggers int) {
				target.SetModifier("base_damage", singer.GetStat("pie")*config.ScalePerPiety)
			},
			func() {
				target.SetModifier("base_damage", -singer.GetStat("pie")*config.ScalePerPiety)
			})
	}
}

func BansheesLament(target interface{}, singer *Character) {
	switch target := target.(type) {
	case *Mob:
		damage := ((utils.Roll(singer.Equipment.Off.SidesDice, singer.Equipment.Off.NumDice, singer.Equipment.Off.PlusDice) + singer.Equipment.Off.Adjustment) / 2) / len(Rooms[target.ParentId].Mobs.Contents)
		singer.Write([]byte(text.Red + "Your song caused " + strconv.Itoa(damage) + " damage to " + target.Name + ".\n" + text.Reset))
		target.ReceiveDamageNoArmor(damage)
		target.AddThreatDamage(damage, singer)
	}
}
