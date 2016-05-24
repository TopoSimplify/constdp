package constdp

//
///*
// description filter subrange
// * assumes subrange is sorted
// param subrange
// param nextint -
// param fixint - fix most interesting int at a given node
// private
// */
//func (self *ConstDP) _filtersubrange(subrange, nextint, fixint) {
//  //called after sorting for binary search
//  var b, i, f
//  var nb, ni, nf
//  var angb, angi, angf
//
//  i = _.index_of(subrange, nextint, true)
//  (i - 1 > 0) && ( b = i - 1)
//  (i + 1 < len(subrange) - 1) && (f = i + 1)
//
//  ni = [subrange[i - 1], subrange[i], subrange[i + 1]]
//
//  b &&
//  (nb = [subrange[b - 1], subrange[b], subrange[b + 1]])
//
//  f &&
//  (nf = [subrange[f - 1], subrange[f], subrange[f + 1]])
//
//  nb &&
//  (angb = vect.prototype.angatpnt(
//    self.pln[nb[1]], self.pln[nb[0]], self.pln[nb[2]])
//  )
//
//  ni &&
//  (angi = vect.prototype.angatpnt(
//    self.pln[ni[1]], self.pln[ni[0]], self.pln[ni[2]])
//  )
//
//  nf &&
//  (angf = vect.prototype.angatpnt(
//    self.pln[nf[1]], self.pln[nf[0]], self.pln[nf[2]])
//  )
//
//  var rmlist = []
//  //note self ordering b i f is important for _rmsubrange
//  if (angb && angb > self._3vdefln) { rmlist.append(b) }
//  if (angi && angi > self._3vdefln) { rmlist.append(i) }
//  if (angf && angf > self._3vdefln) { rmlist.append(f) }
//  len(rmlist) && self._rmsubrange(rmlist, subrange, fixint)
//}
//
//func (self *ConstDP) _rmsubrange(indexlist, subrange, fixint) {
//  var i, index
//  for (i = len(indexlist) - 1 i >= 0 --i) {
//    index = indexlist[i]//preserve fixint
//    (subrange[index] != fixint) && subrange.splice(index, 1)
//  }
//}
//
