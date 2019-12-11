package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

type point struct {
	x int
	y int
}

func loadMap(filename string) [][]bool {
	contents, err := ioutil.ReadFile(filename)

	if err != nil {
		panic("can't read file")
	}

	lines := strings.Split(strings.TrimSpace(string(contents)), "\n")

	result := [][]bool{}

	for _, line := range lines {
		row := make([]bool, len(line))
		result = append(result, row)

		for i, char := range line {
			row[i] = char == '#'
		}
	}

	return result
}

func getAsteroids(m [][]bool) []point {
	result := []point{}

	for y, row := range m {
		for x, c := range row {
			if c {
				result = append(result, point{x, y})
			}
		}
	}

	return result
}

func isVisibleFrom(asteroids []point, origin, p point) bool {
	otherDX := origin.x - p.x
	otherDY := origin.y - p.y
	otherAngle := math.Atan2(float64(otherDY), float64(otherDX))

	for _, yetAnotherAsteroid := range asteroids {
		if origin == yetAnotherAsteroid || p == yetAnotherAsteroid {
			continue
		}

		yetAnotherDX := origin.x - yetAnotherAsteroid.x
		yetAnotherDY := origin.y - yetAnotherAsteroid.y
		yetAnotherAngle := math.Atan2(float64(yetAnotherDY), float64(yetAnotherDX))

		if otherAngle == yetAnotherAngle {
			otherDistance := distance(origin, p)
			yetAnotherDistance := distance(origin, yetAnotherAsteroid)

			if yetAnotherDistance < otherDistance {
				return false
			}
		}

	}

	return true
}

func numVisibleFrom(m [][]bool, p point) int {
	asteroids := getAsteroids(m)
	count := 0

	for _, otherAsteriod := range asteroids {
		if p == otherAsteriod {
			continue
		}

		if isVisibleFrom(asteroids, p, otherAsteriod) {
			count++
		}
	}

	return count
}

func findBestAsteroid(m [][]bool) point {
	var winningAsteroid point
	winningCount := 0

	for _, asteroid := range getAsteroids(m) {
		count := numVisibleFrom(m, asteroid)

		if count > winningCount {
			winningCount = count
			winningAsteroid = asteroid
		}
	}

	return winningAsteroid
}

func distance(a, b point) float64 {
	dx := float64(a.x - b.x)
	dy := float64(a.y - b.y)

	return math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
}

func calculateAngle(origin, p point) float64 {
	dx := float64(p.x - origin.x)
	dy := float64(p.y - origin.y)

	a := math.Atan2(-1*dy, dx)

	a -= math.Pi / 2 // orient with north, rather than the X axis

	a = math.Pi*2 - a // flip to clockwise

	// clamp to [0, 2π]
	for a < 0 {
		a += math.Pi * 2
	}

	for a >= math.Pi*2 {
		a -= math.Pi * 2
	}

	return a
}

func runVaporizer(station point, asteroids []point) []point {
	vaporizations := []point{}

	sort.Slice(asteroids, func(i, j int) bool {
		return distance(station, asteroids[i]) < distance(station, asteroids[j])
	})

	sort.SliceStable(asteroids, func(i, j int) bool {
		return calculateAngle(station, asteroids[i]) < calculateAngle(station, asteroids[j])
	})

	remainingAsteroids := make([]*point, len(asteroids))

	for i := range asteroids {
		remainingAsteroids[i] = &asteroids[i]
	}

	numAsteroids := len(asteroids)

	for numAsteroids > 0 {
		previousAngle := math.NaN()

		for i, asteroid := range remainingAsteroids {
			if asteroid == nil {
				continue
			}

			angle := calculateAngle(station, *asteroid)

			if angle == previousAngle {
				continue
			}

			remainingAsteroids[i] = nil
			numAsteroids--
			vaporizations = append(vaporizations, *asteroid)
			previousAngle = angle
		}
	}

	return vaporizations
}

func main() {
	m := loadMap(os.Args[1])

	best := findBestAsteroid(m)
	fmt.Println(best, numVisibleFrom(m, best))

	vaporizations := runVaporizer(best, getAsteroids(m))

	for i, p := range vaporizations {
		fmt.Println(i+1, p)
	}
}
