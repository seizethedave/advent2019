package main

import (
	"fmt"

	"github.com/seizethedave/advent2019/advent05"
)

func main() {
	memory := []advent05.Word{3, 0, 4, 0, 99}

	err := advent05.Exec(memory)
	if err != nil {
		panic(err)
	}

	fmt.Println()
}
