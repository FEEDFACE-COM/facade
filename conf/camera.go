
package conf

import (
    "fmt"
    "flag"
)

type CameraConfig struct {
    Isometric bool
    Zoom float64
}



func (config *CameraConfig) AddFlags(flags *flag.FlagSet) {
    flags.BoolVar(&config.Isometric,"iso",config.Isometric,"isometric projection" )
    flags.Float64Var(&config.Zoom,"zoom",config.Zoom,"zoom" )
}

func NewCameraConfig() *CameraConfig { return &CameraConfig{Isometric: true, Zoom: 1.0} }

func (config *CameraConfig) Desc() string { 
    tmp := "ppv"
    if config.Isometric { tmp = "iso" }
    return fmt.Sprintf("camera[%s %.2f]",tmp,config.Zoom) 
}
    
