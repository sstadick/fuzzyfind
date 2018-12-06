# fuzzyfind
A clone of python fuzzysearch library in Go.

It is based off the v0.1.0 version of the Python [fuzzysearch library](https://github.com/taleinat/fuzzysearch/blob/v0.1.0/fuzzysearch/fuzzysearch.py) by [Tal Einat](https://github.com/taleinat).

## Synopsis
fuzzyfind uses a modified levenshtein algorithm to find approximate matches of a subsequence in a sequence. My reason for using this tools is for extracting regions of sequencing reads that might have mutations. 

Why is fuzzysearch special? There are TONS of fuzzy / Levenshtein implementations out there, but fuzzysearch is the only one I've found that returns both the start/end index of the match, and the Levenshtein score for the match.

## Development todos
Bring the features in line with the original Python library and Tal's Javascript port as well.

## Futher readings
- [Python Version](https://github.com/taleinat/fuzzysearch) (currently I'm based on v0.1.0)
- [Javascript VErsion](https://github.com/taleinat/levenshtein-search), which has some features not yet in the Python version
- [Tal's Slides](https://taleinat.github.io/playing_with_cython/)
- [Levenshtein](https://en.wikipedia.org/wiki/Levenshtein_distance)
- ngrams - To be researched
- regex engine implementation - To be researched
  - [TRE](https://laurikari.net/tre/) engine
  - By [blog](https://ducktape.blot.im/tre-a-regex-engine-with-approximate-matching) post on TRE
- Local alignment algorithms
- Longest Common string algorithms


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
	result, _ := fuzzyfind.FindNearMatches(needle, haystack, maxDist)
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
