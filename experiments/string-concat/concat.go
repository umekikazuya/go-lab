package stringconcat

import "strings"

// ConcatPlus concatenates parts using the + operator in a loop.
// Each iteration allocates a new string, resulting in O(N²) copy cost.
func ConcatPlus(parts []string) string {
	s := ""
	for _, p := range parts {
		s += p
	}
	return s
}

// ConcatBuilder concatenates parts using strings.Builder.
// The internal buffer grows with a doubling strategy: O(log N) allocations, O(N) copies.
func ConcatBuilder(parts []string) string {
	var b strings.Builder
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

// ConcatBuilderGrow concatenates parts using strings.Builder with pre-allocated capacity.
// A single allocation is performed upfront, eliminating regrowth overhead.
func ConcatBuilderGrow(parts []string) string {
	var b strings.Builder
	total := 0
	for _, p := range parts {
		total += len(p)
	}
	b.Grow(total)
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}
