package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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

type Grid [][]int

func parseGrid(in []string) Grid {
	vals := [][]int{}

	for _, line := range in {
		valLine := []int{}
		for _, r := range line {
			n := int(r) - int('0')
			valLine = append(valLine, n)
		}

		vals = append(vals, valLine)
	}

	return vals
}

func (this Grid) Print() {
	for y := 0; y < len(this); y++ {
		for x := 0; x < len(this[y]); x++ {
			fmt.Printf("%2d ", this[y][x])
		}
		fmt.Printf("\n")
	}
}

func newGrid(ySize, xSize, val int) Grid {
	grid := make([][]int, ySize)

	for y := 0; y < ySize; y++ {
		grid[y] = make([]int, xSize)
		for x := 0; x < xSize; x++ {
			grid[y][x] = val
		}
	}

	return grid
}

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

type PointQueue struct {
	q []Point
}

func NewPointQueue() *PointQueue {
	return &PointQueue{[]Point{}}
}

func (this *PointQueue) Push(p Point) {
	this.q = append(this.q, p)
}

func (this *PointQueue) Pop() Point {
	if len(this.q) == 0 {
		log.Fatalf("empty queue")
	}

	p := this.q[0]
	this.q = this.q[1:]

	return p
}

func (this PointQueue) Size() int {
	return len(this.q)
}

func pathFind(grid Grid) int {
	m := newGrid(len(grid), len(grid[0]), math.MaxInt64)
	queue := NewPointQueue()
	queue.Push(NewPoint(0, 0))
	m[0][0] = 0

	for queue.Size() > 0 {
		point := queue.Pop()

		neighbors := []Point{
			Point{point.X - 1, point.Y},
			Point{point.X + 1, point.Y},
			Point{point.X, point.Y - 1},
			Point{point.X, point.Y + 1}}

		for _, neighbor := range neighbors {
			if neighbor.Y < 0 || neighbor.X < 0 || neighbor.Y >= len(grid) || neighbor.X >= len(grid[0]) {
				continue
			}

			newCost := m[point.Y][point.X] + grid[neighbor.Y][neighbor.X]
			if m[neighbor.Y][neighbor.X] > newCost {
				m[neighbor.Y][neighbor.X] = newCost
				queue.Push(neighbor)
			}
		}
	}

	return m[len(grid)-1][len(grid[0])-1]
}

func main() {
	lines := readFile("15.txt")
	grid := parseGrid(lines)
	pathScore := pathFind(grid)
	fmt.Printf("Lowest scoring path: %d\n", pathScore)
}
