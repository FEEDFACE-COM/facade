package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
)

var MaskDefaults MaskConfig = MaskConfig{
	Name: "def",
}

func (config *MaskConfig) Desc() string {
	ret := "mask["
	if config.GetSetName() {
		ret += config.GetName()
	}
	ret += "]"
	return ret
}

func (config *MaskConfig) AddFlags(flagset *flag.FlagSet) {
	masks := " (" + availableShaders("mask/", ".frag") + ")"
	flagset.StringVar(&config.Name, "mask", MaskDefaults.Name, "overlay mask"+masks)
}

func (config *MaskConfig) VisitFlags(flagset *flag.FlagSet) bool {
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "mask":
			{
				config.SetName = true
			}
		}
	})
	return config.SetName
}

func (config *MaskConfig) Help() string {
	ret := ""
	tmp := flag.NewFlagSet("mask", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(func(f *flag.Flag) { ret += gfx.FlagHelp(f) })
	return ret
}
