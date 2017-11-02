package constdp

import "simplex/node"

//Update hull nodes with dp instance
func (self *ConstDP) selfUpdate() {
    var hull *node.Node
    for _, h := range *self.Hulls.DataView() {
        hull = castAsNode(h)
        hull.Instance = self
    }
}
