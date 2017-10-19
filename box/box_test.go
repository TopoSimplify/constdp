package box

import (
	"testing"
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/geom"
	"github.com/franela/goblin"
)

func TestMBRToPolygon(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("box as polygon", func() {
		g.It("should test mbr to polygon conversion", func() {
			box := mbr.NewMBR(0.25, 0.5, 5, 5)
			pts := make([]*geom.Point, 0)
			for _, pt := range box.AsPolyArray() {
				pts = append(pts, geom.NewPoint(pt))
			}

			ply := geom.NewPolygon(pts)
			g.Assert(box.Area()).Equal(ply.Area())
			g.Assert(MBRToPolygon(box).Area()).Equal(ply.Area())
		})
	})
}
