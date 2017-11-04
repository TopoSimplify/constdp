package constdp

import (
    "simplex/lnr"
    "simplex/node"
    "simplex/opts"
    "github.com/intdxdt/fan"
    "github.com/intdxdt/sset"
    "github.com/intdxdt/rtree"
	"simplex/common"
)



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
    simplifyClass(selfs, junctions)

    var selections map[string]*node.Node
    var hulldb = rtree.NewRTree(rtreeBucketSize)
    var boxes = make([]rtree.BoxObj, 0)

    var deformables = make([]*node.Node, 0)
    for _, self := range selfs {
        self.selfUpdate()
        for _, hull := range self.Hulls {
            deformables = append(deformables, hull)
            boxes = append(boxes, hull)
        }
        node.Clear(&self.Hulls)  // empty deque, this is for future splits
    }
    hulldb.Load(boxes)

    for len(deformables) > 0 {
        // 1. find deformable node
        selections = findDeformableNodes(deformables, hulldb)
        // 2. deform selected nodes
        deformables = deformNodes(selections)
        // 2. remove selected nodes from db
        cleanUpDB(hulldb, selections)
        // 2. add new deformations to db
        updateDB(hulldb, deformables)
        // 3. repeat until there are no deformables
    }
    groupHullsByFC(hulldb)
}

func simplifyClass(selfs []*ConstDP, junctions map[string]*sset.SSet) {
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
            constVerts = common.AsInts(v.Values())
        } else {
            constVerts = make([]int, 0)
        }
        self.Simplify(constVerts)
        return self
    }

    var out = fan.Stream(stream, worker, concurProcs, exit)
    for range out {}
}
