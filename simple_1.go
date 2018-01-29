package constdp

import (
	"simplex/dp"
	"simplex/node"
	"simplex/split"
	"github.com/intdxdt/fan"
)

func deformNodes(nodes map[string]*node.Node) []*node.Node {
	var stream = make(chan interface{}, 4*concurProcs)
	var exit = make(chan struct{})
	defer close(exit)

	go streamDeformNodes(stream, nodes)
	var out = fan.Stream(stream, processDeformNodes, concurProcs, exit)

	var results = make([]*node.Node, 0, len(nodes)*2)
	for sel := range out {
		splits := sel.([]*node.Node)
		results = append(results, splits...)
	}
	return results
}

func streamDeformNodes(stream chan interface{}, nodes map[string]*node.Node) {
	for _, o := range nodes {
		stream <- o
	}
	close(stream)
}

func processDeformNodes(v interface{}) interface{} {
	var hull = v.(*node.Node)
	var self = hull.Instance.(*ConstDP)
	if hull.Range.Size() > 1 {
		var ha, hb = split.AtScoreSelection(hull, self.Score, dp.NodeGeometry)
		return []*node.Node{ha, hb}
	}
	return []*node.Node{hull}
}
