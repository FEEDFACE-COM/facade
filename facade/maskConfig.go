package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
)

var MaskDefaults MaskConfig = MaskConfig{
	Name: "mask",
}

func (config *MaskConfig) Desc() string {
	ret := "mask["
	if config.GetSetName() {
		ret += config.GetName()
	}
	ret += "]"
	return ret
}

var MASK_QUICKFIX bool = MaskDefaults.Name == "mask"
func (config *MaskConfig) AddFlags(flagset *flag.FlagSet, basicOptions bool) {
	masks := " (" + AvailableShaders("mask/", ".frag") + ")"
	if basicOptions {
		flagset.BoolVar(&MASK_QUICKFIX, "mask", MASK_QUICKFIX, "overlay mask?")
	} else {
		flagset.StringVar(&config.Name, "mask", MaskDefaults.Name, "overlay mask"+masks)
	}
}

func (config *MaskConfig) VisitFlags(flagset *flag.FlagSet, basicOptions bool) bool {
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "mask":
			{
				config.SetName = true
				if basicOptions {
					if MASK_QUICKFIX {
						config.Name = "mask"
					} else {
						config.Name = "def"
					}
				}
			}
		}
	})
	return config.SetName
}

func (config *MaskConfig) Help(basicOptions bool) string {
	ret := ""
	tmp := flag.NewFlagSet("mask", flag.ExitOnError)
	config.AddFlags(tmp,basicOptions)
	tmp.VisitAll(func(f *flag.Flag) { ret += gfx.FlagHelp(f) })
	return ret
}
