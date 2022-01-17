package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func parseDiagnostic(in []string) [][]bool {
	output := [][]bool{}

	for _, line := range in {
		outLine := []bool{}

		for _, x := range line {
			switch x {
			case '0':
				outLine = append(outLine, false)
			case '1':
				outLine = append(outLine, true)
			default:
				log.Panic(x)
			}
		}

		output = append(output, outLine)
	}

	return output
}

const DiagnosticLineLength = 12

func main() {
	lines := readFile("03.txt")
	diagnostic := parseDiagnostic(lines)
	epsAcc := [DiagnosticLineLength]int{}

	for _, diagLine := range diagnostic {
		for i, diagBit := range diagLine {
			if diagBit {
				epsAcc[i]++
			}
		}
	}

	diagnosticLines := len(diagnostic)
	epsilon := uint(0)

	for _, val := range epsAcc {
		epsilon = epsilon << 1

		if float64(val)/float64(diagnosticLines) > 0.5 {
			epsilon = epsilon | 1
			fmt.Printf("1")
		} else {
			fmt.Printf("0")
		}
	}

	gamma := epsilon ^ 0xFFF
	fmt.Printf("\n%d\n%d\n", epsilon, gamma)
	fmt.Print(epsilon * gamma)

}
