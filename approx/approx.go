// approx is a package for approximate string matching. It will extract strings withing
// a given edit distance and return their start, end, and distance.
package approx

import (
	"fmt"
	"os"
)

// ApproxFind uses a modified Levenshtein function to find a 'pattern' in a 'text'.
// It is not the most optimal way to do this for longer strings, thus it is recommend
// for use on short patterns only (to be defined lanter). If you are finding that it
// returns many options for a single pattern, ie find 'perl' in 'berd' with a max dist of 2,
// You should concider ammending the Options to make the version you don't want to see cost more.
func ApproxFind(pattern string, text string, maxE int, op Options) ([]Match, error) {
	// Check for empty strings first
	if pattern == "" {
		return nil, fmt.Errorf("pattern to search empty")
	} else if text == "" {
		return nil, fmt.Errorf("text to search is empty")
	}
	// Runify
	p := []rune(pattern)
	t := []rune(text)
	matches, err := approxLeven(p, t, maxE, op)
	if err != nil {
		return matches, err
	}
	return matches, nil
}

// This version makes use of the pigeon hole principle, which is the idea that
// if I am going to have x number of mutations, then if I split my pattern into
// x + 1 regions, I will have at least one region that will match exaclty to the
// text. Then I can extend the match from there.
// If the chunk size of the pattern would be less than 1, or the pattern is longer than the
// text - maxE, ApproxFind is run instead
func approxFindPigeon(pattern string, text string, maxE int, op Options) ([]Match, error) {
	// Check for empty strings first
	if pattern == "" {
		return nil, fmt.Errorf("pattern to search empty")
	} else if text == "" {
		return nil, fmt.Errorf("text to search is empty")
	}

	if len(pattern)/(maxE+1) < 1 {
		fmt.Fprintf(os.Stderr, "Pattern is too short for this method with given max edit distance\n")
		return ApproxFind(pattern, text, maxE, op)
		//return []Match{}, errors.New("Pattern is too short for this method with given max edit distance")
	} else if len(pattern) >= len(text)-maxE {
		fmt.Fprintf(os.Stderr, "Pattern is too long for given text, running ApproxFind Instead\n")
		return ApproxFind(pattern, text, maxE, op)
		//return []Match{}, errors.New("Pattern is too long for the given text")
	}

	// Runify
	p := []rune(pattern)
	t := []rune(text)
	matches, err := approxPigeon(p, t, maxE, op)
	if err != nil {
		return matches, err
	}

	return matches, nil
}

// From a list of matches, grab the longest leftmost match with the best score
// TODO: Return an error when we are given no matches
func BestMatch(matches []Match) (Match, error) {

	if len(matches) == 0 {
		return Match{}, fmt.Errorf("No matches in list")
	}

	// Start with a match that is horrible
	topMatch := Match{
		Start: MaxInt,
		End:   MaxInt,
		Dist:  MaxInt,
	}
	for _, match := range matches {

		if match.Dist < topMatch.Dist {
			// match is better
			topMatch = match
		} else if match.Dist == topMatch.Dist {
			// Match dist equal, see which is leftmost
			if match.Start < topMatch.Start {
				topMatch = match
			} else if match.Start == topMatch.Start {

				// match starts are equal, check which is longer
				if (match.End - match.Start) > (topMatch.End - topMatch.Start) {
					// match is longer
					topMatch = match
				} else if (match.End - match.Start) == (topMatch.End - topMatch.Start) {
					// This case should never be hit, leave topMatch on top
					fmt.Fprintf(os.Stderr, "This state should never be hit")
				}
			}
		}
	}
	return topMatch, nil
}
