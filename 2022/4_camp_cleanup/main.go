package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// The Range struct models a range with inclusive bounds
type Range struct {
	From int
	To   int
}

// Evalutates whether or not the range `r` contains another range `otherRange`.
// Keep in mind that even if one Range does not contain the other, the other
// might contain the first one, so you might need to check "in both directions".
func (r *Range) contains(otherRange Range) bool {
	return otherRange.From >= r.From && otherRange.To <= r.To
}

// Computes whether there is any overlap between two ranges
func (r *Range) overlaps(otherRange Range) bool {
	for i := r.From; i <= r.To; i++ {
		if i >= otherRange.From && i <= otherRange.To {
			return true
		}
	}

	return false
}

// Creates a new range from a line of the input.
// Line is expected to be of shape: "from-to", e.g.: "1-3"
func NewRange(line string) Range {
	fromTo := strings.Split(line, "-")

	// Convert From and To to integers, panic if we encounter an error as input _should_ be perfect
	from, err := strconv.Atoi(fromTo[0])
	if err != nil {
		panic(err)
	}

	to, err := strconv.Atoi(fromTo[1])
	if err != nil {
		panic(err)
	}

	return Range{From: from, To: to}
}

// Solves for part 1
func solve1(line string) int {
	// The line should be of shape "<<range>>,<<range>>"
	ranges := strings.Split(line, ",")
	leftRange := NewRange(ranges[0])
	rightRange := NewRange(ranges[1])

	if leftRange.contains(rightRange) || rightRange.contains(leftRange) {
		return 1
	}

	return 0
}

func solve2(line string) int {
	// The line should be of shape "<<range>>,<<range>>"
	ranges := strings.Split(line, ",")
	leftRange := NewRange(ranges[0])
	rightRange := NewRange(ranges[1])

	if leftRange.overlaps(rightRange) {
		return 1
	}

	return 0
}

func main() {
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var result1 int
	var result2 int

	lines := strings.Split(string(file), "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1] // Drop the last element if its empty
	}

	for _, line := range lines {
		result1 += solve1(line)
		result2 += solve2(line)
	}

	fmt.Printf("Part 1 result: %d\n", result1)
	fmt.Printf("Part 2 result: %d\n", result2)
}
