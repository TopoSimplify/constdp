package igeom

import "simplex/geom"

type IGeom interface {
	Geometry() geom.Geometry
}
