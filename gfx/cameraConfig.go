
package gfx

import (
    "fmt"
    "flag"
)


// config

const (
	cameraZoom  = "zoom"
	cameraIso   = "isometric"
)

type CameraConfig Config

func (config *CameraConfig) Iso()  (bool,bool)    { ret,ok := (*config)[cameraIso].(bool);     return         ret ,ok }
func (config *CameraConfig) Zoom() (float64,bool) { ret,ok := (*config)[cameraZoom].(float64); return float64(ret),ok }


func (config *CameraConfig) SetIso(val bool)     { (*config)[cameraIso]  = val }
func (config *CameraConfig) SetZoom(val float64) { (*config)[cameraZoom] = val }

func (config *CameraConfig) Desc() string { 
    ret := "cam["
    if zoom,ok := config.Zoom(); ok {
        ret += fmt.Sprintf("%.1f",zoom)
    }
    if iso,ok := config.Iso(); ok {
	    if iso { ret += "i" } else { ret += "p" }
    } 
    ret += "]"
    return ret
}

func (config *CameraConfig) ApplyConfig(cfg *CameraConfig) {
	if tmp,ok := cfg.Zoom(); ok { config.SetZoom(tmp) }	
	if tmp,ok := cfg.Iso(); ok { config.SetIso(tmp) }	
}

// state

type CameraState struct {
	Zoom float64
	Isometric bool	
}

var CameraDefaults = CameraState{
	Zoom:        1.0,
	Isometric: false,	
}


func (state *CameraState) AddFlags(flags *flag.FlagSet) {
    flags.Float64Var(&state.Zoom,"zoom",state.Zoom,"zoom" )
    flags.BoolVar(&state.Isometric,"iso",state.Isometric,"isometric projection" )
}

func (state *CameraState) CheckFlags(flags *flag.FlagSet) (*CameraConfig,bool) {
	ok := false
	ret := make(CameraConfig)
	flags.Visit( func(f *flag.Flag) {
		if f.Name == "zoom" { ok = true; ret.SetZoom( state.Zoom ) }
		if f.Name == "iso" {  ok = true; ret.SetIso( state.Isometric ) }
	})
	return &ret,ok
}

   
func (state *CameraState) Desc() string { return state.Config().Desc() }

func (state *CameraState) Config() *CameraConfig {
	ret := make(CameraConfig)
	ret.SetIso(state.Isometric)
	ret.SetZoom(state.Zoom)
	return &ret
}

func (state *CameraState) ApplyConfig(config *CameraConfig) bool {
	changed := false
    if tmp,ok := config.Iso();  ok { if state.Isometric != tmp { changed = true }; state.Isometric = tmp }
    if tmp,ok := config.Zoom(); ok { if state.Zoom      != tmp { changed = true }; state.Zoom      = tmp }
    return changed
}


