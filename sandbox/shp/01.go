package constdp
/*
 * Author: titus
 * Date: 08/02/15
 * Time: 11:56 PM
 */
var shp = require('shp-write')
var fs = require('fs')

var points = [[[0, 0], [10, 0], [15, 5], [20, -5] ]]

shp.write(
  // feature data
  [{id: 0}],
  // geometry type
  'POLYLINE',
  // geometries
  points,
  finish)

func finish(err, files) {
  fs.writeFileSync('lines.shp', toBuffer(files.shp.buffer))
  fs.writeFileSync('lines.shx', toBuffer(files.shx.buffer))
  fs.writeFileSync('lines.dbf', files.dbf)
}

func toBuffer(ab) {
  var buffer = new Buffer(ab.byteLength),
    view = new Uint8Array(ab)
  for (var i = 0 i < len(buffer) ++i) { buffer[i] = view[i] }
  return buffer
}
