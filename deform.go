package constdp

import (
	"simplex/dp"
	"simplex/node"
	"simplex/split"
	"simplex/constrain"
	"github.com/intdxdt/fan"
	"github.com/intdxdt/rtree"
)

func findDeformableNodes(hulls []*node.Node, hulldb *rtree.RTree) map[string]*node.Node {
	var stream = make(chan interface{}, concurProcs)
	var exit = make(chan struct{})
	defer close(exit)

	go func() {
		for _, n := range hulls {
			stream <- n
		}
		close(stream)
	}()

	var worker = func(v interface{}) interface{} {
		var selections = make([]*node.Node, 0)
		var hull = v.(*node.Node)
		var self = hull.Instance.(*ConstDP)

		// find hull neighbours
		// self intersection constraint
		// can self intersect with itself but not with other lines
		constrain.ByFeatureClassIntersection(self.Options(), hull, hulldb, &selections)

		// context_geom geometry constraint
		self.ValidateContextRelation(hull, &selections)

		return selections
	}

	var out = fan.Stream(stream, worker, concurProcs, exit)
	var results = make(map[string]*node.Node)
	for sel := range out {
		selections := sel.([]*node.Node)
		for _, n := range selections {
			results[n.Id()] = n
		}
	}
	return results
}

func deformNodes(nodes map[string]*node.Node) []*node.Node {
	var stream = make(chan interface{}, 4*concurProcs)
	var exit = make(chan struct{})
	defer close(exit)

	go func() {
		for _, o := range nodes {
			stream <- o
		}
		close(stream)
	}()

	var worker = func(v interface{}) interface{} {
		var hull = v.(*node.Node)
		var self = hull.Instance.(*ConstDP)
		if hull.Range.Size() > 1 {
			var ha, hb = split.AtScoreSelection(hull, self.Score, dp.NodeGeometry)
			return []*node.Node{ha, hb}
		}
		//fmt.Println(hull.Geom.WKT())
		return []*node.Node{hull}
	}

	var out = fan.Stream(stream, worker, concurProcs, exit)
	var results = make([]*node.Node, 0, len(nodes)*2)
	for sel := range out {
		splits := sel.([]*node.Node)
		results = append(results, splits...)
	}
	return results
}
