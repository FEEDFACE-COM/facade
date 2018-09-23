
package conf

import (
    "fmt"
    "flag"
)

type MaskConfig struct { 
    Mask bool
}

func (config *MaskConfig) AddFlags(flags *flag.FlagSet) {
    flags.BoolVar(&config.Mask,"mask",config.Mask,"mask" )
    
}

func NewMaskConfig() *MaskConfig { 
    return &MaskConfig{Mask: false} 
}

func (config *MaskConfig) Desc() string { return fmt.Sprintf("mask[]") }
