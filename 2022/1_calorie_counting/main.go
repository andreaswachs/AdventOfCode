package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	// Open the input file
	readFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Prepare to scan the file, line by line
	// credits: https://golangdocs.com/golang-read-file-line-by-line

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var sums []int
	var currentBlockSum int

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line == "" {
			sums = append(sums, currentBlockSum)
			currentBlockSum = 0
		} else {
			// The block keeps going on!
			num, _ := strconv.Atoi(line)
			currentBlockSum += num
		}
	}

	sums = append(sums, currentBlockSum)
	currentBlockSum = 0

	// Sorts the numbers in descending order
	sort.SliceStable(sums, func(i, j int) bool {
		return sums[j] < sums[i]
	})

	fmt.Printf("answer for #1: %d\n", sums[0])
	fmt.Printf("answer for #2: %d\n", sums[0]+sums[1]+sums[2])
}
