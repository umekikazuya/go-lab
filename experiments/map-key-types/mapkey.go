package mapkey

import "strconv"

// CompositeKey represents a composite map key using int + string fields.
// The compiler generates a hash function that sequentially hashes each field.
type CompositeKey struct {
	ID   int
	Code string
}

// IntPairKey represents a composite map key using two int fields.
// Fixed-size (16 bytes on amd64), enabling fast memhash.
type IntPairKey struct {
	X, Y int
}

// StringKey builds a string-based composite key by concatenation.
// This incurs heap allocation for the resulting string.
func StringKey(id int, code string) string {
	return strconv.Itoa(id) + ":" + code
}
