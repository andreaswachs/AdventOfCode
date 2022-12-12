package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yourbasic/graph"
)

const (
	END   = 69 // nice
	BEGIN = 83
)

// canMoveBetween determines if you can move from one point to another based on the elevation.
// The rules are that you can move to the next point if the elevation is at most one more
// lexical value than the one you're currently located on.
func canMoveBetween(from, to rune) bool {
	if from == BEGIN {
		from = 'a'
	}
	if to == END {
		to = 'z'
	}
	return (to - from) <= 1
}

type position struct {
	x int
	y int
}

// Usage: go run main.go <file>
// Optionally: set VISUAL=1 as env variable to see path found step by step for part 1
func main() {
	// Start by reading in the file supplied by the user
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n")

	// We need to store several positons: start/source, end and possibly many positions with elevation 'a' (part 2)
	startPos := position{}
	endPos := position{}
	posWithA := []position{}
	grid := [][]rune{}

	// Go through the map and record appropriate locations
	for y, line := range lines {
		for x, r := range line {
			// Finding S
			if r == BEGIN {
				startPos.x = x
				startPos.y = y
			}

			// Finding E
			if r == END {
				endPos.x = x
				endPos.y = y
			}

			if r == 'a' {
				posWithA = append(posWithA, position{x, y})
			}
		}
		if line == "" {
			break
		}
		grid = append(grid, []rune(line))
	}

	// Set some helpful variables based on the map and also some helper functions to go from
	// coordinates to indices on the graph
	width := len(grid[0])
	height := len(grid)
	edges := width * height
	g := graph.New(edges)

	indexWithCoords := func(x, y int) int {
		return (width * (y + 1)) - (width - x)
	}
	coordsWithIndex := func(index int) (int, int) {
		y := index / width
		x := index - (width * y)
		return x, y
	}

	// We go through each cell in the grid and then take all directly adjacent cells and check
	// if we can go from the cell at (x, y) to the cell at (xx, yy)
	for y, row := range grid {
		for x, _ := range row {
			// Check all cells around that one cell in a cross pattern, excluding the (x, y) cell
			for yy := y - 1; yy <= y+1; yy++ {
				for xx := x - 1; xx <= x+1; xx++ {

					// Only consider (xx, yy) cell if it is within the map and not diagonally positioned in regards to (x, y)
					if yy >= 0 && yy < height && xx >= 0 && xx < width && (yy == y || xx == x) && !(yy == y && xx == x) {
						if canMoveBetween(grid[y][x], grid[yy][xx]) {
							from := indexWithCoords(x, y)
							to := indexWithCoords(xx, yy)
							g.AddCost(from, to, 1)
						}
					}
				}
			}
		}
	}

	sourceIndex := indexWithCoords(startPos.x, startPos.y)
	goalIndex := indexWithCoords(endPos.x, endPos.y)
	path, dist := graph.ShortestPath(g, sourceIndex, goalIndex)

	// If the user wants to visualize the path found for part 1
	if os.Getenv("VISUAL") != "" {
		visualizePath(height, width, path, coordsWithIndex, grid)
	}
	fmt.Printf("Result part 1: %d\n", dist)

	// Part2: find the shortest path form any point with elevation 'a'
	var shortestDist int64 = 9223372036854775807 // Max value for int64, lets hope path is shorter
	for _, pos := range posWithA {
		_, dist := graph.ShortestPath(g, indexWithCoords(pos.x, pos.y), goalIndex)
		// If dist == -1, then no path were found.. we ignore that case
		if dist != -1 && dist < shortestDist {
			shortestDist = dist
		}
	}

	fmt.Printf("Result part 2: %d\n", shortestDist)
}

func visualizePath(height int, width int, path []int, coordsWithIndex func(index int) (int, int), grid [][]rune) {
	visualGrid := [][]rune{}
	for i := 0; i < height; i++ {
		visualGrid = append(visualGrid, []rune(strings.Repeat(".", width)))
	}

	for _, index := range path {
		x, y := coordsWithIndex(index)
		visualGrid[y][x] = grid[y][x]

		for _, row := range visualGrid {
			for _, cell := range row {
				fmt.Printf("%c", cell)
			}
			fmt.Println()
		}

		fmt.Println()
	}
}
