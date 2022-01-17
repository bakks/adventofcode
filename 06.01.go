package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

type LanternTracker struct {
	fish []int
}

const ResetAge = 6
const BirthAge = 8

func NewLanternTracker(fishAges []int) *LanternTracker {
	fish := make([]int, BirthAge+1)
	for _, currFish := range fishAges {
		fish[currFish]++
	}

	return &LanternTracker{fish}
}

func (this *LanternTracker) IncrementDay() {
	zeros := this.fish[0]
	this.fish = append(this.fish[1:], zeros)
	this.fish[ResetAge] += zeros
}

func (this *LanternTracker) Count() int {
	n := 0

	for _, numFish := range this.fish {
		n += numFish
	}

	return n
}

func (this *LanternTracker) Print() {
	fmt.Println(this.fish)
}

func parseFish(in []string) *LanternTracker {
	line := in[0]
	fields := strings.Split(line, ",")

	fish := make([]int, len(fields))
	for i, strAge := range fields {
		age, err := strconv.Atoi(strAge)
		if err != nil {
			log.Fatal(err)
		}
		fish[i] = age
	}

	return NewLanternTracker(fish)
}

func main() {
	lines := readFile("06.txt")
	fish := parseFish(lines)

	for day := 0; day < 256; day++ {
		fish.IncrementDay()
		fish.Print()
		fmt.Printf("After %d days: %d\n", day+1, fish.Count())
	}
}
