
package conf

import (
    "fmt"
    "flag"
)

type FontConfig struct {
    Face string
}



func (config *FontConfig) AddFlags(flags *flag.FlagSet) {
    flags.StringVar(&config.Face,"font",config.Face,"use fontface `font`" )
}

func NewFontConfig() *FontConfig { return &FontConfig{Face: "Monaco"} }

func (config *FontConfig) Describe() string { return fmt.Sprintf("font[%s]",config.Face) }


