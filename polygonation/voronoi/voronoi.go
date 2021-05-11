// MIT License: See https://github.com/pzsz/voronoi/LICENSE.md

// Author: Przemyslaw Szczepaniak (przeszczep@gmail.com)
// Port of Raymond Hill's (rhill@raymondhill.net) javascript implementation 
// of Steven Forune's algorithm to compute Voronoi diagrams

package voronoi

import "math"
import "sort"
import "fmt"

type Voronoi struct {
	cells []*Cell
	edges []*Edge

	cellsMap map[Vertex]*Cell

	beachline        rbTree
	circleEvents     rbTree
	firstCircleEvent *circleEvent
}

type Diagram struct {
	Cells []*Cell
	Edges []*Edge
	//	EdgesVertices map[Vertex]EdgeVertex
}

func (s *Voronoi) getCell(site Vertex) *Cell {
	ret := s.cellsMap[site]
	if ret == nil {
		panic(fmt.Sprintf("Couldn't find cell for site %v", site))
	}
	return ret
}

func (s *Voronoi) createEdge(LeftCell, RightCell *Cell, va, vb Vertex) *Edge {
	edge := newEdge(LeftCell, RightCell)
	s.edges = append(s.edges, edge)
	if va != NO_VERTEX {
		s.setEdgeStartpoint(edge, LeftCell, RightCell, va)
	}

	if vb != NO_VERTEX {
		s.setEdgeEndpoint(edge, LeftCell, RightCell, vb)
	}

	lCell := LeftCell
	rCell := RightCell

	lCell.Halfedges = append(lCell.Halfedges, newHalfedge(edge, LeftCell, RightCell))
	rCell.Halfedges = append(rCell.Halfedges, newHalfedge(edge, RightCell, LeftCell))
	return edge
}

func (s *Voronoi) createBorderEdge(LeftCell *Cell, va, vb Vertex) *Edge {
	edge := newEdge(LeftCell, nil)
	edge.Va.Vertex = va
	edge.Vb.Vertex = vb

	s.edges = append(s.edges, edge)
	return edge
}

func (s *Voronoi) setEdgeStartpoint(edge *Edge, LeftCell, RightCell *Cell, vertex Vertex) {
	if edge.Va.Vertex == NO_VERTEX && edge.Vb.Vertex == NO_VERTEX {
		edge.Va.Vertex = vertex
		edge.LeftCell = LeftCell
		edge.RightCell = RightCell
	} else if edge.LeftCell == RightCell {
		edge.Vb.Vertex = vertex
	} else {
		edge.Va.Vertex = vertex
	}
}

func (s *Voronoi) setEdgeEndpoint(edge *Edge, LeftCell, RightCell *Cell, vertex Vertex) {
	s.setEdgeStartpoint(edge, RightCell, LeftCell, vertex)
}

type Beachsection struct {
	node        *rbNode
	site        Vertex
	circleEvent *circleEvent
	edge        *Edge
}

// rbNodeValue intergface
func (s *Beachsection) bindToNode(node *rbNode) {
	s.node = node
}

// rbNodeValue intergface
func (s *Beachsection) getNode() *rbNode {
	return s.node
}

// Calculate the left break point of a particular beach section,
// given a particular sweep line
func leftBreakPoint(arc *Beachsection, directrix float64) float64 {
	site := arc.site
	rfocx := site.X
	rfocy := site.Y
	pby2 := rfocy - directrix
	// parabola in degenerate case where focus is on directrix
	if pby2 == 0 {
		return rfocx
	}

	lArc := arc.getNode().previous
	if lArc == nil {
		return math.Inf(-1)
	}
	site = lArc.value.(*Beachsection).site
	lfocx := site.X
	lfocy := site.Y
	plby2 := lfocy - directrix
	// parabola in degenerate case where focus is on directrix
	if plby2 == 0 {
		return lfocx
	}
	hl := lfocx - rfocx
	aby2 := 1/pby2 - 1/plby2
	b := hl / plby2
	if aby2 != 0 {
		return (-b+math.Sqrt(b*b-2*aby2*(hl*hl/(-2*plby2)-lfocy+plby2/2+rfocy-pby2/2)))/aby2 + rfocx
	}
	// both parabolas have same distance to directrix, thus break point is midway
	return (rfocx + lfocx) / 2
}

