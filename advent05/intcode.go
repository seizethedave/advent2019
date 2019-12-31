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

	debug = false
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
		mask |= uint8((flags & 1) << i)
		flags /= 10
		i++
	}

	return Header{
		opcode:    opcode,
		paramMask: mask,
	}
}

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
		err := ExecOp(memory, &address)
		if err != nil {
			return err
		}
	}

	return nil
}
