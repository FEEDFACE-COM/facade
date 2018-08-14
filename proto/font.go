
package proto

import (
    "flag"
)

type Font struct {
    Face string    
}
func (font *Font) Desc() string { 
    ret := "font[" + font.Face + ""
    return ret
}

func NewFont() *Font { return &Font{Face: "Monaco"} }

func (font *Font) AddFlags(flags *flag.FlagSet) {
    flags.StringVar(&font.Face, "font",font.Face,"use font `font`" )
}
