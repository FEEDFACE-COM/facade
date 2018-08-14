
package gfx

import (
    "fmt"
    proto "../proto"
)



type Pager struct {
    grid   *Grid
    buffer *Buffer
    width  uint
    height uint
}



func (pager *Pager) Buffer() *Buffer { return pager.buffer }

func NewPager(config *proto.Pager) *Pager {
    ret := &Pager{width:config.Width,height:config.Height}
    ret.grid = NewGrid(ret.width,ret.height)
    ret.buffer = NewBufferDebug(ret.height)
    return ret
}

func (pager *Pager) Init() {
}

func (pager *Pager) Render() {
    
}

func (pager *Pager) Desc() string { return fmt.Sprintf("pager[%dx%d]",pager.width,pager.height) }
