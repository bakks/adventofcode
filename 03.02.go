package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func readFile(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	lines := []string{}
	s := bufio.NewScanner(f)

	for s.Scan() {
		err = s.Err()
		if err != nil {
			log.Fatal(err)
		}

		lines = append(lines, s.Text())
	}

	return lines
}

func parseDiagnostic(in []string) [][]bool {
	output := [][]bool{}

	for _, line := range in {
		outLine := []bool{}

		for _, x := range line {
			switch x {
			case '0':
				outLine = append(outLine, false)
			case '1':
				outLine = append(outLine, true)
			default:
				log.Panic(x)
			}
		}

		output = append(output, outLine)
	}

	return output
}

const DiagnosticLineLength = 12

func mostCommon(diag [][]bool, index int) bool {
	t, f := 0, 0
	for _, diagLine := range diag {
		if diagLine[index] {
			t++
		} else {
			f++
		}
	}

	return t >= f
}

func diagLineToInt(diagLine []bool) int {
	n := uint(0)
	for _, bit := range diagLine {
		n = n << 1
		if bit {
			n = n | 1
		}
	}

	return int(n)
}

func diagRating(diag [][]bool, searchMostCommon bool) int {
	prevSet := diag

	for i := 0; i < len(diag[0]) && len(prevSet) > 1; i++ {
		nextSet := [][]bool{}
		searchBit := mostCommon(prevSet, i)
		if !searchMostCommon {
			searchBit = !searchBit
		}

		for _, diagLine := range prevSet {
			if diagLine[i] == searchBit {
				nextSet = append(nextSet, diagLine)
			}
		}

		prevSet = nextSet
	}

	if len(prevSet) != 1 {
		log.Fatalf("Expected single final number, had %d\n", len(prevSet))
	}

	return diagLineToInt(prevSet[0])
}

func oxygenGeneratorRating(diag [][]bool) int {
	return diagRating(diag, true)
}

func c02ScrubberRating(diag [][]bool) int {
	return diagRating(diag, false)
}

var TestData = [][]bool{
	{false, false, true, false, false},
	{true, true, true, true, false},
	{true, false, true, true, false},
	{true, false, true, true, true},
	{true, false, true, false, true},
	{false, true, true, true, true},
	{false, false, true, true, true},
	{true, true, true, false, false},
	{true, false, false, false, false},
	{true, true, false, false, true},
	{false, false, false, true, false},
	{false, true, false, true, false},
}

func main() {
	lines := readFile("03.txt")
	diagnostic := parseDiagnostic(lines)

	oxyRating := oxygenGeneratorRating(diagnostic)
	c02Rating := c02ScrubberRating(diagnostic)
	lifeRating := oxyRating * c02Rating

	fmt.Printf("Oxygen Generator Rating: %d\n", oxyRating)
	fmt.Printf("C02 Scrubber Rating: %d\n", c02Rating)
	fmt.Printf("Life Support Rating: %d\n", lifeRating)
}
