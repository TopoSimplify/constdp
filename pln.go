package constdp

import (
	"simplex/geom"
	"simplex/geom/mbr"
)

type Polyline struct {
	coords []*geom.Point
	geom   *geom.LineString
	segs   map[[2]int]*geom.Segment
}

func NewPolyline(coords []*geom.Point) *Polyline {
	coords = geom.CloneCoordinates(coords)
	return &Polyline{
		coords: coords,
		geom:   geom.NewLineString(coords, false),
	}
}

func (ln *Polyline) BBox() *mbr.MBR {
	return ln.geom.BBox()
}

func (ln *Polyline) Coordinate(i int) *geom.Point {
	return ln.coords[i]
}

func (ln *Polyline) Segments(rng *Range) {
	for idx := 0; idx < ln.len()-1; idx++ {
		i, j := idx, idx+1
		ln.segs[[2]int{i, j}] = geom.NewSegment(
			ln.coords[i], ln.coords[j],
		)
	}
}

func (ln *Polyline) Segment(rng *Range) *geom.Segment {
	return geom.NewSegment(
		ln.coords[rng.i],
		ln.coords[rng.j],
	)
}

func (ln *Polyline) len() int {
	return len(ln.coords)
}
