package config

import "github.com/ArcCS/Nevermore/permissions"

var AvailableClasses = make([]string, 9)
var ClassPerms = make([]permissions.Permissions, 9)
var StartingGear = make(map[int][]int, 9)

func init() {
	AvailableClasses[0] = "fighter"
	AvailableClasses[1] = "barbarian"
	AvailableClasses[2] = "thief"
	AvailableClasses[3] = "ranger"
	AvailableClasses[4] = "mage"
	AvailableClasses[5] = "cleric"
	AvailableClasses[6] = "paladin"
	AvailableClasses[7] = "bard"
	AvailableClasses[8] = "monk"

	ClassPerms[0] = permissions.Fighter
	ClassPerms[1] = permissions.Barbarian
	ClassPerms[2] = permissions.Thief
	ClassPerms[3] = permissions.Ranger
	ClassPerms[4] = permissions.Mage
	ClassPerms[5] = permissions.Cleric
	ClassPerms[6] = permissions.Paladin
	ClassPerms[7] = permissions.Bard
	ClassPerms[8] = permissions.Monk

	StartingGear[0] = []int{129, 3615}
	StartingGear[1] = []int{227, 3615}
	StartingGear[2] = []int{2628, 3615}
	StartingGear[3] = []int{665, 3615}
	StartingGear[4] = []int{665, 1180, 3343, 3615}
	StartingGear[5] = []int{227, 1179, 242, 3615}
	StartingGear[6] = []int{129, 1179, 3615}
	StartingGear[7] = []int{665, 129, 1179, 3615, 3430, 3848, 3398}
	StartingGear[8] = []int{3343, 3615}
}

type classDef struct {
	Desc              string
	Skills            string
	Stats             string
	Armor             string
	Races             string
	Health            int
	Stamina           int
	Mana              int
	WeaponAdvancement float64
}

var Classes = map[string]classDef{
	"barbarian": {
		Desc:              "Raised in the harsh lands of tribal villages, barbarians are hearty warriors capable of sustaining blow after blow from opponents. The barbarian can bash its opponent, rendering them stunned for a while and unable to attack. The barbarian can circle an opponent, an excellent tactic used while fighting.",
		Skills:            "Bash, Circle, Berserk",
		Stats:             "Strength, Constitution, Dexterity",
		Armor:             "Light, Medium, Heavy",
		Races:             "Half-Giant, Human, Dwarf, Orc",
		Health:            11,
		Stamina:           18,
		Mana:              1,
		WeaponAdvancement: .7,
	},
	"bard": {
		Desc:              "The bard is a very clever and resourceful character class. Bards can be expected to spend their lives in search of knowledge of the arts and sciences. The young bard is usually quite foolish. The more scholarly bards can entertain an entire room with a song which will envigorate all. Also, bards have been known to charm friends and foes so completely that they are beyond harm.",
		Skills:            "Sing",
		Stats:             "Piety",
		Armor:             "Light",
		Races:             "Human, Elf, Dwarf, Half-elf",
		Health:            10,
		Stamina:           12,
		Mana:              10,
		WeaponAdvancement: .5,
	},
	"cleric": {
		Desc:              "The cleric is the most powerful of the classes in the arts of healing. In addition to the ability to heal, the cleric can turn the undead. Clerics can also achieve magical powers only bested by a mage, and in some cases can cast spells that even a mage cannot achieve. However, their offensive spell capabilities do not stretch beyond the lower tier spells their strength lies in their curative abilities. However, clerics gain experience faster than other classes, as they reciece exp bonus for healing during the fight.",
		Skills:            "Turn , Pray",
		Stats:             "Piety",
		Armor:             "Robes",
		Races:             "Human, Half-Elf, Gnome",
		Health:            10,
		Stamina:           11,
		Mana:              10,
		WeaponAdvancement: .4,
	},
	"fighter": {
		Desc:              "The fighter is a master of the fighting arts. As the fighter advances he will achieve great proficiency in the use of weapons. The greater proficiency in weapon use is, the less chances fighter has to shatter weapon when he/she does a critical strike, vital strike. Compared to all other classes fighters gain weapon proficiency much faster than any other class. The fighter is able to bash and circle opponents like a barbarian.",
		Skills:            "Hamstring, Circle",
		Stats:             "Strength , Dexterity",
		Armor:             "Light, Medium, Heavy",
		Races:             "Human, Half-Giant, Dwarf, Orc",
		Health:            12,
		Stamina:           16,
		Mana:              2,
		WeaponAdvancement: 1,
	},
	"mage": {
		Desc:              "A master of the magic arts, the mage will gain the power to unleash incredible amounts of damage by means of spells, but magic alone cannot overcome all enemies. A mage can teach low level spells to other players, once the mage has learned them. In addition, the mage is the only class that is able to enchant things.",
		Skills:            "Teach, Enchant",
		Stats:             "Intelligence, Piety",
		Armor:             "Robes",
		Races:             "Human, Elf, Half - Elf",
		Health:            8,
		Stamina:           10,
		Mana:              14,
		WeaponAdvancement: .4,
	},
	"monk": {
		Desc:              "The monk is the master of self-discipline. By calling upon inner strength, the monk can do grave damage to foes. A monk must spend his time in self-contemplation, growing stronger all the time. He can call upon his strength to heal himself or others, or to hide from his enemies. The path of the monk is a hard one, but those few warriors who chose it will be rewarded with powers beyond most mortal men. As their Chi grows in strength they gain natural resistance to attacks.",
		Skills:            "Meditate, Touch of Death",
		Stats:             "Piety, Constitution, Dexterity",
		Armor:             "Constitution Based",
		Races:             "Human, Dwarven",
		Health:            11,
		Stamina:           15,
		Mana:              2,
		WeaponAdvancement: .5,
	},
	"paladin": {
		Desc:              "The paladin is a brave warrior of faith, and must continue to be good aligned in order to inflict damage. An evil paladin suffers greatly. The paladin is a powerful warrior and healer, and like clerics, can turn the undead. A paladin suffers a small loss if he flees from a fight. The paladin is also required to spend a term serving in the militia to show their interest in benefitting society.",
		Skills:            "Turn, Pray, Shield Slam",
		Stats:             "Strength Piety",
		Armor:             "Light, Medium, Heavy",
		Races:             "Human, Dwarf, Gnome, Orc",
		Health:            13,
		Stamina:           12,
		Mana:              8,
		WeaponAdvancement: .7,
	},
	"ranger": {
		Desc:              "The ranger is a skillful fighter with the abilities to track opponents, to search for hidden exits, monsters, and treasures, and to hide from enemies very well. A ranger can hasten, and thus be allowed to attack faster than other classes. Rangers are necessary for some of the quests, as tracking can be required in some areas. Parties without a ranger can become hopelessly lost.",
		Skills:            "Haste, Snipe, Sneak",
		Stats:             "Dexterity",
		Armor:             "Light, Medium",
		Races:             "Human , Halfling",
		Health:            12,
		Stamina:           14,
		Mana:              5,
		WeaponAdvancement: .7,
	},
	"thief": {
		Desc:              "A thief is a very valuable player in any group of adventurers, and in some cases, necessary. A thief is capable of picking locks and stealing from opponents, and has the ability to sneak undetected from place to place.",
		Skills:            "Steal, Pick, Peek, Backstab, Sneak",
		Stats:             "Dexterity",
		Armor:             "Light, Medium",
		Races:             "Human, Halfling",
		Health:            12,
		Stamina:           14,
		Mana:              5,
		WeaponAdvancement: .7,
	},
}

