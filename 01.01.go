package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("01.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	s := bufio.NewScanner(f)

	last := -1
	n := -1

	for s.Scan() {
		current, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Fatal(err)
		}

		if current > last {
			n = n + 1
		}
		last = current
	}

	fmt.Printf("%d\n", n)

	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}
}
