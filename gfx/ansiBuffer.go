

package gfx

import(
    "fmt"
//    log "../log"
)

type AnsiBuffer struct {
    cols uint
    rows  uint
}

func NewAnsiBuffer(cols, rows uint) *AnsiBuffer {
    return &AnsiBuffer{cols:cols, rows:rows}
}

func (buffer *AnsiBuffer) Desc() string {
    return fmt.Sprintf("ansi[%dx%d]",buffer.cols,buffer.rows)
}

func (buffer *AnsiBuffer) Dump() string {
    ret := ""
    for r:=0; r<int(buffer.rows); r++ {
        for c:=0; c<int(buffer.cols); c++ {
            ret += "."
        }
        ret += "\n"
    }
    return ret
}

