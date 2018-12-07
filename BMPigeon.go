package fuzzyfind

/*This package is based off of notes from Ben Langmeads Computational Genomics class:
* http://nbviewer.jupyter.org/github/BenLangmead/comp-genomics-class/blob/master/notebooks/CG_go_006_BoyerMooreApprox.ipynb
* The modification to use an edit distance instead of a hamming distance is my own.
 */
import (
	"bytes"
	"fmt"
)

// Use Z-algorithm to preprocess given string.  See
// Gusfield for complete description of algorithm.
func zArray(s string) []int {
	Z := make([]int, len(s)+1)
	Z[0] = len(s)

	// Initial comparison of s[1:] with prefix
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			Z[1] += 1
		} else {
			break
		}
	}

	r, l := 0, 0
	if Z[1] > 0 {
		r, l = Z[1], 1
	}

	for k := 2; k < len(s); k++ {
		if k > r {
			// Case 1
			for i := k; i < len(s); i++ {
				if s[i] == s[i-k] {
					Z[k] += 1
				} else {
					break
				}
			}
			r, l = k+Z[k]-1, k
		} else {
			// Case 2
			// Calculate length of beta
			nbeta := r - k + 1
			Zkp := Z[k-l]
			if nbeta > Zkp {
				// Case 2a: Zkp wins
				Z[k] = Zkp
			} else {
				// Case 2b: Compare characters just past r
				nmatch := 0
				for i := r + 1; i < len(s); i++ {
					if s[i] == s[i-k] {
						nmatch += 1
					} else {
						break
					}
				}
				l, r = k, r+nmatch
				Z[k] = r - k + 1
			}
		}
	}
	return Z
}

// Helper function that returns a new string
// that is the reverse of the argument
func reverseString(s string) string {
	var buf bytes.Buffer
	for i := 0; i < len(s); i++ {
		buf.WriteByte(s[len(s)-i-1])
	}
	return buf.String()
}

// Helper function that returns a new int slice
// that is the reverse of the argument
func reverseInts(is []int) []int {
	buf := make([]int, len(is))
	for i := 0; i < len(is); i++ {
		buf[len(is)-i-1] = is[i]
	}
	return buf
}

// Make N array (Gusfield theorem 2.2.2) from Z array
func nArray(s string) []int {
	return reverseInts(zArray(reverseString(s)))
}

// Make L' array (Gusfield theorem 2.2.2) using p and N array.
// L'[i] = largest index j less than n such that N[j] = |P[i:]|
func bigLPrimeArray(length int, n []int) []int {
	lp := make([]int, length)
	for j := 0; j < length-1; j++ {
		i := length - n[j]
		if i < length {
			lp[i] = j + 1
		}
	}
	return lp
}

// Compile L array (Gusfield theorem 2.2.2) using p and L' array.
// L[i] = largest index j less than n such that N[j] >= |P[i:]|
func bigLArray(length int, lp []int) []int {
	l := make([]int, length)
	l[1] = lp[1]
	for i := 2; i < length; i++ {
		l[i] = l[i-1]
		if lp[i] > l[i] {
			l[i] = lp[i]
		}
	}
	return l
}

// Compile lp' array (Gusfield theorem 2.2.4) using N array
func smallLPrimeArray(n []int) []int {
	smallLps := make([]int, len(n))
	for i := 0; i < len(n); i++ {
		if n[i] == i+1 { // prefix matching a suffix
			smallLps[len(n)-i-1] = i + 1
		}
	}
	for i := len(n) - 2; i > -1; i-- { // "smear" them out to the left
		if smallLps[i] == 0 {
			smallLps[i] = smallLps[i+1]
		}
	}
	return smallLps
}

// Return tables needed to apply good suffix rule
func goodSuffixTable(p string) ([]int, []int) {
	n := nArray(p)
	lp := bigLPrimeArray(len(p), n)
	return bigLArray(len(p), lp), smallLPrimeArray(n)
}

// Given pattern string and list with ordered alphabet characters, create
// and return a dense bad character table.  Table is indexed by offset
// then by character.
func denseBadCharTab(p string, amap map[byte]int) [][]int {
	tab := make([][]int, len(p))
	nxt := make([]int, len(amap))
	for i, c := range []byte(p) {
		tab[i] = nxt[:]
		nxt[amap[c]] = i + 1
	}
	return tab
}

// Boyer-Moore preprocessing object
type BoyerMoore struct {
	pattern             string
	amap                map[byte]int
	badChars            [][]int
	bigLS, smallLPrimes []int
}

// Constructor for Boyer-Moore preprocessing object
func NewBoyerMoore(p string, alphabet string) *BoyerMoore {
	m := new(BoyerMoore)
	m.pattern = p
	m.amap = make(map[byte]int)
	for i, c := range []byte(alphabet) {
		m.amap[c] = i
	}
	m.badChars = denseBadCharTab(p, m.amap)
	fmt.Printf("P is: %q\n", p)
	m.bigLS, m.smallLPrimes = goodSuffixTable(p)
	return m
}

// Return # skips given by bad character rule at offset i
func (bm *BoyerMoore) badCharacterRule(i int, c byte) int {
	ci := bm.amap[c]
	return i - (bm.badChars[i][ci] - 1)
}

