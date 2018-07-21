package constdp

import (
	"github.com/TopoSimplify/node"
	"github.com/intdxdt/rtree"
)

func updateDB(hulldb *rtree.RTree, nodes []*node.Node) {
	var boxes = make([]rtree.Obj, len(nodes), len(nodes))
	for i, n := range nodes {
		boxes[i] = n
	}
	hulldb.Load(boxes)
}
