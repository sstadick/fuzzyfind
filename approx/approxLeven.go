package approx

import (
	"fmt"
	"io"
	"os"
)

const MaxInt = int(^uint(0) >> 1)

// LevenContext provides a reusable int Matrix
type LevenContext struct {
	matrix [][]int
}

func (c *LevenContext) getMatrix(height int) [][]int {
	if cap(c.matrix) < height {
		c.matrix = make([][]int, height)
	}
	return c.matrix[:height][:height]
}

// approxLeven is a wrapper for calling the distance function with the context struct
func approxLeven(pattern []rune, text []rune, maxE int, op Options) ([]Match, error) {
	c := LevenContext{}
	return c.ApproxLeven(string(pattern), string(text), maxE, op)
}

// Use the leven alogorithm to find the best match of pattern in text with up to maxE edit dist
// It is expected that the user has already normalized: https://blog.golang.org/normalization
// Leven adapted from https://github.com/texttheater/golang-levenshtein/blob/master/levenshtein/levenshtein.go
// Note this is a little wierd right now because I'm hijacking what I had in order to add a context struct
func (c *LevenContext) ApproxLeven(p string, t string, maxE int, op Options) ([]Match, error) {

	// Check for empty strings first
	if p == "" {
		return nil, fmt.Errorf("pattern to search empty")
	} else if t == "" {
		return nil, fmt.Errorf("text to search is empty")
	}
	pattern := []rune(p)
	text := []rune(t)
	height := len(pattern) + 1
	width := len(text) + 1
	matrix := c.getMatrix(height)

	// Initialize trivial distances (from/to empty string). That is, fill
	// the left column and the top row with row/column indices.
	for i := 0; i < height; i++ {
		matrix[i] = make([]int, width)
		matrix[i][0] = i
	}
	// Set the top row to 0's
	for j := 1; j < width; j++ {
		matrix[0][j] = 0
	}

	// Fill in the remaining cells: for each prefix pair, choose the
	// (edit history, operation) pair with the lowest cost.
	for i := 1; i < height; i++ {
		currentMin := MaxInt
		for j := 1; j < width; j++ {
			delCost := matrix[i-1][j] + op.DelCost
			matchSubCost := matrix[i-1][j-1]
			if !op.Matches(pattern[i-1], text[j-1]) {
				matchSubCost += op.SubCost
			}
			insCost := matrix[i][j-1] + op.InsCost
			matrix[i][j] = min(delCost, min(matchSubCost,
				insCost))
			if matrix[i][j] < currentMin {
				currentMin = matrix[i][j]
			}
		}
		// Check to see if the min for the row is greater than the
		// max allowed
		if currentMin > maxE {
			return nil, nil
		}
	}
	//LogMatrix(pattern, text, matrix)
	// Return a traceback for each alignment less than maxE
	minCols := []int{}
	for j := 0; j <= len(text); j++ {
		if matrix[len(pattern)][j] <= maxE {
			minCols = append(minCols, j)
		}
	}

	// If we somehow end up with no good matches
	//	if len(minCols) == 0 {
	//		return []Match{}
	//	}
	matches, err := trace(matrix, pattern, text, minCols, op)
	if err != nil {
		return matches, fmt.Errorf("can't traceback matches: %v", err)
	}
	return matches, nil

}

// Traceback to find all the lowest edit distances
func trace(matrix [][]int, p []rune, t []rune, minCols []int, op Options) ([]Match, error) {
	// For each min alignment found, do a traceback
	// I need the start, and end releative to the text, and the distance
	// I have the end and the dist, just need the start
	matches := []Match{}
	for _, min := range minCols {
		// Set the 'corner' that we will start looking in
		i, j := len(p), min

		for i > 0 {
			diag, vert, horz := MaxInt, MaxInt, MaxInt
			var delt int
			if i > 0 && j > 0 {
				if op.Matches(p[i-1], t[j-1]) {
					delt = 0
				} else {
					delt = op.SubCost
				}
				diag = matrix[i-1][j-1] + delt
			}
			if i > 0 {
				vert = matrix[i-1][j] + op.DelCost
			}
			if j > 0 {
				horz = matrix[i][j-1] + op.InsCost
			}
			if diag <= vert && diag <= horz {
				// diagonal was best, it was a match or mismatch
				i--
				j--
			} else if vert <= horz {
				// vertical was best, it was an insertion
				i--
			} else {
				// horizontal was best, it was a deletion
				j--
			}
		}
		matches = append(matches, Match{Start: j, End: min, Dist: matrix[len(p)][min]})
	}
	return matches, nil
}

// WriteMatrix writes a visual representation of the given matrix for the given
// strings to the given writer.
func WriteMatrix(pattern []rune, text []rune, matrix [][]int, writer io.Writer) {
	fmt.Fprintf(writer, "    ")
	for _, textRune := range text {
		fmt.Fprintf(writer, "  %c", textRune)
	}
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "  %2d", matrix[0][0])
	for j := range text {
		fmt.Fprintf(writer, " %2d", matrix[0][j+1])
	}
	fmt.Fprintf(writer, "\n")
	for i, patternRune := range pattern {
		fmt.Fprintf(writer, "%c %2d", patternRune, matrix[i+1][0])
		for j := range text {
			fmt.Fprintf(writer, " %2d", matrix[i+1][j+1])
		}
		fmt.Fprintf(writer, "\n")
	}

}

// LogMatrix writes a visual representation of the given matrix for the given
// strings to os.Stderr. This function is deprecated, use
// WriteMatrix(pattern, text, matrix, os.Stderr) instead.
func LogMatrix(pattern []rune, text []rune, matrix [][]int) {
	WriteMatrix(pattern, text, matrix, os.Stderr)
}

func min(a int, b int) int {
	if b < a {
		return b
	}
	return a
}

func max(a int, b int) int {
	if b > a {
		return b
	}
	return a
}
