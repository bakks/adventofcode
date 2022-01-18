package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minArr(input []int) int {
	min := -1
	for _, x := range input {
		if min == -1 || x < min {
			min = x
		}
	}

	return min
}

func maxArr(input []int) int {
	max := -1
	for _, x := range input {
		if max == -1 || x > max {
			max = x
		}
	}

	return max
}

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

type SignalSet struct {
	signals []string
}

func NewSignalSet(signals []string) *SignalSet {
	return &SignalSet{signals: signals}
}

func (this *SignalSet) UniqueSignals() int {
	n := 0

	for _, sig := range this.signals {
		switch len(sig) {
		case 2, 3, 4, 7: // 1, 7, 4, 8
			n++
		}
	}

	return n
}

func contains(in []rune, r rune) bool {
	for _, curr := range in {
		if curr == r {
			return true
		}
	}
	return false
}

func remove(in []rune, toRemove []rune) []rune {
	newSlice := []rune{}
	for _, r := range in {
		if !contains(toRemove, r) {
			newSlice = append(newSlice, r)
		}
	}
	return newSlice
}

func equal(a []rune, b []rune) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}

	return true
}

var ZERO = []rune{'a', 'b', 'c', 'e', 'f', 'g'}
var ONE = []rune{'c', 'f'}
var TWO = []rune{'a', 'c', 'd', 'e', 'g'}
var THREE = []rune{'a', 'c', 'd', 'f', 'g'}
var FOUR = []rune{'b', 'c', 'd', 'f'}
var FIVE = []rune{'a', 'b', 'd', 'f', 'g'}
var SIX = []rune{'a', 'b', 'd', 'e', 'f', 'g'}
var SEVEN = []rune{'a', 'c', 'f'}
var EIGHT = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
var NINE = []rune{'a', 'b', 'c', 'd', 'f', 'g'}

func missingFrom(in []rune) []rune {
	all := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	return remove(all, in)
}

func removeOtherRunes(in, removals []rune) []rune {
	newSlice := []rune{}
	for _, r := range in {
		if contains(removals, r) {
			newSlice = append(newSlice, r)
		}
	}

	return newSlice
}

func printPossibles(in map[rune][]rune) {
	for k, v := range in {
		fmt.Printf("%c: ", k)
		for _, r := range v {
			fmt.Printf("%c,", r)
		}
		fmt.Printf("\n")
	}
}

func (this *SignalSet) ApplyConfiguration(config map[rune]rune) []int {
	out := []int{}

	for _, sig := range this.signals {
		mappedSig := []rune{}
		for _, r := range []rune(sig) {
			mappedSig = append(mappedSig, config[r])
		}

		sort.SliceStable(mappedSig, func(i, j int) bool {
			return mappedSig[i] < mappedSig[j]
		})

		n := -1
		switch {
		case equal(mappedSig, ZERO):
			n = 0
		case equal(mappedSig, ONE):
			n = 1
		case equal(mappedSig, TWO):
			n = 2
		case equal(mappedSig, THREE):
			n = 3
		case equal(mappedSig, FOUR):
			n = 4
		case equal(mappedSig, FIVE):
			n = 5
		case equal(mappedSig, SIX):
			n = 6
		case equal(mappedSig, SEVEN):
			n = 7
		case equal(mappedSig, EIGHT):
			n = 8
		case equal(mappedSig, NINE):
			n = 9
		}

		if n == -1 {
			log.Fatal("bummmmp")
		}

		out = append(out, n)
	}

	return out
}

