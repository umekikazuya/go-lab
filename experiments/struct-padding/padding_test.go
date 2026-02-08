package padding

import (
	"testing"
	"unsafe"
)

// TestSize statically verifies struct sizes predicted by the hypothesis.
func TestSize(t *testing.T) {
	tests := []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{"Unpadded", unsafe.Sizeof(Unpadded{}), 24},
		{"Padded", unsafe.Sizeof(Padded{}), 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("unsafe.Sizeof(%s{}) = %d B, want %d B", tt.name, tt.got, tt.want)
			}
			t.Logf("unsafe.Sizeof(%s{}) = %d B", tt.name, tt.got)
		})
	}

	diff := float64(unsafe.Sizeof(Unpadded{})-unsafe.Sizeof(Padded{})) / float64(unsafe.Sizeof(Unpadded{})) * 100
	t.Logf("Size reduction: %.1f%%", diff)
}

// TestFieldOffsets logs the offset of each field for visual verification.
func TestFieldOffsets(t *testing.T) {
	t.Log("=== Unpadded ===")
	t.Logf("  a bool   offset=%d  size=%d", unsafe.Offsetof(Unpadded{}.a), unsafe.Sizeof(Unpadded{}.a))
	t.Logf("  b int64  offset=%d  size=%d", unsafe.Offsetof(Unpadded{}.b), unsafe.Sizeof(Unpadded{}.b))
	t.Logf("  c bool   offset=%d  size=%d", unsafe.Offsetof(Unpadded{}.c), unsafe.Sizeof(Unpadded{}.c))
	t.Logf("  d int32  offset=%d  size=%d", unsafe.Offsetof(Unpadded{}.d), unsafe.Sizeof(Unpadded{}.d))
	t.Logf("  total=%d", unsafe.Sizeof(Unpadded{}))

	t.Log("=== Padded ===")
	t.Logf("  b int64  offset=%d  size=%d", unsafe.Offsetof(Padded{}.b), unsafe.Sizeof(Padded{}.b))
	t.Logf("  d int32  offset=%d  size=%d", unsafe.Offsetof(Padded{}.d), unsafe.Sizeof(Padded{}.d))
	t.Logf("  a bool   offset=%d  size=%d", unsafe.Offsetof(Padded{}.a), unsafe.Sizeof(Padded{}.a))
	t.Logf("  c bool   offset=%d  size=%d", unsafe.Offsetof(Padded{}.c), unsafe.Sizeof(Padded{}.c))
	t.Logf("  total=%d", unsafe.Sizeof(Padded{}))
}

// ---------------------------------------------------------------------------
// Benchmarks: mass allocation
// ---------------------------------------------------------------------------

const N = 1024 * 1024 // 1M elements

// BenchmarkAllocUnpadded measures allocation throughput for Unpadded structs.
func BenchmarkAllocUnpadded(b *testing.B) {
	for b.Loop() {
		s := make([]Unpadded, N)
		_ = s
	}
}

// BenchmarkAllocPadded measures allocation throughput for Padded structs.
func BenchmarkAllocPadded(b *testing.B) {
	for b.Loop() {
		s := make([]Padded, N)
		_ = s
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: slice traversal (cache-line effect)
// ---------------------------------------------------------------------------

// sink prevents dead-code elimination by the compiler.
var sink int64

// BenchmarkTraverseUnpadded measures traversal throughput over a large
// slice of Unpadded structs to expose cache-line inefficiency.
func BenchmarkTraverseUnpadded(b *testing.B) {
	data := make([]Unpadded, N)
	for i := range data {
		data[i] = Unpadded{a: true, b: int64(i), c: false, d: int32(i)}
	}
	b.ResetTimer()
	for b.Loop() {
		var acc int64
		for i := range data {
			acc += data[i].b + int64(data[i].d)
		}
		sink = acc
	}
}

// BenchmarkTraversePadded measures traversal throughput over a large
// slice of Padded structs to expose cache-line efficiency.
func BenchmarkTraversePadded(b *testing.B) {
	data := make([]Padded, N)
	for i := range data {
		data[i] = Padded{b: int64(i), d: int32(i), a: true, c: false}
	}
	b.ResetTimer()
	for b.Loop() {
		var acc int64
		for i := range data {
			acc += data[i].b + int64(data[i].d)
		}
		sink = acc
	}
}
