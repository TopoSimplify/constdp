package ctx

import (
	"simplex/geom"
	"simplex/geom/mbr"
	"simplex/struct/sset"
	"simplex/constdp/cmp"
)

const (
	Self             = "self"
	SelfVertex       = "self_vertex"
	SelfSegment      = "self_segment"
	SelfSimple       = "self_simple"
	SelfNonVertex    = "self_non_vertex"
	ContextNeighbour = "context_neighbour"
)

type ctxMeta struct {
	SelfVertices    *sset.SSet
	SelfNonVertices *sset.SSet
}

type CtxGeom struct {
	Geom geom.Geometry
	Type string
	I    int
	J    int
	Meta *ctxMeta
}

func NewCtxGeom(g geom.Geometry, i, j int) *CtxGeom {
	return &CtxGeom{
		Geom: g,
		Type: Self,
		I:    i,
		J:    j,
		Meta: &ctxMeta{
			SelfVertices:    sset.NewSSet(cmp.IntCmp, 2),
			SelfNonVertices: sset.NewSSet(cmp.IntCmp, 2),
		},
	}
}

func (o *CtxGeom) String() string {
	return o.Geom.WKT()
}

func (o *CtxGeom) BBox() *mbr.MBR {
	return o.Geom.BBox()
}

// --------------------------------------------------------------------
func (o *CtxGeom) AsSelf() *CtxGeom {
	o.Type = Self
	return o
}

func (o *CtxGeom) IsSelf() bool {
	return o.Type == Self
}

// --------------------------------------------------------------------
func (o *CtxGeom) AsSelfVertex() *CtxGeom {
	o.Type = SelfVertex
	return o
}

func (o *CtxGeom) IsSelfVertex() bool {
	return o.Type == SelfVertex
}

// --------------------------------------------------------------------
func (o *CtxGeom) AsSelfNonVertex() *CtxGeom {
	o.Type = SelfNonVertex
	return o
}

func (o *CtxGeom) IsSelfNonVertex() bool {
	return o.Type == SelfNonVertex
}

// --------------------------------------------------------------------
func (o *CtxGeom) AsSelfSegment() *CtxGeom {
	o.Type = SelfSegment
	return o
}

func (o *CtxGeom) IsSelfSegment() bool {
	return o.Type == SelfSegment
}

// --------------------------------------------------------------------
func (o *CtxGeom) AsSelfSimple() *CtxGeom {
	o.Type = SelfSimple
	return o
}

func (o *CtxGeom) IsSelfSimple() bool {
	return o.Type == SelfSimple
}

// --------------------------------------------------------------------
func (o *CtxGeom) AsContextNeighbour() *CtxGeom {
	o.Type = ContextNeighbour
	return o
}

func (o *CtxGeom) IsContextNeighbour() bool {
	return o.Type == ContextNeighbour
}

// --------------------------------------------------------------------
func (o *CtxGeom) Intersection(other geom.Geometry) []*geom.Point {
	return o.Geom.Intersection(other)
}
