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

type LetterFreq struct {
	freq map[rune]int
}

func NewLetterFreq() *LetterFreq {
	return &LetterFreq{map[rune]int{}}
}

func (this *LetterFreq) Print() {
	for k, v := range this.freq {
		fmt.Printf("%c:%d, ", k, v)
	}
	fmt.Printf("\n")
}

func (this *LetterFreq) Copy() *LetterFreq {
	newMap := map[rune]int{}
	for k, v := range this.freq {
		newMap[k] = v
	}

	return &LetterFreq{newMap}
}

func (this *LetterFreq) Increment(r rune) {
	val, ok := this.freq[r]
	if !ok {
		this.freq[r] = 1
	} else {
		this.freq[r] = val + 1
	}
}

func (this *LetterFreq) Decrement(r rune) {
	_, ok := this.freq[r]
	if !ok {
		log.Fatalf("decrement when it doesn't exist")
	}
	this.freq[r]--
}

func (this *LetterFreq) Merge(that *LetterFreq) {
	for k, v := range that.freq {
		oldVal, _ := this.freq[k]
		this.freq[k] = oldVal + v
	}
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

func (this *LetterFreq) Score() int {
	max := 0
	min := -1

	for _, count := range this.freq {
		if count < min || min == -1 {
			min = count
		}
		if count > max {
			max = count
		}
	}

	return max - min
}

var memo = map[string]*LetterFreq{}

func memoPoly(pair string, depth int, rules map[string]rune) *LetterFreq {
	key := fmt.Sprintf("%s%d", pair, depth)
	result, ok := memo[key]

	if ok {
		return result.Copy()
	}

	newResult := poly(pair, depth, rules)
	memo[key] = newResult.Copy()
	return newResult
}

func poly(pair string, depth int, rules map[string]rune) *LetterFreq {
	if depth < 0 {
		freq := NewLetterFreq()
		return freq
	}

	target, ok := rules[pair]
	if !ok {
		log.Fatalf("No rule for %s\n", pair)
	}

	pair1 := string(pair[0]) + string(target)
	pair2 := string(target) + string(pair[1])
	freq1 := memoPoly(pair1, depth-1, rules)
	freq2 := memoPoly(pair2, depth-1, rules)
	freq1.Merge(freq2)
	freq1.Increment(target)

	return freq1
}

func polymerize(polymer string, rules map[string]rune, steps int) *LetterFreq {
	freq := NewLetterFreq()

	for i := 0; i < len(polymer)-1; i++ {
		pair := polymer[i : i+2]
		newFreq := poly(pair, steps-1, rules)
		freq.Merge(newFreq)
		freq.Increment(rune(pair[0]))
	}

	freq.Increment(rune(polymer[len(polymer)-1]))

	return freq
}

func main() {
	lines := readFile("14.txt")
	polymer, rules := parsePolymer(lines)
	freq := polymerize(polymer, rules, 40)
	freq.Print()
	score := freq.Score()
	fmt.Printf("Score: %d\n", score)
}
