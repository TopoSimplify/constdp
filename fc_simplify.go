package constdp

import (
	"github.com/intdxdt/fan"
	"github.com/TopoSimplify/lnr"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/opts"
	"github.com/TopoSimplify/hdb"
)

//Simplify a feature class of linear geometries
//optional callback for the number of deformables
func SimplifyFeatureClass(selfs []*ConstDP, opts *opts.Opts, callback ... func(n int)) {
	var deformableCallback = func(n int) {}
	if len(callback) > 0 {
		deformableCallback = callback[0]
	}

	var junctions = make(map[string][]int)

	if opts.PlanarSelf {
		instances := make([]*lnr.FC, len(selfs))
		for i, sf := range selfs {
			instances[i] = lnr.NewFC(sf.Coordinates(), sf.Id())
		}
		junctions = lnr.FCPlanarSelfIntersection(instances)
	}

	SimplifyDPs(selfs, junctions)

	var constBln = opts.AvoidNewSelfIntersects || opts.PlanarSelf ||
		opts.GeomRelation || opts.DirRelation || opts.DistRelation

	var selections map[string]*node.Node
	var hulldb = hdb.NewHdb(rtreeBucketSize)
	var deformables []*node.Node

	for _, self := range selfs {
		self.selfUpdate()
		for i := range self.Hulls {
			deformables = append(deformables, self.Hulls[i])
		}
		hulldb.Load(self.Hulls)
		if constBln {
			node.Clear(&self.Hulls) // empty deque, this is for future splits
		}
	}
	//hulldb.Load(boxes)

	for constBln && len(deformables) > 0 {
		deformableCallback(len(deformables))
		// 1. find deformable node
		selections = findDeformableNodes(deformables, hulldb)
		// 2. deform selected nodes
		deformables = deformNodes(selections)
		// 3. remove selected nodes from db
		cleanUpDB(hulldb, selections)
		// 4. add new deformations to db
		hulldb.Load(deformables)
		// 5. repeat until there are no deformables
	}

	groupHullsByFC(hulldb)
}

func SimplifyDPs(selfs []*ConstDP, junctions map[string][]int) {
	var stream = make(chan interface{})
	var exit = make(chan struct{})
	defer close(exit)

	go inputStreamSimplifyDP(stream, selfs)
	var worker = processSimplifyDPs(junctions)
	var out = fan.Stream(stream, worker, concurProcs, exit)
	for range out {
	}
}

func inputStreamSimplifyDP(stream chan interface{}, selfs []*ConstDP) {
	for _, self := range selfs {
		stream <- self
	}
	close(stream)
}

func processSimplifyDPs(junctions map[string][]int) func(v interface{}) interface{} {
	return func(v interface{}) interface{} {
		var self = v.(*ConstDP)
		var constVerts []int
		if v, ok := junctions[self.Id()]; ok {
			constVerts = v
		} else {
			constVerts = make([]int, 0)
		}
		self.Simplify(constVerts)
		return self
	}
}
