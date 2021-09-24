package constdp

import (
	"fmt"
	"github.com/TopoSimplify/geometry"
	"github.com/intdxdt/geom"
)

func parseConstraintFeatures(inputs []string) []geometry.IGeometries {
	//var points = make([]geometry.Point, 0, len(inputs))
	//var polylines = make([]geometry.Polyline, 0, len(inputs))
	//var polygons = make([]geometry.Polygon, 0, len(inputs))
	var geometries = make([]geometry.IGeometries, 0, len(inputs))

	for idx, wkt := range inputs {
		g := geom.ReadGeometry(wkt)

		if g.Type().IsPoint() {
			var pnt = createPoint(g.(geom.Point), idx)
			//points = append(points, pnt)
			geometries = append(geometries, pnt)
		}

		if g.Type().IsLineString() {
			var ln = createPolyline(idx, g.(*geom.LineString))
			geometries = append(geometries, ln)
		}

		if g.Type().IsPolygon() {
			var poly = createPolygon(idx, g.(*geom.Polygon))
			geometries = append(geometries, poly)
		}
	}

	return geometries
}

func createPoint(g geom.Point, id int) geometry.Point {
	var fid = fmt.Sprintf("%v", id)
	return geometry.Point{G: g, Id: fid}
}

func createPolyline(idx int, g *geom.LineString) geometry.Polyline {
	return geometry.Polyline{G: g, Id: composeId(idx)}
}

func createPolygon(idx int, g *geom.Polygon) geometry.Polygon {
	return geometry.Polygon{G: g, Id: composeId(idx)}
}

func composeId(index int) string {
	return fmt.Sprintf("idx:%v", index)
}
