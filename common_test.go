package constdp

import (
	"github.com/franela/goblin"
	"github.com/intdxdt/geom"
	"testing"
)

type TestDat struct {
	pln     string
	relates ReLates
	idxs    []interface{}
	simple  string
}

type ReLates struct {
	geom bool
	dir  bool
	dist bool
}

func TestCmp(t *testing.T) {
	g := goblin.Goblin(t)
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
