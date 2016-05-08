package constdp
import "github.com/intdxdt/simplex/util"
///*
// * User: titus
// * Date: 16/07/13
// * Time: 3:24 PM
// */
//var ConstDouglasPeucker = require("../../lib/core").ConstDouglasPeucker
//var DouglasPeucker = require("../../lib/core").DouglasPeucker
//var _ = require('./../../lib/base')._
//var libgge = require('./../../../libgge')
//
//var poly_bay = [[[20.11, 11.59], [23.02, 17.78], [34.23, 30], [43.87, 25.89], [48.09, 17.16], [50, 10], [48.89, 3.95], [35.53, -0.01], [27.75, 3.09], [35.12, 6.41], [43.91, 7.15], [44.34, 14.4], [35.91, 25.08], [23.31, 10.28], [21.18, 7.46], [20.11, 11.59]]]
//var polygon_bay_geom = libgge.core.geom.Polygon(poly_bay)
//var ship_path = [[-8.7, 40.88], [8.66, 5.96], [11.65, -3.5], [14.02, 9.26], [20.52, 4.08], [26.8, 8.69], [31.21, 16.61], [35.49, 21.9], [38.1, 19.09], [38.77, 9.79], [17.43, -0.77], [21.72, -4.12], [17.8, -10.16], [22.02, -11.26], [17.9, -15.08], [22.53, -17.1], [42.44, -7.84], [58.13, -5.23], [61.65, 8.65]]
//var ship_path_geom = libgge.core.geom.LineString(ship_path)
////var constGeomList = _.map([poly_12, poly_13, poly_14, poly_15], func (wkt) {
////    return libgge.core.geom.readGeometryFromWkt(wkt)
////})
//constGeomList = [polygon_bay_geom]
//
//var res = 0, distThresh = 1
//var dp_obj = new ConstDouglasPeucker(ship_path,res)
//
////dp_obj.simplify()
//dp_obj.simplify()
//
////simplification will be triggered on traverse
//var gen_poly = dp_obj.constTraverse(30,distThresh,constGeomList)
//var gen_poly2 = dp_obj.traverseAt(dp_obj.tree, 30)
//var cur_gen_line_geom   = libgge.core.geom.createGeomFromLineArray(gen_poly)
//var cur_gen_line_geom2   = libgge.core.geom.createGeomFromLineArray(gen_poly2)
//
//console.log(JSON.stringify(ship_path))
//console.log(JSON.stringify(poly_bay))
//console.log(JSON.stringify(libgge.core.geom.geometryToArray(cur_gen_line_geom)))
//console.log(JSON.stringify(libgge.core.geom.geometryToArray(cur_gen_line_geom2)))
//
