package facade

import (
	"flag"
	"fmt"
)

const DEFAULT_FONT = "RobotoMono"

var FontDefaults FontConfig = FontConfig{
	Name: DEFAULT_FONT,
}

func (config *FontConfig) Desc() string {
	ret := "font["
	ret += config.GetName()
	ret += "]"
	return ret
}

func (config *FontConfig) AddFlags(flagset *flag.FlagSet) {
	flagset.StringVar(&config.Name, "font", FontDefaults.Name, "font face")
}

func (config *FontConfig) VisitFlags(flagset *flag.FlagSet) bool {
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "font":
			{
				config.SetName = true
			}
		}
	})
	return config.SetName
}

func (config *FontConfig) Help() string {
	ret := ""
	fun := func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		}
		ret += fmt.Sprintf("  -%-24s %-24s\n", name, f.Usage)
	}

	tmp := flag.NewFlagSet("font", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(fun)
	return ret
}
