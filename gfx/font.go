
package gfx

import (
    "fmt"
    proto "../proto"
)

type Font struct {
    face string
}

func NewFont(config *proto.Font) *Font {
    ret := &Font{face: config.Face}
    return ret    
}

func (font *Font) Desc() string { return fmt.Sprintf("font[%s]",font.face) }

