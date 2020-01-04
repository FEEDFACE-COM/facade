package facade

import (
	"flag"
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
