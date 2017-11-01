package constdp

import (
    "sort"
    "simplex/dp"
    "simplex/knn"
    "simplex/node"
    "simplex/merge"
    "github.com/intdxdt/sset"
    "simplex/db"
)

//Merge segment fragments where possible
func (self *ConstDP) AggregateSimpleSegments(nodeDB *db.DB,
    constVertexSet *sset.SSet, scoreRelation func(float64) bool,
    validateMerge func(node *node.Node, nodeDB *db.DB) bool) {

    var hull *node.Node
    var fragmentSize = 1
    var neighbours *node.Nodes
    var cache = make(map[[4]int]bool)
    var hulls = nodesFromRtreeNodes(nodeDB.All()).Sort().AsDeque()

    for !hulls.IsEmpty() {
        // from left
        hull = popLeftHull(hulls)

        if hull.Range.Size() != fragmentSize {
            continue
        }

        //make sure hull index is not part of vertex with degree > 2
        if constVertexSet.Contains(hull.Range.I()) || constVertexSet.Contains(hull.Range.J()) {
            continue
        }

        nodeDB.Remove(hull)

        // find context neighbours
        neighbours = nodesFromBoxes(knn.FindNodeNeighbours(nodeDB, hull, knn.EpsilonDist))

        // find context neighbours
        var prev, nxt = node.Neighbours(hull, neighbours)

        // find mergeable neighbours contiguous
        var key [4]int
        var mergePrev, mergeNxt *node.Node

        if prev != nil {
            key = cacheKey(prev, hull)
            if !cache[key] {
                addToMergeCache(cache, &key)
                mergePrev = merge.ContiguousFragmentsAtThreshold(self.Score, prev, hull, scoreRelation, dp.NodeGeometry)
            }
        }

        if nxt != nil {
            key = cacheKey(hull, nxt)
            if !cache[key] {
                addToMergeCache(cache, &key)
                mergeNxt = merge.ContiguousFragmentsAtThreshold(self.Score, hull, nxt, scoreRelation, dp.NodeGeometry)
            }
        }

        var merged bool
        //nxt, prev
        if !merged && mergeNxt != nil {
            nodeDB.Remove(nxt)
            if validateMerge(mergeNxt, nodeDB) {
                var h = castAsNode(hulls.First())
                if h == nxt {
                    hulls.PopLeft()
                }
                nodeDB.Insert(mergeNxt)
                merged = true
            } else {
                merged = false
                nodeDB.Insert(nxt)
            }
        }

        if !merged && mergePrev != nil {
            nodeDB.Remove(prev)
            //prev cannot exist since moving from left --- right
            if validateMerge(mergePrev, nodeDB) {
                nodeDB.Insert(mergePrev)
                merged = true
            } else {
                merged = false
                nodeDB.Insert(prev)
            }
        }

        if !merged {
            nodeDB.Insert(hull)
        }
    }
}

func cacheKey(a, b *node.Node) [4]int {
    var ij = [4]int{a.Range.I(), a.Range.J(), b.Range.I(), b.Range.J()}
    sort.Ints(ij[:])
    return ij
}

func addToMergeCache(cache map[[4]int]bool, key *[4]int) {
    cache[*key] = true
}
