package constdp

import (
	"github.com/intdxdt/rtree"
	"github.com/TopoSimplify/ctx"
	"github.com/TopoSimplify/knn"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/constrain"
)

func (self *ConstDP) ValidateMerge(hull *node.Node, hulldb *rtree.RTree) bool {
	var bln = true
	var sideEffects = make([]*node.Node, 0)

	// self intersection constraint
	if self.Opts.AvoidNewSelfIntersects {
		bln = constrain.BySelfIntersection(self.Opts, hull, hulldb, &sideEffects)
	}

	if len(sideEffects) != 0 || !bln {
		return false
	}

	// context geometry constraint
	bln = self.ValidateContextRelation(hull, &sideEffects)
	return bln && (len(sideEffects) == 0)
}

//Constrain for context neighbours
// finds the collapsibility of hull with respect to context hull neighbours
// if hull is deformable, its added to selections
func (self *ConstDP) ValidateContextRelation(hull *node.Node, selections *[]*node.Node) bool {
	if !(self.Opts.GeomRelation || self.Opts.DistRelation || self.Opts.DirRelation) {
		return true
	}
	var bln = true

	// find context neighbours - if valid
	var boxObjs = knn.FindNeighbours(self.ContextDB, hull, self.Opts.MinDist)

	var neighbours = make([]*ctx.ContextGeometry, len(boxObjs))
	for i , o := range boxObjs{
		neighbours[i] = o.(*ctx.ContextGeometry)
	}
	var ctxtgeoms = (&ctx.ContextGeometries{}).SetData(neighbours)

	if bln && self.Opts.GeomRelation {
		bln = constrain.ByGeometricRelation(hull, ctxtgeoms)
	}

	if bln && self.Opts.DistRelation {
		bln = constrain.ByMinDistRelation(self.Options(), hull, ctxtgeoms)
	}

	if bln && self.Opts.DirRelation {
		bln = constrain.BySideRelation(hull, ctxtgeoms)
	}

	if !bln {
		*selections = append(*selections, hull)
	}

	return bln
}
