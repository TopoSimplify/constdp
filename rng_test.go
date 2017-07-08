package constdp

import (
	"github.com/franela/goblin"
	"testing"
)

func TestRange(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Range", func() {
		g.It("int range", func() {
			rng := NewRange(3, 7)
			g.Assert(rng.I()).Equal(3)
			g.Assert(rng.J()).Equal(7)
			g.Assert(rng.Size()).Equal(4)
			g.Assert(rng.Stride()).Equal([]int{3, 4, 5, 6, 7})
			g.Assert(rng.Stride(1)).Equal([]int{3, 4, 5, 6, 7})
			g.Assert(rng.Stride(2)).Equal([]int{3, 5, 7})

			g.Assert(rng.ExclusiveStride()).Equal([]int{4, 5, 6,})
			g.Assert(rng.ExclusiveStride(1)).Equal([]int{4, 5, 6,})
			g.Assert(rng.ExclusiveStride(2)).Equal([]int{4,  6,})

			g.Assert(rng.AsArray()).Equal([]int{3, 7})
			g.Assert(rng.AsArray()).Equal([]int{3, 7})
			g.Assert(rng.String()).Equal("Range(i=3, j=7)")
		})
	})
}
