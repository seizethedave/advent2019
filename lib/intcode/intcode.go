package intcode

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Word int
type Address int

const (
	add         = Word(1)
	mul         = Word(2)
	input       = Word(3)
	output      = Word(4)
	jumpIfTrue  = Word(5)
	jumpIfFalse = Word(6)
	lessThan    = Word(7)
	equals      = Word(8)
	halt        = Word(99)

	debug = false
)

type Header struct {
	opcode    Word
	paramMask byte
}

func readHeader(header Word) Header {
	var mask byte
	i := 0

	for flags := header / 100; flags > 0; flags /= 10 {
		mask |= byte((flags & 1) << i)
		i++
	}

	return Header{
		opcode:    header % 100,
		paramMask: mask,
	}
}

// opref consults the parameter mode definition and returns either the raw value
// v if it is an indirect value, or the value at mem[v] if it is positional.
func (h Header) opref(mem []Word, ptr, index Address) Address {
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

type SimpleOp struct {
	Op
	fun func([]Word, Address)
}

type BinaryOp struct {
	Op
	fun func([]Word, Address, Address, Address)
}

type JumpIfTrueOp struct {
	Op
}

type JumpIfFalseOp struct {
	Op
}

type CompareOp struct {
	Op
	fun func([]Word, Address, Address, Address)
}

type Runnable interface {
	Exec(mem []Word, ptr *Address, header Header)
}

func (op SimpleOp) Exec(mem []Word, ptr *Address, header Header) {
	if debug {
		fmt.Fprintln(os.Stderr, " >", mem[*ptr:*ptr+2])
	}

	op.fun(mem, header.opref(mem, *ptr, 0))
	*ptr += 2
}

func (op BinaryOp) Exec(mem []Word, ptr *Address, header Header) {
	if debug {
		fmt.Fprintln(os.Stderr, " >", mem[*ptr:*ptr+4])
	}

	op.fun(mem, header.opref(mem, *ptr, 0), header.opref(mem, *ptr, 1), Address(mem[*ptr+3]))
	*ptr += 4
}

func (op JumpIfTrueOp) Exec(mem []Word, ptr *Address, header Header) {
	if debug {
		fmt.Fprintln(os.Stderr, " >", mem[*ptr:*ptr+3])
	}
	if mem[header.opref(mem, *ptr, 0)] != 0 {
		*ptr = Address(mem[header.opref(mem, *ptr, 1)])
	} else {
		*ptr += 3
	}
}

func (op JumpIfFalseOp) Exec(mem []Word, ptr *Address, header Header) {
	if debug {
		fmt.Fprintln(os.Stderr, " >", mem[*ptr:*ptr+3])
	}
	if mem[header.opref(mem, *ptr, 0)] == 0 {
		*ptr = Address(mem[header.opref(mem, *ptr, 1)])
	} else {
		*ptr += 3
	}
}

func (op CompareOp) Exec(mem []Word, ptr *Address, header Header) {
	if debug {
		fmt.Fprintln(os.Stderr, " >", mem[*ptr:*ptr+4])
	}

	op.fun(mem, header.opref(mem, *ptr, 0), header.opref(mem, *ptr, 1), Address(mem[*ptr+3]))
	*ptr += 4
}

func ExecOp(mem []Word, ptr *Address) error {
	header := readHeader(mem[*ptr])

	op, ok := ops[header.opcode]
	if !ok {
		return fmt.Errorf("invalid opcode %v", header.opcode)
	}

	op.Exec(mem, ptr, header)
	return nil
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
	fmt.Println(mem[operand])
}

func opLessThan(mem []Word, lhs, rhs, out Address) {
	if mem[lhs] < mem[rhs] {
		mem[out] = 1
	} else {
		mem[out] = 0
	}
}

func opEquals(mem []Word, lhs, rhs, out Address) {
	if mem[lhs] == mem[rhs] {
		mem[out] = 1
	} else {
		mem[out] = 0
	}
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
	jumpIfTrue:  JumpIfTrueOp{},
	jumpIfFalse: JumpIfFalseOp{},
	lessThan: CompareOp{
		fun: opLessThan,
	},
	equals: CompareOp{
		fun: opEquals,
	},
}

// Exec runs a program.
func Exec(memory []Word) error {
	for address := Address(0); memory[address] != halt; {
		err := ExecOp(memory, &address)
		if err != nil {
			return err
		}
	}

	return nil
}
