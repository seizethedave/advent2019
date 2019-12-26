package intcode

import (
  "fmt"
)

type Word int64
type Address int64

const (
  add               = Word(1)
  mul               = Word(2)
  halt              = Word(99)
  instructionLength = Address(4)
)

func opAdd(mem []Word, in1, in2, out Address) {
	mem[mem[out]] = mem[mem[in1]] + mem[mem[in2]]
}

func opMul(mem []Word, in1, in2, out Address) {
	mem[mem[out]] = mem[mem[in1]] * mem[mem[in2]]
}

// Exec runs a program.
func Exec(memory []Word) error {
	for address := Address(0); memory[address] != halt; address += instructionLength {
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