// calculate the right break point of a particular beach section,
// given a particular directrix
func rightBreakPoint(arc *Beachsection, directrix float64) float64 {
	rArc := arc.getNode().next
	if rArc != nil {
		return leftBreakPoint(rArc.value.(*Beachsection), directrix)
	}
	site := arc.site
	if site.Y == directrix {
		return site.X
	}
	return math.Inf(1)
}

func (s *Voronoi) detachBeachsection(arc *Beachsection) {
	s.detachCircleEvent(arc)
	s.beachline.removeNode(arc.node)
}

type BeachsectionPtrs []*Beachsection

func (s *BeachsectionPtrs) appendLeft(b *Beachsection) {
	*s = append(*s, b)
	for id := len(*s) - 1; id > 0; id-- {
		(*s)[id] = (*s)[id-1]
	}
	(*s)[0] = b
}

func (s *BeachsectionPtrs) appendRight(b *Beachsection) {
	*s = append(*s, b)
}

func (s *Voronoi) removeBeachsection(beachsection *Beachsection) {
	circle := beachsection.circleEvent
	x := circle.x
	y := circle.ycenter
	vertex := Vertex{x, y}
	previous := beachsection.node.previous
	next := beachsection.node.next
	disappearingTransitions := BeachsectionPtrs{beachsection}
	abs_fn := math.Abs

	// remove collapsed beachsection from beachline
	s.detachBeachsection(beachsection)

	// there could be more than one empty arc at the deletion point, this
	// happens when more than two edges are linked by the same vertex,
	// so we will collect all those edges by looking up both sides of
	// the deletion point.
	// by the way, there is *always* a predecessor/successor to any collapsed
	// beach section, it's just impossible to have a collapsing first/last
	// beach sections on the beachline, since they obviously are unconstrained
	// on their left/right side.

	// look left
	lArc := previous.value.(*Beachsection)
	for lArc.circleEvent != nil &&
		abs_fn(x-lArc.circleEvent.x) < 1e-9 &&
		abs_fn(y-lArc.circleEvent.ycenter) < 1e-9 {

		previous = lArc.node.previous
		disappearingTransitions.appendLeft(lArc)
		s.detachBeachsection(lArc) // mark for reuse
		lArc = previous.value.(*Beachsection)
	}
	// even though it is not disappearing, I will also add the beach section
	// immediately to the left of the left-most collapsed beach section, for
	// convenience, since we need to refer to it later as this beach section
	// is the 'left' site of an edge for which a start point is set.
	disappearingTransitions.appendLeft(lArc)
	s.detachCircleEvent(lArc)

	// look right
	var rArc = next.value.(*Beachsection)
	for rArc.circleEvent != nil &&
		abs_fn(x-rArc.circleEvent.x) < 1e-9 &&
		abs_fn(y-rArc.circleEvent.ycenter) < 1e-9 {
		next = rArc.node.next
		disappearingTransitions.appendRight(rArc)
		s.detachBeachsection(rArc) // mark for reuse
		rArc = next.value.(*Beachsection)
	}
	// we also have to add the beach section immediately to the right of the
	// right-most collapsed beach section, since there is also a disappearing
	// transition representing an edge's start point on its left.
	disappearingTransitions.appendRight(rArc)
	s.detachCircleEvent(rArc)

	// walk through all the disappearing transitions between beach sections and
	// set the start point of their (implied) edge.
	nArcs := len(disappearingTransitions)

	for iArc := 1; iArc < nArcs; iArc++ {
		rArc = disappearingTransitions[iArc]
		lArc = disappearingTransitions[iArc-1]

		lSite := s.getCell(lArc.site)
		rSite := s.getCell(rArc.site)

		s.setEdgeStartpoint(rArc.edge, lSite, rSite, vertex)
	}

	// create a new edge as we have now a new transition between
	// two beach sections which were previously not adjacent.
	// since this edge appears as a new vertex is defined, the vertex
	// actually define an end point of the edge (relative to the site
	// on the left)
	lArc = disappearingTransitions[0]
	rArc = disappearingTransitions[nArcs-1]
	lSite := s.getCell(lArc.site)
	rSite := s.getCell(rArc.site)

	rArc.edge = s.createEdge(lSite, rSite, NO_VERTEX, vertex)

	// create circle events if any for beach sections left in the beachline
	// adjacent to collapsed sections
	s.attachCircleEvent(lArc)
	s.attachCircleEvent(rArc)
}

