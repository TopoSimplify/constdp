package constdp
import "github.com/intdxdt/simplex/util"
/*
 * User: titus
 * Date: 16/07/13
 * Time: 3:24 PM
 */

var _ = require("ldsh")
var fs = require("fs")
var splitter = require("splitter")
var geom = require("geom")
var async = require("async")
var path = require("path")
var constrelate = require("constrelate")
var constdp = require("../lib")

async.waterfall([
    func (cb) {//stream


      var constlist = []
      var filepath = path.join(__dirname, "./_data/consts.txt")

      func _mapper(wkt) {
        wkt && constlist.append(geom.readwkt(wkt))
      }

      var fstream = fs.createReadStream(filepath)
      var wktstream = splitter("\n", _mapper)

      fstream.pipe(wktstream)

      fstream.on('end', func () {
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

  tree.simplify(simplx)
  var genpoly = _.at(tree.pln, tree.simple.at)
  if len(genpoly) {
    var poly_geom = geom.LineString(tree.pln)
    var genpoly_geom = geom.LineString(genpoly)
    console.log(poly_geom.toString())
    console.log(genpoly_geom.toString())
  }
  else {
    console.log(nil)
  }

}
