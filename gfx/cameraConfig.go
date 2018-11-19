
package gfx

import (
    "fmt"
    "flag"
)

type CameraConfig map[string]interface{}

//type CameraConfig struct {
//    Isometric bool
//    Zoom float64
//}


func (config *CameraConfig) Iso() (bool,bool) { ret,ok := (*config)["isometric"].(bool); return ret,ok }
func (config *CameraConfig) Zoom() (float32,bool) { ret,ok := (*config)["zoom"].(float64); return float32(ret),ok }


func (config *CameraConfig) SetIso(val bool) { (*config)["isometric"] = val }
func (config *CameraConfig) SetZoom(val float32) { (*config)["zoom"] = float64(val) }


func (config *CameraConfig) AddFlags(flags *flag.FlagSet) {
//    flags.BoolVar(&config.Isometric,"iso",config.Isometric,"isometric projection" )
//    flags.Float64Var(&config.Zoom,"zoom",config.Zoom,"zoom" )
}

func NewCameraConfig() *CameraConfig { 
    ret := make(CameraConfig)
    ret.SetZoom(1.0)
    return &ret
}

func (config *CameraConfig) Desc() string { 
    ret := "cam["
    if zoom,ok := config.Zoom(); ok {
        ret += fmt.Sprintf("%.1f",zoom)
    }
    if iso,ok := config.Iso(); ok && iso {
        ret += " iso"
    } 
    ret += "]"
    return ret
}
    
