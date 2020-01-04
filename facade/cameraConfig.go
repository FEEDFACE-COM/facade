package facade

import (
	"flag"
	"fmt"
)

var CameraDefaults = CameraConfig{
	Zoom:      1.0,
	Isometric: false,
}

func (config *CameraConfig) Desc() string {
	ret := "cam["
	if config.GetSetZoom() {
		ret += fmt.Sprintf("%.1f", config.GetZoom())
	}
	if config.GetSetIsometric() {
		if config.GetIsometric() {
			ret += "i"
		} else {
			ret += "p"
		}
	}
	ret += "]"
	return ret
}

func (config *CameraConfig) AddFlags(flagset *flag.FlagSet) {
	flagset.Float64Var(&config.Zoom, "zoom", CameraDefaults.Zoom, "camera zoom")
	flagset.BoolVar(&config.Isometric, "iso", CameraDefaults.Isometric, "camera iso?")
}

func (config *CameraConfig) VisitFlags(flagset *flag.FlagSet) bool {
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "zoom":
			{
				config.SetZoom = true
			}
		case "iso":
			{
				config.SetIsometric = true
			}
		}
	})
	return config.SetZoom || config.SetIsometric
}
