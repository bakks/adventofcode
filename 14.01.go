package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func readFile(filename string) []string {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
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

func parsePolymer(lines []string) (string, map[string]rune) {
	polymer := lines[0]
	rules := map[string]rune{}

	for i := 2; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		rules[fields[0]] = rune(fields[2][0])
	}

	return polymer, rules
}

func polymerize(polymer *string, rules map[string]rune) *string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%c", (*polymer)[0])

	for i := 0; i < len(*polymer)-1; i++ {
		pair := (*polymer)[i : i+2]
		target, ok := rules[pair]
		if ok {
			fmt.Fprintf(&builder, "%c", target)
		}
		fmt.Fprintf(&builder, "%c", pair[1])
	}

	x := builder.String()
	return &x
}

func scorePolymer(polymer *string) int {
	counts := map[rune]int{}

	// find frequency of each letter
	for _, r := range *polymer {
		_, ok := counts[r]
		if !ok {
			counts[r] = 0
		}
		counts[r]++
	}
	fmt.Println(counts)

	max := 0
	min := len(*polymer)

	for _, count := range counts {
		if count < min {
			min = count
		}
		if count > max {
			max = count
		}
	}

	return max - min
}

func main() {
	lines := readFile("14.txt")
	polymer, rules := parsePolymer(lines)
	pptr := &polymer

	for i := 0; i < 10; i++ {
		pptr = polymerize(pptr, rules)
		fmt.Printf("Iteration: %d, Length: %d\n", i, len(*pptr))
		if i < 6 {
			fmt.Println(*pptr)
		}
	}

	score := scorePolymer(pptr)
	fmt.Printf("Score: %d\n", score)
}
