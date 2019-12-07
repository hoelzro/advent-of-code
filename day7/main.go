package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"hoelz.ro/advent-of-code/2019/intcode"
)

func fillPermutations(values []int, prefix []int, result [][]int) [][]int {
	if len(values) == 0 {
		member := make([]int, len(prefix))
		copy(member, prefix)
		return append(result, member)
	}

	for i := range values {
		newPrefix := make([]int, len(prefix)+1)
		copy(newPrefix, prefix)
		newPrefix[len(prefix)] = values[i]

		newValues := make([]int, 0, len(values)-1)
		for j := range values {
			if j != i {
				newValues = append(newValues, values[j])
			}
		}

		result = fillPermutations(newValues, newPrefix, result)
	}

	return result
}

func permutations(values []int) [][]int {
	result := [][]int{}

	return fillPermutations(values, nil, result)
}

func runAmplifiers(program []int, phaseSettings []int) int {
	previousOutput := 0

	for _, phaseSetting := range phaseSettings {
		inputs := make(chan int, 2)
		outputs := make(chan int, 2)

		inputs <- phaseSetting
		inputs <- previousOutput

		intcode.RunProgram(inputs, outputs, program)
		previousOutput = <-outputs
	}

	return previousOutput
}

func part1(program []int) {
	var maxPermutation []int
	maxOutput := 0

	values := []int{0, 1, 2, 3, 4}
	for _, p := range permutations(values) {
		output := runAmplifiers(program, p)

		if output > maxOutput {
			maxPermutation = p
			maxOutput = output
		}
	}

	fmt.Println(maxOutput, maxPermutation)
}

func runAmplifiersFeedback(program []int, phaseSettings []int) int {
	inputs := make([]chan int, len(phaseSettings))
	outputs := make([]chan int, len(phaseSettings))

	for i, phaseSetting := range phaseSettings {
		ch := make(chan int, 2)
		ch <- phaseSetting

		inputs[i] = ch
		if i == 0 {
			outputs[(len(phaseSettings) - 1)] = ch
		} else {
			outputs[(i - 1)] = ch
		}
	}

	inputs[0] <- 0

	wg := &sync.WaitGroup{}

	realOutput := outputs[len(phaseSettings)-1]
	tappedOutput := make(chan int, 2)

	lastValueSeen := 0

	go func(input <-chan int, output chan<- int) {
		for {
			value, ok := <-input
			if !ok {
				close(output)
				break
			}
			lastValueSeen = value
			output <- value
		}

		wg.Done()
	}(realOutput, tappedOutput)
	wg.Add(1)

	inputs[0] = tappedOutput

	for i := range phaseSettings {
		wg.Add(1)

		go func(i int) {
			intcode.RunProgram(inputs[i], outputs[i], program)
			close(outputs[i])
			wg.Done()
		}(i)
	}

	wg.Wait()

	return lastValueSeen
}

func part2(program []int) {
	var maxPermutation []int
	maxOutput := 0

	values := []int{5, 6, 7, 8, 9}
	for _, p := range permutations(values) {
		output := runAmplifiersFeedback(program, p)

		if output > maxOutput {
			maxPermutation = p
			maxOutput = output
		}
	}

	fmt.Println(maxOutput, maxPermutation)
}

func main() {
	programBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic("couldn't read program")
	}

	program := intcode.ParseProgram(strings.TrimSpace(string(programBytes)))

	//part1(program)
	part2(program)
}
