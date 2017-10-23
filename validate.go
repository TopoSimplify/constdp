package constdp

import (
	"simplex/node"
	"simplex/constrain"
	"github.com/intdxdt/rtree"
)


func (self *ConstDP) ValidateMerge(hull *node.Node, hulldb *rtree.RTree) bool {
	var bln = true
	var sideEffects = node.NewNodes()
	var options = self.Options()

	if bln && options.AvoidNewSelfIntersects {
		// self intersection constraint
		bln = constrain.SelfIntersection(self, hull, hulldb, sideEffects)
	}

	if !sideEffects.IsEmpty() || !bln {
		return false
	}

	// context geometry constraint
	bln = constrain.ContextRelation(self, self.ContextDB, hull, sideEffects)
	return sideEffects.IsEmpty() && bln
}
