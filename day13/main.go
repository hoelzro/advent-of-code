package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type vec2 struct {
	x int
	y int
}

type WorldUpdater interface {
	Update(x, y, tileID int)
}

type WorldUpdateOutputter struct {
	valuesBuffer []int
	updater      WorldUpdater
}

func (out *WorldUpdateOutputter) Write(value int) {
	out.valuesBuffer = append(out.valuesBuffer, value)

	if len(out.valuesBuffer) == 3 {
		out.updater.Update(out.valuesBuffer[0], out.valuesBuffer[1], out.valuesBuffer[2])
		out.valuesBuffer = nil
	}
}

const (
	Empty  = 0
	Wall   = 1
	Block  = 2
	Paddle = 3
	Ball   = 4
)

const (
	NoDebug   = 0
	PaintGrid = 1
	PrintShit = 2
)

var debugMode int = NoDebug

type WorldState struct {
	score                int
	tiles                map[vec2]int
	ballPosition         vec2
	ballVelocity         vec2
	previousBallVelocity vec2
	previousBallPosition vec2
	paddlePosition       vec2
	intercept            *int
}

type JoystickInput struct {
	world   *WorldState
	program []int
}

type SimulatedReader struct {
	done *bool
}

func (sr *SimulatedReader) Read(program []int, ip int) (int, bool) {
	if *sr.done {
		return 0, false
	} else {
		return 0, true
	}
}

type SimulatedUpdater struct {
	hasResult            *bool
	intercept            *int
	world                *WorldState
	prevBallPosition     vec2
	prevPrevBallPosition vec2
	paddleX              int
	paddleY              int
}

func (updater *SimulatedUpdater) Update(x, y, tileID int) {
	if tileID != Ball {
		return
	}

	if debugMode == PrintShit {
		fmt.Printf("simulated ball pos: (%d, %d)\n", x, y)
	}

	ballPos := vec2{x, y}

	previousVelocity := vec2{
		x: updater.prevBallPosition.x - updater.prevPrevBallPosition.x,
		y: updater.prevBallPosition.y - updater.prevPrevBallPosition.y,
	}

	currentVelocity := vec2{
		x: ballPos.x - updater.prevBallPosition.x,
		y: ballPos.y - updater.prevBallPosition.y,
	}

	var possibleIntercepts []int

	if y == updater.paddleY {
		possibleIntercepts = append(possibleIntercepts, x-currentVelocity.x)
		possibleIntercepts = append(possibleIntercepts, x)

		if debugMode == PrintShit {
			fmt.Println("[1] ball intersects paddle plane at x = ", x, ", y = ", y)
		}
	} else {
		if updater.prevBallPosition.y+previousVelocity.y == updater.paddleY && previousVelocity.y != currentVelocity.y && abs(updater.prevBallPosition.x-updater.paddleX) <= 1 {
			if debugMode == PrintShit {
				fmt.Println("[2] ball intersects paddle plane at x = ", x, ", y = ", y)
			}

			if updater.prevBallPosition.x == updater.paddleX {
				possibleIntercepts = append(possibleIntercepts, updater.paddleX)
			} else {
				possibleIntercepts = append(possibleIntercepts, updater.paddleX)
			}
		}

		updater.prevPrevBallPosition = updater.prevBallPosition
		updater.prevBallPosition = ballPos
	}

	if possibleIntercepts != nil {
		*updater.hasResult = true
		idx := int(rand.Uint32()) % len(possibleIntercepts)
		*updater.intercept = possibleIntercepts[idx]
	}
}

func (ji *JoystickInput) Read(program []int, ip int) (int, bool) {
	if debugMode == PrintShit {
		fmt.Println("reading value from joystick - paddle is at ", ji.world.paddlePosition)
	}

	if ji.world.previousBallPosition.x == 0 && ji.world.previousBallPosition.y == 0 {
		return 0, true
	}

	if ji.world.intercept != nil {
		return sign(*ji.world.intercept - ji.world.paddlePosition.x), true
	}

	// make a copy of the program, and run it until the ball's position is level with the paddle, then
	// move the paddle towards that spot
	programCopy := make([]int, len(program))
	copy(programCopy, program)

	hasResult := false
	intercept := 0

	input := &SimulatedReader{done: &hasResult}
	output := &WorldUpdateOutputter{
		updater: &SimulatedUpdater{
			world:                ji.world,
			hasResult:            &hasResult,
			intercept:            &intercept,
			prevBallPosition:     ji.world.ballPosition,
			prevPrevBallPosition: ji.world.previousBallPosition,
			paddleX:              ji.world.paddlePosition.x,
			paddleY:              ji.world.paddlePosition.y,
		},
	}

	RunProgram(input, output, programCopy, ip)

	if debugMode == PrintShit {
		fmt.Println("returning ", sign(intercept-ji.world.paddlePosition.x))
	}

	ji.world.intercept = new(int)
	*ji.world.intercept = intercept

	return sign(*ji.world.intercept - ji.world.paddlePosition.x), true
}

type MainUpdater struct {
	stepCounter int
	world       *WorldState
	tiles       map[vec2]int
}

func (updater *MainUpdater) Update(x, y, tileID int) {
	if x == -1 && y == 0 {
		updater.world.score = tileID
	} else {
		updater.tiles[vec2{x, y}] = tileID

		if tileID == Ball {
			updater.stepCounter++

			newPosition := vec2{x, y}

			updater.world.previousBallPosition = updater.world.ballPosition
			updater.world.previousBallVelocity = updater.world.ballVelocity
			updater.world.ballPosition = newPosition
			updater.world.ballVelocity.x = newPosition.x - updater.world.previousBallPosition.x
			updater.world.ballVelocity.y = newPosition.y - updater.world.previousBallPosition.y

			if updater.world.ballVelocity != updater.world.previousBallVelocity && newPosition.y >= 21 {
				updater.world.intercept = nil
			}

			if debugMode == PaintGrid {
				drawGrid(updater.tiles)
				fmt.Println(updater.stepCounter)
				time.Sleep(100 * time.Millisecond)
			}

			if debugMode == PrintShit {
				fmt.Printf("[step = %d] new ball position: %v\n", updater.stepCounter, newPosition)
			}
		} else if tileID == Paddle {
			if debugMode == PrintShit {
				fmt.Println("new paddle position: ", vec2{x, y})
			}
			updater.world.paddlePosition = vec2{x, y}
		}
	}
}

func runGame(program []int) (map[vec2]int, int) {
	tiles := make(map[vec2]int)
	world := &WorldState{tiles: tiles}

	input := &JoystickInput{world: world, program: program}
	output := &WorldUpdateOutputter{
		updater: &MainUpdater{world: world, tiles: tiles},
	}

	RunProgram(input, output, program, 0)

	return tiles, world.score
}

func main() {
	rand.Seed(time.Now().Unix())

	programBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err.Error())
	}

	pristineProgram := ParseProgram(strings.TrimSpace(string(programBytes)))
	program := make([]int, len(pristineProgram))
	copy(program, pristineProgram)

	// Part 1
	tiles, _ := runGame(program)

	count := 0
	for _, tileID := range tiles {
		if tileID == 2 {
			count++
		}
	}
	fmt.Println(count)

	// Part 2
	copy(program, pristineProgram)
	program[0] = 2
	tiles, score := runGame(program)
	fmt.Println("score:", score)
}
