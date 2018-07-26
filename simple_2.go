package constdp

import (
    "github.com/TopoSimplify/hdb"
    "github.com/TopoSimplify/node"
)

func cleanUpDB(hulldb *hdb.Hdb, selections map[string]*node.Node) {
    for i := range selections {
        hulldb.Remove(selections[i])
    }
}
