package main

import (
	"bytes"
	"fmt"

	"github.com/seizethedave/advent2019/advent04"
)

func hasRunOfLength2(digits []byte) bool {
	length := 1
	prev := digits[0]

	// Any time the running digit changes, see if the run length was 2.

	for i := 1; i < len(digits); i++ {
		if digits[i] == prev {
			length++
		} else {
			if length == 2 {
				return true
			}
			length = 1
			prev = digits[i]
		}
	}
	return length == 2
}

func main() {
	// given was 124075, 124444 was next proper value.
	digits := []byte{1, 2, 4, 4, 4, 4}
	end := []byte{5, 8, 0, 7, 6, 9}
	num := 0

	for ; bytes.Compare(digits, end) <= 0; advent04.Incr(digits) {
		if hasRunOfLength2(digits) {
			num++
		}
	}

	fmt.Println(num)
}
