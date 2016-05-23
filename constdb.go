package constdp
import (
    . "simplex/struct/rtree"
)

//in-memory rtree
func constdb(geometries []BoxObj) *RTree{
    return NewRTree(16).Load(geometries)
}

