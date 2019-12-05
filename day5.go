package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

func readInteger(input *bufio.Reader) int {
	line, _, err := input.ReadLine()
	if err != nil {
		panic("couldn't read integer: " + err.Error())
	}
	i, err := strconv.Atoi(string(line))
	if err != nil {
		panic("couldn't read integer: " + err.Error())
	}
	return i
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

func runProgram(input *bufio.Reader, program []int) {
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
			program[program[ip+1]] = readInteger(input)
			ip += 2
		case 4:
			fmt.Println(getValue(program, ip+1, modes[1]))
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

func main() {
	programBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic("couldn't read program")
	}
	program := parseProgram(strings.TrimSpace(string(programBytes)))
	runProgram(bufio.NewReader(os.Stdin), program)
}
