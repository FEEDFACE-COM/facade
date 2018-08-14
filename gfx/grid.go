
package gfx

import (
    "fmt"
//    "flag"
)

type Grid struct {
    width uint
    height uint
}


func NewGrid(width, height uint) *Grid {
    ret := &Grid{width:width, height:height}
    return ret    
}


//func (grid *Grid) AddFlags(flags *flag.FlagSet) {
//    flags.UintVar(&grid.width, "width", grid.width, "pager width")   
//    flags.UintVar(&grid.heigt, "height", grid.height, "pager height")   
//}
//

func (grid *Grid) Render() {
    
    // stuff involving opengl
    
    
}

func (grid *Grid) Desc() string { return fmt.Sprintf("grid[%dx%d]",grid.width,grid.height) }

