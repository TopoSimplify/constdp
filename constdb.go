package constdp

import (
    "simplex/struct/rtree"
    "simplex/geom/mbr"
    "simplex/geom"
)


//in-memory rtree
func NewConstDB(geometries []geom.Geometry) *rtree.RTree {
    var objs = make([]rtree.BoxObj, len(geometries))
    for i := range geometries {
        objs[i] = geometries[i]
    }
    return rtree.NewRTree(16).Load(objs)
}

func SearchDb(db *rtree.RTree, query *mbr.MBR) []geom.Geometry {
    nodes := db.Search(query)
    geoms := make([]geom.Geometry, len(nodes)) //
    for i := range nodes {
        geoms[i] = nodes[i].GetItem().(geom.Geometry)
    }
    return geoms
}