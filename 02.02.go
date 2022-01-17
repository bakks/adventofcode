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
	f, err := os.Open(filename)
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

func main() {
	lines := readFile("02.txt")

	horizontal := 0
	depth := 0
	aim := 0
	n := 0

	for _, line := range lines {
		fields := strings.Split(line, " ")
		command := fields[0]
		units, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Fatal(err)
		}

		switch command {
		case "up":
			aim -= units
		case "down":
			aim += units
		case "forward":
			horizontal += units
			depth += aim * units
		}
		n++
	}

	fmt.Printf("%d actions, %d depth, %d forward, %d total", n, depth, horizontal, depth*horizontal)
}
