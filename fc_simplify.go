package constdp

import (
    "simplex/node"
    "simplex/opts"
    "github.com/intdxdt/fan"
    "github.com/intdxdt/sset"
    "github.com/intdxdt/rtree"
    "simplex/lnr"
)

const rtreeBucketSize = 4
const concurProcs = 8

//Update hull nodes with dp instance
func (self *ConstDP) selfUpdate() {
    var hull *node.Node
    for _, h := range *self.Hulls.DataView() {
        hull = castAsNode(h)
        hull.Instance = self
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
    simplifyClass(selfs, junctions)

    var hull *node.Node
    var selections map[string]*node.Node
    var hlist = make([]*node.Node, 0)
    var hulldb = rtree.NewRTree(rtreeBucketSize)
    var boxes = make([]rtree.BoxObj, 0)

    for _, self := range selfs {
        self.selfUpdate()
        for _, o := range *self.Hulls.DataView() {
            hull = castAsNode(o)
            hlist = append(hlist, hull)
            boxes = append(boxes, hull)
        }
        self.Hulls.Clear() // empty deque, this is for future splits
    }
    hulldb.Load(boxes)

    for len(hlist) > 0 {
        selections = findDeformableNodes(hlist, hulldb)
        hlist = deformNodes(selections)
        cleanUpDB(hulldb, selections)
        updateDB(hulldb, hlist)
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
            constVerts = asInts(v.Values())
        } else {
            constVerts = make([]int, 0)
        }
        self.Simplify(constVerts)
        return self
    }

    var out = fan.Stream(stream, worker, concurProcs, exit)
    for range out {
    }
}
