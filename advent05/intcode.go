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
)

type Header struct {
	opcode    Word
	paramMask uint8
}

func readHeader(header Word) Header {
	opcode := header % 100
	flags := header / 100
	var mask uint8
	i := 0

	for flags != 0 {
		mask |= uint8(((flags % 10) & 1) << i)
		flags /= 10
		i++
	}

	return Header{
		opcode:    opcode,
		paramMask: mask,
	}
}

func (h Header) opref(ptr, index Address, mem []Word) Address {
	immediate := (h.paramMask & (1 << index)) != 0
	if immediate {
		return ptr + 1 + index
	} else {
		return Address(mem[ptr+1+index])
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
	Exec(mem []Word, ptr Address, header Header) (Address, error)
}

func (op SimpleInputOp) Exec(mem []Word, ptr Address, header Header) (Address, error) {
	// ptr+1 contains the address of where to write the input.
	op.fun(mem, Address(mem[ptr+1]))
	return ptr + 2, nil
}

func (op SimpleOutputOp) Exec(mem []Word, ptr Address, header Header) (Address, error) {
	op.fun(mem, header.opref(ptr, 0, mem))
	return ptr + 2, nil
}

func (op BinaryOp) Exec(mem []Word, ptr Address, header Header) (Address, error) {
	op.fun(mem, header.opref(ptr, 0, mem), header.opref(ptr, 1, mem), Address(mem[ptr+3]))
	return ptr + 4, nil
}

func ExecOp(mem []Word, ptr Address) (Address, Runnable, error) {
	header := readHeader(mem[ptr])

	op, ok := ops[header.opcode]
	if !ok {
		return 0, op, fmt.Errorf("invalid opcode %v", header.opcode)
	}

	ptr, err := op.Exec(mem, ptr, header)
	if err != nil {
		return 0, op, err
	}

	return ptr, op, nil
}

func opAdd(mem []Word, lhs, rhs, out Address) {
	mem[out] = mem[lhs] + mem[rhs]
}

func opMul(mem []Word, lhs, rhs, out Address) {
	mem[out] = mem[lhs] * mem[rhs]
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
	mem[operand] = Word(v)
}

func opOutput(mem []Word, operand Address) {
	fmt.Println(operand)
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
