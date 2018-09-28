
package gfx

import (
    "fmt"
    "flag"
)

type FontConfig struct {
    Name string
}



func (config *FontConfig) AddFlags(flags *flag.FlagSet) {
    flags.StringVar(&config.Name,"font",config.Name,"use fontface `font`" )
}

func NewFontConfig() *FontConfig { return &FontConfig{Name: "Monaco"} }

func (config *FontConfig) Desc() string { return fmt.Sprintf("font[%s]",config.Name) }


