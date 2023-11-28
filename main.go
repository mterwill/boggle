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

func main() {
	const boardSize = 4

	fmt.Println("Input your Boggle board, separating tiles with spaces:")
	reader := bufio.NewReader(os.Stdin)
	var board [][]string // can't use rune as there's a Qu board character
	for i := 0; i < boardSize; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("reading input: %s", err)
		}
		line = strings.TrimSpace(line)
		line = strings.ToLower(line)

		tiles := strings.Split(line, " ")
		if got, want := len(tiles), boardSize; got != want {
			log.Fatalf("read %d tiles, want %d space-separated tiles", got, want)
		}
		board = append(board, tiles)
	}

	solution := solve(board)
	common := make(map[string]bool)
	uncommon := make(map[string]bool)
	for _, match := range solution.Matches {
		fmt.Printf("Found %q via %v\n", match.Word, match.Path)
		if match.Common {
			common[match.Word] = true
		} else {
			uncommon[match.Word] = true
		}
	}

	fmt.Printf("\n===\n")
	fmt.Printf("Checked %d different paths and found %d common words:\n", solution.Traversed, len(common))
	for _, word := range sortedKeys(common) {
		fmt.Printf("  %v\n", word)
	}
	fmt.Printf("and %d uncommon words: %v\n", len(uncommon), sortedKeys(uncommon))
}

type Point struct {
	X, Y int
}

type Match struct {
	Path   []Point
	Word   string
	Common bool
}

type Solution struct {
	Matches   []Match
	Traversed int
}

func solve(board [][]string) Solution {
	boardSize := len(board)
	commonWords := parseDictionary(smallDictionary)
	dictionary := parseDictionary(largeDictionary)

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
	// pop removes and returns the first path from the queue.
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

	var solution Solution
	for len(queue) > 0 {
		solution.Traversed++
		path := pop()
		x, y := last(path)

		// Words of at least 3 letters are considered matches.
		word := resolve(path)
		if len(word) >= 3 && dictionary[word] {
			solution.Matches = append(solution.Matches, Match{
				Path:   path,
				Word:   word,
				Common: commonWords[word],
			})
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

	return solution
}

func sortedKeys(in map[string]bool) (out []string) {
	for k := range in {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
