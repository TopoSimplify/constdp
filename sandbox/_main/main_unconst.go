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
var dp = require("../dp")
var trajstream = require("../trajstream")
var constrelate = require("../constrelate")


async.waterfall([
  func (cb) {//traj stream
    var args = {}
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
    var tree = dp(options)
    cb(nil, tree)
  }
],
  func (err, tree) { //simplify
    if err {
      errors.New(err)
    }

    client_simplify(tree)
  })


func client_simplify(tree) {

  var res  = 10000     //res
  var genpoly = tree.simplify(tree.root, res)

  if len(genpoly) {
    var genpoly_geom = geom.LineString(genpoly)
    console.log(genpoly_geom.toString())
  }
  else {
    console.log(nil)
  }
}