func (s *Voronoi) addBeachsection(site Vertex) {
	x := site.X
	directrix := site.Y

	// find the left and right beach sections which will surround the newly
	// created beach section.
	// rhill 2011-06-01: This loop is one of the most often executed,
	// hence we expand in-place the comparison-against-epsilon calls.
	var lNode, rNode *rbNode
	var dxl, dxr float64
	node := s.beachline.root

	for node != nil {
		nodeBeachline := node.value.(*Beachsection)
		dxl = leftBreakPoint(nodeBeachline, directrix) - x
		// x lessThanWithEpsilon xl => falls somewhere before the left edge of the beachsection
		if dxl > 1e-9 {
			// this case should never happen
			// if (!node.rbLeft) {
			//    rNode = node.rbLeft;
			//    break;
			//    }
			node = node.left
		} else {
			dxr = x - rightBreakPoint(nodeBeachline, directrix)
			// x greaterThanWithEpsilon xr => falls somewhere after the right edge of the beachsection
			if dxr > 1e-9 {
				if node.right == nil {
					lNode = node
					break
				}
				node = node.right
			} else {
				// x equalWithEpsilon xl => falls exactly on the left edge of the beachsection
				if dxl > -1e-9 {
					lNode = node.previous
					rNode = node
				} else if dxr > -1e-9 {
					// x equalWithEpsilon xr => falls exactly on the right edge of the beachsection
					lNode = node
					rNode = node.next
					// falls exactly somewhere in the middle of the beachsection
				} else {
					lNode = node
					rNode = node
				}
				break
			}
		}
	}

	var lArc, rArc *Beachsection

	if lNode != nil {
		lArc = lNode.value.(*Beachsection)
	}
	if rNode != nil {
		rArc = rNode.value.(*Beachsection)
	}

	// at this point, keep in mind that lArc and/or rArc could be
	// undefined or null.

	// create a new beach section object for the site and add it to RB-tree
	newArc := &Beachsection{site: site}
	if lArc == nil {
		s.beachline.insertSuccessor(nil, newArc)
	} else {
		s.beachline.insertSuccessor(lArc.node, newArc)
	}

	// cases:
	//

	// [null,null]
	// least likely case: new beach section is the first beach section on the
	// beachline.
	// This case means:
	//   no new transition appears
	//   no collapsing beach section
	//   new beachsection become root of the RB-tree
	if lArc == nil && rArc == nil {
		return
	}

	// [lArc,rArc] where lArc == rArc
	// most likely case: new beach section split an existing beach
	// section.
	// This case means:
	//   one new transition appears
	//   the left and right beach section might be collapsing as a result
	//   two new nodes added to the RB-tree
	if lArc == rArc {
		// invalidate circle event of split beach section
		s.detachCircleEvent(lArc)

		// split the beach section into two separate beach sections
		rArc = &Beachsection{site: lArc.site}
		s.beachline.insertSuccessor(newArc.node, rArc)

		// since we have a new transition between two beach sections,
		// a new edge is born
		lCell := s.getCell(lArc.site)
		newCell := s.getCell(newArc.site)
		newArc.edge = s.createEdge(lCell, newCell, NO_VERTEX, NO_VERTEX)
		rArc.edge = newArc.edge

		// check whether the left and right beach sections are collapsing
		// and if so create circle events, to be notified when the point of
		// collapse is reached.
		s.attachCircleEvent(lArc)
		s.attachCircleEvent(rArc)
		return
	}

	// [lArc,null]
	// even less likely case: new beach section is the *last* beach section
	// on the beachline -- this can happen *only* if *all* the previous beach
	// sections currently on the beachline share the same y value as
	// the new beach section.
	// This case means:
	//   one new transition appears
	//   no collapsing beach section as a result
	//   new beach section become right-most node of the RB-tree
	if lArc != nil && rArc == nil {
		lCell := s.getCell(lArc.site)
		newCell := s.getCell(newArc.site)
		newArc.edge = s.createEdge(lCell, newCell, NO_VERTEX, NO_VERTEX)
		return
	}

	// [null,rArc]
	// impossible case: because sites are strictly processed from top to bottom,
	// and left to right, which guarantees that there will always be a beach section
	// on the left -- except of course when there are no beach section at all on
	// the beach line, which case was handled above.
	// rhill 2011-06-02: No point testing in non-debug version
	//if (!lArc && rArc) {
	//    throw "Voronoi.addBeachsection(): What is this I don't even";
	//    }

	// [lArc,rArc] where lArc != rArc
	// somewhat less likely case: new beach section falls *exactly* in between two
	// existing beach sections
	// This case means:
	//   one transition disappears
	//   two new transitions appear
	//   the left and right beach section might be collapsing as a result
	//   only one new node added to the RB-tree
	if lArc != rArc {
		// invalidate circle events of left and right sites
		s.detachCircleEvent(lArc)
		s.detachCircleEvent(rArc)

		// an existing transition disappears, meaning a vertex is defined at
		// the disappearance point.
		// since the disappearance is caused by the new beachsection, the
		// vertex is at the center of the circumscribed circle of the left,
		// new and right beachsections.
		// http://mathforum.org/library/drmath/view/55002.html
		// Except that I bring the origin at A to simplify
		// calculation
		LeftSite := lArc.site
		ax := LeftSite.X
		ay := LeftSite.Y
		bx := site.X - ax
		by := site.Y - ay
		RightSite := rArc.site
		cx := RightSite.X - ax
		cy := RightSite.Y - ay
		d := 2 * (bx*cy - by*cx)
		hb := bx*bx + by*by
		hc := cx*cx + cy*cy
		vertex := Vertex{(cy*hb-by*hc)/d + ax, (bx*hc-cx*hb)/d + ay}

		lCell := s.getCell(LeftSite)
		cell := s.getCell(site)
		rCell := s.getCell(RightSite)

		// one transition disappear
		s.setEdgeStartpoint(rArc.edge, lCell, rCell, vertex)

		// two new transitions appear at the new vertex location
		newArc.edge = s.createEdge(lCell, cell, NO_VERTEX, vertex)
		rArc.edge = s.createEdge(cell, rCell, NO_VERTEX, vertex)

		// check whether the left and right beach sections are collapsing
		// and if so create circle events, to handle the point of collapse.
		s.attachCircleEvent(lArc)
		s.attachCircleEvent(rArc)
		return
	}
}

