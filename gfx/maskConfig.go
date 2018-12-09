
package gfx

import (
    "flag"
)

const(
	maskMask = "mask"
)

type MaskConfig Config

func (config *MaskConfig) Mask() (bool,bool) { ret,ok := (*config)[maskMask].(bool); return ret,ok }
func (config *MaskConfig) SetMask(val bool) { (*config)[maskMask] = val }
    
func (config *MaskConfig) Desc() string { 
    ret := "mask["
    if val,ok := config.Mask(); ok {
	    if  val {ret += "✓" } else { ret += "✗" }
    }
    ret += "]"
    return ret
}


func (config *MaskConfig) ApplyConfig(cfg *MaskConfig) {
	if tmp,ok := cfg.Mask(); ok { config.SetMask(tmp) }	
}


type MaskState struct {
	Mask bool	
}

var MaskDefaults = MaskState{
	Mask:    false,
}

func (state *MaskState) AddFlags(flags *flag.FlagSet) {
    flags.BoolVar(&state.Mask,"mask",state.Mask,"mask" )  
}


func (state *MaskState) CheckFlags(flags *flag.FlagSet) (*MaskConfig,bool) {
	ok := false
	ret := make(MaskConfig)
	flags.Visit( func(f *flag.Flag) {
		if f.Name == "mask" { ok = true; ret.SetMask( state.Mask ) }
	})
	return &ret,ok
}





func (state *MaskState) Desc() string { return state.Config().Desc() }

func (state *MaskState) Config() *MaskConfig {
	ret := make(MaskConfig)
	ret.SetMask( state.Mask )
	return &ret
}

func (state *MaskState) ApplyConfig(config *MaskConfig) bool {
	changed := false
	if tmp,ok := config.Mask(); ok { if state.Mask != tmp { changed = true }; state.Mask = tmp }
	return changed
}