package config

var AvailableClasses = make([]string, 9)
var AvailableRaces = make([]string, 14)

var AllRaces = []string{"human", "half-giant", "troll", "ogre", "dwarf", "elf", "dark-elf", "half-elf", "half-orc", "orc", "hobbit", "gnome", "sprite", "reni", "god"}


func init(){
	AvailableClasses[0] = "barbarian"
	AvailableClasses[1] = "bard"
	AvailableClasses[2] = "cleric"
	AvailableClasses[3] = "fighter"
	AvailableClasses[4] = "mage"
	AvailableClasses[5] = "monk"
	AvailableClasses[6] = "paladin"
	AvailableClasses[7] = "ranger"
	AvailableClasses[8] = "thief"

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
	AvailableRaces[13] = "reni"
}

type classDef struct {
	Desc string
	Skills string
	Stats string
	Races string
	Health int
	Stamina int
	Mana int
	WeaponAdvancement float64
}

type raceDef struct{
	Desc string
	StrMin int
	StrMax int
	DexMin int
	DexMax int
	ConMin int
	ConMax int
	IntMin int
	IntMax int
	PieMin int
	PieMax int
}


var Classes = map[string]classDef{
	"barbarian": {
		Desc:   "Raised in the harsh lands of tribal villages, barbarians are hearty warriors capable of sustaining blow after blow from opponents. The barbarian can bash its opponent, rendering them stunned for a while and unable to attack. The barbarian can circle an opponent, an excellent tactic used while fighting.",
		Skills: "Bash, Circle, Berserk",
		Stats:  "Strength, Constitution, Dexterity",
		Races:  "Half-Giant, Human, Dwarf, Orc",
		Health: 11,
		Stamina: 18,
		Mana: 1,
		WeaponAdvancement: .7,
	},
	"bard": {
		Desc: "The bard is a very clever and resourceful character class. Bards can be expected to spend their lives in search of knowledge of the arts and sciences. The young bard is usually quite foolish. The more scholarly bards can entertain an entire room with a song which will envigorate all. Also, bards have been known to charm friends and foes so completely that they are beyond harm.",
		Skills: "Sing",
		Stats: "Piety",
		Races: "Human, Elf, Dwarf, Half-elf",
		Health: 10,
		Stamina: 12,
		Mana: 10,
		WeaponAdvancement: .5,
	},
	"cleric": {
		Desc: "The cleric is the most powerful of the classes in the arts of healing. In addition to the ability to heal, the cleric can turn the undead. Clerics can also achieve magical powers only bested by a mage, and in some cases can cast spells that even a mage cannot achieve. However, their offensive spell capabilities do not stretch beyond the lower tier spells their strength lies in their curative abilities. However, clerics gain experience faster than other classes, as they reciece exp bonus for healing during the fight.",
		Skills: "Turn , Pray",
		Stats: "Piety",
		Races: "Human, Half-Elf, Gnome",
		Health: 10,
		Stamina: 11,
		Mana: 10,
		WeaponAdvancement: .4,
	},
	"fighter": {
		Desc: "The fighter is a master of the fighting arts. As the fighter advances he will achieve great proficiency in the use of weapons. The greater proficiency in weapon use is, the less chances fighter has to shatter weapon when he/she does a critical strike, vital strike. Compared to all other classes fighters gain weapon proficiency much faster than any other class. The fighter is able to bash and circle opponents like a barbarian.",
		Skills: "Hamstring, Circle",
		Stats: "Strength , Dexterity",
		Races: "Human, Half-Giant, Dwarf, Orc",
		Health: 12,
		Stamina: 16,
		Mana: 2,
		WeaponAdvancement: 1,
	},
	"mage": {
		Desc: "A master of the magic arts, the mage will gain the power to unleash incredible amounts of damage by means of spells, but magic alone cannot overcome all enemies. A mage can teach low level spells to other players, once the mage has learned them. In addition, the mage is the only class that is able to enchant things.",
		Skills: "Teach, Enchant",
		Stats: "Intelligence, Piety",
		Races: "Human, Elf, Half - Elf",
		Health: 8,
		Stamina: 10,
		Mana: 14,
		WeaponAdvancement: .4,
	},
	"monk": {
		Desc: "The monk is the master of self-discipline. By calling upon inner strength, the monk can do grave damage to foes. A monk must spend his time in self-contemplation, growing stronger all the time. He can call upon his strength to heal himself or others, or to hide from his enemies. The path of the monk is a hard one, but those few warriors who chose it will be rewarded with powers beyond most mortal men. As their Chi grows in strength they gain natural resistance to attacks.",
		Skills: "Meditate, Touch of Death, Chi Focus",
		Stats: "Piety, Constitution, Dexterity",
		Races: "Human, Dwarven",
		Health: 11,
		Stamina: 15,
		Mana: 5,
		WeaponAdvancement: .5,
	},
	"paladin": {
		Desc: "The paladin is a brave warrior of faith, and must continue to be good aligned in order to inflict damage. An evil paladin suffers greatly. The paladin is a powerful warrior and healer, and like clerics, can turn the undead. A paladin suffers a small loss if he flees from a fight. The paladin is also required to spend a term serving in the militia to show their interest in benefitting society.",
		Skills: "Turn, Pray",
		Stats: "Strength Piety",
		Races: "Human, Dwarf, Gnome, Orc",
		Health: 13,
		Stamina: 12,
		Mana: 8,
		WeaponAdvancement: .7,
	},
	"ranger": {
		Desc: "The ranger is a skillful fighter with the abilities to track opponents, to search for hidden exits, monsters, and treasures, and to hide from enemies very well. A ranger can hasten, and thus be allowed to attack faster than other classes. Rangers are necessary for some of the quests, as tracking can be required in some areas. Parties without a ranger can become hopelessly lost.",
		Skills: "Track , Haste, Snipe, Sneak",
		Stats: "Dexterity",
		Races: "Human , Halfling",
		Health: 12,
		Stamina: 14,
		Mana: 5,
		WeaponAdvancement: .7,
	},
	"thief": {
		Desc: "A thief is a very valuable player in any group of adventurers, and in some cases, necessary. A thief is capable of picking locks and stealing from opponents, and has the ability to sneak undetected from place to place.",
		Skills: "Steal, Pick, Peek, Backstab, Sneak",
		Stats: "Dexterity",
		Races: "Human, Halfling",
		Health: 12,
		Stamina: 14,
		Mana: 5,
		WeaponAdvancement: .7,
	},
}

