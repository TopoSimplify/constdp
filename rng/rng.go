package rng

import (
	"fmt"
	"github.com/intdxdt/cmp"
	"github.com/intdxdt/iter"
	"github.com/intdxdt/sset"
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

//clone Range
func (o *Range) Clone() *Range {
	return NewRange(o.i, o.j)
}

//compare equality of two ranges
func (o *Range) Equals(r *Range) bool {
	return  (o.i == r.i) && (o.j == r.j)
}

//as segment
func (o *Range) Contains(idx int) bool {
	return idx == o.i || idx == o.j
}

//As Array
func (o *Range) AsArray() [2]int {
	return [2]int{o.i, o.j}
}

//As Slice
func (o *Range) AsSlice() []int {
	ar := o.AsArray()
	return ar[:]
}

//Size
func (o *Range) Size() int {
	return o.j - o.i
}

//Stride
func (o *Range) Stride(step ...int) []int {
	var s = 1
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

//Split Range at indices
func (o *Range) Split(idxs []int) []*Range {
	var idxset = sset.NewSSet(cmp.Int)
	for _, v := range idxs {
		idxset.Add(v)
	}
	idxs = make([]int, 0)
	for _, o := range idxset.Values() {
		idxs = append(idxs, o.(int))
	}

	i, j := o.I(), o.J()
	sub := make([]*Range, 0)
	for _, idx := range idxs {
		if i < idx && idx < j {
			s := i
			if len(sub) > 0 {
				s = sub[len(sub)-1].J()
			}
			sub = append(sub, NewRange(s, idx))
		}
	}
	//close rest of split
	if len(sub) > 0 {
		e := sub[len(sub)-1].J()
		if e < j {
			sub = append(sub, NewRange(e, j))
		}
	}
	return sub
}
