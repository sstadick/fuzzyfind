# fuzzyfind
A clone of python fuzzysearch library in Go.

It is based off the v0.1.0 version of the Python [fuzzysearch library](https://github.com/taleinat/fuzzysearch/blob/v0.1.0/fuzzysearch/fuzzysearch.py) by [Tal Einat](https://github.com/taleinat).

## Synopsis
fuzzyfind uses a modified levenshtein algorithm to find approximate matches of a subsequence in a sequence. My reason for using this tools is for extracting regions of sequencing reads that might have mutations. 

### If you want to .... just get going:
Use the FuzzyFind method, it will choose the best method for you depending on your pattern and text sizes

### If you want to .... match the same pattern against multiple texts:
This has yet to be implemented. It will likely use boyer moore to create a lookup table for the pattern.

### If you want to .... match multiple patterns against the same text:
This has yet to be implemented. It will likely use a kmer index of the text

### If you want to .... specifically use just the modified levenshtien algorithm:
Use FuzzyFindShort. This should work best on short patterns.

### if you want to .... specifically use the pigeonhole method:
use FuzzyFindPigeon. This will only work when you have a `len(pattern) / (maxDist + 1) >= 1` and should really only be used when greater than 3.
This is also not yet implemented.

## Install
`go get github.com/sstadick/fuzzyfind`

## Usage in code
```Go
package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sstadick/fuzzyfind"
)

func main() {
	haystack := "GACTAGCACTGTAGGGATAACAATTTCACACAGGTGGACAATTACATTGAAAATCACAGATTGGTCACACACACATTGGACATACATAGAAACACACACACATACATTAGATACGAACATAGAAACACACATTAGACGCGTACATAGACACAAACACATTGACAGGCAGTTCAGATGATGACGCCCGACTGATACTCGCGTAGTCGTGGGAGGCAAGGCACACAGGGGATAGG"
	needle := "TGCACTGTAGGGATAACAAT" // distance 1
	maxDist := 2
	result, _ := fuzzyfind.FuzzyFindShort(needle, haystack, maxDist, fuzzyfind.DefaultOptions)
	spew.Dump(result)
}
```

Output

```
([]fuzzyfind.Match) (len=4 cap=4) {
 (fuzzyfind.Match) {
  Start: (int) 6,
  End: (int) 24,
  Dist: (int) 2
 },
 (fuzzyfind.Match) {
  Start: (int) 5,
  End: (int) 24,
  Dist: (int) 1
 },
 (fuzzyfind.Match) {
  Start: (int) 3,
  End: (int) 24,
  Dist: (int) 1
 },
 (fuzzyfind.Match) {
  Start: (int) 3,
  End: (int) 24,
  Dist: (int) 2
 }
}
```

## Futher readings
- [Python Version](https://github.com/taleinat/fuzzysearch) (currently I'm based on v0.1.0)
- [Javascript Version](https://github.com/taleinat/levenshtein-search), which has some features not yet in the Python version
- [Tal's Slides](https://taleinat.github.io/playing_with_cython/)
- [Levenshtein](https://en.wikipedia.org/wiki/Levenshtein_distance)
- ngrams - To be researched
- regex engine implementation - To be researched
  - [TRE](https://laurikari.net/tre/) engine
  - My [blog](https://ducktape.blot.im/tre-a-regex-engine-with-approximate-matching) post on TRE

## Notes
- Allow a custom penalty matrix like smith-waterman and co?
