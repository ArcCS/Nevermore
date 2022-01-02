package objects

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
	"sweet-comfort:": {
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
		"verse":  " draen, throw your spells anew, know the weave will answer you.",
		"effect": "draens-tale",
	},
}

var SongEffects = map[string]map[string]interface{}{
	"draens-tale": {"target": "players", "effect": DraensTale},
	"run-away": {"target": "mobs", "effect": RunAway},
	"sweet-comfort": {"target": "mobs", "effect": SweetComfort},
	"curious-canticle": {"target": "mobs", "effect": CuriousCanticle},
	"celebration-night": {"target": "players", "effect": CelebrationNight},
}

func DraensTale(c *Character, singer *Character){
	c.ApplyEffect(singer.Name + "_draens_tale", "16", "0",
		func() {
			c.ToggleFlag("draens_tale", singer.Name + "_celebration_night")
		},
		func() {
			c.ToggleFlag("draens_tale", singer.Name + "_celebration_night")
		})
}

func RunAway(m *Mob, singer *Character){
	m.ApplyEffect(singer.Name + "_run_away", "16", "0",
		func() {
			m.ToggleFlag("run_away")
		},
		func() {
			m.ToggleFlag("run_away")
		})

}

func SweetComfort(m *Mob, singer *Character){
	m.ApplyEffect(singer.Name + "_sweet_comfort", "16", "0",
		func() {
			m.ToggleFlag("sweet_comfort")
		},
		func() {
			m.ToggleFlag("sweet_comfort")
		})
}

func CuriousCanticle(m *Mob, singer *Character){
	m.ApplyEffect(singer.Name + "_curious_canticle", "16", "0",
		func() {
			m.ToggleFlag("curious_canticle")
		},
		func() {
			m.ToggleFlag("curious_canticle")
		})
}

func CelebrationNight(c *Character, singer *Character){
	c.ApplyEffect(singer.Name + "_celebration_night", "16", "0",
		func() {
			c.ToggleFlag("celebration_night", singer.Name + "_celebration_night")
		},
		func() {
			c.ToggleFlag("celebration_night", singer.Name + "_celebration_night")
		})
}