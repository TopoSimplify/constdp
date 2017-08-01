package db

import (
	"simplex/struct/rtree"
)

func KNN(
	db *rtree.RTree, g rtree.BoxObj, minscore float64,
	score func(rtree.BoxObj, rtree.BoxObj) float64,
	predicate ... func(*rtree.KObj) (bool, bool),
) []rtree.BoxObj {
	var pred func(*rtree.KObj) (bool, bool)
	if len(predicate) > 0 {
		pred = predicate[0]
	} else {
		pred = default_predicate(minscore)
	}

	return db.KNN(g, -1, score, pred)
}

func default_predicate(mindist float64) func(*rtree.KObj) (bool, bool) {
	return func(candidate *rtree.KObj) (bool, bool) {
		dist := candidate.Score()
		if dist <= mindist {
			return true, false
		}
		return false, true
	}
}
