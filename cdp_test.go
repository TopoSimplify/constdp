package constdp

import (
	"fmt"
	"time"
	"testing"
	"strings"
	"simplex/geom"
	"simplex/struct/sset"
	"simplex/constdp/cmp"
	"simplex/constdp/opts"
	"simplex/constdp/offset"
	"github.com/franela/goblin"
)

type testData struct {
	wkt string
	rlt rLates
	res []interface{}
	out string
}

type rLates struct {
	geom bool
	dir  bool
	dist bool
}

var test_constraints_wkt = []string{
	"POLYGON (( 486.7703509316769 722.8914456521738, 475.4937795031055 708.0909456521739, 495.9325652173912 703.8622313664596, 505.7995652173912 713.7292313664595, 519.8952795031055 709.5005170807453, 529.7622795031054 725.005802795031, 510.02827950310547 729.939302795031, 489.58949378881977 729.939302795031, 486.7703509316769 722.8914456521738 ))",
	"POLYGON (( 543.8579937888197 460.00637422360245, 516.371350931677 426.8814456521739, 537.5149223602483 426.8814456521739, 537.5149223602483 440.27237422360247, 562.887208074534 434.6340885093167, 581.2116366459626 454.36808850931675, 562.1824223602483 462.8255170807453, 543.8579937888197 460.00637422360245 ))",
	"POLYGON (( 556.1038053880469 503.83863973683697, 551.5530679870658 503.6964291930563, 550.4153836368205 499.5723234234172, 552.264120705969 496.44369146024263, 555.3927526691436 491.60853297170024, 559.2324373512214 493.5994805846295, 556.3882264756082 499.71453396719784, 556.1038053880469 503.83863973683697 ))",
	"POLYGON (( 266.6254208102681 549.9800181165591, 250.83622260496082 501.7352458225647, 277.15155294713963 490.3319360076206, 299.9581725770279 551.7343734727044, 266.6254208102681 549.9800181165591 ))",
	"POLYGON (( 423.9302472540722 489.78849755284546, 423.9302472540722 512.734076819228, 444.4605023871514 512.734076819228, 444.4605023871514 489.78849755284546, 423.9302472540722 489.78849755284546 ))",
	"POLYGON (( 273.9909412231634 433.60761308851255, 274.819826756752 455.15863696181606, 292.77901331783823 455.15863696181606, 291.39753742852395 432.5024323770611, 273.9909412231634 433.60761308851255 ))",
	"POLYGON (( 219.83535798993222 660.6726799840633, 219.83535798993222 678.1837799505131, 235.5349648564045 678.1837799505131, 235.5349648564045 660.6726799840633, 219.83535798993222 660.6726799840633 ))",
	"POLYGON (( 381.05824388793627 358.75716332113427, 381.05824388793627 381.098911554191, 398.56934385438615 381.098911554191, 398.56934385438615 358.75716332113427, 381.05824388793627 358.75716332113427 ))",
}

func TestConstDP(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("const dp", func() {
		g.It("should test constraint dp algorithm", func() {
			g.Timeout(1 * time.Hour)
			options := &opts.Opts{
				Threshold:              50.0,
				MinDist:                20.0,
				RelaxDist:              30.0,
				KeepSelfIntersects:     true,
				AvoidNewSelfIntersects: true,
				GeomRelation:           false,
				DistRelation:           false,
				DirRelation:            false,
			}
			for _, td := range test_data {
				constraints := make([]geom.Geometry, 0)
				for _, wkt := range test_constraints_wkt {
					g := geom.NewPolygonFromWKT(wkt)
					constraints = append(constraints, g)
				}

				options.GeomRelation = td.rlt.geom
				options.DirRelation  = td.rlt.dir
				options.DistRelation = td.rlt.dist

				coords := geom.NewLineStringFromWKT(td.wkt).Coordinates()
				homo := NewConstDP(coords, constraints, options, offset.MaxOffset)
				homo.Simplify(options)
				ptset := sset.NewSSet(cmp.IntCmp)
				for _, o := range homo.Simple {
					ptset.Add(o.Range.I())
					ptset.Add(o.Range.J())
				}
				fmt.Println(ptset.Values())
				fmt.Println(td.res)
				g.Assert(ptset.Values()).Equal(td.res)
				fmt.Println(strings.Repeat("--", 80))
			}
		})
	})

}

func TestConstSED(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("const sed", func() {
		g.It("should test constraint sed algorithm", func() {

			options := &opts.Opts{
				Threshold:              1.0,
				MinDist:                20.0,
				RelaxDist:              30.0,
				KeepSelfIntersects:     true,
				AvoidNewSelfIntersects: true,
				GeomRelation:           true,
				DistRelation:           false,
				DirRelation:            true,
			}

			constraints := make([]geom.Geometry, 0)
			for _, wkt := range test_constraints_wkt {
				g := geom.NewPolygonFromWKT(wkt)
				constraints = append(constraints, g)
			}

			coords := []*geom.Point{{3.0, 1.6, 0.0}, {3.0, 2.0, 1.0}, {2.4, 2.8, 3.0}, {0.5, 3.0, 4.5}, {1.2, 3.2, 5.0}, {1.4, 2.6, 6.0}, {2.0, 3.5, 10.0}}
			homo := NewConstDP(coords, constraints, options, offset.MaxSEDOffset)
			homo.Simplify(options)
		})
	})
}
