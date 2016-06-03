package constdp

import (
    "simplex/dp"
    "sort"
    "simplex/util/iter"
    "simplex/geom"
    "simplex/struct/slist"
    "simplex/constrelate"
)

type Morph struct {
    pln     []*geom.Point
    queue     *slist.SList
    relations []constrelate.Relation
}

func (self *Morph) IsValid(subrange *[2]int) bool{
    intobj := list.Pop().(*dp.Vertex)
    nextint := intobj.Index()

    subrange := append(subrange, nextint)
    sort.Ints(subrange)

    //subrange is sorted
    self.filter_subrange(subrange, nextint, fixint)

    subpoly := self.subpoly(iter.NewGenerator_AsVals(subrange...))
    subgeom := geom.NewLineString(subpoly)
}

func (self *Morph) HasNext() bool{
    return !self.queue.IsEmpty()
}


func (self *Morph) NextCandidate() *dp.Vertex{
    v, ok := self.queue.Pop().(*dp.Vertex)
    if ok {
        return v
    }
    return nil
}




func (self *ConstDP) candidate_morphs(list *slist.SList) {
    intobj := list.Pop().(*dp.Vertex)
    nextint := intobj.Index()

    subrange := append(subrange, nextint)
    sort.Ints(subrange)

    //subrange is sorted
    self.filter_subrange(subrange, nextint, fixint)

    subpoly := self.subpoly(iter.NewGenerator_AsVals(subrange...))
    subgeom := geom.NewLineString(subpoly)
}
