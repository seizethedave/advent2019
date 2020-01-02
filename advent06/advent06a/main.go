package main

import (
  "fmt"

  "github.com/seizethedave/advent2019/advent06"
)

func main() {
	rootBody := advent06.ParseBodyTree(advent06.Input)
	fmt.Println(advent06.CountOrbits(rootBody, 0))
}
