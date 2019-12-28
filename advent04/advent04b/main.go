package main

import (
	"bytes"
	"fmt"

	"github.com/seizethedave/advent2019/advent04"
)

// compute a list of "run lengths".
// If the input is [1 1 1 1 2 2], there is a run of 4 1's and 2 2's => [4 2].
// [5 6 7 7 7 7] => [1 1 4] etc.
func runLengths(digits []byte) []int {
	var lengths []int
	length := 1
	prev := digits[0]

	for i := 1; i < len(digits); i++ {
		if digits[i] == prev {
			length++
		} else {
			lengths = append(lengths, length)
			length = 1
			prev = digits[i]
		}
	}

	lengths = append(lengths, length)
	return lengths
}

func hasRunOfLength2(digits []byte) bool {
	for _, length := range runLengths(digits) {
		if length == 2 {
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
		if hasRunOfLength2(digits) {
			num++
		}
	}

	fmt.Println(num)
}
