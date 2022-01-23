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

type Grid struct {
	values [][]int
}

type GridPoint struct {
	X int
	Y int
}

func NewPoint(x, y int) GridPoint {
	return GridPoint{X: x, Y: y}
}

func (this GridPoint) Equals(that GridPoint) bool {
	return this.X == that.X && this.Y == that.Y
}

type PointSet struct {
	points []GridPoint
	i      int
}

func NewPointSet() *PointSet {
	return &PointSet{[]GridPoint{}, -1}
}

func (this *PointSet) Next() bool {
	this.i = this.i + 1
	return this.i < len(this.points)
}

func (this *PointSet) Get() GridPoint {
	return this.points[this.i]
}

func (this *PointSet) Size() int {
	return len(this.points)
}

func (this *PointSet) Contains(point GridPoint) bool {
	for _, p := range this.points {
		if p.Equals(point) {
			return true
		}
	}
	return false
}

func (this *PointSet) Add(newPoint GridPoint) bool {
	for _, point := range this.points {
		if point.Equals(newPoint) {
			return false
		}
	}

	this.points = append(this.points, newPoint)
	return true
}

func (this *PointSet) AddSet(newSet *PointSet) {
	for _, newPoint := range newSet.points {
		this.Add(newPoint)
	}
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

func (this *Grid) LowPoints() []GridPoint {
	points := []GridPoint{}
	xSize := len(this.values)
	ySize := len(this.values[0])

	for i := 0; i < xSize; i++ {
		for j := 0; j < ySize; j++ {
			lowPoint := true
			val := this.values[i][j]

			if i-1 >= 0 && this.values[i-1][j] <= val {
				lowPoint = false
			}
			if j-1 >= 0 && this.values[i][j-1] <= val {
				lowPoint = false
			}
			if i+1 < xSize && this.values[i+1][j] <= val {
				lowPoint = false
			}
			if j+1 < ySize && this.values[i][j+1] <= val {
				lowPoint = false
			}

			if lowPoint {
				points = append(points, GridPoint{X: i, Y: j})
			}
		}
	}

	return points
}

func (this *Grid) RiskLevels() []int {
	levels := []int{}
	lowPoints := this.LowPoints()

	for _, point := range lowPoints {
		risk := this.values[point.X][point.Y] + 1
		levels = append(levels, risk)
	}

	return levels
}

func (this *Grid) isValid(newPoint, currPoint, start GridPoint, fullSet *PointSet) bool {
	if newPoint.X < 0 || newPoint.Y < 0 {
		return false
	}
	if newPoint.X >= len(this.values) || newPoint.Y >= len(this.values[0]) {
		return false
	}

	if this.values[newPoint.X][newPoint.Y] == 9 {
		return false
	}

	if fullSet.Contains(newPoint) {
		return false
	}

	return true
}

func (this *Grid) BasinSize(start GridPoint) int {
	currSet := NewPointSet()
	currSet.Add(start)
	fullSet := NewPointSet()
	fullSet.Add(start)

	for true {
		newSet := NewPointSet()

		for currSet.Next() {
			point := currSet.Get()

			p1 := NewPoint(point.X+1, point.Y)
			if this.isValid(p1, point, start, fullSet) {
				newSet.Add(p1)
			}
			p2 := NewPoint(point.X-1, point.Y)
			if this.isValid(p2, point, start, fullSet) {
				newSet.Add(p2)
			}
			p3 := NewPoint(point.X, point.Y+1)
			if this.isValid(p3, point, start, fullSet) {
				newSet.Add(p3)
			}
			p4 := NewPoint(point.X, point.Y-1)
			if this.isValid(p4, point, start, fullSet) {
				newSet.Add(p4)
			}
		}

		if newSet.Size() == 0 {
			return fullSet.Size()
		}

		fullSet.AddSet(newSet)
		currSet = newSet
	}

	return 0
}

func (this *Grid) BasinSizes() []int {
	lowPoints := this.LowPoints()
	sizes := make([]int, len(lowPoints))

	for i, lowPoint := range lowPoints {
		size := this.BasinSize(lowPoint)
		sizes[i] = size
	}

	return sizes
}

func (this *Grid) Mult3Basins() int {
	sizes := this.BasinSizes()
	sort.Ints(sizes)

	top3 := sizes[len(sizes)-3:]
	return top3[0] * top3[1] * top3[2]
}

func getTotalRisk(grid *Grid) int {
	riskLevels := grid.RiskLevels()

	sum := 0
	for _, riskLevel := range riskLevels {
		sum += riskLevel
	}

	return sum
}

func main() {
	lines := readFile("09.txt")
	grid := parseGrid(lines)

	totalRisk := getTotalRisk(grid)
	fmt.Printf("Total risk: %d\n", totalRisk)

	basins := grid.Mult3Basins()
	fmt.Printf("Basin sizes multiplied: %d\n", basins)
}
