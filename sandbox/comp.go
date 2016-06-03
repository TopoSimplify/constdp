package main

import "fmt"

type Drawer struct{}

func (self *Drawer) DrawRect(rec Rect) error {
    fmt.Println("drawing rectangle")
    return nil
}

func (self *Drawer) DrawEllipseInRect(rec Rect) error {
    fmt.Println("drawing ellipse rectangle")
    return nil
}

type Point struct {
    X int
    Y int
}
type Size struct {
    Height int
    Width  int
}
type Rect struct {
    Location Point
    Size     Size
}

// VisualElement that is drawn on the screen
type VisualElement interface {
    // Draw draws the visual element
    Draw(drawer *Drawer) error
}

// Square represents a square
type Square struct {
    // Location of the square
    Location Point
    // Side size
    Side     int
}

// Draw draws a square
func (square *Square) Draw(drawer *Drawer) error {
    return drawer.DrawRect(Rect{
        Location: square.Location,
        Size: Size{
            Height: square.Side,
            Width :  square.Side,
        },
    })
}

// Circle represents a circle shape
type Circle struct {
    // Center of the circle
    Center Point
    // Radius of the circle
    Radius int
}

// Draw draws a circle
func (circle *Circle) Draw(drawer *Drawer) error {
    rect := Rect{
        Location: Point{
            X: circle.Center.X - circle.Radius,
            Y: circle.Center.Y - circle.Radius,
        },
        Size: Size{
            Width  :  2 * circle.Radius,
            Height : 2 * circle.Radius,
        },
    }

    return drawer.DrawEllipseInRect(rect)
}

// Layer contains composition of visual elements
type Layer struct {
    // Elements of visual elements
    Elements []VisualElement
}

// Draw draws a layer
func (layer *Layer) Draw(drawer *Drawer) error {
    for _, element := range layer.Elements {
        if err := element.Draw(drawer); err != nil {
            return err
        }
        fmt.Println()
    }

    return nil
}

func main() {
    circle := &Circle{
        Center: Point{X: 100, Y: 100},
        Radius: 50,
    }

    square := &Square{
        Location: Point{X: 50, Y: 50},
        Side:     20,
    }

    layer := &Layer{
        Elements: []VisualElement{circle, square},
    }
    layer.Draw(&Drawer{})
}
