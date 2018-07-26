package constdp

import (
	"sort"
	"github.com/intdxdt/iter"
	"github.com/TopoSimplify/dp"
	"github.com/TopoSimplify/knn"
	"github.com/TopoSimplify/hdb"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/merge"
)

//Merge segment fragments where possible
func (self *ConstDP) AggregateSimpleSegments(
	db *hdb.Hdb, constVertexSet []int,
	scoreRelation func(float64) bool,
	validateMerge func(*node.Node, *hdb.Hdb) bool,
) {

	var fragmentSize = 1
	var neighbours []*node.Node
	var cache = make(map[[4]int]bool)
	//var objects = common.NodesFromObjects(db.All())
	var objects = db.All()
	sort.Sort(node.Nodes(objects))

	for len(objects) != 0 {
		hull := objects[0]
		objects = objects[1:]
		if hull.Range.Size() != fragmentSize {
			continue
		}

		//make sure hull index is not part of vertex with degree > 2
		if iter.SortedSearchInts(constVertexSet, hull.Range.I) &&
			iter.SortedSearchInts(constVertexSet, hull.Range.J) {
			continue
		}
		var withNext, withPrev = !iter.SortedSearchInts(constVertexSet, hull.Range.J),
			!iter.SortedSearchInts(constVertexSet, hull.Range.I)

		db.Remove(hull)

		// find context neighbours
		//neighbours = common.NodesFromObjects(knn.FindNodeNeighbours(db, hull, knn.EpsilonDist))
		neighbours = knn.FindNodeNeighbours(db, hull, knn.EpsilonDist)

		// find context neighbours
		var prev, nxt = node.Neighbours(hull, neighbours)

		// find mergeable neighbours contiguous
		var key [4]int
		var mergePrev, mergeNxt *node.Node

		if withPrev && prev != nil {
			key = cacheKey(prev, hull)
			if !cache[key] {
				addToMergeCache(cache, &key)
				mergePrev = merge.ContiguousFragmentsAtThreshold(
					self.Score, prev, hull, scoreRelation, dp.NodeGeometry,
				)
			}
		}

		if withNext && nxt != nil {
			key = cacheKey(hull, nxt)
			if !cache[key] {
				addToMergeCache(cache, &key)
				mergeNxt = merge.ContiguousFragmentsAtThreshold(
					self.Score, hull, nxt, scoreRelation, dp.NodeGeometry,
				)
			}
		}

		var merged bool
		//nxt, prev
		if !merged && mergeNxt != nil {
			db.Remove(nxt)
			if validateMerge(mergeNxt, db) {
				if objects[0] == nxt {
					objects = objects[1:]
				}
				db.Insert(mergeNxt)
				merged = true
			} else {
				merged = false
				db.Insert(nxt)
			}
		}

		if !merged && mergePrev != nil {
			db.Remove(prev)
			//prev cannot exist since moving from left --- right
			if validateMerge(mergePrev, db) {
				db.Insert(mergePrev)
				merged = true
			} else {
				merged = false
				db.Insert(prev)
			}
		}

		if !merged {
			db.Insert(hull)
		}
	}
}

func cacheKey(a, b *node.Node) [4]int {
	var ij = [4]int{a.Range.I, a.Range.J, b.Range.I, b.Range.J}
	sort.Ints(ij[:])
	return ij
}

func addToMergeCache(cache map[[4]int]bool, key *[4]int) {
	cache[*key] = true
}
