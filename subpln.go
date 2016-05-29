package constdp
//
////get sub polylines outside range i, j
//func (self *ConstDP) xor_subpln(i, j, N int ) {
//
//  if i == 0 && j != N {
//    sublist.append(_.range(j, N + 1))
//  } else if i > 0 && j < N {
//    sublist.append(_.range(0, i + 1))
//    sublist.append(_.range(j, N + 1))
//  }  else if i > 0 && j == N {
//    sublist.append(_.range(0, i + 1))
//  }
//
//  return _.map(sublist, func (idx) {
//    //perturb begin and end vertices -> disconnect seg intersect at vertex
//    var ln = _.at(self.pln, idx)
//    if _.last(idx) != N {
//      seg = _approx(ln.slice(-2), 0)
//      ln[N(ln) - 1] = seg[0]
//    }
//    if _.first(idx) != 0 {
//      seg = _approx(ln.slice(0, 2), 1)
//      ln[0] = seg[1]
//    }
//    return geom.LineString(ln)
//  })
//
//
//}
//
//func _approx(seg, m) {
//  var n = 10, mid
//  for n > 0 {
//    mid = [(seg[0][0] + seg[1][0]) / 2, (seg[0][1] + seg[1][1]) / 2]
//    seg[m] = mid
//    n -= 1
//  }
//  return seg
//}