type circleEvent struct {
	node    *rbNode
	site    Vertex
	arc     *Beachsection
	x       float64
	y       float64
	ycenter float64
}

func (s *circleEvent) bindToNode(node *rbNode) {
	s.node = node
}

func (s *circleEvent) getNode() *rbNode {
	return s.node
}

func (s *Voronoi) attachCircleEvent(arc *Beachsection) {
	lArc := arc.node.previous
	rArc := arc.node.next
	if lArc == nil || rArc == nil {
		return // does that ever happen?
	}
	LeftSite := lArc.value.(*Beachsection).site
	cSite := arc.site
	RightSite := rArc.value.(*Beachsection).site

	// If site of left beachsection is same as site of
	// right beachsection, there can't be convergence
	if LeftSite == RightSite {
		return
	}

	// Find the circumscribed circle for the three sites associated
	// with the beachsection triplet.
	// rhill 2011-05-26: It is more efficient to calculate in-place
	// rather than getting the resulting circumscribed circle from an
	// object returned by calling Voronoi.circumcircle()
	// http://mathforum.org/library/drmath/view/55002.html
	// Except that I bring the origin at cSite to simplify calculations.
	// The bottom-most part of the circumcircle is our Fortune 'circle
	// event', and its center is a vertex potentially part of the final
	// Voronoi diagram.
	bx := cSite.X
	by := cSite.Y
	ax := LeftSite.X - bx
	ay := LeftSite.Y - by
	cx := RightSite.X - bx
	cy := RightSite.Y - by

	// If points l->c->r are clockwise, then center beach section does not
	// collapse, hence it can't end up as a vertex (we reuse 'd' here, which
	// sign is reverse of the orientation, hence we reverse the test.
	// http://en.wikipedia.org/wiki/Curve_orientation#Orientation_of_a_simple_polygon
	// rhill 2011-05-21: Nasty finite precision error which caused circumcircle() to
	// return infinites: 1e-12 seems to fix the problem.
	d := 2 * (ax*cy - ay*cx)
	if d >= -2e-12 {
		return
	}

	ha := ax*ax + ay*ay
	hc := cx*cx + cy*cy
	x := (cy*ha - ay*hc) / d
	y := (ax*hc - cx*ha) / d
	ycenter := y + by

	// Important: ybottom should always be under or at sweep, so no need
	// to waste CPU cycles by checking

	// recycle circle event object if possible	
	circleEventInst := &circleEvent{
		arc:     arc,
		site:    cSite,
		x:       x + bx,
		y:       ycenter + math.Sqrt(x*x+y*y),
		ycenter: ycenter,
	}

	arc.circleEvent = circleEventInst

	// find insertion point in RB-tree: circle events are ordered from
	// smallest to largest
	var predecessor *rbNode = nil
	node := s.circleEvents.root
	for node != nil {
		nodeValue := node.value.(*circleEvent)
		if circleEventInst.y < nodeValue.y || (circleEventInst.y == nodeValue.y && circleEventInst.x <= nodeValue.x) {
			if node.left != nil {
				node = node.left
			} else {
				predecessor = node.previous
				break
			}
		} else {
			if node.right != nil {
				node = node.right
			} else {
				predecessor = node
				break
			}
		}
	}
	s.circleEvents.insertSuccessor(predecessor, circleEventInst)
	if predecessor == nil {
		s.firstCircleEvent = circleEventInst
	}
}

