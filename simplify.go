package constdp

import (
    "simplex/dp"
    "simplex/geom"
    "simplex/util/iter"
    "simplex/struct/heap"
    "simplex/struct/stack"
    "simplex/struct/bst"
    "sort"
)

//constrained dp simplify
func (self *ConstDP) Simplify(opts *dp.Options) *ConstDP {
    self.opts = opts
    self.Simple.Reset()
    self.Filter(self.Root, self.opts.Threshold)
    for !(self.NodeSet.IsEmpty()) {
        self.genints()
    }
    return self
}

/*
 description generalize ints
 param dpfilter
 param options
 private
 */
func (self *ConstDP) genints() {
    var n = self.AsBSTNode_Item(self.NodeSet.Shift())
    var node = self.AsDPNode_BSTNode_Item(n)
    var fixint = node.Ints.Peek().(*dp.Vertex).Index()

    //early exit
    if node == nil {
        return
    }

    var subrange = []int{node.Key[0], node.Key[1]}
    var i, j = subrange[0], subrange[1]

    var poly = self.subpoly(
        iter.NewGenerator(i, j + 1),
    )
    var subpoly = self.subpoly(
        iter.NewGenerator_AsVals(subrange...),
    )

    var polygeom = geom.NewLineString(poly)
    var subgeom = geom.NewLineString(subpoly)
    var constlist = self.context_neighbours(node)

    //add intersect points with neighbours as constraints
    //self.updateconsts(constlist, polygeom, node, options)
    var nextint int
    if self.opts.Relations != nil && len(constlist) > 0 {
        var comparators = self._cmptors(polygeom, constlist)
        //intlist
        var intfuncs = self._intcandidates(n)
        var intfn  func() *heap.Heap
        intfn, intfuncs = intfuncs[0], intfuncs[1:]
        var curints *heap.Heap = intfn()
        //assume not valid
        var isvalid = false
        //proof otherwise
        for !isvalid {
            if len(subpoly) == len(poly) {
                self.Simple.Add(subrange...)
                isvalid = true
                continue
            }
            //check if subgeom is valid
            isvalid = self._isvalid(subgeom, comparators)

            if isvalid {
                self.Simple.Add(subrange...)
            } else {
                if !curints.IsEmpty() {
                    //index at end is -1
                    intobj := curints.Pop().(*dp.Vertex)
                    nextint = intobj.Index()
                    subrange = append(subrange, nextint)
                    sort.Ints(subrange)
                    //subrange is sorted
                    self.filter_subrange(subrange, nextint, fixint)
                    subpoly = self.subpoly(iter.NewGenerator_AsVals(subrange...))
                    subgeom = geom.NewLineString(subpoly)
                } else {
                    //reset
                    if len(intfuncs) > 0 {
                        subrange = []int{node.Key[0], node.Key[1]}
                        nextint = node.Ints.Peek().(*dp.Vertex).Index()
                        subrange = append(subrange, nextint) //keep top level node int
                        curints, intfuncs = intfuncs[0](), intfuncs[1:]
                    } else {
                        //go to original
                        subrange = iter.NewGenerator(i, j + 1).Values()
                        subpoly = poly
                    }
                }
            }
        }
    } else {
        //keep range interesting index
        self.Simple.Add(subrange...)
    }
}

func (self *ConstDP) _childints(n *bst.Node) *heap.Heap  {
    var node = self.AsDPNode(n)
    var stack = stack.NewStack()
    var nextint = node.Ints.Peek().(*dp.Vertex).Index()
    var intlist = heap.NewHeap(heap.NewHeapType().AsMax())
    var intobj    *dp.Vertex

    //node stack
    if n.Right != nil {
        stack.Add(n.Right)

    }
    if n.Left != nil {
        stack.Add(n.Left)
    }

    for !stack.IsEmpty() {
        node = self.AsDPNode_BSTNode_Any(stack.Pop())
        intobj = node.Ints.Peek().(*dp.Vertex)

        if nextint != intobj.Index() && intobj.Value() > 0.0 {
            intlist.Push(intobj)
        }
        if n.Right != nil {
            stack.Add(n.Right)
        }
        if n.Left != nil {
            stack.Add(n.Left)
        }
    }

    return intlist
}


/*
 description int candidates
 param node
 returns {*[]}
 */
func (self *ConstDP) _intcandidates(n *bst.Node) []func() *heap.Heap {
    var node = self.AsDPNode(n)
    var node_ints = func() *heap.Heap {
        return node.Ints
    }
    var child_ints = func() *heap.Heap {
        return self._childints(n)
    }
    return []func() *heap.Heap{node_ints, child_ints}
}