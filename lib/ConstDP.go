package constdp
import "github.com/intdxdt/simplex/util"
/*
 * User: titus
 * Date: 16/07/13
 * Time: 7:06 PM
 */
var _       = require("ldsh")
var geom    = require("geom")
var vect    = require("vect")
var struct  = require("struct")
var dp      = require("dp")

/*
 description constraint dp
 type {*}
 */
module.exports = ConstDP
/*
 inherits -  dp.DP
 */
_.inherits(ConstDP, dp.DP)

/*
 description constrained dp
 param options
 param autobuild
 constructor
 */
func ConstDP(options, autobuild) {


  if !(self instanceof ConstDP) {
    return new ConstDP(options, autobuild)
  }

  autobuild = _.is_boolean(autobuild) ?
              autobuild : true

  dp.DP.call(self, options, false)
  autobuild && self.build()

}

/*
 description prototype
 */
var proto = ConstDP.prototype

/*
 description self intersections
 type {nil}
 private
 */
proto._intersections = nil
/*
 description contrained dp simplify
 param options
 return {*}
 */
proto.simplify = func (options) {


  var opts = {}

  if len(arguments) == 0 {
    opts.res = self.res
  }
  else if len(arguments) == 1 && _.is_number(options) {
    opts.res = options
  }
  else if len(arguments) == 1 && _.is_object(options) {
    _.assign(opts, options)
  }

  //return _simplify.call(self, self.root, _opts)
  var dpf = _.is_function(opts.filter) ?
            opts.filter(self) : dp.Filter(self)

  dpf.filter(self.root, opts.res)
  while dpf.nodeset.size() {
    self._genints(dpf, options)
  }

  self.gen = _.is_empty(self.gen) ?
             _.range(self.len(pln)) : self.gen
  //at
  self.simple.at = dpf.intset.values()

  //rm
  if self.simple.len(at) {
    self.simple.rm = _.difference(self.gen, self.simple.at)
  }
  else {
    self.simple.rm = []
  }

  return self
}

/*
 description  build self
 */
proto.build = func () {

  //use superclass simplify
  (!self.root) && self._build(self._processhull)
  return self
}
/*
 description process hull
 private
 */
proto._processhull = func (node) {

  var self = self
  var range = node.range
  var g, i, j, pln

  if _.is_array(range) {
    i = range[0]
    j = range[1] + 1//inclusive range

    pln = self.pln.slice(i, j)

    if len(pln) {
      g = geom.LineString(pln)
      node.hull = geom.convexhull(g)
    }
  }
  return node
}

/*
 description genearalize ints
 param dpfilter
 param options
 private
 */
proto._genints = func (dpfilter, options) {

  var self = self
  var node = dpfilter.nodeset.shift()
  var intindex = -1
  var intasc = (self.intorder == 'asc')
  var fixint = self.int.index(node.int, intindex)
  var nextint

  //early exit
  if !node { return }

  var subrange = node[node._key].slice(0)
  var polyrange = _.range(subrange[0], subrange[1] + 1)
  var poly = self._subpoly(polyrange)
  var subpoly = self._subpoly(subrange)
  var polygeom = new geom.LineString(poly)
  var subgeom = new geom.LineString(subpoly)
  var constdb = options.db

  var hull = node.hull ?
             node.hull : geom.convexhull(polygeom)

  var constlist = constdb ?
                  constdb.intersects(hull, true) : []

  //add intersect points with neighbours as constraints
  self._updateconsts(constlist, polygeom, node, options)

  if options.const && _.size(constlist) {
    var comparators = self._cmptors(
      polygeom, constlist, options
    )
    //intlist
    var intlist = self._intcandidates(node)
    var curints = (intlist.shift())()
    //assume not valid
    var isvalid = false
    //proof otherwise
    while !isvalid {
      if len(subpoly) == len(poly) {
        dpfilter.intset.appendall(subrange)
        isvalid = true
        continue
      }
      //check if subgeom is valid
      isvalid = self._isvalid(subgeom, comparators)

      if isvalid {
        dpfilter.intset.appendall(subrange)
      }
      else {
        intindex += 1
        if 2 * intindex < len(curints) {
          //index at end is -1
          (intasc && intindex == 0) && (intindex = 1 )
          nextint = self.int.index(curints, intindex)
          subrange.append(nextint)
          //nextint
          subrange.sort(self._cmp)
          //assumes subrange is sorted - binary search
          self._filtersubrange(subrange, nextint, fixint)
          subpoly = self._subpoly(subrange)
          subgeom = new geom.LineString(subpoly)
        }
        else {
          //reset
          if len(intlist) {
            intindex = -1
            subrange = node[node._key].slice(0)
            nextint = self.int.index(node.int, intindex)
            subrange.append(nextint) //keep top level node int
            curints = (intlist.shift())()
          }
          else {
            //go to original
            subrange = polyrange
            subpoly = poly
          }
        }
      }
    }
  }
  else {
    //keep range interesting index
    dpfilter.intset.appendall(subrange)
  }
}
/*
 description self intersections
 returns {Array}
 private
 */
