package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
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

type Target struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

func (this *Target) hit(x, y int) bool {
	return x >= this.x1 && x <= this.x2 && y >= this.y1 && y <= this.y2
}

func parseTarget(line string) *Target {
	fields := getCaptures(line)
	x1, _ := strconv.Atoi(fields["x1"])
	x2, _ := strconv.Atoi(fields["x2"])
	y1, _ := strconv.Atoi(fields["y1"])
	y2, _ := strconv.Atoi(fields["y2"])

	return &Target{
		x1: x1,
		x2: x2,
		y1: y1,
		y2: y2,
	}
}

func getCaptures(line string) map[string]string {
	var compRegEx = regexp.MustCompile(`target area: x=(?P<x1>[0-9\-]+)\.\.(?P<x2>[0-9\-]+), y=(?P<y1>[0-9\-]+)\.\.(?P<y2>[0-9\-]+)`)
	match := compRegEx.FindStringSubmatch(line)

	paramsMap := make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

func simulate(xVelocity, yVelocity int, target *Target) (int, bool) {
	maxY := 0
	currX := 0
	currY := 0
	i := 0

	for true {
		currX += xVelocity
		currY += yVelocity
		yVelocity -= 1

		if xVelocity > 0 {
			xVelocity--
		} else if xVelocity < 0 {
			xVelocity++
		}

		if currY > maxY {
			maxY = currY
		}
		i++

		if target.hit(currX, currY) {
			return maxY, true
		}

		if currY < target.y1 && yVelocity < 0 {
			break
		}
	}
	return maxY, false
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func search(target *Target) (int, int) {
	xStart := 1
	xEnd := target.x2
	yStart := target.y1
	yEnd := abs(target.y1)
	yMax := math.MinInt
	hits := 0

	for x := xStart; x <= xEnd; x++ {
		for y := yStart; y < yEnd; y++ {
			high, hit := simulate(x, y, target)

			if hit {
				hits++
				//fmt.Printf("Hit xVelocity: %d  yVelocity: %d  high: %d\n", x, y, high)
			}

			if hit && high > yMax {
				yMax = high
			}
		}
	}

	return yMax, hits
}

func main() {
	lines := readFile("17.txt")
	target := parseTarget(lines[0])
	high, hits := search(target)
	fmt.Printf("y max: %d   possible trajectories: %d\n", high, hits)
}