// Given a mismatch at offset i, return amount to shift
// as determined by (weak) good suffix rule.
func (bm *BoyerMoore) goodSuffixRule(i int) int {
	length := len(bm.bigLS)
	if i == length-1 {
		return 0
	}
	i += 1 // i points to leftmost matching position of P
	if bm.bigLS[i] > 0 {
		return length - bm.bigLS[i]
	}
	return length - bm.smallLPrimes[i]
}

// Return amount to shift in case where P matches T
func (bm *BoyerMoore) matchSkip() int {
	return len(bm.smallLPrimes) - bm.smallLPrimes[1]
}

// Utility function for taking max of integers
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Do Boyer-Moore matching
func (bm *BoyerMoore) match(t string) []int {
	i := 0
	occurrences := make([]int, 0)
	p := bm.pattern
	for i < len(t)-len(p)+1 {
		shift, skipGs := 1, 1
		mismatched := false
		for j := len(p) - 1; j >= 0; j-- {
			if p[j] != t[i+j] {
				skipBc := bm.badCharacterRule(j, t[i+j])
				skipGs = bm.goodSuffixRule(j)
				shift = Max(shift, Max(skipBc, skipGs))
				mismatched = true
				break
			}
		}
		if !mismatched {
			occurrences = append(occurrences, i)
			skipGs := Max(shift, skipGs)
			shift = Max(shift, skipGs)
		}
		i += shift
	}
	return occurrences
}

// Return a string splice containing non-overlapping,
// non-empty substrings that cover p.  They should be
// as close to equal-length as possible.
func partition(p string, pieces int) []string {
	base, mod := len(p)/pieces, len(p)%pieces
	fmt.Printf("Base: %d, Mod %d, String, %q\n", base, mod, p)
	idx := 0
	ps := make([]string, 0)
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

// Get the set of letters used in alphabet
func getAlph(p string, t string) string {
	alph := map[rune]int{}
	for _, char := range p {
		alph[char]++
	}
	for _, char := range t {
		alph[char]++
	}
	letters := make([]rune, 0, len(alph))
	for k := range alph {
		letters = append(letters, k)
	}
	return string(letters)
}

// Use the pigeonhole principle together with Boyer-Moore to find
// approximate matches with up to a specified edit distance
func BMApprox(p string, t string, k int) []Match {
	// Check for empty strings first
	if p == "" || t == "" {
		return []Match{}
	}
	alph := getAlph(p, t)
	fmt.Printf("-----------\n")
	fmt.Printf("Alphabet is: %q\n", alph)
	fmt.Printf("Pattern is: %q\n", p)
	fmt.Printf("Text is: %q\n", t)
	fmt.Printf("-----------\n")
	// Short circuit if partitions for pattern will be less than 2
	// Either fall back to raw leven method, or skip this boyer moore nonsense and use a kmer
	// Honestly I just need to write the leven version since my strings will be very close to eachother.
	// Pigeonholing + kmer, what does that look like?
	//  Take pattern and split into k + 1 chunks, see how long the majority of those are, and create kmer index
	//  of the search space. start with first chunk, check for possible star sites, drop the rest in and leven it
	//  or do it the harder way and try to extend, both require a leven with a short circuit
	//  This might be faster than fuzzysearch ngram
	//  Or just do a find on the string for the subseq and start trying to exted from there... depends how long string is
	//  whether the kmer version is worth or not
	// 5 versions:
	//    1. leven only
	//    2. Normal exact search if no editdist
	//    3. pigeon + kmer + leven extend
	//    4. pigeon + boyer + leven extend
	//    5. pigeon + exact + leven extend
	if len(p)/(k+1) >= 2 {
		fmt.Printf("Need to handle this special case\n")
		return []Match{}
	}
	ps := partition(p, k+1)            // split p into k+1 non-empty, non-overlapping substrings
	off := 0                           // offset into p of current partition
	occurrences := make(map[int]Match) // we might see the same occurrence >1 time
	for _, pi := range ps {            // for each partition
		fmt.Printf("Partition is: %q\n", pi)
		bm := NewBoyerMoore(pi, alph) // BM preprocess the partition
		for _, hit := range bm.match(t) {
			if hit-off < 0 {
				continue // pattern falls off left end of T?
			}
			if hit+len(p)-off > len(t) {
				continue // falls off right end?
			}
			// Count mismatches to left and right of the matching partition
			nmm := 0
			for i := 0; i < off; i++ {
				if t[hit-off+i] != p[i] {
					nmm++
					if nmm > k {
						break // exceeded maximum # mismatches
					}
				}
			}
			if nmm <= k {
				for i := off + len(pi); i < len(p); i++ {
					if t[hit-off+i] != p[i] {
						nmm++
						if nmm > k {
							break // exceeded maximum # mismatches
						}
					}
				}
			}
			if nmm <= k {
				occurrences[hit-off] = Match{Start: hit, End: hit + len(p) - off, Dist: nmm} // approximate match
			}
		}
		off += len(pi) // Update offset of current partition
	}
	var occurrence_list []Match
	for _, v := range occurrences {
		occurrence_list = append(occurrence_list, v)
	}
	return occurrence_list
}
