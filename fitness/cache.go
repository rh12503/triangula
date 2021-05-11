package fitness

// A CacheFunction represents a fitness function that caches data for efficiency.
type CacheFunction interface {
	Function

	SetBase(function CacheFunction)
	Cache() []TriFit
	SetCache([]TriFit)
}

// TriFit stores the triangles vertices and its fitness, and is used to cache calculations.
type TriFit struct {
	aX, aY    int16
	bX, bY    int16
	cX, cY    int16
	fitness   float64
	OtherHash uint32
}

// Equals returns if the TriFit is equal to another.
func (t TriFit) Equals(other TriFit) bool {
	return t.aX == other.aX && t.aY == other.aY &&
		t.bX == other.bX && t.bY == other.bY &&
		t.cX == other.cX && t.cY == other.cY
}

// Hash calculates the hash code of a TriFit.
func (t TriFit) Hash() uint64 {
	x := int(t.aX) + int(t.bX) + int(t.cX)
	y := int(t.aY) + int(t.bY) + int(t.cY)

	return uint64((97+x)*97 + y)
}