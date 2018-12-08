package fuzzyfind

import (
	"errors"
	"fmt"
	"os"
)

// FuzzyFindShort uses a modified Levenshtein function to find a 'pattern' in a 'text'.
// It is not the most optimal way to do this for longer strings, thus it is recommend
// for use on short patterns only (to be defined lanter). If you are finding that it
// returns many options for a single pattern, ie find 'perl' in 'berd' with a max dist of 2,
// You should concider ammending the Options to make the version you don't want to see cost more.
func FuzzyFindShort(pattern string, text string, maxE int, op Options) ([]Match, error) {
	// Check for empty strings first
	if pattern == "" || text == "" {
		return []Match{}, errors.New("Pattern and or Text is empty")
	}
	// Runify
	p := []rune(pattern)
	t := []rune(text)
	matches := approxLeven(p, t, maxE, op)
	return matches, nil
}

// This version makes use of the pigeon hole principle, which is the idea that
// if I am going to have x number of mutations, then if I split my pattern into
// x + 1 regions, I will have at least one region that will match exaclty to the
// text. Then I can extend the match from there.
// If the chunk size of the pattern would be less than 1, or the pattern is longer than the
// text - maxE, FuzzyFindShort is run instead
func FuzzyFindPigeon(pattern string, text string, maxE int, op Options) ([]Match, error) {
	// Check for empty strings first
	if pattern == "" || text == "" {
		return []Match{}, errors.New("Pattern and or Text is empty")
	}

	if len(pattern)/(maxE+1) < 1 {
		fmt.Fprintf(os.Stderr, "Pattern is too short for this method with given max edit distance\n")
		return FuzzyFindShort(pattern, text, maxE, op)
		//return []Match{}, errors.New("Pattern is too short for this method with given max edit distance")
	} else if len(pattern) >= len(text)-maxE {
		fmt.Fprintf(os.Stderr, "Pattern is too long for given text, running FuzzyFindShort Instead\n")
		return FuzzyFindShort(pattern, text, maxE, op)
		//return []Match{}, errors.New("Pattern is too long for the given text")
	}

	// Runify
	p := []rune(pattern)
	t := []rune(text)
	matches := approxPigeon(p, t, maxE, op)

	return matches, nil
}
