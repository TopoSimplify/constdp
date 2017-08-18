package constdp

type TestDat struct {
	pln     string
	relates ReLates
	idxs    []interface{}
	simple  string
}

type ReLates struct {
	geom bool
	dir  bool
	dist bool
}
