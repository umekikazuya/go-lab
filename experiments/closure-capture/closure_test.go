package closurecapture

import "testing"

// sink prevents the compiler from eliminating benchmark results via dead-code elimination.
var sink int

// TestAllocations measures heap allocations per call using testing.AllocsPerRun.
//
// Assertion policy:
//   - Hard assertion (t.Errorf): patterns where the compiler outcome is determined
//     by the Go spec or well-established escape analysis rules.
//   - Observation only (t.Logf): exploratory patterns whose allocation count is
//     the primary research question. These must NOT be constrained by a prior hypothesis.
func TestAllocations(t *testing.T) {

	// ---- Section 0: Baseline ----

	t.Run("Baseline_NoCapture", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = NoCapture() })
		if got != 0 {
			t.Errorf("NoCapture: want 0 allocs, got %v", got)
		}
	})

	// ---- Section 1: 2×2 Factorial ----

	t.Run("PatternA_NonEscapingReadOnly", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = PatternA_NonEscapingReadOnly() })
		if got != 0 {
			t.Errorf("PatternA: want 0 allocs, got %v", got)
		}
	})

	// PatternB is the core research question — observe without asserting.
	t.Run("PatternB_NonEscapingMutating", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = PatternB_NonEscapingMutating() })
		t.Logf("PatternB (main experiment — reference capture, non-escaping): allocs/run = %v", got)
	})

	// PatternC: the returned closure is immediately called in the test expression.
	// With inlining enabled the compiler may fold the entire expression and eliminate
	// the heap allocation even though the function signature implies escape.
	// Observation only — this is part of the research result.
	t.Run("PatternC_EscapingReadOnly", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = PatternC_EscapingReadOnly()() })
		t.Logf("PatternC (escaping read-only, inlined call site): allocs/run = %v", got)
	})

	// PatternD: same inlining caveat as PatternC.
	t.Run("PatternD_EscapingMutating", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = PatternD_EscapingMutating()() })
		t.Logf("PatternD (escaping mutating, inlined call site): allocs/run = %v", got)
	})

	// ---- Section 2: Escape Mechanism Variants ----

	t.Run("EscapeViaGlobal", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { EscapeViaGlobal() })
		t.Logf("EscapeViaGlobal: allocs/run = %v", got)
	})

	t.Run("EscapeViaGoroutine", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = EscapeViaGoroutine() })
		t.Logf("EscapeViaGoroutine: allocs/run = %v", got)
	})

	t.Run("EscapeViaInterface", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { EscapeViaInterface() })
		t.Logf("EscapeViaInterface: allocs/run = %v", got)
	})

	// ---- Section 3: IIFE ----

	t.Run("IIFEReadOnly", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = IIFEReadOnly() })
		t.Logf("IIFEReadOnly: allocs/run = %v", got)
	})

	t.Run("IIFEMutating", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = IIFEMutating() })
		t.Logf("IIFEMutating: allocs/run = %v", got)
	})

	t.Run("IIFEMutatingEscaping", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = IIFEMutatingEscaping()() })
		t.Logf("IIFEMutatingEscaping: allocs/run = %v", got)
	})

	// ---- Section 4: Multi-Variable Capture ----

	t.Run("CaptureOneVar", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = CaptureOneVar()() })
		t.Logf("CaptureOneVar: allocs/run = %v", got)
	})

	t.Run("CaptureTwoVars", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = CaptureTwoVars()() })
		t.Logf("CaptureTwoVars: allocs/run = %v", got)
	})

	t.Run("CaptureFourVars", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = CaptureFourVars()() })
		t.Logf("CaptureFourVars: allocs/run = %v", got)
	})

	t.Run("CaptureEightVars", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = CaptureEightVars()() })
		t.Logf("CaptureEightVars: allocs/run = %v", got)
	})

	// ---- Section 5: Nested Closures ----

	t.Run("NestedNeitherEscapes", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = NestedNeitherEscapes() })
		t.Logf("NestedNeitherEscapes: allocs/run = %v", got)
	})

	t.Run("NestedInnerEscapes", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = NestedInnerEscapes()() })
		t.Logf("NestedInnerEscapes: allocs/run = %v", got)
	})

	t.Run("NestedOuterEscapes", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = NestedOuterEscapes()() })
		t.Logf("NestedOuterEscapes: allocs/run = %v", got)
	})

	// ---- Section 6: Pointer Capture ----

	t.Run("CapturePointerNonEscaping", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = CapturePointerNonEscaping() })
		t.Logf("CapturePointerNonEscaping: allocs/run = %v", got)
	})

	t.Run("CapturePointerMutatingNonEscaping", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = CapturePointerMutatingNonEscaping() })
		t.Logf("CapturePointerMutatingNonEscaping: allocs/run = %v", got)
	})

	t.Run("CapturePointerEscaping", func(t *testing.T) {
		got := testing.AllocsPerRun(100, func() { sink = CapturePointerEscaping()() })
		t.Logf("CapturePointerEscaping: allocs/run = %v", got)
	})
}

