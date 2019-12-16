package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type reagent struct {
	name   string
	amount int
}

type reaction struct {
	input  []reagent
	output reagent
}

func parseReagent(input string) reagent {
	reagentRE := regexp.MustCompile(`(\d+)\s*(\w+)`)
	submatches := reagentRE.FindStringSubmatch(input)

	if submatches == nil {
		panic("unable to parse reagent")
	}

	amount, err := strconv.Atoi(submatches[1])

	if err != nil {
		panic(err.Error())
	}

	return reagent{
		name:   submatches[2],
		amount: amount,
	}
}

func loadReactions(input string) []reaction {
	var reactions []reaction

	reactionSidesRE := regexp.MustCompile(`\s*=>\s*`)
	commasRE := regexp.MustCompile(`\s*,\s*`)

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		parts := reactionSidesRE.Split(line, 2)
		output := parseReagent(parts[1])
		inputStringss := commasRE.Split(parts[0], -1)
		var inputs []reagent
		for _, s := range inputStringss {
			inputs = append(inputs, parseReagent(s))
		}
		reactions = append(reactions, reaction{
			input:  inputs,
			output: output,
		})
	}

	return reactions
}

func part1(reactions []reaction) {
	reactionForReagent := make(map[string]reaction)

	for _, reaction := range reactions {
		reactionForReagent[reaction.output.name] = reaction
	}

	fuelOutput := reactionForReagent["FUEL"]

	oreRequired := 0
	stockpile := make(map[string]int)

	q := list.New()
	q.PushBack(fuelOutput)

	for q.Front() != nil {
		e := q.Front()
		q.Remove(e)
		reaction := e.Value.(reaction)

		for _, input := range reaction.input {
			if input.name == "ORE" {
				oreRequired += input.amount
			} else {
				reactionForInput, exists := reactionForReagent[input.name]
				if !exists {
					panic("no matching reaction")
				}

				for stockpile[input.name] < input.amount {
					q.PushBack(reactionForInput)
					stockpile[input.name] += reactionForInput.output.amount
				}
				stockpile[input.name] -= input.amount
			}
		}
	}

	fmt.Println(oreRequired)

}

func part2(reactions []reaction) {
	reactionForReagent := make(map[string]reaction)

	for _, reaction := range reactions {
		reactionForReagent[reaction.output.name] = reaction
	}

	fuelOutput := reactionForReagent["FUEL"]

	fuelYield := 0
	oreRemaining := 1000000000000
	stockpile := make(map[string]int)

	q := list.New()

	for oreRemaining > 0 {
		q.PushBack(fuelOutput)
	qLoop:
		for q.Front() != nil {
			e := q.Front()
			q.Remove(e)
			reaction := e.Value.(reaction)

			for _, input := range reaction.input {
				if input.name == "ORE" {
					oreRemaining -= input.amount
					if oreRemaining < 0 {
						break qLoop
					}
				} else {
					reactionForInput, exists := reactionForReagent[input.name]
					if !exists {
						panic("no matching reaction")
					}

					for stockpile[input.name] < input.amount {
						q.PushBack(reactionForInput)
						stockpile[input.name] += reactionForInput.output.amount
					}
					stockpile[input.name] -= input.amount
				}
			}
		}

		if q.Front() == nil {
			fuelYield++
		}
	}

	fmt.Println(fuelYield)
}

func main() {
	inputBytes, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		panic("Unable to load file: " + err.Error())
	}

	reactions := loadReactions(strings.TrimSpace(string(inputBytes)))

	part1(reactions)
	part2(reactions)
}
