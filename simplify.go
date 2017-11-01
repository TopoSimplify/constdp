package constdp

import (
	"simplex/db"
	"simplex/dp"
	"simplex/node"
	"simplex/split"
	"simplex/constrain"
	"github.com/intdxdt/sset"
)

//Homotopic simplification at a given threshold
func (self *ConstDP) Simplify(constVertices ...[]int) *ConstDP {
	var constVertexSet *sset.SSet
	var constVerts = []int{}

	if len(constVertices) > 0 {
		constVerts = constVertices[0]
	}

	self.SimpleSet.Empty()
	self.Hulls = self.Decompose()

	// constrain hulls to self intersects
	self.Hulls, _, constVertexSet = constrain.ToSelfIntersects(
		self.NodeQueue(), self.Polyline(), self.Options(),
		constVerts, self.Score, self.ScoreRelation,
	)

	var bln bool
	var hull *node.Node
	var selections = node.NewNodes()

	var hulldb = db.NewDB(RtreeBucketSize)
	for !self.Hulls.IsEmpty() {
		// assume popped hull to be valid
		bln = true

		// pop hull in queue
		hull = popLeftHull(self.Hulls)

		// insert hull into hull db
		hulldb.Insert(hull)

		// self intersection constraint
		if bln && self.Opts.AvoidNewSelfIntersects {
			bln = constrain.BySelfIntersection(self.Options(), hull, hulldb, selections)
		}

		if !selections.IsEmpty() {
			split.SplitNodesInDB(
				self.NodeQueue(), hulldb, selections,
				self.Score, dp.NodeGeometry,
			)
		}

		if !bln {
			continue
		}

		// context_geom geometry constraint
		bln = self.ValidateContextRelation(hull, selections)

		if !selections.IsEmpty() {
			split.SplitNodesInDB(
				self.NodeQueue(), hulldb, selections,
				self.Score, dp.NodeGeometry,
			)
		}
	}

	self.AggregateSimpleSegments(
		hulldb, constVertexSet,
		self.ScoreRelation, self.ValidateMerge,
	)

	self.Hulls.Clear()
	self.SimpleSet.Empty()
	for _, h := range nodesFromRtreeNodes(hulldb.All()).Sort().DataView() {
		self.Hulls.Append(h)
		self.SimpleSet.Extend(h.Range.I(), h.Range.J())
	}
	return self
}


