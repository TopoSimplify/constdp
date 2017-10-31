package constdp

import (
	"simplex/dp"
	"simplex/lnr"
	"simplex/node"
	"simplex/opts"
	"simplex/split"
	"simplex/constrain"
	"github.com/intdxdt/sset"
	"github.com/intdxdt/deque"
	"github.com/intdxdt/rtree"
	"github.com/intdxdt/fan"
)

const RtreeBucketSize = 4

//Update hull nodes with dp instance
func (self *ConstDP) selfUpdate() {
	var hull *node.Node
	for _, h := range *self.Hulls.DataView() {
		hull = castAsNode(h)
		hull.Instance = self
	}
}

func deformClassSelections(queue *deque.Deque, hulldb *rtree.RTree, selections *node.Nodes) {
	for _, s := range selections.DataView() {
		self := castConstDP(s.Instance)
		sels := node.NewNodes().Push(s)
		split.SplitNodesInDB(
			self.NodeQueue(), hulldb, sels,
			self.Score, dp.NodeGeometry,
		)
		self.selfUpdate()
		for self.Hulls.Len() > 0 {
			queue.AppendLeft(self.Hulls.Pop())
		}
	}
	selections.Empty() //empty selections
}

// Group hulls in hulldb by instance of ConstDP
func groupHullsBySelf(hulldb *rtree.RTree) {
	var ok bool
	var hull *node.Node
	var selfs = make([]*ConstDP, 0)
	var smap = make(map[string]*node.Nodes)
	for _, h := range nodesFromRtreeNodes(hulldb.All()).DataView() {
		var lst *node.Nodes
		var self = castConstDP(h.Instance)
		var id = self.Id()
		if lst, ok = smap[id]; !ok {
			lst = node.NewNodes()
		}
		lst.Push(h)
		smap[id] = lst
	}

	for _, lst := range smap {
		var self = castConstDP(lst.Get(0).Instance)
		self.Hulls.Clear()
		for _, h := range lst.Sort().DataView() {
			self.Hulls.Append(h)
		}
		selfs = append(selfs, self)
	}

	for _, self := range selfs {
		self.SimpleSet.Empty() //update new simple
		for _, h := range *self.Hulls.DataView() {
			hull = castAsNode(h)
			self.SimpleSet.Extend(hull.Range.I(), hull.Range.J())
		}
	}
}

//Simplify a feature class of linear geometries
func SimplifyFeatureClass(selfs []*ConstDP, opts *opts.Opts) {
	var junctions = make(map[string]*sset.SSet, 0)

	if opts.KeepSelfIntersects {
		instances := make([]lnr.Linear, len(selfs))
		for i, v := range selfs {
			instances[i] = v
		}
		junctions = lnr.FeatureClassSelfIntersection(instances)
	}


	simplifyClass(selfs, opts, junctions)

	var hlist = make([]*node.Node, 0)
	var hulldb = rtree.NewRTree(RtreeBucketSize)
	for _, self := range selfs {
		self.selfUpdate()
		for _, h := range *self.Hulls.DataView() {
			hlist = append(hlist, castAsNode(h))
		}
		self.Hulls.Clear() // empty deque, this is for future splits
	}



	var bln bool
	var self *ConstDP
	var hull *node.Node
	var selections = node.NewNodes()
	var dque = deque.NewDeque(len(hlist))

	for _, h := range hlist {
		dque.Append(h)
	}

	for !dque.IsEmpty() {
		//fmt.Println("queue size :", dque.Len())
		// assume poped hull to be valid
		hull = castAsNode(dque.PopLeft())
		self = hull.Instance.(*ConstDP)

		// insert hull into hull db
		hulldb.Insert(hull)

		// find hull neighbours
		// self intersection constraint
		// can self intersect with itself but not with other lines
		bln = constrain.ByFeatureClassIntersection(self.Options(), hull, hulldb, selections)

		if !selections.IsEmpty() {
			deformClassSelections(dque, hulldb, selections)
		}

		if !bln {
			continue
		}

		// context_geom geometry constraint
		self.ValidateContextRelation(hull, selections)

		if !selections.IsEmpty() {
			deformClassSelections(dque, hulldb, selections)
		}
	}
	groupHullsBySelf(hulldb)
}

func simplifyClass(selfs []*ConstDP, opts *opts.Opts, junctions map[string]*sset.SSet)  {
	var stream = make(chan interface{})
	var exit = make(chan struct{})
	defer close(exit)

	go func() {
		for _, self := range selfs {
			stream <- self
		}
		close(stream)
	}()

	var worker = func(v interface{}) interface{} {
		var self = v.(*ConstDP)
		var constVerts []int
		if v, ok := junctions[self.Id()]; ok {
			constVerts = asInts(v.Values())
		} else {
			constVerts = make([]int, 0)
		}
		self.Simplify(constVerts)
		return self
	}

	var out = fan.Stream(stream, worker, 8, exit)
	//var results = make([]*ConstDP, 0)
	for _ = range out {
		//results = append(results, o.(*ConstDP))
	}
	//return results
}
