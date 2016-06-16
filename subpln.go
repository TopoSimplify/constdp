package constdp

import (
    . "simplex/geom"
    . "simplex/util/iter"
    "simplex/vect"
)

//get sub polylines outside range i, j
func (self *ConstDP) xor_subpln(i, j, N int) []Geometry {
    var pt *Point
    var coords []*Point
    var idx  *Generator
    var geometries = make([]Geometry, 0)

    if i > 0 {
        idx := NewGenerator(0, i + 1)
        pt = self.Perturb(self.Pln[i - 1], self.Pln[i])
        coords = self.subpoly(idx)
        coords[len(coords) - 1] = pt
        geometries = append(geometries, NewLineString(coords))
    }

    if j < N {
        idx = NewGenerator(j, N + 1)
        pt = self.Perturb(self.Pln[j + 1], self.Pln[j])
        coords = self.subpoly(idx)
        coords[0] = pt
        geometries = append(geometries, NewLineString(coords))
    }
    return geometries
}


//Perturb boundary at point vertex b
//returns a point close to b but not b along segment a---b
func (self *ConstDP) Perturb(a , b *Point) *Point {
    eps := 1e-8
    a, b = a.Clone(), b.Clone()
    v := vect.NewVect(&vect.Options{A:a, B:b})

    M := v.M()
    M = M - eps

    D := v.D()

    v2 := vect.NewVect(&vect.Options{A:a, M:&M, D: &D})
    pt := v2.B()
    return &pt
}

