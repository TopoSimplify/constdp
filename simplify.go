package constdp

import (
    "simplex/dp"
    "simplex/geom"
    "simplex/util/iter"
    "fmt"
    "simplex/struct/slist"
    "golang.org/x/mobile/bind/testdata"
    "simplex/struct/stack"
)

//constrained dp simplify
func (self *ConstDP) Simplify(opts *Options) *ConstDP {
    self.opts = opts
    self.Simple.Reset()
    self.Filter(self.Root, self.opts.Threshold)
    for !(self.NodeSet.IsEmpty()) {
        self.genints()
    }
    return self
}

/*
 description generalize ints
 param dpfilter
 param options
 private
 */
func (self *ConstDP) genints() {
    var node = self.NodeSet.Shift().(*Node)
    var fixint = node.Ints.Last().(*Vertex).Index()
    //var nextint int

    fmt.Println(fixint)

    //early exit
    if node == nil {
        return
    }

    var subrange = []int{node.Key[0], node.Key[1]}

    var poly = self.subpoly(
        iter.NewGenerator(subrange[0], subrange[1] + 1),
    )

    var subpoly = self.subpoly(
        iter.NewGenerator_AsVals(subrange...),
    )

    var polygeom    =   geom.NewLineString(poly)
    var subgeom     =   geom.NewLineString(subpoly)
    var db          =   self.opts.Db

    if node.Hull == nil {
        node.Hull = polygeom.ConvexHull()
    }
    var hull        =   node.Hull
    var constlist []geom.Geometry

    if db != nil {
        constlist = SearchDb(db, hull.BBox())
    }

    fmt.Println(polygeom)
    fmt.Println(subgeom)
    fmt.Println(constlist)


    //add intersect points with neighbours as constraints
    //self.updateconsts(constlist, polygeom, node, options)

    if self.opts.Constraints != nil && len(constlist) > 0  {
      var comparators = self._cmptors(polygeom, constlist)
      //intlist
      var intlist = self._intcandidates(node)
      var curints = (intlist.shift())()
      //assume not valid
      var isvalid = false
      //proof otherwise
      for  !isvalid {
        if len(subpoly) == len(poly) {
          dpfilter.intset.appendall(subrange)
          isvalid = true
          continue
        }
        //check if subgeom is valid
        isvalid = self._isvalid(subgeom, comparators)

        if isvalid {
          dpfilter.intset.appendall(subrange)
        } else {
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
          } else {
            //reset
            if len(intlist) {
              intindex = -1
              subrange = node[node._key].slice(0)
              nextint = self.int.index(node.int, intindex)
              subrange.append(nextint) //keep top level node int
              curints = (intlist.shift())()
            }  else {
              //go to original
              subrange = polyrange
              subpoly = poly
            }
          }
        }
      }
    }  else {
      //keep range interesting index
      dpfilter.intset.appendall(subrange)
    }
}


func (self *ConstDP) _childints(node *dp.Node) {

  var stack = stack.NewStack()
  var nextint = node.Ints.Last().(*Vertex).Index()
  var intobj = {}
  var intlist = []
  //node stack
    if node.
  stack.append(node.right)
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
 */
func (self *ConstDP) _intcandidates (node *Node) {
   var node_ints  =     func () { return node.Ints }
   var child_ints =     func () { return self._childints(node) }
  return []*slist.SList{node_ints, child_ints}
}