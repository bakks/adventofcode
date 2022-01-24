package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

type Grid struct {
	values [][]int
}

func (this *Grid) Size() int {
	return len(this.values) * len(this.values[0])
}

func parseGrid(in []string) *Grid {
	vals := [][]int{}

	for _, line := range in {
		valLine := []int{}
		for _, r := range line {
			n := int(r) - int('0')
			valLine = append(valLine, n)
		}

		vals = append(vals, valLine)
	}

	return &Grid{vals}
}

func (this *Grid) safeInc(i, j int) {
	if i < 0 || j < 0 {
		return
	}
	if i >= len(this.values) || j >= len(this.values[i]) {
		return
	}
	if this.values[i][j] == 0 {
		return
	}
	this.values[i][j]++
}

func (this *Grid) RunStep() int {
	// initial increment
	for i, row := range this.values {
		for j, _ := range row {
			this.values[i][j]++
		}
	}

	flashes := 0

	updates := true
	for updates {
		updates = false

		for i, row := range this.values {
			for j, _ := range row {
				val := this.values[i][j]

				if val > 9 {
					// flash
					updates = true
					flashes++
					this.safeInc(i-1, j-1)
					this.safeInc(i-1, j)
					this.safeInc(i-1, j+1)
					this.safeInc(i, j-1)
					this.safeInc(i, j+1)
					this.safeInc(i+1, j-1)
					this.safeInc(i+1, j)
					this.safeInc(i+1, j+1)
					this.values[i][j] = 0
				}
			}
		}

	}

	return flashes
}

func (this *Grid) Print() {
	for _, row := range this.values {
		for _, i := range row {
			fmt.Printf("%d ", i)
		}
		fmt.Printf("\n")
	}
}

func main() {
	lines := readFile("11.txt")
	grid := parseGrid(lines)

	grid.Print()
	fmt.Println()
	flashes := 0

	for i := 0; i < 100; i++ {
		flashes += grid.RunStep()
	}

	grid.Print()
	fmt.Printf("Flashes: %d\n", flashes)

	grid = parseGrid(lines)
	gridSize := grid.Size()
	for i := 0; i < 10000; i++ {
		flashes := grid.RunStep()

		if flashes == gridSize {
			fmt.Printf("Synchronized flash at iteration: %d\n", i+1)
			break
		}
	}
}
