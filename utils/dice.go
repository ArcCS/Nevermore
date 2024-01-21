// Definition for dice rolls

package utils

import (
	crand "crypto/rand"
	"log"
	"math"
	"math/big"
	"math/rand"
	"sort"
	"time"
)

var DieSeed = rand.NewSource(CryptoRandSecure(math.MaxInt))
var generator = rand.New(DieSeed)
var LastUpdate = time.Now().Unix()

func RollMax(dieSides int, numDice int, mod int) int {
	return (dieSides * numDice) + mod
}

func CryptoRandSecure(max int64) int64 {
	nBig, err := crand.Int(crand.Reader, big.NewInt(max))
	if err != nil {
		log.Println(err)
		return time.Now().UnixNano() // fall back to using unix time as a seed
	}
	return nBig.Int64()
}

func RollMin(numDice int, mod int) int {
	return numDice + mod
}

func Roll(dieSides int, numDice int, mod int) int {
	return DiceRoll(dieSides, numDice, mod, 0, true)[0]
}

func DiceRoll(dieSides int, numDice int, mod int, drop int, total bool) []int {
	rolls := make([]int, numDice)

	if time.Now().Unix()-LastUpdate >= 600 {
		DieSeed = rand.NewSource(CryptoRandSecure(math.MaxInt))
		generator = rand.New(DieSeed)
		LastUpdate = time.Now().Unix()
	}

	for i := range rolls {
		rolls[i] = generator.Intn(dieSides) + 1
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
