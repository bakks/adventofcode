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

type VentMap struct {
	vents [][]int
}

func NewVentMap(size int) VentMap {
	vents := make([][]int, size)
	for i, _ := range vents {
		vents[i] = make([]int, size)
	}

	return VentMap{vents}
}

func (this VentMap) AddVent(x1, y1, x2, y2 int) {
	if x1 == x2 {
		yMin := min(y1, y2)
		yMax := max(y1, y2)

		for y := yMin; y <= yMax; y++ {
			this.vents[x1][y]++
		}
	} else if y1 == y2 {
		xMin := min(x1, x2)
		xMax := max(x1, x2)

		for x := xMin; x <= xMax; x++ {
			this.vents[x][y1]++
		}
	} else {
		xMin := min(x1, x2)
		xMax := max(x1, x2)
		yMin := min(y1, y2)
		yMax := max(y1, y2)
		fmt.Printf("%d,%d -> %d,%d\n", x1, y1, x2, y2)

		if yMax-yMin != xMax-xMin {
			log.Fatalf("Not diagonal: %d,%d -> %d,%d\n", x1, y1, x2, y2)
		}

		x, y := x1, y1
		for x != x2 {
			this.vents[x][y]++
			if x1 < x2 {
				x++
			} else {
				x--
			}
			if y1 < y2 {
				y++
			} else {
				y--
			}
		}
		this.vents[x][y]++

	}
}

func (this VentMap) Print(distance int) {
	for x := 0; x < distance; x++ {
		for y := 0; y < distance; y++ {
			n := this.vents[y][x]
			if n == 0 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%d", this.vents[y][x])
			}
		}
		fmt.Printf("\n")
	}
}

func (this VentMap) CountDangers() int {
	n := 0

	for _, ventRow := range this.vents {
		for _, ventCount := range ventRow {
			if ventCount > 1 {
				n++
			}
		}
	}

	return n
}

func parseVents(in []string) VentMap {
	vents := NewVentMap(1000)

	for _, line := range in {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			log.Fatal(line)
		}

		one := strings.Split(fields[0], ",")
		two := strings.Split(fields[2], ",")
		x1, err1 := strconv.Atoi(one[0])
		y1, err2 := strconv.Atoi(one[1])
		x2, err3 := strconv.Atoi(two[0])
		y2, err4 := strconv.Atoi(two[1])

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			log.Fatal(line)
		}

		vents.AddVent(x1, y1, x2, y2)
	}

	return vents
}

func main() {
	lines := readFile("05.txt")
	vents := parseVents(lines)

	fmt.Println(vents.CountDangers())
	vents.Print(10)
}
