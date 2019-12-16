package main

func sign(n int) int {
	if n < 0 {
		return -1
	} else if n > 0 {
		return 1
	} else {
		return 0
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
