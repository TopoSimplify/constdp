package constdp

import (
    "runtime"
    "simplex/db"
    "simplex/nque"
    "simplex/opts"
    "simplex/node"
    "simplex/constrain"
    "github.com/intdxdt/fan"
    "github.com/intdxdt/rtree"
    "github.com/intdxdt/avl"
    "github.com/intdxdt/cmp"
)

func processFeatClassNodes(selfs []*ConstDP, opts *opts.Opts) {
    var ConCur = 8

    var bln bool
    var queue = nque.NewQueue()
    var historyMap = avl.NewAVL(cmp.Str)

    var hlist = make([]*node.Node, 0)
    var hulldb = db.NewDB(RtreeBucketSize)
    var boxes = make([]rtree.BoxObj, 0)

    for _, self := range selfs {
        self.selfUpdate()
        for _, o := range self.Hulls.Nodes() {
            hull := castAsNode(o)
            queue.Append(hull)
            historyMap.Insert(hull.Id())
            hlist = append(hlist, hull)
            boxes = append(boxes, hull)
        }
        self.Hulls.Clear() // empty deque, this is for future splits
    }

    hulldb.Load(boxes)

    var stream = make(chan interface{})
    var exit = make(chan struct{})

    //go pool
    go func() {
        for {
            select {
            case <-exit:
                return
            default:
                var flip = false
                for !queue.IsEmpty() {
                    if flip {
                        stream <- queue.PopLeft()
                    }else{
                        stream <- queue.Pop()
                    }
                    flip  = !flip
                }
                runtime.Gosched()
            }

        }
    }()

    //var historyMap =  make(map[string]bool)
    var worker = func(v interface{}) interface{} {
        var selections = node.NewNodes()
        //fmt.Println("queue size :", queue.Len())
        // assume poped hull to be valid
        var hull = castAsNode(v)
        var self = hull.Instance.(*ConstDP)

        // insert hull into hull db
        //hulldb.Insert(hull)

        //check state in history map
        if !historyMap.Contains(hull.Id()) {
            return struct{}{} //continue
        }

        // find hull neighbours
        // self intersection constraint
        // can self intersect with itself but not with other lines
        bln = constrain.ByFeatureClassIntersection(self.Options(), hull, hulldb, selections)

        if !selections.IsEmpty() {
            deformClassSelections(queue, hulldb, selections, historyMap)
        }

        if !bln {
            return struct{}{}
        }

        // context_geom geometry constraint
        self.ValidateContextRelation(hull, selections)

        if !selections.IsEmpty() {
            deformClassSelections(queue, hulldb, selections, historyMap)
        }
        return struct{}{}
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
                if queue.IsEmpty() && pool.IsIdle() {
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