// ---- Section 0: Baseline ----

func BenchmarkBaseline_NoCapture(b *testing.B) {
	for b.Loop() {
		sink = NoCapture()
	}
}

// ---- Section 1: 2×2 Factorial ----

func BenchmarkPatternA_NonEscapingReadOnly(b *testing.B) {
	for b.Loop() {
		sink = PatternA_NonEscapingReadOnly()
	}
}

func BenchmarkPatternB_NonEscapingMutating(b *testing.B) {
	for b.Loop() {
		sink = PatternB_NonEscapingMutating()
	}
}

func BenchmarkPatternC_EscapingReadOnly(b *testing.B) {
	for b.Loop() {
		sink = PatternC_EscapingReadOnly()()
	}
}

func BenchmarkPatternD_EscapingMutating(b *testing.B) {
	for b.Loop() {
		sink = PatternD_EscapingMutating()()
	}
}

// ---- Section 2: Escape Mechanism Variants ----

func BenchmarkEscapeViaGlobal(b *testing.B) {
	for b.Loop() {
		EscapeViaGlobal()
	}
}

func BenchmarkEscapeViaGoroutine(b *testing.B) {
	for b.Loop() {
		sink = EscapeViaGoroutine()
	}
}

func BenchmarkEscapeViaInterface(b *testing.B) {
	for b.Loop() {
		EscapeViaInterface()
	}
}

// ---- Section 3: IIFE ----

func BenchmarkIIFEReadOnly(b *testing.B) {
	for b.Loop() {
		sink = IIFEReadOnly()
	}
}

func BenchmarkIIFEMutating(b *testing.B) {
	for b.Loop() {
		sink = IIFEMutating()
	}
}

func BenchmarkIIFEMutatingEscaping(b *testing.B) {
	for b.Loop() {
		sink = IIFEMutatingEscaping()()
	}
}

// ---- Section 4: Multi-Variable Capture ----

func BenchmarkCaptureOneVar(b *testing.B) {
	for b.Loop() {
		sink = CaptureOneVar()()
	}
}

func BenchmarkCaptureTwoVars(b *testing.B) {
	for b.Loop() {
		sink = CaptureTwoVars()()
	}
}

func BenchmarkCaptureFourVars(b *testing.B) {
	for b.Loop() {
		sink = CaptureFourVars()()
	}
}

func BenchmarkCaptureEightVars(b *testing.B) {
	for b.Loop() {
		sink = CaptureEightVars()()
	}
}

// ---- Section 5: Nested Closures ----

func BenchmarkNestedNeitherEscapes(b *testing.B) {
	for b.Loop() {
		sink = NestedNeitherEscapes()
	}
}

func BenchmarkNestedInnerEscapes(b *testing.B) {
	for b.Loop() {
		sink = NestedInnerEscapes()()
	}
}

func BenchmarkNestedOuterEscapes(b *testing.B) {
	for b.Loop() {
		sink = NestedOuterEscapes()()
	}
}

// ---- Section 6: Pointer Capture ----

func BenchmarkCapturePointerNonEscaping(b *testing.B) {
	for b.Loop() {
		sink = CapturePointerNonEscaping()
	}
}

func BenchmarkCapturePointerMutatingNonEscaping(b *testing.B) {
	for b.Loop() {
		sink = CapturePointerMutatingNonEscaping()
	}
}

func BenchmarkCapturePointerEscaping(b *testing.B) {
	for b.Loop() {
		sink = CapturePointerEscaping()()
	}
}
