package approx

// A match
type Match struct {
	Start int
	End   int
	Dist  int
}

type MatchFunction func(rune, rune) bool

type Options struct {
	InsCost int
	DelCost int
	SubCost int
	Matches MatchFunction
}

// DefaultOptions is the default options: insertion cost is 1, deletion cost is
// 1, substitution cost is 1, and two runes match iff they are the same.
var DefaultOptions Options = Options{
	InsCost: 1,
	DelCost: 1,
	SubCost: 1,
	Matches: func(sourceCharacter rune, targetCharacter rune) bool {
		return sourceCharacter == targetCharacter
	},
}
