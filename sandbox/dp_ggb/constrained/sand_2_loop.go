package constdp
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
//var ship_path = [[4, 22], [4.3, 15.31], [4.88, 12.13], [6.65, 11.27], [7.93, 13.62], [8, 20], [9.26, 23.87], [11.18, 21.37], [10.93, 12.84], [10.95, 10.98], [12.91, 10.51], [14.27, 13.54], [14.32, 16.7], [14, 22]]
//var const_1=    libgge.core.geom.Point([[6, 14]])
//var const_2=    libgge.core.geom.Point([[9.16342, 19.08378]])
//var const_3=    libgge.core.geom.LineString([[12.1, 13.02], [13.23, 16.1], [12, 18]])
//
//var constGeomList = [const_1,const_2,const_3]
//
//var res = 0, distThresh = 0.5
//var dp_obj = new ConstDouglasPeucker(ship_path,res)
//
////dp_obj.simplify()
//dp_obj.simplify()
//
////simplification will be triggered on traverse
//var gen_poly = dp_obj.constTraverse(11.4,distThresh,constGeomList)
//var gen_poly2 = dp_obj.traverseAt(dp_obj.tree, 11.49)
//
//var cur_gen_line_geom       = libgge.core.geom.createGeomFromLineArray(gen_poly)
//var cur_gen_line_geom2      = libgge.core.geom.createGeomFromLineArray(gen_poly2)
//
//console.log(JSON.stringify(libgge.core.geom.geometryToArray(cur_gen_line_geom)))
//console.log(JSON.stringify(libgge.core.geom.geometryToArray(cur_gen_line_geom2)))
//
