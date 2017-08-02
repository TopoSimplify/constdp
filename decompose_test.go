package constdp

import (
	"testing"
	"simplex/geom"
	"github.com/franela/goblin"
	"simplex/constdp/opts"
	"simplex/constdp/offset"
)

func TestDecompose(t *testing.T) {
	var g = goblin.Goblin(t)

	g.Describe("hull decomposition", func() {
		g.It("should test decomposition of a line", func() {
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

			// self.relates = relations(self)
			wkt := "LINESTRING ( 470 480, 470 450, 490 430, 520 420, 540 440, 560 430, 580 420, 590 410, 630 400, 630 430, 640 460, 630 490, 630 520, 640 540, 660 560, 690 580, 700 600, 730 600, 750 570, 780 560, 790 550, 800 520, 830 500, 840 480, 850 460, 900 440, 920 440, 950 480, 990 480, 1000 520, 1000 570, 990 600, 1010 620, 1060 600 )"
			coords := geom.NewLineStringFromWKT(wkt).Coordinates()
			constraints := make([]geom.Geometry, 0)
			cdp := NewConstDP(coords, constraints, options, offset.MaxOffset)

			hulls := cdp.decompose(120)
			g.Assert(hulls.Len()).Equal(4)
			hulls = cdp.decompose(150)
			g.Assert(hulls.Len()).Equal(1)
			h := hulls.Get(0).(*HullNode)
			g.Assert(h.Range.AsSlice()).Equal([]int{0, len(coords) - 1})
			g.Assert(h.Range.AsSlice()).Equal(h.PRange.AsSlice())
		})
	})
}
