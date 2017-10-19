package constdp

import (
	"github.com/intdxdt/sset"
	"simplex/opts"
	"github.com/intdxdt/deque"
	"github.com/intdxdt/rtree"
)

//Update hull nodes with dp instance
func (self *ConstDP) self_update() {
	var hull *HullNode
	for _, h := range *self.Hulls.DataView() {
		hull = cast_as_hullnode(h)
		hull.DP = self
	}
}

func deform_class_selections(queue *deque.Deque, hulldb *rtree.RTree, selections *HullNodes) {
	for _, s := range selections.list {
		self := s.DP
		sels := NewHullNodes().Push(s)
		self.deform_hulls(hulldb, sels)
		self.self_update()
		for self.Hulls.Len() > 0 {
			queue.AppendLeft(self.Hulls.Pop())
		}
	}
	selections.Empty() //empty selections
}

// Group hulls in hulldb by instance of ConstDP
func group_hulls_by_self(hulldb *rtree.RTree) {
	var ok bool
	var self *ConstDP
	var hull *HullNode
	var selfs = make([]*ConstDP, 0)
	var smap = make(map[string]*HullNodes)

	for _, h := range NewHullNodesFromNodes(hulldb.All()).list {
		var lst *HullNodes
		self = h.DP
		if lst, ok = smap[self.Id]; !ok {
			lst = NewHullNodes()
		}
		lst.Push(h)
		smap[self.Id] = lst
	}

	for _, lst := range smap {
		self = lst.Get(0).DP
		self.Hulls.Clear()
		for _, h := range lst.Sort().list {
			self.Hulls.Append(h)
		}
		selfs = append(selfs, self)
	}

	for _, self := range selfs {
		self.Simple.Empty() //update new simple
		for _, h := range *self.Hulls.DataView() {
			hull = cast_as_hullnode(h)
			self.Simple.Extend(hull.Range.I(), hull.Range.J())
		}
	}
}

//Simplify a feature class of linear geometries
func SimplifyFeatureClass(selfs []*ConstDP, opts *opts.Opts) {
	var junctions = make(map[string]*sset.SSet, 0)
	if opts.KeepSelfIntersects {
		junctions = linear_ftclass_self_intersection(selfs)
	}

	// return common.simple_hulls_as_ptset
	for _, self := range selfs {
		var const_verts []int
		if v, ok := junctions[self.Id]; ok {
			const_verts = as_ints(v.Values())
		} else {
			const_verts = make([]int, 0)
		}
		self.Simplify(opts, const_verts)
	}

	var hlist = make([]*HullNode, 0)
	var hulldb = rtree.NewRTree(8)
	for _, self := range selfs {
		self.self_update()
		for _, h := range *self.Hulls.DataView() {
			hlist = append(hlist, cast_as_hullnode(h))
		}
		self.Hulls.Clear() // empty deque, this is for future splits
	}

	var bln bool
	var self *ConstDP
	var hull *HullNode
	var selections = NewHullNodes()
	var dque = deque.NewDeque(len(hlist))

	for _, h := range hlist {
		dque.Append(h)
	}

	for !dque.IsEmpty() {
		//fmt.Println("queue size :", dque.Len())
		// assume poped hull to be valid
		hull = cast_as_hullnode(dque.PopLeft())
		self = hull.DP

		// insert hull into hull db
		hulldb.Insert(hull)

		// find hull neighbours
		// self intersection constraint
		// can self intersect with itself but not with other lines
		bln = self.constrain_ftclass_intersection(hull, hulldb, selections)

		if !selections.IsEmpty() {
			deform_class_selections(dque, hulldb, selections)
		}

		if !bln {
			continue
		}

		// context_geom geometry constraint
		self.constrain_context_relation(hull, selections)

		if !selections.IsEmpty() {
			deform_class_selections(dque, hulldb, selections)
		}
	}
	group_hulls_by_self(hulldb)
}
