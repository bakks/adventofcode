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

func parseCrabs(in []string) []int {
	fields := strings.Split(in[0], ",")
	crabs := []int{}

	for _, field := range fields {
		crab, err := strconv.Atoi(field)
		if err != nil {
			log.Fatal(err)
		}

		crabs = append(crabs, crab)
	}

	return crabs
}

func calculateFuelCost(crabs []int, target int) int {
	cost := 0

	for _, crab := range crabs {
		cost += abs(crab - target)
	}

	return cost
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func calculateBestTarget(crabs []int) int {
	min := minArr(crabs)
	max := maxArr(crabs)

	lowestCostTarget := -1
	lowestCost := -1

	for i := min; i <= max; i++ {
		cost := calculateFuelCost(crabs, i)
		if lowestCost == -1 || cost < lowestCost {
			lowestCost = cost
			lowestCostTarget = i
		}

		fmt.Printf("target: %5d,  cost: %5d,  lowest cost: %5d\n", i, cost, lowestCost)
	}

	return lowestCostTarget
}

func main() {
	lines := readFile("07.txt")
	crabs := parseCrabs(lines)

	fmt.Println(crabs)
	fmt.Println(calculateBestTarget(crabs))
}
