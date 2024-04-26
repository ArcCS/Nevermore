package utils

import (
	"bufio"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
)

func Sum(input []int) int {
	sum := 0
	for i := range input {
		sum += input[i]
	}
	return sum
}

func RemoveInt(slice []int, s int) []int {
	for i, v := range slice {
		if v == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func RandMapKeySelection(mapList map[string]int) string {
	keys := make([]string, 0, len(mapList))
	for k := range mapList {
		keys = append(keys, k)
	}
	if len(keys) > 0 {
		return keys[rand.Intn(len(keys))]
	}
	return ""
}

func RandListSelection(stringList []string) string {
	return stringList[rand.Intn(len(stringList))]
}

func IntIn(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringIn(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RankMapStringInt(values map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range values {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	ranked := make([]string, len(values))
	for i, kv := range ss {
		ranked[i] = kv.Key
	}
	return ranked
}

func StringInLike(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}

func StringInLower(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("error closing file", err)
		}
	}()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func Btoi(val bool) int {
	if val {
		return 1
	} else {
		return 0
	}
}

func WhereAt(subLoc int, charLoc int) string {
	// Moving backwards
	if subLoc == charLoc {
		return " next to you."
	}
	diff := math.Abs(float64(subLoc - charLoc))
	steps := ""
	direction := ""
	if subLoc > charLoc {
		direction = "in front of you."
	} else {
		direction = "behind you."
	}
	if diff == 1 {
		steps = " a couple steps "
	} else if diff == 2 {
		steps = " a dozen steps "
	} else if diff == 3 {
		steps = " many steps "
	} else if diff == 4 {
		steps = " at the other end of the room "
	}
	return steps + direction
}

func RandString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("rando error", err)
	}
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func Title(s string) string {
	caser := cases.Title(language.AmericanEnglish)
	return caser.String(s)
}
