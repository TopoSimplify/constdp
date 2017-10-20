package main

import (
	"sort"
	"fmt"
	"strings"
	"simplex/pln"
	"simplex/rng"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/sset"
	"simplex/node"
)

type HullNodes []*node.Node

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
	Pln    *pln.Polyline
	Range  *rng.Range
	PRange *rng.Range
	Geom   geom.Geometry
	PtSet  *sset.SSet
}

//sort hulls
func sort_hulls(hulls []*node.Node) []*node.Node {
	sort.Sort(HullNodes(hulls))
	return hulls
}

//reverse sort hulls
func sort_reverse(hulls []*node.Node) []*node.Node {
	sort.Sort(sort.Reverse(HullNodes(hulls)))
	return hulls
}

func main() {

	example := []int{1, 25, 3, 5, 4}
	sort.Sort(sort.Reverse(sort.IntSlice(example)))
	fmt.Println(example)

	a := &node.Node{Range: rng.NewRange(90, 97)}
	b := &node.Node{Range: rng.NewRange(72, 76)}
	c := &node.Node{Range: rng.NewRange(55, 61)}
	hulls := []*node.Node{a, c, b}

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
