package advent04

// Incr increments the given number (embodied by a slice of the digits)
//  to the next number that conforms to the "non-decreasing" rule.
func Incr(digits []byte) {
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
