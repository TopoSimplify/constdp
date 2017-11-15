package main

import (
	"simplex/constdp"
	"github.com/intdxdt/geom"
)

func extractSimpleSegs(forest []*constdp.ConstDP) ([]*geom.LineString, CompStats) {
	var simpleLns = []*geom.LineString{}
	var originalN, simpleN int

	for _, tree := range forest {
		var coords = make([]*geom.Point, 0)
		var simple = tree.Simple()

		originalN += len(tree.Pln.Coordinates)
		simpleN += len(simple)

		for _, i := range simple {
			coords = append(coords, tree.Pln.Coordinates[i])
		}

		simpleLns = append(simpleLns, geom.NewLineString(coords))
	}

	return simpleLns, CompStats{
		Uncompressed: float64(originalN),
		Compressed:   float64(simpleN),
	}
}
