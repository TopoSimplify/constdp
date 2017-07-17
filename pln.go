package constdp

import (
	"simplex/geom"
	"simplex/geom/mbr"
)

//Polyline Type
type Polyline struct {
	coords []*geom.Point
	geom   *geom.LineString
	segs   map[[2]int]*Seg
}

//construct new polyline
func NewPolyline(coords []*geom.Point) *Polyline {
	coords = geom.CloneCoordinates(coords)
	pln := &Polyline{
		coords: coords,
		geom:   geom.NewLineString(coords, false),
		segs:   make(map[[2]int]*Seg, 0),
	}
	pln.buildSegments()
	return pln
}

//Bounding box of polyline
func (ln *Polyline) BBox() *mbr.MBR {
	return ln.geom.BBox()
}

//Coordinates at index i
func (ln *Polyline) Coordinate(i int) *geom.Point {
	return ln.coords[i]
}

//build polyline segments
func (ln *Polyline) buildSegments() {
	for i := 0; i < ln.len()-1; i++ {
		j := i + 1
		ln.segs[[2]int{i, j}] = NewSeg(
			ln.coords[i], ln.coords[j], i, j,
		)
	}
}

//Polyline segments
func (ln *Polyline) Segments() []*Seg {
	lst := make([]*Seg, 0)
	for i := 0; i < ln.len()-1; i++ {
		j := i + 1
		lst = append(lst, ln.segs[[2]int{i, j}])
	}
	return lst
}

//Segment given range
func (ln *Polyline) Segment(rng *Range) *Seg {
	if rng.Size() == 1 {
		return ln.segs[[2]int{rng.i, rng.j}]
	}
	return NewSeg(
		ln.coords[rng.i],
		ln.coords[rng.j],
		rng.i, rng.j,
	)
}

//generates sub polyline from generator indices
func (self *Polyline) SubPolyline(rng *Range) *geom.LineString {
	var poly = make([]*geom.Point, 0)
	for _, i := range rng.Stride() {
		poly = append(poly, self.coords[i])
	}
	return geom.NewLineString(poly)
}

//Length of coordinates in polyline
func (ln *Polyline) len() int {
	return len(ln.coords)
}
