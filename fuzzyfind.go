package fuzzyfind

import (
	"errors"
)

type Candidate struct {
	Start       int
	SubseqIndex int
	Dist        int
}

type Match struct {
	Start int
	End   int
	Dist  int
}

// Check if a rune slice contains a rune
func find(b []rune, c rune) int {
	for idx, char := range b {
		if char == c {
			return idx
		}
	}
	return -1
}

func prepareInitCandidatesMap(subseq []rune, maxLDist int) map[rune]int {
	candidateMap := map[rune]int{}
	for index, char := range subseq[:maxLDist+1] {
		if find(subseq[:index], char) == -1 {
			candidateMap[char] = index
		}
	}
	return candidateMap
}

func FindNearMatches(subsequence string, sequence string, maxLDist int) ([]Match, error) {
	// Ensure substring has at least some values
	if subsequence == "" {
		return []Match{}, errors.New("Subsequence must have values")
	}

	// Normalize the strings
	subseq := []rune(subsequence)
	seq := []rune(sequence)
	matches := []Match{}

	// Prepare some often used things in advance
	initCandidateMap := prepareInitCandidatesMap(subseq, maxLDist)
	subseqLen := len(subseq)

	candidates := []Candidate{}
	//lastCandidateInitIndex := -1

	for index, char := range seq {
		newCandidates := []Candidate{}
		idxInSubseq := initCandidateMap[char]
		if _, ok := initCandidateMap[char]; ok {
			newCandidates = append(newCandidates, Candidate{index, idxInSubseq + 1, idxInSubseq})
		}
		for _, cand := range candidates {
			nextSubseqChars := subseq[cand.SubseqIndex : cand.SubseqIndex+maxLDist-cand.Dist+1]
			idx := find(nextSubseqChars, char)

			// If this sequence char is the candidate's next expected char
			if idx == 0 {
				// If reached the end of the subsequence, return a match
				if cand.SubseqIndex+1 == subseqLen {
					matches = append(matches, Match{cand.Start, index + 1, cand.Dist})

				} else { // otherwise update the candidate and keep it
					newCand := cand
					newCand.SubseqIndex++
					newCandidates = append(newCandidates, newCand)
				}
			} else { // If this sequence char is *not* the candidate's next expected char
				// We can try skipping a sequence or sub-seq char (or both),
				// unless this candidate has already skipped the max allowed
				// number of chars
				if cand.Dist == maxLDist {
					continue
				}

				// Add a candidate skipping a sequence char
				newCandSeqSkip := cand
				newCandSeqSkip.Dist++
				newCandidates = append(newCandidates, newCandSeqSkip)

				if index+1 < len(seq) && cand.SubseqIndex+1 < subseqLen {
					// add a candidate skipping both a sequence char and a
					// subsequence char
					newCandBothSkip := cand
					newCandBothSkip.Dist++
					newCandBothSkip.SubseqIndex++
					newCandidates = append(newCandidates, newCandBothSkip)
				}

				// Try skipping subseq chars
				for nSkipped := 1; nSkipped < maxLDist-cand.Dist+1; nSkipped++ {
					if cand.SubseqIndex+nSkipped == subseqLen {
						// If skipping nSkipped subseq chars reaches the ned
						// of the subseq, add a match
						matches = append(matches, Match{cand.Start, index + 1, cand.Dist + nSkipped})
						break
					} else if subseq[cand.SubseqIndex+nSkipped] == char {
						// otherwise, if skipping nSkipped subseq chars
						// reaches a subseq char identical to this subseq,
						// add a candidate skipping nSkipped subseq chars
						newCandSubSeqSkip := cand
						newCandSubSeqSkip.Dist += nSkipped
						newCandSubSeqSkip.SubseqIndex += nSkipped + 1
						newCandidates = append(newCandidates, newCandSubSeqSkip)
						break
					}
				}
				// note: if the above loop ends without a break, that means that
				// no candidate could be added by skipping subseq chars

			}
		}
		candidates = newCandidates

	}
	for _, cand := range candidates {
		dist := cand.Dist + len(subseq) - cand.SubseqIndex
		if dist <= maxLDist {
			matches = append(matches, Match{cand.Start, len(seq), dist})
		}
	}
	return matches, nil

}
