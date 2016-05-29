package constdp

import (
    "simplex/struct/rtree"
    "simplex/geom/mbr"
    "simplex/geom"
    "bufio"
    "strings"
    "log"
    "os"
)


//in-memory rtree
func NewConstDB() *rtree.RTree {
    return rtree.NewRTree(16)
}

func LoadConstDBFromGeometries(db *rtree.RTree, objs []rtree.BoxObj) *rtree.RTree {
    return db.Load(objs)
}

func LoadConstDBFromFile(db *rtree.RTree, fname string) *rtree.RTree {
    objs := make([]rtree.BoxObj, 0)
    process_file_byline(fname, func(wkt string) {
        objs = append(objs, geom.ReadGeometry(wkt))
    })
    return db.Load(objs)
}

func SearchDb(db *rtree.RTree, query *mbr.MBR) []geom.Geometry {
    nodes := db.Search(query)
    geoms := make([]geom.Geometry, len(nodes)) //
    for i := range nodes {
        geoms[i] = nodes[i].GetItem().(geom.Geometry)
    }
    return geoms
}

func process_file_byline(fname string, process func(string)) {
    // Open an input file, exit on error.
    fid, err := os.Open(fname)
    if err != nil {
        log.Fatal("Error opening input file:", err)
    }
    defer fid.Close()

    scanner := bufio.NewScanner(fid)
    // scanner.Scan() advances to the next token returning false if an error was encountered
    for scanner.Scan() {
        ln := strings.TrimSpace(scanner.Text())
        if ln != "" {
            process(ln)
        }
    }
    // check scanner for errors
    if err := scanner.Err(); err != nil {
        log.Fatal(scanner.Err())
    }
}