package kenlm

import "strings"

// ScoreChars scores a string as a space-separated sequence of Unicode
// runes — the standard input convention for character-level n-gram LMs.
// Whitespace runes in the input are dropped (they would tokenize to empty
// fields after the join).
func (m *Model) ScoreChars(s string) float64 {
	var b strings.Builder
	b.Grow(len(s) * 2)
	first := true
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			continue
		}
		if !first {
			b.WriteByte(' ')
		}
		b.WriteRune(r)
		first = false
	}
	return m.Score(b.String())
}
