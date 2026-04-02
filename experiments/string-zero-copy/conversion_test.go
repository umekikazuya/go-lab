package stringzerocopy

import (
	"fmt"
	"strings"
	"testing"
)

var sizes = []int{8, 64, 512, 4096}

func makeBytes(n int) []byte  { return []byte(strings.Repeat("a", n)) }
func makeString(n int) string { return strings.Repeat("a", n) }
func makeMap(n int) map[string]int {
	m := map[string]int{strings.Repeat("a", n): 1}
	return m
}

// --- Copy patterns (expect allocs/op = 1) ---

func BenchmarkBytesToStringAssign(b *testing.B) {
	for _, n := range sizes {
		bs := makeBytes(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = BytesToStringAssign(bs)
			}
		})
	}
}

func BenchmarkStringToBytesAssign(b *testing.B) {
	for _, n := range sizes {
		s := makeString(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = StringToBytesAssign(s)
			}
		})
	}
}

func BenchmarkBytesToStringConcat(b *testing.B) {
	for _, n := range sizes {
		bs := makeBytes(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = BytesToStringConcat(bs)
			}
		})
	}
}

// --- Zero-copy patterns (expect allocs/op = 0) ---

func BenchmarkBytesToStringMapLookup(b *testing.B) {
	for _, n := range sizes {
		bs := makeBytes(n)
		m := makeMap(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = BytesToStringMapLookup(m, bs)
			}
		})
	}
}

func BenchmarkStringToBytesRange(b *testing.B) {
	for _, n := range sizes {
		s := makeString(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			for b.Loop() {
				StringToBytesRange(s)
			}
		})
	}
}

func BenchmarkBytesToStringCompare(b *testing.B) {
	for _, n := range sizes {
		bs := makeBytes(n)
		target := makeString(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = BytesToStringCompare(bs, target)
			}
		})
	}
}
