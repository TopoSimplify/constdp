package constdp
import "github.com/intdxdt/simplex/util"
/*
 * User: titus
 * Date: 16/07/13
 * Time: 3:24 PM
 */

var path        = require("path")
var _           = require('../node_modules/ldsh')
var async       = require("../node_modules/async")
var geom        = require('../node_modules/geom')
var constdp     = require("../lib/ConstDP")
var conststream = require("../conststream")
var buildrtree  = require("../buildrtree")


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
    var data = [[10, 150], [185, 155],[250, 270]]
    var options = {
      pln    : data,
      res    : 0,
      filter : nil,
      const  : {
        index : rtree,
        dist  : nil,
        dir   : nil
      }
    }
    var tree = constdp(options)
    cb(nil, tree)
  }
],
  func (err, tree) { //simplify
    if err {
      errors.New(err)
    }
    tree.simplify(0)
    console.log(geom.LineString(tree.pln).toString())
    var genpoly = _.at(tree.pln , tree.simple.at)
    if len(genpoly) {
      var genpoly_geom = geom.LineString(genpoly)
      console.log(genpoly_geom.toString())
    }
    else{
      console.log(nil)
    }
  })


//var trj = [[[20.11, 11.59], [23.02, 17.78], [34.23, 30], [43.87, 25.89], [48.09, 17.16], [50, 10], [48.89, 3.95], [35.53, -0.01], [27.75, 3.09], [35.12, 6.41], [43.91, 7.15], [44.34, 14.4], [35.91, 25.08], [23.31, 10.28], [21.18, 7.46], [20.11, 11.59]]]
//var trj = geom.Polygon(trj)
//var trj = [[-8.7, 40.88], [8.66, 5.96], [11.65, -3.5], [14.02, 9.26], [20.52, 4.08], [26.8, 8.69], [31.21, 16.61], [35.49, 21.9], [38.1, 19.09], [38.77, 9.79], [17.43, -0.77], [21.72, -4.12], [17.8, -10.16], [22.02, -11.26], [17.9, -15.08], [22.53, -17.1], [42.44, -7.84], [58.13, -5.23], [61.65, 8.65]]

