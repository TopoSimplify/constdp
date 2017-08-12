package constdp

import (
	"time"
	"testing"
	"simplex/geom"
	"simplex/constdp/ln"
	"simplex/constdp/cmp"
	"simplex/constdp/rng"
	"simplex/constdp/opts"
	"simplex/struct/sset"
	"simplex/struct/rtree"
	"simplex/constdp/offset"
	"github.com/franela/goblin"
)

func create_hulls(indxs [][]int, coords []*geom.Point) []*HullNode {
	ranges := make([]*rng.Range, 0)
	for _, idx := range indxs {
		ranges = append(ranges, rng.NewRange(idx[0], idx[1]))
	}

	n := len(coords) - 1
	pln := ln.NewPolyline(coords)
	hulls := make([]*HullNode, 0)

	for _, r := range ranges {
		pr := rng.NewRange(0, n)
		h := NewHullNode(pln, r, pr)
		hulls = append(hulls, h)
	}
	return hulls
}

func TestSplitHull(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("test split hull", func() {
		g.It("should test split", func() {
			g.Timeout(1 * time.Hour)
			options := &opts.Opts{
				Threshold:              50.0,
				MinDist:                20.0,
				RelaxDist:              30.0,
				KeepSelfIntersects:     true,
				AvoidNewSelfIntersects: true,
				GeomRelation:           true,
				DistRelation:           false,
				DirRelation:            false,
			}

			// self.relates = relations(self)
			constraints := make([]geom.Geometry, 0)
			wkt := "LINESTRING ( 860 390, 810 360, 770 400, 760 420, 800 440, 810 470, 850 500, 810 530, 780 570, 760 530, 720 530, 710 500, 650 450 )"
			coords := geom.NewLineStringFromWKT(wkt).Coordinates()
			n := len(coords) - 1
			homo := NewConstDP(coords, constraints, options, offset.MaxOffset)

			hull := create_hulls([][]int{{0, n}}, coords)[0]
			ha, hb := split_at_score_selection(homo, hull)

			g.Assert(ha.Range.AsSlice()).Equal([]int{0, 8})
			g.Assert(hb.Range.AsSlice()).Equal([]int{8, len(coords) - 1})

			splits := split_at_index(homo, ha, []int{3, 6})
			g.Assert(len(splits)).Equal(3)
			g.Assert(splits[0].Range.AsSlice()).Equal([]int{0, 3})
			g.Assert(splits[1].Range.AsSlice()).Equal([]int{3, 6})
			g.Assert(splits[2].Range.AsSlice()).Equal([]int{6, 8})

			splits = split_at_index(homo, hull, []int{
				ha.Range.I(), ha.Range.J(),
				hb.Range.I(), hb.Range.J(),
			})

			g.Assert(len(splits)).Equal(2)
			splits = split_at_index(homo, hull, []int{
				ha.Range.I(), ha.Range.J(), hb.Range.I(),
				hb.Range.I() - 1, hb.Range.J(),
			})
			g.Assert(len(splits)).Equal(3)

			hulldb := rtree.NewRTree(8)
			boxes := make([]rtree.BoxObj, len(splits))
			for i, v := range splits {
				boxes[i] = v
			}
			hulldb.Load(boxes)

			vertex_set := sset.NewSSet(cmp.IntCmp)
			keep, rm := homo.find_mergeable_contiguous_fragments(splits, hulldb, vertex_set)
			g.Assert(len(keep)).Equal(2)
			g.Assert(len(rm)).Equal(2)

		})
	})
}
