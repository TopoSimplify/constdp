package constdp

import (
    "time"
    "runtime"
    "simplex/db"
    "simplex/lnr"
    "simplex/opts"
    "simplex/node"
    "simplex/constrain"
    "github.com/intdxdt/deque"
    "github.com/intdxdt/sset"
    "github.com/intdxdt/fan"
)

func processFeatClassNodes(selfs []*ConstDP, opts *opts.Opts) {
    var ConCur = 8
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
    var hulldb = db.NewDB(RtreeBucketSize)
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

    var dict = NewMap()

    for _, h := range hlist {
        dque.Append(h)
        dict.Set(h.Id())
    }

    var stream = make(chan interface{})
    var exit = make(chan struct{})

    //go pool
    go func() {
        for {
            for !dque.IsEmpty() {
                stream <- dque.PopLeft()
            }
            runtime.Gosched()
            time.Sleep(8 * time.Millisecond)
        }
    }()

    var worker = func(v interface{}) interface{} {
        //fmt.Println("queue size :", dque.Len())
        // assume poped hull to be valid
        hull = castAsNode(v)
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
            return bln
        }

        // context_geom geometry constraint
        self.ValidateContextRelation(hull, selections)

        if !selections.IsEmpty() {
            deformClassSelections(dque, hulldb, selections)
        }
        return bln
    }

    var done = make(chan struct{})
    var pool = fan.NewPool(stream, worker, ConCur, exit)

    var out = pool.Start()
    go func() {
        for {
            select {
            case <-out:
            default:
                //halting condition all data served and processed
                if dque.IsEmpty() && pool.IsIdle() {
                    close(exit)
                    close(done)
                    return
                }
            }
        }
    }()

    <-done
    groupHullsBySelf(hulldb)
}
