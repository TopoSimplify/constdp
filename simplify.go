package constdp

import (
	"sort"
	"simplex/struct/rtree"
	"simplex/constdp/opts"
)

//homotopic simplification at a given threshold
func (self *ConstDP) Simplify(opts *opts.Opts) *ConstDP {
	self.Opts   = opts
	self.Simple = make([]*HullNode, 0)
	self.Hulls  = self.decompose(opts.Threshold)

	// constrain hulls to self intersects
	self.Hulls, _ = self.constrain_to_selfintersects(opts)

	// for _, h := range *self.Hulls.DataView() {
	// 	fmt.Println(h)
	// }

	var bln bool
	var hull *HullNode
	var hlist []*HullNode

	var hulldb = rtree.NewRTree(8)
	for !self.Hulls.IsEmpty() {
		// assume poped hull to be valid
		bln = true

		// pop hull in queue
		hull = self.Hulls.PopLeft().(*HullNode)

		// insert hull into hull db
		hulldb.Insert(hull)

		// self intersection constraint
		hlist, bln = self.constrain_self_intersection(hull, hulldb, bln)
		self.deform_hull(hulldb, hlist)

		if !bln {
			continue
		}

		// context_geom geometry constraint
		hlist, bln = self.constrain_context_relation(hull, hulldb, bln)
		self.deform_hull(hulldb, hlist)

	}

	self.Simple = self.merge_simple_segments(self, opts, hulldb)
	return self

	hdata := make([]*HullNode, 0)
	for _, h := range hulldb.All() {
		hdata = append(hdata, h.GetItem().(*HullNode))
	}
	sort.Sort(HullNodes(hdata))
	self.Simple = hdata
	return self
}
