package constdp

import (
	"simplex/dp"
	"simplex/opts"
	"simplex/node"
	"simplex/merge"
	"simplex/split"
	"simplex/constrain"
	"github.com/intdxdt/sset"
	"github.com/intdxdt/rtree"
)

//Homotopic simplification at a given threshold
func (self *ConstDP) Simplify(opts *opts.Opts, const_vertices ...[]int) *ConstDP {
	var const_vertex_set *sset.SSet
	var const_verts = []int{}

	if len(const_vertices) > 0 {
		const_verts = const_vertices[0]
	}

	self.SimpleSet.Empty()
	self.Opts = opts
	self.Hulls = self.Decompose()

	//debug_print_ptset(self.Hulls)

	//debug_print_hulls(self.Hulls)
	// constrain hulls to self intersects
	self.Hulls, _, const_vertex_set = constrain.ToSelfIntersects(self, const_verts, self.ScoreRelation)
	//debug_print_ptset(self.Hulls)

	var bln bool
	var hull *node.Node
	var selections = node.NewNodes()

	var hulldb = rtree.NewRTree(8)
	for !self.Hulls.IsEmpty() {
		// assume popped hull to be valid
		bln = true

		// pop hull in queue
		hull = popLeftHull(self.Hulls)

		// insert hull into hull db
		hulldb.Insert(hull)

		// self intersection constraint
		if bln && self.Opts.AvoidNewSelfIntersects {
			bln = constrain.SelfIntersection(self, hull, hulldb, selections)
		}

		if !selections.IsEmpty() {
			split.SplitNodesInDB(self, hulldb, selections, dp.NodeGeometry)
		}

		if !bln {
			continue
		}

		// context_geom geometry constraint
		bln = constrain.ContextRelation(self, self.ContextDB, hull, selections)
		if !selections.IsEmpty() {
			split.SplitNodesInDB(self, hulldb, selections, dp.NodeGeometry)
		}
	}

	merge.SimpleSegments(self, hulldb, const_vertex_set, self.ScoreRelation, self.ValidateMerge)

	self.Hulls.Clear()
	self.SimpleSet.Empty()
	for _, h := range nodesFromRtreeNodes(hulldb.All()).Sort().DataView() {
		self.Hulls.Append(h)
		self.SimpleSet.Extend(h.Range.I(), h.Range.J())
	}
	return self
}
