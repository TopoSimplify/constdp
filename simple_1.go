package constdp

import (
	"github.com/intdxdt/fan"
	"github.com/intdxdt/iter"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/split"
	"github.com/TopoSimplify/common"
)

func deformNodes(id *iter.Igen, nodes map[int]*node.Node) []node.Node {
	var stream = make(chan interface{}, 4*concurProcs)
	var exit = make(chan struct{})
	defer close(exit)

	go streamDeformNodes(stream, nodes)
	var out = fan.Stream(stream, processDeformNodes(id), concurProcs, exit)

	var results = make([]node.Node, 0, len(nodes)*2)
	for sel := range out {
		splits := sel.([]node.Node)
		results = append(results, splits...)
	}
	return results
}

func streamDeformNodes(stream chan interface{}, nodes map[int]*node.Node) {
	for i := range nodes {
		stream <- nodes[i]
	}
	close(stream)
}

func processDeformNodes(id *iter.Igen) func(v interface{}) interface{} {
	return func(v interface{}) interface{} {
		var hull = v.(*node.Node)
		var self = hull.Instance.(*ConstDP)
		if hull.Range.Size() > 1 {
			var ha, hb = split.AtScoreSelection(id, hull, self.Score, common.Geometry)
			return []node.Node{ha, hb}
		}
		return []node.Node{*hull}
	}
}
