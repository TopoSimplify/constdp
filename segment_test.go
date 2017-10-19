package constdp

import (
	"testing"
	"github.com/intdxdt/geom"
	"github.com/franela/goblin"
	"simplex/rng"
	"simplex/pln"
)

func TestHullSeg(t *testing.T) {
	var g = goblin.Goblin(t)

	create_hulls := func(ranges [][]int, coords []*geom.Point) []*HullNode {
		pln := pln.New(coords)
		hulls := make([]*HullNode, 0)
		for _, r := range ranges {
			i, j := r[0], r[len(r)-1]
			h := NewHullNode(pln, rng.NewRange(i, j))
			hulls = append(hulls, h)
		}
		return hulls
	}

	g.Describe("hull decomposition", func() {
		g.It("should test decomposition of a line", func() {
			wkt := "LINESTRING ( 670 550, 680 580, 750 590, 760 630, 830 640, 870 630, 890 610, 920 580, 910 540, 890 500, 900 460, 870 420, 860 390, 810 360, 770 400, 760 420, 800 440, 810 470, 850 500, 820 560, 780 570, 760 530, 720 530, 707.3112236920351 500.3928552814154, 650 450 )"
			coords := geom.NewLineStringFromWKT(wkt).Coordinates()
			homo := &ConstDP{Pln: pln.New(coords)}
			ranges := [][]int{{0, 12}, {12, 18}, {18, len(coords) - 1}}
			hulls := create_hulls(ranges, coords)

			for i, r := range ranges {
				s := hull_segment(homo, hulls[i])
				a, b := coords[r[0]][:2], coords[r[1]][:2]
				g.Assert(r).Equal(s.Range().AsSlice())
				g.Assert(s.A[:2]).Equal(a)
				g.Assert(s.B[:2]).Equal(b)
			}
		})
	})
}
