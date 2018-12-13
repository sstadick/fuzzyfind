[![Documentation](https://godoc.org/github.com/sstadick/fuzzyfind?status.svg)](http://godoc.org/github.com/sstadick/fuzzyfind)
[![Go Report Card](https://goreportcard.com/badge/github.com/sstadick/fuzzyfind)](https://goreportcard.com/report/github.com/sstadick/fuzzyfind)

# approx
fuzzyfind/approx is an approximate string matching tool. It is heavily inspired by Python [fuzzysearch library](https://github.com/taleinat/fuzzysearch/blob/v0.1.0/fuzzysearch/fuzzysearch.py) by [Tal Einat](https://github.com/taleinat).

What makes this library differnt than other \*fuzzy\* libraries is that it returns the start and end index's of your pattern in the text that you search. The long term goal is to add more variants of this search for different use cases. For now though, ApproxFind will work well.

## Synopsis
ApproxFind uses a modified levenshtein algorithm to find approximate matches of a subsequence in a sequence. My reason for using this tools is for extracting regions of sequencing reads that might have mutations. 


## Benchmarks
```
Short Pattern Short Text:
BenchmarkShortPShortTApproxFind-8    	  300000	      4301 ns/op
BenchmarkShortPShortTapproxPigeon-8   	  100000	     13708 ns/op

Short Pattern Long  Text:
BenchmarkShortPLongTApproxFind-8     	  100000	     18044 ns/op
BenchmarkShortPLongTapproxPigeon-8    	  100000	     17393 ns/op

Long  Pattern Short Text:
BenchmarkLongPShortTApproxFind-8     	  300000	      4567 ns/op
BenchmarkLongPShortTapproxPigeon-8    	  200000	     10575 ns/op

Long  Pattern Long  Text:
BenchmarkLongPLongTApproxFind-8      	  100000	     18173 ns/op
BenchmarkLongPLongTapproxPigeon-8     	   30000	     46241 ns/op

```

### This is a work in progress
For now just use the ApproxFind function, it has solid performance and is the most correct.

## Install
`go get github.com/sstadick/fuzzyfind`

## Usage in code
```Go
package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sstadick/fuzzyfind/approx"
)

func main() {
	haystack := "GACTAGCACTGTAGGGATAACAATTTCACACAGGTGGACAATTACATTGAAAATCACAGATTGGTCACACACACATTGGACATACATAGAAACACACACACATACATTAGATACGAACATAGAAACACACATTAGACGCGTACATAGACACAAACACATTGACAGGCAGTTCAGATGATGACGCCCGACTGATACTCGCGTAGTCGTGGGAGGCAAGGCACACAGGGGATAGG"
	needle := "TGCACTGTAGGGATAACAAT" // distance 1
	maxDist := 2
	result, _ := approx.ApproxFind(needle, haystack, maxDist, approx.DefaultOptions)
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

### If you want to .... just get going:
Use the ApproxFind method, it will choose the best method for you depending on your pattern and text sizes

### If you want to .... match the same pattern against multiple texts:
This has yet to be implemented. It will likely use boyer moore to create a lookup table for the pattern.

### If you want to .... match multiple patterns against the same text:
This has yet to be implemented. It will likely use a kmer index of the text

### If you want to .... specifically use just the modified levenshtien algorithm:
Use ApproxFind. This should work best on short patterns.

### if you want to .... specifically use the pigeonhole method:
use approxPigeon. This will only work when you have a `len(pattern) / (maxDist + 1) >= 1` and should really only be used when greater than 3. Use this one at your own risk as it is not as well tested as ApproxFind, and seems to be slower in almost all cases right now.

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
