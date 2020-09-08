package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"strings"
)

const DEFAULT_FONT = "monaco"

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
	fonts := " (" + availableFonts() + ")"
	flagset.StringVar(&config.Name, "font", FontDefaults.Name, "typeface"+fonts)
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
	tmp := flag.NewFlagSet("font", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(func(f *flag.Flag) { ret += gfx.FlagHelp(f) })
	return ret
}

func availableFonts() string {
	ret := ""
	for name, _ := range FontAsset {
		ret += name
		ret += ", "
	}
	ret = strings.TrimSuffix(ret, ", ")
	return ret
}
