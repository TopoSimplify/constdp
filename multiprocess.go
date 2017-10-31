package constdp

//import (
//	"simplex/lnr"
//	"simplex/opts"
//	"simplex/node"
//	"simplex/constrain"
//	"github.com/intdxdt/rtree"
//	"github.com/intdxdt/deque"
//	"github.com/intdxdt/sset"
//)

//func processFeatClassNodes(selfs []*ConstDP, opts *opts.Opts) {
//	var junctions = make(map[string]*sset.SSet, 0)
//
//	if opts.KeepSelfIntersects {
//		instances := make([]lnr.Linear, len(selfs))
//		for i, v := range selfs {
//			instances[i] = v
//		}
//		junctions = lnr.FeatureClassSelfIntersection(instances)
//	}


//	simplifyClass(selfs, opts, junctions)
//
//	var hlist = make([]*node.Node, 0)
//	var hulldb = NewCRTree(RtreeBucketSize)
//	for _, self := range selfs {
//		self.selfUpdate()
//		for _, h := range *self.Hulls.DataView() {
//			hlist = append(hlist, castAsNode(h))
//		}
//		self.Hulls.Clear() // empty deque, this is for future splits
//	}
//
//	var bln bool
//	var self *ConstDP
//	var hull *node.Node
//	var selections = node.NewNodes()
//	var dque = deque.NewDeque(len(hlist))
//
//	for _, h := range hlist {
//		dque.Append(h)
//	}
//
//	var results = make([]interface{}, 0)
//	var stream = make(chan interface{})
//	var exit = make(chan struct{})
//	//go pool
//	go func() {
//		for !dque.IsEmpty() {
//			stream <- dque.PopLeft()
//		}
//	}()
//
//	var worker = func(v interface{}) interface{} {
//		//fmt.Println("queue size :", dque.Len())
//		// assume poped hull to be valid
//		hull = castAsNode(v)
//		self = hull.Instance.(*ConstDP)
//
//		// insert hull into hull db
//		hulldb.Insert(hull)
//
//		// find hull neighbours
//		// self intersection constraint
//		// can self intersect with itself but not with other lines
//		bln = constrain.ByFeatureClassIntersection(self.Options(), hull, hulldb, selections)
//
//		if !selections.IsEmpty() {
//			deformClassSelections(dque, hulldb, selections)
//		}
//
//		if !bln {
//			continue
//		}
//
//		// context_geom geometry constraint
//		self.ValidateContextRelation(hull, selections)
//
//		if !selections.IsEmpty() {
//			deformClassSelections(dque, hulldb, selections)
//		}
//	}
//
//	var done = make(chan struct{})
//	var pool = NewPool(stream, worker, ConCur, exit)
//
//	groupHullsBySelf(hulldb)
//}
