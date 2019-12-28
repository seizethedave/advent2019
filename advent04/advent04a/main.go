package main

import (
	"bytes"
	"fmt"

	"github.com/seizethedave/advent2019/advent04"
)

func hasRepeat(digits []byte) bool {
	for i := 0; i < len(digits)-1; i++ {
		if digits[i] == digits[i+1] {
			return true
		}
	}
	return false
}

func main() {
	// given was 124075, 124444 was next proper value.
	digits := []byte{1, 2, 4, 4, 4, 4}
	end := []byte{5, 8, 0, 7, 6, 9}
	num := 0

	for ; bytes.Compare(digits, end) <= 0; advent04.Incr(digits) {
		if hasRepeat(digits) {
			num++
		}
	}

	fmt.Println(num)
}
