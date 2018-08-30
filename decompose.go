package constdp

import (
	"github.com/intdxdt/iter"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/common"
	"github.com/TopoSimplify/decompose"
)

func (self *ConstDP) Decompose(id *iter.Igen) []node.Node {
	return decompose.DouglasPeucker(
		id,
		self.Polyline(),
		self.Score,
		self.ScoreRelation,
		common.Geometry,
		self,
	)
}