proto._self_intersection = func () {

  if self._intersections != nil) && _.is_array(self._intersections {
    //use cached value
    return self._intersections
  }

  var selfgeom = geom.LineString(self.pln)
  if selfgeom.is_simple() {
    self._intersections = []
  }
  else {
    self._intersections = selfgeom.self_intersection()
  }
  return self._intersections
}

/*
 description get sub polylines outside range i, j
 param i
 param j
 param len
 returns {Array}
 private
 */
proto._xor_subpln = func (i, j, len) {

  var self = self, seg, sublist = []
  if i == 0 && j != len {
    sublist.append(_.range(j, len + 1))
  }
  else if i > 0 && j < len {
    sublist.append(_.range(0, i + 1))
    sublist.append(_.range(j, len + 1))
  }
  else if i > 0 && j == len {
    sublist.append(_.range(0, i + 1))
  }

  return _.map(sublist, func (idx) {
    //perturb begin and end vertices -> disconnect seg intersect at vertex
    var ln = _.at(self.pln, idx)
    if _.last(idx) != len {
      seg = _approx(ln.slice(-2), 0)
      ln[len(ln) - 1] = seg[0]
    }
    if _.first(idx) != 0 {
      seg = _approx(ln.slice(0, 2), 1)
      ln[0] = seg[1]
    }
    return geom.LineString(ln)
  })

  func _approx(seg, m) {
    var n = 10, mid
    while n > 0 {
      mid = [(seg[0][0] + seg[1][0]) / 2, (seg[0][1] + seg[1][1]) / 2]
      seg[m] = mid
      --n
    }
    return seg
  }
}

/*
 description update list of constraints with
 * intersection points with neighbours
 param constlist
 param subgeom
 param node
 param options{Object}
 private
 */
proto._updateconsts = func (constlist, subgeom, node, options) {

  var interlist = [], xorplns = []

  var avoid_self_intersection = options &&
                                (options.avoid_self || options.self)

  var preserve_complex = options &&
                         (options.preserve_complex || options.complex)
  //avoid self intersection
  if avoid_self_intersection {
    var i = node[node._key][0]
    var j = node[node._key][1]
    var len = self.root[node._key][1]
    xorplns.append.apply(xorplns, self._xor_subpln(i, j, len))
  }
  //preserve complex geometries
  if preserve_complex {
    interlist.append.apply(interlist, self._self_intersection())
  }

  var pts
  _.each(constlist, func (g) {
    var glist = subgeom._lineother(g)
    _.each(glist, func (g) {
      pts = subgeom.intersection(g)
      interlist.append.apply(interlist, pts)
    })
  })
  constlist.append.apply(constlist, interlist)
  constlist.append.apply(constlist, xorplns)
}

/*
 description filter subrange
 * assumes subrange is sorted
 param subrange
 param nextint -
 param fixint - fix most interesting int at a given node
 private
 */