func (s *Voronoi) detachCircleEvent(arc *Beachsection) {
	circle := arc.circleEvent
	if circle != nil {
		if circle.node.previous == nil {
			if circle.node.next != nil {
				s.firstCircleEvent = circle.node.next.value.(*circleEvent)
			} else {
				s.firstCircleEvent = nil
			}
		}
		s.circleEvents.removeNode(circle.node) // remove from RB-tree
		arc.circleEvent = nil
	}
}

// Bounding Box
type BBox struct {
	Xl, Xr, Yt, Yb float64
}

// Create new Bounding Box
func NewBBox(xl, xr, yt, yb float64) BBox {
	return BBox{xl, xr, yt, yb}
}

// connect dangling edges (not if a cursory test tells us
// it is not going to be visible.
// return value:
//   false: the dangling endpoint couldn't be connected
//   true: the dangling endpoint could be connected
func connectEdge(edge *Edge, bbox BBox) bool {
	// skip if end point already connected
	vb := edge.Vb.Vertex
	if vb != NO_VERTEX {
		return true
	}

	// make local copy for performance purpose
	va := edge.Va.Vertex
	xl := bbox.Xl
	xr := bbox.Xr
	yt := bbox.Yt
	yb := bbox.Yb
	LeftSite := edge.LeftCell.Site
	RightSite := edge.RightCell.Site
	lx := LeftSite.X
	ly := LeftSite.Y
	rx := RightSite.X
	ry := RightSite.Y
	fx := (lx + rx) / 2
	fy := (ly + ry) / 2

	var fm, fb float64

	// get the line equation of the bisector if line is not vertical
	if !equalWithEpsilon(ry, ly) {
		fm = (lx - rx) / (ry - ly)
		fb = fy - fm*fx
	}

	// remember, direction of line (relative to left site):
	// upward: left.X < right.X
	// downward: left.X > right.X
	// horizontal: left.X == right.X
	// upward: left.X < right.X
	// rightward: left.Y < right.Y
	// leftward: left.Y > right.Y
	// vertical: left.Y == right.Y

	// depending on the direction, find the best side of the
	// bounding box to use to determine a reasonable start point

	// special case: vertical line
	if equalWithEpsilon(ry, ly) {
		// doesn't intersect with viewport
		if fx < xl || fx >= xr {
			return false
		}
		// downward
		if lx > rx {
			if va == NO_VERTEX {
				va = Vertex{fx, yt}
			} else if va.Y >= yb {
				return false
			}
			vb = Vertex{fx, yb}
			// upward
		} else {
			if va == NO_VERTEX {
				va = Vertex{fx, yb}
			} else if va.Y < yt {
				return false
			}
			vb = Vertex{fx, yt}
		}
		// closer to vertical than horizontal, connect start point to the
		// top or bottom side of the bounding box
	} else if fm < -1 || fm > 1 {
		// downward
		if lx > rx {
			if va == NO_VERTEX {
				va = Vertex{(yt - fb) / fm, yt}
			} else if va.Y >= yb {
				return false
			}
			vb = Vertex{(yb - fb) / fm, yb}
			// upward
		} else {
			if va == NO_VERTEX {
				va = Vertex{(yb - fb) / fm, yb}
			} else if va.Y < yt {
				return false
			}
			vb = Vertex{(yt - fb) / fm, yt}
		}
		// closer to horizontal than vertical, connect start point to the
		// left or right side of the bounding box
	} else {
		// rightward
		if ly < ry {
			if va == NO_VERTEX {
				va = Vertex{xl, fm*xl + fb}
			} else if va.X >= xr {
				return false
			}
			vb = Vertex{xr, fm*xr + fb}
			// leftward
		} else {
			if va == NO_VERTEX {
				va = Vertex{xr, fm*xr + fb}
			} else if va.X < xl {
				return false
			}
			vb = Vertex{xl, fm*xl + fb}
		}
	}
	edge.Va.Vertex = va
	edge.Vb.Vertex = vb
	return true
}

