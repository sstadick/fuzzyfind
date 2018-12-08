package fuzzyfind

import (
	"testing"
)

type TestCase struct {
	Pattern     string
	Text        string
	Description string
	MaxDist     int
	Expected    []Match
}

var ExactTestCases = []TestCase{
	TestCase{
		Pattern:     "and",
		Text:        "The dog and the cat.",
		Description: "Sanity Check",
		Expected: []Match{
			Match{
				Start: 8,
				End:   11,
				Dist:  0,
			},
		},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "GATTACATATATACA",
		Description: "Exact match at start",
		Expected: []Match{
			Match{
				Start: 0,
				End:   7,
				Dist:  0,
			},
		},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "GATTACA",
		Description: "Identical Sequences",
		Expected: []Match{
			Match{
				Start: 0,
				End:   7,
				Dist:  0,
			},
		},
	},
	TestCase{
		Pattern:     "",
		Text:        "GATTACATATATACA",
		Description: "Empty pattern",
		Expected:    []Match{},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "",
		Description: "Empty Text",
		Expected:    []Match{},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "ATCGACTGAGATTACA",
		Description: "Find sequence at end",
		Expected: []Match{
			Match{
				Start: 9,
				End:   16,
				Dist:  0,
			},
		},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "GATTACAATCGACTGA",
		Description: "Find sequence at start",
		Expected: []Match{
			Match{
				Start: 0,
				End:   7,
				Dist:  0,
			},
		},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "GATTACAATCGGATTACAACTGA",
		Description: "Find sequence twice",
		Expected: []Match{
			Match{
				Start: 0,
				End:   7,
				Dist:  0,
			},
			Match{
				Start: 11,
				End:   18,
				Dist:  0,
			},
		},
	},
}

