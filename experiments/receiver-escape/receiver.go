package receiver

// Sumer is an interface satisfied by value receivers.
// Used to force dynamic dispatch and observe escape behavior of value types.
type Sumer interface {
	Sum() float64
}

// PSumer is an interface satisfied only by pointer receivers.
// Used to force dynamic dispatch and observe escape behavior of pointer types.
type PSumer interface {
	PSum() float64
}

// Small is a 24-byte struct (3 x float64).
type Small struct {
	X, Y, Z float64
}

func (s Small) Sum() float64   { return s.X + s.Y + s.Z }
func (s *Small) PSum() float64 { return s.X + s.Y + s.Z }

//go:noinline
func (s Small) SumNoInline() float64 { return s.X + s.Y + s.Z }

//go:noinline
func (s *Small) PSumNoInline() float64 { return s.X + s.Y + s.Z }

// Medium is a 64-byte struct (8 x float64).
type Medium struct {
	Data [8]float64
}

func (m Medium) Sum() float64 {
	var acc float64
	for _, v := range m.Data {
		acc += v
	}
	return acc
}

func (m *Medium) PSum() float64 {
	var acc float64
	for _, v := range m.Data {
		acc += v
	}
	return acc
}

//go:noinline
func (m Medium) SumNoInline() float64 {
	var acc float64
	for _, v := range m.Data {
		acc += v
	}
	return acc
}

//go:noinline
func (m *Medium) PSumNoInline() float64 {
	var acc float64
	for _, v := range m.Data {
		acc += v
	}
	return acc
}

// Large is a 128-byte struct (16 x float64).
type Large struct {
	Data [16]float64
}

func (l Large) Sum() float64 {
	var acc float64
	for _, v := range l.Data {
		acc += v
	}
	return acc
}

func (l *Large) PSum() float64 {
	var acc float64
	for _, v := range l.Data {
		acc += v
	}
	return acc
}

//go:noinline
func (l Large) SumNoInline() float64 {
	var acc float64
	for _, v := range l.Data {
		acc += v
	}
	return acc
}

//go:noinline
func (l *Large) PSumNoInline() float64 {
	var acc float64
	for _, v := range l.Data {
		acc += v
	}
	return acc
}

// XLarge is a 256-byte struct (32 x float64).
type XLarge struct {
	Data [32]float64
}

func (x XLarge) Sum() float64 {
	var acc float64
	for _, v := range x.Data {
		acc += v
	}
	return acc
}

func (x *XLarge) PSum() float64 {
	var acc float64
	for _, v := range x.Data {
		acc += v
	}
	return acc
}

//go:noinline
func (x XLarge) SumNoInline() float64 {
	var acc float64
	for _, v := range x.Data {
		acc += v
	}
	return acc
}

//go:noinline
func (x *XLarge) PSumNoInline() float64 {
	var acc float64
	for _, v := range x.Data {
		acc += v
	}
	return acc
}
