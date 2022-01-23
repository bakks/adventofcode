package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
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

func isOpener(r rune) bool {
	return r == '(' || r == '{' || r == '<' || r == '['
}

var openToClose = map[rune]rune{
	'(': ')',
	'{': '}',
	'<': '>',
	'[': ']',
}

var closeScores = map[rune]int{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4,
}

// if corrupted, return the syntax error rune
// if incomplete, return the current stack
func parseLine(line string) (rune, []rune) {
	stack := []rune{}

	for _, r := range line {
		if isOpener(r) {
			// push
			stack = append(stack, r)
		} else {
			// pop and check
			if len(stack) == 0 {
				//fmt.Printf("Expected nothing, found %c\n", r)
				return r, nil
			}
			popped := stack[len(stack)-1]
			if openToClose[popped] != r {
				//fmt.Printf("Expected %c, found %c\n", openToClose[popped], r)
				return r, nil
			}

			stack = stack[:len(stack)-1]
		}
	}

	return '0', stack
}

func findSyntaxErrors(in []string) []rune {
	errors := []rune{}
	for _, line := range in {
		errRune, _ := parseLine(line)
		if errRune != '0' {
			errors = append(errors, errRune)
		}
	}

	return errors
}

func calcSyntaxScore(runes []rune) int {
	sum := 0
	for _, r := range runes {
		switch r {
		case ')':
			sum += 3
		case '}':
			sum += 1197
		case ']':
			sum += 57
		case '>':
			sum += 25137
		}
	}

	return sum
}

func calcCompletionScore(lines []string) int {
	completions := [][]rune{}
	scores := []int{}

	for _, line := range lines {
		_, stack := parseLine(line)

		if stack == nil {
			continue
		}

		completion := []rune{}

		for i := len(stack) - 1; i >= 0; i-- {
			r := openToClose[stack[i]]
			completion = append(completion, r)
		}

		fmt.Printf("%s : %s\n", string(stack), string(completion))

		completions = append(completions, completion)

		total := 0
		for _, r := range completion {
			total *= 5
			total += closeScores[r]
		}
		scores = append(scores, total)
	}

	sort.Ints(scores)
	return scores[len(scores)/2]
}

func main() {
	lines := readFile("10.txt")
	errorRunes := findSyntaxErrors(lines)
	score := calcSyntaxScore(errorRunes)
	fmt.Printf("Corruption score: %d\n", score)

	completionScore := calcCompletionScore(lines)
	fmt.Printf("Completion score: %d\n", completionScore)
}
