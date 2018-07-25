package constdp

import (
	"sort"
	"github.com/TopoSimplify/dp"
	"github.com/TopoSimplify/knn"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/merge"
	"github.com/intdxdt/rtree"
	"github.com/intdxdt/iter"
)

//Merge segment fragments where possible
func (self *ConstDP) AggregateSimpleSegments(
	nodeDB *rtree.RTree, constVertexSet []int,
	scoreRelation func(float64) bool,
	validateMerge func(*node.Node, *rtree.RTree) bool,
) {

	var fragmentSize = 1
	var neighbours []*rtree.Obj
	var cache = make(map[[4]int]bool)
	//var objects = common.NodesFromObjects(nodeDB.All())
	var objects = nodeDB.All()
	sort.Sort(node.NodeObjects(objects))

	for len(objects) != 0 {
		obj := objects[0]
		hull := obj.Object.(*node.Node)
		objects = objects[1:]
		if hull.Range.Size() != fragmentSize {
			continue
		}

		//make sure hull index is not part of vertex with degree > 2
		if iter.SortedSearchInts(constVertexSet, hull.Range.I) && iter.SortedSearchInts(constVertexSet, hull.Range.J) {
			continue
		}
		var withNext, withPrev = !iter.SortedSearchInts(constVertexSet, hull.Range.J),
			!iter.SortedSearchInts(constVertexSet, hull.Range.I)

		nodeDB.RemoveObj(obj)

		// find context neighbours
		//neighbours = common.NodesFromObjects(knn.FindNodeNeighbours(nodeDB, hull, knn.EpsilonDist))
		neighbours = knn.FindNodeNeighbours(nodeDB, hull, knn.EpsilonDist)

		// find context neighbours
		var prev, nxt = node.Neighbours(obj, neighbours)

		// find mergeable neighbours contiguous
		var key [4]int
		var mergePrev, mergeNxt *rtree.Obj

		if withPrev && prev != nil {
			key = cacheKey(prev.Object.(*node.Node), hull)
			if !cache[key] {
				addToMergeCache(cache, &key)
				o := merge.ContiguousFragmentsAtThreshold(
					self.Score, prev.Object.(*node.Node), hull, scoreRelation, dp.NodeGeometry,
				)
				mergePrev = rtree.Object(prev.Id, o.Bounds(), o)
			}
		}

		if withNext && nxt != nil {
			key = cacheKey(hull, nxt.Object.(*node.Node))
			if !cache[key] {
				addToMergeCache(cache, &key)
				o := merge.ContiguousFragmentsAtThreshold(
					self.Score, hull, nxt.Object.(*node.Node), scoreRelation, dp.NodeGeometry,
				)
				mergeNxt = rtree.Object(nxt.Id, o.Bounds(), o)
			}
		}

		var merged bool
		//nxt, prev
		if !merged && mergeNxt != nil {
			nodeDB.RemoveObj(nxt)
			if validateMerge(mergeNxt.Object.(*node.Node), nodeDB) {
				if objects[0] == nxt {
					objects = objects[1:]
				}
				nodeDB.Insert(mergeNxt)
				merged = true
			} else {
				merged = false
				nodeDB.Insert(nxt)
			}
		}

		if !merged && mergePrev != nil {
			nodeDB.RemoveObj(prev)
			//prev cannot exist since moving from left --- right
			if validateMerge(mergePrev.Object.(*node.Node), nodeDB) {
				nodeDB.Insert(mergePrev)
				merged = true
			} else {
				merged = false
				nodeDB.Insert(prev)
			}
		}

		if !merged {
			nodeDB.Insert(obj)
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
