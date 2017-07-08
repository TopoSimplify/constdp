package main

//import (
//    . "simplex/geom"
//    "fmt"
//)
//
//var data = []*Point{{0.5, 1.0}, {1.5, 2.0}, {2.5, 1.5}, {3.5, 2.5}, {4.0, 1.5}, {3.0, 1.0} }
//
//func main() {
//    //var gstr = 'LINESTRING (520.3891360357894 542.3912033070129, 506.3024618690045 551.4232473315985, 499.8456492240652 555.3948968460392, 492.961552805167 552.5004635914114, 489.3155900796462 547.0315195031302, 494.7910190818659 540.6453203655232, 503.2430235819369 542.0539877822016, 506.3024618690045 551.4232473315985, 505.72509579166825 560.3502151427206, 505.2252456091915 568.0786679640912)'
//    //var g = geom.readwkt(gstr)
//    var opts = &dp.Options{Polyline: data, Threshold: 0}
//    var tree = NewConstDP(opts, true)
//    var tree_str = tree.Print()
//    fmt.Println(tree_str)
//
//    var o = NewLineString(tree.Coordinates())
//    fmt.Println(o.String())
//    opts = &dp.Options{Threshold:0.6}
//    o = NewLineString(tree.Simplify(opts).At())
//
//    fmt.Println (o.String())
//    fmt.Println(o.Coordinates())
//}