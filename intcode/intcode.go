package intcode

import (
	"strconv"
	"strings"
)

func ParseProgram(program string) []int {
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

const (
	ModePosition  = 0
	ModeImmediate = 1
)

func getValue(program []int, ip, mode int) int {
	if mode == ModePosition {
		return program[program[ip]]
	} else if mode == ModeImmediate {
		return program[ip]
	}
	panic("Unrecognized mode")
}

func RunProgram(input <-chan int, output chan<- int, originalProgram []int) {
	program := make([]int, len(originalProgram))
	copy(program, originalProgram)

	ip := 0

programLoop:
	for {
		modes := make(map[int]int)
		opcode := program[ip] % 100
		mode := program[ip] / 100
		param := 1

		for mode != 0 {
			if mode%10 == 1 {
				modes[param] = ModeImmediate
			}
			mode /= 10
			param++
		}

		switch opcode {
		case 1:
			program[program[ip+3]] = getValue(program, ip+1, modes[1]) + getValue(program, ip+2, modes[2])
			ip += 4
		case 2:
			program[program[ip+3]] = getValue(program, ip+1, modes[1]) * getValue(program, ip+2, modes[2])
			ip += 4
		case 3:
			var ok bool
			program[program[ip+1]], ok = <-input
			if !ok {
				break programLoop
			}
			ip += 2
		case 4:
			output <- getValue(program, ip+1, modes[1])
			ip += 2
		case 5:
			if getValue(program, ip+1, modes[1]) != 0 {
				ip = getValue(program, ip+2, modes[2])
			} else {
				ip += 3
			}
		case 6:
			if getValue(program, ip+1, modes[1]) == 0 {
				ip = getValue(program, ip+2, modes[2])
			} else {
				ip += 3
			}
		case 7:
			if getValue(program, ip+1, modes[1]) < getValue(program, ip+2, modes[2]) {
				program[program[ip+3]] = 1
			} else {
				program[program[ip+3]] = 0
			}
			ip += 4
		case 8:
			if getValue(program, ip+1, modes[1]) == getValue(program, ip+2, modes[2]) {
				program[program[ip+3]] = 1
			} else {
				program[program[ip+3]] = 0
			}
			ip += 4
		case 99:
			break programLoop
		default:
			panic("invalid opcode: " + strconv.Itoa(opcode))
		}
	}

}
