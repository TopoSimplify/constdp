package main

import "fmt"

type  Generator struct {
    start, stop, step, v int
    Next                 bool
    from_vals            bool
    vals                 []int
}

func NewGenerator(args ...int) *Generator {
    self := &Generator{}
    if len(args) == 1 {
        self.start, self.stop, self.step = 0, args[0], 1
    } else if len(args) == 2 {
        self.start, self.stop, self.step = args[0], args[1], 1
    } else if len(args) == 3 {
        self.start, self.stop, self.step = args[0], args[1], args[2]
    }
    self.v = self.start - self.step
    self.update_next(self.start)
    return self
}

func NewGenerator_AsVals(args ...int) *Generator {
    self := &Generator{}
    self.vals = args
    if len(self.vals) > 0 {
        self.start = 0
        self.stop = len(self.vals)
        self.step = 1
    }
    self.from_vals = true
    self.v = self.start - self.step
    self.update_next(self.start)
    return self
}

func (self *Generator) update_next(v int) {
    if self.step > 0 {
        self.Next = v < self.stop
    } else {
        self.Next = v > self.stop
    }
}

func (self *Generator) Val() int {
    self.v += self.step
    if (self.step > 0  && self.v > self.stop) ||
        (self.step < 0  && self.v < self.stop) {
        panic("generator out of range")
    }
    self.update_next(self.v + self.step)

    if self.from_vals {
        return self.vals[self.v]
    }
    return self.v
}

func intGen(i, j int) func() (bool, int) {
    var n = i - 1;
    return func() (bool, int) {
        n += 1
        if n < j {
            return true, n
        } else {
            return false, 0
        }
    }
}

func main() {
    rng := [...]int{4, 6, 8, 9, 6, 7}
    var gen = NewGenerator_AsVals(rng[:]...)
    for gen.Next {
        fmt.Println(gen.Val());
    }
}
