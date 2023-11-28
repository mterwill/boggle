package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSolve(t *testing.T) {
	solution := solve([][]string{
		{"h", "r", "s", "p"},
		{"e", "f", "u", "n"},
		{"i", "o", "r", "e"},
		{"i", "r", "o", "y"},
	})

	if got, want := solution.Traversed, 12_029_640; got != want {
		t.Errorf("traversed %d paths, expected %d", got, want)
	}

	// Sanity check the common words.
	gotCommon := make(map[string]bool)
	for _, match := range solution.Matches {
		if match.Common {
			gotCommon[match.Word] = true
		}
	}
	wantCommon := []string{
		"ensure",
		"for",
		"four",
		"fun",
		"her",
		"our",
		"pure",
		"roof",
		"run",
		"sue",
		"sun",
		"sure",
	}
	if diff := cmp.Diff(sortedKeys(gotCommon), wantCommon); diff != "" {
		t.Errorf("common words differ (-got +want):\n%s", diff)
	}
}
