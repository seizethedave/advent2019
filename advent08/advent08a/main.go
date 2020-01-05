package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	width  = 25
	height = 6
	area   = width * height
)

// To run:
// cat ../input.txt | go run main.go
func main() {
	reader := bufio.NewReader(os.Stdin)

	i := 0
	zeros := 0
	ones := 0
	twos := 0

	fewestZeros := 0
	fewestSet := false
	result := 0

	for {
		r, _, err := reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		num := int(r - '0')

		switch num {
		case 0:
			zeros++
		case 1:
			ones++
		case 2:
			twos++
		}

		if i++; i >= area {
			// layer boundary.
			if !fewestSet || zeros < fewestZeros {
				fewestZeros = zeros
				result = ones * twos
				fewestSet = true
			}

			i = 0
			zeros = 0
			ones = 0
			twos = 0
		}
	}

	fmt.Println(result)
}
