package constdp

//Update hull nodes with dp instance
func (self *ConstDP) selfUpdate() {
	for i := range self.Hulls {
		self.Hulls[i].Instance = self
	}
}
