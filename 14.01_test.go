package main

import "testing"
import "fmt"

func BenchmarkStuff(b *testing.B) {
	lines := ReadFile("14.txt")
	polymer, rules := parsePolymer(lines)

	for i := 0; i < b.N; i++ {
		pptr := &polymer
		for i := 0; i < 18; i++ {
			pptr = polymerize(pptr, rules)
			fmt.Printf("Iteration: %d, Length: %d\n", i, len(*pptr))
		}
	}
}
