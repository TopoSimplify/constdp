package constdp

import (
	"sort"
	"simplex/struct/rtree"
	"simplex/struct/deque"
	"simplex/struct/sset"
	"simplex/constdp/cmp"
)

type nodes []*HullNode

func (s nodes) Len() int {
	return len(s)
}

func (s nodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s nodes) Less(i, j int) bool {
	return s[i].Range.I() < s[j].Range.I()
}

type HullNodes struct {
	list nodes
}

//Get at index
func (self *HullNodes) Get(index int) *HullNode {
	return self.list[index]
}

//Sort a slice of hull nodes
func (self *HullNodes) Sort() *HullNodes {
	sort.Sort(self.list)
	return self
}

//Reverse Sort a slice of hull nodes
func (self *HullNodes) Reverse() *HullNodes {
	sort.Sort(sort.Reverse(self.list))
	return self
}

func (self *HullNodes) Push(v *HullNode) *HullNodes {
	self.list = append(self.list, v)
	return self
}

func (self *HullNodes) Extend(vals ...*HullNode) *HullNodes {
	for _, h := range vals {
		self.list = append(self.list, h)
	}
	return self
}

func (self *HullNodes) Pop() *HullNodes {
	if !self.IsEmpty() {
		n := len(self.list) - 1
		self.list[n] = nil
		self.list = self.list[:n]
	}
	return self
}

func (self *HullNodes) Len() int {
	return len(self.list)
}

func (self *HullNodes) IsEmpty() bool {
	return self.Len() == 0
}

func (self *HullNodes) Empty() *HullNodes {
	for i := range self.list {
		self.list[i] = nil
	}
	self.list = self.list[:0]
	return self
}

func (self *HullNodes) AsDeque() *deque.Deque {
	queue := deque.NewDeque()
	for _, h := range self.list {
		queue.Append(h)
	}
	return queue
}

func (self *HullNodes) AsPointSet() *sset.SSet {
	var set = sset.NewSSet(cmp.IntCmp)
	for _, o := range self.list {
		set.Extend(o.Range.I(), o.Range.J())
	}
	return set
}

func NewHullNodes(size ...int) *HullNodes {
	var n = 0
	if len(size) > 0 {
		n = size[0]
	}
	return &HullNodes{list: make(nodes, n)}
}

//HullNodes from Rtree boxes
func NewHullNodesFromBoxes(iter []rtree.BoxObj) *HullNodes {
	var self = NewHullNodes(len(iter))
	for i, h := range iter {
		self.list [i] = h.(*HullNode)
	}
	return self
}

//HullNodes from Rtree nodes
func NewHullNodesFromNodes(iter []*rtree.Node) *HullNodes {
	var self = NewHullNodes(len(iter))
	for i, h := range iter {
		self.list[i] = h.GetItem().(*HullNode)
	}
	return self
}
