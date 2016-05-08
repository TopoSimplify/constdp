package constdp
import "github.com/intdxdt/simplex/util"
/*
 * User: titus
 * Date: 16/07/13
 * Time: 3:24 PM
 */

var _     = require("ldsh")
var fs = require("fs")
var splitter = require("splitter")
var geom  = require('geom')
var async = require("async")
var path  = require("path")

var constdp     = require("../lib")
var constrelate = require("constrelate")


async.waterfall([
    func (cb) {//stream

      var constlist = []
      var filepath = path.join(__dirname, "_data/consts2.txt")

      func _mapper(wkt) {
        wkt && constlist.append(geom.readwkt(wkt))
      }
      var fstream = fs.createReadStream(filepath)
      var wktstream = splitter("\n", _mapper)
      fstream.pipe(wktstream)
      fstream.on('end', func () {
        cb(nil, constlist)
      })
    },
    func (constlist, cb) {//rtree
      //constlist = []
      var rtree = constdp.constdb(constlist)
      cb(nil, rtree)
    },
    func (rtree, cb) {//constdp
      var data = 'LINESTRING ( 520.3891360357894 542.3912033070129, 506.3024618690045 551.4232473315985, 499.8456492240652 555.3948968460392, 492.961552805167 552.5004635914114, 489.3155900796462 547.0315195031302, 494.7910190818659 540.6453203655232, 503.2430235819369 542.0539877822016, 506.3024618690045 551.4232473315985, 505.72509579166825 560.3502151427206, 505.2252456091915 568.0786679640912, 509.1982380315184 573.1744625927278, 510.22036282943196 578.5066706448671, 506.9538876603224 582.6378010057998, 500.5170101211947 582.253509809434, 492.2547493993293 573.7991034893856, 486.1060902574759 569.7640459275444, 481.13065522694603 565.0658697854723, 477.55561113833613 565.056478772063, 477.07661476306924 570.8172868053748, 478.90063032561653 576.393069064855, 488.0275462393051 588.113950554013, 494.11875860924937 590.8324257774667, 503.20704849575554 590.8039889285739, 513.3907651994501 588.2100233531045, 519.7227675184544 584.0205966373184, 528.1782382221293 579.2651445083028, 528.2675859252843 570.3930616942043, 532.7013978168333 560.3489116165815, 531.7413381584337 549.6831287580544 )'
      //weak
      //var data = 'LINESTRING ( 520.3891360357894 542.3912033070129, 506.3024618690045 551.4232473315985, 499.8456492240652 555.3948968460392, 492.961552805167 552.5004635914114, 489.3155900796462 547.0315195031302, 494.7910190818659 540.6453203655232, 503.2430235819369 542.0539877822016, 506.3024618690045 551.4232473315985, 506.7617420621395 568.2268811420811, 511.46930921762095 577.7380882521355, 502.9188300984811 587.2492953621899, 492.7129644375179 573.3818864739396, 476.7516232431989 561.9217555683406, 477.30048564327086 576.944986410244, 486.99543198158756 589.5401303711341, 503.30312129484696 593.8783184995004, 528.1782382221293 579.2651445083028, 528.2675859252843 570.3930616942043, 535.0623346410276 560.118075831373, 531.7413381584337 549.6831287580544 )'
      //data = geom.readwkt(data)
      //data = geom.geom2array(data)
      data = [[520.3891360357894,542.3912033070129],[506.3024618690045,551.4232473315985],[499.8456492240652,555.3948968460392],[492.961552805167,552.5004635914114],[489.3155900796462,547.0315195031302]]
      var options = {
        pln    : data,
        res    : 0,
        filter : nil
      }
      var tree = new constdp.ConstDP(options)
      cb(nil, tree, rtree)
    }
  ],
  func (err, tree, rtree) { //simplify
    if err) { errors.New(err }
    client_simplify(tree, rtree)
  })


func client_simplify(tree, rtree) {


  console.log(tree.print())
  var _const =  {
    mindist : 0
  }
  /*
    const comparators : signature
      func (geom, constlist, options)
   */
  _.assign(_const, constrelate)
  //delete _const.dir
  var simplx = {
    res     : 20,     //res
    db      : rtree,  //const index
    filter  : nil,   //node filter
    const   : _const  //const comparators
  }

  tree.simplify(simplx)
  var genpoly = _.at(tree.pln , tree.simple.at)
  if len(genpoly) {
    var poly_geom     = geom.LineString(tree.pln)
    var genpoly_geom  = geom.LineString(genpoly)
    console.log(poly_geom.toString())
    console.log(genpoly_geom.toString())
  }
  else {
   errors.New("genpoly is empty")
  }

}
