package constdp

import (
	"simplex/struct/deque"
	"simplex/struct/stack"
	"simplex/constdp/rng"
)

//Douglas Peucker decomposition at a given threshold
func (self *ConstDP) decompose(threshold float64) *deque.Deque {
	var pln  = self.Pln
	var hque = deque.NewDeque()
	var rg, prg *rng.Range

	rg = rng.NewRange(0, pln.Len()-1)
	stk := stack.NewStack()
	stk.Add([2]*rng.Range{rg, prg})
	// dtn = dict()

	for !stk.IsEmpty() {
		ranges := stk.Pop().([2]*rng.Range)
		rg, prg := ranges[0], ranges[1]
		k, val := self.MaximumOffset(self, rg)
		//dtn[rg] = (k, val)

		if val <= threshold {
			if prg == nil {
				prg = rg
			}
			node := NewHullNode(pln, rg, prg)
			hque.Append(node)
		} else {
			stk.Add([2]*rng.Range{rng.NewRange(k, rg.J()), rg}) // right
			stk.Add([2]*rng.Range{rng.NewRange(rg.I(), k), rg}) // left
		}
	}
	return hque
}
