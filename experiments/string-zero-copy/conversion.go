package stringzerocopy

// sink prevents the compiler from eliminating dead code.
var sink int

// BytesToStringAssign converts []byte to string via assignment.
// The result escapes to the caller, forcing a heap allocation.
//
//go:noinline
func BytesToStringAssign(b []byte) string {
	s := string(b)
	return s
}

// StringToBytesAssign converts string to []byte via assignment.
// The result escapes to the caller, forcing a heap allocation.
//
//go:noinline
func StringToBytesAssign(s string) []byte {
	b := []byte(s)
	return b
}

// BytesToStringMapLookup performs a map lookup using string(b) as the key.
// The compiler recognizes this pattern and avoids allocating a new string.
//
//go:noinline
func BytesToStringMapLookup(m map[string]int, b []byte) int {
	return m[string(b)]
}

// StringToBytesRange iterates over []byte(s) using a range loop.
// The compiler recognizes this pattern and avoids allocating a new slice.
//
//go:noinline
func StringToBytesRange(s string) {
	for _, v := range []byte(s) {
		sink += int(v)
	}
}

// BytesToStringCompare compares string(b) against a string literal.
// The compiler recognizes this pattern and avoids allocating a new string.
//
//go:noinline
func BytesToStringCompare(b []byte, target string) bool {
	return string(b) == target
}

// BytesToStringConcat concatenates a prefix with string(b).
// A new string must be allocated for the result.
//
//go:noinline
func BytesToStringConcat(b []byte) string {
	return "prefix:" + string(b)
}