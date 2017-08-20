package constdp

import (
	"fmt"
	"simplex/struct/deque"
)

func debug_print_ptset(hulls *deque.Deque) {
	fmt.Println(simple_hulls_as_ptset(debug_dque(hulls)))
}

func debug_print_hulls(dq *deque.Deque) {
	for _, h := range debug_dque(dq) {
		fmt.Println(h.Geom.WKT())
	}
}

func debug_dque(dq *deque.Deque) []*HullNode {
	hulls := make([]*HullNode, 0)
	for _, o := range *dq.DataView() {
		hulls = append(hulls, o.(*HullNode))
	}
	return hulls
}
