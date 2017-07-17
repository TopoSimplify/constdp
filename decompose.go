package constdp

import (
	"simplex/struct/deque"
	"simplex/struct/stack"
)

//Douglas Peucker decomposition at a given threshold
func (self *ConstDP) dp_decompose(threshold float64) *deque.Deque {
	pln := self.Pln
	hque := deque.NewDeque()
	var rng, prng *Range
	rng = NewRange(0, pln.len()-1)
	stk := stack.NewStack()
	stk.Add([2]*Range{rng, prng})
	// dtn = dict()

	for !stk.IsEmpty() {
		ranges := stk.Pop().([2]*Range)
		rng, prng := ranges[0], ranges[1]
		k, val := self.MaximumOffset(self, rng)
		//dtn[rng] = (k, val)

		if val <= threshold {
			if prng == nil {
				prng = rng
			}
			node := NewHullNode(pln, rng, prng)
			hque.Append(node)
		} else {
			stk.Add([2]*Range{NewRange(k, rng.j), rng}) // right
			stk.Add([2]*Range{NewRange(rng.i, k), rng}) // left
		}
	}
	return hque
}
