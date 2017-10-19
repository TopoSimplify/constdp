package constdp

import (
	"testing"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/rtree"
	"simplex/constdp/box"
	"github.com/franela/goblin"
)

func TestDB(t *testing.T) {
	g := goblin.Goblin(t)
	wkts := []string{
		"POINT ( 190 310 )", "POINT ( 220 400 )", "POINT ( 260 200 )", "POINT ( 260 340 )",
		"POINT ( 260 290 )", "POINT ( 310 280 )", "POINT ( 350 250 )", "POINT ( 350 330 )",
		"POINT ( 380 370 )", "POINT ( 400 240 )", "POINT ( 410 310 )",
		"POLYGON (( 160 340, 160 380, 180 380, 180 340, 160 340 ))",
		"POLYGON (( 180 240, 180 280, 210 280, 210 240, 180 240 ))",
		"POLYGON (( 280 370, 280 400, 300 400, 300 370, 280 370 ))",
		"POLYGON (( 340 210, 340 230, 360 230, 360 210, 340 210 ))",
		"POLYGON (( 410 340, 410 430, 420 430, 420 340, 410 340 ))",
	}
	g.Describe("rtree knn", func() {
		score_fn := func(q, item rtree.BoxObj) float64 {
			g := q.(geom.Geometry)
			var other geom.Geometry
			if o, ok := item.(*mbr.MBR); ok {
				other = box.MBRToPolygon(o)
			} else {
				other = item.(geom.Geometry)
			}
			return g.Distance(other)
		}
		g.It("should test k nearest neighbour", func() {
			objs := make([]rtree.BoxObj, 0)
			for _, wkt := range wkts {
				objs = append(objs, geom.NewGeometry(wkt))
			}
			tree := rtree.NewRTree(8)
			tree.Load(objs)
			q := geom.NewGeometry("POLYGON (( 370 300, 370 330, 400 330, 400 300, 370 300 ))")

			results := find_knn(tree, q, 15, score_fn)

			g.Assert(len(results) == 2)
			results = find_knn(tree, q, 20, score_fn)
			g.Assert(len(results) == 3)
		})
	})
}
