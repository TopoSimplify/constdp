package constdp

func ExceptHull(hulls []*HullNode, hull *HullNode) []*HullNode {
	hdict := make(map[[2]int]*HullNode, 0)
	for _, h := range hulls {
		hdict[h.Range.AsArray()] = h
	}

	delete(hdict, hull.Range.AsArray())

	hulls = make([]*HullNode, 0)
	for _, v := range hdict {
		hulls = append(hulls, v)
	}
	return hulls
}
