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
	set := [3]int{}
	n := 0
	i := 0

	for s.Scan() {
		currLine, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Fatal(err)
		}

		set[i%3] = currLine
		currSum := set[0] + set[1] + set[2]

		fmt.Printf("%d %d : ", currLine, currSum)

		if i >= 3 && currSum > last {
			n++
			fmt.Printf("greater")
		}
		fmt.Printf("\n")
		last = currSum
		i++
	}

	fmt.Printf("%d\n", n)

	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}
}
