package main

import (
	"fmt"
	"strconv"
	"strings"
)

type word int64
type address int64

const (
	inputProgram = "1,0,0,3,1,1,2,3,1,3,4,3,1,5,0,3,2,13,1,19,1,19,10,23,1,23,13,27,1,6,27,31,1,9,31,35,2,10,35,39,1,39,6,43,1,6,43,47,2,13,47,51,1,51,6,55,2,6,55,59,2,59,6,63,2,63,13,67,1,5,67,71,2,9,71,75,1,5,75,79,1,5,79,83,1,83,6,87,1,87,6,91,1,91,5,95,2,10,95,99,1,5,99,103,1,10,103,107,1,107,9,111,2,111,10,115,1,115,9,119,1,13,119,123,1,123,9,127,1,5,127,131,2,13,131,135,1,9,135,139,1,2,139,143,1,13,143,0,99,2,0,14,0"
	target       = 19690720

	add               = word(1)
	mul               = word(2)
	halt              = word(99)
	instructionLength = address(4)
)

func opAdd(mem []word, in1, in2, out address) {
	mem[mem[out]] = mem[mem[in1]] + mem[mem[in2]]
}

func opMul(mem []word, in1, in2, out address) {
	mem[mem[out]] = mem[mem[in1]] * mem[mem[in2]]
}

func exec(memory []word) error {
	for address := address(0); memory[address] != halt; address += instructionLength {
		if memory[address] == add {
			opAdd(memory, address+1, address+2, address+3)
		} else if memory[address] == mul {
			opMul(memory, address+1, address+2, address+3)
		} else {
			return fmt.Errorf("invalid instruction %d at address %d", memory[address], address)
		}
	}

	return nil
}

func main() {
	var memory []word

	for _, s := range strings.Split(inputProgram, ",") {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}
		memory = append(memory, word(n))
	}

	for noun := word(0); noun <= 99; noun++ {
		for verb := word(0); verb <= 99; verb++ {
			trial := make([]word, len(memory))
			copy(trial, memory)
			trial[1] = noun
			trial[2] = verb
			err := exec(trial)
			if err != nil {
				panic(err)
			}

			if trial[0] == target {
				fmt.Println(100*noun + verb)
				return
			}
		}
	}
}
