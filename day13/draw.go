package main

import (
	"fmt"
)

func drawGrid(tiles map[vec2]int) {
	minX := 0
	minY := 0
	maxX := 0
	maxY := 0

	for v := range tiles {
		if v.x < minX {
			minX = v.x
		}
		if v.x > maxX {
			maxX = v.x
		}
		if v.y < minY {
			minY = v.y
		}
		if v.y > maxY {
			maxY = v.y
		}
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			ch := " "
			switch tiles[vec2{x, y}] {
			case Empty:
				ch = " "
			case Wall:
				ch = "█"
			case Block:
				ch = "▀"
			case Paddle:
				ch = "_"
			case Ball:
				ch = "⬤"
			}
			fmt.Print(ch)
		}
		fmt.Println("")
	}
}
