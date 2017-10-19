package igeom

import "github.com/intdxdt/geom"

type IGeom interface {
	Geometry() geom.Geometry
}
