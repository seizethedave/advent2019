package main

import (
	"bytes"
	"fmt"
)

// incr increments the given number (embodied by a slice of the digits)
//  to the next number that conforms to the "non-decreasing" rule.
func incr(digits []byte) {
	// increment
	place := len(digits) - 1
	digits[place]++

	// carry
	for ; place >= 0; place-- {
		if digits[place] == 10 {
			if place == 0 {
				return
			}
			digits[place] = 0
			digits[place-1]++
		} else {
			break
		}
	}

	// enforce non-decreasing rule.
	minDigit := digits[place]
  place++

	for ; place < len(digits); place++ {
		if digits[place] < minDigit {
			digits[place] = minDigit
		} else if digits[place] > minDigit {
			minDigit = digits[place]
		}
	}
}

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

	for ; bytes.Compare(digits, end) <= 0; incr(digits) {
    if hasRepeat(digits) {
			num++
		}
	}

	fmt.Println(num)
}
