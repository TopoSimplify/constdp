package main

import (
    "fmt"
    "simplex/geom"
    "simplex/struct/rtree"
)

func main() {
    wkt := "LINESTRING ( 325 293, 479 425, 661 349, 559 182, 453 233, 349 141, 234 143 )"
    var g = geom.NewLineStringFromWKT(wkt)
    fmt.Println(g.Coordinates())

    var ply = geom.NewPolygon(g.Coordinates())
    fmt.Println(g.BBox())
    fmt.Println(ply.BBox())
    fmt.Println(g.Coordinates()[0].BBox())

    db := rtree.NewRTree(16)
    db.Insert(g)
    db.Insert(ply)
    db.Insert(g.Coordinates()[0])

    var glist = make([]geom.Geometry, 0)
    glist = append(glist, g)
    glist = append(glist, ply)
    glist = append(glist, g.Coordinates()[0])

    fmt.Println(glist)

}
