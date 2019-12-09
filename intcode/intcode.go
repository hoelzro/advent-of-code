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

func getValue(program *[]int, ip, mode, relativeBase int) int {
	if mode == ModePosition {
		return (*program)[(*program)[ip]]
	} else if mode == ModeImmediate {
		return (*program)[ip]
	} else if mode == ModeRelative {
		return (*program)[(*program)[ip]+relativeBase]
	}
	panic("Unrecognized mode")
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
			program[program[ip+3]] = getValue(&program, ip+1, modes[1], relativeBase) + getValue(&program, ip+2, modes[2], relativeBase)
			ip += 4
		case 2:
			program[program[ip+3]] = getValue(&program, ip+1, modes[1], relativeBase) * getValue(&program, ip+2, modes[2], relativeBase)
			ip += 4
		case 3:
			var ok bool
			program[program[ip+1]], ok = <-input
			if !ok {
				break programLoop
			}
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
			if getValue(&program, ip+1, modes[1], relativeBase) < getValue(&program, ip+2, modes[2], relativeBase) {
				program[program[ip+3]] = 1
			} else {
				program[program[ip+3]] = 0
			}
			ip += 4
		case 8:
			if getValue(&program, ip+1, modes[1], relativeBase) == getValue(&program, ip+2, modes[2], relativeBase) {
				program[program[ip+3]] = 1
			} else {
				program[program[ip+3]] = 0
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
