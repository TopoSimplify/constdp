package constdp

import (
    "time"
    "testing"
    "simplex/opts"
    "simplex/offset"
    "github.com/intdxdt/geom"
    "github.com/franela/goblin"
)

func TestConstDP_FC(t *testing.T) {
    g := goblin.Goblin(t)

    var wkts = []string{
        "POLYGON (( 435.6413255044321 1244.880520473631, 435.6413255044321 1313.5981136783437, 529.8098791553348 1313.5981136783437, 529.8098791553348 1244.880520473631, 435.6413255044321 1244.880520473631 ))",
        "POLYGON (( 700 827.4847691561165, 700 900, 763.9587152602818 900, 763.9587152602818 827.4847691561165, 700 827.4847691561165 ))",
    }
    var constraints = make([]geom.Geometry, 0)
    for _, wkt := range wkts {
        constraints = append(constraints, geom.NewGeometry(wkt))
    }

    var extractSimpleSegs = func(forest []*ConstDP) []*geom.LineString {
        var simpleLns = []*geom.LineString{}
        for _, tree := range forest {
            var coords = make([]*geom.Point, 0)
            for _, i := range tree.SimpleSet.Values() {
                coords = append(coords, tree.Pln.Coordinates[i.(int)])
            }
            simpleLns = append(simpleLns, geom.NewLineString(coords))
        }
        return simpleLns
    }

    var simplifyForest = func(lns []*geom.LineString, opts *opts.Opts) []*geom.LineString {
        var forest = []*ConstDP{}
        for _, l := range lns {
            dp := NewConstDP(l.Coordinates(), constraints, opts, offset.MaxOffset)
            forest = append(forest, dp)
        }

        SimplifyFeatureClass(forest, opts)
        return extractSimpleSegs(forest)
    }

    var simplifyInIsolation = func(lns []*geom.LineString, opts *opts.Opts) []*geom.LineString {
        forest := []*ConstDP{}
        for _, l := range lns {
            dp := NewConstDP(l.Coordinates(), constraints, opts, offset.MaxOffset)
            forest = append(forest, dp)
        }

        for _, tree := range forest {
            tree.Simplify()
        }

        return extractSimpleSegs(forest)
    }

    options := &opts.Opts{
        Threshold:              300.0,
        MinDist:                20.0,
        RelaxDist:              30.0,
        PlanarSelf:             true,
        AvoidNewSelfIntersects: true,
        GeomRelation:           true,
        DirRelation:            true,
        DistRelation:           false,
    }

    g.Describe("const dp fc", func() {
        g.It("should test constraint dp fc algorithm case 1", func() {
            g.Timeout(1 * time.Hour)
            wkts = []string{
                "LINESTRING ( 300 0, 300 400, 600 600, 600 1000, 900 1000, 900 700, 1300 700, 1400 400, 1600 200, 1300 0, 800 100, 300 0 )",
                "LINESTRING ( 100 200, 0 300, 100 500, 0 700, 400 700, 300 1100, 700 1200, 900 1300, 1100 1100, 1100 900, 1400 800, 1500 600, 1800 400, 1600 0, 1100 -200, 600 -200 )",
                "LINESTRING ( 0 100, -400 500, -300 800, 100 900, 200 1100, 200 1400, 600 1600, 900 1500, 1100 1300, 1600 1300, 1700 900, 1900 600, 1800 -200 )",
            }
            var plns = make([]*geom.LineString, 0)
            for _, wkt := range wkts {
                plns = append(plns, geom.NewLineStringFromWKT(wkt))
            }
            l0, l1, l2 := plns[0], plns[1], plns[2]

            gs := simplifyInIsolation(plns, options)
            g0, g1, g2 := gs[0], gs[1], gs[2]

            g.Assert(l0.Intersects(l1)).IsFalse()
            g.Assert(g0.Intersects(g1)).IsTrue()

            g.Assert(l1.Intersects(l2)).IsFalse()
            g.Assert(g1.Intersects(g2)).IsTrue()

            g.Assert(l0.Intersects(l2)).IsFalse()
            g.Assert(g0.Intersects(g2)).IsFalse()

            gs = simplifyForest(plns, options)
            g0, g1, g2 = gs[0], gs[1], gs[2]

            g.Assert(l0.Intersects(l1)).IsFalse()
            g.Assert(g0.Intersects(g1)).IsFalse()

            g.Assert(l1.Intersects(l2)).IsFalse()
            g.Assert(g1.Intersects(g2)).IsFalse()

            g.Assert(l0.Intersects(l2)).IsFalse()
            g.Assert(g0.Intersects(g2)).IsFalse()
        })

        g.It("should test constraint dp fc algorithm case 2", func() {
            g.Timeout(1 * time.Hour)
            wkts = []string{
                "LINESTRING ( 300 0, 300 400, 600 600, 600 1000, 900 1000, 900 700, 1300 700, 1400 400, 1600 200, 1300 0, 800 100, 300 0 )",
                "LINESTRING ( 100 200, 0 300, 100 500, 100 700, 300 800, 300 1100, 333.48668893714955 1263.8423649803672, 400 1300, 800 1100, 1100 1100, 1100 900, 1200 900, 1300 700, 1600 700, 1500 500, 1700 400, 1630.634600565117 122.53840226046754, 1600 0, 1100 -200, 600 -200 )",
                "LINESTRING ( 100 -100, -100 0, -100 100, -200 200, -200 400, -400 500, -500 400, -600 300, -500 100, -300 100, -200 400, -300 700, -200 800, -200 900, 0 800, 300 1100, 300 1300, 600 1400, 900 1500, 1100 1300, 1400 900, 1700 900, 1800 600, 1800 -200 )",
            }
            var constraints = make([]geom.Geometry, 0)
            constraints = append(constraints, geom.NewPoint([]float64{1400, 1000}))
            var plns = make([]*geom.LineString, 0)
            for _, wkt := range wkts {
                plns = append(plns, geom.NewLineStringFromWKT(wkt))
            }
            l0, l1, l2 := plns[0], plns[1], plns[2]

            l0l1 := l0.Intersection(l1)[0]
            l1l2 := l1.Intersection(l2)[0]

            g0g1 := geom.NewPointFromWKT("POINT ( 1300 700 )")
            g1g2 := geom.NewPointFromWKT("POINT ( 300 1100 )")

            g.Assert(l0l1.Distance(g0g1)).Equal(0.0)
            g.Assert(l1l2.Distance(g1g2)).Equal(0.0)

            //gs := simplify_forest(plns, options)
            var forest = []*ConstDP{}
            for _, l := range plns {
                dp := NewConstDP(l.Coordinates(), constraints, options, offset.MaxOffset)
                forest = append(forest, dp)
            }

            SimplifyFeatureClass(forest, options)
            gs := extractSimpleSegs(forest)

            g0, g1, g2 := gs[0], gs[1], gs[2]
            s0s1 := g0.Intersection(g1)[0]
            s1s2 := g1.Intersection(g2)[0]

            g.Assert(s0s1.Distance(g0g1)).Equal(0.0)
            g.Assert(s1s2.Distance(g1g2)).Equal(0.0)
        })
    })
}
