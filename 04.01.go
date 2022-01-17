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

const BingoBoardSize = 5

type BingoBoard struct {
	board [][]int
	marks [][]bool
}

func (this BingoBoard) mark(mark int) {
	for i := 0; i < len(this.board); i++ {
		for j := 0; j < len(this.board[i]); j++ {
			if this.board[i][j] == mark {
				this.marks[i][j] = true
				return
			}
		}
	}
}

func (this BingoBoard) sumUnmarked() int {
	x := 0
	for i := 0; i < len(this.board); i++ {
		for j := 0; j < len(this.board); j++ {
			if !this.marks[i][j] {
				x += this.board[i][j]
			}
		}
	}

	return x
}

func (this BingoBoard) hasWon() bool {
	// check horizontals
	for i := 0; i < len(this.board); i++ {
		j := 0
		for ; j < len(this.board); j++ {
			if !this.marks[i][j] {
				break
			}
		}
		if j == len(this.board) {
			return true
		}
	}

	// check verticals
	for i := 0; i < len(this.board); i++ {
		j := 0
		for ; j < len(this.board); j++ {
			if !this.marks[j][i] {
				break
			}
		}
		if j == len(this.board) {
			return true
		}
	}

	return false
}

func BingoBoardFromFileInput(lines []string) BingoBoard {
	if len(lines) != BingoBoardSize {
		log.Fatalf("Expected board if size %d, received %d lines\n", BingoBoardSize, len(lines))
	}

	board := [][]int{}

	for _, line := range lines {
		numbers := strings.Fields(line)
		if len(numbers) != BingoBoardSize {
			log.Fatalf("Line of size %d : %s\n", len(numbers), line)
		}

		boardLine := []int{}
		for i := 0; i < BingoBoardSize; i++ {
			n, err := strconv.Atoi(numbers[i])
			if err != nil {
				log.Fatal(err)
			}
			boardLine = append(boardLine, n)
		}
		board = append(board, boardLine)
	}

	marks := make([][]bool, BingoBoardSize)
	for i := range marks {
		marks[i] = make([]bool, BingoBoardSize)
	}

	return BingoBoard{board, marks}
}

func parseBingo(in []string) ([]int, []BingoBoard) {
	moves := []int{}

	// parse the top line containing the moves
	for _, stringNum := range strings.Split(in[0], ",") {
		n, err := strconv.Atoi(stringNum)
		if err != nil {
			log.Fatal(err)
		}
		moves = append(moves, n)
	}

	boards := []BingoBoard{}

	// iterate through specific boards
	for i := 1; i < len(in); i += BingoBoardSize + 1 {
		boardLines := in[i+1 : i+1+BingoBoardSize]
		board := BingoBoardFromFileInput(boardLines)
		boards = append(boards, board)
	}

	return moves, boards
}

func main() {
	lines := readFile("04.txt")
	moves, boards := parseBingo(lines)
	log.Printf("Found %d boards\n", len(boards))

	for i, move := range moves {
		log.Printf("Move: %d\n", move)
		for _, board := range boards {
			board.mark(move)

			if i >= BingoBoardSize && board.hasWon() {
				fmt.Print(board)
				sum := board.sumUnmarked()
				score := sum * move
				fmt.Printf("Sum: %d, Score: %d\n", sum, score)
				return
			}
		}
	}

	log.Printf("No winner")
}
