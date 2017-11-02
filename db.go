package constdp

import (
    "simplex/node"
    "github.com/intdxdt/rtree"
)

func cleanUpDB(hulldb *rtree.RTree, selections map[string]*node.Node) {
    for _, n := range selections {
        hulldb.Remove(n)
    }
}

func updateDB(hulldb *rtree.RTree, nodes []*node.Node) {
    for _, n := range nodes {
        hulldb.Insert(n)
    }
}
