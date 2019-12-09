package intcode

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
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
	ModeRelative  = 2
)

// XXX use a map instead?
func expandIfNeeded(program *[]int, index int) {
	if index > 1000000 {
		panic("I forbid you from growing beyond a megabyte")
	}
	if index >= len(*program) {
		newProgram := make([]int, index+1)
		copy(newProgram, *program)
		*program = newProgram
	}
}

func getValue(program *[]int, ip, mode, relativeBase int) int {
	if mode == ModePosition {
		expandIfNeeded(program, ip)
		expandIfNeeded(program, (*program)[ip])
		return (*program)[(*program)[ip]]
	} else if mode == ModeImmediate {
		expandIfNeeded(program, ip)
		return (*program)[ip]
	} else if mode == ModeRelative {
		expandIfNeeded(program, ip+relativeBase)
		return (*program)[(*program)[ip]+relativeBase]
	}
	panic("Unrecognized mode")
}

func setValue(program *[]int, ip, mode, relativeBase, value int) {
	if mode == ModePosition {
		expandIfNeeded(program, (*program)[ip])
		(*program)[(*program)[ip]] = value
	} else if mode == ModeRelative {
		expandIfNeeded(program, ip+relativeBase)
		(*program)[(*program)[ip]+relativeBase] = value
	} else {
		panic("Unrecognized mode")
	}
}

func RunProgram(input <-chan int, output chan<- int, originalProgram []int) {
	program := make([]int, len(originalProgram))
	copy(program, originalProgram)

	relativeBase := 0
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
			} else if mode%10 == 2 {
				modes[param] = ModeRelative
			}
			mode /= 10
			param++
		}

		switch opcode {
		case 1:
			expandIfNeeded(&program, ip+3)
			expandIfNeeded(&program, program[ip+3])
			setValue(&program, ip+3, modes[3], relativeBase,
				getValue(&program, ip+1, modes[1], relativeBase)+getValue(&program, ip+2, modes[2], relativeBase))
			ip += 4
		case 2:
			expandIfNeeded(&program, ip+3)
			expandIfNeeded(&program, program[ip+3])
			setValue(&program, ip+3, modes[3], relativeBase,
				getValue(&program, ip+1, modes[1], relativeBase)*getValue(&program, ip+2, modes[2], relativeBase))
			ip += 4
		case 3:
			expandIfNeeded(&program, ip+1)

			value, ok := <-input

			if !ok {
				break programLoop
			}

			setValue(&program, ip+1, modes[1], relativeBase, value)

			ip += 2
		case 4:
			output <- getValue(&program, ip+1, modes[1], relativeBase)
			ip += 2
		case 5:
			if getValue(&program, ip+1, modes[1], relativeBase) != 0 {
				ip = getValue(&program, ip+2, modes[2], relativeBase)
			} else {
				ip += 3
			}
		case 6:
			if getValue(&program, ip+1, modes[1], relativeBase) == 0 {
				ip = getValue(&program, ip+2, modes[2], relativeBase)
			} else {
				ip += 3
			}
		case 7:
			expandIfNeeded(&program, ip+3)
			expandIfNeeded(&program, program[ip+3])
			if getValue(&program, ip+1, modes[1], relativeBase) < getValue(&program, ip+2, modes[2], relativeBase) {
				setValue(&program, ip+3, modes[3], relativeBase, 1)
			} else {
				setValue(&program, ip+3, modes[3], relativeBase, 0)
			}
			ip += 4
		case 8:
			expandIfNeeded(&program, ip+3)
			expandIfNeeded(&program, program[ip+3])
			if getValue(&program, ip+1, modes[1], relativeBase) == getValue(&program, ip+2, modes[2], relativeBase) {
				setValue(&program, ip+3, modes[3], relativeBase, 1)
			} else {
				setValue(&program, ip+3, modes[3], relativeBase, 0)
			}
			ip += 4
		case 9:
			relativeBase += getValue(&program, ip+1, modes[1], relativeBase)
			ip += 2
		case 99:
			break programLoop
		default:
			panic("invalid opcode: " + strconv.Itoa(opcode))
		}
	}
}

func StdoutWriter(wg *sync.WaitGroup) chan<- int {
	wg.Add(1)
	ch := make(chan int, 10)

	go func() {
		for {
			value, ok := <-ch
			if !ok {
				break
			}
			fmt.Println(value)
		}
		wg.Done()
	}()

	return ch
}

func ConstantReader(value int) <-chan int {
	ch := make(chan int, 10)

	go func() {
		for {
			ch <- value
		}
	}()

	return ch
}
