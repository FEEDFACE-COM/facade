
package conf

import (
    "fmt"
    "flag"
)

type CameraConfig struct {
    Isometric bool
}



func (config *CameraConfig) AddFlags(flags *flag.FlagSet) {
    flags.BoolVar(&config.Isometric,"iso",config.Isometric,"isometric projection" )
}

func NewCameraConfig() *CameraConfig { return &CameraConfig{Isometric: true} }

func (config *CameraConfig) Desc() string { 
    tmp := "ppv"
    if config.Isometric { tmp = "iso" }
    return fmt.Sprintf("camera[%s]",tmp) 
}
    
