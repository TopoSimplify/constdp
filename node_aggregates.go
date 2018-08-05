package constdp

import (
	"sort"
	"github.com/intdxdt/iter"
		"github.com/TopoSimplify/knn"
	"github.com/TopoSimplify/hdb"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/merge"
	"github.com/TopoSimplify/common"
)

//Merge segment fragments where possible
func (self *ConstDP) AggregateSimpleSegments(
	id *iter.Igen, db *hdb.Hdb, constVertexSet []int,
	scoreRelation func(float64) bool, validateMerge func(*node.Node, *hdb.Hdb) bool) {

	var fragmentSize = 1
	var neighbours []*node.Node
	var cache = make(map[[4]int]bool)
	//var objects = common.NodesFromObjects(db.All())
	var objects = db.All()
	sort.Sort(node.NodePtrs(objects))

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
		var mergePrev, mergeNxt node.Node
		var mergeprevBln, mergenxtBln bool

		if withPrev && prev != nil {
			key = cacheKey(prev, hull)
			if !cache[key] {
				addToMergeCache(cache, &key)
				mergeprevBln, mergePrev = merge.ContiguousFragmentsAtThreshold(
					id, self.Score, prev, hull, scoreRelation, common.Geometry,
				)
			}
		}

		if withNext && nxt != nil {
			key = cacheKey(hull, nxt)
			if !cache[key] {
				addToMergeCache(cache, &key)
				mergenxtBln, mergeNxt = merge.ContiguousFragmentsAtThreshold(
					id, self.Score, hull, nxt, scoreRelation, common.Geometry,
				)
			}
		}

		var insertList []node.Node
		var merged bool
		//nxt, prev
		if !merged && mergenxtBln {
			db.Remove(nxt)
			if validateMerge(&mergeNxt, db) {
				if objects[0] == nxt {
					objects = objects[1:]
				}
				insertList = append(insertList, mergeNxt)
				merged = true
			} else {
				merged = false
				insertList = append(insertList, *nxt)
			}
		}

		if !merged && mergeprevBln {
			db.Remove(prev)
			//prev cannot exist since moving from left --- right
			if validateMerge(&mergePrev, db) {
				insertList = append(insertList, mergePrev)
				merged = true
			} else {
				merged = false
				insertList = append(insertList, *prev)
			}
		}

		if !merged {
			insertList = append(insertList, *hull)
		}

		db.Load(insertList)
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
