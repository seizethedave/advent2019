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

	black       = '0'
	white       = '1'
	transparent = '2'
)

func drawImage(img []byte) {
	for i, b := range img {
		if b == black {
			fmt.Print("@")
		} else {
			fmt.Print(" ")
		}

		if (i+1)%width == 0 {
			fmt.Println()
		}
	}
}

// To run:
// cat ../input.txt | go run main.go
func main() {
	reader := bufio.NewReader(os.Stdin)
	layerBuf := make([]byte, area)
	var layers [][]byte

	for {
		_, err := io.ReadFull(reader, layerBuf)
		if err == io.EOF {
			break
		}

		layer := make([]byte, area)
		copy(layer, layerBuf)
		layers = append(layers, layer)
	}

	// Walk them in LIFO order and render a final depiction.
	img := make([]byte, area)

	for i := len(layers) - 1; i >= 0; i-- {
		for j := 0; j < area; j++ {
			if layers[i][j] != transparent {
				img[j] = layers[i][j]
			}
		}
	}

	drawImage(img)
}
