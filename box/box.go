package box

import (
	"simplex/geom"
	"simplex/geom/mbr"
)

func MBRToPolygon(o *mbr.MBR) *geom.Polygon {
	coords := make([]*geom.Point, 0)
	for _, a := range o.AsPolyArray() {
		coords = append(coords, geom.NewPoint(a))
	}
	return geom.NewPolygon(coords)
}