type classTitles struct {
	Male   map[int]string
	Female map[int]string
}

var ClassTitles = map[string]classTitles{
	"barbarian": {
		Male:   map[int]string{1: "Barbarian", 5: "Savage", 10: "Berserker", 15: "Devastator", 20: "Annihilator"},
		Female: map[int]string{1: "Barbarian", 5: "Savage", 10: "Berserker", 15: "Devastator", 20: "Ravager"},
	},
	"bard": {
		Male:   map[int]string{1: "Jongleur", 5: "Skald", 10: "Minstrel", 15: "Muse", 20: "Loremaster"},
		Female: map[int]string{1: "Jongleur", 5: "Skald", 10: "Minstrel", 15: "Muse", 20: "Loremistress"},
	},
	"cleric": {
		Male:   map[int]string{1: "Cleric", 5: "Adept", 10: "Priest", 15: "High-Priest", 20: "Prophet"},
		Female: map[int]string{1: "Cleric", 5: "Adept", 10: "Priestess", 15: "High-Priestess", 20: "Prophetess"},
	},
	"fighter": {
		Male:   map[int]string{1: "Fighter", 5: "Warrior", 10: "Myrmidon", 15: "Champion", 20: "Warlord"},
		Female: map[int]string{1: "Fighter", 5: "Warrior", 10: "Myrmidon", 15: "Champion", 20: "Warmistress"},
	},
	"mage": {
		Male:   map[int]string{1: "Apprentice", 5: "Mage", 10: "Wizard", 15: "Arch-wizard", 20: "Weavemaster"},
		Female: map[int]string{1: "Apprentice", 5: "Mage", 10: "Wizardress", 15: "Arch-wizardress", 20: "Weavemistress"},
	},
	"monk": {
		Male:   map[int]string{1: "Initiate", 5: "Disciple", 10: "Immaculate", 15: "Master", 20: "Philosopher"},
		Female: map[int]string{1: "Initiate", 5: "Disciple", 10: "Immaculate", 15: "Master", 20: "Philosopher"},
	},
	"paladin": {
		Male:   map[int]string{1: "Cavalier", 5: "Warder", 10: "Holy Warrior", 15: "Lord", 20: "Crusader"},
		Female: map[int]string{1: "Cavalier", 5: "Warder", 10: "Holy Warrior", 15: "Lady", 20: "Crusader"},
	},
	"ranger": {
		Male:   map[int]string{1: "Ranger", 5: "Scout", 10: "Pathfinder", 15: "Trailblazer", 20: "Waymaker"},
		Female: map[int]string{1: "Ranger", 5: "Scout", 10: "Pathfinder", 15: "Trailblazer", 20: "Waymaker"},
	},
	"thief": {
		Male:   map[int]string{1: "Rogue", 5: "Pickpocket", 10: "Nightblade", 15: "Shadow Warrior", 20: "Shadowmaster"},
		Female: map[int]string{1: "Rogue", 5: "Pickpocket", 10: "Nightblade", 15: "Shadow Warrior", 20: "Shadowmistress"},
	},
}

func ClassTitle(class int, gender string, tier int) string {
	if class == 99 {
		return "Deity"
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
	if gender == "m" {
		return ClassTitles[AvailableClasses[class]].Male[selectTier]
	} else {
		return ClassTitles[AvailableClasses[class]].Female[selectTier]
	}

}

var TextGender = map[string]string{"m": "male", "f": "female"}
var TextSubPronoun = map[string]string{"m": "he", "f": "she"}
var TextPosPronoun = map[string]string{"m": "his", "f": "her"}
var TextDescPronoun = map[string]string{"m": "him", "f": "her"}
