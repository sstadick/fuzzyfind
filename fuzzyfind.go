package fuzzyfind

import (
    "fmt"
    "os"
    "strings"
)


type Candidate struct {
    Start int
    SubseqIndex int
    Dist int
}

type Match struct {
    Start int
    End int
    Dist int
}


func prepareInitCandidatesMap(subsequence string, maxLDist int) {
    candidateMap := make(Map[rune]int)
    for index, char := range subsequence[:maxLDist + 1] {
        if ! strings.ContainsRune(subsequence[:index] , char) {
            fmt.Fprintf(os.stderr("%s not in %s", char subsequence[:index]))
            candidateMap[char] = index
        }
    } 
}
