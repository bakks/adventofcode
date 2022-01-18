package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func main() {
	lines := readFile("08.txt")
	_, outputSignals := parseSignals(lines)

	n := 0
	for _, outputSignal := range outputSignals {
		n += outputSignal.UniqueSignals()
	}

	fmt.Println(n)
}
