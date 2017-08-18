package constdp

import (
	"fmt"
	"simplex/struct/deque"
)

func debug_print_ptset(hulls *deque.Deque){
	fmt.Println(simple_hulls_as_ptset(debug_dque(hulls)))
}

func debug_dque(dq *deque.Deque) []*HullNode{
	hulls := make([]*HullNode, 0)
	for _, o := range *dq.DataView() {
		hulls = append(hulls, o.(*HullNode))
	}
	return hulls
}
