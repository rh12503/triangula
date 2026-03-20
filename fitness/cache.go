package fitness

// triangleKey is a normalized representation of a triangle's vertices, used as a map key
// for delta-based fitness evaluation. Vertices are sorted lexicographically (by X then Y)
// so the same geometric triangle always maps to the same key.
type triangleKey struct {
	aX, aY, bX, bY, cX, cY int16
}

// makeTriangleKey creates a triangleKey with vertices sorted lexicographically.
func makeTriangleKey(ax, ay, bx, by, cx, cy int16) triangleKey {
	// Sort three vertices by X, then Y
	if ax > bx || (ax == bx && ay > by) {
		ax, ay, bx, by = bx, by, ax, ay
	}
	if bx > cx || (bx == cx && by > cy) {
		bx, by, cx, cy = cx, cy, bx, by
	}
	if ax > bx || (ax == bx && ay > by) {
		ax, ay, bx, by = bx, by, ax, ay
	}
	return triangleKey{aX: ax, aY: ay, bX: bx, bY: by, cX: cx, cY: cy}
}

// triangleVarData stores per-triangle variance and area for delta-based fitness.
type triangleVarData struct {
	variance float64
	area     float64
}

// A CacheFunction represents a fitness function that caches data for efficiency.
type CacheFunction interface {
	Function

	SetBase(function CacheFunction)
	Cache() []CacheData
	SetCache([]CacheData)
}

type CacheData interface {
	Equals(data CacheData) bool
	Hash() uint64
	Data() float64
	CachedHash() uint32
	SetCachedHash(uint322 uint32)
}

// TriangleCacheData stores the triangles vertices and its fitness, and is used to cache calculations.
type TriangleCacheData struct {
	aX, aY  int16
	bX, bY  int16
	cX, cY  int16
	fitness float64
	hash    uint32
}

func (t TriangleCacheData) Data() float64 {
	return t.fitness
}

// Equals returns if the TriangleCacheData is equal to another.
func (t TriangleCacheData) Equals(other CacheData) bool {
	tri := other.(*TriangleCacheData)
	return t.aX == tri.aX && t.aY == tri.aY &&
		t.bX == tri.bX && t.bY == tri.bY &&
		t.cX == tri.cX && t.cY == tri.cY
}

// Hash calculates the hash code of a TriangleCacheData.
// Uses a polynomial hash with distinct prime multipliers per coordinate
// to reduce collisions compared to the previous commutative sum-based hash.
func (t TriangleCacheData) Hash() uint64 {
	h := uint64(17)
	h = h*31 + uint64(uint16(t.aX))
	h = h*37 + uint64(uint16(t.aY))
	h = h*41 + uint64(uint16(t.bX))
	h = h*43 + uint64(uint16(t.bY))
	h = h*47 + uint64(uint16(t.cX))
	h = h*53 + uint64(uint16(t.cY))
	return h
}

func (t TriangleCacheData) CachedHash() uint32 {
	return t.hash
}

func (t *TriangleCacheData) SetCachedHash(hash uint32) {
	t.hash = hash
}

type PolygonCacheData struct {
	coords  []int16
	fitness float64
	hash    uint32
}

func (p PolygonCacheData) CachedHash() uint32 {
	return p.hash
}

func (p *PolygonCacheData) SetCachedHash(hash uint32) {
	p.hash = hash
}

func (p PolygonCacheData) Data() float64 {
	return p.fitness
}

// Equals returns if the TriangleCacheData is equal to another.
func (p PolygonCacheData) Equals(other CacheData) bool {

	poly := other.(*PolygonCacheData)
	if len(poly.coords) != len(p.coords) {
		return false
	}
	for i, v := range poly.coords {
		if v != p.coords[i] {
			return false
		}
	}
	return true
}

// Hash calculates the hash code of a TriangleCacheData.
func (p PolygonCacheData) Hash() uint64 {

	hash := uint64(1)

	for i := 0; i < len(p.coords); i++ {
		hash = hash*97 + uint64(p.coords[i])
	}

	return hash
}
