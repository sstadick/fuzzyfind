package approx

import (
	"fmt"
	"index/suffixarray"
	"os"
	"sort"
)

// Return a slice of rune slices containing non-overlapping,
// non-empty substrings that cover p.  They should be
// as close to equal-length as possible.
func partition(p []rune, pieces int) [][]rune {
	base, mod := len(p)/pieces, len(p)%pieces
	idx := 0
	ps := [][]rune{}
	modAdjust := 1
	for i := 0; i < pieces; i++ {
		if i >= mod {
			modAdjust = 0
		}
		newIdx := idx + base + modAdjust
		ps = append(ps, p[idx:newIdx])
		idx = newIdx
	}
	return ps
}

// Break the pattern up into chunks, search for the chunks with exact match that uses a suffix index
// Extend chunks when a match occurs, this is really only worth doing when the patterns and strings
// get pretty long. For very repetative sequences, this can end up doing more work than a regular leven
func approxPigeon(pattern []rune, text []rune, maxE int, op Options) ([]Match, error) {
	partitions := partition(pattern, maxE+1)
	//fmt.Printf("Parts: %q\n", partitions)
	offset := 0
	occurances := make(map[Match]int)
	// Should i be concerned about this conversion messing with my index's?
	// I would still like to make my own kmer index possibly
	textIndex := suffixarray.New([]byte(string(text)))

	for _, part := range partitions {
		hits := textIndex.Lookup([]byte(string(part)), -1)

		for _, hit := range hits {
			//fmt.Printf("Hit %s at %d || offset %d\n", string(part), hit, offset)
			if hit-offset < 0 {
				//fmt.Printf("Fell off left side\n")
				continue // pattern falls off left end of T?
			}
			if hit+len(pattern)-offset > len(text) {
				//fmt.Printf("Fell off right side\n")
				continue // falls off right end?
			}
			rightIdx := hit + len(pattern) - offset
			leftIdx := hit - offset
			if rightIdx+maxE <= len(text)-1 {
				rightIdx += maxE
			}
			if leftIdx-maxE >= 0 {
				leftIdx -= maxE
			}

			possible, err := approxLeven(pattern, text[leftIdx:rightIdx], maxE, op)
			if err != nil {
				fmt.Fprintf(os.Stderr, "bad thing in approxPigeon, find %q in %q", pattern, text[leftIdx:rightIdx])
			}
			for _, p := range possible {
				occurances[Match{
					Start: p.Start + leftIdx,
					End:   p.End + leftIdx,
					Dist:  p.Dist,
				}]++
			}

			// TODO: Fix this failed attempt to extend one side at a time
			// Extend right only if we aren't already at far right
			//			rightMost := hit + len(part) // set the default right most
			//			totalE := 0
			//			if offset+len(part) != len(pattern) {
			//				rightStart := hit + len(part)
			//				rightEnd := hit + len(pattern) - offset
			//				fmt.Printf("Part %s, Offset: %d Rightstart: %d Rightend: %d\n", string(part), offset, rightStart, rightEnd)
			//				rightMatches := approxLeven(text[rightStart:rightEnd], pattern[offset:], maxE, op)
			//				fmt.Printf("Searching %s (%d:%d)\n", string(text[rightStart:rightEnd]), rightStart, rightEnd)
			//				fmt.Printf("Pattern %s (%d:%d)\n", string(pattern[offset:]), offset, len(pattern))
			//				if len(rightMatches) == 0 {
			//					continue
			//				}
			//				lowestRightMatch := Match{Dist: MaxInt}
			//				for _, m := range rightMatches {
			//					if m.Dist < lowestRightMatch.Dist {
			//						lowestRightMatch = m
			//					}
			//				}
			//				rightMost = lowestRightMatch.End + rightStart
			//				totalE += lowestRightMatch.Dist
			//			}
			//			// Extend left only if we aren't already on the far left
			//			leftMost := hit + offset // Set the default leftmost
			//			if offset != 0 {
			//				leftEnd := hit
			//				leftStart := hit - offset
			//				// If dist allowed is now zero, I should have a faster fallback
			//				leftMatches := approxLeven(text[leftStart:leftEnd], pattern[:offset], maxE-totalE, op)
			//				fmt.Printf("Searching %s (%d:%d)\n", string(text[leftStart:leftEnd]), leftStart, leftEnd)
			//				fmt.Printf("Pattern %s (%d:%d)\n", string(pattern[:offset]), 0, offset)
			//				fmt.Printf("Left Foudn: %v\n", leftMatches)
			//				if len(leftMatches) == 0 {
			//					continue
			//				}
			//				lowestLeftMatch := Match{Dist: MaxInt}
			//				for _, m := range leftMatches {
			//					if m.Dist < lowestLeftMatch.Dist {
			//						lowestLeftMatch = m
			//					}
			//				}
			//				leftMost = lowestLeftMatch.Start + leftStart
			//				totalE += lowestLeftMatch.Dist
			//			}
			//			// If we've gotten this far we have a match
			//			occurances[Match{
			//				Start: leftMost,
			//				End:   rightMost,
			//				Dist:  totalE,
			//			}]++

		}
		offset += len(part)

	}
	matches := make([]Match, 0, len(occurances))
	for k := range occurances {
		matches = append(matches, k)
	}
	sort.Slice(matches, func(a, b int) bool {
		if matches[a].Start == matches[b].Start {
			return matches[a].End < matches[b].End
		}
		return matches[a].Start < matches[b].Start
	})
	return matches, nil
}
