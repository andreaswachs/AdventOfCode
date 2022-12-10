package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Feels a little jank that this is hard coded
const crtOutputLines = 7

// cpu struct represents our simple CPU from the problem with the x register and a cycle counter
// the two last fields are for solving the problem, for part 1 and two respectively
type cpu struct {
	x           int
	cycle       int
	part1Answer int
	crtOutput   [crtOutputLines][]rune
}

// Performs an instruction put in a string, which is either noop or addx n
func (c *cpu) performInstruction(instruction string) {
	switch {
	case instruction == "noop":
		c.doCycle()
	case instruction[:4] == "addx":
		n, err := strconv.Atoi(instruction[5:])
		if err != nil {
			panic(err)
		}
		c.doCycle()
		c.doCycle()
		c.x += n
	}
}

// Performs a cycle and does computations for both parts of the problem
func (c *cpu) doCycle() {
	// This seems to only be relevant for part1
	// Adds the signal strength (cycle * x register value) to the answer
	// for cycles 20 and every 40th cycle after that
	if c.cycle == 20 || (c.cycle%20 == 0 && (c.cycle/20)%2 == 1) {
		c.part1Answer += c.x * c.cycle
	}

	// Calculate the CRT output (part 2)
	// We need to adjust the cycle since it begins at 1,
	// but our calculations needs for it to have begun at 0
	row := (c.cycle - 1) / 40
	index := (c.cycle - 1) % 40

	// if the index is within c.x - 1 <= index <= c.x + 1 then the output is a # else a .
	if index >= c.x-1 && index <= c.x+1 {
		c.crtOutput[row] = append(c.crtOutput[row], '#')
	} else {
		c.crtOutput[row] = append(c.crtOutput[row], '.')
	}

	c.cycle++
}

func NewCpu() *cpu {
	return &cpu{
		x:           1,
		cycle:       1,
		part1Answer: 0,
		crtOutput:   [crtOutputLines][]rune{},
	}
}

func main() {
	// Scan the input file
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	c := NewCpu()

	for fileScanner.Scan() {
		// Get the instruction
		instruction := fileScanner.Text()
		// Perform the instruction
		c.performInstruction(instruction)
	}

	// Part 1 answer
	fmt.Printf("Part 1 answer: %d\n", c.part1Answer)

	// Part 2 answer
	fmt.Printf("Part 2 answer:\n")

	for _, row := range c.crtOutput {
		for _, char := range row {
			fmt.Printf("%c", char)
		}
		fmt.Println()
	}
}
