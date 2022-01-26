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

func contains(stack []string, target string) bool {
	for _, s := range stack {
		if s == target {
			return true
		}
	}

	return false
}

func canVisitTwice(label string) bool {
	return IsUpper(label)
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
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

		if canVisitTwice(target.label) || !contains(path, target.label) {
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
