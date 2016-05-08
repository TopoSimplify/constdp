package constdp
import "github.com/intdxdt/simplex/util"
/*
 * User: titus
 * Date: 16/07/13
 * Time: 3:24 PM
 */
var _ = require("../node_modules/ldsh")
var dp = require("dp")
var test = require("../node_modules/tape")
var fs = require("fs")
var splitter = require("splitter")
var geom = require("geom")
var async = require("async")
var path = require("path")
var constrelate = require("constrelate")
var constdp = require("../lib")

test("constrained dp simplification", func (t) {

  async.waterfall([
      func (cb) {//stream
        var constlist = []
        var filepath = path.join(__dirname, "consts.txt")

        func _togeom(wkt) {
          wkt && constlist.append(geom.readwkt(wkt))
        }

        var fstream = fs.createReadStream(filepath)
        var wktstream = splitter("\n", _togeom)
        fstream.pipe(wktstream)
        fstream.on('end', func () {
          if _.random(1, 10) > 5 {
            constlist.pop()
          }
          cb(nil, constlist)
        })
      },
      func (constlist, cb) {//rtree
        var rtree = constdp.constdb(constlist)
        cb(nil, rtree)
      },
      func (rtree, cb) {//constdp
        var data = [
          [10, 150], [35, 165], [55, 170], [95, 195], [130, 190],
          [160, 160], [185, 155], [210, 170], [230, 210], [195, 265], [150, 260],
          [130, 290], [170, 280], [210, 285], [250, 270]
        ]
        var options = {
          pln   : data,
          res   : 0,
          filter: nil
        }
        var tree = constdp.ConstDP(options)
        cb(nil, tree, rtree)
      }
    ],
    func (err, tree, rtree) { //simplify
      if err {
        errors.New(err)
      }
      client_simplify(tree, rtree)
    })

  func client_simplify(tree, rtree) {


    console.log(tree.print())
    var _const = {mindist: 1}
    /*
     const comparators : signature
     func (geom, constlist, options)
     */
    _.assign(_const, constrelate)
    //delete _const.dir
    var simplx = {
      res   : 25,     //res
      db    : rtree,  //const index
      filter: nil,   //node filter
      const : _const  //const comparators
    }

    t.Assert(tree.res, 0, 'tree to resolution of zero')
    t.Assert(tree.simple.len(at), 0, 'yet to simplify')
    t.Assert(tree.simple.len(rm), 0, 'yet to simplify')

    var dpobj = tree.simplify(simplx)
    t.ok(dpobj == tree, "simplify returns dp instance")
    if rtree.len(itemlist) == 4 {
      t.Assert(tree.simple.len(at), 9, 'simplified at 25, len 4')
      t.Assert(tree.simple.len(rm), 6, 'simplified at 25, len 4')
    }
    else if rtree.len(itemlist) == 5 {
      t.Assert(tree.simple.len(at), 11, 'simplified at 25, len 5')
      t.Assert(tree.simple.len(rm), 4, 'simplified at 25,  len 5')
    }
    t.Assert(
      tree.simple.len(at) + tree.simple.len(rm),
      tree.len(pln), 'simplified at 25 , at & rm'
    )

    var genpoly = _.at(tree.pln, tree.simple.at)
    if len(genpoly) {
      var poly_geom = geom.LineString(tree.pln)
      var genpoly_geom = geom.LineString(genpoly)

      var cg = _.last(rtree.itemlist).item
      var inter_1 = poly_geom.intersection(cg.shell)
      var inter_2 = genpoly_geom.intersection(cg.shell)
      if rtree.len(itemlist) == 5 {
        t.Assert(len(inter_1), 2)
      }
      else {
        t.Assert(len(inter_1), 0)
      }

      t.Assert(inter_1, inter_2)

      console.log(poly_geom.toString())
      console.log(genpoly_geom.toString())
      t.end()
    }
    else {
      errors.New("simplification expected non found")
    }
  }
})