func (this *SignalSet) DeduceConfiguration() map[rune]rune {
	possible := map[rune][]rune{
		'a': []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
		'b': []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
		'c': []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
		'd': []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
		'e': []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
		'f': []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
		'g': []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
	}

	for i := 0; i < 10; i++ {
		for _, sig := range this.signals {
			sigRunes := []rune(sig)
			switch len(sig) {
			case 2: // this is a 1 digit
				// remove other mapping possibilities for c and f than these two
				possible['a'] = remove(possible['a'], sigRunes)
				possible['b'] = remove(possible['b'], sigRunes)
				possible['c'] = removeOtherRunes(possible['c'], sigRunes)
				possible['d'] = remove(possible['d'], sigRunes)
				possible['e'] = remove(possible['e'], sigRunes)
				possible['f'] = removeOtherRunes(possible['f'], sigRunes)
				possible['g'] = remove(possible['g'], sigRunes)

			case 3: // this is a 7 digit
				if len(possible['c']) <= 2 && len(possible['f']) <= 2 { // if we've narrowed down with a 1 digit
					// then the third rune not found must be 'a' (top of 7)
					aRune := remove(sigRunes, append(possible['c'], possible['f']...))
					if len(aRune) != 1 {
						log.Fatal("bad")
					}
					possible['a'] = aRune
					possible['b'] = remove(possible['b'], aRune)
					possible['c'] = remove(possible['c'], aRune)
					possible['d'] = remove(possible['d'], aRune)
					possible['e'] = remove(possible['e'], aRune)
					possible['f'] = remove(possible['f'], aRune)
					possible['g'] = remove(possible['g'], aRune)
				}

			case 4: // this is a 4 digit
				if len(possible['c']) <= 2 && len(possible['f']) <= 2 { // if we've narrowed down with a 1 digit
					// then the other runes in the signal must be b or d
					bdRunes := remove(sigRunes, append(possible['c'], possible['f']...))
					possible['a'] = remove(possible['a'], bdRunes)
					possible['b'] = removeOtherRunes(possible['b'], bdRunes)
					possible['c'] = remove(possible['c'], bdRunes)
					possible['d'] = removeOtherRunes(possible['d'], bdRunes)
					possible['e'] = remove(possible['e'], bdRunes)
					possible['f'] = remove(possible['f'], bdRunes)
					possible['g'] = remove(possible['g'], bdRunes)
				}

			case 5: // digits where the 3 horizontal bars are guaranteed (2, 3, 5)
				missing := missingFrom(sigRunes)
				possible['a'] = remove(possible['a'], missing)
				possible['d'] = remove(possible['d'], missing)
				possible['g'] = remove(possible['g'], missing)

			case 6:
				missing := missingFrom(sigRunes)
				possible['f'] = remove(possible['f'], missing)
			}

			// if we have a definitive mapping, remove that option from other possibilities
			for k, options := range possible {
				if len(options) == 1 {
					for k2, options2 := range possible {
						if k != k2 {
							possible[k2] = remove(options2, options)
						}
					}
				}
			}
		}
	}

	config := map[rune]rune{}
	for k, v := range possible {
		if len(v) != 1 {
			log.Fatalf("did not deduce!")
		}
		config[v[0]] = k
	}
	printPossibles(possible)

	return config
}

func delimiterPosition(fields []string) int {
	for i, field := range fields {
		if field == "|" {
			return i
		}
	}

	log.Fatal("not found")
	return -1
}

func parseSignals(in []string) ([]*SignalSet, []*SignalSet) {
	configSignals := []*SignalSet{}
	outputSignals := []*SignalSet{}

	for _, line := range in {
		fields := strings.Fields(line)
		delimiter := delimiterPosition(fields)

		c := fields[:delimiter]
		o := fields[delimiter+1:]
		configSignals = append(configSignals, NewSignalSet(c))
		outputSignals = append(outputSignals, NewSignalSet(o))
	}

	return configSignals, outputSignals
}

func toInt(in []int) int {
	x := 0

	for _, n := range in {
		x = x*10 + n
	}

	return x
}

func main() {
	lines := readFile("08.txt")
	configSignals, outputSignals := parseSignals(lines)

	total := 0

	for i, configSignal := range configSignals {
		config := configSignal.DeduceConfiguration()
		output := outputSignals[i].ApplyConfiguration(config)
		fmt.Println(outputSignals[i])
		fmt.Println(output)
		total += toInt(output)
	}

	fmt.Printf("Total: %d\n", total)
}
