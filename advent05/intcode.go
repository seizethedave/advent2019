package advent05

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Word int64
type Address int64

const (
	add    = Word(1)
	mul    = Word(2)
	input  = Word(3)
	output = Word(4)
	halt   = Word(99)

	instructionLength = Address(4)
)

const (
	paramModePosition  = 0
	paramModeImmediate = 1
)

func paramMode(header Word, index Address) int {
	modeDigits := header / 100

	for i := Address(0); i < index; i++ {
		modeDigits /= 10
	}

	return int(modeDigits % 10)
}

func opref(ptr, index Address, header Word, mem []Word) Address {
	if paramMode(header, index) == paramModePosition {
		return Address(mem[ptr+index])
	} else {
		return ptr + index
	}
}

type Op struct {
	header Word
}

type SimpleInputOp struct {
	Op
	fun func([]Word, Address)
}

type SimpleOutputOp struct {
	Op
	fun func([]Word, Address)
}

type BinaryOp struct {
	Op
	fun func([]Word, Address, Address, Address)
}

type Runnable interface {
	Exec(mem []Word, ptr Address) (Address, error)
}

func (op SimpleInputOp) Exec(mem []Word, ptr Address) (Address, error) {
	// ptr+1 contains the address of where to write the input.
	op.fun(mem, Address(mem[ptr+1]))
	return ptr + 2, nil
}

func (op SimpleOutputOp) Exec(mem []Word, ptr Address) (Address, error) {
	header := mem[ptr]
	op.fun(mem, opref(ptr, 1, header, mem))
	return ptr + 2, nil
}

func (op BinaryOp) Exec(mem []Word, ptr Address) (Address, error) {
	header := mem[ptr]
	op.fun(mem, opref(ptr, 1, header, mem), opref(ptr, 2, header, mem), ptr+3)
	return ptr + 4, nil
}

func ExecOp(mem []Word, ptr Address) (Address, Runnable, error) {
	opcode := mem[ptr] % 100
	op, ok := ops[opcode]
	if !ok {
		return 0, op, fmt.Errorf("invalid opcode %v", opcode)
	}

	ptr, err := op.Exec(mem, ptr)
	if err != nil {
		return 0, op, err
	}

	return ptr, op, nil
}

func opAdd(mem []Word, lhs, rhs Address, out Address) {
	mem[mem[out]] = mem[lhs] + mem[rhs]
}

func opMul(mem []Word, lhs, rhs Address, out Address) {
	mem[mem[out]] = mem[lhs] * mem[rhs]
}

func opInput(mem []Word, operand Address) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter input: ")
	s, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		panic(err)
	}
	mem[mem[operand]] = Word(v)
}

func opOutput(mem []Word, operand Address) {
	fmt.Print(operand)
}

var ops = map[Word]Runnable{
	add: BinaryOp{
		fun: opAdd,
	},
	mul: BinaryOp{
		fun: opMul,
	},
	input: SimpleInputOp{
		fun: opInput,
	},
	output: SimpleOutputOp{
		fun: opOutput,
	},
}

// Exec runs a program.
func Exec(memory []Word) error {
	for address := Address(0); memory[address] != halt; {
		addr, _, err := ExecOp(memory, address)
		if err != nil {
			return err
		}

		address = addr
	}

	return nil
}
