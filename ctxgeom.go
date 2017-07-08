package constdp

import (
	"simplex/geom/mbr"
	"simplex/geom"
)

var ctx = struct {
	Self             string
	SelfVertex       string
	SelfSegment      string
	SelfSimple       string
	SelfNonVertex    string
	ContextNeighbour string
}{
	Self:             "cdp",
	SelfVertex:       "self_vertex",
	SelfSegment:      "self_segment",
	SelfSimple:       "self_simple",
	SelfNonVertex:    "self_non_vertex",
	ContextNeighbour: "context_neighbour",
}

type CtxGeom struct {
	Geom geom.Geometry
	Type string
	I    int
	J    int
}

func NewCtxGeom(g geom.Geometry, i , j int) *CtxGeom {

	return &CtxGeom{
		Geom: g,
		Type: ctx.Self,
		I:    i,
		J:    j,
	}
}

func (o *CtxGeom) String() string {
	return o.Geom.WKT()
}

func (o *CtxGeom) BBox() *mbr.MBR {
	return o.Geom.BBox()
}

// -------------------------------------------
func (o *CtxGeom) AsSelf() *CtxGeom {
	o.Type = ctx.Self
	return o
}

func (o *CtxGeom) IsSelf() bool {
	return o.Type == ctx.Self
}

// -------------------------------------------
func (o *CtxGeom) AsSelfVertex() *CtxGeom {
	o.Type = ctx.SelfVertex
	return o
}

func (o *CtxGeom) IsSelfVertex() bool {
	return o.Type == ctx.SelfVertex
}

// -------------------------------------------
func (o *CtxGeom) AsSelfNonVertex() *CtxGeom {
	o.Type = ctx.SelfNonVertex
	return o
}

func (o *CtxGeom) IsSelfNonVertex() bool {
	return o.Type == ctx.SelfNonVertex
}

// -------------------------------------------
func (o *CtxGeom) AsSelfSegment() *CtxGeom {
	o.Type = ctx.SelfSegment
	return o
}

func (o *CtxGeom) IsSelfSegment() bool {
	return o.Type == ctx.SelfSegment
}

// -------------------------------------------
func (o *CtxGeom) AsSelfSimple() *CtxGeom {
	o.Type = ctx.SelfSimple
	return o
}

func (o *CtxGeom) IsSelfSimple() bool {
	return o.Type == ctx.SelfSimple
}

// -------------------------------------------
func (o *CtxGeom) AsContextNeighbour() *CtxGeom {
	o.Type = ctx.ContextNeighbour
	return o
}

func (o *CtxGeom) IsContextNeighbour() bool {
	return o.Type == ctx.ContextNeighbour
}

// -------------------------------------------
func (o *CtxGeom) Intersection(other geom.Geometry) []*geom.Point {
	return o.Geom.Intersection(other)
}
