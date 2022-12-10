package main

import (
	"bufio"
	"crypto"
	_ "crypto/sha256"
	"fmt"
	"hash"
	"math"
	"os"
	"strconv"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Position struct {
	X int
	Y int
}

func (p *Position) Distance(other Position) float64 {
	return math.Sqrt(math.Pow(float64(p.X-other.X), 2) + math.Pow(float64(p.Y-other.Y), 2))
}

func (p *Position) NeedsToMove(other *Position) bool {
	// We need to move if we are not diagonally adjacent to the other position

	// Check if we're on top of the other position
	if p.X == other.X && p.Y == other.Y {
		return false
	}

	// Check if we are on the same row
	if p.Y == other.Y {
		// Check if we are adjacent
		if p.X == other.X+1 || p.X == other.X-1 {
			return false
		}
	}

	// Check if we are on the same column
	if p.X == other.X {
		// Check if we are adjacent
		if p.Y == other.Y+1 || p.Y == other.Y-1 {
			return false
		}
	}

	// Check to see if we're on the diagonals
	if p.X == other.X+1 && p.Y == other.Y+1 {
		return false
	}
	if p.X == other.X-1 && p.Y == other.Y-1 {
		return false
	}
	if p.X == other.X+1 && p.Y == other.Y-1 {
		return false
	}
	if p.X == other.X-1 && p.Y == other.Y+1 {
		return false
	}

	// We need to move
	return true
}

type RopeKnot struct {
	Name        rune
	Position    *Position
	NextKnot    *RopeKnot
	HasVisited  map[string]Position
	MovementLog []Position
	hasher      hash.Hash
}

func (r *RopeKnot) GetLastKnot() *RopeKnot {
	if r.NextKnot == nil {
		return r
	}

	return r.NextKnot.GetLastKnot()
}

func (r *RopeKnot) AddKnot(knot *RopeKnot) {
	// Add the knot to the end of the rope
	if r.NextKnot == nil {
		r.NextKnot = knot
		return
	}

	r.NextKnot.AddKnot(knot)
}

func (r *RopeKnot) SavePositionHash(oldPosition Position) {
	buffer := []byte(fmt.Sprintf("%06d,%06d", oldPosition.X, oldPosition.Y))

	r.hasher.Write(buffer)
	visitedBitstring := r.hasher.Sum(nil)
	r.hasher.Reset()

	visiedHash := string(visitedBitstring)
	r.HasVisited[visiedHash] = oldPosition
}

func (r *RopeKnot) Move(d Direction) {
	// Move the rope end in the given direction by the given distance
	// and update the movement log
	oldPosition := Position{
		X: r.Position.X,
		Y: r.Position.Y,
	}

	r.Position = &Position{
		X: r.Position.X,
		Y: r.Position.Y,
	}

	switch d {
	case Up:
		r.Position.Y++
	case Down:
		r.Position.Y--
	case Left:
		r.Position.X--
	case Right:
		r.Position.X++
	}

	// Do "memoising" computations with the old position
	r.MovementLog = append(r.MovementLog, oldPosition)
	r.SavePositionHash(oldPosition)

	// Update the movement subscribers
	r.UpdateNextKnot()
}

func (r *RopeKnot) UpdateNextKnot() {
	if r.NextKnot != nil {
		r.NextKnot.Update(r)
	}
}

func (r *RopeKnot) Update(other *RopeKnot) {
	// Check to see if this rope end needs to move
	if !r.Position.NeedsToMove(other.Position) {
		// We don't need to propagate the update if we don't need to move
		return
	}

	// Get the latest position of the other rope end
	otherLatestPosition := other.MovementLog[len(other.MovementLog)-1]

	// Add the current position to the movement log
	posCopy := Position{
		X: r.Position.X,
		Y: r.Position.Y,
	}

	r.MovementLog = append(r.MovementLog, posCopy)

	// Move this rope end to the other rope ends latest position
	r.Position = &otherLatestPosition

	// Log the position
	r.SavePositionHash(otherLatestPosition)

	// Update the movement subscribers
	r.UpdateNextKnot()
}

func NewRopeEnd(name rune) *RopeKnot {
	return &RopeKnot{
		Name:        name,
		Position:    &Position{X: 0, Y: 0},
		NextKnot:    nil,
		HasVisited:  make(map[string]Position),
		MovementLog: []Position{},
		hasher:      crypto.SHA256.New(),
	}
}

func part1(r RopeKnot, head RopeKnot) int {
	// Find the sum of visited positions
	visualize := true
	n := 8
	grid := make([][]rune, n)

	if visualize {
		fmt.Println()
		for i := range grid {
			grid[i] = make([]rune, n)
			for j := range grid[i] {
				grid[i][j] = '.'
			}
		}
	}

	sum := 0
	for _, v := range r.HasVisited {
		if visualize {
			grid[v.Y][v.X] = '#'
		}
		sum++
	}

	if visualize {
		grid[r.Position.Y][r.Position.X] = r.Name
		grid[head.Position.Y][head.Position.X] = head.Name

		for i := len(grid) - 1; i >= 0; i-- {
			row := grid[i]
			for _, col := range row {
				fmt.Printf("%c", col)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	return sum
}

func part2(head RopeKnot) int {
	// Find the sum of visited positions
	visualize := false
	n := 80
	offset := 20
	grid := make([][]rune, n)

	if visualize {
		fmt.Println()
		for i := range grid {
			grid[i] = make([]rune, n)
			for j := range grid[i] {
				grid[i][j] = '.'
			}
		}
	}

	sum := 0
	tail := head.GetLastKnot()
	for _, v := range tail.HasVisited {
		if visualize {
			grid[offset+v.Y][offset+v.X] = '#'
		}
		sum++
	}

	if visualize {
		for r := head.NextKnot; r != nil; r = r.NextKnot {
			grid[offset+r.Position.Y][offset+r.Position.X] = r.Name
		}

		grid[offset+head.Position.Y][offset+head.Position.X] = head.Name

		for i := len(grid) - 1; i >= 0; i-- {
			row := grid[i]
			for _, col := range row {
				fmt.Printf("%c", col)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	return sum
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	ropeHead := NewRopeEnd('H')

	ends := os.Args[2]
	numEnds, err := strconv.Atoi(ends)
	if err != nil {
		panic(err)
	}

	for i := 0; i < numEnds; i++ {
		knot := NewRopeEnd(rune(49 + i))
		knot.SavePositionHash(Position{X: 0, Y: 0}) // Hacky way of making sure we remember to also have visited (0, 0)
		ropeHead.AddKnot(knot)
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		direction, distance := parseDirectionAndDistance(fileScanner.Text())
		// Move the amount of distance given, so we're not making jumps
		for i := 0; i < distance; i++ {
			if numEnds > 1 {
				part2(*ropeHead)
			} else {
				part1(*ropeHead.GetLastKnot(), *ropeHead)
			}
			ropeHead.Move(direction)
		}
	}

	if numEnds == 1 {
		fmt.Printf("Part 1 result: %d\n", part1(*ropeHead.GetLastKnot(), *ropeHead))
	} else {
		fmt.Printf("Part 2 result: %d\n", part2(*ropeHead))
	}
}

// Uninteresting helper functions

func parseDirectionAndDistance(s string) (Direction, int) {
	var d Direction
	var distance int

	switch s[0] {
	case 'U':
		d = Up
	case 'D':
		d = Down
	case 'L':
		d = Left
	case 'R':
		d = Right
	}

	fmt.Sscanf(s[1:], "%d", &distance)

	return d, distance
}
