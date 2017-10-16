package constdp

import (
	"fmt"
	"simplex/struct/deque"
)

func debug_print_ptset(hulls *deque.Deque) {
	fmt.Println(debug_dque(hulls).AsPointSet())
}

func debug_print_hulls(dq *deque.Deque) {
	for _, h := range debug_dque(dq).list {
		fmt.Println(h.Geom.WKT())
	}
}

func debug_dque(dq *deque.Deque) *HullNodes {
	hulls := NewHullNodes()
	for _, o := range *dq.DataView() {
		hulls.Push(o.(*HullNode))
	}
	return hulls
}
