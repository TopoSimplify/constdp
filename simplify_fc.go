import constrain
from sort import sort_hulls
from collections import deque
from deform import deform_hulls
from pylib.structs.sset import SSet
from pylib.structs.rtree import RTree
from constdp.intersection import linear_ftclass_self_intersection

//Update hull nodes with dp instance
func self_update(self){
    for hull := range self.hulls{
        hull.dp = self
    }
}



func deform_ftclass_selections(queue, hulldb, selections){
    for s in selections{
        self = s.DP
        deform_hulls(self, hulldb, [s])
        self_update(self=self)
        for  len(self.hulls){
            queue.appendleft(self.hulls.pop())
        }
    }
    del selections[:] //empty selections
}

// Group hulls in hulldb by instance of ConstDP
func group_hulls_by_self(hulldb){
    smap = make(map[string][]*HullNode)
    for h in hulldb.all() {
        self = h.dp
        lst = smap.get(self.id, [])
        lst.append(h)
        smap[self.id] = lst
    }

    selfs = []
    for key, lst in smap.items(){
        self = lst[0].dp
        self.hulls.clear()
        self.hulls.extend(sort_hulls(lst))
        selfs.append(self)
    }

    for self in selfs{
        self.simple.empty() //update new simple
        for h in self.hulls{
            self.simple.add(*h.range)
        }
    }
}

//Simplify a feature class of linear geometries
func simplify_featureclass(selfs, opts){
    junctions = make(map[]SSet,0)
    if opts.keep_selfintersects{
        junctions = linear_ftclass_self_intersection(selfs)
    }
    // return common.simple_hulls_as_ptset
    for self := range selfs{
        const_verts = junctions.get(self.id, SSet())
        self.simplify(opts=opts, const_verts=const_verts.values())
    }

    queue = []
    hulldb = RTree(8, attribute=('minx', 'miny', 'maxx', 'maxy'))
    for self in selfs{
        self_update(self)
        queue.extend([h for h in self.hulls])
        self.hulls.clear()  // empty deque, this is for future splits
    }

    selections = []
    queue = deque(queue)
    while len(queue)>0:
        // assume poped hull to be valid
        hull = queue.popleft()
        self = hull.dp
        // insert hull into hull db
        hulldb.insert(hull)
        // find hull neighbours
        // self intersection constraint
        // can self intersect with itself but not with other lines
        bln = constrain.self_ftclass_intersection(self=self, hull=hull, hulldb=hulldb, selections=selections)

        if len(selections) > 0{
            deform_ftclass_selections(queue=queue, hulldb=hulldb, selections=selections)
        }

        if ! bln{
            continue
        }

        // context_geom geometry constraint
        constrain.context_relation(self=self, hull=hull, selections=selections)

        if len(selections) > 0{
            deform_ftclass_selections(queue=queue, hulldb=hulldb, selections=selections)
        }

    group_hulls_by_self(hulldb=hulldb)
}


