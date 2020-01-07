package facade

import (
	"flag"
	"fmt"
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
	fun := func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		}
		ret += fmt.Sprintf("  -%-24s %-24s\n", name, f.Usage)
	}

	tmp := flag.NewFlagSet("mask", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(fun)
	return ret
}
