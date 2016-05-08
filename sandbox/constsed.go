package constdp
import "github.com/intdxdt/simplex/util"
/*
 * User: titus
 * Date: 16/07/13
 * Time: 3:24 PM
 */

var path = require("path")
var _ = require('ldsh')
var async = require("async")
var geom = require('geom')
var constsed = require("../constsed/constsed")
var conststream = require("../conststream")
var buildrtree = require("../buildrtree")
var constrelate = require("../constrelate")


async.waterfall([
    func (cb) {//stream
      var stream = conststream(path.join(__dirname, "../_data/consts.txt"))
      cb(nil, stream)
    },
    func (stream, cb) {//rtree
      buildrtree(stream, func (rtree) {
        cb(nil, rtree)
      })
    },
    func (rtree, cb) {//constdp
      var data = [
        [10.0, 150.0, 6.5], [35.0, 165.0, 6.8], [55.0, 170.0, 7.0],
        [95.0, 195.0, 7.5], [130.0, 190.0, 7.9], [160.0, 160.0, 8.1],
        [185.0, 155.0, 8.6], [210.0, 170.0, 9.5], [230.0, 210.0, 10.1],
        [195.0, 265.0, 10.3], [150.0, 260.0, 10.5], [130.0, 290.0, 10.8],
        [170.0, 280.0, 11.8], [210.0, 285.0, 11.9], [250.0, 270.0, 12.5],
        [280.0, 280.0, 12.8]
      ]
      var options = {
        pln    : data,
        res    : 0,
        filter : nil
      }
      var tree = constsed(options)
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
  var _const = {mindist : 1}
  /*
   const comparators : signature
   func (geom, constlist, options)
   */
  _.assign(_const, constrelate)
  var simplx = {
    res    : 40,     //res
    db     : rtree,  //const index
    filter : nil,   //node filter
    const  : _const  //const comparators
  }

  tree.simplify(simplx)
  var genpoly = _.at(tree.pln , tree.simple.at)
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
