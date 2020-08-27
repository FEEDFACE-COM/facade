package facade

import (
	"flag"
	"fmt"
)

const ENABLE_CAMERA_ISOMETRIC = false

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
	if ENABLE_CAMERA_ISOMETRIC {
    	flagset.BoolVar(&config.Isometric, "iso", CameraDefaults.Isometric, "camera isometric?")
    }
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

func (config *CameraConfig) Help() string {
	ret := ""
	fun := func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		}
		ret += fmt.Sprintf("  -%-24s %-24s\n", name, f.Usage)
	}

	tmp := flag.NewFlagSet("camera", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(fun)
	return ret
}
