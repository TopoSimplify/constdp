package main

import "github.com/intdxdt/geom"

type Ln struct {
	Id      string
	Geom    *geom.LineString
	Partnum int
}

func NewLn(id string, geom *geom.LineString, partnum int) *Ln {
	return &Ln{Id: id, Geom: geom, Partnum: partnum}
}
