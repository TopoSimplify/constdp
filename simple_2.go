package constdp

import (
    "github.com/TopoSimplify/hdb"
    "github.com/TopoSimplify/node"
)

func cleanUpDB(hulldb *hdb.Hdb, selections map[*node.Node]struct{}) {
    for o := range selections {
        hulldb.Remove(o)
    }
}