// line-clipping code taken from:
//   Liang-Barsky function by Daniel White
//   http://www.skytopia.com/project/articles/compsci/clipping.html
// Thanks!
// A bit modified to minimize code paths
func clipEdge(edge *Edge, bbox BBox) bool {
	ax := edge.Va.X
	ay := edge.Va.Y
	bx := edge.Vb.X
	by := edge.Vb.Y
	t0 := float64(0)
	t1 := float64(1)
	dx := bx - ax
	dy := by - ay

	// left
	q := ax - bbox.Xl
	if dx == 0 && q < 0 {
		return false
	}
	r := -q / dx
	if dx < 0 {
		if r < t0 {
			return false
		} else if r < t1 {
			t1 = r
		}
	} else if dx > 0 {
		if r > t1 {
			return false
		} else if r > t0 {
			t0 = r
		}
	}
	// right
	q = bbox.Xr - ax
	if dx == 0 && q < 0 {
		return false
	}
	r = q / dx
	if dx < 0 {
		if r > t1 {
			return false
		} else if r > t0 {
			t0 = r
		}
	} else if dx > 0 {
		if r < t0 {
			return false
		} else if r < t1 {
			t1 = r
		}
	}

	// top
	q = ay - bbox.Yt
	if dy == 0 && q < 0 {
		return false
	}
	r = -q / dy
	if dy < 0 {
		if r < t0 {
			return false
		} else if r < t1 {
			t1 = r
		}
	} else if dy > 0 {
		if r > t1 {
			return false
		} else if r > t0 {
			t0 = r
		}
	}
	// bottom        
	q = bbox.Yb - ay
	if dy == 0 && q < 0 {
		return false
	}
	r = q / dy
	if dy < 0 {
		if r > t1 {
			return false
		} else if r > t0 {
			t0 = r
		}
	} else if dy > 0 {
		if r < t0 {
			return false
		} else if r < t1 {
			t1 = r
		}
	}

	// if we reach this point, Voronoi edge is within bbox

	// if t0 > 0, va needs to change
	// rhill 2011-06-03: we need to create a new vertex rather
	// than modifying the existing one, since the existing
	// one is likely shared with at least another edge
	if t0 > 0 {
		edge.Va.Vertex = Vertex{ax + t0*dx, ay + t0*dy}
	}

	// if t1 < 1, vb needs to change
	// rhill 2011-06-03: we need to create a new vertex rather
	// than modifying the existing one, since the existing
	// one is likely shared with at least another edge
	if t1 < 1 {
		edge.Vb.Vertex = Vertex{ax + t1*dx, ay + t1*dy}
	}

	return true
}

