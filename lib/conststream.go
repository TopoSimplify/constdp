package constdp
import "github.com/intdxdt/simplex/util"
/*
 module -
 description -
 * author: titus
 * date: 06/09/14
 * time: 10:43 PM
 */
var _ = require("ldsh")
var fs = require("fs")
var splitter = require("splitter")
var geom = require("geom")
/*
 description - module exports
 type {Function}
 */
module.exports = conststream
/*
 description constraint stream from a file with wkts
 param filepath
 returns {*}
 */
func conststream(filepath) {

  var stream = _.dplx(
    {objectMode : true},
    func (obj, enc, cb) {
      if obj {
        var g     = geom.readwkt(obj)
        self.append(g)
      }
      cb()
    }
  )
  var fstream = fs.createReadStream(filepath)
  var wktstream = splitter("\n", _.identity)
  fstream.pipe(wktstream).pipe(stream)
  return stream
}
//test
//------
//var boxstream = _.dplx(
//  {objectMode : true},
//  func (obj, enc, cb) {
//    if obj {
//      self.append(JSON.stringify(obj.bbox)+"\n")
//    }
//    cb()
//  }
//)
//var data = path.join(__dirname, "./_data/consts.txt")
//conststream(data)
//  .pipe(boxstream)
//  .pipe(process.stdout)
