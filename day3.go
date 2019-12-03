package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

func distanceFromOrigin(p point) int {
	return int(math.Abs(float64(p.x)) + math.Abs(float64(p.y)))
}

func parseWire(line string) []point {
	pieces := strings.Split(line, ",")

	currentPos := point{x: 0, y: 0}
	wire := make([]point, 0)
	wire = append(wire, currentPos)

	for _, piece := range pieces {
		var dx, dy int

		switch piece[0] {
		case 'R':
			dx = 1
		case 'L':
			dx = -1
		case 'U':
			dy = -1
		case 'D':
			dy = 1
		default:
			panic("Unrecognized direction: " + piece)
		}

		run, err := strconv.Atoi(piece[1:])

		if err != nil {
			panic("Unrecognized direction: " + piece)
		}

		for offset := 1; offset <= run; offset++ {
			pos := point{x: currentPos.x + dx, y: currentPos.y + dy}
			wire = append(wire, pos)
			currentPos = pos
		}
	}

	return wire
}

func findIntersectionClosestToOrigin(wireA, wireB []point) point {
	var closest point

	commonPoints := make(map[point]bool)

	for _, p := range wireA {
		commonPoints[p] = true
	}

	for _, p := range wireB {
		_, present := commonPoints[p]

		if present {
			closestDistance := distanceFromOrigin(closest)
			pDistance := distanceFromOrigin(p)

			if closest.x == 0 && closest.y == 0 || pDistance < closestDistance {
				closest = p
			}
		}
	}

	return closest
}

func findIntersectionDistanceClosestToWireStarts(wireA, wireB []point) int {
	var closest point
	closestDistance := 0

	distanceFromStartA := make(map[point]int)

	for distance, p := range wireA {
		_, present := distanceFromStartA[p]
		if !present {
			distanceFromStartA[p] = distance
		}
	}

	seen := make(map[point]bool)

	for distanceB, p := range wireB {
		_, present := seen[p]

		if present {
			continue
		}

		distanceA, present := distanceFromStartA[p]

		if present {
			distanceP := distanceA + distanceB

			if closest.x == 0 && closest.y == 0 || distanceP < closestDistance {
				closest = p
				closestDistance = distanceP
			}
		}
	}

	return closestDistance
}

func main() {
	wire1 := parseWire(os.Args[1])
	wire2 := parseWire(os.Args[2])

	// Part 1
	intersection := findIntersectionClosestToOrigin(wire1, wire2)
	fmt.Println(distanceFromOrigin(intersection))

	// Part 2
	distance := findIntersectionDistanceClosestToWireStarts(wire1, wire2)
	fmt.Println(distance)
}
