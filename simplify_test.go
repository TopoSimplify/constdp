package constdp

import (
	"fmt"
	"time"
	"bytes"
	"testing"
	"github.com/intdxdt/iter"
	"github.com/intdxdt/geom"
	"github.com/franela/goblin"
	"github.com/TopoSimplify/opts"
	"github.com/TopoSimplify/offset"
)

//@formatter:off
func cmpSlices(a, b []interface{}) bool {
	bln := len(a) == len(b)
	for i := range a {
		if !bln {
			break
		}
		bln = a[i].(int) == b[i].(int)
	}
	return bln
}

func printArray(a []interface{}) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, v := range a {
		if i < len(a)-1 {
			buf.WriteString(fmt.Sprintf("%v, ", v))
		} else {
			buf.WriteString(fmt.Sprintf("%v", v))
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func TestConstDP(t *testing.T) {
	var id = iter.NewIgen()
	var g = goblin.Goblin(t)

	g.Describe("const dp", func() {
		g.It("should test constraint dp algorithm", func() {
			g.Timeout(1 * time.Hour)
			options := &opts.Opts{
				Threshold:              50.0,
				MinDist:                20.0,
				RelaxDist:              30.0,
				PlanarSelf:             true,
				NonPlanarSelf:          false,
				AvoidNewSelfIntersects: true,
				GeomRelation:           false,
				DirRelation:            false,
				DistRelation:           false,
			}
			for i, td := range testData {
				if i < 8 {
					continue
				}
				var constraints = make([]geom.Geometry, 0)

				for _, wkt := range datConstraints {
					constraints = append(constraints, geom.NewPolygonFromWKT(wkt))
				}

				options.GeomRelation = td.relates.geom
				options.DirRelation = td.relates.dir
				options.DistRelation = td.relates.dist

				var coords = geom.NewLineStringFromWKT(td.pln).Coordinates
				var dp = NewConstDP(id.Next(),coords, constraints, options, offset.MaxOffset)
				var ptset = dp.Simplify(id).SimpleSet

				var simplx = dp.Coordinates()
				var indices = make([]int, 0, ptset.Size())
				for _, v := range ptset.Values() {
					indices = append(indices, v.(int))
				}
				simplx.Idxs = indices

				//fmt.Println(i,td.relates, td.pln)
				if !cmpSlices(ptset.Values(), td.idxs) {
					fmt.Println("debug:", i)
					fmt.Println("original:", td.idxs)
					fmt.Println("expected:", ptset.Values())
					fmt.Println("expected:", printArray(ptset.Values()))
					fmt.Println(td.pln)
					fmt.Println(td.simple)
					fmt.Println("new simple:")
					fmt.Println(geom.NewLineString(simplx).WKT())

					fmt.Println(td.relates)
				}

				g.Assert(ptset.Values()).Equal(td.idxs)
			}
		})
	})

}

func TestConstSED(t *testing.T) {
	var g = goblin.Goblin(t)
	var id = iter.NewIgen()

	g.Describe("const sed", func() {
		g.It("should test constraint sed algorithm", func() {
			g.Timeout(1 * time.Hour)
			var options = &opts.Opts{
				Threshold:              0.0,
				MinDist:                20.0,
				RelaxDist:              30.0,
				PlanarSelf:             true,
				AvoidNewSelfIntersects: true,
				GeomRelation:           true,
				DistRelation:           false,
				DirRelation:            true,
			}

			var constraints = make([]geom.Geometry, 0)
			//for _, wkt := range datConstraints {
			//	constraints = append(constraints, geom.NewPolygonFromWKT(wkt))
			//}

			var coords = geom.Coordinates([]geom.Point{
				{3.0, 1.6, 0.0}, {3.0, 2.0, 1.0}, {2.4, 2.8, 3.0}, {0.5, 3.0, 4.5},
				{1.2, 3.2, 5.0}, {1.4, 2.6, 6.0}, {2.0, 3.5, 10.0},
			})
			var inst = NewConstDP(id.Next(),coords, constraints, options, offset.MaxSEDOffset).Simplify(id)
			var ptset = make([]int, 0)
			for _, i := range inst.SimpleSet.Values() {
				ptset = append(ptset, i.(int))
			}
			g.Assert(ptset).Equal([]int{0, 1, 2, 3, 4, 5, 6})

			inst.Opts.Threshold = 1.0
			inst.Simplify(id)
			ptset = make([]int, 0)
			for _, i := range inst.SimpleSet.Values() {
				ptset = append(ptset, i.(int))
			}
			g.Assert(ptset).Equal([]int{0, 2, 3, 6})

			inst.Opts.Threshold = 1.25
			inst.Simplify(id)
			ptset = make([]int, 0)
			for _, i := range inst.SimpleSet.Values() {
				ptset = append(ptset, i.(int))
			}
			g.Assert(ptset).Equal([]int{0, 3, 6})
		})
	})
}
