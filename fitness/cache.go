package fitness

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
func (t TriangleCacheData) Hash() uint64 {
	x := int(t.aX) + int(t.bX) + int(t.cX)
	y := int(t.aY) + int(t.bY) + int(t.cY)

	return uint64((97+x)*97 + y)
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