func equalWithEpsilon(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func lessThanWithEpsilon(a, b float64) bool {
	return b-a > 1e-9
}

func greaterThanWithEpsilon(a, b float64) bool {
	return a-b > 1e-9
}

// Connect/cut edges at bounding box
func (s *Voronoi) clipEdges(bbox BBox) {
	// connect all dangling edges to bounding box
	// or get rid of them if it can't be done
	abs_fn := math.Abs

	// iterate backward so we can splice safely
	for i := len(s.edges) - 1; i >= 0; i-- {
		edge := s.edges[i]
		// edge is removed if:
		//   it is wholly outside the bounding box
		//   it is actually a point rather than a line
		if !connectEdge(edge, bbox) || !clipEdge(edge, bbox) || (abs_fn(edge.Va.X-edge.Vb.X) < 1e-9 && abs_fn(edge.Va.Y-edge.Vb.Y) < 1e-9) {
			edge.Va.Vertex = NO_VERTEX
			edge.Vb.Vertex = NO_VERTEX
			s.edges[i] = s.edges[len(s.edges)-1]
			s.edges = s.edges[0 : len(s.edges)-1]
		}
	}
}

func (s *Voronoi) closeCells(bbox BBox) {
	// prune, order halfedges, then add missing ones
	// required to close cells
	xl := bbox.Xl
	xr := bbox.Xr
	yt := bbox.Yt
	yb := bbox.Yb
	cells := s.cells
	abs_fn := math.Abs

	for _, cell := range cells {
		// trim non fully-defined halfedges and sort them counterclockwise
		if cell.prepare() == 0 {
			continue
		}

		// close open cells
		// step 1: find first 'unclosed' point, if any.
		// an 'unclosed' point will be the end point of a halfedge which
		// does not match the start point of the following halfedge
		halfedges := cell.Halfedges
		nHalfedges := len(halfedges)

		// special case: only one site, in which case, the viewport is the cell
		// ...
		// all other cases
		iLeft := 0
		for iLeft < nHalfedges {
			iRight := (iLeft + 1) % nHalfedges
			endpoint := halfedges[iLeft].GetEndpoint()
			startpoint := halfedges[iRight].GetStartpoint()
			// if end point is not equal to start point, we need to add the missing
			// halfedge(s) to close the cell
			if abs_fn(endpoint.X-startpoint.X) >= 1e-9 || abs_fn(endpoint.Y-startpoint.Y) >= 1e-9 {
				// if we reach this point, cell needs to be closed by walking
				// counterclockwise along the bounding box until it connects
				// to next halfedge in the list
				va := endpoint
				vb := endpoint
				// walk downward along left side
				if equalWithEpsilon(endpoint.X, xl) && lessThanWithEpsilon(endpoint.Y, yb) {
					if equalWithEpsilon(startpoint.X, xl) {
						vb = Vertex{xl, startpoint.Y}
					} else {
						vb = Vertex{xl, yb}
					}

					// walk rightward along bottom side
				} else if equalWithEpsilon(endpoint.Y, yb) && lessThanWithEpsilon(endpoint.X, xr) {
					if equalWithEpsilon(startpoint.Y, yb) {
						vb = Vertex{startpoint.X, yb}
					} else {
						vb = Vertex{xr, yb}
					}
					// walk upward along right side
				} else if equalWithEpsilon(endpoint.X, xr) && greaterThanWithEpsilon(endpoint.Y, yt) {
					if equalWithEpsilon(startpoint.X, xr) {
						vb = Vertex{xr, startpoint.Y}
					} else {
						vb = Vertex{xr, yt}
					}
					// walk leftward along top side
				} else if equalWithEpsilon(endpoint.Y, yt) && greaterThanWithEpsilon(endpoint.X, xl) {
					if equalWithEpsilon(startpoint.Y, yt) {
						vb = Vertex{startpoint.X, yt}
					} else {
						vb = Vertex{xl, yt}
					}
				} else {
					//			break
				}

				// Create new border edge. Slide it into iLeft+1 position
				edge := s.createBorderEdge(cell, va, vb)
				cell.Halfedges = append(cell.Halfedges, nil)
				halfedges = cell.Halfedges
				nHalfedges = len(halfedges)

				copy(halfedges[iLeft+2:len(halfedges)], halfedges[iLeft+1:len(halfedges)-1])
				halfedges[iLeft+1] = newHalfedge(edge, cell, nil)

			}
			iLeft++
		}
	}
}

func (s *Voronoi) gatherVertexEdges() {
	vertexEdgeMap := make(map[Vertex][]*Edge)

	for _, edge := range s.edges {
		vertexEdgeMap[edge.Va.Vertex] = append(
			vertexEdgeMap[edge.Va.Vertex], edge)
		vertexEdgeMap[edge.Vb.Vertex] = append(
			vertexEdgeMap[edge.Vb.Vertex], edge)
	}

	for vertex, edgeSlice := range vertexEdgeMap {
		for _, edge := range edgeSlice {
			if vertex == edge.Va.Vertex {
				edge.Va.Edges = edgeSlice
			}
			if vertex == edge.Vb.Vertex {
				edge.Vb.Edges = edgeSlice
			}
		}
	}
}

// Compute voronoi diagram. If closeCells == true, edges from bounding box will be 
// included in diagram.
func ComputeDiagram(sites []Vertex, bbox BBox, closeCells bool) *Diagram {
	s := &Voronoi{
		cellsMap: make(map[Vertex]*Cell),
	}

	// Initialize site event queue
	sort.Sort(VerticesByY{sites})

	pop := func() *Vertex {
		if len(sites) == 0 {
			return nil
		}

		site := sites[0]
		sites = sites[1:]
		return &site
	}

	site := pop()

	// process queue
	xsitex := math.SmallestNonzeroFloat64
	xsitey := math.SmallestNonzeroFloat64
	var circle *circleEvent

	// main loop
	for {
		// we need to figure whether we handle a site or circle event
		// for this we find out if there is a site event and it is
		// 'earlier' than the circle event
		circle = s.firstCircleEvent

		// add beach section
		if site != nil && (circle == nil || site.Y < circle.y || (site.Y == circle.y && site.X < circle.x)) {
			// only if site is not a duplicate
			if site.X != xsitex || site.Y != xsitey {
				// first create cell for new site
				nCell := newCell(*site)
				s.cells = append(s.cells, nCell)
				s.cellsMap[*site] = nCell
				// then create a beachsection for that site
				s.addBeachsection(*site)
				// remember last site coords to detect duplicate
				xsitey = site.Y
				xsitex = site.X
			}
			site = pop()
			// remove beach section
		} else if circle != nil {
			s.removeBeachsection(circle.arc)
			// all done, quit
		} else {
			break
		}
	}

	// wrapping-up:
	//   connect dangling edges to bounding box
	//   cut edges as per bounding box
	//   discard edges completely outside bounding box
	//   discard edges which are point-like
	s.clipEdges(bbox)

	//   add missing edges in order to close opened cells
	if closeCells {
		s.closeCells(bbox)
	} else {
		for _, cell := range s.cells {
			cell.prepare()
		}
	}

	s.gatherVertexEdges()

	result := &Diagram{
		Edges: s.edges,
		Cells: s.cells,
	}
	return result
}
