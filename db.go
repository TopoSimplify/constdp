package constdp

import (
	"simplex/struct/rtree"
)

func dbKNN(db *rtree.RTree, g rtree.BoxObj, mindist float64) []rtree.BoxObj {
	return db.KNN(g, -1, score, predicate(mindist))
}

func score(qg, item rtree.BoxObj) float64 {
	//convert qg to geometry type
	//and item to box or polygon geometry
	return qg.BBox().Distance(item.BBox())
}

func pred(_ *rtree.KObj) bool {
	return true
}

func predicate(mindist float64, _fn ...func(*rtree.KObj) bool) func(*rtree.KObj) (bool, bool) {
	fn := pred
	if len(_fn) > 0 {
		fn = _fn[0]
	}
	predicate := func(candidate *rtree.KObj) (bool, bool) {
		dist := candidate.Score()
		if dist <= mindist && fn(candidate) {
			return true, false
		}
		return false, true
	}
	return predicate
}
