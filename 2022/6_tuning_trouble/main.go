package main

import (
	"fmt"
	"os"
	"strings"
)

// packetType describes which kind of "packet" were looking for
// this differentiates part 1 and 2
type packetType uint8

const (
	startOfPacket packetType = iota
	startOfMessage
)

// uniquenesCheckBuffer is a mappign from runes 'a' -> 'z', that can tell us with reasonable lookup time
// wheter a sequence of runes contain the same runes. You need to remember to reset it between checks!!
var uniquenessCheckBuffer map[rune]bool

func initCheckBuffer() {
	uniquenessCheckBuffer = make(map[rune]bool, 25)
}

func resetCheckBuffer() {
	for i := 'a'; i <= 'z'; i++ {
		uniquenessCheckBuffer[i] = false
	}
}

func checkUniqueness(r rune) bool {
	if !uniquenessCheckBuffer[r] {
		uniquenessCheckBuffer[r] = true
		return true
	}

	return false
}

// For part 1, we need to look for unique sequences of runes with a length of 4, and
// for part 2, we need to look for unique sequences of runes with a length of 14,
// so I've deduceted 1 from the result due to my method of solving this in solve()
func getPacketOffset(pt packetType) int {
	switch pt {
	case startOfPacket:
		return 3
	case startOfMessage:
		return 13
	}

	// In an Advent of Code world, where input is perfect,
	// this should never be reached on purpose
	panic(fmt.Errorf("Unknown packet type"))
}

// solve solves the problem depending on packettype that were looking for (depending on part 1 or 2)
// Solve traverses the input from left to right and looks at 4/14 letters at the time.
// If it finds a sequence of 4/14 letters that are unique in that sequence, it returns the index of the
// character after this sequence, since that sequence either signalled the start of the packet or message.
func solve(input string, pt packetType) (int, error) {
	bound := len(input) - 4 // should this be 4?
	additionalPackets := getPacketOffset(pt)

	// Go thorough all the characeters but the last 4 ones
	for i := 0; i <= bound; i++ {
		resetCheckBuffer()
		// Go through runes i, .. i+4
		for j := i; j <= i+additionalPackets; j++ {
			r := rune(input[j]) // to inspect current rune for debugging
			if checkUniqueness(r) {
				if j == i+additionalPackets {
					return j + 1, nil
				}
			} else {
				break
			}
		}
	}

	return 0, fmt.Errorf("could not find solution")
}

func main() {
	initCheckBuffer()

	contents, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	input := strings.TrimSpace(string(contents))

	result1, err := solve(input, startOfPacket)
	if err != nil {
		panic(err)
	}

	result2, err := solve(input, startOfMessage)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 2: %d\n", result2)
	fmt.Printf("Part 1: %d\n", result1)
}
