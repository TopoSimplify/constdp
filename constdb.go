package constdp

import (
    "io"
    "log"
    "bufio"
    "strings"
    "simplex/geom"
    "simplex/geom/mbr"
    "simplex/struct/rtree"
)


//in-memory rtree
func NewConstDB() *rtree.RTree {
    return rtree.NewRTree(8)
}

func LoadConstDBFromGeometries(db *rtree.RTree, geoms []geom.Geometry) *rtree.RTree {
    var gs = make([]rtree.BoxObj, 0)
    for _, g := range geoms {
        gs = append(gs, g)
    }
    return db.Load(gs)
}

func LoadConstDBFromFile(db *rtree.RTree, r io.Reader) *rtree.RTree {
    objs := make([]rtree.BoxObj, 0)
    var err = process_file_byline(r, func(wkt string) {
        objs = append(objs, geom.ReadGeometry(wkt))
    })
    if err != nil {
        log.Fatal(err)
    }

    return db.Load(objs)
}

func SearchDb(db *rtree.RTree, query *mbr.MBR) []geom.Geometry {
    nodes := db.Search(query)
    geoms := make([]geom.Geometry, len(nodes))
    for i := range nodes {
        geoms[i] = nodes[i].GetItem().(geom.Geometry)
    }
    return geoms
}

func process_file_byline(r io.Reader, process func(string)) error {
    // Open an input file, exit on error.
    scanner := bufio.NewScanner(r)
    // scanner.Scan() advances to the next token returning false if an error was encountered
    for scanner.Scan() {
        ln := strings.TrimSpace(scanner.Text())
        if ln != "" {
            process(ln)
        }
    }
    // check scanner for errors
    return scanner.Err()
}