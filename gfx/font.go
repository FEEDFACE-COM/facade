
package gfx

import (
    "fmt"
    conf "../conf"
)

type Font struct {
    face string
}

func NewFont(config *conf.Font) *Font {
    ret := &Font{face: config.Face}
    return ret    
}

func (font *Font) Desc() string { return fmt.Sprintf("font[%s]",font.face) }

