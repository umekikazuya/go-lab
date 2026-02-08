package padding

// Unpadded is a struct with suboptimal field ordering that causes
// excessive padding due to alignment gaps between fields.
//
// Memory layout (amd64/arm64):
//
//	a bool    offset 0:  1 B + 7 B padding  (align next int64 to 8)
//	b int64   offset 8:  8 B
//	c bool    offset 16: 1 B + 3 B padding  (align next int32 to 4)
//	d int32   offset 20: 4 B
//	                      Total: 24 B
type Unpadded struct {
	a bool
	b int64
	c bool
	d int32
}

// Padded is the same struct with fields reordered by descending size
// to minimize padding.
//
// Memory layout (amd64/arm64):
//
//	b int64   offset 0:  8 B
//	d int32   offset 8:  4 B
//	a bool    offset 12: 1 B
//	c bool    offset 13: 1 B + 2 B padding  (struct tail alignment to 8)
//	                      Total: 16 B
type Padded struct {
	b int64
	d int32
	a bool
	c bool
}
