package constdp

import (
	"sync"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/hdb"
	"github.com/intdxdt/geom"
	"github.com/TopoSimplify/constrain"
)

func findDeformableNodes(hulls []node.Node, hulldb *hdb.Hdb) map[int]*node.Node {
	return processConstSimplification(hulls, hulldb, ConcurProcs)
}

//process
func processConstSimplification(nodeHulls []node.Node, db *hdb.Hdb, concurrency int) map[int]*node.Node {
	var wg sync.WaitGroup
	//var cache = newCacheMap(len(nodeHulls))
	var nodes = make([]*node.Node, len(nodeHulls))
	for i := range nodeHulls {
		nodes[i] = &nodeHulls[i]
	}

	var hulls = chunkTasks(nodes, concurrency)
	//set up number of of clones to wait for
	wg.Add(len(hulls))

	var out = make(chan []*node.Node, 2*concurrency)

	//assume only one worker
	var fn = func(in []*node.Node) {
		defer wg.Done()

		//perform fn here...
		for i := range in {
			var hull = in[i]
			var selections []*node.Node
			var self = hull.Instance.(*ConstDP)

			//if hull is segment
			if hull.Range.Size() == 1 {
				continue
			}

			//if hull geometry is line then points are collinear
			if _, ok := hull.Geom.(*geom.LineString); ok {
				continue
			}

			// self intersection constraint
			if self.Opts.AvoidNewSelfIntersects {
				constrain.ByFeatureClassIntersection(self.Opts, hull, db, &selections)
			}

			// context_geom geometry constraint
			self.ValidateContextRelation(hull, &selections)
			if len(selections) > 0 {
				//for _, o := range selections {
				//	var key = cacheKey(hull, o)
				//	if !cache.HasKey(&key) {
				//		cache.Set(&key)
				//	}
				//}
				//out <- worker(in[i], db)
				out <- selections
			}
		}

	}

	//now expand one worker into clones of workers
	go func() {
		for i := range hulls {
			go fn(hulls[i])
		}
	}()

	//wait for all the clones to be done
	//in a new go routine
	go func() {
		wg.Wait()
		close(out)
	}()

	var results = make(map[int]*node.Node)
	for selections := range out {
		for _, n := range selections {
			results[n.Id] = n
		}
	}
	//return out chan to whoever want to read from it
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
		chunks = append(chunks, vals[idx:stop])
		idx = stop
	}
	return chunks
}
