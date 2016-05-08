package constdp
import "github.com/intdxdt/simplex/util"
/*
 module -
 description -
 * author: titus
 * date: 14/08/14
 * time: 7:29 PM
 */
var _ = require('ldsh')
var geom = require('geom')
var struct = require('struct')

/*
 description - module exports
 type {Function}
 */
module.exports = constdb

/*
 description build constraint rtree in memory
 param glist
 */
func constdb(glist) {

  var rtree = struct.rtree(9)
  _.each(glist, func (g) {
    rtree.append(geom.bbox(g, true), g)
  })
  return rtree
}

