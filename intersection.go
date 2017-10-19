package constdp

import (
	"github.com/intdxdt/geom"
	"simplex/constdp/ln"
	"simplex/constdp/seg"
	"github.com/intdxdt/sset"
	"github.com/intdxdt/rtree"
	"simplex/constdp/cmp"
	"simplex/constdp/ctx"
)

type kvCount struct {
	count   int
	indxset *sset.SSet
}

type selfInter struct {
	keyset *sset.SSet
	point  *geom.Point
}

func update_kv_count(dict map[[2]float64]*kvCount, o *geom.Point, index int) {
	k := [2]float64{o[0], o[1]}
	v, ok := dict[k]
	if !ok {
		v = &kvCount{
			count:   0,
			indxset: sset.NewSSet(cmp.IntCmp, 8),
		}
		dict[k] = v
	}
	v.indxset.Add(index)
	v.count += 1
}

func linear_ftclass_self_intersection(ftcls []*ConstDP) map[string]*sset.SSet {
	var coord_dict = make(map[[2]float64]map[string]int, 0)
	for _, self := range ftcls {
		pln := self.Pln
		n := self.Pln.Len()
		for i := 0; i < n; i++ {
			var dat map[string]int
			var ok bool
			var pt = pln.Coords[i]
			var key = [2]float64{pt.X(), pt.Y()}

			if dat, ok = coord_dict[key]; !ok {
				dat = make(map[string]int, 0)
			}

			if _, ok = dat[self.Id]; !ok {
				dat[self.Id] = i
			}
			coord_dict[key] = dat
		}
	}

	var fc_junctions = make(map[string]*sset.SSet, 0)
	for _, o := range coord_dict {
		if len(o) > 1 {
			for sid, idx := range o {
				var ok bool
				var s *sset.SSet

				if s, ok = fc_junctions [sid]; !ok {
					s = sset.NewSSet(cmp.IntCmp)
				}
				s.Add(idx)
				fc_junctions[sid] = s
			}
		}
	}
	return fc_junctions
}

func linear_self_intersection(pln *ln.Polyline) []*ctx.CtxGeom {
	var tree = *rtree.NewRTree(8)
	var dict = make(map[[2]float64]*kvCount)

	var data = make([]rtree.BoxObj, 0)
	for _, s := range pln.Segments() {
		data = append(data, s)
	}
	tree.Load(data)

	self_intersects := make(map[string]*selfInter, 0)

	for _, d := range data {
		s := d.(*seg.Seg)
		res := tree.Search(s.BBox())
		update_kv_count(dict, s.A, s.I)
		update_kv_count(dict, s.B, s.J)

		for _, node := range res {
			other_seg := node.GetItem().(*seg.Seg)

			if s == other_seg {
				continue
			}
			seg_g, other_seg_g := s.Segment, other_seg.Segment
			intersects := seg_g.Intersection(other_seg_g)

			if len(intersects) == 0 {
				continue
			}

			for _, pt := range intersects {
				if s.A.Equals2D(pt) || s.B.Equals2D(pt) {
					continue
				}
				skey := sset.NewSSet(cmp.IntCmp, 8).Extend(s.I, s.J, other_seg.I, other_seg.J)

				k := skey.String()
				v, ok := self_intersects[k]
				if !ok {
					v = &selfInter{
						keyset: skey,
						point:  pt,
					}
				}
				self_intersects[k] = v
			}
		}
	}

	results := make([]*ctx.CtxGeom, 0)
	for _, val := range self_intersects {
		cg := ctx.NewCtxGeom(val.point, 0, -1).AsSelfNonVertex()
		cg.Meta.SelfNonVertices = val.keyset
		results = append(results, cg)
	}

	for k, v := range dict {
		if v.count > 2 {
			i := v.indxset.First().(int)
			cg := ctx.NewCtxGeom(geom.NewPoint(k[:]), i, i).AsSelfVertex()
			cg.Meta.SelfVertices = v.indxset
			results = append(results, cg)
		}
	}

	return sort_context_geoms(results)
}
