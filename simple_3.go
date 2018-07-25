package constdp

import (
	"github.com/TopoSimplify/node"
	"github.com/intdxdt/rtree"
)

func updateDB(hulldb *rtree.RTree, nodes []*node.Node) {
	var boxes = make([]*rtree.Obj, 0, len(nodes))
	for i := range nodes {
		boxes = append(boxes, rtree.Object(i, nodes[i].Bounds(), nodes[i]))
	}
	hulldb.Load(boxes)
}