var MismatchTestCases = []TestCase{
	TestCase{
		Pattern:     "and",
		Text:        "The dog amd the cat.",
		Description: "Single mm in middle",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 8,
				End:   11,
				Dist:  1,
			},
		},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "CATTACATATATACA",
		Description: "Single mm at start",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 0,
				End:   7,
				Dist:  1,
			},
		},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "GATTACA",
		Description: "Identical Sequences, mm allowed",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 0,
				End:   7,
				Dist:  0,
			},
		},
	},
	TestCase{
		Pattern:     "ACATCC",
		Text:        "GATTACATATATACATCG",
		Description: "Single mm at end",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 12,
				End:   18,
				Dist:  1,
			},
		},
	},
	TestCase{
		Pattern:     "ACATCC",
		Text:        "GATTACATATATGCATCT",
		Description: "Double mm",
		MaxDist:     2,
		Expected: []Match{
			Match{
				Start: 4,
				End:   10,
				Dist:  2,
			},
			Match{
				Start: 8,
				End:   14,
				Dist:  2,
			},
			Match{
				Start: 12,
				End:   18,
				Dist:  2,
			},
		},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "GACTACAATCGGATTTCAACTGA",
		Description: "Find sequence twice, with single mm in each",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 0,
				End:   7,
				Dist:  1,
			},
			Match{
				Start: 11,
				End:   18,
				Dist:  1,
			},
		},
	},
}
var EditTestCases = []TestCase{
	TestCase{
		Pattern:     "and",
		Text:        "The dog amd the cat.",
		Description: "Single mm in middle",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 8,
				End:   11,
				Dist:  1,
			},
		},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "CATTACATATATACA",
		Description: "Single mm at start",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 0,
				End:   7,
				Dist:  1,
			},
		},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "GATTACA",
		Description: "Identical Sequences, mm allowed",
		MaxDist:     1,
		Expected: []Match{

			Match{
				Start: 0,
				End:   6,
				Dist:  1,
			},
			Match{
				Start: 0,
				End:   7,
				Dist:  0,
			},
		},
	},
	TestCase{
		Pattern:     "ACATCC",
		Text:        "GATTACATATATACATCG",
		Description: "Single mm at end",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 12,
				End:   17,
				Dist:  1,
			},
			Match{
				Start: 12,
				End:   18,
				Dist:  1,
			},
		},
	},
	TestCase{
		Pattern:     "ACATCC",
		Text:        "GATTACATATATGCATCT",
		Description: "Double mm",
		MaxDist:     2,
		Expected: []Match{
			Match{
				Start: 4,
				End:   8,
				Dist:  2,
			},
			Match{
				Start: 4,
				End:   9,
				Dist:  2,
			},
			Match{
				Start: 4,
				End:   10,
				Dist:  2,
			},
			Match{
				Start: 8,
				End:   14,
				Dist:  2,
			},
			Match{
				Start: 12,
				End:   17,
				Dist:  2,
			},
			Match{
				Start: 12,
				End:   18,
				Dist:  2,
			},
		},
	},
	TestCase{
		Pattern:     "GATTACA",
		Text:        "GACTACAATCGGATTTCAACTGA",
		Description: "Find sequence twice, with single mm in each",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 0,
				End:   7,
				Dist:  1,
			},
			Match{
				Start: 11,
				End:   18,
				Dist:  1,
			},
		},
	},
	TestCase{
		Pattern:     "GGGTTLTTSS",
		Text:        "XXXXXXXXXXXXXXXXXXXGGGTTVTTSSAAAAAAAAAAAAAGGGTTVTTSSAAAAAAAAAAAAAAAAAAAAAABBBBBBBBBBBBBBBBBBBBBBBBBGGGTTLTTSS",
		Description: "Protein Search max dist 0",
		MaxDist:     0,
		Expected: []Match{
			Match{
				Start: 99,
				End:   109,
				Dist:  0,
			},
		},
	},
	TestCase{
		Pattern:     "GGGTTLTTSS",
		Text:        "XXXXXXXXXXXXXXXXXXXGGGTTVTTSSAAAAAAAAAAAAAGGGTTVTTSSAAAAAAAAAAAAAAAAAAAAAABBBBBBBBBBBBBBBBBBBBBBBBBGGGTTLTTSS",
		Description: "Protein Search max dist 1",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 19,
				End:   29,
				Dist:  1,
			},
			Match{
				Start: 42,
				End:   52,
				Dist:  1,
			},
			Match{
				Start: 99,
				End:   108,
				Dist:  1,
			},
			Match{
				Start: 99,
				End:   109,
				Dist:  0,
			},
		},
	},
	TestCase{
		Pattern:     "perl",
		Text:        "pearl",
		Description: "Test match with 1 insertion",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 0,
				End:   5,
				Dist:  1,
			},
		},
	},
	TestCase{
		Pattern:     "perl",
		Text:        "prl",
		Description: "Test match with 1 deletion",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 0,
				End:   3,
				Dist:  1,
			},
		},
	},
	TestCase{
		Pattern:     "perl",
		Text:        "perd",
		Description: "Test match with 1 substitution",
		MaxDist:     1,
		Expected: []Match{
			Match{
				Start: 0,
				End:   3,
				Dist:  1,
			},
			Match{
				Start: 0,
				End:   4,
				Dist:  1,
			},
		},
	},
	TestCase{
		Pattern:     "perl",
		Text:        "pesarl",
		Description: "Test match with 2 insertions",
		MaxDist:     2,
		Expected: []Match{
			Match{
				Start: 0,
				End:   2,
				Dist:  2,
			},
			Match{
				Start: 0,
				End:   3,
				Dist:  2,
			},
			Match{
				Start: 0,
				End:   4,
				Dist:  2,
			},
			Match{
				Start: 2,
				End:   6,
				Dist:  2,
			},
		},
	},
	TestCase{
		Pattern:     "perl",
		Text:        "er",
		Description: "Test match with 2 deletions",
		MaxDist:     2,
		Expected: []Match{
			Match{
				Start: 0,
				End:   2,
				Dist:  2,
			},
		},
	},
	TestCase{
		Pattern:     "perl",
		Text:        "berd",
		Description: "Test match with 2 substitutions",
		MaxDist:     2,
		Expected: []Match{
			Match{
				Start: 0,
				End:   3,
				Dist:  2,
			},
			Match{
				Start: 0,
				End:   4,
				Dist:  2,
			},
		},
	},
	TestCase{
		Pattern:     "perl",
		Text:        "prjl",
		Description: "Test match with 1I1D",
		MaxDist:     2,
		Expected: []Match{
			Match{
				Start: 0,
				End:   2,
				Dist:  2,
			},
			Match{
				Start: 0,
				End:   3,
				Dist:  2,
			},
			Match{
				Start: 0,
				End:   4,
				Dist:  2,
			},
		},
	},
	TestCase{
		Pattern:     "perl",
		Text:        "ped",
		Description: "Test match with 1D1S",
		MaxDist:     2,
		Expected: []Match{
			Match{
				Start: 0,
				End:   2,
				Dist:  2,
			},
			Match{
				Start: 0,
				End:   3,
				Dist:  2,
			},
		},
	},
	TestCase{
		Pattern:     "perl",
		Text:        "derJl",
		Description: "Test match with 1I1S",
		MaxDist:     2,
		Expected: []Match{
			Match{
				Start: 0,
				End:   3,
				Dist:  2,
			},
			Match{
				Start: 0,
				End:   4,
				Dist:  2,
			},
			Match{
				Start: 0,
				End:   5,
				Dist:  2,
			},
		},
	},
}

