package stringconcat

import (
	"fmt"
	"strings"
	"testing"
)

// makeParts generates n fixed-length (8 byte) string parts.
func makeParts(n int) []string {
	parts := make([]string, n)
	for i := range parts {
		parts[i] = strings.Repeat("a", 8)
	}
	return parts
}

var sizes = []int{2, 4, 8, 16, 32, 64}

func BenchmarkConcatPlus(b *testing.B) {
	for _, n := range sizes {
		parts := makeParts(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = ConcatPlus(parts)
			}
		})
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	for _, n := range sizes {
		parts := makeParts(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = ConcatBuilder(parts)
			}
		})
	}
}

func BenchmarkConcatBuilderGrow(b *testing.B) {
	for _, n := range sizes {
		parts := makeParts(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = ConcatBuilderGrow(parts)
			}
		})
	}
}
