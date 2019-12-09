package main

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"hoelz.ro/advent-of-code/2019/intcode"
)

func main() {
	programBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic("couldn't load program")
	}
	program := intcode.ParseProgram(strings.TrimSpace(string(programBytes)))

	// Part 1
	wg := &sync.WaitGroup{}
	r := intcode.ConstantReader(1)
	w := intcode.StdoutWriter(wg)
	intcode.RunProgram(r, w, program)
	close(w)
	wg.Wait()

	// Part 2
	wg = &sync.WaitGroup{}
	r = intcode.ConstantReader(2)
	w = intcode.StdoutWriter(wg)
	intcode.RunProgram(r, w, program)
	close(w)
	wg.Wait()
}
