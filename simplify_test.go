package constdp

import (
	"fmt"
	"time"
	"testing"
	"github.com/intdxdt/geom"
	"simplex/opts"
	"simplex/offset"
	"github.com/franela/goblin"
)
//@formatter:off

func cmp_slices(a, b []interface{}) bool {
	bln := len(a) == len(b)
	for  i := range a {
		if !bln {break}
		bln = a[i].(int) == b[i].(int)
	}
	return bln
}

func TestConstDP(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("const dp", func() {
		g.It("should test constraint dp algorithm", func() {
			g.Timeout(1 * time.Hour)
			options := &opts.Opts{
				Threshold:              50.0,
				MinDist:                20.0,
				RelaxDist:              30.0,
				KeepSelfIntersects:     true,
				AvoidNewSelfIntersects: true,
				GeomRelation:           false,
				DirRelation:            false,
				DistRelation:           false,
			}

			for i, td := range test_data {
				var constraints = make([]geom.Geometry, 0)

				for _, wkt := range dat_constraints {
					constraints = append(constraints, geom.NewPolygonFromWKT(wkt))
				}

				options.GeomRelation = td.relates.geom
				options.DirRelation  = td.relates.dir
				options.DistRelation = td.relates.dist

				var coords = geom.NewLineStringFromWKT(td.pln).Coordinates()
				var dp = NewConstDP(coords, constraints, options, offset.MaxOffset)
				var ptset = dp.Simplify(options).SimpleSet

				simplx := make([]*geom.Point, 0)
				for _, i := range ptset.Values() {
					simplx = append(simplx, coords[i.(int)])
				}
				//fmt.Println(i,td.relates, td.pln)
				if !cmp_slices(ptset.Values(), td.idxs){
					fmt.Println("debug:", i)
					fmt.Println(ptset.Values())
					fmt.Println(td.idxs)
					fmt.Println(td.pln)
					fmt.Println(td.relates)
				}
				g.Assert(ptset.Values()).Equal(td.idxs)
			}
		})
	})

}

func TestConstSED(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("const sed", func() {
		g.It("should test constraint sed algorithm", func() {

			options := &opts.Opts{
				Threshold:              1.0,
				MinDist:                20.0,
				RelaxDist:              30.0,
				KeepSelfIntersects:     true,
				AvoidNewSelfIntersects: true,
				GeomRelation:           true,
				DistRelation:           false,
				DirRelation:            true,
			}

			constraints := make([]geom.Geometry, 0)
			for _, wkt := range dat_constraints {
				constraints = append(constraints, geom.NewPolygonFromWKT(wkt))
			}

			coords := []*geom.Point{
				{3.0, 1.6, 0.0}, {3.0, 2.0, 1.0}, {2.4, 2.8, 3.0},
				{0.5, 3.0, 4.5}, {1.2, 3.2, 5.0}, {1.4, 2.6, 6.0},
				{2.0, 3.5, 10.0}}
			homo := NewConstDP(coords, constraints, options, offset.MaxSEDOffset)
			homo.Simplify(options)
		})
	})
}
