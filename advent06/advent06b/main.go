package main

import (
	"fmt"

	"github.com/seizethedave/advent2019/advent06"
)

func commonAncestorDistance(node1, node2 *advent06.Body) (int, int) {
	distance := make(map[string]int)
	dist1 := 0
	dist2 := 0

	for {
		dist := 0

		if node1.Orbiting != nil {
			// haven't made it to the root yet.
			node1 = node1.Orbiting
			dist = distance[node1.Id]
			if dist != 0 {
				// an ancestor of node2.
				return dist1, dist
			}
			distance[node1.Id] = dist1
			dist1++
		}

		// repeat for node2.

		if node2.Orbiting != nil {
			// haven't made it to the root yet.
			node2 = node2.Orbiting
			dist = distance[node2.Id]
			if dist != 0 {
				// an ancestor of node1.
				return dist, dist2
			}
			distance[node2.Id] = dist2
			dist2++
		}
	}

	panic("no common ancestor found")
}

func main() {
	_, bodyMap := advent06.ParseBodyTree(advent06.Input)
	d1, d2 := commonAncestorDistance(bodyMap["SAN"], bodyMap["YOU"])
	fmt.Println(d1 + d2)
}