// General function for checking matches
func checkMatches(tCase TestCase, matches []Match, t *testing.T) {
	// Check the zero cases
	if (len(matches) == 0 && len(tCase.Expected) != 0) || (len(matches) != 0 && len(tCase.Expected) == 0) {
		//fmt.Printf("%v\t%v\n", matches, tCase)
		t.Errorf("Bad Number of Matchs: %s\n Found: %v\n Expected %v\n", tCase.Description,
			matches, tCase.Expected)
		return
	}

	// Check that the number of matches returned is expected
	if len(matches) != len(tCase.Expected) {
		t.Errorf("Mismatch in number of matches found vs expected for %s:\nFound: %v\nExpected: %v\n",
			tCase.Description, matches, tCase.Expected)
		return
	}

	// Test that the matches we got are good
	for i, result := range matches {
		// First check if we even expected matches
		if len(tCase.Expected) == 0 {
			//fmt.Printf("%v\t%v\n", matches, tCase)
			t.Errorf("Bad Match: %s\n Found: %v\n Expected %v\n", tCase.Description,
				result, tCase.Expected)
		} else if !(result == tCase.Expected[i]) {
			//fmt.Printf("%v\t%v\n", matches, tCase)
			t.Errorf("Bad Match: %s\n Found: %v\n Expected %v\n", tCase.Description,
				result, tCase.Expected[i])
		}
	}
}

func TestNaiveFind(t *testing.T) {
	for _, tCase := range ExactTestCases {
		matches := naiveFind(tCase.Pattern, tCase.Text)
		checkMatches(tCase, matches, t)
	}
}

func TestNaiveFuzzyFind(t *testing.T) {
	for _, tCase := range ExactTestCases {
		matches := naiveFuzzyFind(tCase.Pattern, tCase.Text, tCase.MaxDist)
		//fmt.Printf("%v\t%v\n", matches, tCase)
		checkMatches(tCase, matches, t)
	}
	for _, tCase := range MismatchTestCases {
		matches := naiveFuzzyFind(tCase.Pattern, tCase.Text, tCase.MaxDist)
		//fmt.Printf("%v\t%v\n", matches, tCase)
		checkMatches(tCase, matches, t)
	}
}

//func TestBMApprox(t *testing.T) {
//	for _, tCase := range ExactTestCases {
//		matches := BMApprox(tCase.Pattern, tCase.Text, tCase.MaxDist)
//		fmt.Printf("%v\t%v\n", matches, tCase)
//		checkMatches(tCase, matches, t)
//	}
//	for _, tCase := range MismatchTestCases {
//		matches := BMApprox(tCase.Pattern, tCase.Text, tCase.MaxDist)
//		fmt.Printf("%v\t%v\n", matches, tCase)
//		checkMatches(tCase, matches, t)
//	}
//}

func TestFuzzyFindShort(t *testing.T) {

	for _, tCase := range ExactTestCases {
		matches, _ := FuzzyFindShort(tCase.Pattern, tCase.Text, tCase.MaxDist, DefaultOptions)
		//fmt.Printf("%v\t%v\n", matches, tCase)
		checkMatches(tCase, matches, t)
	}
	for _, tCase := range EditTestCases {
		matches, _ := FuzzyFindShort(tCase.Pattern, tCase.Text, tCase.MaxDist, DefaultOptions)
		//fmt.Printf("%v\t%v\n", matches, tCase)
		checkMatches(tCase, matches, t)
	}
}

func TestFuzzyFindPigeon(t *testing.T) {

	for _, tCase := range ExactTestCases {
		matches, _ := FuzzyFindPigeon(tCase.Pattern, tCase.Text, tCase.MaxDist, DefaultOptions)
		//fmt.Printf("%v\t%v\n", matches, tCase)
		checkMatches(tCase, matches, t)
	}
	for _, tCase := range EditTestCases {
		matches, _ := FuzzyFindPigeon(tCase.Pattern, tCase.Text, tCase.MaxDist, DefaultOptions)
		//fmt.Printf("%v\t%v\n", matches, tCase)
		checkMatches(tCase, matches, t)
	}
}

