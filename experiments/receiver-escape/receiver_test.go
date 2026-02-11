package receiver

import (
	"testing"
	"unsafe"
)

// TestStructSizes verifies expected struct sizes.
func TestStructSizes(t *testing.T) {
	tests := []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{"Small", unsafe.Sizeof(Small{}), 24},
		{"Medium", unsafe.Sizeof(Medium{}), 64},
		{"Large", unsafe.Sizeof(Large{}), 128},
		{"XLarge", unsafe.Sizeof(XLarge{}), 256},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("unsafe.Sizeof(%s{}) = %d B, want %d B", tt.name, tt.got, tt.want)
			}
			t.Logf("unsafe.Sizeof(%s{}) = %d B", tt.name, tt.got)
		})
	}
}

// sink prevents dead-code elimination.
var sink float64

// ---------------------------------------------------------------------------
// Small (24 B)
// ---------------------------------------------------------------------------

func BenchmarkSmallValue(b *testing.B) {
	s := Small{X: 1.0, Y: 2.0, Z: 3.0}
	for b.Loop() {
		sink = s.Sum()
	}
}

func BenchmarkSmallPointer(b *testing.B) {
	s := Small{X: 1.0, Y: 2.0, Z: 3.0}
	for b.Loop() {
		sink = s.PSum()
	}
}

// ---------------------------------------------------------------------------
// Medium (64 B)
// ---------------------------------------------------------------------------

func BenchmarkMediumValue(b *testing.B) {
	m := Medium{}
	for i := range m.Data {
		m.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = m.Sum()
	}
}

func BenchmarkMediumPointer(b *testing.B) {
	m := Medium{}
	for i := range m.Data {
		m.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = m.PSum()
	}
}

// ---------------------------------------------------------------------------
// Large (128 B)
// ---------------------------------------------------------------------------

func BenchmarkLargeValue(b *testing.B) {
	l := Large{}
	for i := range l.Data {
		l.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = l.Sum()
	}
}

func BenchmarkLargePointer(b *testing.B) {
	l := Large{}
	for i := range l.Data {
		l.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = l.PSum()
	}
}

// ---------------------------------------------------------------------------
// XLarge (256 B)
// ---------------------------------------------------------------------------

func BenchmarkXLargeValue(b *testing.B) {
	x := XLarge{}
	for i := range x.Data {
		x.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = x.Sum()
	}
}

func BenchmarkXLargePointer(b *testing.B) {
	x := XLarge{}
	for i := range x.Data {
		x.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = x.PSum()
	}
}

// ===========================================================================
// NoInline variants — simulate non-inlineable methods to expose escape cost
// ===========================================================================

// ---------------------------------------------------------------------------
// Small (24 B) NoInline
// ---------------------------------------------------------------------------

func BenchmarkSmallValueNoInline(b *testing.B) {
	s := Small{X: 1.0, Y: 2.0, Z: 3.0}
	for b.Loop() {
		sink = s.SumNoInline()
	}
}

func BenchmarkSmallPointerNoInline(b *testing.B) {
	s := Small{X: 1.0, Y: 2.0, Z: 3.0}
	for b.Loop() {
		sink = s.PSumNoInline()
	}
}

// ---------------------------------------------------------------------------
// Medium (64 B) NoInline
// ---------------------------------------------------------------------------

func BenchmarkMediumValueNoInline(b *testing.B) {
	m := Medium{}
	for i := range m.Data {
		m.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = m.SumNoInline()
	}
}

func BenchmarkMediumPointerNoInline(b *testing.B) {
	m := Medium{}
	for i := range m.Data {
		m.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = m.PSumNoInline()
	}
}

// ---------------------------------------------------------------------------
// Large (128 B) NoInline
// ---------------------------------------------------------------------------

func BenchmarkLargeValueNoInline(b *testing.B) {
	l := Large{}
	for i := range l.Data {
		l.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = l.SumNoInline()
	}
}

func BenchmarkLargePointerNoInline(b *testing.B) {
	l := Large{}
	for i := range l.Data {
		l.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = l.PSumNoInline()
	}
}

// ---------------------------------------------------------------------------
// XLarge (256 B) NoInline
// ---------------------------------------------------------------------------

func BenchmarkXLargeValueNoInline(b *testing.B) {
	x := XLarge{}
	for i := range x.Data {
		x.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = x.SumNoInline()
	}
}

func BenchmarkXLargePointerNoInline(b *testing.B) {
	x := XLarge{}
	for i := range x.Data {
		x.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = x.PSumNoInline()
	}
}

// ===========================================================================
// Interface variants — dynamic dispatch forces heap escape
// ===========================================================================

// callViaInterface dispatches Sum() through the Sumer interface (value receiver path).
//
//go:noinline
func callViaInterface(s Sumer) float64 {
	return s.Sum()
}

// callViaPInterface dispatches PSum() through the PSumer interface (pointer receiver path).
//
//go:noinline
func callViaPInterface(s PSumer) float64 {
	return s.PSum()
}

// ---------------------------------------------------------------------------
// Small (24 B) Interface
// ---------------------------------------------------------------------------

func BenchmarkSmallValueIface(b *testing.B) {
	s := Small{X: 1.0, Y: 2.0, Z: 3.0}
	for b.Loop() {
		sink = callViaInterface(s)
	}
}

func BenchmarkSmallPointerIface(b *testing.B) {
	s := &Small{X: 1.0, Y: 2.0, Z: 3.0}
	for b.Loop() {
		sink = callViaPInterface(s)
	}
}

// ---------------------------------------------------------------------------
// Medium (64 B) Interface
// ---------------------------------------------------------------------------

func BenchmarkMediumValueIface(b *testing.B) {
	m := Medium{}
	for i := range m.Data {
		m.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = callViaInterface(m)
	}
}

func BenchmarkMediumPointerIface(b *testing.B) {
	m := &Medium{}
	for i := range m.Data {
		m.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = callViaPInterface(m)
	}
}

// ---------------------------------------------------------------------------
// Large (128 B) Interface
// ---------------------------------------------------------------------------

func BenchmarkLargeValueIface(b *testing.B) {
	l := Large{}
	for i := range l.Data {
		l.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = callViaInterface(l)
	}
}

func BenchmarkLargePointerIface(b *testing.B) {
	l := &Large{}
	for i := range l.Data {
		l.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = callViaPInterface(l)
	}
}

// ---------------------------------------------------------------------------
// XLarge (256 B) Interface
// ---------------------------------------------------------------------------

func BenchmarkXLargeValueIface(b *testing.B) {
	x := XLarge{}
	for i := range x.Data {
		x.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = callViaInterface(x)
	}
}

func BenchmarkXLargePointerIface(b *testing.B) {
	x := &XLarge{}
	for i := range x.Data {
		x.Data[i] = float64(i)
	}
	for b.Loop() {
		sink = callViaPInterface(x)
	}
}
