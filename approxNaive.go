package fuzzyfind

// This file contains functions meant for learning
// various string matching algorithms

// A naive find algorithm
func naiveFind(subsequence string, sequence string) []Match {
	// Check for empty strings first
	if subsequence == "" || sequence == "" {
		return []Match{}
	}
	subseq := []rune(subsequence)
	seq := []rune(sequence)
	matches := []Match{}

	for i := 0; i <= len(seq)-len(subseq); i++ {
		match := true
		for j := 0; j <= len(subseq)-1; j++ {
			if seq[i+j] != subseq[j] {
				match = false
				break
			}
		}
		if match {
			matches = append(matches, Match{i, i + len(subseq), 0})
		}
	}
	return matches
}

// Naive Fuzzy find allowing mismatches
func naiveFuzzyFind(subsequence string, sequence string, maxMM int) []Match {
	// Check for empty strings first
	if subsequence == "" || sequence == "" {
		return []Match{}
	}
	subseq := []rune(subsequence)
	seq := []rune(sequence)
	matches := []Match{}
	for i := 0; i <= len(seq)-len(subseq); i++ {
		match := true
		mm := 0 // Initialize the mismatch counter
		for j := 0; j <= len(subseq)-1; j++ {
			if seq[i+j] != subseq[j] {
				mm++
				if mm > maxMM {
					match = false
					break
				}
			}
		}
		if match {
			matches = append(matches, Match{i, i + len(subseq), mm})
		}
	}
	return matches
}
