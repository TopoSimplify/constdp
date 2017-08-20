package main

import (
	"simplex/geom"
	"simplex/constdp/ln"
	"simplex/constdp/rng"
	"simplex/struct/sset"
	"sort"
	"fmt"
	"strings"
)

type HullNodes []*HullNode

func (s HullNodes) Len() int {
	return len(s)
}

func (s HullNodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s HullNodes) Less(i, j int) bool {
	return s[i].Range.I() < s[j].Range.I()
}

//hull node
type HullNode struct {
	Pln    *ln.Polyline
	Range  *rng.Range
	PRange *rng.Range
	Geom   geom.Geometry
	PtSet  *sset.SSet
}

//sort hulls
func sort_hulls(hulls []*HullNode) []*HullNode {
	sort.Sort(HullNodes(hulls))
	return hulls
}

//reverse sort hulls
func sort_reverse(hulls []*HullNode) []*HullNode {
	sort.Sort(sort.Reverse(HullNodes(hulls)))
	return hulls
}

func main() {

	example := []int{1, 25, 3, 5, 4}
	sort.Sort(sort.Reverse(sort.IntSlice(example)))
	fmt.Println(example)

	a := &HullNode{Range: rng.NewRange(90, 97)}
	b := &HullNode{Range: rng.NewRange(72, 76)}
	c := &HullNode{Range: rng.NewRange(55, 61)}
	hulls := []*HullNode{a, c, b}

	hulls = sort_hulls(hulls)
	for _, o := range hulls {
		fmt.Println(o.Range)
	}

	hulls = sort_reverse(hulls)
	fmt.Println(strings.Repeat("-", 50))
	for _, o := range hulls {
		fmt.Println(o.Range)
	}
}
