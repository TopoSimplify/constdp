package constdp

import (
	"simplex/struct/deque"
	"simplex/struct/stack"
	"simplex/constdp/rng"
)

const size = 2

//Douglas Peucker decomposition at a given threshold
func (self *ConstDP) decompose(threshold float64) *deque.Deque {
	var pln   = self.Pln
	var score = self.score
	var hque  = deque.NewDeque()
	var rg, prg *rng.Range

	rg = pln.Range()
	s := stack.NewStack().Push([size]*rng.Range{rg, prg})

	for !s.IsEmpty() {
		ranges  := s.Pop().([size]*rng.Range)
		rg, prg := ranges[0], ranges[1]
		k, val  := score(self, rg)
		//dtn[rg] = (k, val)

		if val <= threshold {
			if prg == nil {
				prg = rg.Clone()
			}
			hque.Append(NewHullNode(pln, rg, prg))
		} else {
			s.Push(
				[size]*rng.Range{rng.NewRange(k, rg.J()), rg}, // right
				[size]*rng.Range{rng.NewRange(rg.I(), k), rg}, // left
			)
		}
	}
	return hque
}
