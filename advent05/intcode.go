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

type Op struct {
	header Word
}

const (
	paramModePosition  = 0
	paramModeImmediate = 1
)

func (op *Op) paramMode(index int) int {
	modeDigits := op.header / 100

	for i := 0; i < index; i++ {
		modeDigits /= 10
	}

	return int(modeDigits % 10)
}

type SimpleOp struct {
	Op
	fun func([]Word, Address)
}

type UnaryOp struct {
	Op
	fun func([]Word, Address, Address)
}

type BinaryOp struct {
	Op
	fun func([]Word, Address, Address, Address)
}

type Runnable interface {
	Exec(mem []Word, ptr Address) (Address, error)
}

func (op SimpleOp) Exec(mem []Word, ptr Address) (Address, error) {
	op.header = mem[ptr]
	op.fun(mem, ptr+1)
	return ptr + 2, nil
}
func (op UnaryOp) Exec(mem []Word, ptr Address) (Address, error) {
	op.header = mem[ptr]
	op.fun(mem, ptr+1, ptr+2)
	return ptr + 3, nil
}
func (op BinaryOp) Exec(mem []Word, ptr Address) (Address, error) {
	op.header = mem[ptr]
	op.fun(mem, ptr+1, ptr+2, ptr+3)
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

func opAdd(mem []Word, in1, in2, out Address) {
	mem[mem[out]] = mem[mem[in1]] + mem[mem[in2]]
}

func opMul(mem []Word, in1, in2, out Address) {
	mem[mem[out]] = mem[mem[in1]] * mem[mem[in2]]
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
	fmt.Print(mem[mem[operand]])
}

var ops = map[Word]Runnable{
	add: BinaryOp{
		fun: opAdd,
	},
	mul: BinaryOp{
		fun: opMul,
	},
	input: SimpleOp{
		fun: opInput,
	},
	output: SimpleOp{
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
