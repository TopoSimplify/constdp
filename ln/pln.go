package ln

import (
	"simplex/geom"
	"simplex/geom/mbr"
	"simplex/constdp/seg"
	"simplex/constdp/rng"
)

//Linear interface
type Linear interface {
	Coordinates() []*geom.Point
	Polyline() *Polyline
	MaximumOffset(Linear, *rng.Range) (int, float64)
}

//Polyline Type
type Polyline struct {
	Coords []*geom.Point
	Geom   *geom.LineString
	Segs   map[[2]int]*seg.Seg
}

//construct new polyline
func NewPolyline(coords []*geom.Point) *Polyline {
	coords = geom.CloneCoordinates(coords)
	pln := &Polyline{
		Coords: coords,
		Geom:   geom.NewLineString(coords, false),
		Segs:   make(map[[2]int]*seg.Seg, 0),
	}
	pln.buildSegments()
	return pln
}

//Bounding box of polyline
func (ln *Polyline) BBox() *mbr.MBR {
	return ln.Geom.BBox()
}

//polyline
func (ln *Polyline) Polyline() *Polyline {
	return ln
}

//Coordinates at index i
func (ln *Polyline) Coordinates(i int) []*geom.Point {
	return ln.Coords
}

//Coordinates at index i
func (ln *Polyline) Coordinate(i int) *geom.Point {
	return ln.Coords[i]
}

//build polyline segments
func (ln *Polyline) buildSegments() {
	for i := 0; i < ln.Len()-1; i++ {
		j := i + 1
		ln.Segs[[2]int{i, j}] = seg.NewSeg(
			ln.Coords[i], ln.Coords[j], i, j,
		)
	}
}

//Polyline segments
func (ln *Polyline) Segments() []*seg.Seg {
	lst := make([]*seg.Seg, 0)
	for i := 0; i < ln.Len()-1; i++ {
		j := i + 1
		lst = append(lst, ln.Segs[[2]int{i, j}])
	}
	return lst
}

//Segment given range
func (ln *Polyline) Segment(rng *rng.Range) *seg.Seg {
	if rng.Size() == 1 {
		return ln.Segs[[2]int{rng.I(), rng.J()}]
	}
	return seg.NewSeg(
		ln.Coords[rng.I()],
		ln.Coords[rng.J()],
		rng.I(), rng.J(),
	)
}

//generates sub polyline from generator indices
func (self *Polyline) SubPolyline(rng *rng.Range) *Polyline {
	var poly = make([]*geom.Point, 0)
	for _, i := range rng.Stride() {
		poly = append(poly, self.Coords[i])
	}
	return NewPolyline(poly)
}

//Length of coordinates in polyline
func (ln *Polyline) Len() int {
	return len(ln.Coords)
}
