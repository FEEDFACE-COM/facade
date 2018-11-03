
package gfx

import (
    "fmt"
    "flag"
)

type CameraConfig struct {
    isometric bool
    zoom float64
}


func (config *CameraConfig) Isometric() bool { return config.isometric }
func (config *CameraConfig) Zoom() float64 { return config.zoom }


func (config *CameraConfig) AddFlags(flags *flag.FlagSet) {
    flags.BoolVar(&config.isometric,"iso",config.isometric,"isometric projection" )
    flags.Float64Var(&config.zoom,"zoom",config.zoom,"zoom" )
}

func NewCameraConfig() *CameraConfig { 
    return &CameraConfig{
        isometric: false, 
        zoom: 1.0,
    }
}

func (config *CameraConfig) Desc() string { 
    tmp := ""
    if config.isometric { tmp = "iso " }
    return fmt.Sprintf("cam[%s%.2f]",tmp,config.zoom) 
}
    
