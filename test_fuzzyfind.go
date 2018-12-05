package fuzzyfind

import (
	"testing"
)

func TestprepareInitCandidatesMap(t *testing.T) {
	result := prepareInitCandidatesMap([]rune("GATTACA"), 2)
	if result['G'] != 0 && result['A'] != 1 && result['T'] != 2 {
		t.Error("Error making candidate dict")
	}
}
