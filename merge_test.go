package constdp

import (
	"time"
	"testing"
	"simplex/geom"
	"simplex/constdp/cmp"
	"simplex/constdp/opts"
	"simplex/struct/sset"
	"simplex/struct/rtree"
	"simplex/constdp/offset"
	"github.com/franela/goblin"
)

//@formatter:off
func TestMergeHull(t *testing.T) {
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
			wkt     := "LINESTRING ( 860 390, 810 360, 770 400, 760 420, 800 440, 810 470, 850 500, 810 530, 780 570, 760 530, 720 530, 710 500, 650 450 )"
			coords  := linear_coords(wkt)
			n       := len(coords) - 1
			homo    := NewConstDP(coords, constraints, options, offset.MaxOffset)

			hull    := create_hulls([][]int{{0, n}}, coords)[0]
			ha, hb  := split_at_score_selection(homo, hull)
			splits  := split_at_index(homo, hull, []int{
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

			splits  = split_at_index(homo, hull, []int{0, 5, 6, 7, 8, 12,})
			g.Assert(len(splits)).Equal(5)

			hulldb = rtree.NewRTree(8)
			boxes = make([]rtree.BoxObj, len(splits))
			for i, v := range splits {
				boxes[i] = v
			}
			hulldb.Load(boxes)

			vertex_set = sset.NewSSet(cmp.IntCmp)
			keep, rm = homo.find_mergeable_contiguous_fragments(splits, hulldb, vertex_set)
			g.Assert(len(keep)).Equal(3)
			g.Assert(len(rm)).Equal(4)
		})
	})
}
