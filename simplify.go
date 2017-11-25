package constdp

import (
	"simplex/node"
	"simplex/common"
	"simplex/constrain"
	"github.com/intdxdt/sset"
	"github.com/intdxdt/rtree"
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
	self.selfUpdate()

	var selections map[string]*node.Node
	var hulldb = rtree.NewRTree(rtreeBucketSize)
	var boxes = make([]rtree.BoxObj, 0)

	var deformables = make([]*node.Node, 0)
	for _, hull := range self.Hulls {
		deformables = append(deformables, hull)
		boxes = append(boxes, hull)
	}
	node.Clear(&self.Hulls) // empty deque, this is for future splits

	hulldb.Load(boxes)

	for len(deformables) > 0 {
		// 0. find deformable node
		selections = findDeformableNodes(deformables, hulldb)
		// 1. deform selected nodes
		deformables = deformNodes(selections)
		// 2. remove selected nodes from db
		cleanUpDB(hulldb, selections)
		// 3. add new deformations to db
		updateDB(hulldb, deformables)
		// 4. repeat until there are no deformables
	}

	self.AggregateSimpleSegments(
		hulldb, constVertexSet,
		self.ScoreRelation, self.ValidateMerge,
	)

	node.Clear(&self.Hulls)
	self.SimpleSet.Empty()

	for _, h := range common.NodesFromRtreeNodes(hulldb.All()) {
		self.Hulls = append(self.Hulls, h)
		self.SimpleSet.Extend(h.Range.I, h.Range.J)
	}
	return self
}
