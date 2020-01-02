package main

import (
  "fmt"

  "github.com/seizethedave/advent2019/advent06"
)

func countOrbits(body *advent06.Body, distance int) int {
	count := 0
	for _, orbiter := range body.Orbiters {
		count += 1 + distance + countOrbits(orbiter, distance+1)
	}
	return count
}

func main() {
	rootBody := advent06.ParseBodyTree(advent06.Input)
	fmt.Println(countOrbits(rootBody, 0))
}
