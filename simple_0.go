package constdp

import (
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/constrain"
	"github.com/intdxdt/fan"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/rtree"
)

func findDeformableNodes(hulls []*node.Node, hulldb *rtree.RTree) map[string]*node.Node {
	var stream = make(chan interface{}, concurProcs)
	var exit = make(chan struct{})
	defer close(exit)

	go inputStreamFindDeform(stream, hulls)

	var worker = processFindDeformables(hulldb)
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

func inputStreamFindDeform(stream chan interface{}, hulls []*node.Node) {
	for _, n := range hulls {
		stream <- n
	}
	close(stream)
}

func processFindDeformables(hulldb *rtree.RTree) func(v interface{}) interface{} {
	return func(v interface{}) interface{} {
		var selections = make([]*node.Node, 0)
		var hull = v.(*node.Node)
		var self = hull.Instance.(*ConstDP)

		//if hull is segment
		if hull.Range.Size() == 1 {
			return selections
		}

		//if hull geometry is line then points are collinear
		if _, ok := hull.Geom.(*geom.LineString); ok {
			return selections
		}

		// self intersection constraint
		if self.Opts.AvoidNewSelfIntersects {
			constrain.ByFeatureClassIntersection(self.Opts, hull, hulldb, &selections)
		}

		// context_geom geometry constraint
		self.ValidateContextRelation(hull, &selections)

		return selections
	}
}
