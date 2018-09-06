package constdp

import (
	"github.com/intdxdt/iter"
	"github.com/TopoSimplify/hdb"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/constrain"
)

//Line simplification at a given threshold
func (self *ConstDP) Simplify(id *iter.Igen, constVertices ...[]int) *ConstDP {
	var constVertexSet []int
	var constVerts []int

	if len(constVertices) > 0 {
		constVerts = constVertices[0]
	}

	//state
	self.State().MarkDirty()

	self.SimpleSet.Empty()
	self.Hulls = self.Decompose(id)

	// constrain hulls to self intersects
	self.Hulls, _, constVertexSet = constrain.ToSelfIntersects(
		id, self.Hulls, self.Polyline(), self.Options(), constVerts,
	)

	var db = hdb.NewHdb()
	var selections map[*node.Node]struct{}
	var deformables = make([]node.Node, 0, len(self.Hulls))

	for i := range self.Hulls {
		deformables = append(deformables, self.Hulls[i])
	}
	// empty deque, this is for future splits
	node.Clear(&self.Hulls)
	db.Load(deformables)

	for len(deformables) > 0 {
		// 0. find deformable node
		selections = findDeformableNodes(deformables, db)
		// 1. deform selected nodes
		if len(selections) > 0 {
			deformables = deformNodes(id, selections)
			// 2. remove selected nodes from db
			cleanUpDB(db, selections)
			// 3. add new deformations to db
			db.Load(deformables)
		} else {
			deformables = deformables[:0]
		}
		// 4. repeat until there are no deformables
	}

	self.AggregateSimpleSegments(
		id, db, constVertexSet,  self.ValidateMerge,
	)

	node.Clear(&self.Hulls)
	self.SimpleSet.Empty()

	var hull *node.Node
	var nodes = db.All()
	for i := range nodes {
		hull = nodes[i]
		self.Hulls = append(self.Hulls, *hull)
		self.SimpleSet.Extend(hull.Range.I, hull.Range.J)
	}

	//state
	self.State().MarkClean()

	return self
}
