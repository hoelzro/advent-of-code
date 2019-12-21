package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"hoelz.ro/advent-of-code/2019/intcode"
)

func affectedByTractorBeam(program []int, x, y int) bool {
	input := make(chan int, 2)
	input <- x
	input <- y

	output := make(chan int, 1)

	intcode.RunProgram(input, output, program)

	result := <-output

	close(input)
	close(output)

	return result == 1
}

func main() {
	programBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err.Error())
	}

	program := intcode.ParseProgram(strings.TrimSpace(string(programBytes)))

	count := 0

	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {
			if affectedByTractorBeam(program, x, y) {
				count++
			}
		}
	}

	fmt.Println(count)
}
