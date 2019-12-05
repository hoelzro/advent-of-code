package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"hoelz.ro/advent-of-code/2019/intcode"
)

func main() {
	programBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic("couldn't read program")
	}
	program := intcode.ParseProgram(strings.TrimSpace(string(programBytes)))
	intcode.RunProgram(bufio.NewReader(os.Stdin), program)
}
