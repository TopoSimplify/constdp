package constdp

import (
	"github.com/intdxdt/iter"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/common"
	"github.com/TopoSimplify/decompose"
	"github.com/TopoSimplify/offset"
)

func (self *ConstDP) Decompose(id *iter.Igen) []node.Node {
	var score = self.Score
	var relation = self.ScoreRelation

	if self.SquareScore != nil {
		score = self.SquareScore
		relation = self.SquareScoreRelation
	}

	var decomp = offset.EpsilonDecomposition{
		ScoreFn:  score,
		Relation: relation,
	}

	return decompose.DouglasPeucker(
		id, self.Polyline(), decomp, common.Geometry, self,
	)
}