type classTitles struct{
	Male map[int]string
	Female map[int]string
}

var ClassTitles = map[string]classTitles{
	"barbarian": {
		Male: map[int]string{1: "Barbarian", 5: "Savage", 10: "Berserker", 15: "Devastator", 20: "Annihilator"},
		Female: map[int]string{1: "Barbarian", 5: "Savage", 10: "Berserker", 15: "Devastator", 20: "Ravager"},
	},
	"bard": {
		Male: map[int]string{1: "Jongleur", 5: "Skald", 10: "Minstrel", 15: "Muse", 20: "Loremaster"},
		Female: map[int]string{1: "Jongleur", 5: "Skald", 10: "Minstrel", 15: "Muse", 20: "Loremistress"},
	},
	"cleric": {
		Male: map[int]string{1: "Cleric", 5: "Adept", 10: "Priest", 15: "High-Priest", 20: "Prophet"},
		Female: map[int]string{1: "Cleric", 5: "Adept", 10: "Priestess", 15: "High-Priestess", 20: "Prophetess"},
	},
	"fighter": {
		Male: map[int]string{1: "Fighter", 5: "Warrior", 10: "Myrmidon", 15: "Champion", 20: "Warlord"},
		Female: map[int]string{1: "Fighter", 5: "Warrior", 10: "Myrmidon", 15: "Champion", 20: "Warmistress"},
	},
	"mage": {
		Male: map[int]string{1: "Apprentice", 5: "Mage", 10: "Wizard", 15: "Arch-wizard", 20: "Weavemaster"},
		Female: map[int]string{1: "Apprentice", 5: "Mage", 10: "Wizardress", 15: "Arch-wizardress", 20: "Weavemistress"},
	},
	"monk": {
		Male: map[int]string{1: "Initiate", 5: "Disciple", 10: "Immaculate", 15: "Master", 20: "Philosopher"},
		Female: map[int]string{1: "Initiate", 5: "Disciple", 10: "Immaculate", 15: "Master", 20: "Philosopher"},
	},
	"paladin": {
		Male: map[int]string{1: "Cavalier", 5: "Warder", 10: "Holy Warrior", 15: "Lord", 20: "Crusader"},
		Female: map[int]string{1: "Cavalier", 5: "Warder", 10: "Holy Warrior", 15: "Lady", 20: "Crusader"},
	},
	"ranger": {
		Male: map[int]string{1: "Ranger", 5: "Scout", 10: "Pathfinder", 15: "Trailblazer", 20: "Waymaker"},
		Female: map[int]string{1: "Ranger", 5: "Scout", 10: "Pathfinder", 15: "Trailblazer", 20: "Waymaker"},
	},
	"thief": {
		Male: map[int]string{1: "Rogue", 5: "Pickpocket", 10: "Nightblade", 15: "Shadow Warrior", 20: "Shadowmaster"},
		Female: map[int]string{1: "Rogue", 5: "Pickpocket", 10: "Nightblade", 15: "Shadow Warrior", 20: "Shadowmistress"},
	},
}

