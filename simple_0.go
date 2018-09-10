package constdp

import (
	"sync"
	"github.com/intdxdt/geom"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/hdb"
	"github.com/TopoSimplify/constrain"
)

func findDeformableNodes(hulls []node.Node, hulldb *hdb.Hdb) map[*node.Node]struct{} {
	return processConstSimplification(hulls, hulldb, ConcurProcs)
}

func processConstSimplification(nodeHulls []node.Node, db *hdb.Hdb, concurrency int) map[*node.Node]struct{} {
	var wg sync.WaitGroup
	var nodes = make([]*node.Node, len(nodeHulls))
	for i := range nodeHulls {
		nodes[i] = &nodeHulls[i]
	}

	var hulls = chunkTasks(nodes, concurrency)

	wg.Add(len(hulls))

	var out = make(chan *node.Node, 2*concurrency)

	var fn = func(in []*node.Node) {
		defer wg.Done()

		var selections []*node.Node

		for i := range in {
			var bln = true
			var hull = in[i]
			var self = hull.Instance.(*ConstDP)

			if hull.Range.Size() == 1 { //segment
				continue
			}

			//if hull geometry is line then points are collinear
			if _, ok := hull.Geom.(*geom.LineString); ok {
				continue
			}

			// self intersection constraint
			if self.Opts.AvoidNewSelfIntersects {
				bln = constrain.ByFeatureClassIntersection(self.Opts, hull, db, &selections)
			}

			// short circuit, if invalid, skip context validation
			if bln {
				bln = self.ValidateContextRelation(hull, &selections)
			}

			for _, o := range selections {
				out <- o
			}
			selections = selections[:0]
		}
	}

	go func() {
		for i := range hulls {
			go fn(hulls[i])
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	var results = make(map[*node.Node]struct{})
	for o := range out {
		results[o] = struct{}{}
	}
	return results
}

func chunkTasks(vals []*node.Node, concurrency int) [][]*node.Node {
	var n = len(vals)
	var chunkSize = n / concurrency
	if chunkSize == 0 {
		chunkSize = 1
	}
	var idx = 0
	var chunks = make([][]*node.Node, 0, concurrency+3)

	for idx < n {
		var stop = idx + chunkSize
		//if stop > n || (stop < n && stop+chunkSize > n) {
		if stop > n {
			stop = n
		}
		var task []*node.Node
		for _, o := range vals[idx:stop] {
			task = append(task, o)
		}
		chunks = append(chunks, task)
		idx = stop
	}
	return chunks
}



func chunkInstances(vals []*ConstDP, concurrency int) [][]*ConstDP {
	var n = len(vals)
	var chunkSize = n / concurrency
	if chunkSize == 0 {
		chunkSize = 1
	}
	var idx = 0
	var chunks = make([][]*ConstDP, 0, concurrency+3)

	for idx < n {
		var stop = idx + chunkSize
		//if stop > n || (stop < n && stop+chunkSize > n) {
		if stop > n {
			stop = n
		}
		var task []*ConstDP
		for _, o := range vals[idx:stop] {
			task = append(task, o)
		}
		chunks = append(chunks, task)
		idx = stop
	}
	return chunks
}

