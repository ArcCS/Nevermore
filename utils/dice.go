// Definition for dice rolls

package utils

import (
	"math/rand"
	"sort"
	"time"
)

func RollMax(dieSides int, numDice int, mod int) int {
	return (dieSides * numDice) + mod
}

func RollMin(numDice int, mod int) int {
	return numDice + mod
}

func Roll(dieSides int, numDice int, mod int) int {
	return DiceRoll(dieSides, numDice, mod, 0, true)[0]
}

func DiceRoll(dieSides int, numDice int, mod int, drop int, total bool) []int {
	rolls := make([]int, numDice)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := range rolls {
		rolls[i] = r1.Intn(dieSides) + 1
	}

	if drop > 0 {
		sort.Sort(sort.Reverse(sort.IntSlice(rolls)))
		rolls = rolls[:len(rolls)-drop]
	} else if total {
		// If we're asking for the total, just send it back as the first element
		rolls[0] = Sum(rolls) + mod
	}
	return rolls
}
