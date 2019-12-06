package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseOrbits(input *bufio.Scanner) [][]string {
	result := [][]string{}

	for input.Scan() {
		line := input.Text()

		chunks := strings.Split(line, ")")
		result = append(result, chunks)
	}

	return result
}

type orbitTree struct {
	label    string
	parent   *orbitTree
	children []*orbitTree
}

func buildOrbitTree(orbitList [][]string) *orbitTree {
	nameToNodeMap := make(map[string]*orbitTree)

	for _, pair := range orbitList {
		orbitee := pair[0]
		orbiter := pair[1]

		orbiteeNode := nameToNodeMap[orbitee]
		if orbiteeNode == nil {
			orbiteeNode = &orbitTree{
				label: orbitee,
			}
			nameToNodeMap[orbitee] = orbiteeNode
		}

		orbiterNode := nameToNodeMap[orbiter]
		if orbiterNode == nil {
			orbiterNode = &orbitTree{
				label: orbiter,
			}
			nameToNodeMap[orbiter] = orbiterNode
		}

		orbiterNode.parent = orbiteeNode
		orbiteeNode.children = append(orbiteeNode.children, orbiterNode)
	}

	return nameToNodeMap["COM"]
}

func calculateNumPaths(orbits *orbitTree, accum int) int {
	childrenPaths := 0
	for _, child := range orbits.children {
		childrenPaths += calculateNumPaths(child, accum+1)
	}
	return accum + childrenPaths
}

func findNodeByLabel(orbits *orbitTree, label string) *orbitTree {
	if orbits.label == label {
		return orbits
	}

	for _, child := range orbits.children {
		found := findNodeByLabel(child, label)
		if found != nil {
			return found
		}
	}

	return nil
}

func pathToRoot(orbits *orbitTree) []string {
	path := []string{}

	parent := orbits.parent

	for parent != nil {
		path = append(path, parent.label)
		parent = parent.parent
	}

	return path
}

func reversePath(path []string) []string {
	reversed := make([]string, len(path))

	for i, node := range path {
		reversed[len(reversed)-i-1] = node
	}

	return reversed
}

func findClosestCommonAncestor(nodeA, nodeB *orbitTree) string {
	pathA := reversePath(pathToRoot(nodeA))
	pathB := reversePath(pathToRoot(nodeB))

	i := 0

	for pathA[i] == pathB[i] {
		i++
	}

	return pathA[i-1]
}

func findDistanceToAncestor(orbits *orbitTree, ancestor string) int {
	distance := 0
	parent := orbits.parent

	for parent != nil {
		if parent.label == ancestor {
			return distance
		}

		// increment *after* the check, because distance to the body
		// we're currently orbiting is 0
		distance++
		parent = parent.parent
	}

	panic("not found")
}

func main() {
	var input *os.File
	if len(os.Args) < 2 {
		input = os.Stdin
	} else {
		input, _ = os.Open(os.Args[1])
	}

	orbits := parseOrbits(bufio.NewScanner(input))
	tree := buildOrbitTree(orbits)
	fmt.Println(calculateNumPaths(tree, 0))

	youNode := findNodeByLabel(tree, "YOU")
	santaNode := findNodeByLabel(tree, "SAN")

	common := findClosestCommonAncestor(youNode, santaNode)
	fmt.Println(findDistanceToAncestor(youNode, common) + findDistanceToAncestor(santaNode, common))
}
