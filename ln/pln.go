package ln

import (
	"github.com/intdxdt/geom"
	"github.com/intdxdt/mbr"
	"simplex/constdp/seg"
	"simplex/constdp/rng"
)

//Linear interface
type Linear interface {
	Coordinates() []*geom.Point
	Polyline() *Polyline
	Score(Linear, *rng.Range) (int, float64)
}

//Polyline Type
type Polyline struct {
	Coords []*geom.Point
	Geom   *geom.LineString
}

//construct new polyline
func NewPolyline(coords []*geom.Point) *Polyline {
	var n = len(coords)
	return &Polyline{
		Coords: coords[:n:n],
		Geom:   geom.NewLineString(coords),
	}
}

//Bounding box of polyline
func (ln *Polyline) BBox() *mbr.MBR {
	return ln.Geom.BBox()
}

//polyline
func (ln *Polyline) Polyline() *Polyline {
	return ln
}

//Coordinates
func (ln *Polyline) Coordinates() []*geom.Point {
	return ln.Coords
}

//Coordinates at index i
func (ln *Polyline) Coordinate(i int) *geom.Point {
	return ln.Coords[i]
}

//Polyline segments
func (ln *Polyline) Segments() []*seg.Seg {
	var i, j int
	var lst = make([]*seg.Seg, 0)
	for i = 0; i < ln.Len()-1; i++ {
		j = i + 1
		lst = append(lst, seg.NewSeg(ln.Coords[i], ln.Coords[j], i, j))
	}
	return lst
}

//Range of entire polyline
func (ln *Polyline) Range() *rng.Range {
	return rng.NewRange(0, ln.Len()-1)
}

//Segment given range
func (ln *Polyline) Segment(rng *rng.Range) *seg.Seg {
	var i, j = rng.I(), rng.J()
	return seg.NewSeg(ln.Coords[i], ln.Coords[j], i, j)
}

//generates sub polyline from generator indices
func (self *Polyline) SubPolyline(rng *rng.Range) *Polyline {
	return NewPolyline(self.Coords[rng.I():rng.J()+1])
}

//Length of coordinates in polyline
func (ln *Polyline) Len() int {
	return len(ln.Coords)
}
