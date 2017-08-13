package constdp

import (
	"testing"
	"simplex/geom"
	"simplex/constdp/ln"
	"github.com/franela/goblin"
	"fmt"
	"simplex/constdp/opts"
	"time"
	"simplex/constdp/ctx"
)

func TestRelate(t *testing.T) {
	var g = goblin.Goblin(t)

	g.Describe("test relate", func() {
		g.It("should test relate", func() {
			g.Timeout(1 * time.Hour)
			options := &opts.Opts{
				Threshold:              50.0,
				MinDist:                20.0,
				RelaxDist:              30.0,
				KeepSelfIntersects:     true,
				AvoidNewSelfIntersects: true,
				GeomRelation:           true,
				DistRelation:           false,
				DirRelation:            false,
			}
			wkt := "LINESTRING ( 670 550, 680 580, 750 590, 760 630, 830 640, 870 630, 890 610, 920 580, 910 540, 890 500, 900 460, 870 420, 860 390, 810 360, 770 400, 760 420, 800 440, 810 470, 850 500, 820 560, 780 570, 760 530, 720 530, 707.3112236920351 500.3928552814154, 650 450 )"
			coords := geom.NewLineStringFromWKT(wkt).Coordinates()
			cdp := &ConstDP{Pln: ln.NewPolyline(coords), Opts: options}
			ranges := [][]int{{0, 12}, {12, 18}, {18, len(coords) - 1}}

			hulls := create_hulls(ranges, coords)
			neib := geom.NewPolygonFromWKT("POLYGON ((674.7409300316725 422.8229196659948, 674.7409300316725 446.72732507918226, 691.3886409444281 446.72732507918226, 691.3886409444281 422.8229196659948, 674.7409300316725 422.8229196659948))")
			const_geom := ctx.NewCtxGeom(neib, 0, -1).AsContextNeighbour()
			for _, h := range hulls {
				fmt.Println(h.Geom.WKT())
				g.Assert(cdp.is_geom_relate_valid(h, const_geom)).IsTrue()
				g.Assert(cdp.is_dir_relate_valid(h, const_geom)).IsTrue()
				g.Assert(cdp.is_dist_relate_valid(h, const_geom)).IsTrue()
			}

			neib = geom.NewPolygonFromWKT("POLYGON ((800 614.9282601093252, 800 640, 816.138388266816 640, 816.138388266816 614.9282601093252, 800 614.9282601093252))")
			const_geom = ctx.NewCtxGeom(neib, 0, -1).AsContextNeighbour()
			g.Assert(cdp.is_geom_relate_valid(hulls[0], const_geom)).IsFalse()
			g.Assert(cdp.is_geom_relate_valid(hulls[1], const_geom)).IsTrue()
			g.Assert(cdp.is_geom_relate_valid(hulls[2], const_geom)).IsTrue()

			neib = geom.NewPolygonFromWKT("POLYGON ((749.9625484910762 464.581584548546, 749.9625484910762 486.30832777325406, 762.1390749137147 486.30832777325406, 762.1390749137147 464.581584548546, 749.9625484910762 464.581584548546))")
			const_geom = ctx.NewCtxGeom(neib, 0, -1).AsContextNeighbour()
			g.Assert(cdp.is_geom_relate_valid(hulls[0], const_geom)).IsFalse()
			g.Assert(cdp.is_geom_relate_valid(hulls[1], const_geom)).IsTrue()
			g.Assert(cdp.is_geom_relate_valid(hulls[2], const_geom)).IsFalse()

			g.Assert(cdp.is_dist_relate_valid(hulls[0], const_geom)).IsFalse()
			g.Assert(cdp.is_dist_relate_valid(hulls[1], const_geom)).IsTrue()
			g.Assert(cdp.is_dist_relate_valid(hulls[2], const_geom)).IsFalse()

			g.Assert(cdp.is_score_relate_valid(25)).IsTrue()
			g.Assert(cdp.is_score_relate_valid(50.1)).IsFalse()
		})
	})
}
