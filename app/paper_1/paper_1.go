package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"runtime"
	"simplex/lnr"
	"simplex/opts"
	"simplex/offset"
	"simplex/constdp"
	"github.com/intdxdt/geom"
	"github.com/jonas-p/go-shp"
	"github.com/intdxdt/random"
	"github.com/intdxdt/math"
)

var constraints []geom.Geometry

func init() {
	constraints = make([]geom.Geometry, 0)
}

func genRanges(parts []int32, N int32) [][2]int {
	ints := make([]int, len(parts))
	for i, v := range parts {
		ints[i] = int(v)
	}
	q := int(N)
	n := len(ints) - 1
	if ints[n] < q {
		ints = append(ints, q)
	}
	rngs := make([][2]int, 0)
	for k := 0; k < len(ints)-1; k++ {
		i, j := ints[k], ints[k+1]
		rngs = append(rngs, [2]int{i, j})
	}
	return rngs
}

func flattenLinearShp(shape *shp.PolyLine) []*Ln {
	var geoms = make([]*Ln, 0)
	var ints = genRanges(shape.Parts, shape.NumPoints)
	var partnum = 0
	var partid = random.String(10)
	var shpPoints = shape.Points
	for _, ij := range ints {
		i, j := ij[0], ij[1]
		pts := shpPoints[i:j]
		coords := make([]*geom.Point, len(pts))
		for i, pt := range pts {
			coords[i] = geom.NewPointXY(pt.X, pt.Y)
		}
		ln := NewLn(partid, geom.NewLineString(coords), partnum)
		geoms = append(geoms, ln)
		partnum += 1
	}
	return geoms
}

func readShapes(fname string) []*Ln {
	file, err := shp.Open(fname)
	if err != nil {
		log.Fatalln("Failed to open shapefile: " + fname + " (" + err.Error() + ")")
	}
	defer file.Close()

	var lns = make([]*Ln, 0)
	for file.Next() {
		_, shape := file.Shape()
		ln := shape.(*shp.PolyLine)
		tokens := flattenLinearShp(ln)
		for _, ln := range tokens {
			lns = append(lns, ln)
		}
	}

	if file.Err() != nil {
		log.Fatalln("Error while getting shapes for %s: %v", fname, file.Err())
	}

	return lns
}

func constructDpInstances(plns []*Ln, constGeoms []geom.Geometry,
	opts *opts.Opts, score lnr.ScoreFn) []*constdp.ConstDP {

	var instances = make([]*constdp.ConstDP, 0)
	for _, ln := range plns {
		dp := constdp.NewConstDP(ln.Geom.Coordinates(), constGeoms, opts, score)
		dp.Meta["oid"], dp.Meta["partnum"] = ln.Id, ln.Partnum
		instances = append(instances, dp)
	}

	return instances

}

func writeOutput(lns []*geom.LineString, fname string) {
	f, err := os.Create(fname)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	for _, ln := range lns {
		f.WriteString(ln.WKT() + "\n")
	}
	f.Sync()
}

func benchMark(NumIters int, options opts.Opts, inputShp string, output string) LnrStats {
	var stats = LnrStats{
		InputPath:   inputShp,
		Options:     options,
		ResultsPath: output,
	} //InputPath, ResultsPath, Options
	var shapes = readShapes(stats.InputPath)
	stats.NumLines = len(shapes) // stats - NumLines
	stats.NumIters = NumIters    // stats - NumIters

	for i := 0; i < NumIters; i++ {
		fmt.Println("threshold :", options.Threshold, fmt.Sprintf("iteration : #%v...", i))
		var forest = constructDpInstances(shapes, constraints, &options, offset.MaxOffset)
		var t0, t1 time.Time

		t0 = time.Now()
		if i == 0 {
			constdp.SimplifyFeatureClass(forest, &options, func(n int) {
				stats.NumDeformables = append(stats.NumDeformables, n)
			})
		} else {
			constdp.SimplifyFeatureClass(forest, &options)
		}
		t1 = time.Now()

		stats.EllapsedTime = append(stats.EllapsedTime, t1.Sub(t0).Seconds()) //stats - EllapsedTime
		if i == 0 {
			var lns, comp = extractSimpleSegs(forest)
			writeOutput(lns, output)
			stats.CompStats = comp //stats - CompStats
		}
		time.Sleep(10 * time.Second)
	}

	return stats
}

type Task struct {
	inputFname string
	thresholds []float64
	outPrefix  string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var options = opts.Opts{
		Threshold:              0,
		MinDist:                0,
		RelaxDist:              3.0,
		KeepSelfIntersects:     true,
		AvoidNewSelfIntersects: true,
		GeomRelation:           true,
		DirRelation:            true,
		DistRelation:           false,
	}

	var nb = Task{
		thresholds: []float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, math.MaxFloat64},
		outPrefix:  "nb",
		inputFname: "/media/titus/dat/data/geonb_nbrn-rrnb_shp/RoadSegmentEntity.shp",
	}
	var pitkin = Task{
		//thresholds: []float64{5, 10, 15, 20, 25, 30, 35, 40, math.MaxFloat64},
		thresholds: []float64{ 35, 40, math.MaxFloat64},
		outPrefix:  "pitkin",
		inputFname: "/media/titus/dat/tmp/pitkin_county_contours_prj.shp",
	}

	var NumIters = 10
	var tasks = []Task{nb, pitkin}
	for _, task := range tasks[1:] {
		for _, t := range task.thresholds {
			fmt.Println("Processing  Threshold : ", t)
			options.Threshold = t
			options.MinDist = t
			var outputFname = fmt.Sprintf("%v/contours_%vm.wkt", task.outPrefix, options.Threshold)

			var obj = benchMark(NumIters, options, task.inputFname, outputFname)
			var gobPath = fmt.Sprintf("%v/%vm.gob", task.outPrefix, options.Threshold)
			Save(gobPath, obj)

			var gobj = LnrStats{}
			Load(gobPath, &gobj)

			fmt.Println("Number of Iterations:", gobj.NumIters)
			fmt.Println("Number of Polylines :", gobj.NumLines)
			fmt.Println("compression ratio   :", gobj.CompStats.CompressionRatio())
			fmt.Println("compression ratio   :", gobj.CompStats.SpaceSavings(), "%")
			var m, std = gobj.AvgElapsedTime()
			fmt.Println("ellapsed time       :", m, "±", std, "secs")
			fmt.Println("deformables  :", gobj.PrintDeformables())

			runtime.Gosched()
			time.Sleep(10 * time.Second)
		}
	}

}
