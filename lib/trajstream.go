package constdp
import "github.com/intdxdt/simplex/util"
/*
 module -
 description -
 * author: titus
 * date: 29/09/14
 * time: 7:09 AM
 */
var fs = require("fs")
var _ = require("ldsh")
var splitter = require("splitter")
var geom = require("geom")
/*
 description - module exports
 type {Function}
 */
module.exports = trajstream
/*
 description traj stream
 param filepath
 returns {*}
 */
func trajstream(filepath) {

  var stream = _.dplx(
    {objectMode : true},
    func (obj, enc, cb) {
      if obj {
        var g = geom.readwkt(obj)
        var o = geom.geom2array(g)
        self.append(o)
      }
      cb()
    }
  )
  var fstream = fs.createReadStream(filepath)
  var _mapper = func (obj) {
    return obj
  }
  var wktstream = splitter("\n", _mapper)
  fstream.pipe(wktstream).pipe(stream)
  return stream
}
