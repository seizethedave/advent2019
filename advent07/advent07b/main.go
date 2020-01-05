package main

import (
	"fmt"
	"io"
	"sync"

	"github.com/seizethedave/advent2019/lib/intcode"
)

func visitPermutations(prefix, items []intcode.Word, callback func([]intcode.Word)) {
	if len(items) == 0 {
		callback(prefix)
		return
	}

	for i, atom := range items {
		// append modifies the shared items/prefix backing arrays, so we gotta copy 'em.
		subPrefix := make([]intcode.Word, len(prefix))
		copy(subPrefix, prefix)
		subItems := make([]intcode.Word, len(items))
		copy(subItems, items)
		visitPermutations(append(subPrefix, atom), append(subItems[:i], subItems[i+1:]...), callback)
	}
}

func execInGroup(wg *sync.WaitGroup, interp *intcode.Interpreter, code []intcode.Word, signalDone bool) {
	myCode := make([]intcode.Word, len(code))
	copy(myCode, code)
	if err := interp.Exec(myCode); err != nil {
		panic(err)
	}
	if signalDone {
		wg.Done()
	}
}

func main() {
	code := []intcode.Word{3, 8, 1001, 8, 10, 8, 105, 1, 0, 0, 21, 38, 59, 76, 89, 106, 187, 268, 349, 430, 99999, 3, 9, 1002, 9, 3, 9, 101, 2, 9, 9, 1002, 9, 4, 9, 4, 9, 99, 3, 9, 1001, 9, 5, 9, 1002, 9, 5, 9, 1001, 9, 2, 9, 1002, 9, 3, 9, 4, 9, 99, 3, 9, 1001, 9, 4, 9, 102, 4, 9, 9, 1001, 9, 3, 9, 4, 9, 99, 3, 9, 101, 4, 9, 9, 1002, 9, 5, 9, 4, 9, 99, 3, 9, 1002, 9, 3, 9, 101, 5, 9, 9, 1002, 9, 3, 9, 4, 9, 99, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 99, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 99, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 99, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 99, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 99}
	bestScore := intcode.Word(0)

	visitPermutations([]intcode.Word{}, []intcode.Word{5, 6, 7, 8, 9}, func(perm []intcode.Word) {
		rB, wA := io.Pipe()
		rC, wB := io.Pipe()
		rD, wC := io.Pipe()
		rE, wD := io.Pipe()
		rA, wE := io.Pipe()

		interpA := &intcode.Interpreter{
			Id:           "A",
			InputStream:  rA,
			OutputStream: wA,
		}
		interpB := &intcode.Interpreter{
			Id:           "B",
			InputStream:  rB,
			OutputStream: wB,
		}
		interpC := &intcode.Interpreter{
			Id:           "C",
			InputStream:  rC,
			OutputStream: wC,
		}
		interpD := &intcode.Interpreter{
			Id:           "D",
			InputStream:  rD,
			OutputStream: wD,
		}
		interpE := &intcode.Interpreter{
			Id:           "E",
			InputStream:  rE,
			OutputStream: wE,
		}

		var wg sync.WaitGroup
		// Want to wake when all but the last amp has halted.
		wg.Add(4)
		go execInGroup(&wg, interpA, code, true)
		go execInGroup(&wg, interpB, code, true)
		go execInGroup(&wg, interpC, code, true)
		go execInGroup(&wg, interpD, code, true)
		go execInGroup(&wg, interpE, code, false)

		fmt.Fprintln(wE, perm[0])
		fmt.Fprintln(wA, perm[1])
		fmt.Fprintln(wB, perm[2])
		fmt.Fprintln(wC, perm[3])
		fmt.Fprintln(wD, perm[4])

		fmt.Fprintln(wE, 0)
		wg.Wait()

		if finalSignal := interpA.ScanWord(); finalSignal > bestScore {
			bestScore = finalSignal
		}
	})

	fmt.Println(bestScore)
}
