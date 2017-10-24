package constdp

import (
	"simplex/knn"
	"simplex/node"
	"simplex/constrain"
	"github.com/intdxdt/rtree"
)

func (self *ConstDP) ValidateMerge(hull *node.Node, hulldb *rtree.RTree) bool {
	var bln = true
	var sideEffects = node.NewNodes()
	var options = self.Options()

	if options.AvoidNewSelfIntersects {
		// self intersection constraint
		bln = constrain.BySelfIntersection(self, hull, hulldb, sideEffects)
	}

	if !sideEffects.IsEmpty() || !bln {
		return false
	}

	// context geometry constraint
	bln = self.ValidateContextRelation(hull, sideEffects)
	return sideEffects.IsEmpty() && bln
}

//Constrain for context neighbours
// finds the collapsibility of hull with respect to context hull neighbours
// if hull is deformable, its added to selections
func (self *ConstDP) ValidateContextRelation(hull *node.Node, selections *node.Nodes) bool {
	if !self.checkContextRelations() {
		return true
	}

	var bln = true
	var options = self.Options()
	// find context neighbours - if valid
	var ctxs = knn.FindNeighbours(self.ContextDB, hull, options.MinDist)
	for _, contxt := range ctxs {
		if !bln {
			break
		}

		var cg = castAsContextGeom(contxt)
		if bln && options.GeomRelation {
			bln = constrain.ByGeometricRelation(self, hull, cg)
		}

		if bln && options.DistRelation {
			bln = constrain.ByMinDistRelation(self, hull, cg)
		}

		if bln && options.DirRelation {
			bln = constrain.BySideRelation(self, hull, cg)
		}
	}

	if !bln {
		selections.Push(hull)
	}

	return bln
}

func (self *ConstDP) checkContextRelations() bool {
	return self.Opts.GeomRelation || self.Opts.DistRelation || self.Opts.DirRelation
}
