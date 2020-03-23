package utils

import "strings"

func Sum(input []int) int {
	sum := 0
	for i := range input {
		sum += input[i]
	}
	return sum
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
	return -1    //not found.
}



func Btoi(val bool) int {
	if val{
		return 1
	}else{
		return 0
	}
}