package constdp

import (
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/constrain"
	"github.com/TopoSimplify/hdb"
)

//Line simplification at a given threshold
func (self *ConstDP) Simplify(constVertices ...[]int) *ConstDP {
	var constVertexSet []int
	var constVerts []int

	if len(constVertices) > 0 {
		constVerts = constVertices[0]
	}

	self.SimpleSet.Empty()
	self.Hulls = self.Decompose()

	// constrain hulls to self intersects
	self.Hulls, _, constVertexSet = constrain.ToSelfIntersects(
		self.NodeQueue(), self.Polyline(), self.Options(), constVerts,
	)
	self.selfUpdate()

	var selections map[string]*node.Node
	var db = hdb.NewHdb(rtreeBucketSize)

	var deformables = make([]*node.Node, len(self.Hulls))
	copy(deformables, self.Hulls)
	// empty deque, this is for future splits
	node.Clear(&self.Hulls)
	db.Load(deformables)

	for len(deformables) > 0 {
		// 0. find deformable node
		selections = findDeformableNodes(deformables, db)
		// 1. deform selected nodes
		deformables = deformNodes(selections)
		// 2. remove selected nodes from db
		cleanUpDB(db, selections)
		// 3. add new deformations to db
		db.Load(deformables)
		// 4. repeat until there are no deformables
	}

	self.AggregateSimpleSegments(
		db, constVertexSet, self.ScoreRelation, self.ValidateMerge,
	)

	node.Clear(&self.Hulls)
	self.SimpleSet.Empty()

	var nodes = db.All()
	for _, h := range nodes {
		self.Hulls = append(self.Hulls, h)
		self.SimpleSet.Extend(h.Range.I, h.Range.J)
	}
	return self
}
