package main

import "fmt"

type Areal interface {
    Area()
}

func ProcessArea(fn func()) {
    fmt.Println("Areal :")
    fn()
}

type Rectangle struct {
    w int
    h int
}

func (self *Rectangle ) Area() {
    fmt.Println("Rectangular Area : ", self.w * self.h)
}

type Square struct {
    w int
}

func (self *Square ) Area() {
    fmt.Println("Square Area : ", self.w * self.w)
}

func comStaticArea() {
    fmt.Println("static :", 7 * 8)
}

func main() {
    sq := &Square{4}
    rect := &Rectangle{4, 7}
    ProcessArea(sq.Area)
    ProcessArea(rect.Area)
    ProcessArea(comStaticArea)
}
