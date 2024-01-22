package config

var LevelCap = 25
var SkillCap = 45000000

// TierExpLevels Leveling Values
var TierExpLevels = map[int]int{
	2:  450,
    	3:  1500, // 1050
    	4:  5000, // 3500
    	5:  10000, // 5000
    	6:  25000, // 15000
    	7:  125000, // 100K
    	8:  350000, // 225K
    	9:  750000, // 400K
    	10: 1450000, // 700K
    	11: 2550000, // 1.1m
    	12: 4050000, // 1.5m
    	13: 6250000, // 2.2m
    	14: 9500000, // 3.25m
    	15: 14200000, // 4.7m
    	16: 19200000, // 5m
    	17: 25200000, // 6m
    	18: 32800000, // 7.6m
    	19: 41600000, // 8.8m
    	20: 51000000, // 9.4m
    	21: 70000000, // 19m
    	22: 95000000, // 25m
    	23: 122000000, // 27m
	24: 152800000, // 30.8m
	25: 185000000, // 32.2m
}

func MaxLoss(tier int) int {
	return TierExpLevels[tier+1] - TierExpLevels[tier]
}

/* Old Values
// TierExpLevels Leveling Values
var TierExpLevels = map[int]int{
	2:  650,
	3:  2800,
	4:  7500,
	5:  15000,
	6:  30000,
	7:  84000,
	8:  250000,
	9:  550000,
	10: 1100000,
	11: 1900000,
	12: 3000000,
	13: 4500000,
	14: 6500000,
	15: 8600000,
	16: 11200000,
	17: 14650000,
	18: 19400000,
	19: 25300000,
	20: 32500000,
	21: 41100000,
	22: 51200000,
	23: 63100000,
	24: 76900000,
	25: 92500000,
}
*/

var GoldPerLevel = map[int]int{
	2:  100,
	3:  333,
	4:  666,
	5:  1333,
	6:  2666,
	7:  5333,
	8:  10666,
	9:  21333,
	10: 42666,
	11: 84333,
	12: 142666,
	13: 210333,
	14: 250666,
	15: 275333,
	16: 300666,
	17: 345333,
	18: 480666,
	19: 596333,
	20: 723666,
	21: 876333,
	22: 1010666,
	23: 1212333,
	24: 1433666,
	25: 1566333,
}

var TextTiers = []string{
	"zero",
	"first",
	"second",
	"third",
	"fourth",
	"fifth",
	"sixth",
	"seventh",
	"eighth",
	"ninth",
	"tenth",
	"eleventh",
	"twelfth",
	"thirteenth",
	"fourteenth",
	"fifteenth",
	"sixteenth",
	"seventeenth",
	"eighteenth",
	"nineteenth",
	"twentieth",
	"twenty-first",
	"twenty-second",
	"twenty-third",
	"twenty-fourth",
	"twenty-fifth",
	"twenty-sixth",
	"twenty-seventh",
	"twenty-eighth",
	"twenty-ninth",
	"thirtieth",
	"thirty-first",
	"thirty-second",
	"thirty-third",
	"thirty-fourth",
	"thirty-fifth",
	"thirty-sixth",
	"thirty-seventh",
	"thirty-eighth",
	"thirty-ninth",
	"fortieth",
	"forty-first",
	"forty-second",
	"forty-third",
	"forty-fourth",
	"forty-fifth",
	"forty-sixth",
	"forty-seventh",
	"forty-eighth",
	"forty-ninth",
	"fiftieth",
}

var PrintNumbers = []string{
	"0",
	"1st",
	"2nd",
	"3rd",
	"4th",
	"5th",
	"6th",
	"7th",
	"8th",
	"9th",
	"10th",
	"11th",
	"12th",
	"13th",
	"14th",
	"15th",
	"16th",
	"17th",
	"18th",
	"19th",
	"20th",
	"21st",
	"22nd",
	"23rd",
	"24th",
	"25th",
	"26th",
	"27th",
	"28th",
	"29th",
	"30th",
	"31st",
	"32nd",
	"33rd",
	"34th",
	"35th",
	"36th",
	"37th",
	"38th",
	"39th",
	"40th",
}

var TextNumbers = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
	"ten",
	"eleven",
	"twelve",
	"thirteen",
	"fourteen",
	"fifteen",
	"sixteen",
	"seventeen",
	"eighteen",
	"nineteen",
	"twenty",
	"twenty-one",
	"twenty-two",
	"twenty-three",
	"twenty-four",
	"twenty-five",
	"twenty-six",
	"twenty-seven",
	"twenty-eight",
	"twenty-nine",
	"thirty",
	"thirty-one",
	"thirty-two",
	"thirty-three",
	"thirty-four",
	"thirty-five",
	"thirty-six",
	"thirty-seven",
	"thirty-eight",
	"thirty-nine",
	"forty",
	"forty-one",
	"forty-two",
	"forty-three",
	"forty-four",
	"forty-five",
	"forty-six",
	"forty-seven",
	"forty-eight",
	"forty-nine",
	"fifty",
}
