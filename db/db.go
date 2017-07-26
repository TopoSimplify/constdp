package db

import (
	"simplex/struct/rtree"
)

func KNN(db *rtree.RTree, g rtree.BoxObj, minscore float64,
	score func(query, item rtree.BoxObj) float64) []rtree.BoxObj {
	return db.KNN(g, -1, score, predicate(minscore))
}

func predicate(mindist float64) func(*rtree.KObj) (bool, bool) {
	return func(candidate *rtree.KObj) (bool, bool) {
		dist := candidate.Score()
		if dist <= mindist {
			return true, false
		}
		return false, true
	}
}
