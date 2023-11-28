package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	// Source: https://scrabutility.com/TWL06.txt
	//go:embed all.txt
	largeDictionary string

	// Source: https://www.ef.edu/english-resources/english-vocabulary/top-3000-words/
	//go:embed common.txt
	smallDictionary string
)

func parseDictionary(input string) map[string]bool {
	input = strings.ToLower(input)
	words := strings.Fields(input)
	dict := make(map[string]bool, len(words))
	for _, word := range words {
		dict[word] = true
	}
	return dict
}

const boardSize = 4

func main() {
	fmt.Println("Input your Boggle board, separating tiles with spaces:")
	reader := bufio.NewReader(os.Stdin)
	var board [boardSize][boardSize]string // can't use rune as there's a Qu board character
	for i := 0; i < 4; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("reading input: %s", err)
		}
		line = strings.TrimSpace(line)
		line = strings.ToLower(line)

		tiles := strings.Split(line, " ")
		if got, want := len(tiles), 4; got != want {
			log.Fatalf("read %d tiles, want %d space-separated tiles", got, want)
		}
		for j, tile := range tiles {
			board[i][j] = tile
		}
	}

	commonWords := parseDictionary(smallDictionary)
	dictionary := parseDictionary(largeDictionary)

	type Point struct {
		X, Y int
	}
	var queue [][]Point

	// resolve converts the points in path to letters on the board.
	resolve := func(path []Point) string {
		var s strings.Builder
		for _, point := range path {
			s.WriteString(board[point.X][point.Y])
		}
		return s.String()
	}
	// last returns the (x, y) coordinates of the last Point in path.
	last := func(path []Point) (int, int) {
		last := path[len(path)-1]
		return last.X, last.Y
	}
	// pop removes and returns the first Point from the queue.
	pop := func() []Point {
		defer func() {
			queue = queue[1:]
		}()
		return queue[0]
	}
	// push adds another path to visit to the back of the queue.
	// For convenience, it also handles appending the next Point to the current path.
	push := func(current []Point, next Point) {
		var path []Point // avoid clobbering backing arrays
		path = append(path, current...)
		path = append(path, next)
		queue = append(queue, path)
	}
	// ok returns true if p is on the board and hasn't already been visited in path.
	ok := func(path []Point, p Point) bool {
		for _, visited := range path {
			if visited.X == p.X && visited.Y == p.Y {
				return false
			}
		}
		return p.X >= 0 && p.X < boardSize && p.Y >= 0 && p.Y < boardSize
	}

	// Initial seed - traverse each character on the board
	for x := 0; x < boardSize; x++ {
		for y := 0; y < boardSize; y++ {
			push(nil, Point{x, y})
		}
	}

	type Match struct {
		Path []Point
		Word string
	}
	var allMatches []Match
	commonMatches := make(map[string]bool)
	uncommonMatches := make(map[string]bool)

	var traversed int // Keep track of total number of options traversed to report back later for fun.
	for len(queue) > 0 {
		traversed++
		path := pop()
		x, y := last(path)

		// Words of at least 3 letters are considered matches.
		word := resolve(path)
		if len(word) >= 3 && dictionary[word] {
			allMatches = append(allMatches, Match{
				Path: path,
				Word: word,
			})
			if commonWords[word] {
				commonMatches[word] = true
			} else {
				uncommonMatches[word] = true
			}
		}

		// Traverse touching tiles.
		for _, candidate := range []Point{
			{x, y - 1},     // up
			{x, y + 1},     // down
			{x - 1, y},     // left
			{x + 1, y},     // right
			{x - 1, y - 1}, // diag top left
			{x + 1, y - 1}, // diag top right
			{x - 1, y + 1}, // diag lower left
			{x + 1, y + 1}, // diag lower right
		} {
			if ok(path, candidate) {
				push(path, candidate)
			}
		}
	}

	for _, match := range allMatches {
		fmt.Printf("Found %q via %v\n", match.Word, match.Path)
	}

	fmt.Printf("\n===\n")
	fmt.Printf("Checked %d different paths and found %d common words:\n", traversed, len(commonMatches))
	for _, word := range sortedKeys(commonMatches) {
		fmt.Printf("  %v\n", word)
	}
	fmt.Printf("and %d uncommon words: %v\n", len(uncommonMatches), sortedKeys(uncommonMatches))
}

func sortedKeys(in map[string]bool) (out []string) {
	for k := range in {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
