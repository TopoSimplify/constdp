package constdp

import (
	"simplex/geom"
	"simplex/util/math"
	"simplex/struct/sset"
)

type HullSideTangent struct {
	aseg *Seg
	bseg *Seg
	rtan *Seg
	ltan *Seg
	hull *sset.SSet
	side *HullCollapseSidedness
}

type HullCollapseSidedness struct {
	hullset *sset.SSet
	keys    *sset.SSet
}

func (self *HullCollapseSidedness) index(i int) int {
	return int(self.PtAt(i)[2])
}

func (self *HullCollapseSidedness) PtAt(i int) *geom.Point {
	hpt := self.hullset.Get(i).(*HPt)
	return hpt.Point
}

//Hull collapse sidedness measures the correctness of a contiguous
//hull collaps - goal is to prevent line flips against preceeding hull
func NewHullCollapseSidedness(hull_ptset *sset.SSet) *HullCollapseSidedness {
	indx_set := sset.NewSSet(int_cmp)
	for _, o := range hull_ptset.Values() {
		p := o.(*HPt).Point
		indx_set.Add(int(p[2]))
	}
	return &HullCollapseSidedness{hullset: hull_ptset, keys: indx_set}
}

func (self *HullCollapseSidedness) key_knn(key int) (int, int) {
	n := self.hullset.Size()
	idx := self.index_of_key(key)
	idx_prev, idx_next := idx-1, idx+1

	if idx_prev < 0 {
		idx_prev = -1
	}

	if idx_next > (n - 1) {
		idx_next = 0
	}

	prev, next := self.hullset.Get(idx_prev).(*HPt), self.hullset.Get(idx_next).(*HPt)

	return int(prev.Point[2]), int(next.Point[2])
}

func (self *HullCollapseSidedness) index_of_key(key int) int {
	return self.keys.IndexOf(key)
}

func (self *HullCollapseSidedness) _side_tangents(pt *geom.Point, hull *sset.SSet) *HullSideTangent {
	n := hull.Size()
	hpts := make([]*HPt, n)
	coords := make([]*geom.Point, hull.Size())
	for i, o := range hull.Values() {
		hpts[i] = o.(*HPt)
		coords[i] = hpts[i].Point
	}

	pt_key := int(pt[2])
	r, l := geom.TangentPointToPoly(pt, coords)

	if r == 0 && l == 0 {
		hpt := hpts[0] //hull A and B starts at pt
		if math.FloatEqual(pt[0], hpt.Point[0]) && math.FloatEqual(pt[1], hpt.Point[1]) {
			r, l = 1, n-1
		} else { //hull A ends and B starts at pt
			r, l = 0, n-2
		}

		rpt, lpt := hpts[r], hpts[l]
		hseg := NewSeg(pt, rpt.Point, 0, -1)
		s := hseg.SideOf(lpt.Point)
		if s.IsRight() { // if lpt is on right then swap
			r, l = l, r
		}
	}

	rtan := NewSeg(pt, hpts[r].Point, 0, -1)
	ltan := NewSeg(pt, hpts[l].Point, 0, -1)

	prv, nxt := self.key_knn(pt_key)

	fk := self.index_of_key(pt_key)
	nk := self.index_of_key(nxt)
	pk := self.index_of_key(prv)

	if pt_key == self.keys.Get(-1).(int) {
		nk, pk = pk, nk
	}

	aseg := NewSeg(self.PtAt(fk), self.PtAt(nk), 0, -1)
	bseg := NewSeg(self.PtAt(fk), self.PtAt(pk), 0, -1)

	return &HullSideTangent{
		aseg: aseg, bseg: bseg,
		rtan: rtan, ltan: ltan,
		side: self, hull: hull,
	}
}

//checks if hullsidedness of self line(i---j)
//is consistent with original hullset and hull
func (self *HullCollapseSidedness) IsValid(hull *sset.SSet) bool {
	sa, sb := self.hullset.First().(*HPt).Point, self.hullset.Last().(*HPt).Point
	ha, hb := hull.First().(*HPt).Point, hull.Last().(*HPt).Point
	var pt *geom.Point
	var feq = math.FloatEqual
	if (feq(sa[0], ha[0]) && feq(sa[1], ha[1])) || (feq(sa[0], hb[0]) && feq(sa[1], hb[1])) {
		pt = sa
	} else {
		pt = sb
	}

	tangents := self._side_tangents(pt, hull)
	ltan := tangents.ltan
	rtan := tangents.rtan

	aseg := tangents.aseg
	bseg := tangents.bseg

	side_a := rtan.SideOf(aseg.B)
	side_b := ltan.SideOf(bseg.B)
	if side_a.IsRight() && side_b.IsLeft() {
		return false
	}
	return true
}