func ClassTitle(class int64, gender string, tier int64) string{
	if class >= 50 && class < 60 {
		return "Builder"
	}
	if class >= 60 && class < 100 {
		return "DungeonMaster"
	}
	if class == 100 {
		return "GameMaster"
	}
	var selectTier int
	switch {
	case tier < 5:
		selectTier = 1
	case tier < 10:
		selectTier = 5
	case tier < 15:
		selectTier = 10
	case tier < 20:
		selectTier = 15
	case tier < 25:
		selectTier = 20
	}
	if gender == "m"{
		return ClassTitles[AvailableClasses[class]].Male[selectTier]
	}else{
		return ClassTitles[AvailableClasses[class]].Female[selectTier]
	}

}

var TextGender = map[string]string{"m": "male", "f": "female"}

var RaceDefs = map[string]raceDef{
	"dark-elf": {
		Desc: "The dark elves were the only breed of elf to have never seen the Light of the Two Trees.  They are frequently scorned by the other breeds of elves and many consider it a disgrace to even be seen with them.",
		StrMin: 4, StrMax: 28,
		DexMin: 12, DexMax: 40,
		ConMin: 5, ConMax: 20,
		IntMin: 12, IntMax: 34,
		PieMin: 4, PieMax: 28,
	},
	"dwarf": {
		Desc: "A dwarf is a stocky and short demihuman, standing about 4 feet tall.  Dwarves are sturdy fighters, and are known to be stubborn and practical.",
		StrMin: 10, StrMax: 34,
		DexMin: 3, DexMax: 24,
		ConMin: 12, ConMax: 38,
		IntMin: 4, IntMax: 24,
		PieMin: 8, PieMax: 30,
	},
	"elf": {
		Desc: "Somewhat shorter than humans, the elf is of weaker constitution and higher intelligence.",
		StrMin: 4, StrMax: 26,
		DexMin: 12, DexMax: 40,
		ConMin: 5, ConMax: 22,
		IntMin: 12, IntMax: 36,
		PieMin: 4, PieMax: 26,
	},
	"gnome": {
		Desc: "A cousin of the dwarf, gnomes are small demihumans which can become very capable clerics and paladins.",
		StrMin: 5, StrMax: 20,
		DexMin: 4, DexMax: 23,
		ConMin: 6, ConMax: 20,
		IntMin: 12, IntMax: 42,
		PieMin: 12, PieMax: 45,
	},
	"half-giant": {
		Desc: "A cross between the giant and human races, a half-giant is brutally strong and makes a very good warrior.",
		StrMin: 18, StrMax: 45,
		DexMin: 2, DexMax: 16,
		ConMin: 14, ConMax: 42,
		IntMin: 2, IntMax: 22,
		PieMin: 2, PieMax: 25,
	},
	"half-elf": {
		Desc: "A cross between the elven and human races, a half-elf can become a master in any class.",
		StrMin: 5, StrMax: 28,
		DexMin: 9, DexMax: 34,
		ConMin: 4, ConMax: 26,
		IntMin: 5, IntMax: 32,
		PieMin: 5, PieMax: 30,
	},
	"hobbit": {
		Desc: "Small and agile, the hobbit specializes in dexterity, and thus makes a good thief, or ranger. They are also known to other races as halflings, but they prefer to be called by their chosen name of hobbit.",
		StrMin: 4, StrMax: 24,
		DexMin: 14, DexMax: 43,
		ConMin: 4, ConMax: 23,
		IntMin: 5, IntMax: 30,
		PieMin: 5, PieMax: 30,
	},
	"half-orc": {
		Desc: "The result of a failed attempt to make an orc that is closer to an elf, these half breeds are hated by both orcs and elves.",
		StrMin: 8, StrMax: 32,
		DexMin: 4, DexMax: 34,
		ConMin: 8, ConMax: 32,
		IntMin: 5, IntMax: 32,
		PieMin: 2, PieMax: 20,
	},
	"human": {
		Desc: "What is man? Who knows? And if you are actually reading this, perhaps you should stop mudding for about a week, and read philosophy.",
		StrMin: 5, StrMax: 30,
		DexMin: 5, DexMax: 30,
		ConMin: 5, ConMax: 30,
		IntMin: 5, IntMax: 30,
		PieMin: 5, PieMax: 30,
	},
	"ogre":  {
		Desc: "Large and strong, this powerful race can also excel at physical combat but are generally not well versed in the magical arts.",
		StrMin: 17, StrMax: 43,
		DexMin: 3, DexMax: 28,
		ConMin: 14, ConMax: 43,
		IntMin: 1, IntMax: 16,
		PieMin: 1, PieMax: 20,
	},
	"orc": {
		Desc: "Orcs are fierce warriors, who in their homelands prefer banding together for hunting and raiding. Orcs are strong and make good warriors. They were created in mockery of elves and like elves they do not die naturally. They are weakened by the sun and prefer the dark.",
		StrMin: 12, StrMax: 36,
		DexMin: 5, DexMax: 32,
		ConMin: 12, ConMax: 36,
		IntMin: 3, IntMax: 24,
		PieMin: 3, PieMax: 22,
	},
	"renis": {
		Desc: "The Renis are a scholarly race, once responsible for maintaining all of the knowledge of the Allied Races. Renis are a tall, slender people, half again as tall as humans, though weighing slightly less. Renis are covered in very short fur, usually pale blue in color, but often grey, green or even rarely black. Renis ears end in points, similar to those of elves, however the points are more severe. Like most races, the Renis have hair, which is always the colors of a precious gemstone.",
		StrMin: 3, StrMax: 20,
		DexMin: 4, DexMax: 25,
		ConMin: 3, ConMax: 20,
		IntMin: 17, IntMax: 45,
		PieMin: 8, PieMax: 38,
	},
	"troll": {
		Desc: "Trolls are an evil race. Large, strong, ugly and stupid, they enjoy to hoard treasure, kill for pleasure and eat raw flesh. Trolls generally prefer to travel alone, but can sometimes be found in groups of three or more.",
		StrMin: 15, StrMax: 40,
		DexMin: 2, DexMax: 22,
		ConMin: 17, ConMax: 45,
		IntMin: 2, IntMax: 20,
		PieMin: 3, PieMax: 23,
	},
	"sprite": {
		Desc: "Tiny and mischievious, agile and swift,Sprites are like in  personalities to elves and gnomes. Sprites are woodland creatures who are connected to the life force of all things, making them excellent magi and rangers. Sprites naturally levitate and can detect-invisible objects and creatures.",
		StrMin: 1, StrMax: 16,
		DexMin: 17, DexMax: 45,
		ConMin: 4, ConMax: 18,
		IntMin: 7, IntMax: 38,
		PieMin: 2, PieMax: 33,
	},
}



