package cluster

import (
	"hash/fnv"
	"strings"
	"unicode"
)

func Simhash(text string, weight int) uint64 {
	if weight < 1 {
		weight = 1
	}
	var vec [64]int
	for _, tok := range tokenize(text) {
		h := hashToken(tok)
		for i := 0; i < 64; i++ {
			if (h>>i)&1 == 1 {
				vec[i] += weight
			} else {
				vec[i] -= weight
			}
		}
	}
	var out uint64
	for i := 0; i < 64; i++ {
		if vec[i] > 0 {
			out |= 1 << i
		}
	}
	return out
}

func Hamming(a, b uint64) int {
	x := a ^ b
	count := 0
	for x != 0 {
		x &= x - 1
		count++
	}
	return count
}

func hashToken(tok string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(tok))
	return h.Sum64()
}

func tokenize(text string) []string {
	var out []string
	var b strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		if b.Len() > 0 {
			out = append(out, b.String())
			b.Reset()
		}
	}
	if b.Len() > 0 {
		out = append(out, b.String())
	}
	return out
}
