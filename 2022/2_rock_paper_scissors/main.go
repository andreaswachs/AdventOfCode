package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//   Encryption card:
//   INPUT		OPPONENT	PLAYER		POINTS FOR USE
//   A			ROCK					1
//   B			PAPER					2
//   C			SCISSOR					3
//   X						ROCK		1
//   Y						PAPER		2
//   Z						SCISSOR 	3
//
// 	Poitns for rounds:
// 	OUTCOME		POINTS
// 	LOST		0
// 	DRAW		6
// 	WIN			6

func main() {
	// Open the input file

	readFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	var pointsPart1 int
	var pointsPart2 int

	// Prepare to scan the file, line by line
	// credits: https://golangdocs.com/golang-read-file-line-by-line
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		pointsPart1 += pointsForRound(line)
		bb := expectedPointsForRound(line)
		pointsPart2 += bb
	}

	fmt.Printf("Points part 1: %d\n", pointsPart1)
	fmt.Printf("Points part 2: %d\n", pointsPart2)

}

func pointsForRound(inputLine string) int {
	inputs := strings.Split(inputLine, " ")

	return pointsForAction(inputs[1]) + pointsForRoundOutcome(inputLine)
}

// This function returns the results foreach line for part two
func expectedPointsForRound(inputLine string) int {
	switch inputLine {
	case "A X": // rock lose
		return 3 // paper + lose
	case "A Y": // rock draw
		return 4
	case "A Z": // rock win
		return 8
	case "B X": // paper lose
		return 1
	case "B Y": // paper draw
		return 5
	case "B Z": // paper win
		return 9
	case "C X": // scissor lose
		return 2
	case "C Y": // scissor draw
		return 6
	case "C Z": // scissor win
		return 7
	}

	panic(inputLine)
}

func pointsForAction(action string) int {
	switch action {
	case "X":
		return 1
	case "Y":
		return 2
	case "Z":
		return 3
	}
	// If we reach this point then something went wrong
	panic(action)
}

func pointsForRoundOutcome(input string) int {
	switch input {
	case "A X": // rock rock
		return 3
	case "A Y": // rock paper
		return 6
	case "A Z": // rock scissor
		return 0
	case "B X": // paper rock
		return 0
	case "B Y": // paper paper
		return 3
	case "B Z": // paper scissor
		return 6
	case "C X": // scissor rock
		return 6
	case "C Y": // scissor paper
		return 0
	case "C Z": // scissor scissor
		return 3
	}

	panic(input)
}