func BenchmarkShortPShortTFuzzyFindShort(b *testing.B) {
	// Two mutations, one I one D
	pattern := "TCGTCGTAGCGTC"
	text := "TATAACTCGTCGTAGCGTCAGATGT"
	for i := 0; i < b.N; i++ {
		matches, _ := FuzzyFindShort(pattern, text, 2, DefaultOptions)
		for range matches {

		}
	}
}

func BenchmarkShortPLongTFuzzyFindShort(b *testing.B) {
	// Two mutations, one I one D
	pattern := "TCGTCGGCAGCGTC"
	text := "ACTCANTTATGCATGACTGGCAACAGTCATGTATAACTCGTCGTAGCGTCAGATGTGTATAAGAGACAGCTGTTCTCTCTCTCATCCCAAAACCTTTTGATTCCACTTCTTCCACCA"
	for i := 0; i < b.N; i++ {
		matches, _ := FuzzyFindShort(pattern, text, 2, DefaultOptions)
		for range matches {

		}
	}
}

func BenchmarkLongPShortTFuzzyFindShort(b *testing.B) {
	// Two mutations, one I one D
	pattern := "TCGTCGTAGCGTCGTAGCG"
	text := "TATAACTCGTCGTAGCGTCAGATGT"
	for i := 0; i < b.N; i++ {
		matches, _ := FuzzyFindShort(pattern, text, 2, DefaultOptions)
		for range matches {

		}
	}
}

func BenchmarkLongPLongTFuzzyFindShort(b *testing.B) {
	// Two mutations, one I one D
	pattern := "TCGTCGGCAGCGTC"
	text := "ACTCANTTATGCATGACTGGCAACAGTCATGTATAACTCGTCGTAGCGTCAGATGTGTATAAGAGACAGCTGTTCTCTCTCTCATCCCAAAACCTTTTGATTCCACTTCTTCCACCA"
	for i := 0; i < b.N; i++ {
		matches, _ := FuzzyFindShort(pattern, text, 2, DefaultOptions)
		for range matches {

		}
	}
}

func BenchmarkShortPShortTFuzzyFindPigeon(b *testing.B) {
	// Two mutations, one I one D
	pattern := "TCGTCGTAGCGTC"
	text := "TATAACTCGTCGTAGCGTCAGATGT"
	for i := 0; i < b.N; i++ {
		matches, _ := FuzzyFindPigeon(pattern, text, 2, DefaultOptions)
		for range matches {

		}
	}
}

func BenchmarkShortPLongTFuzzyFindPigeon(b *testing.B) {
	// Two mutations, one I one D
	pattern := "TCGTCGGCAGCGTC"
	text := "ACTCANTTATGCATGACTGGCAACAGTCATGTATAACTCGTCGTAGCGTCAGATGTGTATAAGAGACAGCTGTTCTCTCTCTCATCCCAAAACCTTTTGATTCCACTTCTTCCACCA"
	for i := 0; i < b.N; i++ {
		matches, _ := FuzzyFindPigeon(pattern, text, 2, DefaultOptions)
		for range matches {

		}
	}
}

func BenchmarkLongPShortTFuzzyFindPigeon(b *testing.B) {
	// Two mutations, one I one D
	pattern := "TCGTCGTAGCGTCGTAGCG"
	text := "TATAACTCGTCGTAGCGTCAGATGT"
	for i := 0; i < b.N; i++ {
		matches, _ := FuzzyFindPigeon(pattern, text, 2, DefaultOptions)
		for range matches {

		}
	}
}

func BenchmarkLongPLongTFuzzyFindPigeon(b *testing.B) {
	// Two mutations, one I one D
	pattern := "TCGTCGTAGCGTCAGATGTGTATAAGAGAC"
	text := "ACTCANTTATGCATGACTGGCAACAGTCATGTATAACTCGTCGTAGCGTCAGATGTGTATAAGAGACAGCTGTTCTCTCTCTCATCCCAAAACCTTTTGATTCCACTTCTTCCACCA"
	for i := 0; i < b.N; i++ {
		matches, _ := FuzzyFindPigeon(pattern, text, 2, DefaultOptions)
		for range matches {

		}
	}
}
