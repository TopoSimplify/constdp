package constdp

import (
	"simplex/geom"
)

type Linear interface {
	Coordinates() []*geom.Point
}
