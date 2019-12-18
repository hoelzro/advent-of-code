package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"hoelz.ro/advent-of-code/2019/intcode"
)

type vec2 struct {
	x int
	y int
}

const (
	Robot    = 1
	Space    = 2
	Scaffold = 3
)

func parseMap(mapInput string) map[vec2]int {
	lines := strings.Split(mapInput, "\n")

	result := make(map[vec2]int)

	for y, line := range lines {
		for x, ch := range line {
			pos := vec2{x, y}
			switch ch {
			case '.':
				result[pos] = Space
			case '^':
				fallthrough
			case '<':
				fallthrough
			case '>':
				fallthrough
			case 'v':
				result[pos] = Robot
			case '#':
				result[pos] = Scaffold
			default:
				panic("Unexpected character")
			}
		}
	}

	return result
}

func main() {
	programBytes, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		panic(err.Error())
	}

	program := intcode.ParseProgram(strings.TrimSpace(string(programBytes)))

	wg := &sync.WaitGroup{}

	wg.Add(1)

	// Part 1
	mapBuilder := &strings.Builder{}

	output := make(chan int)
	go func() {
		for {
			value, ok := <-output
			if !ok {
				break
			}
			fmt.Fprint(mapBuilder, string(value))
		}
		wg.Done()
	}()

	intcode.RunProgram(nil, output, program)
	close(output)

	wg.Wait()

	environs := parseMap(mapBuilder.String())

	intersections := []vec2{}

	for pos, feature := range environs {
		if feature == Scaffold {
			x := pos.x
			y := pos.y

			neighbors := []vec2{
				{x + 1, y},
				{x - 1, y},
				{x, y + 1},
				{x, y - 1},
			}

			allNeighborsScaffold := true

			for _, neighbor := range neighbors {
				if environs[neighbor] != Scaffold {
					allNeighborsScaffold = false
					break
				}
			}

			if allNeighborsScaffold {
				intersections = append(intersections, pos)
			}
		}
	}

	alignmentParameterSum := 0

	for _, intersection := range intersections {
		alignmentParameterSum += intersection.x * intersection.y
	}

	fmt.Println(alignmentParameterSum)
}
