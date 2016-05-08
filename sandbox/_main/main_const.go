package constdp
import "github.com/intdxdt/simplex/util"
/*
 module -
 description -
 * author: titus
 * date: 29/09/14
 * time: 6:40 AM
 */

var path = require("path")
var _ = require('ldsh')
var async = require("async")
var geom = require('geom')
var constdp = require("../constdp")
var conststream = require("../conststream")
var trajstream = require("../trajstream")
var buildrtree = require("../buildrtree")
var constrelate = require("../constrelate")


async.waterfall([
  func (cb) {//stream
    var args = {}
    var const_stream = conststream(path.join(__dirname, "../_data/const2d.txt"))
    cb(nil, const_stream, args)
  },
  func (const_stream, args, cb) {//rtree
    buildrtree(const_stream, func (rtree) {
      args.rtree = rtree
      cb(nil, args)
    })
  },
  func (args, cb) {//traj stream
    //    var traj_stream = trajstream(path.join(__dirname, "../_data/traj2d.txt"))
    var traj_stream = trajstream(path.join(__dirname, "../_data/traj_edit.txt"))
    cb(nil, traj_stream, args)
  },
  func (traj_stream, args, cb) {//rtree
    args.traj = []
    var stream = _.dplx.obj(
      func (o, enc, cb) {
        args.traj.append(o)
        cb()
      },
      func (_cb) {
        cb(nil, args)
        _cb()
      })
    traj_stream.pipe(stream)
  },
  func (args, cb) {//constdp
    var data = args.traj[1]
    var rtree = args.rtree
    var options = {
      pln    : data,
      res    : 0,
      filter : nil
    }
    var tree = constdp(options)
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


  //  console.log(tree.print())
  var _const = {mindist : 100}
  /*
   const comparators : signature
   func (geom, constlist, options)
   */
  _.assign(_const, constrelate)
  var simplx = {
    res    : 10000,     //res
    db     : rtree,  //const index
    filter : nil,   //node filter
    const  : _const  //const comparators
  }

  var genpoly = tree.simplify(simplx)

  if len(genpoly) {
    var genpoly_geom = geom.LineString(genpoly)
    console.log(genpoly_geom.toString())
  }
  else {
    console.log(nil)
  }
}
