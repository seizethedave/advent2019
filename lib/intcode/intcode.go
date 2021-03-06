package intcode

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Interpreter struct {
	Id           string
	relativeBase Address
	scanner      *bufio.Scanner
	InputStream  io.Reader
	OutputStream io.Writer
	InputFunc    func() Word
	OutputFunc   func(Word)
}

type Word int
type Address int

const (
	add           = Word(1)
	mul           = Word(2)
	input         = Word(3)
	output        = Word(4)
	jumpIfTrue    = Word(5)
	jumpIfFalse   = Word(6)
	lessThan      = Word(7)
	equals        = Word(8)
	adjustRelBase = Word(9)
	halt          = Word(99)

	positional = byte(0)
	immediate  = byte(1)
	relative   = byte(2)

	debug = false
)

type Header struct {
	opcode     Word
	paramModes []byte
}

func readHeader(header Word) Header {
	var modes []byte

	for flags := header / 100; flags > 0; flags /= 10 {
		modes = append(modes, byte(flags%10))
	}

	return Header{
		opcode:     header % 100,
		paramModes: modes,
	}
}

// opref consults the parameter mode definition and returns either the raw value
// v if it is an indirect value, or the value at mem[v] if it is positional.
func (h Header) opref(it *Interpreter, mem []Word, ptr, index Address) Address {
	mode := positional
	if int(index) < len(h.paramModes) {
		mode = h.paramModes[index]
	}
	switch mode {
	case immediate:
		return ptr + 1 + index
	case positional:
		return Address(mem[ptr+1+index])
	case relative:
		return it.relativeBase + Address(mem[ptr+1+index])
	default:
		panic("unknown param mode")
	}
}

type Op struct {
	header Word
}

type SimpleOp struct {
	Op
	fun func(*Interpreter, []Word, Address)
}

type BinaryOp struct {
	Op
	fun func(*Interpreter, []Word, Address, Address, Address)
}

type JumpIfTrueOp struct {
	Op
}

type JumpIfFalseOp struct {
	Op
}

type CompareOp struct {
	Op
	fun func(*Interpreter, []Word, Address, Address, Address)
}

type Runnable interface {
	Exec(it *Interpreter, mem []Word, ptr *Address, header Header)
}

func (op SimpleOp) Exec(it *Interpreter, mem []Word, ptr *Address, header Header) {
	if debug {
		fmt.Fprintln(os.Stderr, " >", it.Id, mem[*ptr:*ptr+2], "@", it.relativeBase)
	}

	op.fun(it, mem, header.opref(it, mem, *ptr, 0))
	*ptr += 2
}

func (op BinaryOp) Exec(it *Interpreter, mem []Word, ptr *Address, header Header) {
	if debug {
		fmt.Fprintln(os.Stderr, " >", it.Id, mem[*ptr:*ptr+4], it.relativeBase)
	}

	op.fun(it, mem,
		header.opref(it, mem, *ptr, 0),
		header.opref(it, mem, *ptr, 1),
		header.opref(it, mem, *ptr, 2),
	)
	*ptr += 4
}

func (op JumpIfTrueOp) Exec(it *Interpreter, mem []Word, ptr *Address, header Header) {
	if debug {
		fmt.Fprintln(os.Stderr, " >", it.Id, mem[*ptr:*ptr+3], it.relativeBase)
	}
	if mem[header.opref(it, mem, *ptr, 0)] != 0 {
		*ptr = Address(mem[header.opref(it, mem, *ptr, 1)])
	} else {
		*ptr += 3
	}
}

func (op JumpIfFalseOp) Exec(it *Interpreter, mem []Word, ptr *Address, header Header) {
	if debug {
		fmt.Fprintln(os.Stderr, " >", it.Id, mem[*ptr:*ptr+3], it.relativeBase)
	}
	if mem[header.opref(it, mem, *ptr, 0)] == 0 {
		*ptr = Address(mem[header.opref(it, mem, *ptr, 1)])
	} else {
		*ptr += 3
	}
}

func (op CompareOp) Exec(it *Interpreter, mem []Word, ptr *Address, header Header) {
	if debug {
		fmt.Fprintln(os.Stderr, " >", it.Id, mem[*ptr:*ptr+4], it.relativeBase)
	}

	op.fun(it, mem,
		header.opref(it, mem, *ptr, 0),
		header.opref(it, mem, *ptr, 1),
		header.opref(it, mem, *ptr, 2),
	)
	*ptr += 4
}

func opAdd(it *Interpreter, mem []Word, lhs, rhs, out Address) {
	mem[out] = mem[lhs] + mem[rhs]
}

func opMul(it *Interpreter, mem []Word, lhs, rhs, out Address) {
	mem[out] = mem[lhs] * mem[rhs]
}

func opInput(it *Interpreter, mem []Word, operand Address) {
	mem[operand] = it.InputFunc()
}

func opOutput(it *Interpreter, mem []Word, operand Address) {
	it.OutputFunc(mem[operand])
}

func opLessThan(it *Interpreter, mem []Word, lhs, rhs, out Address) {
	if mem[lhs] < mem[rhs] {
		mem[out] = 1
	} else {
		mem[out] = 0
	}
}

func opEquals(it *Interpreter, mem []Word, lhs, rhs, out Address) {
	if mem[lhs] == mem[rhs] {
		mem[out] = 1
	} else {
		mem[out] = 0
	}
}

func opAdjustRelBase(it *Interpreter, mem []Word, operand Address) {
	it.relativeBase += Address(mem[operand])
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
	adjustRelBase: SimpleOp{
		fun: opAdjustRelBase,
	},
}

func (it *Interpreter) ExecOp(mem []Word, ptr *Address) error {
	header := readHeader(mem[*ptr])

	op, ok := ops[header.opcode]
	if !ok {
		return fmt.Errorf("invalid opcode %v", header.opcode)
	}

	op.Exec(it, mem, ptr, header)
	return nil
}

// Exec runs a program.
func Exec(memory []Word) error {
	it := &Interpreter{
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
	}
	return it.Exec(memory)
}

func (it *Interpreter) ScanWord() Word {
	if !it.scanner.Scan() {
		panic("no scan")
	}
	v, err := strconv.Atoi(it.scanner.Text())
	if err != nil {
		panic(err)
	}
	return Word(v)
}

func (it *Interpreter) Exec(memory []Word) error {
	if it.InputFunc == nil {
		it.InputFunc = func() Word {
			return it.ScanWord()
		}
	}

	if it.OutputFunc == nil {
		it.OutputFunc = func(value Word) {
			fmt.Fprintln(it.OutputStream, value)
		}
	}

	it.scanner = bufio.NewScanner(it.InputStream)

	for address := Address(0); memory[address] != halt; {
		err := it.ExecOp(memory, &address)
		if err != nil {
			return err
		}
	}

	return nil
}
