package mapkey

import (
	"strconv"
	"testing"
)

const mapSize = 10_000

// sink prevents dead-code elimination by the compiler.
var sink int

// ---------------------------------------------------------------------------
// Test: verify all three key strategies produce correct lookups
// ---------------------------------------------------------------------------

func TestMapKeyLookup(t *testing.T) {
	t.Run("StringKey", func(t *testing.T) {
		m := make(map[string]int, mapSize)
		for i := range mapSize {
			m[StringKey(i, "abc")] = i
		}
		for i := range mapSize {
			v, ok := m[StringKey(i, "abc")]
			if !ok || v != i {
				t.Fatalf("lookup failed for key %d", i)
			}
		}
		t.Logf("map[string]int: %d entries OK", mapSize)
	})

	t.Run("CompositeKey", func(t *testing.T) {
		m := make(map[CompositeKey]int, mapSize)
		for i := range mapSize {
			m[CompositeKey{ID: i, Code: "abc"}] = i
		}
		for i := range mapSize {
			v, ok := m[CompositeKey{ID: i, Code: "abc"}]
			if !ok || v != i {
				t.Fatalf("lookup failed for key %d", i)
			}
		}
		t.Logf("map[CompositeKey]int: %d entries OK", mapSize)
	})

	t.Run("IntPairKey", func(t *testing.T) {
		m := make(map[IntPairKey]int, mapSize)
		for i := range mapSize {
			m[IntPairKey{X: i, Y: i * 7}] = i
		}
		for i := range mapSize {
			v, ok := m[IntPairKey{X: i, Y: i * 7}]
			if !ok || v != i {
				t.Fatalf("lookup failed for key %d", i)
			}
		}
		t.Logf("map[IntPairKey]int: %d entries OK", mapSize)
	})
}

// ---------------------------------------------------------------------------
// Benchmarks: Insert (key construction included)
// ---------------------------------------------------------------------------

func BenchmarkInsert_StringKey(b *testing.B) {
	for b.Loop() {
		m := make(map[string]int, mapSize)
		for i := range mapSize {
			m[StringKey(i, "abc")] = i
		}
	}
}

func BenchmarkInsert_CompositeKey(b *testing.B) {
	for b.Loop() {
		m := make(map[CompositeKey]int, mapSize)
		for i := range mapSize {
			m[CompositeKey{ID: i, Code: "abc"}] = i
		}
	}
}

func BenchmarkInsert_IntPairKey(b *testing.B) {
	for b.Loop() {
		m := make(map[IntPairKey]int, mapSize)
		for i := range mapSize {
			m[IntPairKey{X: i, Y: i * 7}] = i
		}
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: Lookup (key construction included)
// ---------------------------------------------------------------------------

func BenchmarkLookup_StringKey(b *testing.B) {
	m := make(map[string]int, mapSize)
	for i := range mapSize {
		m[StringKey(i, "abc")] = i
	}
	b.ResetTimer()
	for b.Loop() {
		var acc int
		for i := range mapSize {
			acc += m[StringKey(i, "abc")]
		}
		sink = acc
	}
}

func BenchmarkLookup_CompositeKey(b *testing.B) {
	m := make(map[CompositeKey]int, mapSize)
	for i := range mapSize {
		m[CompositeKey{ID: i, Code: "abc"}] = i
	}
	b.ResetTimer()
	for b.Loop() {
		var acc int
		for i := range mapSize {
			acc += m[CompositeKey{ID: i, Code: "abc"}]
		}
		sink = acc
	}
}

func BenchmarkLookup_IntPairKey(b *testing.B) {
	m := make(map[IntPairKey]int, mapSize)
	for i := range mapSize {
		m[IntPairKey{X: i, Y: i * 7}] = i
	}
	b.ResetTimer()
	for b.Loop() {
		var acc int
		for i := range mapSize {
			acc += m[IntPairKey{X: i, Y: i * 7}]
		}
		sink = acc
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: Lookup (pre-built key — isolate pure map access cost for string)
// ---------------------------------------------------------------------------
// Only string keys benefit from this separation.
// Struct keys have near-zero construction cost, so splitting them
// would just add slice-iteration noise without isolating anything.

func BenchmarkLookupPrebuilt_StringKey(b *testing.B) {
	keys := make([]string, mapSize)
	m := make(map[string]int, mapSize)
	for i := range mapSize {
		k := StringKey(i, "abc")
		keys[i] = k
		m[k] = i
	}
	b.ResetTimer()
	for b.Loop() {
		var acc int
		for _, k := range keys {
			acc += m[k]
		}
		sink = acc
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: Key construction cost only (isolated)
// ---------------------------------------------------------------------------

func BenchmarkKeyBuild_String(b *testing.B) {
	for b.Loop() {
		for i := range mapSize {
			sink = len(strconv.Itoa(i) + ":" + "abc")
		}
	}
}

func BenchmarkKeyBuild_CompositeKey(b *testing.B) {
	for b.Loop() {
		for i := range mapSize {
			k := CompositeKey{ID: i, Code: "abc"}
			sink = k.ID
		}
	}
}

func BenchmarkKeyBuild_IntPairKey(b *testing.B) {
	for b.Loop() {
		for i := range mapSize {
			k := IntPairKey{X: i, Y: i * 7}
			sink = k.X
		}
	}
}
