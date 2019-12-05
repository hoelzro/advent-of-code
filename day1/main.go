package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [masses...]\n", os.Args[0])
		return
	}

	totalFuel := 0

	for _, arg := range os.Args[1:] {
		mass, err := strconv.Atoi(arg)

		if err != nil {
			fmt.Fprintln(os.Stderr, "couldn't parse %s as an integer: %v\n", arg, err)
			os.Exit(1)
		}

		fuel := (mass / 3) - 2

		fuelForModule := 0

		for fuel > 0 {
			fuelForModule += fuel
			fuel = (fuel / 3) - 2
		}

		totalFuel += fuelForModule
	}

	fmt.Printf("total fuel required: %d\n", totalFuel)
}
