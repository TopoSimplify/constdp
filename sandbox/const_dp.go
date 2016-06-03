package main

import (
    "fmt"
    "simplex/dp"
    . "simplex/geom"
    . "simplex/constdp"
    "simplex/constrelate"
)

func main() {
    var ln = `LINESTRING ( 520.3891360357894 542.3912033070129, 506.3024618690045 551.4232473315985, 499.8456492240652 555.3948968460392, 492.961552805167 552.5004635914114, 489.3155900796462 547.0315195031302, 494.7910190818659 540.6453203655232, 503.2430235819369 542.0539877822016, 506.3024618690045 551.4232473315985, 505.72509579166825 560.3502151427206, 505.2252456091915 568.0786679640912, 509.1982380315184 573.1744625927278, 510.22036282943196 578.5066706448671, 506.9538876603224 582.6378010057998, 500.5170101211947 582.253509809434, 492.2547493993293 573.7991034893856, 486.1060902574759 569.7640459275444, 481.13065522694603 565.0658697854723, 477.55561113833613 565.056478772063, 477.07661476306924 570.8172868053748, 478.90063032561653 576.393069064855, 488.0275462393051 588.113950554013, 494.11875860924937 590.8324257774667, 503.20704849575554 590.8039889285739, 513.3907651994501 588.2100233531045, 519.7227675184544 584.0205966373184, 528.1782382221293 579.2651445083028, 528.2675859252843 570.3930616942043, 532.7013978168333 560.3489116165815, 531.7413381584337 549.6831287580544 )`
    var fname = "/home/titus/01/dev/godev/src/simplex/constdp/sandbox/consts2.txt"

    data := NewLineStringFromWKT(ln).Coordinates()

    var opts = &dp.Options{Polyline: data, Threshold: 0}
    opts.Relations = []constrelate.Relation{
        constrelate.NewGeometryRelate(),
        constrelate.NewQuadRelate(),
        constrelate.NewMinDistanceRelate(0),
    }
    opts.Db = LoadConstDBFromFile(NewConstDB(), fname)

    var tree = NewConstDP(opts, true)

    var o = NewLineString(tree.Coordinates())
    fmt.Println(o.String())

    opts.Threshold = 20
    o = NewLineString(tree.Simplify(opts).At())

    fmt.Println(o.String())
}
