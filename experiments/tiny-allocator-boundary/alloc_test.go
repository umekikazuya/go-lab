package tinyallocator

import (
	"fmt"
	"runtime"
	"runtime/metrics"
	"testing"
	"unsafe"
)

// --- Benchmarks: ns/op and allocs/op across the 16B boundary ---

// benchSize runs one allocation benchmark, naming the sub-benchmark after the
// runtime size of T so the boundary is readable in the output.
func benchSize[T any](b *testing.B, note ...string) {
	var zero T
	name := fmt.Sprintf("%02dB", unsafe.Sizeof(zero))
	if len(note) > 0 {
		name += "_" + note[0]
	}
	b.Run(name, func(b *testing.B) {
		for b.Loop() {
			_ = Alloc[T]()
		}
	})
}

// BenchmarkNoscan sweeps pointer-free sizes in 1B steps across the
// maxTinySize=16 boundary, then jumps to the neighbouring size classes.
func BenchmarkNoscan(b *testing.B) {
	benchSize[B1](b)
	benchSize[B2](b)
	benchSize[B3](b)
	benchSize[B4](b)
	benchSize[B5](b)
	benchSize[B6](b)
	benchSize[B7](b)
	benchSize[B8](b)
	benchSize[B9](b)
	benchSize[B10](b)
	benchSize[B11](b)
	benchSize[B12](b)
	benchSize[B13](b)
	benchSize[B14](b)
	benchSize[B15](b)
	benchSize[B16](b)
	benchSize[B17](b)
	benchSize[B18](b)
	benchSize[B20](b)
	benchSize[B24](b)
	benchSize[B25](b)
	benchSize[B32](b)
	benchSize[B33](b)
}

// BenchmarkPointer is the control group: identical sizes, but each type holds
// a pointer, so the tiny allocator is never used.
func BenchmarkPointer(b *testing.B) {
	benchSize[P8](b)
	benchSize[P15](b, "from15") // 15B of payload, padded up to 16B
	benchSize[P16](b)
	benchSize[P24](b)
	benchSize[P32](b)
}

// --- Direct verification of the runtime counters ---

// counters holds the deltas observed while allocating measureN objects.
//
// runtime.MemStats has no exported tiny-allocation field (memstats.tinyallocs
// is folded into Mallocs), so the tiny counter is read from runtime/metrics
// instead.
type counters struct {
	mallocs   uint64 // total objects allocated
	tiny      uint64 // objects served by the tiny allocator
	heapBytes uint64 // heap growth
}

const (
	measureN = 1 << 16

	metricObjects = "/gc/heap/allocs:objects"
	metricTiny    = "/gc/heap/tiny/allocs:objects"
)

// readCounters samples the two allocation-count metrics plus HeapAlloc.
func readCounters() counters {
	samples := []metrics.Sample{
		{Name: metricObjects},
		{Name: metricTiny},
	}
	metrics.Read(samples)

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	return counters{
		mallocs:   samples[0].Value.Uint64(),
		tiny:      samples[1].Value.Uint64(),
		heapBytes: ms.HeapAlloc,
	}
}

// measure allocates measureN values of T, keeping every pointer alive, and
// returns the counter deltas.
func measure[T any]() counters {
	keep := make([]*T, measureN) // preallocated: excluded from the delta

	runtime.GC()
	before := readCounters()

	for i := range keep {
		keep[i] = Alloc[T]()
	}

	after := readCounters()
	runtime.KeepAlive(keep)

	return counters{
		mallocs:   after.mallocs - before.mallocs,
		tiny:      after.tiny - before.tiny,
		heapBytes: after.heapBytes - before.heapBytes,
	}
}

// TestTinyPacking measures how many bytes of heap each object actually costs.
//
// Being routed to the tiny path is not the same as benefiting from it. A tiny
// block is 16B wide, so a run of identical S-byte objects packs
// floor(16/align(S)) per block; once that count drops to 1 the tiny path costs
// exactly what a 16B size class costs. wantBytes below is that arithmetic, and
// the tiny counter reports only the objects that were packed into an already
// open block — allocations that had to open a fresh block are counted as
// ordinary mallocs instead.
func TestTinyPacking(t *testing.T) {
	tests := []struct {
		name       string
		measure    func() counters
		wantBytes  float64 // heap bytes per object
		wantPacked bool    // some objects shared a block
	}{
		{"noscan/01B", measure[B1], 1.0, true}, // 16 per block
		{"noscan/03B", measure[B3], 3.2, true}, // 5 per block
		{"noscan/05B", measure[B5], 16.0 / 3, true},
		{"noscan/06B", measure[B6], 8.0, true},   // align 2 -> 2 per block
		{"noscan/08B", measure[B8], 8.0, true},   // 2 per block
		{"noscan/09B", measure[B9], 16.0, false}, // 1 per block: no gain
		{"noscan/15B", measure[B15], 16.0, false},
		{"noscan/16B", measure[B16], 16.0, false}, // size class 16
		{"noscan/17B", measure[B17], 24.0, false}, // size class 24
		{"pointer/08B", measure[P8], 8.0, false},
		{"pointer/16B", measure[P16], 16.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.measure()
			got := float64(c.heapBytes) / measureN
			t.Logf("mallocs=%d tiny=%d heap=%dB (%.2f B/obj)",
				c.mallocs, c.tiny, c.heapBytes, got)

			if diff := got - tt.wantBytes; diff < -0.05 || diff > 0.05 {
				t.Errorf("heap = %.2f B/obj, want %.2f", got, tt.wantBytes)
			}

			// The tiny counter lives in the mcache and is flushed lazily, so
			// only a lower bound is asserted for the packed cases.
			if tt.wantPacked && c.tiny < measureN/4 {
				t.Errorf("tiny allocs = %d, want >= %d", c.tiny, measureN/4)
			}
			if !tt.wantPacked && c.tiny != 0 {
				t.Errorf("tiny allocs = %d, want 0", c.tiny)
			}
		})
	}
}

// TestAllocEscapes guarantees the benchmarks measure heap allocations rather
// than stack slots elided by escape analysis.
func TestAllocEscapes(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
	}{
		{"noscan/01B", func() { _ = Alloc[B1]() }},
		{"noscan/16B", func() { _ = Alloc[B16]() }},
		{"pointer/08B", func() { _ = Alloc[P8]() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testing.AllocsPerRun(1000, tt.fn); got != 1 {
				t.Errorf("allocs per run = %v, want 1", got)
			}
		})
	}
}
