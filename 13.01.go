package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

type Grid struct {
	grid [][]bool
}

type Point struct {
	X int
	Y int
}

type Fold struct {
	FoldY  bool
	Crease int
}

func NewGrid(points []Point) *Grid {
	maxY := 0
	maxX := 0

	for _, point := range points {
		if point.Y > maxY {
			maxY = point.Y
		}
		if point.X > maxX {
			maxX = point.X
		}
	}

	maxY++
	maxX++

	grid := make([][]bool, maxY)
	for i := 0; i < maxY; i++ {
		grid[i] = make([]bool, maxX)
	}

	for _, point := range points {
		grid[point.Y][point.X] = true
	}

	return &Grid{grid}
}

func (this *Grid) Print() {
	for _, line := range this.grid {
		for _, point := range line {
			x := "."
			if point {
				x = "#"
			}
			fmt.Printf(x)
		}
		fmt.Printf("\n")
	}
}

func (this *Grid) Fold(fold Fold) {
	if fold.FoldY {
		for y := fold.Crease + 1; y < len(this.grid); y++ {
			for x, val := range this.grid[y] {
				if val {
					newY := fold.Crease - (y - fold.Crease)
					this.grid[newY][x] = true
				}
			}
		}

		this.grid = this.grid[0:fold.Crease]

	} else {
		for x := fold.Crease + 1; x < len(this.grid[0]); x++ {
			for y, _ := range this.grid {
				val := this.grid[y][x]
				if val {
					newX := fold.Crease - (x - fold.Crease)
					this.grid[y][newX] = true
				}
			}
		}

		for y, row := range this.grid {
			this.grid[y] = row[0:fold.Crease]
		}
	}
}

func (this *Grid) Count() int {
	sum := 0

	for _, row := range this.grid {
		for _, point := range row {
			if point {
				sum++
			}
		}
	}

	return sum
}

func parseOrigami(lines []string) ([]Point, []Fold) {
	index := 0
	points := []Point{}
	folds := []Fold{}

	for i, line := range lines {
		index = i
		if line == "" {
			break
		}

		fields := strings.Split(line, ",")
		x, _ := strconv.Atoi(fields[0])
		y, _ := strconv.Atoi(fields[1])
		point := Point{X: x, Y: y}
		points = append(points, point)
	}

	for i := index + 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		assignment := strings.Split(fields[2], "=")
		foldY := assignment[0] == "y"
		crease, _ := strconv.Atoi(assignment[1])
		fold := Fold{foldY, crease}
		folds = append(folds, fold)
	}

	return points, folds
}

func main() {
	lines := readFile("13.txt")
	points, folds := parseOrigami(lines)

	grid := NewGrid(points)

	for _, fold := range folds {
		grid.Fold(fold)
	}
	grid.Print()

	fmt.Printf("Count: %d\n", grid.Count())
}
