package main

import (
	"simplex/opts"
	"github.com/intdxdt/math"
	"github.com/montanaflynn/stats"
	"log"
	"bytes"
	"fmt"
)

type CompStats struct {
	Uncompressed float64
	Compressed   float64
}

func (o CompStats) CompressionRatio() float64 {
	return math.Round(o.Uncompressed/o.Compressed, 1)
}

func (o CompStats) SpaceSavings() float64 {
	var r = o.Compressed / o.Uncompressed
	return math.Round((1.0-r)*100.0, 1)
}

type LnrStats struct {
	NumIters       int
	InputPath      string
	ResultsPath    string
	EllapsedTime   []float64
	NumLines       int
	SimpleSize     int
	NumDeformables []int
	Options        opts.Opts
	CompStats      CompStats
}

func (o LnrStats) PrintDeformables() string{
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, v := range o.NumDeformables {
		buf.WriteString(fmt.Sprintf("%v", v))
		if i < len(o.NumDeformables)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("]")
	return buf.String()
}
func (o LnrStats) AvgElapsedTime() (float64, float64) {
	var s float64
	for _, v := range o.EllapsedTime {
		s += v
	}
	var dat = stats.Float64Data(o.EllapsedTime)
	var m, err = dat.Mean()
	if err != nil {
		log.Fatal(err)
	}
	std, err := dat.StandardDeviationSample()
	if err != nil {
		log.Fatal(err)
	}
	return math.Round(m), math.Round(std)
}
