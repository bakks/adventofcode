package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
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

type Cave struct {
	label string
	edges []*Cave
}

func NewCave(label string) *Cave {
	return &Cave{
		label: label,
		edges: []*Cave{},
	}
}

func (this *Cave) AddEdge(link *Cave) {
	this.edges = append(this.edges, link)
}

type CaveSystem struct {
	caves map[string]*Cave
}

func NewCaveSystem() *CaveSystem {
	return &CaveSystem{caves: map[string]*Cave{}}
}

func (this *CaveSystem) AddEdge(label1, label2 string) {
	labels := []string{label1, label2}
	for _, label := range labels {
		_, ok := this.caves[label]

		if !ok {
			// this cave doesn't exist
			this.caves[label] = NewCave(label)
		}
	}

	this.caves[label1].AddEdge(this.caves[label2])
	this.caves[label2].AddEdge(this.caves[label1])
}

func parseCaveSystem(lines []string) *CaveSystem {
	caves := NewCaveSystem()

	for _, line := range lines {
		fields := strings.Split(line, "-")
		fmt.Println(fields)
		caves.AddEdge(fields[0], fields[1])
	}

	return caves
}

type CavePath struct {
	path []string
}

func NewCavePath(path []string) *CavePath {
	return &CavePath{path}
}

func (this *CavePath) Print() {
	for _, label := range this.path {
		fmt.Printf("%s,", label)
	}
	fmt.Printf("\n")
}

func contains(stack []string, target string) int {
	found := 0

	for _, s := range stack {
		if s == target {
			found++
		}
	}

	return found
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func canVisitTarget(path []string, target string) bool {
	if target == "start" {
		return false
	}
	if target == "end" {
		return true
	}
	if IsUpper(target) {
		return true
	}

	// lowercase target, only ok if we haven't visited a node twice
	found := map[string]int{}
	alreadyHasDoubleVisit := false

	for _, node := range path {
		if IsUpper(node) {
			continue
		}

		_, ok := found[node]
		if !ok {
			found[node] = 1
		} else {
			found[node]++
			alreadyHasDoubleVisit = true
		}
	}

	if !alreadyHasDoubleVisit {
		return true
	}

	_, ok := found[target]
	return !ok
}

func (this *CaveSystem) search(currNode *Cave, currPath []string) [][]string {
	path := make([]string, len(currPath))
	copy(path, currPath)
	if currNode.label == "end" {
		return [][]string{path}
	}

	paths := [][]string{}

	for _, target := range currNode.edges {
		var result [][]string = nil

		if canVisitTarget(currPath, target.label) {
			result = this.search(target, append(path, target.label))
		}

		if result != nil {
			paths = append(paths, result...)
		}
	}

	return paths
}

func (this *CaveSystem) FindPath() []*CavePath {
	start := this.caves["start"]
	paths := this.search(start, []string{"start"})
	cavePaths := []*CavePath{}

	for _, path := range paths {
		cavePaths = append(cavePaths, NewCavePath(path))
	}

	return cavePaths
}

func main() {
	lines := readFile("12.txt")
	caveSystem := parseCaveSystem(lines)
	paths := caveSystem.FindPath()

	for _, path := range paths {
		path.Print()
	}
	fmt.Printf("%d paths\n", len(paths))
}
