// Package tinyallocator measures the allocation cost boundary of Go's tiny
// allocator.
//
// runtime.mallocgc routes an object to the tiny allocator only when both
// conditions hold:
//
//	size < maxTinySize (16)
//	the type is noscan (contains no pointers)
//
// Objects on that path are packed into a shared 16B block, so several
// allocations can share one span slot. Everything else takes the regular
// mcache path and consumes a full size-class slot.
package tinyallocator

// Alloc heap-allocates a single T and returns it.
//
// The pointer is returned across a non-inlinable call boundary, which forces
// escape analysis to move the object to the heap regardless of T's size.
//
//go:noinline
func Alloc[T any]() *T {
	return new(T)
}

// --- noscan types: eligible for the tiny allocator while size < 16 ---

type (
	B1  [1]byte
	B2  [2]byte
	B3  [3]byte
	B4  [4]byte
	B5  [5]byte
	B6  [6]byte
	B7  [7]byte
	B8  [8]byte
	B9  [9]byte
	B10 [10]byte
	B11 [11]byte
	B12 [12]byte
	B13 [13]byte
	B14 [14]byte
	B15 [15]byte
	B16 [16]byte
	B17 [17]byte
	B18 [18]byte
	B20 [20]byte
	B24 [24]byte
	B25 [25]byte
	B32 [32]byte
	B33 [33]byte
)

// --- pointer-bearing types: never eligible, whatever the size ---
//
// A struct holding a pointer cannot be 15B: the pointer forces 8B alignment,
// so P15 is padded up to 16B. The smallest pointer-bearing size is 8B, which
// is why P8 is the control for the noscan 8B case.

type (
	P8  struct{ p *int }
	P15 struct {
		p *int
		_ [7]byte
	}
	P16 struct {
		p *int
		_ [8]byte
	}
	P24 struct {
		p *int
		_ [16]byte
	}
	P32 struct {
		p *int
		_ [24]byte
	}
)
