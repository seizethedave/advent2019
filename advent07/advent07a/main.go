package main

import (
	"fmt"

	"github.com/seizethedave/advent2019/lib/intcode"
)

const (
	searchMin = 0
	searchMax = 4
)

func multiplexResults(results []intcode.Word) func() intcode.Word {
	calls := 0
	return func() intcode.Word {
		defer func() { calls++ }()
		return results[calls]
	}
}

func main() {
	code := []intcode.Word{3, 8, 1001, 8, 10, 8, 105, 1, 0, 0, 21, 38, 59, 76, 89, 106, 187, 268, 349, 430, 99999, 3, 9, 1002, 9, 3, 9, 101, 2, 9, 9, 1002, 9, 4, 9, 4, 9, 99, 3, 9, 1001, 9, 5, 9, 1002, 9, 5, 9, 1001, 9, 2, 9, 1002, 9, 3, 9, 4, 9, 99, 3, 9, 1001, 9, 4, 9, 102, 4, 9, 9, 1001, 9, 3, 9, 4, 9, 99, 3, 9, 101, 4, 9, 9, 1002, 9, 5, 9, 4, 9, 99, 3, 9, 1002, 9, 3, 9, 101, 5, 9, 9, 1002, 9, 3, 9, 4, 9, 99, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 99, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 99, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 99, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 99, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 99}

	bestScore := intcode.Word(0)

	for m1 := searchMin; m1 <= searchMax; m1++ {
		var output1 intcode.Word

		interp1 := &intcode.Interpreter{
			IO: &intcode.InterpreterIO{
				InputFunc: multiplexResults([]intcode.Word{intcode.Word(m1), 0}),
				OutputFunc: func(value intcode.Word) {
					output1 = value
				},
			},
		}
		interp1.Exec(code)

		for m2 := searchMin; m2 <= searchMax; m2++ {
			var output2 intcode.Word
			interp2 := &intcode.Interpreter{
				IO: &intcode.InterpreterIO{
					InputFunc: multiplexResults([]intcode.Word{intcode.Word(m2), output1}),
					OutputFunc: func(value intcode.Word) {
						output2 = value
					},
				},
			}
			interp2.Exec(code)

			for m3 := searchMin; m3 <= searchMax; m3++ {
				var output3 intcode.Word
				interp3 := &intcode.Interpreter{
					IO: &intcode.InterpreterIO{
						InputFunc: multiplexResults([]intcode.Word{intcode.Word(m3), output2}),
						OutputFunc: func(value intcode.Word) {
							output3 = value
						},
					},
				}
				interp3.Exec(code)

				for m4 := searchMin; m4 <= searchMax; m4++ {
					var output4 intcode.Word
					interp4 := &intcode.Interpreter{
						IO: &intcode.InterpreterIO{
							InputFunc: multiplexResults([]intcode.Word{intcode.Word(m4), output3}),
							OutputFunc: func(value intcode.Word) {
								output4 = value
							},
						},
					}
					interp4.Exec(code)

					for m5 := searchMin; m5 <= searchMax; m5++ {
						var output5 intcode.Word
						interp5 := &intcode.Interpreter{
							IO: &intcode.InterpreterIO{
								InputFunc: multiplexResults([]intcode.Word{intcode.Word(m5), output4}),
								OutputFunc: func(value intcode.Word) {
									output5 = value
								},
							},
						}
						interp5.Exec(code)
						fmt.Println("got output from amp 5", output5, m1, m2, m3, m4, m5)

						if output5 > bestScore {
							bestScore = output5
						}
					}
				}
			}
		}
	}

	fmt.Println(bestScore)
}
