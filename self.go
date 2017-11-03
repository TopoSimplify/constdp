package constdp

//Update hull nodes with dp instance
func (self *ConstDP) selfUpdate() {
	for _, hull := range self.Hulls {
		hull.Instance = self
	}
}