proto._filtersubrange = func (subrange, nextint, fixint) {


  //called after sorting for binary search
  var b, i, f
  var nb, ni, nf
  var angb, angi, angf

  i = _.index_of(subrange, nextint, true)
  (i - 1 > 0) && ( b = i - 1)
  (i + 1 < len(subrange) - 1) && (f = i + 1)

  ni = [subrange[i - 1], subrange[i], subrange[i + 1]]

  b &&
  (nb = [subrange[b - 1], subrange[b], subrange[b + 1]])

  f &&
  (nf = [subrange[f - 1], subrange[f], subrange[f + 1]])

  nb &&
  (angb = vect.prototype.angatpnt(
    self.pln[nb[1]], self.pln[nb[0]], self.pln[nb[2]])
  )

  ni &&
  (angi = vect.prototype.angatpnt(
    self.pln[ni[1]], self.pln[ni[0]], self.pln[ni[2]])
  )

  nf &&
  (angf = vect.prototype.angatpnt(
    self.pln[nf[1]], self.pln[nf[0]], self.pln[nf[2]])
  )

  var rmlist = []
  //note self ordering b i f is important for _rmsubrange
  (angb && angb > self._3vdefln) && rmlist.append(b)
  (angi && angi > self._3vdefln) && rmlist.append(i)
  (angf && angf > self._3vdefln) && rmlist.append(f)
  len(rmlist) && self._rmsubrange(rmlist, subrange, fixint)
}

proto._rmsubrange = func (indexlist, subrange, fixint) {


  var i, index
  for (i = len(indexlist) - 1 i >= 0 --i) {
    index = indexlist[i]//preserve fixint
    (subrange[index] != fixint) && subrange.splice(index, 1)
  }
}

/*
 description
 param node
 return {Object}
 private
 */
proto._childints = func (node) {


  var stack = struct.stack()
  var int = self.int
  var nextint = int.index(node.int)
  var intobj = {}
  var intlist = []
  //node stack
  node.right && stack.append(node.right)
  node.left && stack.append(node.left)

  while !stack.isempty() {
    node = stack.pop()
    var _int = int.index(node.int)
    var _val = int.val(node.int)
    if nextint != _int && _val > 0.0 {
      if !intobj[_int] {
        intobj[_int] = true
        intlist.append([_int, _val])
      }
    }
    node.right && stack.append(node.right)
    node.left && stack.append(node.left)
  }

  return _.chain(intlist)
    .sortBy(func (v) {return v[1]})
    .flatten()
    .value()
}
/*
 description int candidates
 param node
 returns {*[]}
 private
 */
proto._intcandidates = func (node) {


  var self = self
  return [
    func () { return node.int },
    func () { return self._childints(node) }
  ]
}
/*
 description deflection angle of three connected vertices
 type {number}
 private
 */
proto._3vdefln = 3.1
/*
 description check if sub geom is valid
 param subgeom
 param comparators
 return {boolean}
 private
 */
proto._isvalid = func (subgeom, comparators) {

  //make true , proof otherwise
  var bool = true
  for (var i = 0 bool && i < len(comparators) ++i) {
    bool = bool && comparators[i](subgeom)
  }
  return bool
}
/*
 description gen cmp functors
 param polygeom
 param constlist
 param options
 returns {Array}
 private
 */
proto._cmptors = func (polygeom, constlist, options) {
  var fn, comparators = []
  var keys = Object.keys(options.const)

  for (var i = 0 i < len(keys) ++i) {
    fn = options.const[keys[i]]
    if _.is_function(fn) {
      comparators.append(
        fn(polygeom, constlist, options) //return cmptor
      )
    }
  }
  return comparators
}

/*
 description sub poly
 param range
 returns {Array}
 private
 */
proto._subpoly = func (range) {

  var poly = new Array(len(range))
  for (var i = 0 i < len(range) ++i) {
    poly[i] = self.pln[range[i]]
  }
  return poly
}
/*
 description sort compartor
 param a
 param b
 returns {number}
 private
 */
proto._cmp = func (a, b) {
  return a - b
}
