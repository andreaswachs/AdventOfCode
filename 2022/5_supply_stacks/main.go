package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type CraneType uint8

const (
	CrateMover9000 CraneType = iota
	CrateMover9001
)
// This data structure reflects the stacks from the problem and the operations implemented
// are performed from the view of the crane.
type Stack[T any] struct {
	items []T
	Version CraneType
}

func NewStack[T any](craneType CraneType) Stack[T] {
	return Stack[T]{
		Version: craneType,
	}
}

func (s *Stack[T]) Print() {
	log.Printf("Stack: %v", s.items)
}

func (s *Stack[T]) Push(item T) {
	log.Printf("Pushing %v onto the stack", item)
	s.items = append(s.items, item)
	log.Printf("Stack after push: %v", s.items)
}

func (s *Stack[T]) PushMany(items []T) {
	log.Printf("Log before PushMany(%v): %v", items, s.items)

	for i := len(items)-1; i >= 0; i-- {
		s.items = append(s.items, items[i])
	}

	log.Printf("Log after PushMany(%v): %v", items, s.items)
}

func (s *Stack[T]) Grab() (T, error) {
	if len(s.items) == 0 {
		var defaultValue T
		return defaultValue, fmt.Errorf("Stack was empty")
	}
	buffer := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return buffer, nil
}

func (s *Stack[T]) GrabMany(count int) []T {
	buffer := []T{}
	for i := 0; i < count; i++ {
		item, err := s.Grab()
		if err != nil {
			break;
		}
		buffer = append(buffer, item)
	}

	if s.Version == CrateMover9000 {
		reversedBuffer := []T{}
	 	for i := len(buffer)-1; i >= 0; i-- {
	 		reversedBuffer = append(reversedBuffer, buffer[i])
	 	}
		return reversedBuffer
	}
	return buffer
}

func getStacks(input []string) []int {
	stacksNumbersLine := input[len(input)-1]
	numberOfStacks := strings.Split(strings.TrimSpace(strings.Replace(stacksNumbersLine, "   ", " ", -1)), " ")

	var result []int
	for _, v := range numberOfStacks {
		asInt, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}

		result = append(result, asInt)

	}

	return result
}

func indexWithStack(stack int) int {
	return 1 + (4 * (stack - 1))
}

func parseInstruction(input string) (int, int, int) {
	tokens := strings.Split(input, " ")

	quantity, err := strconv.Atoi(tokens[1])
	if err != nil {
		panic(err)
	}

	from, err := strconv.Atoi(tokens[3])
	if err != nil {
		panic(err)
	}

	to, err := strconv.Atoi(tokens[5])
	if err != nil {
		panic(err)
	}

	return quantity, from, to
}

func solve(readFile []byte, craneType CraneType, part int) {
	// Set up the stacks
	stacksAndInstructions := strings.Split(string(readFile), "\n\n")
	stacksHeader := strings.Split(stacksAndInstructions[0], "\n")

	stackIndices := getStacks(stacksHeader)
	log.Printf("Stack indicies: %v", stackIndices)

	stacks := make(map[int]Stack[string])

	// Initialize mapping for the stacks
	for _, v := range stackIndices {
		stacks[v] = NewStack[string](craneType)
	}

	// going throuh each stack from the bottom
	for i := len(stacksHeader)-2; i >= 0; i-- {
		for _, v := range stackIndices {
			index := indexWithStack(v)

			q := string(stacksHeader[i][index])
			if q != " " {
				stack := stacks[v]
				stack.Push(q)
				stacks[v] = stack
			}
		}
	}


	// Debug logging to inspect stack state
	log.Println("Inspecting stack state after creation")
	for _, v := range stackIndices {
		stack := stacks[v]
		stack.Print()
	}


	// Execute the instructions
	instructions := strings.Split(stacksAndInstructions[1], "\n")

	for _, instruction := range instructions {
		quantity, from, to := parseInstruction(instruction)

		stackFrom := stacks[from]
		stackTo := stacks[to]

		stackTo.PushMany(stackFrom.GrabMany(quantity))
		stacks[from] = stackFrom
		stacks[to] = stackTo

	}
	// Debug logging to inspect stack state
	log.Println("Inspecting stack state after instructions")
	for _, v := range stackIndices {
		stack := stacks[v]
		stack.Print()
	}

	// Collect the result for part1
	fmt.Printf("Result part %d: ", part)
	for _, v := range stackIndices {
		stack := stacks[v]
		crate, _ := stack.Grab()
		fmt.Printf("%s", crate)
	}
	fmt.Println()
}


func main() {
	if os.Getenv("LOGS") != "1" {
		log.SetOutput(io.Discard)
	}

	readFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	solve(readFile, CrateMover9000, 1)
	solve(readFile, CrateMover9001, 2)
}
