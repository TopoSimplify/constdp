package cmp

import (
	"testing"
	"github.com/franela/goblin"
	"simplex/geom"
)

func TestCmp(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("cmp int", func() {
		g.It("should test comparison", func() {
			g.Assert(IntCmp(0, 0)).Equal(0)
			g.Assert(IntCmp(0, 1)).Equal(-1)
			g.Assert(IntCmp(1, 0)).Equal(1)
			g.Assert(IntCmp(1, 3)).Equal(-2)
			g.Assert(IntCmp(3, 3)).Equal(0)
		})
	})
	g.Describe("cmp point index", func() {
		g.It("should test comparison of points by index", func() {
			a, b, c := &geom.Point{0, 0, 0}, &geom.Point{0, 0, 1}, &geom.Point{0, 0, 2}
			g.Assert(PointIndexCmp(a, b)).Equal(-1)
			g.Assert(PointIndexCmp(b, a)).Equal(1)
			g.Assert(PointIndexCmp(b, b)).Equal(0)
			g.Assert(PointIndexCmp(c, c)).Equal(0)
		})
	})
}
