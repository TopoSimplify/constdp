package constdp

import (
    "github.com/intdxdt/rtree"
)

func cleanUpDB(hulldb *rtree.RTree, selections map[string]*rtree.Obj) {
    for i := range selections {
        hulldb.RemoveObj(selections[i])
    }
}
