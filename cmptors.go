package constdp
//
//
///*
// description check if sub geom is valid
// param subgeom
// param comparators
// return {boolean}
// private
// */
//func _isvalid(subgeom, comparators) {
//
//  //make true , proof otherwise
//  var bool = true
//  for (var i = 0 bool && i < len(comparators) ++i) {
//    bool = bool && comparators[i](subgeom)
//  }
//  return bool
//}
//
//
///*
// description gen cmp functors
// param polygeom
// param constlist
// param options
// returns {Array}
// private
// */
//func _cmptors(polygeom, constlist, options) {
//  var fn, comparators = []
//  var keys = Object.keys(options.const)
//
//  for (var i = 0 i < len(keys) ++i) {
//    fn = options.const[keys[i]]
//    if _.is_function(fn) {
//      comparators.append(
//        fn(polygeom, constlist, options) //return cmptor
//      )
//    }
//  }
//  return comparators
//}

