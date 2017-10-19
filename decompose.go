package constdp

import (
	"github.com/intdxdt/deque"
	"github.com/intdxdt/stack"
	"simplex/constdp/rng"
)


//Douglas Peucker decomposition at a given threshold
func (self *ConstDP) decompose() *deque.Deque {
	var pln  = self.Pln
	var hque = deque.NewDeque()
	var k int
	var val float64
	var rg   = pln.Range()

	s := stack.NewStack().Push(rg)

	for !s.IsEmpty() {
		rg = s.Pop().(*rng.Range)
		k, val  = self.Score(self, rg)
		if self.is_score_relate_valid(val) {
			hque.Append(NewHullNode(pln, rg))
		} else {
			s.Push(
				rng.NewRange(k, rg.J()), // right
				rng.NewRange(rg.I(), k), // left
			)
		}
	}
	return hque
}
