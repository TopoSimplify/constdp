package constdp

type HullNodes []*HullNode

func (s HullNodes) Len() int {
	return len(s)
}

func (s HullNodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s HullNodes) Less(i, j int) bool {
	return s[i].Range.i < s[j].Range.i
}
