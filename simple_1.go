package constdp

import (
	"github.com/intdxdt/fan"
	"github.com/intdxdt/iter"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/split"
	"github.com/TopoSimplify/common"
)

func deformNodes(id *iter.Igen, nodes map[*node.Node]struct{}) []node.Node {
	var results = make([]node.Node, 0, len(nodes)*2)
	var self *ConstDP
	var ha, hb node.Node
	for hull := range nodes {
		self = hull.Instance.(*ConstDP)
		if hull.Range.Size() > 1 {
			ha, hb = split.AtScoreSelection(id, hull, self.Score, common.Geometry)
			results = append(results, ha, hb)
			hull.Instance.State().MarkDirty() //after split mark instance as dirty
		} else {
			results = append(results, *hull)
		}
	}
	return results
}

func deformNodes_(id *iter.Igen, nodes map[*node.Node]struct{}) []node.Node {
	var stream = make(chan interface{}, 4*ConcurProcs)
	var exit = make(chan struct{})
	defer close(exit)

	go streamDeformNodes(stream, nodes)
	var out = fan.Stream(stream, processDeformNodes(id), ConcurProcs, exit)

	var results = make([]node.Node, 0, len(nodes)*2)
	for sel := range out {
		splits := sel.([]node.Node)
		results = append(results, splits...)
	}
	return results
}

func streamDeformNodes(stream chan interface{}, nodes map[*node.Node]struct{}) {
	for o := range nodes {
		stream <- o
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
