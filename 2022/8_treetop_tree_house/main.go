package main

/*
	This solution is very ugly :D
*/

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Tree struct {
	Height int
}

func NewTree(height int) Tree {
	return Tree{Height: height}
}

func loadGrid() [][]Tree {
	readFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Prepare the grid by converting the input into a 2D array of runes
	input := string(readFile)
	lines := strings.Split(input, "\n")
	var grid [][]Tree
	for _, line := range lines {
		if line == "" {
			continue
		}

		treeLine := make([]Tree, len(line))
		gridLine := strings.Split(line, "")
		for x, trees := range gridLine {
			height, err := strconv.Atoi(trees)
			if err != nil {
				panic(err)
			}
			treeLine[x] = NewTree(height)
		}
		grid = append(grid, treeLine)
	}

	return grid
}

func withinBounds(maxX, maxY, x, y int) bool {
	return (x <= maxX && x >= 0) && (y <= maxY && y >= 0)
}

func atBounds(maxX, maxY, x, y int) bool {
	return (x == 0 || x == maxX) || (y == 0 || y == maxY)
}

type Pos struct {
	X int
	Y int
}

func NewPos(x, y int) Pos {
	return Pos{X: x, Y: y}
}

type Direction uint8

const (
	North Direction = iota
	East
	South
	West
)

func directionToString(d Direction) string {
	switch d {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	}

	panic(d)
}

// Gets the relative direction from source x and source y to x and y.
// We can be sure that one of the axis are equal
func getRelativeDirection(sx, sy, x, y int) Direction {
	if sx == x {
		if sy > y {
			return South
		}
		return North
	}

	// sy == y
	if sy > y {
		return West
	}

	return East
}

func part1(grid [][]Tree, x, y int) bool {
	maxY := len(grid) - 1
	maxX := len(grid[0]) - 1
	targetHeight := grid[y][x].Height

	paths := make([][]Tree, 4)
	// Go in from the left
	for xx := 0; xx < x; xx++ {
		paths[0] = append(paths[0], NewTree(grid[y][xx].Height))
	}

	// Go to the right
	for xx := x + 1; xx <= maxX; xx++ {
		paths[1] = append(paths[1], NewTree(grid[y][xx].Height))
	}

	// Go from the top
	for yy := 0; yy < y; yy++ {
		paths[2] = append(paths[2], NewTree(grid[yy][x].Height))
	}

	// Go towards the bottom
	for yy := y + 1; yy <= maxY; yy++ {
		paths[3] = append(paths[3], NewTree(grid[yy][x].Height))
	}

	all := func(tree []Tree, p func(Tree) bool) bool {
		for _, v := range tree {
			if !p(v) {
				return false
			}
		}

		return true
	}

	p := func(p Tree) bool {
		return p.Height < targetHeight
	}

	return all(paths[0], p) || all(paths[1], p) || all(paths[2], p) || all(paths[3], p)
}

func viewingDistance(grid [][]Tree, x, y int) uint64 {
	maxY := len(grid) - 1
	maxX := len(grid[0]) - 1
	height := grid[y][x].Height

	var distUp uint64 = 0
	var distDown uint64 = 0
	var distLeft uint64 = 0
	var distRight uint64 = 0

	// Go towards the left
	for xx := x - 1; xx >= 0; xx-- {
		distLeft++
		if grid[y][xx].Height >= height {
			break
		}
	}

	// Go towards the right
	for xx := x + 1; xx <= maxX; xx++ {
		distRight++
		if grid[y][xx].Height >= height {
			break
		}
	}

	// Go upwards
	for yy := y - 1; yy >= 0; yy-- {
		distUp++
		if grid[yy][x].Height >= height {
			break
		}
	}

	// Go down
	for yy := y + 1; yy <= maxY; yy++ {
		distDown++
		if grid[yy][x].Height >= height {
			break
		}
	}

	return distUp * distDown * distLeft * distRight
}

func main() {
	grid := loadGrid()
	boundX := len(grid[0]) - 2
	boundY := len(grid) - 2
	visibleTrees := 2*len(grid[0]) + (len(grid)-2)*2

	for y := 1; y <= boundY; y++ {
		for x := 1; x <= boundX; x++ {
			if part1(grid, x, y) {
				visibleTrees++
			}
		}
	}

	log.Printf("Part1: %v", visibleTrees)

	var highestDistance uint64
	for y := 1; y <= boundY; y++ {
		for x := 1; x <= boundX; x++ {
			dist := viewingDistance(grid, x, y)

			if dist > highestDistance {
				highestDistance = dist
			}
		}
	}
	log.Printf("Part2: %v", highestDistance)
}
