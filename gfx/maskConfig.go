
package gfx

import (
    "fmt"
    "flag"
)

type MaskConfig struct { 
    mask bool
}

func (config *MaskConfig) Mask() bool { return config.mask }



func (config *MaskConfig) AddFlags(flags *flag.FlagSet) {
    flags.BoolVar(&config.mask,"mask",config.mask,"mask" )
    
}

func NewMaskConfig() *MaskConfig { 
    return &MaskConfig{mask: false} 
}

func (config *MaskConfig) Desc() string { 
    ret := fmt.Sprintf("mask[]")
    if config.mask {
        ret = fmt.Sprintf("mask[✓]") 
    }
    return ret
}
