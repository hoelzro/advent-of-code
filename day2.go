package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseProgram(program string) []int {
	var result []int

	for _, opcodeStr := range strings.Split(program, ",") {
		opcode, err := strconv.Atoi(opcodeStr)

		if err != nil {
			panic(err.Error())
		}

		result = append(result, opcode)
	}

	return result
}

func runProgram(program []int) {
	ip := 0

programLoop:
	for {
		opcode := program[ip]

		switch opcode {
		case 1:
			program[program[ip+3]] = program[program[ip+1]] + program[program[ip+2]]
			ip += 4
		case 2:
			program[program[ip+3]] = program[program[ip+1]] * program[program[ip+2]]
			ip += 4
		case 99:
			break programLoop
		}
	}

}

func main() {
	program := parseProgram(os.Args[1])

	pristineProgram := make([]int, len(program))
	copy(pristineProgram, program)

	// Part 1
	// 1202 step
	program[1] = 12
	program[2] = 2

	runProgram(program)

	fmt.Println(program[0])

	// Part 2
outerLoop:
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			copy(program, pristineProgram)
			program[1] = noun
			program[2] = verb
			runProgram(program)

			if program[0] == 19690720 {
				fmt.Println(100*noun + verb)
				break outerLoop
			}
		}
	}
}
