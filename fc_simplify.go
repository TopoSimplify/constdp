package constdp

import (
	"sync"
	"github.com/intdxdt/iter"
	"github.com/TopoSimplify/lnr"
	"github.com/TopoSimplify/hdb"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/opts"
)

//Simplify a feature class of linear geometries
//optional callback for the number of deformables
func SimplifyFeatureClass(id *iter.Igen, selfs []*ConstDP, opts *opts.Opts, callback ... func(n int)) {
	var deformableCallback = func(_ int) {}
	if len(callback) > 0 {
		deformableCallback = callback[0]
	}

	var junctions = make(map[int][]int)

	if opts.PlanarSelf {
		instances := make([]*lnr.FC, len(selfs))
		for i, sf := range selfs {
			instances[i] = lnr.NewFC(sf.Coordinates(), sf.Id())
		}
		junctions = lnr.FCPlanarSelfIntersection(instances)
	}

	SimplifyDPs(id, selfs, junctions)

	var constrained = opts.AvoidNewSelfIntersects ||
		opts.PlanarSelf ||
		opts.GeomRelation ||
		opts.DirRelation ||
		opts.DistRelation

	if constrained {
		var selections map[*node.Node]struct{}
		var hulldb = hdb.NewHdb()
		var deformables []node.Node

		for _, self := range selfs {
			for i := range self.Hulls {
				deformables = append(deformables, self.Hulls[i])
			}
			node.Clear(&self.Hulls) // empty deque, this is for future splits
		}
		hulldb.Load(deformables)

		for len(deformables) > 0 {
			deformableCallback(len(deformables))
			// 0. find deformable node
			selections = findDeformableNodes(deformables, hulldb)
			// 1. deform selected nodes
			if len(selections) > 0 {
				deformables = deformNodes(id, selections)
				// 2. remove selected nodes from db
				cleanUpDB(hulldb, selections)
				// 3. add new deformations to db
				hulldb.Load(deformables)
				// 4. repeat until there are no deformables
			} else {
				deformables = deformables[:0]
			}
		}
		groupHullsByFC(hulldb)
	}
}

func SimplifyDPs(id *iter.Igen, selfs []*ConstDP, junctions map[int][]int) {
	var wg sync.WaitGroup
	wg.Add(ConcurProcs)

	var stream = make(chan *ConstDP)
	var out = make(chan *ConstDP, 2*ConcurProcs)

	go func() {
		for s := range selfs {
			stream <- selfs[s]
		}
		close(stream)
	}()

	//assume only one worker reading from input chan
	var fn = func(idx int) {
		defer wg.Done()
		for self := range stream {
			var constVerts []int
			if v, ok := junctions[self.Id()]; ok {
				constVerts = v
			} else {
				constVerts = make([]int, 0)
			}

			self.Simplify(id, constVerts)
			out <- self
		}
	}

	//now expand one worker into clones of workers
	go func() {
		for i := 0; i < ConcurProcs; i++ {
			go fn(i)
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	for range out {
	}
}
