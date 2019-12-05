package main

import (
	"bufio"
	"fmt"
	"os"

	"hoelz.ro/advent-of-code/2019/intcode"
)

func main() {
	program := intcode.ParseProgram(os.Args[1])

	pristineProgram := make([]int, len(program))
	copy(pristineProgram, program)

	// Part 1
	// 1202 step
	program[1] = 12
	program[2] = 2

	programInput := bufio.NewReader(os.Stdin)

	intcode.RunProgram(programInput, program)

	fmt.Println(program[0])

	// Part 2
outerLoop:
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			copy(program, pristineProgram)
			program[1] = noun
			program[2] = verb
			intcode.RunProgram(programInput, program)

			if program[0] == 19690720 {
				fmt.Println(100*noun + verb)
				break outerLoop
			}
		}
	}
}
