package main

import (
	"fmt"

	"io/ioutil"
	"os"
	"strings"

	"hoelz.ro/advent-of-code/2019/intcode"
)

type point struct {
	x int
	y int
}

func turn(currentDirection point, turn int) point {
	if turn == 1 { // clockwise
		if currentDirection.x == 1 { // facing right
			return point{
				x: 0,
				y: 1,
			}
		} else if currentDirection.x == -1 { // facing left
			return point{
				x: 0,
				y: -1,
			}
		} else if currentDirection.y == 1 { // facing down
			return point{
				x: -1,
				y: 0,
			}
		} else { // facing up
			return point{
				x: 1,
				y: 0,
			}
		}
	} else { // counterclockwise
		if currentDirection.x == 1 { // facing right
			return point{
				x: 0,
				y: -1,
			}
		} else if currentDirection.x == -1 { // facing left
			return point{
				x: 0,
				y: 1,
			}
		} else if currentDirection.y == 1 { // facing down
			return point{
				x: 1,
				y: 0,
			}
		} else { // facing up
			return point{
				x: -1,
				y: 0,
			}
		}
	}
}

func move(pos point, delta point) point {
	return point{
		x: pos.x + delta.x,
		y: pos.y + delta.y,
	}
}

func runPainter(program []int, initialIsWhite bool) map[point]int {
	grid := make(map[point]int)
	robotPosition := point{}
	robotDirection := point{x: 0, y: -1}

	if initialIsWhite {
		grid[robotPosition] = 1
	}

	input := make(chan int)
	output := make(chan int)
	done := make(chan bool)

	go func() {
		intcode.RunProgram(input, output, program)
		done <- true
	}()

mainLoop:
	for {
		select {
		case input <- grid[robotPosition]:
			newColor := <-output
			turnDirection := <-output
			grid[robotPosition] = newColor
			robotDirection = turn(robotDirection, turnDirection)
			robotPosition = move(robotPosition, robotDirection)
		case <-done:
			break mainLoop
		}
	}

	return grid
}

func drawGrid(grid map[point]int) {
	minX := 0
	minY := 0
	maxX := 0
	maxY := 0

	for p := range grid {
		if p.x < minX {
			minX = p.x
		}

		if p.y < minY {
			minY = p.y
		}

		if p.x > maxX {
			maxX = p.x
		}

		if p.y > maxY {
			maxY = p.y
		}
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			p := point{x, y}

			if grid[p] == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("")
	}
}

func main() {
	programBytes, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		panic("unable to read program")
	}

	program := intcode.ParseProgram(strings.TrimSpace(string(programBytes)))

	grid := runPainter(program, false)

	fmt.Println(len(grid))

	grid = runPainter(program, true)
	drawGrid(grid)
}
