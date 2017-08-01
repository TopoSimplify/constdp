package constdp

import "sort"

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


func sort_hulls(hulls []*HullNode)[]*HullNode{
	sort.Sort(HullNodes(hulls))
	return hulls
}