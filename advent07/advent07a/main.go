package main

import (
	"fmt"

	"github.com/seizethedave/advent2019/lib/intcode"
)

func multiplexResults(results []intcode.Word) func() intcode.Word {
	calls := 0
	return func() intcode.Word {
		defer func() { calls++ }()
		return results[calls]
	}
}

func visitPermutations(prefix, items []intcode.Word, callback func([]intcode.Word)) {
	if len(items) == 0 {
		callback(prefix)
		return
	}

	for i, atom := range items {
		// append modifies the shared items/prefix backing arrays, so we gotta copy 'em.
		subPrefix := make([]intcode.Word, len(prefix))
		copy(subPrefix, prefix)
		subItems := make([]intcode.Word, len(items))
		copy(subItems, items)
		visitPermutations(append(subPrefix, atom), append(subItems[:i], subItems[i+1:]...), callback)
	}
}

func main() {
	code := []intcode.Word{3, 8, 1001, 8, 10, 8, 105, 1, 0, 0, 21, 38, 59, 76, 89, 106, 187, 268, 349, 430, 99999, 3, 9, 1002, 9, 3, 9, 101, 2, 9, 9, 1002, 9, 4, 9, 4, 9, 99, 3, 9, 1001, 9, 5, 9, 1002, 9, 5, 9, 1001, 9, 2, 9, 1002, 9, 3, 9, 4, 9, 99, 3, 9, 1001, 9, 4, 9, 102, 4, 9, 9, 1001, 9, 3, 9, 4, 9, 99, 3, 9, 101, 4, 9, 9, 1002, 9, 5, 9, 4, 9, 99, 3, 9, 1002, 9, 3, 9, 101, 5, 9, 9, 1002, 9, 3, 9, 4, 9, 99, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 99, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 99, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 99, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 99, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 99}

	bestScore := intcode.Word(0)

	visitPermutations([]intcode.Word{}, []intcode.Word{0, 1, 2, 3, 4}, func(perm []intcode.Word) {
		var output intcode.Word

		for i, phase := range perm {
			var signal intcode.Word

			if i == 0 {
				signal = 0
			} else {
				signal = output
			}

			interp := &intcode.Interpreter{
				IO: &intcode.InterpreterIO{
					InputFunc: multiplexResults([]intcode.Word{phase, signal}),
					OutputFunc: func(value intcode.Word) {
						output = value
					},
				},
			}
			interp.Exec(code)
		}

		if output > bestScore {
			bestScore = output
		}
	})

	fmt.Println(bestScore)
}
