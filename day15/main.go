package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"hoelz.ro/advent-of-code/2019/intcode"
)

type vec2 struct {
	x int
	y int
}

const (
	North = 1
	South = 2
	West  = 3
	East  = 4
)

var Directions []int = []int{North, South, East, West}

const (
	StatusWall         = 0
	StatusSuccess      = 1
	StatusOxygenSystem = 2
)

const (
	Unknown      = 0
	Wall         = 1
	Traversable  = 2
	OxygenSystem = 3
)

func drawMap(botPosition vec2, environs map[vec2]int) {
	fmt.Print("\033[2J")

	minX := 0
	minY := 0
	maxX := 0
	maxY := 0

	for pos := range environs {
		if pos.x < minX {
			minX = pos.x
		}
		if pos.y < minY {
			minY = pos.y
		}
		if pos.x > maxX {
			maxX = pos.x
		}
		if pos.y > maxY {
			maxY = pos.y
		}
	}

	for y := minY - 2; y <= maxY+2; y++ {
		for x := minX - 2; x <= maxX+2; x++ {
			pos := vec2{x, y}

			if x == 0 && y == 0 {
				fmt.Print("\033[41m")
			}

			if pos == botPosition {
				fmt.Print("D")
			} else {
				switch environs[pos] {
				case Unknown:
					fmt.Print("█")
				case Wall:
					fmt.Print("█")
				case Traversable:
					fmt.Print(" ")
				case OxygenSystem:
					fmt.Print("O")
				}
			}
			if x == 0 && y == 0 {
				fmt.Print("\033[0m")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func applyMove(pos vec2, direction int) vec2 {
	switch direction {
	case North:
		return vec2{x: pos.x, y: pos.y - 1}
	case South:
		return vec2{x: pos.x, y: pos.y + 1}
	case East:
		return vec2{x: pos.x + 1, y: pos.y}
	case West:
		return vec2{x: pos.x - 1, y: pos.y}
	}

	panic("the impossible happened")
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

type explorer struct {
	environs   map[vec2]int
	currentPos vec2
	input      chan<- int
	output     <-chan int
}

func (e *explorer) findAdjacentPos(environs map[vec2]int, pos vec2) vec2 {
	dx := abs(pos.x - e.currentPos.x)
	dy := abs(pos.y - e.currentPos.y)

	if (dx == 1 && dy == 0) || (dy == 1 && dx == 0) {
		return e.currentPos
	}

	for _, dir := range Directions {
		newPos := applyMove(pos, dir)
		// XXX alternatively, pick the one closest/easiest to get to from currentPos
		if environs[newPos] == Traversable || environs[newPos] == OxygenSystem {
			return newPos
		}
	}

	panic("Couldn't find adjacent position!")
}

func statusString(status int) string {
	switch status {
	case StatusWall:
		return "wall"
	case StatusSuccess:
		return "sucess"
	case StatusOxygenSystem:
		return "oxygen"
	}
	panic("unknown status")
}

func directionString(direction int) string {
	switch direction {
	case North:
		return "North"
	case South:
		return "South"
	case East:
		return "East"
	case West:
		return "West"
	}
	panic("unknown direction")
}

func (e *explorer) move(direction int) int {
	e.input <- direction
	status := <-e.output

	if status != StatusWall {
		e.currentPos = applyMove(e.currentPos, direction)
	}
	return status
}

func findPath(environs map[vec2]int, start, destination vec2) []int {
	if start == destination {
		return []int{}
	}

	type dijkstraItem struct {
		direction   int
		newPosition vec2
		previous    *dijkstraItem
	}

	q := list.New()
	seen := make(map[vec2]bool)

	seen[start] = true

	for _, dir := range Directions {
		newPos := applyMove(start, dir)
		if environs[newPos] == Traversable || environs[newPos] == OxygenSystem {
			q.PushBack(&dijkstraItem{
				direction:   dir,
				newPosition: newPos,
			})
		}
	}

	for q.Front() != nil {
		next := q.Remove(q.Front()).(*dijkstraItem)

		if seen[next.newPosition] {
			continue
		}
		seen[next.newPosition] = true

		if next.newPosition == destination {
			path := list.New()
			length := 0
			for next != nil {
				path.PushFront(next)
				next = next.previous
				length++
			}
			result := make([]int, 0, length)
			for path.Front() != nil {
				next := path.Remove(path.Front()).(*dijkstraItem)
				result = append(result, next.direction)
			}
			return result
		}

		for _, dir := range Directions {
			newPos := applyMove(next.newPosition, dir)
			if environs[newPos] == Traversable || environs[newPos] == OxygenSystem {
				q.PushBack(&dijkstraItem{
					direction:   dir,
					newPosition: newPos,
					previous:    next,
				})
			}
		}
	}

	return nil
}

func (e *explorer) moveTo(environs map[vec2]int, destination vec2) {
	path := findPath(environs, e.currentPos, destination)

	if path == nil {
		panic(fmt.Sprintf("couldn't find path between %v and %v!", e.currentPos, destination))
	}

	for _, dir := range path {
		status := e.move(dir)
		if status == StatusWall {
			panic("I didn't think there'd be a wall there!")
		}
	}
}

func calculateDirection(posA, posB vec2) int {
	dx := posB.x - posA.x
	dy := posB.y - posA.y

	if dx != 0 && dy != 0 {
		panic("WTF")
	}

	if dx == -1 {
		return West
	} else if dx == 1 {
		return East
	} else {
		if dy == -1 {
			return North
		} else if dy == 1 {
			return South
		} else {
			panic("invalid pair of points")
		}
	}
}

func (e *explorer) exploreEnvirons() vec2 {
	environs := e.environs
	environs[e.currentPos] = Traversable

	// XXX variable name
	//     maintain an explicit queue ordered by distance from current pos?
	stack := list.New()
	for _, dir := range Directions {
		stack.PushFront(applyMove(e.currentPos, dir))
	}

	for stack.Front() != nil {
		pos := stack.Remove(stack.Front()).(vec2)
		if environs[pos] != Unknown {
			continue
		}

		adjacentPos := e.findAdjacentPos(environs, pos)

		e.moveTo(environs, adjacentPos)
		if e.currentPos != adjacentPos {
			panic("wtf")
		}
		dir := calculateDirection(adjacentPos, pos)

		status := e.move(dir)
		switch status {
		case StatusWall:
			environs[pos] = Wall
		case StatusSuccess:
			environs[pos] = Traversable
		case StatusOxygenSystem:
			environs[pos] = OxygenSystem
		}

		if status != StatusWall {
			for _, dir := range Directions {
				newPos := applyMove(pos, dir)
				stack.PushFront(newPos)
			}
		}

		drawMap(e.currentPos, environs)
		time.Sleep(time.Millisecond * 100)
	}

	for pos, feature := range environs {
		if feature == OxygenSystem {
			return pos
		}
	}
	panic("couldn't find oxygen system!")
}

func drawOxygenMap(environs map[vec2]int, hasOxygen map[vec2]bool) {
	fmt.Print("\033[2J")

	minX := 0
	minY := 0
	maxX := 0
	maxY := 0

	for pos := range environs {
		if pos.x < minX {
			minX = pos.x
		}
		if pos.y < minY {
			minY = pos.y
		}
		if pos.x > maxX {
			maxX = pos.x
		}
		if pos.y > maxY {
			maxY = pos.y
		}
	}

	for y := minY - 2; y <= maxY+2; y++ {
		for x := minX - 2; x <= maxX+2; x++ {
			pos := vec2{x, y}

			switch environs[pos] {
			case Unknown:
				fmt.Print("█")
			case Wall:
				fmt.Print("█")
			case Traversable:
				if hasOxygen[pos] {
					fmt.Print("\033[46m \033[0m")
				} else {
					fmt.Print(" ")
				}
			case OxygenSystem:
				fmt.Print("\033[46mO\033[0m")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func fillWithOxygen(environs map[vec2]int) int {
	var oxygenSystemPos vec2

	for pos, feature := range environs {
		if feature == OxygenSystem {
			oxygenSystemPos = pos
			break
		}
	}

	hasOxygen := make(map[vec2]bool)
	hasOxygen[oxygenSystemPos] = true

	numSteps := 0

	for {
		drawOxygenMap(environs, hasOxygen)

		vacuumPresent := false

		for pos, feature := range environs {
			if feature == Traversable || feature == OxygenSystem {
				if !hasOxygen[pos] {
					vacuumPresent = true
					break
				}
			}
		}

		if !vacuumPresent {
			break
		}

		addOxygen := make(map[vec2]bool)

		for pos := range hasOxygen {
			for _, dir := range Directions {
				adjacentPos := applyMove(pos, dir)
				if environs[adjacentPos] == Traversable || environs[adjacentPos] == OxygenSystem {
					addOxygen[adjacentPos] = true
				}
			}
		}

		for pos := range addOxygen {
			hasOxygen[pos] = true
		}

		time.Sleep(time.Millisecond * 100)

		numSteps++
	}

	return numSteps
}

func main() {
	programBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err.Error())
	}
	program := intcode.ParseProgram(strings.TrimSpace(string(programBytes)))

	input := make(chan int)
	output := make(chan int)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		explorer := &explorer{
			environs: make(map[vec2]int),
			input:    input,
			output:   output,
		}
		oxygenSystem := explorer.exploreEnvirons()
		// Part 1
		fmt.Println(len(findPath(explorer.environs, vec2{}, oxygenSystem)))
		// Part 2
		fmt.Println(fillWithOxygen(explorer.environs))
		close(input)
		wg.Done()
	}()

	intcode.RunProgram(input, output, program)

	wg.Wait()
}
