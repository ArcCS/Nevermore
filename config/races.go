package config

var AvailableRaces = make([]string, 14)

// var AllRaces = []string{"human", "half-giant", "troll", "ogre", "dwarf", "elf", "dark-elf", "half-elf", "half-orc", "orc", "hobbit", "gnome", "sprite", "renis", "god"}
const HUMAN, HALF_GIANT, TROLL, OGRE, DWARF, ELF, DARK_ELF, HALF_ELF, HALF_ORC, ORC, HOBBIT, GNOME, SPRITE, RENIS = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13

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
	Darkvision bool
}

var RaceDefs = map[string]raceDef{
	"dark-elf": {
		Desc: "Dark-Elves tend to be split between Drow, the underdark dwelling-typically evil worshipping elves, and the Grey-Elves, those who have moved to the surface and keep the traits of their former lineage but typically worship good entities.  They are frequently scorned by the other breeds of elves and many consider it a disgrace to even be seen with them.  Physically drow tend to have gray dark skin, white hair, and red eyes.  Grey-elves tend to have pale skin, white hair, and blue eyes.  Variations in tones are present with blueish undertones and even sometimes being born with traits from their original ancient elven ancestry. Like elves and dwarves, they have long life spans.",
		// StrShift: -2, DexShift: 1, ConShift: -1, IntShift: 1, PieShift: 1,
		StrMin: 4, StrMax: 28,
		DexMin: 12, DexMax: 40,
		ConMin: 5, ConMax: 20,
		IntMin: 12, IntMax: 34,
		PieMin: 4, PieMax: 28,
		Darkvision: true,
	},
	"dwarf": {
		Desc: "A dwarf is a stocky and short demihuman, standing about 4 feet tall.  Dwarves are sturdy fighters, and are known to be stubborn and practical.  Their unground dwelling kin tend to be known as greedy for their obsession with precious gems and metals. However this may be a misrepresentation of their true nature as they are known to be loyal and trustworthy friends.  Dwarves skin tones come in a variety of colors from pale to dark brown, with hair colors ranging from black to red to white.  They have long life spans.",
		// StrShift: 2, DexShift: -2, ConShift: 2, IntShift: -2, PieShift: 0,
		StrMin: 10, StrMax: 34,
		DexMin: 3, DexMax: 24,
		ConMin: 12, ConMax: 38,
		IntMin: 4, IntMax: 24,
		PieMin: 8, PieMax: 30,
		Darkvision: true,
	},
	"elf": {
		Desc: "Somewhat shorter than humans, elves tend to be brilliant and agile.  They are known to be very wise and are often found in the magical arts.  They are also known to be very good rangers. Their skin tones tend to be lighter, with multiple variations of hair color and eye color. They have long life spans.",
		// StrShift: -3, DexShift: 2, ConShift: -3, IntShift: 2, PieShift: 2,
		StrMin: 4, StrMax: 26,
		DexMin: 12, DexMax: 40,
		ConMin: 5, ConMax: 22,
		IntMin: 12, IntMax: 36,
		PieMin: 4, PieMax: 26,
		Darkvision: true,
	},
	"gnome": {
		Desc: "A cousin of the dwarf, gnomes are small and more agile than their dwarf counterparts.  They tend to be more intelligent and inclined to magical studies especially in the craft of artificing.   They have multiple variations of skin tones, eye colors, and hair colors. They have long life spans. ",
		// StrShift: -3, DexShift: 2, ConShift: -3, IntShift: 1, PieShift: 3,
		StrMin: 5, StrMax: 20,
		DexMin: 4, DexMax: 23,
		ConMin: 6, ConMax: 20,
		IntMin: 12, IntMax: 42,
		PieMin: 12, PieMax: 45,
		Darkvision: false,
	},
	"half-giant": {
		Desc: "A cross between the giant and human races, a half-giant is brutally strong and makes a very good warrior. They are typically known for not being as intelligent as humans, or pious.  They share many traits from their human lineage with a smattering of traits from the giant side of their blood, usually affecting their skin tones and hair color.  They tend to stand from 8ft to 13ft tall.  They are short lived people with the oldest of their living to no more than 80 years old",
		// StrShift: 4, DexShift: -2, ConShift: 4, IntShift: -4, PieShift: -2,
		StrMin: 18, StrMax: 45,
		DexMin: 2, DexMax: 16,
		ConMin: 14, ConMax: 42,
		IntMin: 2, IntMax: 22,
		PieMin: 2, PieMax: 25,
		Darkvision: false,
	},
	"half-elf": {
		Desc: "A cross between the elven and human races, they can share traits in either of their lineages. They have a typical human height range, tend to be smarter and more agile than their human counter parts, but not as pious as their families.  They have a medium to long life span, being much shorter due to their human blood.",
		// StrShift: -1, DexShift: 1, ConShift: -1, IntShift: 1, PieShift: 0,
		StrMin: 5, StrMax: 28,
		DexMin: 9, DexMax: 34,
		ConMin: 4, ConMax: 26,
		IntMin: 5, IntMax: 32,
		PieMin: 5, PieMax: 30,
		Darkvision: false,
	},
	"hobbit": {
		Desc: "Small and agile, the hobbit specializes in dexterity, and thus makes a good thief, or ranger. They are also known to other races as halflings, but they prefer to be called by their chosen name of hobbit.  They are hardy people and pleasure seekers. They tend to have joyous outlooks and liven the world around them.  Often being found as bards, or entertainers.  They have a medium long life span, from 200-300 years.",
		// StrShift: -3, DexShift: 3, ConShift: -2, IntShift: 1, PieShift: 1,
		StrMin: 4, StrMax: 24,
		DexMin: 14, DexMax: 43,
		ConMin: 4, ConMax: 23,
		IntMin: 5, IntMax: 30,
		PieMin: 5, PieMax: 30,
		Darkvision: false,
	},
	"half-orc": {
		Desc: "Half orcs are an interbreed between either an orc and an elf or an orc and a human.  They can display traits from either side of their blood but do tend to have a more orcish appearance.  They are typically not as intelligent as their elven or human counterparts, but are more agile and stronger.  Human-Orcs tend to have a short lifespan, no longer than 80-90 years, while elven-orcs tend to have a medium lifespan up to 200 years.",
		// StrShift: 2, DexShift: 1, ConShift: 1, IntShift: -2, PieShift: -2,
		StrMin: 8, StrMax: 32,
		DexMin: 4, DexMax: 34,
		ConMin: 8, ConMax: 32,
		IntMin: 5, IntMax: 32,
		PieMin: 2, PieMax: 20,
		Darkvision: false,
	},
	"human": {
		Desc: "Humans are the truly average, short lived race of the realm of Altin.  They are veratile and can be found dedicating themselves to many professions and walks of life.  They have short lifespans, typically no more than 100 years.  They have a wide range of skin tones, eye colors, and hair colors.",
		// StrShift: 0, DexShift: 0, ConShift: 0, IntShift: 0, PieShift: 0,
		StrMin: 5, StrMax: 30,
		DexMin: 5, DexMax: 30,
		ConMin: 5, ConMax: 30,
		IntMin: 5, IntMax: 30,
		PieMin: 5, PieMax: 30,
		Darkvision: false,
	},
	"ogre": {
		Desc: "Large and strong, this powerful race can also excel at physical combat but are generally not well versed in the magical arts.  Ogres are not known for their intellectual prowess, dexterity or piety.  The ogres were welcomed to the city of all races long ago as their people grew more social and involved with the people of Nexus. They are a short lived people with the oldest of their living to no more than 80 years old.",
		//StrShift: 3, DexShift: -1, ConShift: 3, IntShift: -4, PieShift: -1,
		StrMin: 17, StrMax: 43,
		DexMin: 3, DexMax: 28,
		ConMin: 14, ConMax: 43,
		IntMin: 1, IntMax: 16,
		PieMin: 1, PieMax: 20,
		Darkvision: false,
	},
	"orc": {
		Desc: "Orcs are fierce warriors, who in their homelands prefer banding together for hunting and raiding. Orcs are strong and make good warriors. They have variations of grey and green skin, tusks, and tend to be taller than humans, ranging from 6ft to 7ft tall.  They have a short life span, typically no more than 80 years.  At one point they were magically endowed with long life spans, however as the presence of the Gods and magic waned, so did their life spans. ",
		// StrShift: 2, DexShift: 2, ConShift: 2, IntShift: -3, PieShift: -3,
		StrMin: 12, StrMax: 36,
		DexMin: 5, DexMax: 32,
		ConMin: 12, ConMax: 36,
		IntMin: 3, IntMax: 24,
		PieMin: 3, PieMax: 22,
		Darkvision: true,
	},
	"renis": {
		Desc: "The Renis are a scholarly race, once responsible for maintaining all of the knowledge of the Allied Races. Renis are a tall, slender people, half again as tall as humans, though weighing slightly less. Renis are covered in very short fur, usually pale blue in color, but often grey, green or even rarely black. Renis ears end in points, similar to those of elves, however the points are more severe. Like most races, the Renis have hair, which is always the colors of a precious gemstone.",
		// StrShift: -3, DexShift: 0, ConShift: -3, IntShift: 5, PieShift: 1,
		StrMin: 3, StrMax: 20,
		DexMin: 4, DexMax: 25,
		ConMin: 3, ConMax: 20,
		IntMin: 17, IntMax: 45,
		PieMin: 8, PieMax: 38,
		Darkvision: true,
	},
	"troll": {
		Desc: "Trolls used to be considered an evil race of people, typically known for hoarding treasure and killing for pleasure and eat raw flesh.  Overtime they have grown more social and contribute to the collective goals of society in altin. Despite their involvement in society trolls still1 generally prefer to travel alone, but can sometimes be found in groups of three or more.  They tend to be very strong, but not very intelligent.  Their skin comes in various hues of green.  They have a short life spans, typically no more than 50 years.",
		// StrShift: 3, DexShift: 0, ConShift: 3, IntShift: -3, PieShift: -3,
		StrMin: 15, StrMax: 40,
		DexMin: 2, DexMax: 22,
		ConMin: 17, ConMax: 45,
		IntMin: 2, IntMax: 20,
		PieMin: 3, PieMax: 23,
		Darkvision: true,
	},
	"sprite": {
		Desc: "Tiny and mischievous, agile and swift. The sprites are a magical race of people that tend to be very smart.  Sprites are woodland creatures who are connected to the life force of all things, making them excellent magi and rangers.  They tend to be around a foot tall, and use magic to allow them to handle objects like their other heoric friends.  Sprites tend to have colorful skin tones, hair colors, and eye colors.  They have long life spans.",
		// StrShift: -4, DexShift: 4, ConShift: -4, IntShift: 4, PieShift: 0,
		StrMin: 1, StrMax: 16,
		DexMin: 17, DexMax: 45,
		ConMin: 4, ConMax: 18,
		IntMin: 7, IntMax: 38,
		PieMin: 2, PieMax: 33,
		Darkvision: true,
	},
}
