package main

import (
	"fmt"
	"os"
	"strconv"
)

func isValidPasswordV1(candidate int) bool {
	s := strconv.Itoa(candidate)

	adjacentDigitConditionSatisfied := false

	var previousDigit rune

	for i, digit := range s {
		if i > 0 {
			if digit == previousDigit {
				adjacentDigitConditionSatisfied = true
			} else if digit < previousDigit {
				return false
			}
		}
		previousDigit = digit
	}

	return adjacentDigitConditionSatisfied
}

func isValidPasswordV2(candidate int) bool {
	s := []byte(strconv.Itoa(candidate))

	currentGroup := []byte{s[0]}
	groups := [][]byte{}

	for _, digit := range s[1:] {
		if digit == currentGroup[0] {
			currentGroup = append(currentGroup, digit)
		} else {
			if digit < currentGroup[0] {
				return false
			}

			groups = append(groups, currentGroup)
			currentGroup = []byte{digit}
		}
	}

	groups = append(groups, currentGroup)

	for _, group := range groups {
		if len(group) == 2 {
			return true
		}
	}

	return false
}

func main() {
	start, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Couldn't parse " + os.Args[1])
	}
	end, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("Couldn't parse " + os.Args[2])
	}

	for candidate := start; candidate <= end; candidate++ {
		if isValidPasswordV2(candidate) {
			fmt.Println(candidate)
		}
	}
}
