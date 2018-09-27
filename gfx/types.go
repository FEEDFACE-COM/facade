
// +build linux,arm

package gfx

import (
    "fmt"
)


func (c *Coord) Desc() string { return fmt.Sprintf("%+d/%+d",c.X,c.Y) }
type Coord struct {
    X int
    Y int    
}


func (d *Dim) Desc() string { return fmt.Sprintf("%dx%d",d.W,d.H) }
type Dim struct {
    W int
    H int
}

func (s *Size) Desc() string { return fmt.Sprintf("%5.1fx%5.1f",s.W,s.H) }
type Size struct {
    W float32
    H float32  
}

func (p *Point) Desc() string { return fmt.Sprintf("(%5.1f %5.1f)",p.X,p.Y) }
type Point struct {
    X float32
    Y float32    
}
    