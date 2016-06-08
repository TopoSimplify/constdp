package constdp

import (
    "simplex/dp"
    "simplex/struct/heap"
    "simplex/struct/stack"
    "simplex/struct/bst"
    . "simplex/interest"
    "simplex/geom"
)

//constrained dp simplify
func (self *ConstDP) Simplify(opts *dp.Options) *ConstDP {
    var n *bst.Node
    var node *dp.Node
    var constlist []geom.Geometry
    var candidates *IntCandidates

    self.opts = opts
    self.Simple.Reset()
    self.Filter(self.Root, self.opts.Threshold)


    for !(self.NodeSet.IsEmpty()) {
        n = self.AsBSTNode_Item(self.NodeSet.Shift())
        node = self.AsDPNode_BSTNode_Item(n)
        //early exit
        if node == nil {
            break
        }

        constlist = self.context_neighbours(node)

        //var comparators = self._cmptors(polygeom, constlist)
        candidates = self.int_candidates(n)

        //update homos
        self.Simple.AddSet(
            self.homos.UpdateHomotopy(
                self.Pln,
                candidates,
                self.opts.Relations,
                constlist,
            ).FindSpatialFit(node.Key[0], node.Key[1]),
        )
    }
    return self
}

func (self *ConstDP) _childints(n *bst.Node) *heap.Heap {
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


//lazy evaluation of int candidates
func (self *ConstDP) int_candidates(n *bst.Node) *IntCandidates {
    var node = self.AsDPNode(n)
    var node_ints = func() *heap.Heap {
        return node.Ints
    }
    var child_ints = func() *heap.Heap {
        return self._childints(n)
    }
    var functors = []IntFunctor{node_ints, child_ints, }
    return NewIntCandidates(functors)
}