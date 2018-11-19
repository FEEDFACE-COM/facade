
package gfx

import (
    "flag"
)

type MaskConfig map[string]interface{}

func (config *MaskConfig) Mask() (bool,bool) { ret,ok := (*config)["mask"].(bool); return ret,ok }
func (config *MaskConfig) SetMask(val bool) { (*config)["mask"] = val }
    

func NewMaskConfig() *MaskConfig {
    ret := make(MaskConfig)
    return &ret
}


func (config *MaskConfig) AddFlags(flags *flag.FlagSet) {
//    flags.BoolVar(&config.Mask,"mask",config.Mask,"mask" )  
}




func (config *MaskConfig) Desc() string { 
    ret := "mask["
    if val,ok := config.Mask(); ok && val {
        ret += "✓"
    }
    ret += "]"
    return ret
}

