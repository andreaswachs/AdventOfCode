package main

import (
	"bufio"
	"fmt"
	"os"
)

// toPriority calculates priority for badges/item types represeented as runes
// Letters a-z has priority 1..26 and letters A-Z has priority 27..53
func toPriority(r rune) int {
	if r >= 97 {
		// Calculate for lowercase rune
		return int(r) - 96
	}

	// Calculate for uppercase rune
	return (int(r) - 64) + 26
}

// calculateBadgePriority calculates the common badge between all the items in the provided rucksacks
// This is done using a slice of hash tables (due to _acceptable_ lookup times),
// where we mark all occurances for each badge within each rucksack
// This is done for all but the last rucksacks. The last rucksack is gone through while
// using the table to see which of the badges in the last rucksack appears in all of the other rucksacks
func calculateBadgePriority(rucksacks []string) int {
	// Generalize the problem to N many rucksacks

	// Setting up the "table"
	var table []map[rune]bool
	for i := 0; i < 3; i++ {
		table = append(table, make(map[rune]bool))
	}

	for i, rucksack := range rucksacks {
		// If we've reached the final rucksack, we skip that one
		if i >= len(rucksack)-1 {
			break
		}
		// We account for all the types in each rucksack, make note of which we encounter
		for _, v := range rucksack {
			table[i][v] = true
		}
	}

	// Go through the last rucksack
	brokeLoop := false
	for _, v := range rucksacks[len(rucksacks)-1] {
		for _, t := range table {
			if !t[v] {
				brokeLoop = true
				break
			}
		}
		if !brokeLoop {
			return toPriority(v)
		}
		brokeLoop = false
	}

	return 0
}

// calculatePriorty calculates the total value of the priority for the item type
// that is found within the two compartments of the rucksack.
// This is done using a hash map (rune -> bool), that can account for occurences
// of a given rune ("item type"). This means that if we insert all the runes found
// in the first compartment, we can then check to see when we hit a match within the
// second compartment, since bool's default value is false, we can trust that if
// map[key] = true, then that rune exists in the first compartment
func calculatePriority(rucksack string) int {
	firstCompartment := make(map[rune]bool)

	for i, v := range rucksack {
		if i <= (len(rucksack)-1)/2 {
			// Mark all occurances of a given item type in the first compartment
			firstCompartment[v] = true
		} else {
			// Check occurance of a given item type in the second compartment
			if firstCompartment[v] {
				return toPriority(v)
			}
		}
	}

	return 0
}

func main() {
	// Open the input file with name given in first argument
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Get ready to scan the input
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var result1 int
	var result2 int

	// Define a helper function to easily solve part 2
	calcBadgePrio := func(sliceBuffer *[]string, buffer string) int {
		if len(*sliceBuffer) == 3 {
			sliceBufferBuffer := *sliceBuffer
			*sliceBuffer = []string{buffer}
			return calculateBadgePriority(sliceBufferBuffer)
		} else {
			*sliceBuffer = append(*sliceBuffer, buffer)
			return 0
		}
	}

	// For part 2, we have to work on 3 input lines at the time, so we store a slice
	var sliceBuffer []string
	for {
		if fileScanner.Scan() {
			buffer := fileScanner.Text()

			// Capture points for part 1
			result1 += calculatePriority(buffer)

			// Capture points for part 2
			result2 += calcBadgePrio(&sliceBuffer, buffer)
			continue
		}

		// We set the for loop up like this as this will get triggered the first time that the input is empty,
		// since we need to make sure to call calcBadgePrio() for the last slice of buffered inputs
		// Capture points for part 2 for the last set of 3 rucksacks
		result2 += calcBadgePrio(&sliceBuffer, "")
		break
	}

	fmt.Printf("Part 1: %d\n", result1)
	fmt.Printf("Part 2: %d\n", result2)
}
