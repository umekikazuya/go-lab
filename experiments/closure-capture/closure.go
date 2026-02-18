package closurecapture

// ---- Section 0: Baseline ----

// NoCapture returns a value with no closure at all.
// Absolute zero point: 0 allocs guaranteed.
func NoCapture() int {
	x := 42
	return x
}

// ---- Section 1: 2×2 Factorial ----

// PatternA_NonEscapingReadOnly creates a closure that reads a local variable
// but is called immediately and never stored. Control group: 0 allocs.
func PatternA_NonEscapingReadOnly() int {
	x := 42
	f := func() int { return x }
	return f()
}

// PatternB_NonEscapingMutating creates a closure that mutates a local variable
// via reference capture but is called locally and never returned.
// Main experiment: does reference capture alone force a heap allocation?
func PatternB_NonEscapingMutating() int {
	x := 0
	f := func() { x++ }
	f()
	return x
}

// PatternC_EscapingReadOnly returns a closure that reads a local variable.
// Control group: the closure escapes, so ≥1 alloc is expected.
func PatternC_EscapingReadOnly() func() int {
	x := 42
	return func() int { return x }
}

// PatternD_EscapingMutating returns a closure that mutates a local variable.
// Control group: escape + mutation, ≥1 alloc expected.
func PatternD_EscapingMutating() func() int {
	x := 0
	f := func() int { x++; return x }
	return f
}

// ---- Section 2: Escape Mechanism Variants ----

// globalSink and ifaceSink prevent the compiler from eliding assignments.
var globalSink func() int
var ifaceSink interface{}

// EscapeViaGlobal causes a closure to escape by assigning it to a global variable.
func EscapeViaGlobal() {
	x := 42
	globalSink = func() int { return x }
}

// EscapeViaGoroutine causes a variable to escape via goroutine capture.
// A buffered channel is used for synchronization without extra allocation.
func EscapeViaGoroutine() int {
	x := 42
	ch := make(chan int, 1)
	go func() { ch <- x }()
	return <-ch
}

// EscapeViaInterface causes a closure to escape by boxing it into an interface{}.
func EscapeViaInterface() {
	x := 42
	f := func() int { return x }
	ifaceSink = f
}

// ---- Section 3: IIFE (Immediately Invoked Function Expression) ----

// IIFEReadOnly is an immediately invoked closure that reads a local variable.
func IIFEReadOnly() int {
	x := 42
	return func() int { return x }()
}

// IIFEMutating is an immediately invoked closure that mutates a local variable.
// The funcval is never stored — extreme non-escaping case.
func IIFEMutating() int {
	x := 0
	func() { x++ }()
	return x
}

// IIFEMutatingEscaping: the IIFE itself is immediately called, but the closure
// it produces is returned and escapes.
func IIFEMutatingEscaping() func() int {
	x := 0
	return func() func() int {
		x++
		return func() int { return x }
	}()
}

// ---- Section 4: Multi-Variable Capture (all escaping) ----

// CaptureOneVar returns a closure over 1 captured variable.
func CaptureOneVar() func() int {
	a := 1
	return func() int { return a }
}

// CaptureTwoVars returns a closure over 2 captured variables.
func CaptureTwoVars() func() int {
	a, b := 1, 2
	return func() int { return a + b }
}

// CaptureFourVars returns a closure over 4 captured variables.
func CaptureFourVars() func() int {
	a, b, c, d := 1, 2, 3, 4
	return func() int { return a + b + c + d }
}

// CaptureEightVars returns a closure over 8 captured variables.
func CaptureEightVars() func() int {
	a, b, c, d, e, f, g, h := 1, 2, 3, 4, 5, 6, 7, 8
	return func() int { return a + b + c + d + e + f + g + h }
}

// ---- Section 5: Nested Closures ----

// NestedNeitherEscapes: both the outer and inner closures are used locally.
func NestedNeitherEscapes() int {
	x := 0
	f := func() {
		g := func() { x++ }
		g()
	}
	f()
	return x
}

// NestedInnerEscapes: the outer closure is called locally; the inner closure
// it creates is returned and escapes.
func NestedInnerEscapes() func() int {
	x := 0
	outer := func() func() int {
		x++
		return func() int { return x }
	}
	return outer()
}

// NestedOuterEscapes: the outer closure is returned (escapes); the inner
// closure it creates is used only locally.
func NestedOuterEscapes() func() int {
	x := 0
	return func() int {
		inner := func() { x++ }
		inner()
		return x
	}
}

// ---- Section 6: Pointer Capture ----

// CapturePointerNonEscaping: closure reads via *int; the closure does not escape.
func CapturePointerNonEscaping() int {
	x := 42
	p := &x
	f := func() int { return *p }
	return f()
}

// CapturePointerMutatingNonEscaping: closure mutates via *int; the closure does not escape.
func CapturePointerMutatingNonEscaping() int {
	x := 0
	p := &x
	f := func() { *p++ }
	f()
	return x
}

// CapturePointerEscaping: closure captures a *int and is returned (escapes).
func CapturePointerEscaping() func() int {
	x := 0
	p := &x
	return func() int { *p++; return *p }
}
