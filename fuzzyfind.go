package fuzzyfind

// FuzzyFindShort uses a modified Levenshtein function to find a 'pattern' in a 'text'.
// It is not the most optimal way to do this for longer strings, thus it is recommend
// for use on short patterns only (to be defined lanter). If you are finding that it
// returns many options for a single pattern, ie find 'perl' in 'berd' with a max dist of 2,
// You should concider ammending the Options to make the version you don't want to see cost more.
func FuzzyFindShort(pattern string, text string, maxE int, op Options) []Match {
	// Check for empty strings first
	if pattern == "" || text == "" {
		return []Match{}
	}
	// Runify
	p := []rune(pattern)
	t := []rune(text)
	matches := approxLeven(p, t, maxE, op)
	return matches
}
