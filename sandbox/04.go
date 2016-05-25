package main

import (
    "fmt"
    "simplex/constdp"
    "simplex/dp"
    . "simplex/geom"
    //"simplex/constrelate"
)

func main() {
    //const_geom := constrelate.NewGeometryRelate()
    //const_dir := constrelate.NewQuadRelate()
    //const_dist := constrelate.NewMinDistanceRelate(0.3)
    //consts := []*constrelate.Constraint{const_dir, const_geom, const_dist}

    var data = []*Point{{0.5, 1.0}, {1.5, 2.0}, {2.5, 1.5}, {3.5, 2.5}, {4.0, 1.5}, {3.0, 1.0} }
    wkt0 := "LINESTRING ( 3.3972882547242254 2.1881448635223006, 3.5430863803963715 1.7339276258513834 )"
    wkt1 := "POLYGON (( 1.3897602166231338 1.526445677779483, 1.3897602166231338 1.7171047651969051, 1.557988823167918 1.7171047651969051, 1.557988823167918 1.526445677779483, 1.3897602166231338 1.526445677779483 ))"
    var g0 = ReadGeometry(wkt0)
    var g1 = ReadGeometry(wkt1)

    db := constdp.NewConstDB([]Geometry{g0, g1})

    var opts = &dp.Options{Polyline: data, Threshold: 0}
    var tree = constdp.NewConstDP(opts, true)
    var tree_str = tree.Print()
    fmt.Println(tree_str)
    opts1 :=      &dp.Options{
        Threshold   : 0.6,
        Db          : db,
    }
    if opts1.Constraints == nil {
        fmt.Println("Nil Constraints")
    }

    fmt.Println(NewLineString(data))
    //var simplx = &dp.Options{
    //    Threshold   : 0.6,
    //    Db          : db,
    //    Constraints : consts,
    //}
    //tree.Simplify(simplx)
}
