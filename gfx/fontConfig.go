
package gfx

import (
    "fmt"
    "flag"
)

type FontConfig struct {
    name string
}

func (config *FontConfig) Name() string { return config.name }

func (config *FontConfig) AddFlags(flags *flag.FlagSet) {
    flags.StringVar(&config.name,"font",config.name,"render with `font`" )
}

func NewFontConfig() *FontConfig { return &FontConfig{name: "RobotoMono"} }

func (config *FontConfig) Desc() string { return fmt.Sprintf("font[%s]",config.name) }


