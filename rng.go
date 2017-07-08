package constdp

import (
	"fmt"
	"simplex/util/iter"
)

//Range
type Range struct {
	i int
	j int
}

//New Range
func NewRange(i, j int) *Range {
	return &Range{i, j}
}

//Stringer interface
func (o *Range) String() string {
	return fmt.Sprintf("Range(i=%v, j=%v)", o.i, o.j)
}

//get I
func (o *Range) I() int {
	return o.i
}

//get J
func (o *Range) J() int {
	return o.j
}

//As Array
func (o *Range) AsArray() []int {
	return []int{o.i, o.j}
}

//Size
func (o *Range) Size() int {
	return o.j - o.i
}

//Stride
func (o *Range) Stride(step ...int) []int {
	s := 1
	if len(step) > 0 {
		s = step[0]
	}
	return iter.NewGenerator(o.i, o.j+1, s).Values()
}

//Exclusive stride
func (o *Range) ExclusiveStride(step ...int) []int {
	s := 1
	if len(step) > 0 {
		s = step[0]
	}
	return iter.NewGenerator(o.i+1, o.j, s).Values()
}
