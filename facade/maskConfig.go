package facade

import (
	"flag"
	gfx "../gfx"
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
	flagset.StringVar(&config.Name, "mask", MaskDefaults.Name, "mask shader")
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
	tmp.VisitAll( func (f *flag.Flag) { ret += gfx.FlagHelp(f) } )
	return ret
}
