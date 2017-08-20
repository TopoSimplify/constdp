package opts

//Opts
type Opts struct {
	Threshold              float64
	MinDist                float64
	RelaxDist              float64
	KeepSelfIntersects     bool
	AvoidNewSelfIntersects bool
	GeomRelation           bool
	DistRelation           bool
	DirRelation            bool
}

