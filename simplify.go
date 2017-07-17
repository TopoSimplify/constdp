package constdp

import (
    "simplex/struct/rtree"
    "simplex/struct/deque"
    "simplex/struct/sset"
    "fmt"
    "sort"
    "simplex/geom"
)

func (self *ConstDP) split_hulls_at_selfintersects(dphulls *deque.Deque) *deque.Deque{
    db := rtree.NewRTree(8)
    self_inters := LinearSelfIntersection(self.Pln)
    data := make([]rtree.BoxObj, 0)
    for _, v := range *dphulls.DataView(){
        h := v.(*HullNode)
        data = append(data, h)
    }
    db.Load(data)
    at_vertex_set := sset.NewSSet(IntCmp)

    for _,inter := range self_inters{
        if inter.IsSelfVertex(){
            at_vertex_set.Union(inter.Meta.SelfVertices)
        }
    }

    for _, inter := range self_inters{
        if !inter.IsSelfVertex(){
            continue
        }

        hulls := dbKNN(db, inter, 1.e-5)
        for _,h := range hulls{
            hull := h.(*HullNode)
            idxs  := inter.Meta.SelfVertices.Values()
            indices := make([]int, 0)
            for _, o := range idxs{
                indices = append(indices, o.(int))
            }
            hsubs := self.split_hull_at_index( hull, indices)

            if len(hsubs) > 0{
                db.Remove(hull)
            }

            keep, rm := self.merge_contig_fragments(
                 hsubs, db, at_vertex_set,
            )

            for _, h := range rm{
                db.Remove(h)
            }

            for _, h := range keep{
                db.Insert(h)
            }
        }
    }

    hdata := make([]*HullNode, 0)
    for _, h:= range db.All(){
        hdata= append(hdata, h.GetItem().(*HullNode))
    }
    sort.Sort(HullNodes(hdata))
    hulls := deque.NewDeque()
    for _, hn := range hdata{
        hulls.Append(hn)
    }
    return hulls
}


//homotopic simplification at a given threshold
func  (self *ConstDP) Simplify(opts *Opts) *ConstDP{
    self.Simple = make([]*HullNode, 0)
    self.Hulls  = self.dp_decompose( opts.Threshold)

    // split hulls by self intersects
    if opts.KeepSelfIntersects {
        self.Hulls = self.split_hulls_at_selfintersects( self.Hulls)
    }

    for h := range *self.Hulls.DataView(){
        fmt.Println(h)
    }

    hulldb := rtree.NewRTree(8)
    for self.Hulls.Len() > 0 {
        // assume poped hull to be valid
        bln := true

        // pop hull in queue
        hull := self.Hulls.PopLeft().(*HullNode)

        // insert hull into hull db
        hulldb.Insert(hull)

        if bln && self.Opts.AvoidNewSelfIntersects{
            // find hull neighbours
            hlist := find_hull_deformation_list(hulldb, hull, self.Opts)
            for _,h := range hlist{
                bln = !(h == hull)
                self.deform_hull(hulldb, h)
            }
        }

        if !bln{
            continue
        }

        // find context neighbours - if valid
        ctxs := dbKNN(self.CtxDB, hull, self.Opts.MinDist)
        i := 0
        for  bln && i < len(ctxs) {
            ctx := ctxs[i].(geom.Geometry)
            if bln && self.Opts.GeomRelation {
                bln = self.is_geom_relate_valid( hull, ctx)
            }

            if bln && self.Opts.DistRelation {
                bln = self.is_dist_relate_valid( hull, ctx)
            }

            if bln && self.Opts.DirRelation {
                bln = self.is_dir_relate_valid(hull, ctx)
            }

            i += 1
        }

        if !bln {
            self.deform_hull(hulldb, hull)
        }}
    hdata := make([]*HullNode, 0)
    for _, h := range hulldb.All(){
        hdata = append(hdata, h.GetItem().(*HullNode))
    }
    sort.Sort(HullNodes(hdata))
    self.Simple = hdata
    return self
}

func (self *ConstDP) deform_hull( hulldb *rtree.RTree, hull *HullNode){
    // split hull at maximum_offset offset
    ha, hb := self.split_hull( hull)
    hulldb.Remove(hull)

    self.Hulls.AppendLeft(hb)
    self.Hulls.AppendLeft(ha)
}

func (self *ConstDP) is_geom_relate_valid(hull *HullNode, ctx geom.Geometry) bool{
    seg     := self.HullSegment( hull)
    subpln  := self.Pln.SubPolyline(hull.Range)

    ln_geom  := subpln.geom
    seg_geom := seg
    ctx_geom := ctx

    ln_g_inter  := ln_geom.Intersects(ctx_geom)
    seg_g_inter := seg_geom.Intersects(ctx_geom)

    bln := true
    if seg_g_inter && (! ln_g_inter){
        bln = false
    }else if (! seg_g_inter) && ln_g_inter{
        bln = false
    }
    // both intersects & disjoint
    return bln
}


//is distance relate valid ?
func (self *ConstDP) is_dist_relate_valid(hull *HullNode, ctx geom.Geometry) bool{
    mindist := self.Opts.MinDist
    seg     := self.HullSegment( hull)
    ln_geom := hull.Pln.geom

    seg_geom := seg
    ctx_geom := ctx

    _or := ln_geom.Distance(ctx_geom)  // original relate
    dr  := seg_geom.Distance(ctx_geom)  // new relate

    bln := dr >= mindist
    if !bln && _or < mindist{
        bln = (dr >= _or)
    }
    return bln
}


func (self *ConstDP) is_dir_relate_valid(hull *HullNode, ctx geom.Geometry) bool{
    subpln  := self.Pln.SubPolyline(hull.Range)
    segment := NewPolyline([]*geom.Point{
        self.Pln.coords[hull.Range.i],
        self.Pln.coords[hull.Range.j],
    })

    lnr   := DirectionRelate(subpln, ctx)
    segr  := DirectionRelate(segment, ctx)

    return lnr == segr
}
