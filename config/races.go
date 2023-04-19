package config

var AvailableRaces = make([]string, 14)

var AllRaces = []string{"human", "half-giant", "troll", "ogre", "dwarf", "elf", "dark-elf", "half-elf", "half-orc", "orc", "hobbit", "gnome", "sprite", "renis", "god"}

func init() {
	AvailableRaces[0] = "human"
	AvailableRaces[1] = "half-giant"
	AvailableRaces[2] = "troll"
	AvailableRaces[3] = "ogre"
	AvailableRaces[4] = "dwarf"
	AvailableRaces[5] = "elf"
	AvailableRaces[6] = "dark-elf"
	AvailableRaces[7] = "half-elf"
	AvailableRaces[8] = "half-orc"
	AvailableRaces[9] = "orc"
	AvailableRaces[10] = "hobbit"
	AvailableRaces[11] = "gnome"
	AvailableRaces[12] = "sprite"
	AvailableRaces[13] = "renis"
}

type raceDef struct {
	Desc       string
	StrMin     int
	StrMax     int
	DexMin     int
	DexMax     int
	ConMin     int
	ConMax     int
	IntMin     int
	IntMax     int
	PieMin     int
	PieMax     int
	MinAge     int
	Darkvision bool
}

var RaceDefs = map[string]raceDef{
	"dark-elf": {
		Desc:   "The dark elves were the only breed of elf to have never seen the Light of the Two Trees.  They are frequently scorned by the other breeds of elves and many consider it a disgrace to even be seen with them.",
		StrMin: 4, StrMax: 28,
		DexMin: 12, DexMax: 40,
		ConMin: 5, ConMax: 20,
		IntMin: 12, IntMax: 34,
		PieMin: 4, PieMax: 28,
		MinAge:     70,
		Darkvision: true,
	},
	"dwarf": {
		Desc:   "A dwarf is a stocky and short demihuman, standing about 4 feet tall.  Dwarves are sturdy fighters, and are known to be stubborn and practical.",
		StrMin: 10, StrMax: 34,
		DexMin: 3, DexMax: 24,
		ConMin: 12, ConMax: 38,
		IntMin: 4, IntMax: 24,
		PieMin: 8, PieMax: 30,
		MinAge:     30,
		Darkvision: true,
	},
	"elf": {
		Desc:   "Somewhat shorter than humans, the elf is of weaker constitution and higher intelligence.",
		StrMin: 4, StrMax: 26,
		DexMin: 12, DexMax: 40,
		ConMin: 5, ConMax: 22,
		IntMin: 12, IntMax: 36,
		PieMin: 4, PieMax: 26,
		MinAge:     8,
		Darkvision: true,
	},
	"gnome": {
		Desc:   "A cousin of the dwarf, gnomes are small demihumans which can become very capable clerics and paladins.",
		StrMin: 5, StrMax: 20,
		DexMin: 4, DexMax: 23,
		ConMin: 6, ConMax: 20,
		IntMin: 12, IntMax: 42,
		PieMin: 12, PieMax: 45,
		MinAge:     20,
		Darkvision: false,
	},
	"half-giant": {
		Desc:   "A cross between the giant and human races, a half-giant is brutally strong and makes a very good warrior.",
		StrMin: 18, StrMax: 45,
		DexMin: 2, DexMax: 16,
		ConMin: 14, ConMax: 42,
		IntMin: 2, IntMax: 22,
		PieMin: 2, PieMax: 25,
		MinAge:     17,
		Darkvision: false,
	},
	"half-elf": {
		Desc:   "A cross between the elven and human races, a half-elf can become a master in any class.",
		StrMin: 5, StrMax: 28,
		DexMin: 9, DexMax: 34,
		ConMin: 4, ConMax: 26,
		IntMin: 5, IntMax: 32,
		PieMin: 5, PieMax: 30,
		MinAge:     24,
		Darkvision: false,
	},
	"hobbit": {
		Desc:   "Small and agile, the hobbit specializes in dexterity, and thus makes a good thief, or ranger. They are also known to other races as halflings, but they prefer to be called by their chosen name of hobbit.",
		StrMin: 4, StrMax: 24,
		DexMin: 14, DexMax: 43,
		ConMin: 4, ConMax: 23,
		IntMin: 5, IntMax: 30,
		PieMin: 5, PieMax: 30,
		MinAge:     35,
		Darkvision: false,
	},
	"half-orc": {
		Desc:   "The result of a failed attempt to make an orc that is closer to an elf, these half breeds are hated by both orcs and elves.",
		StrMin: 8, StrMax: 32,
		DexMin: 4, DexMax: 34,
		ConMin: 8, ConMax: 32,
		IntMin: 5, IntMax: 32,
		PieMin: 2, PieMax: 20,
		MinAge:     20,
		Darkvision: false,
	},
	"human": {
		Desc:   "What is man? Who knows? And if you are actually reading this, perhaps you should stop mudding for about a week, and read philosophy.",
		StrMin: 5, StrMax: 30,
		DexMin: 5, DexMax: 30,
		ConMin: 5, ConMax: 30,
		IntMin: 5, IntMax: 30,
		PieMin: 5, PieMax: 30,
		MinAge:     18,
		Darkvision: false,
	},
	"ogre": {
		Desc:   "Large and strong, this powerful race can also excel at physical combat but are generally not well versed in the magical arts.",
		StrMin: 17, StrMax: 43,
		DexMin: 3, DexMax: 28,
		ConMin: 14, ConMax: 43,
		IntMin: 1, IntMax: 16,
		PieMin: 1, PieMax: 20,
		MinAge:     18,
		Darkvision: false,
	},
	"orc": {
		Desc:   "Orcs are fierce warriors, who in their homelands prefer banding together for hunting and raiding. Orcs are strong and make good warriors. They were created in mockery of elves and like elves they do not die naturally. They are weakened by the sun and prefer the dark.",
		StrMin: 12, StrMax: 36,
		DexMin: 5, DexMax: 32,
		ConMin: 12, ConMax: 36,
		IntMin: 3, IntMax: 24,
		PieMin: 3, PieMax: 22,
		MinAge:     20,
		Darkvision: true,
	},
	"renis": {
		Desc:   "The Renis are a scholarly race, once responsible for maintaining all of the knowledge of the Allied Races. Renis are a tall, slender people, half again as tall as humans, though weighing slightly less. Renis are covered in very short fur, usually pale blue in color, but often grey, green or even rarely black. Renis ears end in points, similar to those of elves, however the points are more severe. Like most races, the Renis have hair, which is always the colors of a precious gemstone.",
		StrMin: 3, StrMax: 20,
		DexMin: 4, DexMax: 25,
		ConMin: 3, ConMax: 20,
		IntMin: 17, IntMax: 45,
		PieMin: 8, PieMax: 38,
		MinAge:     50,
		Darkvision: true,
	},
	"troll": {
		Desc:   "Trolls are an evil race. Large, strong, ugly and stupid, they enjoy to hoard treasure, kill for pleasure and eat raw flesh. Trolls generally prefer to travel alone, but can sometimes be found in groups of three or more.",
		StrMin: 15, StrMax: 40,
		DexMin: 2, DexMax: 22,
		ConMin: 17, ConMax: 45,
		IntMin: 2, IntMax: 20,
		PieMin: 3, PieMax: 23,
		MinAge:     18,
		Darkvision: true,
	},
	"sprite": {
		Desc:   "Tiny and mischievous, agile and swift,Sprites are like in  personalities to elves and gnomes. Sprites are woodland creatures who are connected to the life force of all things, making them excellent magi and rangers. Sprites naturally levitate and can detect-invisible objects and creatures.",
		StrMin: 1, StrMax: 16,
		DexMin: 17, DexMax: 45,
		ConMin: 4, ConMax: 18,
		IntMin: 7, IntMax: 38,
		PieMin: 2, PieMax: 33,
		MinAge:     30,
		Darkvision: true,
	},
}
